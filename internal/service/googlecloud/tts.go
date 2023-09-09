package googlecloud

import (
	"context"
	"fmt"
	"path/filepath"

	"google.golang.org/api/option"

	tts "cloud.google.com/go/texttospeech/apiv1"
	ttspb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

type Service interface {
	TextToSpeech(ctx context.Context, request TTSRequest) ([]byte, error)
}

type service struct {
}

func NewTTS() Service {
	return &service{}
}

type Language string

const (
	LanguageEnUs Language = "en-US"

	certPath = "../config/"
)

type TTSRequest struct {
	Content  string
	Language Language
}

// TextToSpeech transcripts text to speech using google cloud api.
func (s *service) TextToSpeech(ctx context.Context, request TTSRequest) ([]byte, error) {
	certFile := filepath.Join(certPath, "refined-byte-398412-c0c8f53884e1.json")
	cli, err := tts.NewClient(ctx, option.WithCredentialsFile(certFile))
	if err != nil {
		return nil, fmt.Errorf("new client failed:%v", err)
	}
	defer cli.Close()
	r := &ttspb.SynthesizeSpeechRequest{
		Input: &ttspb.SynthesisInput{
			InputSource: &ttspb.SynthesisInput_Text{Text: request.Content},
		},
		Voice: &ttspb.VoiceSelectionParams{
			LanguageCode: string(request.Language),
			SsmlGender:   ttspb.SsmlVoiceGender_FEMALE,
		},
		AudioConfig: &ttspb.AudioConfig{
			AudioEncoding: ttspb.AudioEncoding_MP3,
		},
	}

	resp, err := cli.SynthesizeSpeech(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("synthesize speech failed:%v", err)
	}

	return resp.GetAudioContent(), nil
}
