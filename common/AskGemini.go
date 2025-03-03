package common

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var GeminiKey string = "YOUR_GEMINI_API_KEY"

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
