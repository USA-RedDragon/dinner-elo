package server

import (
	"net/http"

	"github.com/USA-RedDragon/dinner-elo/internal/config"
	"github.com/USA-RedDragon/dinner-elo/internal/server/controllers/auth"
	"github.com/USA-RedDragon/dinner-elo/internal/server/controllers/match"
	"github.com/gin-gonic/gin"
)

func applyRoutes(r *gin.Engine, config *config.Config) {
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"state": "OK"})
	})

	r.GET("/auth", auth.Auth)
	r.GET("/match", requireLogin(config), match.Match)
}
