package repository

import (
	"palu-wiki-backend/internal/models" // Import models package

	"gorm.io/gorm"
)

type GuideRepository struct {
	db *gorm.DB
}

func NewGuideRepository(db *gorm.DB) *GuideRepository {
	return &GuideRepository{db: db}
}

// CreateGuide creates a new guide entry in the database.
func (r *GuideRepository) CreateGuide(guide *models.Guide) error {
	return r.db.Create(guide).Error
}

// GetGuideByID retrieves a guide by its ID.
func (r *GuideRepository) GetGuideByID(id uint) (*models.Guide, error) {
	var guide models.Guide
	if err := r.db.First(&guide, id).Error; err != nil {
		return nil, err
	}
	return &guide, nil
}

// GetAllGuides retrieves all guides from the database.
func (r *GuideRepository) GetAllGuides() ([]models.Guide, error) {
	var guides []models.Guide
	if err := r.db.Find(&guides).Error; err != nil {
		return nil, err
	}
	return guides, nil
}

// UpdateGuide updates an existing guide entry.
func (r *GuideRepository) UpdateGuide(guide *models.Guide) error {
	return r.db.Save(guide).Error
}

// DeleteGuide deletes a guide entry by its ID.
func (r *GuideRepository) DeleteGuide(id uint) error {
	return r.db.Delete(&models.Guide{}, id).Error
}

// SearchGuidesByTitle searches for guides by title (case-insensitive, partial match).
func (r *GuideRepository) SearchGuidesByTitle(title string) ([]models.Guide, error) {
	var guides []models.Guide
	if err := r.db.Where("title LIKE ?", "%"+title+"%").Find(&guides).Error; err != nil {
		return nil, err
	}
	return guides, nil
}

// GetGuidesByCategory retrieves guides by category.
func (r *GuideRepository) GetGuidesByCategory(category string) ([]models.Guide, error) {
	var guides []models.Guide
	if err := r.db.Where("category = ?", category).Find(&guides).Error; err != nil {
		return nil, err
	}
	return guides, nil
}
