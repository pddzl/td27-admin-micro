package api

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
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

func result(w http.ResponseWriter, status int, code int, data interface{}, msg string) {
	httpx.WriteJson(w, status, Response{
		Code: code,
		Data: data,
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

func Custom(
	w http.ResponseWriter,
	status int,
	code int,
	data interface{},
	msg string,
) {
	result(w, status, code, data, msg)
}
