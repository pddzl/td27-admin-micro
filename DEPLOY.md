# Deploying td27-admin-micro to Kubernetes

Complete procedure for the target topology: a **single-node k8s cluster** with
**host-local PostgreSQL + etcd** (outside k8s), **node-local images** (no
registry), and **ArgoCD GitOps** as the only deployment path.

```
ingress(nginx) → web:8500 (nginx, serves SPA + proxies /api/*) → gateway:8888
              → gRPC basis:8080 (discovered via etcd) → host postgres / host etcd
```

---

## 0. Prerequisites (one-time per node)

- Kubernetes single node with containerd (kubeadm / k3s / RKE2) + `kubectl`.
- `podman` or `docker` on the node for building/saving images.
- [ingress-nginx](https://kubernetes.github.io/ingress-nginx/deploy/) installed
  (the Ingress uses `ingressClassName: nginx`).
- [ArgoCD](https://argo-cd.readthedocs.io/en/stable/getting_started/) installed
  into the `argocd` namespace (CLI optional but handy).
- PostgreSQL 16 and etcd 3.5 running **on the host**, reachable from pods at the
  node's LAN IP on ports 5432 and 2379 (step 1).
- A GitHub **deploy key** (read access) for
  `git@github.com:pddzl/td27-admin-micro.git`.
- A local copy of the gitignored `deploy/k8s-configs/{db.env,basis.yaml,gateway.yaml}`
  — these are never in git; copy them to the node manually (step 4).

## 1. Run the host datastores

Migrations are NOT run here — the in-cluster `td27-migrations` Job applies
`scripts/0*.sql` on every sync. The DB itself just has to exist and listen on
the host's LAN interface.

```bash
# PostgreSQL (POSTGRES_DB auto-creates td27_micro; or CREATE DATABASE manually)
sudo mkdir -p /var/lib/td27/pg /var/lib/td27/etcd
podman run -d --name td27-postgres --restart=unless-stopped -p 5432:5432 \
  -v /var/lib/td27/pg:/var/lib/postgresql/data \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=<password> -e POSTGRES_DB=td27_micro \
  docker.m.daocloud.io/library/postgres:16-alpine

# etcd — single member; client on 0.0.0.0, advertised as the host IP
HOST_IP=<node-lan-ip>
podman run -d --name td27-etcd --restart=unless-stopped -p 2379:2379 \
  -v /var/lib/td27/etcd:/etcd-data quay.io/coreos/etcd:v3.5.17 \
  etcd --data-dir=/etcd-data \
  --listen-client-urls=http://0.0.0.0:2379 \
  --advertise-client-urls=http://${HOST_IP}:2379
```

Make sure both ports are reachable from the pod network (firewalld/ufw allow
5432,2379 from the node/CIDR).

## 2. Point the manifests at your node (in git — selfHeal means git is truth)

Edit and commit+push:

- [deploy/k8s/datastores.yaml](deploy/k8s/datastores.yaml) — replace BOTH
  `203.0.113.10` placeholders (RFC-5737 documentation IP, fails fast if
  forgotten) with the node LAN IP in the two Endpoints objects.
- [deploy/k8s/ingress.yaml](deploy/k8s/ingress.yaml) — replace the
  `td27.example.com` placeholder host with your real hostname.

## 3. Build and import the node-local images (BEFORE first sync)

Images are built OUT of this repo (no Dockerfiles committed). Build three
images on the node, tagged to match `newTag` in
[deploy/k8s/kustomization.yaml](deploy/k8s/kustomization.yaml) (`v0.1.0`
initially), then import:

**Image contract**

| Image | Entrypoint | Args | Listens | Notes |
|---|---|---|---|---|
| `localhost/td27-basis` | the basis binary | `-f /etc/td27/basis.yaml` | gRPC :8080 | mounts config Secret at `/etc/td27`, writes uploads to `/app/resource/upload` |
| `localhost/td27-gateway` | the gateway binary | `-f /etc/td27/gateway.yaml` | HTTP :8888, `/health` | mounts config Secret at `/etc/td27` |
| `localhost/td27-web` | nginx | — | HTTP :8500 | serves the SPA and proxies `/api/*` to `http://$BACKEND:8888` (prefix stripped); `BACKEND=gateway` is injected as env |

Example for basis (gateway identical, web = nginx + `envsubst` proxy config):

```bash
GONOSUMCHECK='*' GOFLAGS='-mod=mod' CGO_ENABLED=0 go build -o basis-bin ./rpc/basis
podman build -t localhost/td27-basis:v0.1.0 -f - . <<'EOF'
FROM docker.m.daocloud.io/library/alpine:3.21
COPY basis-bin /usr/local/bin/basis
ENTRYPOINT ["/usr/local/bin/basis"]
EOF
```

Import into containerd (adjust `CTR` for k3s/RKE2):

```bash
make image-import TAG=v0.1.0                       # containerd (kubeadm)
make image-import TAG=v0.1.0 CTR="sudo k3s ctr"    # k3s
```

Order matters: ArgoCD starts syncing as soon as the Application is registered
(step 5). With `imagePullPolicy: Never` a missing image means
`ImagePullBackOff` — build first, or let the sync retry backoff recover it.

## 4. Create the out-of-band Secrets

Copy the gitignored config sources to the node, then:

```bash
make secret
```

Creates `td27-db-credentials` (from `db.env`), `basis-config`, and
`gateway-config` (whole-file YAML) in namespace `td27`. K8s hostnames in the
configs must be `postgres` and `etcd-client`; `Mode: prod`.

## 5. Bootstrap GitOps

```bash
# Register the repo — the URL string must EXACTLY match repoURL in
# deploy/argocd/application.yaml (ssh form, no https://)
argocd repo add git@github.com:pddzl/td27-admin-micro.git \
  --ssh-private-key-path ~/.ssh/td27_deploy_key

# Register AppProject + Application; auto-sync kicks in immediately
make app
```

ArgoCD syncs the ROOT kustomization into namespace `td27`
(`prune: false`, `selfHeal: true`). The `td27-migrations` Job is a Sync hook:
ArgoCD recreates it on every sync, so `scripts/0*.sql` re-run each release —
keep them safe to re-apply.

## 6. Verify

```bash
kubectl -n td27 get pods,job,endpoints
argocd app get td27            # Synced / Healthy
kubectl -n td27 logs job/td27-migrations   # psql applied 001–003
```

Then through the ingress (from any machine that resolves the host):

```bash
curl -H "Host: <your-ingress-host>" http://<node-ip>/health   # via web → gateway
# or add the host to /etc/hosts and browse to http://<your-ingress-host>/
```

Login is captcha-gated: `POST /captcha` then `POST /login` with
`captcha_id` + `captcha`.

## 7. Day-2 operations

| Task | Command / action |
|---|---|
| New release | build 3 images with new TAG → `make release TAG=vX.Y.Z` (imports, bumps `newTag`, commits, pushes → ArgoCD rolls pods) |
| Config change | edit `deploy/k8s-configs/*.yaml` → `make reconfig` (updates Secrets + rollout restart) |
| Restart all | `make restart` (`restartedAt` is ignored by ArgoCD — no drift) |
| Rollback | revert the tag-bump commit; ArgoCD syncs back |
| Migration change | edit `scripts/`, push — the hook recreates the Job on next sync; scripts must be re-apply-safe |
| Manual apply (no ArgoCD) | `kubectl apply -k .` from repo root (hook annotations are inert) |

## 8. Troubleshooting

| Symptom | Cause / fix |
|---|---|
| Pods `ImagePullBackOff` | image not in containerd — rebuild + `make image-import` (tag must match `newTag`) |
| basis CrashLoop | etcd unreachable — check Endpoints IP in `datastores.yaml` and host firewall |
| gateway panics at boot | `basis.rpc` missing in etcd — basis not running/registered yet |
| `td27-migrations` fails | postgres unreachable or wrong creds — check Endpoints IP + `db.env` |
| App `OutOfSync` forever | manual cluster edits revert via selfHeal — change git, not the cluster |
| Stale Endpoints after IP change | edit both `ip:` values in git; ArgoCD applies on sync |
