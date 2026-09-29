package sysManagement

import (
	"net/http"

	"td27/api/gateway/internal/middleware"
	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
)

type LogoutHandler struct {
	svcCtx *svc.ServiceContext
}

func NewLogoutHandler(svcCtx *svc.ServiceContext) *LogoutHandler {
	return &LogoutHandler{svcCtx: svcCtx}
}

func (h *LogoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-token") == "" {
		api.FailWithRequest(w, http.StatusUnauthorized, "missing authorization header")
		return
	}

	middleware.RevokeRequestToken(r)

	api.OkWithMessage(w, "logout success")
}
