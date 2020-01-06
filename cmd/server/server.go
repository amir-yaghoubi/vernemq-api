package main

import (
	"fmt"

	_ "github.com/amir-yaghoubi/vernemq-api/cmd/server/docs"
	authRoute "github.com/amir-yaghoubi/vernemq-api/cmd/server/routes/auth"
	"github.com/amir-yaghoubi/vernemq-api/internal/auth"
	lruRepo "github.com/amir-yaghoubi/vernemq-api/internal/auth/repository/lru"
	redisRepo "github.com/amir-yaghoubi/vernemq-api/internal/auth/repository/redis"
	"github.com/amir-yaghoubi/vernemq-api/internal/conf"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

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

// @title Vernemq API
// @version 1.0
// @description API wrapper around vernemq broker

// @contact.name Amirhossein Yaghoubi
// @contact.url https://yaghoubi.dev
// @contact.email amir.yaghoubi.dev@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9595
// @BasePath /api/v1/
func main() {
	cfg := conf.ParseConfig()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	authRepo := redisRepo.New(client)
	cachedAuthRepo, _ := lruRepo.New(authRepo, cfg.AuthConfig.LruSize)
	authSrv, err := auth.New(cfg.AuthConfig.LruSize, cachedAuthRepo, logger)
	if err != nil {
		logger.Fatal(err)
	}

	prometheus := ginprometheus.NewPrometheus("vmq_auth")
	prometheus.Use(router)

	router.Use(gin.Recovery())

	apiV1 := router.Group("api/v1/")

	authRoute.Use(cfg.CacheControl, apiV1, authSrv, logger)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.WithField("port", cfg.Port).Info("starting API server")
	if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"port":  cfg.Port,
		}).Error("cannot start API")
	}
}
