package sysTool

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/pathvar"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysTool/cron_pb"
)

type CronHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCronHandler(svcCtx *svc.ServiceContext) *CronHandler {
	return &CronHandler{svcCtx: svcCtx}
}

func (h *CronHandler) GetCron(w http.ResponseWriter, r *http.Request) {
	idStr := pathvar.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid id")
		return
	}
	resp, err := h.svcCtx.CronClient.GetCron(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CronHandler) ListCron(w http.ResponseWriter, r *http.Request) {
	var req common_pb.PageReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CronClient.ListCron(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CronHandler) CreateCron(w http.ResponseWriter, r *http.Request) {
	var req cron_pb.CreateCronReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CronClient.CreateCron(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CronHandler) UpdateCron(w http.ResponseWriter, r *http.Request) {
	var req cron_pb.UpdateCronReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CronClient.UpdateCron(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CronHandler) ToggleCronStatus(w http.ResponseWriter, r *http.Request) {
	var req cron_pb.ToggleCronStatusReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CronClient.ToggleCronStatus(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CronHandler) ExecuteCronNow(w http.ResponseWriter, r *http.Request) {
	var req cron_pb.ExecuteCronNowReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CronClient.ExecuteCronNow(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CronHandler) DeleteCron(w http.ResponseWriter, r *http.Request) {
	var req common_pb.IdReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.CronClient.DeleteCron(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *CronHandler) DeleteByIds(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Ids []int64 `json:"ids" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.svcCtx.CronClient.DeleteByIds(context.Background(), &common_pb.IdsReq{Ids: req.Ids})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}
