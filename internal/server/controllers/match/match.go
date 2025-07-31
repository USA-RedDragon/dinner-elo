package match

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Match(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Match found for user", "user": c.MustGet("user"), "match": "example-match-id"})
}
