package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	dataURLPattern = regexp.MustCompile(`data:image/[A-Za-z0-9.+-]+;base64,[A-Za-z0-9+/=_-]+`)
	keyedBase64Pattern = regexp.MustCompile(`("(?:[^"\\]|\\.)*(?:b64_json|image_base64|image_data|base64|b64)(?:[^"\\]|\\.)*"\s*:\s*")([A-Za-z0-9+/=_\-\r\n\t ]{256,})("?)`)
)

func RedactTextContent(text string) string {
	if strings.TrimSpace(text) == "" {
		return text
	}

	changed := false
	redacted := dataURLPattern.ReplaceAllStringFunc(text, func(match string) string {
		changed = true
		return redactDataURLValue(match)
	})

	if looksLikeBase64Payload(strings.TrimSpace(redacted)) {
		return fmt.Sprintf("[omitted image base64, %d chars]", len(compactBase64(strings.TrimSpace(redacted))))
	}

	if !changed {
		return text
	}

	return redacted
}

func RedactCapturedBody(body []byte, contentType string) []byte {
	if len(body) == 0 {
		return nil
	}

	switch {
	case LooksLikeJSON(contentType, body):
		return RedactCapturedJSON(body)
	case IsSSE(contentType):
		return redactSSEBody(body)
	default:
		return redactTextFallback(body)
	}
}

func RedactCapturedJSON(raw []byte) []byte {
	if len(raw) == 0 {
		return nil
	}

	var payload any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return redactTextFallback(raw)
	}

	redacted, changed := redactValue(payload, "")
	if !changed {
		return bytes.Clone(raw)
	}

	encoded, err := json.Marshal(redacted)
	if err != nil {
		return bytes.Clone(raw)
	}
	return encoded
}

func redactSSEBody(raw []byte) []byte {
	if len(raw) == 0 {
		return nil
	}

	lines := strings.Split(string(raw), "\n")
	changed := false

	for i, line := range lines {
		suffix := ""
		if strings.HasSuffix(line, "\r") {
			suffix = "\r"
			line = strings.TrimSuffix(line, "\r")
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" || payload == "[DONE]" {
			continue
		}

		redacted := RedactCapturedJSON([]byte(payload))
		if bytes.Equal(redacted, []byte(payload)) {
			continue
		}

		lines[i] = "data: " + string(redacted) + suffix
		changed = true
	}

	if !changed {
		return bytes.Clone(raw)
	}

	return []byte(strings.Join(lines, "\n"))
}

func redactValue(value any, key string) (any, bool) {
	switch typed := value.(type) {
	case map[string]any:
		changed := false
		for childKey, childValue := range typed {
			redacted, childChanged := redactValue(childValue, childKey)
			typed[childKey] = redacted
			changed = changed || childChanged
		}
		return typed, changed
	case []any:
		changed := false
		for i, childValue := range typed {
			redacted, childChanged := redactValue(childValue, key)
			typed[i] = redacted
			changed = changed || childChanged
		}
		return typed, changed
	case string:
		redacted, changed := redactStringValue(key, typed)
		if changed {
			return redacted, true
		}
		return typed, false
	default:
		return value, false
	}
}

func redactStringValue(key, value string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return value, false
	}

	lowerKey := strings.ToLower(strings.TrimSpace(key))
	lowerValue := strings.ToLower(trimmed)
	if strings.Contains(lowerValue, "data:image/") {
		redacted := RedactTextContent(trimmed)
		if redacted != trimmed {
			return redacted, true
		}
	}

	if !shouldRedactImageField(lowerKey) {
		return value, false
	}
	if !looksLikeBase64Payload(trimmed) {
		return value, false
	}

	return fmt.Sprintf("[omitted image base64, %d chars]", len(compactBase64(trimmed))), true
}

func shouldRedactImageField(key string) bool {
	switch {
	case key == "b64_json":
		return true
	case key == "image_base64":
		return true
	case key == "image_data":
		return true
	case strings.Contains(key, "base64"):
		return true
	case strings.Contains(key, "b64"):
		return true
	case strings.Contains(key, "image"):
		return true
	default:
		return false
	}
}

func redactDataURLValue(value string) string {
	mime := "image"
	payloadLen := 0

	if comma := strings.Index(value, ","); comma >= 0 {
		header := strings.TrimPrefix(value[:comma], "data:")
		payload := compactBase64(value[comma+1:])
		if semi := strings.Index(header, ";"); semi >= 0 {
			header = header[:semi]
		}
		if header != "" {
			mime = header
		}
		payloadLen = len(payload)
	} else {
		payloadLen = len(value)
	}

	return fmt.Sprintf("[omitted image data %s, %d chars]", mime, payloadLen)
}

func looksLikeBase64Payload(value string) bool {
	compact := compactBase64(value)
	if len(compact) < 256 {
		return false
	}
	if strings.ContainsAny(compact, ".:/\\") {
		return false
	}

	for _, r := range compact {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '+', r == '/', r == '=', r == '-', r == '_':
		default:
			return false
		}
	}

	return true
}

func compactBase64(value string) string {
	replacer := strings.NewReplacer("\n", "", "\r", "", "\t", "", " ", "")
	return replacer.Replace(value)
}

func redactTextFallback(raw []byte) []byte {
	if len(raw) == 0 {
		return nil
	}

	text := string(raw)
	changed := false

	text = dataURLPattern.ReplaceAllStringFunc(text, func(match string) string {
		changed = true
		return redactDataURLValue(match)
	})

	text = keyedBase64Pattern.ReplaceAllStringFunc(text, func(match string) string {
		submatch := keyedBase64Pattern.FindStringSubmatch(match)
		if len(submatch) != 4 {
			return match
		}
		changed = true
		return submatch[1] + fmt.Sprintf("[omitted image base64, %d chars]", len(compactBase64(submatch[2]))) + submatch[3]
	})

	if !changed {
		return bytes.Clone(raw)
	}

	return []byte(text)
}
