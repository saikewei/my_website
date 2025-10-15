package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRouters(router gin.IRouter) {
	authGroup := router.Group("/auth")
	{
		// authGroup.GET("/create-password", func(ctx *gin.Context) {
		// 	createPassword("Hyc65319436")
		// 	ctx.JSON(http.StatusOK, gin.H{"message": "Create password endpoint"})
		// })
		authGroup.PUT("/change-password", changePassword)
	}
}

func changePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be at least 6 characters long"})
		return
	}

	err := ChangePasswordStore(&req)
	if err != nil {
		switch err {
		case ErrIncoreectOldPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Old password is incorrect"})
		case ErrPasswordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "No password in database"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
