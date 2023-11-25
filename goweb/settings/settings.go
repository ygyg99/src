package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	// 设置默认值
	viper.SetDefault("fileDir", "./")
	// 读取配置文件
	//viper.SetConfigFile("./goweb/config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")              // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")                // (配合远程配置中心来使用的)如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./goweb")             // 还可以在工作目录中查找配置
	viper.AddConfigPath("./")             // 还可以在工作目录中查找配置
	err = viper.ReadInConfig()                 // 查找并读取配置文件
	if err != nil {                            // 处理读取配置文件的错误
		fmt.Printf("ReadComfig failed, err:%v \n", err)
		return
		//panic(fmt.Errorf("fatal error config file: %s", err))

	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})
	return
}

