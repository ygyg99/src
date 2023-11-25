package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	// 直接调用routes层时容易导致循环引用
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	// 原步骤中有一个获取用户ID部分，加入没加入的话，会卡死在这一步
	r.POST(url, CreatePostHandler)
	// 设置要传入的测试用例，一个json
	body := `{
		"community_id": 1,
		"title": "test",
		"content": "just a test"
	}`

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	r.ServeHTTP(w, req)

	// 判断是否有达到预期的输出(没有用户ID传入时会卡死在那一步)
	assert.Equal(t, 200, w.Code)

	// 方法1，判断响应内容中是否包含指定的字符串
	// assert.Contains(t, w.Body.String(), "需要登录")

	// 方法2，将响应的内容反序列化到Response 然后判断字段与预期是否一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed, err: %v", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)
}
