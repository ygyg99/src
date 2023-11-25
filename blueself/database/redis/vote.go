package redis

import (
	"errors"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeted = errors.New("投票重复")
)

func VoteForPost(c *gin.Context, userID, postID string, value float64) error {
	// 1. 判断投票的限制
	// 去redis取帖子发布时间
	postTime := rdb.ZScore(c, GetRedisKey(KeyPostTimeZSet), postID).Val()

	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2和3需要放到一个pipeline事务里面去
	// 2. 更新分数
	// 先当前用户给当前帖子的投票记录
	ov := rdb.ZScore(c, GetRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	
	// 更新部分：如果这一次投票的值与之前保存的值一致，则提示不允许重复投票
	if value == ov{
		return ErrVoteRepeted
	}
	
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值

	// 定义一个pipeline
	pipeline := rdb.TxPipeline()

	pipeline.ZIncrBy(c, GetRedisKey(KeyPostScoreZset), op*diff*scorePerVote, postID).Result()

	// 3. 记录用户为该帖子投过票的数据
	if value == 0 {
		pipeline.ZRem(c, GetRedisKey(KeyPostVotedZSetPrefix+postID), userID).Result()
	} else {
		z := &redis.Z{
			Score:  value,
			Member: userID,
		}
		pipeline.ZAdd(c, GetRedisKey(KeyPostVotedZSetPrefix+postID), z).Result()
	}
	_, err := pipeline.Exec(c)
	return err

}
