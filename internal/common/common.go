package common

import "context"

type ResponseOptions struct {
	Context      context.Context
	SystemPrompt string
	Temperature  float64
	Model        string
	MaxTokens    int64
}

func NewResponseOptions() *ResponseOptions {
	return &ResponseOptions{
		Context:      context.TODO(),
		SystemPrompt: "You are a helpful assistant.",
		Temperature:  1.0,
		Model:        "",
		MaxTokens:    2048,
	}
}

type AIProvider interface {
	AvaiableModels() []string

	Response(messages []string, opt *ResponseOptions) (string, error)
}
