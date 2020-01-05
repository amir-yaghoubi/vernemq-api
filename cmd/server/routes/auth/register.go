package auth

import (
	"errors"
	"github.com/amir-yaghoubi/vernemq-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type registerRequest struct {
	PeerAddress  string `json:"peer_addr"`
	PeerPort     string `json:"peer_port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Mountpoint   string `json:"mountpoint"`
	ClientID     string `json:"client_id"`
	CleanSession bool   `json:"clean_session"`
}

type registerResponse struct {
	Result    string         `json:"result"`
	Modifiers auth.Modifiers `json:"modifiers"`
}

func registerRoute(service *auth.Service, logger *logrus.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := registerRequest{}
		err := c.ShouldBindBodyWith(&req, binding.JSON)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		modifiers, err := service.Authenticate(req.ClientID, req.Username, req.Password)
		if err != nil {
			if errors.Is(err, auth.ErrUnAuthorizeAccess) {
				c.JSON(200, errorResponse{Result: errorMessage{Error: "not_allowed"}})
				return
			}

			c.JSON(500, errorResponse{Result: errorMessage{Error: err.Error()}})
			return
		}

		c.JSON(200, registerResponse{Result: "ok", Modifiers: *modifiers})
	}
}
