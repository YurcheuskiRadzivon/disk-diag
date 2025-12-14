package export

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ExportTXT exports data as formatted text
func (s *service) ExportTXT(data interface{}, title string) ([]byte, string, error) {
	var sb strings.Builder

	sb.WriteString("=" + strings.Repeat("=", 60) + "\n")
	sb.WriteString(fmt.Sprintf("  %s\n", title))
	sb.WriteString("=" + strings.Repeat("=", 60) + "\n\n")

	// Convert data to map for iteration
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to process data: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, "", fmt.Errorf("failed to parse data: %w", err)
	}

	formatValue(&sb, result, 0)

	return []byte(sb.String()), "text/plain", nil
}

func formatValue(sb *strings.Builder, data interface{}, indent int) {
	prefix := strings.Repeat("  ", indent)

	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			sb.WriteString(fmt.Sprintf("%s%s: ", prefix, key))
			if isScalar(val) {
				sb.WriteString(fmt.Sprintf("%v\n", val))
			} else {
				sb.WriteString("\n")
				formatValue(sb, val, indent+1)
			}
		}
	case []interface{}:
		for i, item := range v {
			sb.WriteString(fmt.Sprintf("%s[%d]:\n", prefix, i))
			formatValue(sb, item, indent+1)
		}
	default:
		sb.WriteString(fmt.Sprintf("%s%v\n", prefix, v))
	}
}

func isScalar(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return false
	default:
		return true
	}
}
