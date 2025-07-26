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
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
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
