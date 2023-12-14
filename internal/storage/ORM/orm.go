package orm

import (
	"fmt"
	"url-shortener/internal/storage"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

type Url struct {
	gorm.Model
	Url   string `gorm:"not null"`
	Alias string `gorm:"unique;not null;index:idx_alias"`
}

func New(storagePath string) (*Storage, error) {
	const op = "stotage.ORM.New"

	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.AutoMigrate(&Url{})

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "stotage.ORM.SaveURL"

	url := Url{Url: urlToSave, Alias: alias}
	result := s.db.Create(&url)
	if result.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}

	return int64(url.ID), nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "stotage.ORM.GetURL"

	var url Url

	result := s.db.First(&url, "alias = ?", alias)
	if result.Error != nil {
		return "", fmt.Errorf("%s: %w", op, result.Error)
	} else if result.RowsAffected == 0 {
		return "", storage.ErrURLNotFound
	}

	return url.Url, nil

}
