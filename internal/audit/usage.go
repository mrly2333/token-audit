package audit

import "encoding/json"

func ExtractUsageTotals(raw []byte) (promptTokens, completionTokens, totalTokens int64) {
	if len(raw) == 0 {
		return 0, 0, 0
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return 0, 0, 0
	}

	promptTokens = firstInt64(payload, "prompt_tokens", "input_tokens")
	completionTokens = firstInt64(payload, "completion_tokens", "output_tokens")
	totalTokens = firstInt64(payload, "total_tokens")
	if totalTokens == 0 && (promptTokens > 0 || completionTokens > 0) {
		totalTokens = promptTokens + completionTokens
	}

	return promptTokens, completionTokens, totalTokens
}

func firstInt64(payload map[string]any, keys ...string) int64 {
	for _, key := range keys {
		value, ok := payload[key]
		if !ok {
			continue
		}
		if converted, ok := int64Value(value); ok {
			return converted
		}
	}
	return 0
}

func int64Value(value any) (int64, bool) {
	switch typed := value.(type) {
	case int:
		return int64(typed), true
	case int32:
		return int64(typed), true
	case int64:
		return typed, true
	case float64:
		return int64(typed), true
	case json.Number:
		parsed, err := typed.Int64()
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
