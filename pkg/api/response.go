package api

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
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

var protoMarshaler = &protojson.MarshalOptions{
	UseProtoNames: true,
}

// Clean protojson output: convert numeric strings back to actual numbers
// protojson serializes int64 as strings like "3". We convert them to 3.
// ISO date strings like "2026-04-06T18:54:30Z" are NOT converted (contain non-digit chars).
func cleanNumbers(v interface{}) interface{} {
	switch x := v.(type) {
	case string:
		if len(x) > 0 && x[0] >= '0' && x[0] <= '9' {
			allDigits := true
			for _, c := range x {
				if c < '0' || c > '9' {
					allDigits = false
					break
				}
			}
			if allDigits {
				var num json.Number
				if err := json.Unmarshal([]byte(x), &num); err == nil {
					if i, err := num.Int64(); err == nil {
						return i
					}
				}
			}
		}
		return x
	case map[string]interface{}:
		for k, val := range x {
			x[k] = cleanNumbers(val)
		}
		return x
	case []interface{}:
		for i, val := range x {
			x[i] = cleanNumbers(val)
		}
		return x
	default:
		return v
	}
}

func result(w http.ResponseWriter, status int, code int, data interface{}, msg string) {
	// Convert protobuf to JSON using protojson (nice timestamps, camelCase fields)
	var jsonData interface{}
	if msg, ok := data.(proto.Message); ok {
		b, err := protoMarshaler.Marshal(msg)
		if err != nil {
			httpx.WriteJson(w, http.StatusInternalServerError, Response{Code: ERROR_RES, Msg: "marshal error: " + err.Error()})
			return
		}
		// Parse back to generic map
		var raw interface{}
		if err := json.Unmarshal(b, &raw); err == nil {
			jsonData = raw
		}
	} else {
		jsonData = data
	}

	// Clean: convert string numbers to actual numbers
	jsonData = cleanNumbers(jsonData)

	httpx.WriteJson(w, status, Response{Code: code, Data: jsonData, Msg: msg})
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
