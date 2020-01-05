package auth

import (
	"errors"
	"github.com/amir-yaghoubi/vernemq-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type subscribeRequest struct {
	ClientID   string              `json:"client_id"`
	Username   string              `json:"username"`
	Mountpoint string              `json:"mountpoint"`
	Topics     []auth.Subscription `json:"topics"`
}

type subscribeResponse struct {
	Result string              `json:"result"`
	Topics []auth.Subscription `json:"topics"`
}

func subscribeRoute(service *auth.Service, logger *logrus.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := subscribeRequest{}
		err := c.ShouldBindBodyWith(&req, binding.JSON)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		subs, err := service.AuthorizeSubscribe(req.ClientID, req.Username, req.Topics)
		if err != nil {
			if errors.Is(err, auth.ErrUnAuthorizeAccess) {
				c.JSON(200, errorResponse{Result: errorMessage{Error: "not_allowed"}})
				return
			}

			c.JSON(500, errorResponse{Result: errorMessage{Error: err.Error()}})
			return
		}

		c.JSON(200, subscribeResponse{Result: "ok", Topics: subs})
	}
}
