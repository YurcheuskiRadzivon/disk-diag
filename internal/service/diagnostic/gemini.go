package diagnostic

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
	"github.com/google/generative-ai-go/genai"
)

func (s *service) AnalyzeAI(data *models.SmartInfo) (*models.HealthReport, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга данных: %v", err)
	}

	prompt := fmt.Sprintf(`
	Ты инженер по восстановлению данных и диагностике SSD/NVMe.
	Проанализируй этот JSON лог от утилиты smartctl.
	
	Данные диска:
	%s
	
	Твоя задача:
	1. Проверить критические параметры (NVMe Critical Warning, Media Errors, Temperature, Spare, Percentage Used).
	2. Дать понятное объяснение для пользователя.
	3. Оценить здоровье от 0 до 100.
	
	Верни ответ СТРОГО в формате JSON, который соответствует этой схеме (не добавляй markdown, просто верни JSON):
	{
		"method": "Gemini AI",
		"status": "Healthy" (или "Warning", "Critical"),
		"health_score": 95,
		"summary": "Текст вывода на русском языке (максимум 2 предложения)",
		"problems": ["Текст проблемы 1", "Текст проблемы 2"],
		"metrics": {
			"temperature_c": int,
			"life_remaining_percent": int,
			"data_written_tb": float,
			"power_on_hours": int,
			"media_errors": int,
			"unsafe_shutdowns": int
		}
	}
	Если проблем нет, массив "problems" должен быть пустым.
	`, string(jsonData))

	resp, err := s.geminiModel.GenerateContent(s.ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к Gemini: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("пустой ответ от Gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	txt, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("ответ не является текстом")
	}

	cleanJSON := strings.TrimSpace(string(txt))
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimPrefix(cleanJSON, "```")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")

	var report models.HealthReport
	if err := json.Unmarshal([]byte(cleanJSON), &report); err != nil {
		return nil, fmt.Errorf("AI вернул некорректный JSON: %v. Ответ: %s", err, cleanJSON)
	}

	report.Method = "Gemini AI"

	return &report, nil
}
