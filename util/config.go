package util

import (
	"strings"

	"github.com/spf13/viper"
)

// config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type config struct {
	Domain       string `mapstructure:"DOMAIN"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DatabaseName string `mapstructure:"DATABASE_NAME"`
	PostgresUser string `mapstructure:"POSTGRES_USER"`
	PostgresPwd  string `mapstructure:"POSTGRES_PWD"`
	PostgresPort string `mapstructure:"POSTGRES_PORT"`
	MysqlUser    string `mapstructure:"MYSQL_USER"`
	MysqlPwd     string `mapstructure:"MYSQL_PWD"`
	MysqlPort    string `mapstructure:"MYSQL_PORT"`
}

// ComposeConfig 组合成需要的字段
type ComposeConfig struct {
	DBSource string
	DBDriver string
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (composeConfig ComposeConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// 读取基本配置字段
	var conf config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return
	}

	composeConfig.DBDriver = conf.DBDriver
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
	composeConfig.DBSource = dbSource.String()

	return
}
