package mysql

import (
	"blueself/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "yg"

//把每一步数据库操作封装成函数，等待logic层根据业务需求调用

func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InserUser 向数据库中加入用户记录
func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	user.Password = encode(user.Password)
	//执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encode(oP string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oP)))
}

func Login(p *models.User) (err error) {
	oP := p.Password
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(p, sqlStr, p.Username)
	// 判断用户名是否存在
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败了
		return err
	}
	// 判断密码是否正确
	if encode(oP) != p.Password {
		return ErrWrongPassword
	}
	return
}

// GetUserByID 根据id获取用户信息
func GetUserByID(id int64)(user *models.User, err error){
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}