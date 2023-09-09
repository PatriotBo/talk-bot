package openai

import (
	"context"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
)

const (
	transcriptionModel = "whisper-1" // ID of the model to use. Only whisper-1 is currently available.
	languageEN         = "en"        // default language when using transcription
)

// AudioRequest represents a request structure for audio API.
type AudioRequest struct {
	Reader   io.Reader
	FilePath string
	Prompt   string
}

// SpeechToText calls openai API to transcript audio to text. It only supports english for now.
func (s *serviceImpl) SpeechToText(ctx context.Context, request *AudioRequest) (string, error) {
	req := openai.AudioRequest{
		Model:    transcriptionModel,
		Reader:   request.Reader,
		FilePath: request.FilePath,
		Language: languageEN,
		Prompt:   request.Prompt,
		Format:   openai.AudioResponseFormatText,
	}
	resp, err := s.cli.CreateTranscription(ctx, req)
	if err != nil {
		return "", err
	}

	fmt.Printf("SpeechToText resp:%+v \n", resp)
	return resp.Text, nil
}
