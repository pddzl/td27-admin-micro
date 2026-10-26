---
name: naming-conventions
description: Naming rules for every layer of td27-admin-micro (proto wire, DB, Go, gateway JSON, frontend TS/Vue). Use when writing, renaming, or reviewing code that crosses a layer boundary — API fields, request payloads, model tags, table columns — or when a "field missing / renders N/A / silently ignored" bug appears.
---

# Naming Conventions

The single source of truth for names that cross a boundary is the **proto
field name**. Every layer maps its local style onto it; nobody invents a
second spelling of the same field.

## The one rule

**The wire is snake_case end to end.** A field called `req_param` in proto
is `req_param` in JSON, in gateway request structs, in frontend API types,
and in Vue row props. Only Go *identifiers* follow Go casing
(`ReqParam`), and only frontend *local variables* follow TS casing
(`reqParam` as a local, never as a wire key).

## Layer rules

| Layer | Style | Example | Notes |
|---|---|---|---|
| Proto fields | snake_case | `req_param = 7` | Defines the wire; changing it = proto change → regen → all consumers |
| DB tables/columns | snake_case | `sys_monitor_operation_log.created_at` | sqlx `db:` tags must match exactly |
| Go identifiers | Go style | `ReqParam`, `listOperationLog` | Never embed casing style of other layers into identifiers |
| Model `json:` tags | ignored on the wire | `json:"reqParam"` | pb types are the wire; do NOT rely on model json tags |
| Gateway request structs | `json:` tag snake_case | `json:"page_size"` | Must mirror proto field names; camelCase here silently drops the field |
| Frontend API types | mirror wire exactly | `req_param: string` | Copy field list from the proto, not from imagination |
| Frontend locals/functions | camelCase | `const reqParam = ...` | Wire keys stay snake_case even inside camelCase code |
| Vue row/template props | snake_case | `scope.row.created_at` | camelCase here renders `undefined` → "N/A" |
| Config YAML keys | hyphenated (mapstructure) | `db-name`, `signing-key` | Only config keys use hyphens |

## Boundaries

- **DB → Go**: sqlx `db:"..."` tags must equal the column names. Soft-delete
  and timestamp columns are always `created_at` / `updated_at` /
  `deleted_at`.
- **Go → wire**: pb-generated structs serialize proto field names
  (snake_case). Gateway handlers must never re-wrap responses under new
  key styles.
- **Frontend → gateway**: request payloads use the same snake_case keys as
  the gateway's `json:"..."` tags (e.g. `page_size`, `captcha_id`).
- **API client types are contracts**: when touching an endpoint, diff the
  frontend type against the proto message first.

## URI conventions

- **Pattern**: `/<singular-resource>/<action>` — kebab-case for multi-word
  parts (`/operation-log/delete-by-ids`, `/service-token/:id`). Never
  camelCase or underscores in URIs (004 renamed `deleteByIds`, `getUserInfo`,
  `switchOpen` for exactly this reason).
- **Action suffix decides the permission action** ([authz.go](file:///Users/paul/Mine/td27-admin-micro/api/gateway/internal/middleware/authz.go)):
  `/list`, `/get`, `/flat`, `/el-tree`, `/by-page`, `/get-el-tree-depts` →
  read; `/create` → create; `/delete`, `/delete-by-ids` → delete;
  `/execute`, `/cleanup`, `/reload-policy` → execute; anything else →
  update. Pick the suffix deliberately — it is the authz contract, not
  decoration.
- **One canonical URI per resource**: no aliases, no legacy second names
  (`/opl/*` was retired by migration 011).
- **Path params**: `:id` style (`/file/download/:id`), and the
  `sys_management_api` row records the parametrized form.
- **Permission rows mirror URIs**: every registered route needs
  `sys_management_api` + api-domain `sys_management_permission` rows with
  `resource` equal to the registered path (and the action from the suffix
  rule); a mismatch 403s even for root.

## Bugs this prevents (do not repeat)

1. **BUG-001**: gateway serialized camelCase while the frontend consumed
   snake_case — two sides "fixed" the same bug against different snapshots.
2. **page_size silently dropped**: gateway `json:"pageSize"` vs frontend
   `page_size` — every explicit page size fell back to 10 with no error.
3. **operation-log page N/A + expand crash**: page read `createdAt` /
   `reqParam` from snake_case wire — `formatDateTime(undefined)` rendered
   N/A and `JSON.parse(undefined)` threw.

All three were silent: `code: 0`, no error, wrong data. A key mismatch
never fails loudly — it fails by omission.

## Review checklist

- [ ] New gateway request field: `json:` tag matches the proto field name
      exactly (grep the proto, don't recall it).
- [ ] New frontend API type: field list copied from the proto message.
- [ ] New table column: `db:` tag matches, and any SELECT lists the column
      explicitly (no `SELECT *`).
- [ ] Vue template reads `row.<snake_case_key>` — search the component for
      camelCase row props (`row\.[a-z]+[A-Z]` is a red flag).
- [ ] Renaming a field: proto first, then regen, then gateway, then
      frontend — never edit one layer's spelling alone.
