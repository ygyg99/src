package mysql

import (
	"blueself/models"
	"blueself/settings"
	"testing"
)

// 还需要对数据库进行初始化(写测试连的数据库，不要连接正式数据库！！)
func init() {
	dbCfg := settings.MySQLConfig{
		Host:         "",
		User:         "",
		Password:     "",
		DB:           "",
		Port:         0,
		MaxOpenConns: 0,
		MaxIdleConns: 0,
	}
	if dbCfg.DB == "nil"{
		return
	}
	err := Init()
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {

	post := models.Post{
		ID:          10,
		AuthorID:    123,
		CommunityID: 1,
		Status:      1,
		Title:       "test",
		Content:     "test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err: %v", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
