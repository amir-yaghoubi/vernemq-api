package auth

import (
	"fmt"
	"github.com/amir-yaghoubi/vernemq-api/internal/auth"
	"github.com/amir-yaghoubi/vernemq-api/internal/conf"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type errorMessage struct {
	Error string `json:"error"`
}

type errorResponse struct {
	Result errorMessage `json:"result"`
}

func cacheControl(ttl time.Duration) func(c *gin.Context) {
	age := fmt.Sprintf("max-age=%d", int(ttl.Seconds()))
	return func(c *gin.Context) {
		c.Header("Cache-Control", age)
	}
}

// Use register auth routes on gin router
func Use(cacheConfig conf.CacheControl, e *gin.Engine, service *auth.Service, logger *logrus.Logger) {
	grp := e.Group("auth/")
	grp.POST("user/", updateUserRoute(service, logger))
	grp.DELETE("user/", deleteUserRoute(service, logger))
	grp.POST("register/", cacheControl(cacheConfig.Register), registerRoute(service, logger))
	grp.POST("publish/", cacheControl(cacheConfig.Publish), publishRoute(service, logger))
	grp.POST("subscribe/", cacheControl(cacheConfig.Subscribe), subscribeRoute(service, logger))
}
