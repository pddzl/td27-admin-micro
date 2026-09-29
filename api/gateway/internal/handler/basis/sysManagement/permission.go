package sysManagement

import (
	"context"
	"net/http"
	"strconv"
	
	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/permission_pb"
)

type PermissionHandler struct {
	svcCtx *svc.ServiceContext
}

func NewPermissionHandler(svcCtx *svc.ServiceContext) *PermissionHandler {
	return &PermissionHandler{svcCtx: svcCtx}
}

func (h *PermissionHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp, err := h.svcCtx.PermissionClient.GetPermission(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) ListPermission(w http.ResponseWriter, r *http.Request) {
	var req common_pb.PageReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.PermissionClient.ListPermission(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.PermissionClient.GetAllPermissions(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) GetPermissionsByRoleId(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp, err := h.svcCtx.PermissionClient.GetPermissionsByRoleId(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req permission_pb.CreatePermissionReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.PermissionClient.CreatePermission(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	var req permission_pb.UpdatePermissionReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.PermissionClient.UpdatePermission(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	var req common_pb.IdReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.PermissionClient.DeletePermission(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	var req permission_pb.CheckPermissionReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.PermissionClient.CheckPermission(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) ReloadPolicy(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.PermissionClient.ReloadPolicy(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *PermissionHandler) LintPolicies(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.PermissionClient.LintPolicies(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}
