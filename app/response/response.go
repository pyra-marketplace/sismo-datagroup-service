package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, &ResponseBody{
		Msg:  "success",
		Data: data,
	})
}

func WithError(c *gin.Context, err Error) {
	c.JSON(err.StatusCode(), &ResponseBody{
		Code: err.ErrorCode(),
		Msg:  err.Msg(),
		Data: nil,
	})
	c.Abort()
}

func Middleware(c *gin.Context) {
	c.Next()
	if c.Errors != nil {
		for _, err := range c.Errors {
			if respErr, ok := err.Err.(Error); ok {
				WithError(c, respErr)
				return
			}
		}
	}
}
