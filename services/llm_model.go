package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LLMResponse struct {
	Description string `json:"description"` // Thay đổi theo cấu trúc JSON thực tế từ LLM
}

func AnalyzeImageWithLLM(pixels []byte) (string, error) {
	// Tạo yêu cầu đến LLM API (OpenAI ví dụ)
	url := "https://api.openai.com/v1/chat/completions" // Thay thế bằng URL thực tế của mô hình LLM

	// Chuyển đổi pixels thành yêu cầu JSON
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo", // Chọn mô hình bạn muốn sử dụng
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": string(pixels),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Gửi yêu cầu POST với API key
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer YOUR_SESSION_CODE") // Thay thế bằng API key của bạn
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to LLM API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API returned non-200 status: %s", resp.Status)
	}

	// Giải mã phản hồi JSON
	var llmResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&llmResponse); err != nil {
		return "", fmt.Errorf("failed to decode LLM response: %v", err)
	}

	return llmResponse.Choices[0].Message.Content, nil
}
