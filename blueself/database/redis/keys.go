package redis

// redis key

// 命名注意使用命名空间的方式，方便查询和拆分

const(
	KeyPrefix = "blueself:"
	KeyPostTimeZSet = "post:time" //帖子以发帖时间为分数
	KeyPostScoreZset = "post:score" //帖子及投票的分数
	KeyPostVotedZSetPrefix = "post:voted:" //记录用户及投票类型  参数是postID
	
	KeyCommunitySetPrefix = "community:" //set;保存每个分区下帖子的id
)

// GetRedisKey 给RedisKey加上前缀
func GetRedisKey(key string)string  {
	return KeyPrefix+key
}