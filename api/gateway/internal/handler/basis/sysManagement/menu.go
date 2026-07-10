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
	id, _ := strconv.ParseInt(idStr, 10, 64)
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

func (h *MenuHandler) GetElTreeMenus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id int64 `json:"id"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get all menus as tree
	treeResp, err := h.svcCtx.MenuClient.GetMenuTree(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	// Transform MenuTreeResp to flat menu nodes (unwrap Menu wrapper)
	var transform func(nodes []*menu_pb.MenuTreeResp) []map[string]interface{}
	transform = func(nodes []*menu_pb.MenuTreeResp) []map[string]interface{} {
		out := make([]map[string]interface{}, 0, len(nodes))
		for _, n := range nodes {
			if n.Menu == nil {
				continue
			}
			item := map[string]interface{}{
				"id":        n.Menu.Id,
				"menu_name": n.Menu.MenuName,
				"icon":      n.Menu.Icon,
				"path":      n.Menu.Path,
				"component": n.Menu.Component,
				"redirect":  n.Menu.Redirect,
				"parentId":  n.Menu.ParentId,
				"sort":      n.Menu.Sort,
				"hidden":    n.Menu.Hidden,
				"keepAlive": n.Menu.KeepAlive,
				"affix":     n.Menu.Affix,
				"alwaysShow": n.Menu.AlwaysShow,
				"title":     n.Menu.Title,
			}
			if len(n.Children) > 0 {
				item["children"] = transform(n.Children)
			}
			out = append(out, item)
		}
		return out
	}

	api.OkWithDetailed(w, map[string]interface{}{
		"list":    transform(treeResp.Tree),
		"menuIds": []int64{},
	}, "获取成功")
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

func menuTreeToResp(tree []*menu_pb.MenuTreeResp) []map[string]interface{} {
	var result []map[string]interface{}
	for _, node := range tree {
		if node.Menu == nil {
			continue
		}
		item := map[string]interface{}{
			"id":          node.Menu.Id,
			"menu_name":   node.Menu.MenuName,
			"icon":        node.Menu.Icon,
			"path":        node.Menu.Path,
			"component":   node.Menu.Component,
			"redirect":    node.Menu.Redirect,
			"parentId":    node.Menu.ParentId,
			"sort":        node.Menu.Sort,
			"hidden":      node.Menu.Hidden,
			"keepAlive":   node.Menu.KeepAlive,
			"affix":       node.Menu.Affix,
			"alwaysShow":  node.Menu.AlwaysShow,
			"title":       node.Menu.Title,
			"createdAt":   node.Menu.CreatedAt,
			"updatedAt":   node.Menu.UpdatedAt,
		}
		if len(node.Children) > 0 {
			item["children"] = menuTreeToResp(node.Children)
		}
		result = append(result, item)
	}
	return result
}

func (h *MenuHandler) ListMenu(w http.ResponseWriter, r *http.Request) {
	roleIds, _ := r.Context().Value(middleware.RoleIdsKey).([]interface{})
	ids := make([]int64, 0, len(roleIds))
	for _, v := range roleIds {
		if id, ok := v.(float64); ok {
			ids = append(ids, int64(id))
		}
	}
	resp, err := h.svcCtx.MenuClient.GetUserMenus(context.Background(), &menu_pb.GetUserMenusReq{RoleIds: ids})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, menuTreeToResp(resp.Tree))
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
