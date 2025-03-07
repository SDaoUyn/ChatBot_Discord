package common

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"io/ioutil"
	"path/filepath"
)

var GeminiKey string = "YOUR_GEMINI_KEY"

func AskGemini(question string) (string, error) {
	// Tạo client với API key
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiKey))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo client Gemini: %v", err)
	}
	defer client.Close()

	// Tạo model Gemini
	model := client.GenerativeModel("gemini-2.0-flash")

	// Cấu hình generative configs
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(500)

	// Tạo và gửi prompt
	prompt := genai.Text(question)
	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("lỗi gọi API Gemini: %v", err)
	}

	// Kiểm tra và lấy phản hồi
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("không nhận được phản hồi từ Gemini")
	}

	// Trích xuất text từ phản hồi
	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("không thể đọc nội dung phản hồi")
	}

	return string(text), nil
}

func getImageFormat(imagePath string) string {
	ext := filepath.Ext(imagePath)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg" // Mặc định là JPEG nếu không xác định được
	}
}

func AskGeminiWithImage(question string, imagePath string) (string, error) {
	// Tạo client với API key
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiKey))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo client Gemini: %v", err)
	}
	defer client.Close()

	// Tạo model Gemini
	model := client.GenerativeModel("gemini-2.0-flash")

	// Cấu hình generative configs
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(500)

	// Đọc nội dung file hình ảnh
	imgBytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("lỗi đọc file hình ảnh: %v", err)
	}

	// Xác định format hình ảnh
	imageFormat := getImageFormat(imagePath)

	// Tạo prompt kết hợp văn bản và hình ảnh
	prompt := []genai.Part{
		genai.Text(question),
		genai.ImageData(imageFormat, imgBytes),
	}

	// Gửi prompt có hình ảnh
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", fmt.Errorf("lỗi gọi API Gemini: %v", err)
	}

	// Kiểm tra và lấy phản hồi
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("không nhận được phản hồi từ Gemini")
	}

	// Trích xuất text từ phản hồi
	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("không thể đọc nội dung phản hồi")
	}

	return string(text), nil
}

// Hàm hỗ trợ xử lý nhiều hình ảnh
func AskGeminiMultipleImages(question string, imagePaths []string) (string, error) {
	// Tạo client với API key
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiKey))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo client Gemini: %v", err)
	}
	defer client.Close()

	// Tạo model Gemini
	model := client.GenerativeModel("gemini-2.0-flash")

	// Cấu hình generative configs
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(500)

	// Tạo prompt chứa văn bản và nhiều hình ảnh
	prompt := []genai.Part{
		genai.Text(question),
	}

	// Thêm các hình ảnh vào prompt
	for _, imagePath := range imagePaths {
		// Đọc nội dung file hình ảnh
		imgBytes, err := ioutil.ReadFile(imagePath)
		if err != nil {
			return "", fmt.Errorf("lỗi đọc file hình ảnh %s: %v", imagePath, err)
		}

		// Xác định format hình ảnh
		imageFormat := getImageFormat(imagePath)

		prompt = append(prompt, genai.ImageData(imageFormat, imgBytes))
	}

	// Gửi prompt có nhiều hình ảnh
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", fmt.Errorf("lỗi gọi API Gemini: %v", err)
	}

	// Kiểm tra và lấy phản hồi
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("không nhận được phản hồi từ Gemini")
	}

	// Trích xuất text từ phản hồi
	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("không thể đọc nội dung phản hồi")
	}

	return string(text), nil
}
