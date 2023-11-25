package controller

import (
	"blueself/logic"
	"blueself/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型的断言
		if !ok {
			zap.L().Error("vote validator failed and trans failed", zap.Error(err))
			ResponseError(c, CodeInvalidParam)
		} else {
			zap.L().Error("vote validator failed", zap.Error(errs))
			errData := removeTopStruct(errs.Translate(trans)) //翻译并去除错误提示中的结构体前缀
			ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		}
		return
	}
	// 获取当前请求的用户的ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 具体的业务逻辑
	if err := logic.PostVote(c, userID, p); err != nil {
		zap.L().Error("logic.PostVote failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
