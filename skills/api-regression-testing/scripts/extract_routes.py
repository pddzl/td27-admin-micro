#!/usr/bin/env python3
"""Extract registered routes (and optional gRPC services) from go-zero sources.

Parses rest.Route blocks shaped like:

    server.AddRoute(rest.Route{
        Method:  http.MethodPost,        // or "POST"
        Path:    "/menu/user-menus",
        Handler: jwtMiddleware.Handle(menuHandler.GetUserMenus),
    })

and emits JSON: [{"method","path","auth","handler"}]. Auth is classified from
the handler wrapper: "public" (no jwt middleware), "jwt", or "jwt+record"
(record/audit middleware wrapped around jwt).

With --proto, also lists services and their rpcs from .proto files:
{"services": [{"name","rpcs":[{"name","request","response"}]}]}.

Usage:
    extract_routes.py HANDLER_FILE_OR_DIR [--proto PROTO_DIR] [--out OUT.json]

If a directory is given, every file named handler.go underneath it is parsed.
Output goes to stdout (pretty JSON) unless --out is set.
"""
import argparse
import json
import re
import sys
from pathlib import Path

METHOD_RE = re.compile(r'Method:\s+(?:http\.Method(\w+)|"(\w+)")\s*,')
PATH_RE = re.compile(r'Path:\s+"([^"]+)"\s*,')
HANDLER_RE = re.compile(r'Handler:\s+(.+?)\s*,\s*$')
BLOCK_START_RE = re.compile(r'rest\.Route\{')
SERVICE_RE = re.compile(r'^\s*service\s+(\w+)\s*\{')
RPC_RE = re.compile(r'^\s*rpc\s+(\w+)\s*\(\s*([\w.]+)\s*\)\s*returns\s*\(\s*([\w.]+)\s*\)')
MAX_WINDOW = 6  # non-comment lines to look ahead for Path/Handler after Method


def classify(handler: str) -> str:
    h = handler or ""
    has_jwt = "jwt" in h.lower()
    has_record = "record" in h.lower() or "audit" in h.lower()
    if has_jwt and has_record:
        return "jwt+record"
    if has_jwt:
        return "jwt"
    if not h:
        return "unknown"
    return "public"


def extract_file(path: Path, routes: list) -> None:
    lines = [ln for ln in path.read_text(encoding="utf-8", errors="replace").splitlines()]
    pending = {}
    window = 0
    for raw in lines:
        stripped = raw.strip()
        if not stripped or stripped.startswith("//"):
            continue
        if BLOCK_START_RE.search(raw):  # new Route block resets state
            if pending.get("method") and pending.get("path"):
                routes.append(pending)
            pending, window = {}, 0

        m = METHOD_RE.search(raw)
        p = PATH_RE.search(raw)
        h = HANDLER_RE.search(raw)
        if m and "method" not in pending:
            pending["method"] = (m.group(1) or m.group(2)).upper()
            window = 0
        elif p and "path" not in pending and pending.get("method"):
            pending["path"] = p.group(1)
            window = 0
        elif h and "handler" not in pending and pending.get("path"):
            pending["handler"] = h.group(1).strip()
            pending["auth"] = classify(pending["handler"])
            routes.append(pending)
            pending, window = {}, 0
            continue
        else:
            window += 1

        # Route with Method but no Handler captured within the window
        if pending.get("method") and pending.get("path") and window > MAX_WINDOW:
            pending.setdefault("handler", "")
            pending.setdefault("auth", classify(""))
            routes.append(pending)
            pending, window = {}, 0

        if window > MAX_WINDOW:  # stale pending, drop it
            pending, window = {}, 0

    if pending.get("method") and pending.get("path"):
        pending.setdefault("handler", "")
        pending.setdefault("auth", classify(""))
        routes.append(pending)


def collect_sources(source: Path) -> list:
    if source.is_file():
        return [source]
    return sorted(source.rglob("handler.go"))


def extract_proto(proto_dir: Path) -> list:
    services = []
    for pf in sorted(proto_dir.rglob("*.proto")):
        current = None
        for raw in pf.read_text(encoding="utf-8", errors="replace").splitlines():
            s = SERVICE_RE.match(raw)
            if s:
                current = {"name": s.group(1), "file": str(pf), "rpcs": []}
                continue
            if current is not None:
                if raw.rstrip().endswith("}") and not raw.strip().startswith("rpc"):
                    services.append(current)
                    current = None
                    continue
                r = RPC_RE.match(raw)
                if r:
                    current["rpcs"].append(
                        {"name": r.group(1), "request": r.group(2), "response": r.group(3)}
                    )
    return services


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("source", type=Path, help="route registration .go file, or a directory to search for handler.go")
    ap.add_argument("--proto", type=Path, help="optional proto root dir to also list gRPC services")
    ap.add_argument("--out", type=Path, help="write JSON here instead of stdout")
    args = ap.parse_args()

    files = collect_sources(args.source)
    if not files:
        print(f"error: no source files found under {args.source}", file=sys.stderr)
        return 1

    routes = []
    for f in files:
        extract_file(f, routes)
    routes.sort(key=lambda r: (r["path"], r["method"]))

    result = {"source_files": [str(f) for f in files], "route_count": len(routes), "routes": routes}
    if args.proto:
        result["services"] = extract_proto(args.proto)
        result["service_count"] = len(result["services"])

    payload = json.dumps(result, indent=2, ensure_ascii=False)
    if args.out:
        args.out.write_text(payload, encoding="utf-8")
        print(f"wrote {args.out}", file=sys.stderr)
    else:
        print(payload)

    auth_counts = {}
    for r in routes:
        auth_counts[r["auth"]] = auth_counts.get(r["auth"], 0) + 1
    print(f"routes: {len(routes)}  by auth: {json.dumps(auth_counts)}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
