package sysMonitor

import (
	"context"
	"net/http"
	"strconv"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysMonitor/operation_log_pb"
)

type OperationLogHandler struct {
	svcCtx *svc.ServiceContext
}

func NewOperationLogHandler(svcCtx *svc.ServiceContext) *OperationLogHandler {
	return &OperationLogHandler{svcCtx: svcCtx}
}

func (h *OperationLogHandler) ListOperationLog(w http.ResponseWriter, r *http.Request) {
	var flat struct {
		Page     int     `json:"page"`
		PageSize int     `json:"pageSize"`
		Path     *string `json:"path"`
		Method   *string `json:"method"`
		Status   *int32  `json:"status"`
	}
	if err := api.DecodeAndValidate(r.Body, &flat); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req := &operation_log_pb.ListOperationLogReq{
		Page: &common_pb.PageReq{
			Page:     uint32(flat.Page),
			PageSize: uint32(flat.PageSize),
		},
		Path:   flat.Path,
		Method: flat.Method,
		Status: flat.Status,
	}
	resp, err := h.svcCtx.OperationLogClient.ListOperationLog(context.Background(), req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *OperationLogHandler) CleanupExpiredLogs(w http.ResponseWriter, r *http.Request) {
	var req operation_log_pb.CleanupExpiredLogsReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.OperationLogClient.CleanupExpiredLogs(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *OperationLogHandler) DeleteOperationLog(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id int64 `json:"id"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	_, err := h.svcCtx.OperationLogClient.Delete(context.Background(), &common_pb.IdReq{Id: req.Id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithMessage(w, "删除成功")
}

func (h *OperationLogHandler) DeleteOperationLogByIds(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Ids []interface{} `json:"ids"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ids := make([]int64, 0, len(req.Ids))
	for _, v := range req.Ids {
		switch val := v.(type) {
		case float64:
			ids = append(ids, int64(val))
		case string:
			id, _ := strconv.ParseInt(val, 10, 64)
			if id > 0 {
				ids = append(ids, id)
			}
		}
	}
	_, err := h.svcCtx.OperationLogClient.DeleteByIds(context.Background(), &common_pb.IdsReq{Ids: ids})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithMessage(w, "删除成功")
}
