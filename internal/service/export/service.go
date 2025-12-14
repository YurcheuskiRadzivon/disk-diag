package export

import (
	"encoding/json"
	"fmt"
)

// Service handles export operations
type Service interface {
	ExportJSON(data interface{}) ([]byte, string, error)
	ExportTXT(data interface{}, title string) ([]byte, string, error)
	ExportPDF(data interface{}, title string) ([]byte, string, error)
}

type service struct{}

// New creates a new export service
func New() Service {
	return &service{}
}

// ExportJSON exports data as JSON
func (s *service) ExportJSON(data interface{}) ([]byte, string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return bytes, "application/json", nil
}
