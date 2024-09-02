package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func ManageThumbnail(file *multipart.FileHeader) (string, error) {
	staticPath := "/web/uploads/thumbnails"
	thumbnailsDir := filepath.Join(os.Getenv("ROOT_PATH"), staticPath)
	err := os.MkdirAll(thumbnailsDir, 0755)
	if err != nil {
		return "", err
	}

	id := uuid.New()
	fileName := fmt.Sprintf("%s%s", id, filepath.Ext(file.Filename))
	filePath := filepath.Join(thumbnailsDir, fileName)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return staticPath + "/" + fileName, nil
}

func ManagePreviews(filePath string) (string, error) {
	id := uuid.New()
	staticPath := "/web/uploads/previews/" + id.String()
	videoDir := filepath.Join(os.Getenv("ROOT_PATH"), staticPath)
	err := os.MkdirAll(videoDir, 0755)
	if err != nil {
		return "", err
	}

	ffmpegPath := filepath.Join(os.Getenv("ROOT_PATH"), "ffmpeg-convert.sh")
	cmd := exec.Command("sh", ffmpegPath, filePath, videoDir)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to convert video to HLS: %v", err)
	}

	return staticPath + "/master.m3u8", nil
}

func ManageVideos(filePath, courseId string) (string, error) {
	id := uuid.New()
	staticPath := filepath.Join("/web/uploads/videos/", courseId, id.String())
	videoDir := filepath.Join(os.Getenv("ROOT_PATH"), staticPath)
	err := os.MkdirAll(videoDir, 0755)
	if err != nil {
		return "", err
	}

	ffmpegPath := filepath.Join(os.Getenv("ROOT_PATH"), "ffmpeg-convert.sh")
	cmd := exec.Command("sh", ffmpegPath, filePath, videoDir)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to convert video to HLS: %v", err)
	}

	return staticPath + "/master.m3u8", nil
}

func GetVideoLength(filePath string) (string, error) {
	cmd := exec.Command("sh", filepath.Join(os.Getenv("ROOT_PATH"), "get-video-length.sh"), filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to get video length: %v", err)
	}
	length := strings.Trim(string(output), "\n")
	return length, nil
}
