package model

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(data interface{}) Response {
	return Response{Code: 0, Message: "ok", Data: data}
}

func Fail(code int, message string) Response {
	return Response{Code: code, Message: message, Data: nil}
}

func JSON(c *gin.Context, status int, resp Response) {
	c.JSON(status, resp)
}
