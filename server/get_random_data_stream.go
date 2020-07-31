package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/bugscatcher/cache-service/pb"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

func (g GRPCHandler) GetRandomDataStream(stream pb.CacheService_GetRandomDataStreamServer) error {
	rand.Seed(time.Now().Unix())
	ctx := stream.Context()
	for i := 0; i < g.Conf.NumberOfRequests; i++ {
		go func() {
			randomUrl := g.Conf.Urls[rand.Intn(len(g.Conf.Urls))]
			status, err := g.Redis.Get(randomUrl).Result()
			if err == redis.Nil {
				log.Info().Str("url", randomUrl).Msg("adding url to cache")
				resp, err := http.Get(randomUrl)
				if err != nil {
					log.Error().Err(err).Str("url", randomUrl).Msg("can't get info for url")
					return
				}

				status = resp.Status
				//minimal supported value is 1ms
				randomDuration := time.Duration(rand.Intn(g.Conf.MaxTimeout-g.Conf.MinTimeout)+g.Conf.MinTimeout) * 1000 * 1000
				_, err = g.Redis.Set(randomUrl, status, randomDuration).Result()
				if err != nil {
					log.Error().Err(err).Str("url", randomUrl).Msg("can't set in redis")
				}
			} else if err != nil {
				log.Error().Err(err).Str("url", randomUrl).Msg("error on get from redis")
			}

			if err := stream.Send(&pb.GetRandomDataResponse{Data: fmt.Sprintf("url: %s, status: %v, time: %v", randomUrl, status, time.Now())}); err != nil {
				log.Error().Err(err).Msg("can't send message")
			}
		}()
	}

	<-ctx.Done()
	log.Info().Msg("client disconnected from server")
	return nil
}
