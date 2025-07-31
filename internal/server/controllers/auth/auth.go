package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/USA-RedDragon/dinner-elo/internal/apis"
	"github.com/USA-RedDragon/dinner-elo/internal/config"
	"github.com/USA-RedDragon/dinner-elo/internal/store"
	"github.com/USA-RedDragon/dinner-elo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-nulltype"
	"gorm.io/gorm"
)

func Auth(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	config, ok := c.MustGet("config").(*config.Config)
	if !ok {
		slog.Error("Failed to get config from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
		return
	}

	store, ok := c.MustGet("store").(store.Store)
	if !ok {
		slog.Error("Failed to get store from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
		return
	}

	urldata := url.Values{}
	urldata.Set("code", code)
	urldata.Set("client_id", config.Auth.ClientID)
	urldata.Set("client_secret", config.Auth.ClientSecret)
	urldata.Set("grant_type", "authorization_code")
	urldata.Set("scope", "user:email")
	urldata.Set("redirect_uri", config.HTTP.URL+"/auth")

	resp, err := utils.HTTPRequest(c, http.MethodPost, config.Auth.TokenURL, strings.NewReader(urldata.Encode()), map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		slog.Error("Failed to request token", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to request token"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Debug("Token request data", "data", urldata.Encode())
		slog.Error("Failed to request token", "status", resp.Status)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to request token"})
		return
	}

	tokenResponse := &struct {
		AccessToken string `json:"access_token"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		slog.Error("Failed to decode response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
		return
	}

	id, err := apis.GetSSOUserID(c, config.Auth.UserURL, tokenResponse.AccessToken)
	if err != nil {
		slog.Error("Failed to get SSO user ID", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
		return
	}

	user, err := store.FindUserBySSOID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create user
			err := store.CreateUser(nulltype.NullInt64Of(id))
			if err != nil {
				slog.Error("Failed to create user", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
				return
			}
			user, err = store.FindUserBySSOID(id)
			if err != nil {
				slog.Error("Failed to find user", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
				return
			}
		} else {
			slog.Error("Failed to register or login user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
			return
		}
	}

	token, err := utils.GenerateJWT(config.Auth.JWTSecret, user.ID)
	if err != nil {
		slog.Error("Failed to generate JWT", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
		return
	}

	c.SetCookie("token", token, 3600, "/", "", true, true)

	c.Redirect(http.StatusFound, config.HTTP.URL)
}
