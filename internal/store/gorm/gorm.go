package gorm

import (
	"fmt"

	"github.com/USA-RedDragon/dinner-elo/internal/config"
	"github.com/USA-RedDragon/dinner-elo/internal/store/models"
	"github.com/glebarez/sqlite"
	"github.com/mattn/go-nulltype"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func NewGormStore(cfg *config.Config) (*Gorm, error) {
	var dialect gorm.Dialector
	switch cfg.Storage.Type {
	case config.StorageTypeSQLite:
		dialect = sqlite.Open(cfg.Storage.DSN)
	case config.StorageTypePostgres:
		dialect = postgres.Open(cfg.Storage.DSN)
	case config.StorageTypeMySQL:
		dialect = mysql.Open(cfg.Storage.DSN)
	default:
		return nil, config.ErrInvalidStorageType
	}

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&models.Dinner{}, &models.User{}, &models.Match{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Gorm{
		db: db,
	}, nil
}

func (s *Gorm) CreateUser(ssoUserID nulltype.NullInt64) error {
	user := &models.User{
		SSOUserID: ssoUserID,
	}
	return s.db.Create(user).Error
}

func (s *Gorm) FindUserBySSOID(ssoUserID int64) (*models.User, error) {
	var user models.User
	err := s.db.Where("sso_user_id = ?", ssoUserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Gorm) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
