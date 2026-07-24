## Context

Dict create currently returns `SuccessResp{success:true}` but the frontend needs the created dict for `push(res.data)`. The dict INSERT now has `RETURNING id` but the response still returns empty data.

## Decisions

Change the gateway handler to: after creating the dict, fetch it by ID to return full dict data to the frontend. This avoids proto changes (CreateDict proto returns SuccessResp, not DictResp).

## Risks

None — the fetch-by-ID path is already implemented and tested.
