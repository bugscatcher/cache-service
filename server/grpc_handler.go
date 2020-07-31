package server

import (
	"github.com/bugscatcher/cache-service/configs"
	"github.com/go-redis/redis"
)

type GRPCHandler struct {
	Redis *redis.Client
	Conf  configs.Config
}
