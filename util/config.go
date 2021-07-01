package util

import (
	"strings"

	// 文件：支持多种类型的文件，例如 JSON、TOML、YAML、ENV 或 INI；还支持环境变量以及远程读取...
	"github.com/spf13/viper"
)

// config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type baseConfig struct {
	Domain         string `mapstructure:"DOMAIN"`
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DatabaseName   string `mapstructure:"DATABASE_NAME"`
	PostgresUser   string `mapstructure:"POSTGRES_USER"`
	PostgresPwd    string `mapstructure:"POSTGRES_PWD"`
	PostgresPort   string `mapstructure:"POSTGRES_PORT"`
	MysqlUser      string `mapstructure:"MYSQL_USER"`
	MysqlPwd       string `mapstructure:"MYSQL_PWD"`
	MysqlPort      string `mapstructure:"MYSQL_PORT"`
	Server_Address string `mapstructure:"SERVER_ADDRESS"`
}

// Config 组合成需要的字段
type Config struct {
	DBSource      string
	DBDriver      string
	ServerAddress string
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// 告诉 viper 配置文件的位置 (调用此函数时的相对位置)
	viper.AddConfigPath(path)
	// 查找具有此特定名称的配置文件
	viper.SetConfigName("app")
	// 配置文件的类型
	viper.SetConfigType("env")

	// 希望 viper 从环境变量中读取值 并调用 viper.AutomaticEnv() 告诉 viper 自动覆盖值
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// 读取基本配置字段
	var conf baseConfig
	// json转结构体
	err = viper.Unmarshal(&conf)
	if err != nil {
		return
	}

	config.DBDriver = conf.DBDriver
	config.ServerAddress = conf.Server_Address
	// 官方推荐拼接字符串使用buffer.Builder
	var dbSource strings.Builder
	dbSource.WriteString("postgresql://")
	dbSource.WriteString(conf.PostgresUser)
	dbSource.WriteString(":")
	dbSource.WriteString(conf.PostgresPwd)
	dbSource.WriteString("@")
	dbSource.WriteString(conf.Domain)
	dbSource.WriteString(":")
	dbSource.WriteString(conf.PostgresPort)
	dbSource.WriteString("/")
	dbSource.WriteString(conf.DatabaseName)
	dbSource.WriteString("?sslmode=disable")
	config.DBSource = dbSource.String()

	return
}
