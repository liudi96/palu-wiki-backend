package gemini

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	model *genai.GenerativeModel
}

func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()
	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	model := genaiClient.GenerativeModel("gemini-2.5-flash") // Using gemini-2.5-flash for text generation
	return &Client{model: model}, nil
}

func (c *Client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// Combine system instruction and user prompt with stronger emphasis and constraints
	fullPrompt := fmt.Sprintf("你现在是【幻兽帕鲁】游戏的专属攻略AI助手。你的唯一职责是根据用户提出的问题，提供极其准确、详细、全面且最新的幻兽帕鲁游戏攻略。严禁回答任何与幻兽帕鲁无关的内容。如果用户的问题与幻兽帕鲁无关，或者你无法提供相关攻略，请直接回答“抱歉，我只能提供幻兽帕鲁相关的攻略信息。”\n\n用户问题：%s", prompt)

	resp, err := c.model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated from Gemini API")
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			result += string(txt)
		}
	}
	return result, nil
}

func (c *Client) Close() {
	if c.model != nil {
		// In a real application, you might want to close the underlying client
		// if it has a Close method. For genai.GenerativeModel, there's no direct Close.
		// The underlying genai.Client might have one, but it's not exposed here.
		// For now, we just log.
		log.Println("Gemini client closed (if applicable).")
	}
}
