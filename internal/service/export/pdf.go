package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// ExportPDF exports data as PDF
func (s *service) ExportPDF(data interface{}, title string) ([]byte, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, title)
	pdf.Ln(15)

	// Content
	pdf.SetFont("Arial", "", 10)

	// Convert data to readable format
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to process data: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, "", fmt.Errorf("failed to parse data: %w", err)
	}

	writePDFContent(pdf, result, 0)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), "application/pdf", nil
}

func writePDFContent(pdf *gofpdf.Fpdf, data interface{}, indent int) {
	leftMargin := 10.0 + float64(indent)*5

	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			pdf.SetX(leftMargin)
			if isScalarPDF(val) {
				text := fmt.Sprintf("%s: %v", key, val)
				// Truncate long text
				if len(text) > 100 {
					text = text[:100] + "..."
				}
				pdf.MultiCell(180-float64(indent)*5, 5, text, "", "", false)
			} else {
				pdf.SetFont("Arial", "B", 10)
				pdf.Cell(50, 5, key+":")
				pdf.Ln(6)
				pdf.SetFont("Arial", "", 10)
				writePDFContent(pdf, val, indent+1)
			}
		}
	case []interface{}:
		for i, item := range v {
			pdf.SetX(leftMargin)
			pdf.SetFont("Arial", "I", 10)
			pdf.Cell(30, 5, fmt.Sprintf("Item %d:", i+1))
			pdf.Ln(6)
			pdf.SetFont("Arial", "", 10)
			writePDFContent(pdf, item, indent+1)
			pdf.Ln(3)
		}
	default:
		pdf.SetX(leftMargin)
		text := fmt.Sprintf("%v", v)
		if len(text) > 100 {
			text = text[:100] + "..."
		}
		// Handle non-ASCII characters
		text = strings.Map(func(r rune) rune {
			if r > 127 {
				return '?'
			}
			return r
		}, text)
		pdf.MultiCell(180-float64(indent)*5, 5, text, "", "", false)
	}
}

func isScalarPDF(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return false
	default:
		return true
	}
}
