package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"td27/api/gateway/internal/svc"
	"time"

	"td27/rpc/basis/types/sysMonitor/operation_log_pb"
)

type OperationRecordMiddleware struct {
	svcCtx *svc.ServiceContext
}

func NewOperationRecordMiddleware(svcCtx *svc.ServiceContext) *OperationRecordMiddleware {
	return &OperationRecordMiddleware{svcCtx: svcCtx}
}

type responseWriter struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (m *OperationRecordMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		body, _ := io.ReadAll(r.Body)
		reqParam := string(body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		rw := &responseWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
			status:         http.StatusOK,
		}
		now := time.Now()

		next(rw, r)

		userId, _ := r.Context().Value(UserIdKey).(float64)
		username, _ := r.Context().Value(UsernameKey).(string)

		req := &operation_log_pb.CreateOperationLogReq{
			Ip:        strings.Split(r.RemoteAddr, ":")[0],
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    int32(rw.status),
			UserAgent: r.UserAgent(),
			ReqParam:  reqParam,
			RespData:  rw.body.String(),
			RespTime:  time.Since(now).Milliseconds(),
			UserId:    int64(userId),
			UserName:  username,
		}

		go m.svcCtx.OperationLogClient.CreateOperationLog(context.Background(), req)
	}
}
