package redis

import (
	"blueself/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func CreatePost(c *gin.Context, post_ID, CommunityID int64) (err error) {
	// 事务操作，利用pipeline
	pipeline := rdb.TxPipeline()

	// 帖子时间
	pipeline.ZAdd(c, GetRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: post_ID,
	})

	// 帖子分数
	pipeline.ZAdd(c, GetRedisKey(KeyPostScoreZset), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: post_ID,
	})
	// 把帖子id加到社区的set
	cKey := GetRedisKey(KeyCommunitySetPrefix+strconv.Itoa(int(CommunityID)))
	pipeline.SAdd(c, cKey, post_ID)
	_, err = pipeline.Exec(c)
	return
}

func GetIDsFormKey(c *gin.Context,key string, Page, Size int64)([]string, error){
	start := (Page - 1) * Size
	end := start + Size - 1
	// 3. ZREVRANGE查询
	return rdb.ZRevRange(c, key, start, end).Result()
}

func GetPostIDsInOrder(c *gin.Context, p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 1. 根据用户请求中携带的order的参数确定要查询的redis key
	key := GetRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScoreZset)
	}
	// 2. 确定查询的索引起始点
	return GetIDsFormKey(c, key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(c *gin.Context, ids []string) (data []int64, err error) {
	/* 方法一
	data = make([]int64, 0, len(ids))
	for _, id := range ids{
		key := GetRedisKey(KeyPostVotedZSetPrefix+id)
		// 查找key中分数是1的元素数量->统计每篇帖子的赞成票的数量
		v := rdb.ZCount(c, key,"1","1").Val()
		data = append(data, v)
	}
	return
	*/

	// 利用pipeline形式 一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := GetRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(c, key, "1", "1")
	}
	cmders, err := pipeline.Exec(c)
	if err != nil {
		return
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区根据ids查询每篇帖子的投赞成票的数据
func GetCommunityPostIDsInOrder(c *gin.Context, p *models.ParamPostList) ([]string, error) {
	// 使用zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset，
	// 针对新的zset按之前的逻辑取数据

	// 社区的key
	cKey := GetRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少zinterstore执行的次数(不能每次都执行)
	key := p.Order + strconv.Itoa(int(p.CommunityID))

	// 这个分支只有当pipeline分支过期了之后才会执行，否则不会执行
	if rdb.Exists(c, key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		
		// key是拼接之后存的地址, Keys中{cKey, p.Orde}是放需要合并的key
		pipeline.ZInterStore(c, key, &redis.ZStore{
			Keys:      []string{cKey, p.Order},
			Aggregate: "MAX",
		})	//zinterstore计算
		pipeline.Expire(c, key, 60*time.Second)
		_,err := pipeline.Exec(c)
		if err != nil{
			return nil, err
		}
	}
	
	// 存在的话直接根据key查询ids
	// 从redis获取id
	return GetIDsFormKey(c, key, p.Page, p.Size)

	// 下面这部分就封装到上面的函数里面去了
	/*
	// 1. 根据用户请求中携带的order的参数确定要查询的redis key
	key := GetRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScoreZset)
	}
	// 2. 确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// 3. ZREVRANGE查询
	return rdb.ZRevRange(c, key, start, end).Result()
	*/
}
