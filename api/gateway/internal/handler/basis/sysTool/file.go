package sysTool

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/pathvar"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysTool/file_pb"
)

type FileHandler struct {
	svcCtx *svc.ServiceContext
}

func NewFileHandler(svcCtx *svc.ServiceContext) *FileHandler {
	return &FileHandler{svcCtx: svcCtx}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	var req file_pb.UploadFileReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.FileClient.UploadFile(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	idStr := pathvar.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid id")
		return
	}
	resp, err := h.svcCtx.FileClient.GetFile(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *FileHandler) ListFile(w http.ResponseWriter, r *http.Request) {
	var req common_pb.PageReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.FileClient.ListFile(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}

func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	idStr := pathvar.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid id")
		return
	}
	resp, err := h.svcCtx.FileClient.DownloadFile(context.Background(), &common_pb.IdReq{Id: id})
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", resp.Mime)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, resp.FileName))
	w.Header().Set("Content-Length", strconv.FormatInt(resp.Size, 10))
	w.Write(resp.FileContent)
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	var req common_pb.IdReq
	if err := api.DecodeAndValidate(r.Body, &req); err != nil {
		api.FailWithRequest(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.svcCtx.FileClient.DeleteFile(context.Background(), &req)
	if err != nil {
		api.FailWithMessage(w, err.Error())
		return
	}
	api.OkWithData(w, resp)
}
