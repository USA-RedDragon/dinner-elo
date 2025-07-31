package store

import (
	"github.com/USA-RedDragon/dinner-elo/internal/config"
	"github.com/USA-RedDragon/dinner-elo/internal/store/gorm"
	"github.com/USA-RedDragon/dinner-elo/internal/store/models"
	"github.com/mattn/go-nulltype"
)

type Store interface {
	CreateUser(ssoUserID nulltype.NullInt64) error
	FindUserBySSOID(ssoUserID int64) (*models.User, error)
	FindUserByID(id uint) (*models.User, error)
}

func NewStore(cfg *config.Config) (Store, error) {
	switch cfg.Storage.Type {
	case config.StorageTypeSQLite:
		return gorm.NewGormStore(cfg)
	case config.StorageTypePostgres:
		return gorm.NewGormStore(cfg)
	case config.StorageTypeMySQL:
		return gorm.NewGormStore(cfg)
	default:
		return nil, config.ErrInvalidStorageType
	}
}
