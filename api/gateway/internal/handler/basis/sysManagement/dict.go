package sysManagement

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/pathvar"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/dict_pb"
)

type DictHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDictHandler(svcCtx *svc.ServiceContext) *DictHandler {
	return &DictHandler{svcCtx: svcCtx}
}

func (h *DictHandler) GetDict(w http.ResponseWriter, r *http.Request) {
	idStr := pathvar.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.svcCtx.DictClient.GetDict(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DictHandler) GetDictByENName(w http.ResponseWriter, r *http.Request) {
	enName := pathvar.Vars(r)["en_name"]

	resp, err := h.svcCtx.DictClient.GetDictByENName(context.Background(), &dict_pb.GetDictByENNameReq{EnName: enName})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DictHandler) ListDict(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	resp, err := h.svcCtx.DictClient.ListDict(context.Background(), &common_pb.PageReq{
		Page:     uint32(req.Page),
		PageSize: uint32(req.PageSize),
	})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DictHandler) CreateDict(w http.ResponseWriter, r *http.Request) {
	var req dict_pb.CreateDictReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DictClient.CreateDict(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DictHandler) UpdateDict(w http.ResponseWriter, r *http.Request) {
	var req dict_pb.UpdateDictReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DictClient.UpdateDict(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *DictHandler) DeleteDict(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id int64 `json:"id" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.DictClient.DeleteDict(context.Background(), &common_pb.IdReq{Id: req.Id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}


