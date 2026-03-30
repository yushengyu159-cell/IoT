package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(r *ghttp.Request, message string, data ...interface{}) {
	resp := Response{Code: 200, Message: message}
	if len(data) > 0 {
		resp.Data = data[0]
	}
	r.Response.WriteJson(resp)
}

func Error(r *ghttp.Request, code int, message string) {
	resp := Response{Code: code, Message: message}
	r.Response.WriteJson(resp)
}

func BadRequest(r *ghttp.Request, message string) {
	Error(r, 400, message)
}

func Unauthorized(r *ghttp.Request, message string) {
	Error(r, 401, message)
}

func Forbidden(r *ghttp.Request, message string) {
	Error(r, 403, message)
}

func InternalError(r *ghttp.Request, message string) {
	Error(r, 500, message)
}
