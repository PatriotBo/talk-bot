package openai

import "context"

// Config openai api configuration
type Config struct {
	AuthToken string `yaml:"authToken"`
	BaseURL   string `yaml:"baseURL"`
}

// Service openai api service
type Service interface {
	ChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (
		openai.ChatCompletionResponse, error)
}

type serviceImpl struct {
	cli *openai.Client
}

// New create a new openai api client
func New(cfg Config) Service {
	config := openai.DefaultConfig(cfg.AuthToken)
	config.BaseURL = cfg.BaseURL
	return &serviceImpl{
		cli: openai.NewClientWithConfig(config),
	}
}
