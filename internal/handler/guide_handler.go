package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"palu-wiki-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GuideHandler struct {
	guideRepo *repository.GuideRepository
}

func NewGuideHandler(guideRepo *repository.GuideRepository) *GuideHandler {
	return &GuideHandler{guideRepo: guideRepo}
}

// GetGuides retrieves a list of guides for the frontend.
func (h *GuideHandler) GetGuides(c *gin.Context) {
	guides, err := h.guideRepo.GetAllGuides()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch guides: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": guides})
}

// GetGuideByID retrieves a single guide by its ID for the frontend.
func (h *GuideHandler) GetGuideByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid guide ID"})
		return
	}

	guide, err := h.guideRepo.GetGuideByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Guide not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch guide: %v", err)})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": guide})
}
