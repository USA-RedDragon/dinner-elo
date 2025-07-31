package apis

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/USA-RedDragon/dinner-elo/internal/utils"
)

type UserResponse struct {
	ID int64 `json:"id"`
}

func GetSSOUserID(ctx context.Context, url, token string) (int64, error) {
	resp, err := utils.HTTPRequest(ctx, http.MethodGet, url, nil, map[string]string{
		"Authorization":        "Bearer " + token,
		"Accept":               "application/vnd.github+json",
		"X-GitHub-Api-Version": "2022-11-28",
	})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	response := UserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, err
	}
	return response.ID, nil
}
