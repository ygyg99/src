package main

// 通用脚手架

import (

	//"GoProject/goweb/pkg/snowflake"

	"blueself/controller"
	"blueself/database/mysql"
	"blueself/database/redis"
	"blueself/logger"
	"blueself/pkg/snowflake"
	"blueself/routes"
	"blueself/settings"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @title bluebell_selfmade
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 8081
// @BasePath 这里写base path

func main() {
	// 1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("Init settings failed, err: %v \n", err)
		return
	}

	// 2、初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("Init logger failed, err: %v \n", err)
		return
	}
	// 将缓冲区的日志追加到日志文件中
	defer zap.L().Sync()

	// 3.初始化MySQL
	if err := mysql.Init(); err != nil {
		fmt.Printf("Init mysql failed, err: %v \n", err)
		return
	}
	defer mysql.Close()

	// 4.初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("Init redis failed, err: %v \n", err)
		return
	}

	defer redis.Close()

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("Init snowflake failed, err: %v \n", err)
		return
	}
	//初始化gin框架内置的检验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("Init validator trans failed, err: %v \n", err)
		return
	}

	// 5.注册路由
	r := routes.Setup()

	// 6.启动服务（优雅关机）

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
