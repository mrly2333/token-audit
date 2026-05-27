package audit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
)

type LimitedBuffer struct {
	max       int64
	total     int64
	truncated bool
	buf       bytes.Buffer
}

func NewLimitedBuffer(max int64) *LimitedBuffer {
	return &LimitedBuffer{max: max}
}

func (b *LimitedBuffer) Write(p []byte) (int, error) {
	if b == nil {
		return len(p), nil
	}

	b.total += int64(len(p))
	if b.max <= 0 {
		b.truncated = true
		return len(p), nil
	}

	remaining := b.max - int64(b.buf.Len())
	if remaining > 0 {
		toWrite := p
		if int64(len(toWrite)) > remaining {
			toWrite = toWrite[:remaining]
		}
		_, _ = b.buf.Write(toWrite)
	}
	if b.total > b.max {
		b.truncated = true
	}

	return len(p), nil
}

func (b *LimitedBuffer) Bytes() []byte {
	if b == nil {
		return nil
	}
	return bytes.Clone(b.buf.Bytes())
}

func (b *LimitedBuffer) Total() int64 {
	if b == nil {
		return 0
	}
	return b.total
}

func (b *LimitedBuffer) Truncated() bool {
	if b == nil {
		return false
	}
	return b.truncated
}

func SanitizeHeaders(headers http.Header) map[string][]string {
	clean := make(map[string][]string, len(headers))
	for key, values := range headers {
		switch strings.ToLower(key) {
		case "authorization", "cookie", "set-cookie":
			continue
		}
		clean[key] = append([]string(nil), values...)
	}
	return clean
}

func TokenMetadata(authorization, secret string) (string, string) {
	if authorization == "" || secret == "" {
		return "", ""
	}

	token := strings.TrimSpace(authorization)
	parts := strings.SplitN(token, " ", 2)
	if len(parts) == 2 {
		token = strings.TrimSpace(parts[1])
	}
	if token == "" {
		return "", ""
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil)), previewToken(token)
}

func VerifyFingerprint(token, fingerprint, secret string) bool {
	if token == "" || fingerprint == "" || secret == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(token))
	expected := hex.EncodeToString(mac.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(expected), []byte(fingerprint)) == 1
}

func previewToken(token string) string {
	switch {
	case len(token) <= 4:
		return "****"
	case len(token) <= 10:
		return token[:2] + "..." + token[len(token)-2:]
	default:
		return token[:6] + "..." + token[len(token)-4:]
	}
}

func LooksLikeJSON(contentType string, body []byte) bool {
	if strings.Contains(strings.ToLower(contentType), "json") {
		return true
	}
	trimmed := bytes.TrimSpace(body)
	return json.Valid(trimmed)
}

func IsSSE(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "text/event-stream")
}

func ParseRequest(body []byte, contentType string) (jsonBody []byte, model string, stream *bool, userText string) {
	if len(body) == 0 || !LooksLikeJSON(contentType, body) {
		return nil, "", nil, ""
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, "", nil, ""
	}

	if modelValue, ok := payload["model"].(string); ok {
		model = modelValue
	}
	if streamValue, ok := payload["stream"].(bool); ok {
		stream = &streamValue
	}

	texts := collectUserTexts(payload["messages"])
	if len(texts) == 0 {
		texts = extractTextParts(payload["input"])
	}
	userText = RedactTextContent(strings.Join(compactStrings(texts), "\n\n"))

	return bytes.Clone(body), model, stream, userText
}

func ParseResponse(body []byte, contentType string) (jsonBody []byte, model string, assistantText string, usageJSON []byte) {
	if len(body) == 0 || !LooksLikeJSON(contentType, body) {
		return nil, "", "", nil
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, "", "", nil
	}

	texts := collectResponseTexts(payload)
	model = extractModel(payload)
	usageJSON = marshalJSON(payload["usage"])
	if len(usageJSON) == 0 {
		if responseObj, ok := payload["response"].(map[string]any); ok {
			usageJSON = marshalJSON(responseObj["usage"])
		}
	}

	return bytes.Clone(body), model, RedactTextContent(strings.Join(compactStrings(texts), "\n\n")), usageJSON
}

type SSEAccumulator struct {
	pending   bytes.Buffer
	textParts []string
	usageJSON []byte
	model     string
}

func (a *SSEAccumulator) Feed(chunk []byte) {
	if len(chunk) == 0 {
		return
	}

	_, _ = a.pending.Write(chunk)
	for {
		line, ok := a.nextLine()
		if !ok {
			return
		}
		a.consumeLine(line)
	}
}

func (a *SSEAccumulator) Finalize() (assistantText string, usageJSON []byte, model string) {
	if a.pending.Len() > 0 {
		a.consumeLine(strings.TrimSuffix(a.pending.String(), "\r"))
		a.pending.Reset()
	}
	return RedactTextContent(strings.Join(compactStrings(a.textParts), "")), bytes.Clone(a.usageJSON), a.model
}

