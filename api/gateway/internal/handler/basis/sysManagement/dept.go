package sysManagement

import (
	"context"
	"github.com/zeromicro/go-zero/rest/pathvar"
	"net/http"
	"strconv"
	"td27/api/gateway/internal/svc"
	"td27/pkg/api"

	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/dept_pb"
)

type DeptHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeptHandler(svcCtx *svc.ServiceContext) *DeptHandler {
	return &DeptHandler{svcCtx: svcCtx}
}

func (h *DeptHandler) GetDept(w http.ResponseWriter, r *http.Request) {
	idStr := pathvar.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.svcCtx.DeptClient.GetDept(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DeptHandler) ListDept(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := uint32(1)
	pageSize := uint32(20)

	if p, err := strconv.ParseUint(pageStr, 10, 32); err == nil && pageStr != "" {
		page = uint32(p)
	}
	if ps, err := strconv.ParseUint(pageSizeStr, 10, 32); err == nil && pageSizeStr != "" {
		pageSize = uint32(ps)
	}

	resp, err := h.svcCtx.DeptClient.ListDept(context.Background(), &dept_pb.ListDeptReq{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func deptTreeToResp(tree []*dept_pb.DeptTreeResp) []map[string]interface{} {
	var result []map[string]interface{}
	for _, node := range tree {
		if node.Dept == nil {
			continue
		}
		item := map[string]interface{}{
			"id":       node.Dept.Id,
			"deptName": node.Dept.DeptName,
			"parentId": node.Dept.ParentId,
			"sort":     node.Dept.Sort,
			"status":   node.Dept.Status,
		}
		if len(node.Children) > 0 {
			item["children"] = deptTreeToResp(node.Children)
		}
		result = append(result, item)
	}
	return result
}

func (h *DeptHandler) GetElTreeDepts(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.DeptClient.GetDeptTree(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, map[string]interface{}{
		"tree": deptTreeToResp(resp.Tree),
		"ids":  []uint64{},
	})
}

func (h *DeptHandler) GetDeptTree(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.DeptClient.GetDeptTree(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DeptHandler) CreateDept(w http.ResponseWriter, r *http.Request) {
	var req dept_pb.CreateDeptReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DeptClient.CreateDept(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DeptHandler) UpdateDept(w http.ResponseWriter, r *http.Request) {
	var req dept_pb.UpdateDeptReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DeptClient.UpdateDept(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DeptHandler) DeleteDept(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id uint64 `json:"id" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DeptClient.DeleteDept(context.Background(), &common_pb.IdReq{Id: req.Id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}
