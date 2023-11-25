package routes

import (
	"GoProject/goweb/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/",func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
	return r
}