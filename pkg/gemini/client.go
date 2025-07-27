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
	// Start a new chat session for each request to ensure the system instruction is applied
	// In a real application, you might want to manage chat sessions per user.
	cs := c.model.StartChat()

	// Send the system instruction as the first message in the chat history
	// This acts as a persistent prompt for the AI's role.
	// Note: The system instruction should ideally be set once per model or chat session.
	// For simplicity in this example, we're adding it with each request.
	// A more robust solution would involve managing chat history per user.

	// Send the user's prompt
	resp, err := cs.SendMessage(ctx, genai.Text("你是一个专门负责幻兽帕鲁游戏攻略的AI助手，请根据用户的问题提供准确、详细和最新的幻兽帕鲁游戏攻略。"), genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to send message to Gemini chat: %w", err)
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
