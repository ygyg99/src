package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// 防止冲突，起两个别名
	ratelimit2 "github.com/juju/ratelimit"
	ratelimit1 "go.uber.org/ratelimit"
)

func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func heiHandler(c *gin.Context) {
	c.String(http.StatusOK, "ha")
}

// 基于漏桶的限流中间件(这一块有点问题)
func rateLimit1() func(ctx *gin.Context) {
	// 生成一个限流器(传入的水滴的时间间隔 单位:ns)
	rl := ratelimit1.New(100)

	return func(ctx *gin.Context) {
		//取水滴
		if rl.Take().Sub(time.Now()) > 0 {
			// time.Sleep(rl.Take().Sub(time.Now())) // 需要等这么长时间，下一滴水才会滴下来
			ctx.String(http.StatusOK, "rate limit...")
			ctx.Abort()
			return
		}
		// 能够取到水滴直接放行
		ctx.Next()
	}
}

// 基于令牌桶的限流中间件
func rateLimit2(fillInterval time.Duration, cap int64) func(ctx *gin.Context) {
	rl := ratelimit2.NewBucket(fillInterval, cap)
	return func(ctx *gin.Context) {
		// rl.Take()      //可以欠账
		if rl.TakeAvailable(1) == 1{
			ctx.Next()
			return
		}//不可以欠账，有令牌的时候才能取出来;=1能取到令牌，！=1不能取到
		ctx.String(http.StatusOK, "rate limit...")
		ctx.Abort()
	}
}

func main() {
	r := gin.Default()
	r.GET("/ping", rateLimit1(), pingHandler)
	r.GET("/hei", rateLimit2(2*time.Second, 1), heiHandler)
	r.Run(":8081")
}
