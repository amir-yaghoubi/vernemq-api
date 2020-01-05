package main

import (
	"fmt"
	authRoute "github.com/amir-yaghoubi/vernemq-api/cmd/server/routes/auth"
	"github.com/amir-yaghoubi/vernemq-api/internal/auth"
	lruRepo "github.com/amir-yaghoubi/vernemq-api/internal/auth/repository/lru"
	redisRepo "github.com/amir-yaghoubi/vernemq-api/internal/auth/repository/redis"
	"github.com/amir-yaghoubi/vernemq-api/internal/conf"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type PublishAcl struct {
	Pattern string
	MaxQos  uint8
}

type SubAcl struct {
	Pattern string
	MaxQos  uint8
}

type AuthModel struct {
	Username   string
	Password   string
	PublishACL []PublishAcl
	SubAcl     []SubAcl
}

func setLogLevel(logger *logrus.Logger, level string) {
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
		break
	case "info":
		logger.SetLevel(logrus.InfoLevel)
		break
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
		break
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
		break
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
		break
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	cfg := conf.ParseConfig()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	authRepo := redisRepo.New(client)
	cachedAuthRepo, _ := lruRepo.New(authRepo, cfg.AuthConfig.LruSize)
	authSrv := auth.New(cachedAuthRepo, logger)

	prometheus := ginprometheus.NewPrometheus("vmq_auth")
	prometheus.Use(router)

	router.Use(gin.Recovery())
	authRoute.Use(cfg.CacheControl, router, authSrv, logger)

	logger.WithField("port", cfg.Port).Info("starting API server")
	if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"port":  cfg.Port,
		}).Error("cannot start API")
	}
}
