package models

const(
	OrderTime = "time"
	OrderScore = "score"
)

// ParamSignup 定义注册请求的参数结构体
type ParamSignup struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 定义登录请求的参数结构体
type ParamLogin struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// ParamVoteData 定义投票请求的参数结构体
type ParamVoteData struct {
	// UserID 从请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`  //帖子ID
	Direction *int8   `json:"direction,string" binding:"oneof=1 0 -1"` //赞成（1）反对（-1）取消投票(0)
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct{
	Page int64 `json:"page" form:"page"`
	Size int64 `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
	CommunityID int `json:"community_id" form:"community_id"` //可以为空
}

// ParamCommunityPostList 按社区获取帖子列表query string参数
type ParamCommunityPostList struct{
	*ParamPostList
	CommunityID int `json:"community_id" form:"community_id"` //可以为空

}