package sysManagement

import (
	"context"
	"net/http"
	"strconv"

	"td27/api/gateway/internal/middleware"
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

func flattenMenuTree(tree []*menu_pb.MenuTreeResp) []*menu_pb.MenuResp {
	var result []*menu_pb.MenuResp
	for _, node := range tree {
		if node.Menu != nil {
			result = append(result, node.Menu)
		}
		if len(node.Children) > 0 {
			result = append(result, flattenMenuTree(node.Children)...)
		}
	}
	return result
}

func (h *MenuHandler) ListMenu(w http.ResponseWriter, r *http.Request) {
	roleIds, _ := r.Context().Value(middleware.RoleIdsKey).([]interface{})
	ids := make([]uint64, 0, len(roleIds))
	for _, v := range roleIds {
		if id, ok := v.(float64); ok {
			ids = append(ids, uint64(id))
		}
	}
	resp, err := h.svcCtx.MenuClient.GetUserMenus(context.Background(), &menu_pb.GetUserMenusReq{RoleIds: ids})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, flattenMenuTree(resp.Tree))
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
