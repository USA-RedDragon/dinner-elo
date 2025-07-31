package server

import (
	"log/slog"
	"net/http"

	"github.com/USA-RedDragon/dinner-elo/internal/config"
	"github.com/USA-RedDragon/dinner-elo/internal/store"
	"github.com/USA-RedDragon/dinner-elo/internal/utils"
	"github.com/gin-gonic/gin"
)

func applyMiddleware(r *gin.Engine, config *config.Config) {
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.TrustedPlatform = "X-Real-IP"

	err := r.SetTrustedProxies(config.HTTP.TrustedProxies)
	if err != nil {
		slog.Error("Failed to set trusted proxies", "error", err.Error())
	}
}

func requireLogin(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		possibleTokenCookies := c.Request.CookiesNamed("token")
		if len(possibleTokenCookies) == 0 {
			// No token found in cookies
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		if len(possibleTokenCookies) > 1 {
			// Multiple tokens found, this is unexpected
			slog.Error("Multiple tokens found in cookies, this is unexpected", "count", len(possibleTokenCookies))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		tokenCookie := possibleTokenCookies[0]
		if tokenCookie.Value == "" {
			// Token is empty, unauthorized
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		uid, err := utils.VerifyJWT(config.Auth.JWTSecret, tokenCookie.Value)
		if err != nil {
			// Token verification failed, unauthorized
			slog.Error("Failed to verify JWT", "error", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		// Token is valid, check if user exists in the database
		store, ok := c.MustGet("store").(store.Store)
		if !ok {
			slog.Error("Failed to get store from context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		dbUser, err := store.FindUserByID(uid)
		if err != nil {
			slog.Error("Failed to find user by ID", "error", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		if dbUser != nil && dbUser.ID == uid {
			c.Set("user", dbUser)
			c.Next()
			return
		}

		// JWT is valid but user does not exist in the database
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
	}
}
