package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nati-d/spill-backend/internal/telegram"
	"net/http"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		initData := c.GetHeader("X-Telegram-Init-Data")
		if initData == "" {
			initData = c.PostForm("initData")
		}

		if initData == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, err := telegram.ValidateInitData(initData)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("tg_user", user)
		c.Next()
	}
}