package sysManagement

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/pathvar"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/button_pb"
)

type ButtonHandler struct {
	svcCtx *svc.ServiceContext
}

func NewButtonHandler(svcCtx *svc.ServiceContext) *ButtonHandler {
	return &ButtonHandler{svcCtx: svcCtx}
}

func (h *ButtonHandler) GetButton(w http.ResponseWriter, r *http.Request) {
	idStr := pathvar.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.svcCtx.ButtonClient.GetButton(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *ButtonHandler) ListButton(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.svcCtx.ButtonClient.ListButton(context.Background(), &common_pb.PageReq{
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *ButtonHandler) GetButtonsByPagePath(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PagePath string `json:"page_path" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.ButtonClient.GetButtonsByPagePath(context.Background(), &button_pb.GetButtonsByPagePathReq{PagePath: req.PagePath})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *ButtonHandler) GetUserButtons(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoleIds []int64 `json:"role_ids" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.ButtonClient.GetUserButtons(context.Background(), &button_pb.GetUserButtonsReq{RoleIds: req.RoleIds})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *ButtonHandler) BatchCheckPermission(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ButtonCodes []string `json:"button_codes" validate:"required"`
		RoleIds     []int64 `json:"role_ids" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.ButtonClient.BatchCheckPermission(context.Background(), &button_pb.BatchCheckPermissionReq{
		ButtonCodes: req.ButtonCodes,
		RoleIds:     req.RoleIds,
	})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *ButtonHandler) CreateButton(w http.ResponseWriter, r *http.Request) {
	var req button_pb.CreateButtonReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.ButtonClient.CreateButton(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *ButtonHandler) UpdateButton(w http.ResponseWriter, r *http.Request) {
	var req button_pb.UpdateButtonReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.ButtonClient.UpdateButton(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}

func (h *ButtonHandler) DeleteButton(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id int64 `json:"id" validate:"required"`
	}
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svcCtx.ButtonClient.DeleteButton(context.Background(), &common_pb.IdReq{Id: req.Id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}

	api.OkWithData(w, resp)
}
