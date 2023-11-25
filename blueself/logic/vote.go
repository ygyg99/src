package logic

import (
	"blueself/database/redis"
	"blueself/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 投票功能
// 1. 用户投票的数据
// 2.

// 本项目使用简化版的投票分数
// 投一票 +432分  86400/200 -> 需要200张赞成票可以给帖子续一天 来自《redis》实战

/*
投票的几种情况
direction=1时，
	1、之前没有投过票，现在投赞成票
	2、之前投反对票，现在投赞成票

direction=0时，
	1、之前投赞成票，现在取消
	2、之前投反对票，现在取消

direction=-1时，
	1、之前没有投票，现在投反对票
	2、之前投赞成票，现在投反对票


投票的限制：
每个帖子自发表日起，一个星期内允许用户投票，超过一个星期就不允许投票了
	1、到期之后将redis中的赞成票及反对票存到MySQL中
	2、到期后删除那个KeyPostVotedZSetPF
*/

func PostVote(c *gin.Context,userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("PostVote",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", *p.Direction))
	return redis.VoteForPost(c, strconv.Itoa(int(int(userID))), p.PostID, float64(*p.Direction))

}
