package message

import (
	"fmt"
	"os/exec"
)

// ConvertAMR using ffmpeg command to convert .amr file to other type of audio file.
// AmrPath is the path of .arm file,as output is the file path of new audio file
// 微信语音文件默认为 .amr格式，需要转换为 .mp3格式
func ConvertAMR(amrPath string, output string) error {
	cmd := exec.Command("ffmpeg", "-i", amrPath, output)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("convert file failed:%v", err)
	}
	return nil
}
