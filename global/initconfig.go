package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig(env string) {

	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s_pro.yaml", configFilePrefix)
	if env == "dev" {
		configFileName = fmt.Sprintf("%s_dev.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&ServerConfig); err != nil {
		panic(err)
	}
	//配置热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&ServerConfig); err != nil {
			panic(fmt.Errorf("unmarshal conf failed err:%s\n", err))
		}
	})
	return
}
