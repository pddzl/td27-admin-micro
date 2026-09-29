package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"td27/api/gateway/internal/svc"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"td27/rpc/basis/types/sysMonitor/operation_log_pb"
)

// maxLoggedBodySize caps request/response payloads persisted to the operation log,
// so large uploads/downloads don't bloat memory and the DB.
const maxLoggedBodySize = 64 * 1024

func truncateForLog(s string) string {
	if len(s) > maxLoggedBodySize {
		return s[:maxLoggedBodySize] + "...(truncated)"
	}
	return s
}

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
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			reqParam = "<multipart form data omitted>"
		}

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
			ReqParam:  truncateForLog(reqParam),
			RespData:  truncateForLog(rw.body.String()),
			RespTime:  time.Since(now).Milliseconds(),
			UserId:    int64(userId),
			UserName:  username,
		}

		go func() {
			if _, err := m.svcCtx.OperationLogClient.CreateOperationLog(context.Background(), req); err != nil {
				logx.Errorf("create operation log failed for %s %s: %v", req.Method, req.Path, err)
			}
		}()
	}
}
