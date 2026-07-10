package api

import (
	"encoding/json"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	SUCCESS   = 0
	ERROR_REQ = 4
	ERROR_RES = 7
)

func marshalData(v interface{}) ([]byte, error) {
	if msg, ok := v.(proto.Message); ok {
		return protojson.Marshal(msg)
	}
	return json.Marshal(v)
}

func writeJson(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func result(w http.ResponseWriter, status int, code int, data interface{}, msg string) {
	dataBytes, err := marshalData(data)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, Response{Code: ERROR_RES, Msg: "marshal error: " + err.Error()})
		return
	}
	writeJson(w, status, struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
		Msg  string          `json:"msg"`
	}{
		Code: code,
		Data: dataBytes,
		Msg:  msg,
	})
}

func Ok(w http.ResponseWriter) {
	result(w, http.StatusOK, SUCCESS, map[string]interface{}{}, "操作成功")
}

func OkWithMessage(w http.ResponseWriter, msg string) {
	result(w, http.StatusOK, SUCCESS, map[string]interface{}{}, msg)
}

func OkWithData(w http.ResponseWriter, data interface{}) {
	result(w, http.StatusOK, SUCCESS, data, "查询成功")
}

func OkWithDetailed(w http.ResponseWriter, data interface{}, msg string) {
	result(w, http.StatusOK, SUCCESS, data, msg)
}

func Fail(w http.ResponseWriter) {
	result(w, http.StatusOK, ERROR_RES, map[string]interface{}{}, "操作失败")
}

func FailReq(w http.ResponseWriter, msg string) {
	result(w, http.StatusOK, ERROR_REQ, map[string]interface{}{}, msg)
}

func FailWithDetailed(w http.ResponseWriter, data interface{}, msg string) {
	result(w, http.StatusOK, ERROR_RES, data, msg)
}

func FailWithMessage(w http.ResponseWriter, msg string) {
	result(w, http.StatusOK, ERROR_RES, map[string]interface{}{}, msg)
}

func FailWithRequest(w http.ResponseWriter, status int, msg string) {
	result(w, status, ERROR_REQ, map[string]interface{}{}, msg)
}

func Custom(w http.ResponseWriter, status int, code int, data interface{}, msg string) {
	result(w, status, code, data, msg)
}
