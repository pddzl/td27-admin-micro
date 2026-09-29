package sysManagement

import (
	"context"
	"net/http"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/sysManagement/user_pb"
)

type LoginHandler struct {
	svcCtx *svc.ServiceContext
}

func NewLoginHandler(svcCtx *svc.ServiceContext) *LoginHandler {
	return &LoginHandler{svcCtx: svcCtx}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username  string `json:"username" validate:"required"`
		Password  string `json:"password" validate:"required"`
		CaptchaId string `json:"captcha_id" validate:"required"`
		Captcha   string `json:"captcha" validate:"required"`
	}

	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, err.Error())
		return
	}

	store := GetCaptchaStore()
	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid or expired captcha")
		return
	}

	resp, err := h.svcCtx.UserClient.Login(context.Background(), &user_pb.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		api.FailWithRequest(w, http.StatusOK, err.Error())
		return
	}

	api.OkWithDetailed(w, resp, "登录成功")
}

func (h *LoginHandler) Health(w http.ResponseWriter, r *http.Request) {
	api.OkWithMessage(w, "health")
}
