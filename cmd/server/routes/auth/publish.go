package auth

import (
	"errors"
	"github.com/amir-yaghoubi/vernemq-auth/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type publishRequest struct {
	ClientID   string `json:"client_id"`
	Username   string `json:"username"`
	Mountpoint string `json:"mountpoint"`
	Topic      string `json:"topic"`
	Qos        uint8  `json:"qos"`
	Retain     bool   `json:"retain"`
}

type publishResponse struct {
	Result string `json:"result"`
}

func publishRoute(service *auth.Service, logger *logrus.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := publishRequest{}
		err := c.ShouldBindBodyWith(&req, binding.JSON)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		isAllowed, err := service.AuthorizePublish(req.ClientID, req.Username, req.Topic, req.Qos, req.Retain)
		if err != nil {
			if errors.Is(err, auth.ErrUnAuthorizeAccess) {
				c.JSON(200, errorResponse{Result: errorMessage{Error: "not_allowed"}})
				return
			}

			c.JSON(500, errorResponse{Result: errorMessage{Error: err.Error()}})
			return
		}

		if !isAllowed {
			c.JSON(200, errorResponse{Result: errorMessage{Error: "not_allowed"}})
			return
		}

		c.JSON(200, publishResponse{Result: "ok"})
	}
}
