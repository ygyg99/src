package controller

import (
	"blueself/database/mysql"
	"blueself/logic"
	"blueself/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	
	//1. 获取参数，参数校验
	var p models.ParamSignup
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("signup with invalid param", zap.Error(err))
		// 判断err是不是validator(实际上是采用了一个类型转换的方式)
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//手动校验参数
	// if len(p.Username) ==0 || len(p.Password)==0||len(p.RePassword)==0 || p.Password != p.RePassword {
	// 	//请求参数有误，直接返回响应
	// 	zap.L().Error("signup with invalid param")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }
	// fmt.Println(p)
	
	//2. 业务处理
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3. 返回响应
	ResponseSuccess(c, nil)

}

func LogInHandler(c *gin.Context) {
	//1. 获取参数,并在数据库中检验参数是否合理
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("login with invalid params", zap.Error(err))
		// 判断err是不是validator(实际上是采用了一个类型转换的方式)
		Err, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(Err.Translate(trans)))
		return
	}
	// 2.业务处理
	user, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("Login failed", zap.String("user:", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		if errors.Is(err, mysql.ErrWrongPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // id值大于1<<53-1  int64类型的最大值是1<<63-1
		"user_name": user.Username,
		"token":     user.Token,
	})

}
