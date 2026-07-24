package sysManagement

import (
	"context"
	"net/http"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/dict_detail_pb"
)

type DictDetailHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDictDetailHandler(svcCtx *svc.ServiceContext) *DictDetailHandler {
	return &DictDetailHandler{svcCtx: svcCtx}
}

func (h *DictDetailHandler) CreateDictDetail(w http.ResponseWriter, r *http.Request) {
	var req dict_detail_pb.CreateDictDetailReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.svcCtx.DictDetailClient.CreateDictDetail(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	// Fetch created detail from flat list
	flatResp, err := h.svcCtx.DictDetailClient.FlatDictDetails(context.Background(), &dict_detail_pb.FlatDictDetailsReq{DictId: req.DictId})
	if err != nil {
		api.OkWithMessage(w, "创建成功")
		return
	}
	if len(flatResp.List) > 0 {
		api.OkWithDetailed(w, flatResp.List[len(flatResp.List)-1], "创建成功")
	} else {
		api.OkWithMessage(w, "创建成功")
	}
}

func (h *DictDetailHandler) UpdateDictDetail(w http.ResponseWriter, r *http.Request) {
	var req dict_detail_pb.UpdateDictDetailReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DictDetailClient.UpdateDictDetail(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithDetailed(w, resp, "更新成功")
}

func (h *DictDetailHandler) DeleteDictDetail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id int64 `json:"id" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DictDetailClient.DeleteDictDetail(context.Background(), &common_pb.IdReq{Id: req.Id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DictDetailHandler) FlatDictDetails(w http.ResponseWriter, r *http.Request) {
	var req dict_detail_pb.FlatDictDetailsReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DictDetailClient.FlatDictDetails(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DictDetailHandler) ListDictDetail(w http.ResponseWriter, r *http.Request) {
	var req dict_detail_pb.ListDictDetailReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DictDetailClient.ListDictDetail(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}
