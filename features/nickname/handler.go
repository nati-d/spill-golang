// features/nickname/handler.go
package nickname

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register these two routes on your authenticated group
func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/nickname/suggestions", getSuggestions) // → GET  /nickname/suggestions
	r.POST("/nickname/reserve", reserveNickname)   // → POST /nickname/reserve
}

func getSuggestions(c *gin.Context) {
	suggestions, err := GenerateThree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate suggestions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
	})
}

func reserveNickname(c *gin.Context) {
	var req struct {
		Nickname string `json:"nickname" binding:"required,min=3,max=30"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nickname"})
		return
	}

	if Reserve(req.Nickname) {
		// Optionally save to profile here or in auth service
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "nickname reserved!",
		})
	} else {
		c.JSON(http.StatusConflict, gin.H{
			"error": "nickname already taken",
		})
	}
}
