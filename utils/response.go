package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	LoginBinding = 601 //用户需绑定
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{code, data, msg})
}
func StatusResult(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(code, Response{code, data, msg})
}
func Ok(c *gin.Context) {
	StatusResult(200, map[string]interface{}{}, "操作成功", c)
}
func OkWithData(data interface{}, c *gin.Context) {
	StatusResult(200, data, "操作成功", c)
}

/*
Fail
*/
//403 服务器拒绝请求
func FailForbidden(message string,c *gin.Context) {
	StatusResult(http.StatusForbidden, map[string]interface{}{}, "操作失败", c)
}
//401	鉴权
func FailUnauthorized(message string, c *gin.Context){
	StatusResult(http.StatusUnauthorized, map[string]interface{}{}, message, c)
}
//404  找不到资源
func FailNotFound ( message string, c *gin.Context) {
	StatusResult(http.StatusNotFound, map[string]interface{}{}, message, c)
}
//500	执行中服务器内部错误
func FailInternalServerError(message string, c *gin.Context) {
	StatusResult(http.StatusInternalServerError, map[string]interface{}{}, message, c)
}
//501   参数错误
func FailNotImplemented(message string, c *gin.Context) {
	StatusResult(http.StatusNotImplemented, map[string]interface{}{}, message, c)
}
//502   参数错误
func FailStatusBadGateway(message string, c *gin.Context) {
	StatusResult(http.StatusBadGateway, map[string]interface{}{}, message, c)
}

func FailCodeMsg(code int,message string, c *gin.Context) {
	StatusResult(code, map[string]interface{}{}, message, c)
}

