package sysMonitor

import (
	"context"
	"net/http"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysMonitor/dashboard_pb"
)

type DashboardHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDashboardHandler(svcCtx *svc.ServiceContext) *DashboardHandler {
	return &DashboardHandler{svcCtx: svcCtx}
}

func (h *DashboardHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.DashboardClient.GetStatistics(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *DashboardHandler) GetRecentOperations(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.DashboardClient.GetRecentOperations(context.Background(), &dashboard_pb.RecentOpsReq{Limit: 10})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *DashboardHandler) GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svcCtx.DashboardClient.GetSystemInfo(context.Background(), &common_pb.Empty{})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}
