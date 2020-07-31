package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/bugscatcher/cache-service/configs"
	"github.com/bugscatcher/cache-service/pb"
	"github.com/bugscatcher/cache-service/server"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	conf, err := configs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("read config")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password:     conf.Redis.Password,
		DB:           conf.Redis.Database,
		MaxRetries:   20,
		DialTimeout:  conf.Redis.Timeout,
		ReadTimeout:  conf.Redis.Timeout,
		WriteTimeout: conf.Redis.Timeout,
		PoolSize:     conf.Redis.MaxConnections,
		MaxConnAge:   conf.Redis.MaxConnectionAge,
	})

	handler := &server.GRPCHandler{
		Redis: redisClient,
		Conf:  conf,
	}

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPCServer.Port))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to listen")
		}
		grpcServer := grpc.NewServer()
		pb.RegisterCacheServiceServer(grpcServer, handler)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("listen server")
		}
		log.Info().Msg("grpc server is started")
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	_ = <-ch
}
