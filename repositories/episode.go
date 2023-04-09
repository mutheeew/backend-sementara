package repositories

import (
	"dumbflix/models"

	"gorm.io/gorm"
)

type EpisodeRepository interface {
	FindEpisode() ([]models.Episode, error)
	GetEpisode(ID int) (models.Episode, error)
	CreateEpisode(Episode models.Episode) (models.Episode, error)
	DeleteEpisode(Episode models.Episode) (models.Episode, error)
	UpdateEpisode(Episode models.Episode) (models.Episode, error)
}

func RepositoryEpisode(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindEpisode() ([]models.Episode, error) {
	var episode []models.Episode
	err := r.db.Preload("Film").Preload("Film_Category").Find(&episode).Error

	return episode, err
}

func (r *repository) GetEpisode(ID int) (models.Episode, error) {
	var episode models.Episode
	err := r.db.Preload("Film").First(&episode, ID).Error

	return episode, err
}

func (r *repository) CreateEpisode(episode models.Episode) (models.Episode, error) {
	err := r.db.Preload("Film").Preload("Film_Category").Create(&episode).Error

	return episode, err
}

func (r *repository) DeleteEpisode(episode models.Episode) (models.Episode, error) {
	err := r.db.Preload("Film").Preload("Film_Category").Delete(&episode).Scan(&episode).Error

	return episode, err
}

func (r *repository) UpdateEpisode(episode models.Episode) (models.Episode, error) {
	err := r.db.Preload("Film").Preload("Film.Category").Save(&episode).Error

	return episode, err
}
