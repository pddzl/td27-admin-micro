package sysManagement

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/pathvar"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/api_pb"
)

type APIHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAPIHandler(svcCtx *svc.ServiceContext) *APIHandler {
	return &APIHandler{svcCtx: svcCtx}
}

func (h *APIHandler) GetAPI(w http.ResponseWriter, r *http.Request) {
	idStr := pathvar.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.svcCtx.APIClient.GetAPI(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *APIHandler) ListAPI(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, _ := strconv.ParseInt(pageStr, 10, 32)
	pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 32)
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	resp, err := h.svcCtx.APIClient.ListAPI(context.Background(), &common_pb.PageReq{
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *APIHandler) GetAPIsByGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GroupEn string `json:"group_en" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.APIClient.GetAPIsByGroup(context.Background(), &api_pb.GetAPIsByGroupReq{GroupEn: req.GroupEn})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *APIHandler) CreateAPI(w http.ResponseWriter, r *http.Request) {
	var req api_pb.CreateAPIReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.svcCtx.APIClient.CreateAPI(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithMessage(w, "创建成功")
}

func (h *APIHandler) UpdateAPI(w http.ResponseWriter, r *http.Request) {
	var req api_pb.UpdateAPIReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.APIClient.UpdateAPI(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithDetailed(w, resp, "更新成功")
}

func (h *APIHandler) DeleteAPI(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id int64 `json:"id" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.svcCtx.APIClient.DeleteAPI(context.Background(), &common_pb.IdReq{Id: req.Id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithMessage(w, "删除成功")
}

func (h *APIHandler) DeleteByIds(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Ids []int64 `json:"ids" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.svcCtx.APIClient.DeleteByIds(context.Background(), &common_pb.IdsReq{Ids: req.Ids})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithMessage(w, "删除成功")
}

func (h *APIHandler) GetAPITree(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id int64 `json:"id"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.APIClient.GetAPITree(context.Background(), &api_pb.APITreeReq{RoleId: req.Id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	// Transform to frontend format: {list: [{key, children}], checkedIds: [...]}
	type apiNode struct {
		Id          uint       `json:"id"`
		Key         string     `json:"key"`
		Path        string     `json:"path,omitempty"`
		Method      string     `json:"method,omitempty"`
		GroupEN     string     `json:"group_en,omitempty"`
		GroupCN     string     `json:"group_cn,omitempty"`
		Description string     `json:"description,omitempty"`
		Children    []apiNode  `json:"children,omitempty"`
	}
	tree := make([]apiNode, 0, len(resp.List))
	for _, item := range resp.List {
		node := apiNode{
			Key:     item.GroupEn,
			GroupEN: item.GroupEn,
			GroupCN: item.GroupCn,
		}
		for _, api := range item.Apis {
			node.Children = append(node.Children, apiNode{
				Id:          uint(api.Id),
				Key:         api.Path + ":" + api.Method,
				Path:        api.Path,
				Method:      api.Method,
				GroupEN:     api.GroupEn,
				GroupCN:     api.GroupCn,
				Description: api.Description,
			})
		}
		tree = append(tree, node)
	}
	api.OkWithDetailed(w, map[string]interface{}{
		"list":        tree,
		"checkedIds": resp.CheckedIds,
	}, "获取成功")
}
