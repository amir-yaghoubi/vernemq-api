package auth

import (
	"net/http"

	"github.com/amir-yaghoubi/vernemq-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type deleteRequest struct {
	Username string `json:"username" binding:"required" example:"USERNAME"`
}

type deleteResponse struct {
	IsDeleted bool `json:"is_deleted" example:"false"`
}

// Delete User godoc
// @Summary Delete a user
// @Description Delete a user by username
// @Tags users
// @Accept  json
// @Produce  json
// @Param payload body deleteRequest true "Delete payload"
// @Success 200 {object} deleteResponse
// @Failure 400 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Router /auth/user [delete]
func deleteUserRoute(service *auth.Service, logger *logrus.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := deleteRequest{}
		err := c.ShouldBindBodyWith(&req, binding.JSON)
		if err != nil {
			c.JSON(400, errorMessage{Error: err.Error()})
			return
		}

		isDeleted, err := service.DeleteUser(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorMessage{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, deleteResponse{IsDeleted: isDeleted})
	}
}
