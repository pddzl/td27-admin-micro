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
	ERROR_REQ = 4 // request error
	ERROR_RES = 7 // response/business error
)

func marshalJSON(v interface{}) ([]byte, error) {
	if msg, ok := v.(proto.Message); ok {
		return protojson.Marshal(msg)
	}
	return json.Marshal(v)
}

func result(w http.ResponseWriter, status int, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	resp := struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
		Msg  string          `json:"msg"`
	}{
		Code: code,
		Msg:  msg,
	}

	if data != nil {
		b, err := marshalJSON(data)
		if err != nil {
			resp.Code = ERROR_RES
			resp.Msg = "marshal error: " + err.Error()
			resp.Data = json.RawMessage("null")
		} else {
			resp.Data = json.RawMessage(b)
		}
	} else {
		resp.Data = json.RawMessage("{}")
	}

	full, _ := json.Marshal(resp)
	w.Write(full)
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

func Custom(
	w http.ResponseWriter,
	status int,
	code int,
	data interface{},
	msg string,
) {
	result(w, status, code, data, msg)
}
