package handler

import (
	"log"
	"net/http"

	"palu-wiki-backend/pkg/gemini"

	"github.com/gin-gonic/gin"
)

type GeminiHandler struct {
	geminiClient *gemini.Client
}

func NewGeminiHandler(client *gemini.Client) *GeminiHandler {
	return &GeminiHandler{geminiClient: client}
}

type GenerateContentRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type GenerateContentResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func (h *GeminiHandler) GenerateContent(c *gin.Context) {
	var req GenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenerateContentResponse{Error: err.Error()})
		return
	}

	content, err := h.geminiClient.GenerateContent(c.Request.Context(), req.Prompt)
	if err != nil {
		log.Printf("Gemini API call failed: %v", err) // Add detailed error logging
		c.JSON(http.StatusInternalServerError, GenerateContentResponse{Error: "AI服务内部错误，请检查后端日志。" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenerateContentResponse{Content: content})
}
