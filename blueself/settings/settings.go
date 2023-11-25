package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)
var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`

	MaxBackups int    `mapstructure:"max_backups"`
}
func Init() (err error) {
	// 设置默认值
	//viper.SetDefault("fileDir", "./")
	// 读取配置文件
	//viper.SetConfigFile("./goweb/config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")              // 配置文件名称(无扩展名)
	// viper.SetConfigType("yaml")                // (配合远程配置中心来使用的)如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./blueself")             // 还可以在工作目录中查找配置
	viper.AddConfigPath("./")             // 还可以在工作目录中查找配置
	err = viper.ReadInConfig()                 // 查找并读取配置文件
	if err != nil {                            // 处理读取配置文件的错误
		fmt.Printf("ReadComfig failed, err:%v \n", err)
		return
		//panic(fmt.Errorf("fatal error config file: %s", err))

	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
	// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})
	return
}

