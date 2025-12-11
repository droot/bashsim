package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/droot/bashsim/session"

	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
	model  string
}

func New(ctx context.Context, model string) (*Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GOOG_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY or GOOG_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Client{
		client: client,
		model:  model,
	}, nil
}

func (c *Client) GenerateResponse(ctx context.Context, history []session.Entry, currentInput string) (string, error) {
	var contents []*genai.Content

	for _, entry := range history {
		userText := entry.Input
		if userText == "" {
			userText = " " // Should not happen for input usually, but safe guard
		}
		contents = append(contents, &genai.Content{
			Role: "user",
			Parts: []*genai.Part{{Text: userText}},
		})

		modelText := entry.Output
		if modelText == "" {
			modelText = " " // Replaced empty output with space to satisfy API
		}
		contents = append(contents, &genai.Content{
			Role: "model",
			Parts: []*genai.Part{{Text: modelText}},
		})
	}

	contents = append(contents, &genai.Content{
		Role: "user",
		Parts: []*genai.Part{{Text: currentInput}},
	})

	resp, err := c.client.Models.GenerateContent(ctx, c.model, contents, &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: "You are a bash shell simulator. You reply with standard output and standard error of the command provided. Do not use markdown blocks unless the command output itself contains them. Be concise and accurate. Do not explain your actions."}},
		},
	})
	if err != nil {
		return "", err
	}

	if resp == nil || len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response candidates")
	}

	var out string
	for _, part := range resp.Candidates[0].Content.Parts {
		out += part.Text
	}

	return out, nil
}
