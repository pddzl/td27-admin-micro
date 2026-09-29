package sysManagement

import (
	"context"
	"net/http"
	"strconv"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/role_pb"
)

type RoleHandler struct {
	svcCtx *svc.ServiceContext
}

func NewRoleHandler(svcCtx *svc.ServiceContext) *RoleHandler {
	return &RoleHandler{svcCtx: svcCtx}
}

func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp, err := h.svcCtx.RoleClient.GetRole(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *RoleHandler) GetRoleWithChildren(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp, err := h.svcCtx.RoleClient.GetRoleWithChildren(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *RoleHandler) ListRole(w http.ResponseWriter, r *http.Request) {
	var req common_pb.PageReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.RoleClient.ListRole(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.RoleClient.GetAllRoles(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req role_pb.CreateRoleReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.RoleClient.CreateRole(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var req role_pb.UpdateRoleReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.RoleClient.UpdateRole(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	var req common_pb.IdReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.RoleClient.DeleteRole(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *RoleHandler) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	var req role_pb.AssignPermissionsReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.RoleClient.AssignPermissions(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *RoleHandler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp, err := h.svcCtx.RoleClient.GetRolePermissions(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}
