package middleware

import (
	"strings"

	"td27/rpc/basis/types/sysManagement/permission_pb"
)

// readSuffixes lists POST endpoints that are queries, not mutations. They map
// to the "read" action so read-only roles can browse data.
var readSuffixes = []string{
	"/list", "/flat", "/el-tree", "/by-page", "/get", "/check", "/validate",
	"/batch-check", "/user-buttons", "/user-menus", "/get-el-tree-depts",
}

// httpMethodToAction maps an HTTP request to a Casbin action. POST paths follow
// the gateway's verb-suffixed convention (/create, /delete, /update, ...).
func httpMethodToAction(method, path string) permission_pb.Action {
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		return permission_pb.Action_ACTION_READ
	}

	// POST: query endpoints first, then mutation verbs from the last path segment
	for _, suffix := range readSuffixes {
		if strings.HasSuffix(path, suffix) {
			return permission_pb.Action_ACTION_READ
		}
	}
	switch {
	case strings.HasSuffix(path, "/delete") || strings.HasSuffix(path, "/delete-by-ids"):
		return permission_pb.Action_ACTION_DELETE
	case strings.HasSuffix(path, "/create"):
		return permission_pb.Action_ACTION_CREATE
	case strings.HasSuffix(path, "/execute") || strings.HasSuffix(path, "/cleanup") || strings.HasSuffix(path, "/reload-policy"):
		return permission_pb.Action_ACTION_EXECUTE
	default:
		return permission_pb.Action_ACTION_UPDATE
	}
}

// roleIdsFromClaims extracts role ids from JWT claims. After JSON decoding the
// claims arrive as []interface{} of float64.
func roleIdsFromClaims(v interface{}) []int64 {
	items, ok := v.([]interface{})
	if !ok {
		return nil
	}
	ids := make([]int64, 0, len(items))
	for _, item := range items {
		if f, ok := item.(float64); ok {
			ids = append(ids, int64(f))
		}
	}
	return ids
}
