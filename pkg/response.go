package pkg

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Success(data interface{}) *Response {
	return &Response{Code: 0, Data: data, Msg: "success"}
}

func SuccessWithMsg(data interface{}, msg string) *Response {
	return &Response{Code: 0, Data: data, Msg: msg}
}

func Error(code int, msg string) *Response {
	return &Response{Code: code, Data: nil, Msg: msg}
}

func WriteJson(w http.ResponseWriter, status int, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