func (a *SSEAccumulator) nextLine() (string, bool) {
	data := a.pending.Bytes()
	idx := bytes.IndexByte(data, '\n')
	if idx < 0 {
		return "", false
	}

	line := string(bytes.TrimSuffix(data[:idx], []byte{'\r'}))
	rest := bytes.Clone(data[idx+1:])
	a.pending.Reset()
	_, _ = a.pending.Write(rest)
	return line, true
}

func (a *SSEAccumulator) consumeLine(line string) {
	if line == "" || !strings.HasPrefix(line, "data:") {
		return
	}

	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if payload == "" || payload == "[DONE]" {
		return
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		return
	}

	if a.model == "" {
		a.model = extractModel(decoded)
	}
	a.textParts = append(a.textParts, collectSSETexts(decoded)...)

	if usage := marshalJSON(decoded["usage"]); len(usage) > 0 {
		a.usageJSON = usage
		return
	}
	if responseObj, ok := decoded["response"].(map[string]any); ok {
		if usage := marshalJSON(responseObj["usage"]); len(usage) > 0 {
			a.usageJSON = usage
		}
	}

	// Claude streaming: message_start event contains message.usage
	if eventType, _ := decoded["type"].(string); eventType == "message_start" {
		if message, ok := decoded["message"].(map[string]any); ok {
			if a.model == "" {
				a.model = extractModel(message)
			}
			if usage := marshalJSON(message["usage"]); len(usage) > 0 {
				a.usageJSON = usage
			}
		}
	}
}

func compactStrings(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}

func collectUserTexts(messages any) []string {
	list, ok := messages.([]any)
	if !ok {
		return nil
	}

	var texts []string
	for _, item := range list {
		message, ok := item.(map[string]any)
		if !ok {
			continue
		}
		role, _ := message["role"].(string)
		if role != "user" {
			continue
		}
		texts = append(texts, extractTextParts(message["content"])...)
	}
	return texts
}

func collectResponseTexts(payload map[string]any) []string {
	var texts []string
	if choices, ok := payload["choices"].([]any); ok {
		for _, item := range choices {
			choice, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if message, ok := choice["message"].(map[string]any); ok {
				texts = append(texts, extractTextParts(message["content"])...)
			}
			if text, ok := choice["text"].(string); ok {
				texts = append(texts, text)
			}
		}
	}
	if outputText, ok := payload["output_text"].(string); ok {
		texts = append(texts, outputText)
	}
	if output, ok := payload["output"].([]any); ok {
		for _, item := range output {
			outputItem, ok := item.(map[string]any)
			if !ok {
				continue
			}
			texts = append(texts, extractTextParts(outputItem["content"])...)
		}
	}
	// Claude format: top-level content array
	if content, ok := payload["content"].([]any); ok {
		for _, item := range content {
			contentBlock, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if text, ok := contentBlock["text"].(string); ok && text != "" {
				texts = append(texts, text)
			}
		}
	}
	return texts
}

func collectSSETexts(payload map[string]any) []string {
	var texts []string

	if choices, ok := payload["choices"].([]any); ok {
		for _, item := range choices {
			choice, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if delta, ok := choice["delta"].(map[string]any); ok {
				texts = append(texts, extractTextParts(delta["content"])...)
			}
			if text, ok := choice["text"].(string); ok {
				texts = append(texts, text)
			}
		}
	}

	if eventType, _ := payload["type"].(string); strings.HasSuffix(eventType, ".delta") {
		if delta, ok := payload["delta"].(string); ok {
			texts = append(texts, delta)
		}
	}

	// Claude streaming: content_block_delta event
	if eventType, _ := payload["type"].(string); eventType == "content_block_delta" {
		if delta, ok := payload["delta"].(map[string]any); ok {
			if text, ok := delta["text"].(string); ok && text != "" {
				texts = append(texts, text)
			}
		}
	}

	return texts
}

func extractModel(payload map[string]any) string {
	if payload == nil {
		return ""
	}
	if model, _ := payload["model"].(string); model != "" {
		return model
	}
	if responseObj, ok := payload["response"].(map[string]any); ok {
		if model, _ := responseObj["model"].(string); model != "" {
			return model
		}
	}
	return ""
}

func extractTextParts(value any) []string {
	switch typed := value.(type) {
	case nil:
		return nil
	case string:
		return []string{typed}
	case []any:
		var texts []string
		for _, item := range typed {
			texts = append(texts, extractTextParts(item)...)
		}
		return texts
	case map[string]any:
		if text, ok := typed["text"].(string); ok {
			return []string{text}
		}
		if content, ok := typed["content"]; ok {
			return extractTextParts(content)
		}
		if delta, ok := typed["delta"]; ok {
			return extractTextParts(delta)
		}
		return nil
	default:
		return nil
	}
}

func marshalJSON(value any) []byte {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil || len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil
	}
	return raw
}
