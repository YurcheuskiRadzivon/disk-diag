package diagnostic

import (
	"context"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Service interface {
	AnalyzeManual(data *models.SmartInfo) (*models.HealthReport, error)
	AnalyzeAI(data *models.SmartInfo) (*models.HealthReport, error)
	Close()
}

type service struct {
	ctx          context.Context
	geminiClient *genai.Client
	geminiModel  *genai.GenerativeModel
}

func NewService(ctx context.Context, apiKey string) (Service, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	model.ResponseMIMEType = "application/json"

	srv := &service{
		ctx:          ctx,
		geminiClient: client,
		geminiModel:  model,
	}

	return srv, nil
}

func (s *service) Close() {
	if s.geminiClient != nil {
		s.geminiClient.Close()
	}
}
