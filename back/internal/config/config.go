package config

import (
	"log"

	"github.com/spf13/viper"
)

// AppConfig 结构体用于映射配置文件中的所有项
type AppConfig struct {
	Database struct {
		Dsn string `mapstructure:"dsn"`
	} `mapstructure:"database"`
	Storage struct {
		PhotoPath string `mapstructure:"photo_path"`
	} `mapstructure:"storage"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

// C 是一个全局变量，用于在程序各处访问配置
var C AppConfig

// LoadConfig 加载位于指定路径的配置文件
func LoadConfig() {
	viper.AddConfigPath("./config") // 配置文件所在的路径
	viper.SetConfigName("config")   // 配置文件的名称 (不带扩展名)
	viper.SetConfigType("yaml")     // 配置文件的类型

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Println("Configuration loaded successfully")
}
