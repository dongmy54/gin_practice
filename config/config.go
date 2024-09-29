package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() {
	// 配置名称
	viper.SetConfigName("settings")
	// 设置配置文件类型 (YAML)
	viper.SetConfigType("yaml")
	// 配置路径 注意这里的路径执行时候的相对路径 这个是在main.go执行的 所以相对路径是相对于main.go的
	viper.AddConfigPath("config/")

	// 设置默认值
	viper.SetDefault("app.debug", true)
	viper.SetDefault("app.name", "gin_practice")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	} else {
		// 打印配置信息
		debug := viper.GetBool("app.debug")     // 获取bool类型配置
		app_name := viper.GetString("app.name") // 获取string类型配置
		fmt.Println("debug:", debug)
		fmt.Println("app_name:", app_name)
	}

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	// 有时候可以反向写出配置文件
	viper.SafeWriteConfigAs("config/write_config.yaml")
}
