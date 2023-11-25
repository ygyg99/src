package logic

import (
	"blueself/database/mysql"
	"blueself/models"
	"blueself/pkg/jwt"
	"blueself/pkg/snowflake"
)

//存放业务逻辑的代码

func SignUp(p *models.ParamSignup) (err error) {
	//判断用户是否存在

	if err := mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询出错
		return err
	}

	//生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//密码加密

	//保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {

	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user。UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT Token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
