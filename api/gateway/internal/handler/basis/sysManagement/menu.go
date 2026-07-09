package sysManagement

import (
	"context"
	"net/http"
	"strconv"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/menu_pb"
)

type MenuHandler struct {
	svcCtx *svc.ServiceContext
}

func NewMenuHandler(svcCtx *svc.ServiceContext) *MenuHandler {
	return &MenuHandler{svcCtx: svcCtx}
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	resp, err := h.svcCtx.MenuClient.GetMenu(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *MenuHandler) GetMenuTree(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.MenuClient.GetMenuTree(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *MenuHandler) GetUserMenus(w http.ResponseWriter, r *http.Request) {
	var req menu_pb.GetUserMenusReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.MenuClient.GetUserMenus(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *MenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var req menu_pb.CreateMenuReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.MenuClient.CreateMenu(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *MenuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	var req menu_pb.UpdateMenuReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.MenuClient.UpdateMenu(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *MenuHandler) ListMenu(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.MenuClient.ListMenu(context.Background(), &common_pb.PageReq{Page: 1, PageSize: 999})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *MenuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	var req common_pb.IdReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.MenuClient.DeleteMenu(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}
