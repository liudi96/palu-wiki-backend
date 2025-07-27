package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"palu-wiki-backend/internal/models"
	"palu-wiki-backend/internal/repository"
	"palu-wiki-backend/pkg/gemini"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
)

type UpdateService struct {
	repo         *repository.GuideRepository
	geminiClient *gemini.Client
}

func NewUpdateService(repo *repository.GuideRepository, geminiClient *gemini.Client) *UpdateService {
	return &UpdateService{repo: repo, geminiClient: geminiClient}
}

// FetchSteamNews fetches news from the Palworld Steam page.
func (s *UpdateService) FetchSteamNews(url string) ([]models.OfficialUpdate, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Steam news: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch Steam news: status code %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Steam news HTML: %w", err)
	}

	var updates []models.OfficialUpdate
	doc.Find(".apphub_Card").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".apphub_CardContentNewsTitle").Text()
		content := s.Find(".apphub_CardContentNewsBody").Text()
		dateStr := s.Find(".apphub_CardContentNewsDate").Text()
		link, _ := s.Attr("href")

		// Parse date string (example format: "Posted: January 23, 2024")
		// This might need adjustment based on actual date format on the page
		parsedDate, err := time.Parse("Posted: January 2, 2006", dateStr)
		if err != nil {
			log.Printf("Failed to parse date string '%s': %v", dateStr, err)
			parsedDate = time.Now() // Fallback to current time
		}

		updates = append(updates, models.OfficialUpdate{
			Title:       title,
			Content:     content,
			PublishDate: parsedDate,
			SourceURL:   link,
			Processed:   false,
		})
	})

	return updates, nil
}

// TODO: Implement FetchTwitterUpdates for Twitter updates
// func (s *UpdateService) FetchTwitterUpdates(url string) ([]models.OfficialUpdate, error) {
// 	// Implementation for Twitter scraping (might require specific libraries or API access)
// 	return nil, fmt.Errorf("Twitter update fetching not yet implemented")
// }

// ProcessUpdates processes fetched updates and stores them in the database.
func (s *UpdateService) ProcessUpdates(updates []models.OfficialUpdate) error {
	for _, update := range updates {
		// Check if update already exists to avoid duplicates
		var existingUpdate models.OfficialUpdate
		result := repository.DB.Where("source_url = ?", update.SourceURL).First(&existingUpdate)
		if result.Error == gorm.ErrRecordNotFound {
			// Update does not exist, create it
			if err := repository.DB.Create(&update).Error; err != nil {
				log.Printf("Failed to save new official update '%s': %v", update.Title, err)
				return fmt.Errorf("failed to save new official update: %w", err)
			}
			log.Printf("Saved new official update: %s", update.Title)

			// Now, use Gemini to generate or update a guide based on this official update
			prompt := fmt.Sprintf("你现在是【幻兽帕鲁】游戏的专属攻略创作AI。你的唯一职责是根据提供的攻略主题和相关信息，创作一篇极其准确、详细、全面且最新的幻兽帕鲁游戏攻略。严禁创作任何与幻兽帕鲁无关的内容。如果主题与幻兽帕鲁无关，或者你无法创作相关攻略，请直接回答“抱歉，我只能创作幻兽帕鲁相关的攻略。”\n\n攻略主题：%s\n相关内容：%s", update.Title, update.Content)

			generatedContent, err := s.geminiClient.GenerateContent(context.Background(), prompt)
			if err != nil {
				log.Printf("Failed to generate guide content with Gemini for update '%s': %v", update.Title, err)
				// Continue processing other updates even if Gemini fails for one
			} else {
				newGuide := models.Guide{
					Title:         "【AI生成】" + update.Title,
					Content:       generatedContent,
					Category:      "官方更新",
					SourceURL:     update.SourceURL,
					IsAIGenerated: true,
					Version:       "1.0", // Or derive from update
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				if err := repository.DB.Create(&newGuide).Error; err != nil {
					log.Printf("Failed to save AI generated guide for update '%s': %v", update.Title, err)
				} else {
					log.Printf("Saved AI generated guide for update: %s", update.Title)
				}
			}

		} else if result.Error != nil {
			return fmt.Errorf("failed to check existing official update: %w", result.Error)
		} else {
			// Official update already exists, skip or update if content has changed
			log.Printf("Official update '%s' already exists, skipping.", update.Title)
		}
	}
	return nil
}
