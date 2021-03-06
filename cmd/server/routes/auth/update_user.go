package auth

import (
	"net/http"

	"github.com/amir-yaghoubi/vernemq-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type publishACL struct {
	Pattern       string `json:"pattern" example:"foo/bar/baz"`
	MaxQos        *uint8 `json:"max_qos,omitempty"`
	AllowedRetain *bool  `json:"allowed_retain,omitempty"`
}

type subACL struct {
	Pattern string `json:"pattern" example:"foo/bar/baz"`
	MaxQos  *uint8 `json:"max_qos,omitempty"`
}

type updateUserRequest struct {
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	ClientID   *string      `json:"client_id,omitempty"`
	Mountpoint string       `json:"mountpoint"`
	PublishACL []publishACL `json:"publish_acl,omitempty"`
	SubACL     []subACL     `json:"subscribe_acl,omitempty"`
}

type updateUserResponse struct {
	Result string `json:"result" example:"ok"`
}

// Insert or Update User
// @Summary Insert or Update User
// @Description Insert or Update User ACL
// @Tags users
// @Accept  json
// @Produce  json
// @Param payload body updateUserRequest true "Payload"
// @Success 200 {object} updateUserResponse
// @Failure 400 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Router /auth/user [post]
func updateUserRoute(service *auth.Service, logger *logrus.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := updateUserRequest{}
		err := c.ShouldBindBodyWith(&req, binding.JSON)
		if err != nil {
			c.JSON(400, errorMessage{Error: err.Error()})
			return
		}

		var subs []auth.SubACL
		var pubs []auth.PublishACL

		if req.SubACL != nil {
			subs = make([]auth.SubACL, 0, len(req.SubACL))
			for i := range req.SubACL {
				subs = append(subs, auth.SubACL{Pattern: req.SubACL[i].Pattern, MaxQos: req.SubACL[i].MaxQos})
			}
		}
		if req.PublishACL != nil {
			pubs = make([]auth.PublishACL, 0, len(req.PublishACL))
			for i := range req.PublishACL {
				pubs = append(pubs, auth.PublishACL{Pattern: req.PublishACL[i].Pattern, MaxQos: req.PublishACL[i].MaxQos, AllowedRetain: req.PublishACL[i].AllowedRetain})
			}
		}

		err = service.UpdateUser(auth.User{
			Username:   req.Username,
			Password:   req.Password,
			ClientID:   req.ClientID,
			Mountpoint: req.Mountpoint,
			SubACL:     subs,
			PublishACL: pubs,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, errorMessage{Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, updateUserResponse{Result: "ok"})
	}
}
