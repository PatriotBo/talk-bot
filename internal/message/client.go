package message

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/util"
)

type FileType string

const (
	AMR FileType = "amr"
	MP3 FileType = "mp3"
)

type Message struct {
	oa     *officialaccount.OfficialAccount
	client *http.Client
}

func NewMessage(oa *officialaccount.OfficialAccount) *Message {
	return &Message{
		oa:     oa,
		client: http.DefaultClient,
	}
}

func (m *Message) GetVoice(ctx context.Context, msg *message.MixMessage) (string, error) {
	if msg.MsgType != message.MsgTypeVoice {
		return "", errors.New("voice message required")
	}

	url, err := m.oa.GetMaterial().GetMediaURL(msg.MediaID)
	if err != nil {
		return "", fmt.Errorf("get media url failed:%v", err)
	}

	by, err := util.HTTPGetContext(ctx, url)
	if err != nil {
		return "", fmt.Errorf("get voice failed:%v", err)
	}

	amrFilename := VoicePath(msg.MediaID, "user", AMR)
	if err = os.WriteFile(amrFilename, by, 0644); err != nil {
		return amrFilename, fmt.Errorf("write voice file:%v", err)
	}

	mp3Filename := VoicePath(msg.MediaID, "user", MP3)
	return mp3Filename, ConvertAMR(amrFilename, mp3Filename)
}

func (m *Message) UploadVoice(_ context.Context, filename string) (string, error) {
	media, err := m.oa.GetMaterial().MediaUpload(material.MediaTypeVoice, filename)
	if err != nil {
		return "", err
	}
	return media.MediaID, nil
}

func VoicePath(mediaID, prefix string, ft FileType) string {
	return filepath.Join("../voice/",
		fmt.Sprintf("%s_%s_%d.%s", prefix, mediaID, time.Now().Unix(), ft))
}
