package controller

import (
	"blueself/logic"
	"blueself/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]


// CreatePostHandler创建帖子的处理函数
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 从c获取当前用户的ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 2. 创建帖子
	if err := logic.CreatePost(c,p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context){
	// 1. 获取参数（帖子的ID）
	post_Str := c.Param("id")
	post_id, err := strconv.ParseInt(post_Str, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据ID取出帖子数据（查数据库）
	data, err := logic.GetPostByID(post_id)
	if err != nil{
		zap.L().Error("logic.GetPostByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的接口
func  GetPostListHandler(c *gin.Context){
	// 获取分页参数
	pageNum, pageSize, err := getPageInfo(c)
	if err != nil{
		zap.L().Error("controller.getPageInfo() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 获取数据
	data, err := logic.GetPostList(pageNum, pageSize)
	if err != nil{
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}





func GetPostListHandler2(c *gin.Context){
	//  GET请求参数:/api/v1/posts2?page=1&size=10&seiz=10&order=time
	// 获取分页参数
	// 初始化结构体时指定初试参数
	p := &models.ParamPostList{
		Page:  1,
		Size: 10,
		Order: models.OrderTime ,
	}
	if err := c.ShouldBindQuery(p);err != nil{
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取数据
	data, err := logic.GetPostListNew(c, p) //更新：合二为一
	if err != nil{
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}



/*
下面的两部分函数将合并为一个handler（上面）函数进行整合

// GetPostListHandler2 升级版帖子接口(根据前端传来参数动态获取)
// 按创建时间 / 分数 排序
// 1. 获取参数
// 2. 去redis查询id列表
// 3. 根据id去数据库查询帖子详细信息
func GetPostListHandler2(c *gin.Context){
	//  GET请求参数:/api/v1/posts2?page=1&size=10&seiz=10&order=time
	// 获取分页参数
	// 初始化结构体时指定初试参数
	p := &models.ParamPostList{
		Page:  1,
		Size: 10,
		Order: models.OrderTime ,
	}
	if err := c.ShouldBindQuery(p);err != nil{
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取数据
	data, err := logic.GetPostList2(c, p)
	if err != nil{
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 根据传入的社区ID查社区的帖子
func GetCommunityPostListHandler(c *gin.Context){
	//  GET请求参数:/api/v1/posts2?page=1&size=10&seiz=10&order=time
	// 获取分页参数
	// 初始化结构体时指定初试参数
	p := &models.ParamCommunityPostList{
		ParamPostList: &models.ParamPostList{Page:  1,
		Size: 10,
		Order: models.OrderTime ,
		},
		CommunityID: 1,
		
	}
	if err := c.ShouldBindQuery(p);err != nil{
		zap.L().Error("ParamCommunityPostList with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
 
	// 获取数据
	data, err := logic.GetCommunityPostList2(c, p)
	if err != nil{
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}
*/