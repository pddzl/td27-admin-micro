package sysManagement

import (
	"net/http"

	"td27/api/gateway/internal/middleware"
	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/role_pb"
)

// ElevationHandler JIT role elevation endpoints. Requester/approver/revoker
// identities always come from the JWT, never from the request body.
type ElevationHandler struct {
	svcCtx *svc.ServiceContext
}

func NewElevationHandler(svcCtx *svc.ServiceContext) *ElevationHandler {
	return &ElevationHandler{svcCtx: svcCtx}
}

func userIdFromCtx(r *http.Request) int64 {
	userId, _ := r.Context().Value(middleware.UserIdKey).(float64)
	return int64(userId)
}

func (h *ElevationHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req role_pb.CreateElevationRequestReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	userId := userIdFromCtx(r)
	if userId == 0 {
		api.FailWithRequest(w, http.StatusUnauthorized, "missing user identity")
		return
	}
	req.UserId = userId

	if _, err := h.svcCtx.RoleClient.CreateElevationRequest(r.Context(), &req); err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithMessage(w, "elevation request submitted")
}

func (h *ElevationHandler) Decide(w http.ResponseWriter, r *http.Request) {
	var req role_pb.DecideElevationReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	userId := userIdFromCtx(r)
	if userId == 0 {
		api.FailWithRequest(w, http.StatusUnauthorized, "missing user identity")
		return
	}
	req.DecidedBy = userId

	if _, err := h.svcCtx.RoleClient.DecideElevation(r.Context(), &req); err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithMessage(w, "elevation request decided")
}

func (h *ElevationHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	var req role_pb.RevokeElevationReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	userId := userIdFromCtx(r)
	if userId == 0 {
		api.FailWithRequest(w, http.StatusUnauthorized, "missing user identity")
		return
	}
	req.RevokedBy = userId

	if _, err := h.svcCtx.RoleClient.RevokeElevation(r.Context(), &req); err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithMessage(w, "elevation revoked")
}

func (h *ElevationHandler) List(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.RoleClient.ListElevations(r.Context(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *ElevationHandler) My(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromCtx(r)
	if userId == 0 {
		api.FailWithRequest(w, http.StatusUnauthorized, "missing user identity")
		return
	}

	resp, err := h.svcCtx.RoleClient.GetMyElevations(r.Context(), &common_pb.IdReq{Id: userId})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}
