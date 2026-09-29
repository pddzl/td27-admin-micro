package sysTool

import (
	"context"
	"net/http"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysTool/cache_pb"
)

type CacheHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCacheHandler(svcCtx *svc.ServiceContext) *CacheHandler {
	return &CacheHandler{svcCtx: svcCtx}
}

func (h *CacheHandler) GetCache(w http.ResponseWriter, r *http.Request) {
	var req cache_pb.GetCacheReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CacheClient.GetCache(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CacheHandler) SetCache(w http.ResponseWriter, r *http.Request) {
	var req cache_pb.SetCacheReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CacheClient.SetCache(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CacheHandler) DeleteCache(w http.ResponseWriter, r *http.Request) {
	var req cache_pb.DeleteCacheReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CacheClient.DeleteCache(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CacheHandler) ListCache(w http.ResponseWriter, r *http.Request) {
	var req common_pb.PageReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CacheClient.ListCache(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CacheHandler) CleanupExpired(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.CacheClient.CleanupExpired(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}
