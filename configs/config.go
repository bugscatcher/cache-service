package configs

import (
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Port             int           `mapstructure:"port"`
	Host             string        `mapstructure:"host"`
	Timeout          time.Duration `mapstructure:"timout"`
	Password         string        `mapstructure:"password"`
	Database         int           `mapstructure:"database"`
	MaxConnections   int           `mapstructure:"max_connections"`
	MaxConnectionAge time.Duration `mapstructure:"max_connection_age"`
}

type GRPCServerConfig struct {
	Port int `mapstructure:"port"`
}

type Config struct {
	Redis            RedisConfig      `mapstructure:"redis"`
	GRPCServer       GRPCServerConfig `mapstructure:"grpc_server"`
	Urls             []string         `mapstructure:"urls"`
	MinTimeout       int              `mapstructure:"mintimeout"`
	MaxTimeout       int              `mapstructure:"maxtimeout"`
	NumberOfRequests int              `mapstructure:"numberofrequests"`
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	_ = err

	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.timeout", time.Minute)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.max_connections", 10*runtime.NumCPU())
	viper.SetDefault("redis.max_connection_age", time.Minute)

	viper.SetDefault("grpc_server.port", 6565)

	viper.SetDefault("urls", []string{"https://golang.org", "https://www.google.com", "https://www.bbc.co.uk", "https://www.github.com", "https://www.gitlab.com", "https://www.duckduckgo.com", "https://www.atlasian.com", "https://www.twitter.com", "https://www.facebook.com"})
	viper.SetDefault("mintimeout", 10)
	viper.SetDefault("maxtimeout", 100)
	viper.SetDefault("numberofrequests", 3)
}

func New() (Config, error) {
	var conf Config
	err := viper.Unmarshal(&conf)
	return conf, err
}
