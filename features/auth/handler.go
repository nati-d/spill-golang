package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nati-d/spill-backend/internal/telegram"
)

func TelegramLogin(c *gin.Context) {
	initData := c.PostForm("init_data")
	if initData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing init_data"})
		return
	}

	tgUser, err := telegram.ValidateInitData(initData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid"})
		return
	}

	profile, suggestions, err := AuthService.LoginOrRegister(c.Request.Context(), tgUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":                 profile,
		"nickname_suggestions": suggestions,
	})
}

func RegisterProfileRoutes(r *gin.RouterGroup) {
	r.PATCH("/profile/nickname", func(c *gin.Context) {
		var body struct {
			Nickname string `json:"nickname" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
			return
		}

		tgUser := c.MustGet("tg_user").(telegram.TelegramUser)
		if err := AuthService.SetNickname(c.Request.Context(), tgUser.ID, body.Nickname); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}
