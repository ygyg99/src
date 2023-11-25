package logic

import (
	"blueself/database/mysql"
	"blueself/database/redis"
	"blueself/models"
	"blueself/pkg/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePost(c *gin.Context, p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err = redis.CreatePost(c, p.ID, p.CommunityID)
	if err != nil {
		return
	}
	err = mysql.CreatePost(p)
	return
	// 3. 返回
}

func GetPostByID(id int64) (data *models.ApiPostDetail, err error) {
	// 查询并拼接相关的数据
	post, err := mysql.GetPostByID(id)
	if err != nil {
		zap.L().Error("mysql.GetPostByID failed", zap.Int64("post_id", id), zap.Error(err))
		return
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}

	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed",
			zap.Int64("CommunityID", post.CommunityID),
			zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(pageNum, pageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, er := mysql.GetUserByID(post.AuthorID)
		if er != nil {
			zap.L().Error("mysql.GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(er))
			continue
		}

		// 根据社区id查询社区详细信息
		community, er := mysql.GetCommunityDetailByID(post.CommunityID)
		if er != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("CommunityID", post.CommunityID),
				zap.Error(er))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(c *gin.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 1. 去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(c, p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}

	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 2. 根据id去数据库查询帖子详细信息
	// 返回的数据需要按照给定的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提取查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(c, ids)

	// 将梯子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, er := mysql.GetUserByID(post.AuthorID)
		if er != nil {
			zap.L().Error("mysql.GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(er))
			continue
		}

		// 根据社区id查询社区详细信息
		community, er := mysql.GetCommunityDetailByID(post.CommunityID)
		if er != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("CommunityID", post.CommunityID),
				zap.Error(er))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList2(c *gin.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
		

		// 1. 去redis查询id列表
		ids, err := redis.GetCommunityPostIDsInOrder(c, p)
		if err != nil {
			return
		}
		if len(ids) == 0 {
			zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
			return
		}
	
		zap.L().Debug("GetPostList2", zap.Any("ids", ids))
		// 2. 根据id去数据库查询帖子详细信息
		// 返回的数据需要按照给定的顺序返回
		posts, err := mysql.GetPostListByIDs(ids)
		if err != nil {
			return
		}
		// 提取查询好每篇帖子的投票数
		voteData, err := redis.GetPostVoteData(c, ids)
	
		// 将梯子的作者及分区信息查询出来填充到帖子中
		for idx, post := range posts {
			// 根据作者id查询作者信息
			user, er := mysql.GetUserByID(post.AuthorID)
			if er != nil {
				zap.L().Error("mysql.GetUserByID failed",
					zap.Int64("author_id", post.AuthorID),
					zap.Error(er))
				continue
			}
	
			// 根据社区id查询社区详细信息
			community, er := mysql.GetCommunityDetailByID(post.CommunityID)
			if er != nil {
				zap.L().Error("mysql.GetCommunityDetailByID failed",
					zap.Int64("CommunityID", post.CommunityID),
					zap.Error(er))
				continue
			}
			postdetail := &models.ApiPostDetail{
				AuthorName:      user.Username,
				VoteNum:         voteData[idx],
				Post:            post,
				CommunityDetail: community,
			}
			data = append(data, postdetail)
		}
		return
}

// GetPostListNew 将上面两个查询逻辑合二为一
func GetPostListNew(c *gin.Context, p *models.ParamPostList)(data []*models.ApiPostDetail, err error){
	if p.CommunityID == 0 {
		// 原来的逻辑
		data, err = GetPostList2(c, p)
	} else {
		// 下面那个函数的逻辑
		data, err = GetCommunityPostList2(c, p)
	}
	if err != nil{
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}