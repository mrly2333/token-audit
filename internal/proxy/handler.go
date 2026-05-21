package proxy

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"newapi-audit-proxy/internal/audit"
	"newapi-audit-proxy/internal/config"
)

var hopByHopHeaders = map[string]struct{}{
	"Connection":          {},
	"Proxy-Connection":    {},
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"Te":                  {},
	"Trailer":             {},
	"Transfer-Encoding":   {},
	"Upgrade":             {},
}

type Handler struct {
	cfg        config.Config
	logger     *log.Logger
	store      *audit.Store
	client     *http.Client
	upstream   *url.URL
	captureSet map[string]struct{}
}

func New(cfg config.Config, store *audit.Store, logger *log.Logger) (*Handler, error) {
	upstream, err := url.Parse(cfg.UpstreamBase)
	if err != nil {
		return nil, fmt.Errorf("parse upstream base: %w", err)
	}

	captureSet := make(map[string]struct{}, len(cfg.CapturePaths))
	for _, path := range cfg.CapturePaths {
		captureSet[path] = struct{}{}
	}

	return &Handler{
		cfg:    cfg,
		logger: logger,
		store:  store,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy:              http.ProxyFromEnvironment,
				DisableCompression: true,
				ForceAttemptHTTP2:  true,
			},
		},
		upstream:   upstream,
		captureSet: captureSet,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startedAt := time.Now()
	record := audit.Record{
		StartedAt:          startedAt,
		Method:             r.Method,
		Path:               r.URL.Path,
		QueryString:        r.URL.RawQuery,
		RemoteAddr:         r.RemoteAddr,
		RequestHost:        r.Host,
		UpstreamBase:       h.cfg.UpstreamBase,
		RequestHeaders:     audit.SanitizeHeaders(r.Header),
		RequestContentType: r.Header.Get("Content-Type"),
	}

	record.TokenFingerprint, record.TokenPreview = audit.TokenMetadata(r.Header.Get("Authorization"), h.cfg.HMACSecret)
	_, record.IsCapturePath = h.captureSet[r.URL.Path]

	defer func() {
		if !record.IsCapturePath {
			return
		}
		record.FinishedAt = time.Now()
		record.DurationMS = record.FinishedAt.Sub(startedAt).Milliseconds()
		h.store.InsertAsync(record)
	}()

	bodyReader, cleanup, err := h.prepareRequestBody(&record, r)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		h.logger.Printf("capture request body failed: %v", err)
		h.writeLocalError(w, &record, http.StatusBadRequest, "failed to read request body")
		return
	}

	upstreamReq, err := h.buildUpstreamRequest(r, bodyReader)
	if err != nil {
		h.logger.Printf("build upstream request failed: %v", err)
		h.writeLocalError(w, &record, http.StatusBadGateway, "failed to build upstream request")
		return
	}

	upstreamResp, err := h.client.Do(upstreamReq)
	if err != nil {
		if isExpectedDisconnect(r.Context(), err) {
			record.StatusCode = 499
			record.ErrorText = "client canceled upstream request"
			return
		}
		h.logger.Printf("upstream request failed: %v", err)
		record.ErrorText = fmt.Sprintf("upstream request failed: %v", err)
		h.writeLocalError(w, &record, http.StatusBadGateway, "upstream request failed")
		return
	}
	defer upstreamResp.Body.Close()

	record.StatusCode = upstreamResp.StatusCode
	record.ResponseHeaders = audit.SanitizeHeaders(upstreamResp.Header)
	record.ResponseContentType = upstreamResp.Header.Get("Content-Type")
	record.ResponseIsSSE = audit.IsSSE(record.ResponseContentType)

	copyResponseHeaders(w.Header(), upstreamResp.Header)
	w.WriteHeader(upstreamResp.StatusCode)

	flusher, _ := w.(http.Flusher)
	responseBuffer := audit.NewLimitedBuffer(h.cfg.MaxBodyBytes)
	var sseAccumulator audit.SSEAccumulator

	if err := streamResponse(w, upstreamResp.Body, flusher, func(chunk []byte) {
		if record.IsCapturePath {
			_, _ = responseBuffer.Write(chunk)
		}
		if record.ResponseIsSSE {
			sseAccumulator.Feed(chunk)
		}
	}); err != nil {
		if isExpectedDisconnect(r.Context(), err) {
			record.ErrorText = fmt.Sprintf("client disconnected during response stream: %v", err)
		} else {
			record.ErrorText = fmt.Sprintf("stream response failed: %v", err)
			h.logger.Printf("stream response failed: %v", err)
		}
	}

	record.ResponseBytes = responseBuffer.Total()
	record.ResponseTruncated = responseBuffer.Truncated()
	rawResponseBody := responseBuffer.Bytes()

	if record.ResponseIsSSE {
		var responseModel string
		record.AssistantText, record.UsageJSON, responseModel = sseAccumulator.Finalize()
		if responseModel != "" {
			record.Model = responseModel
		}
	} else if record.IsCapturePath {
		var responseModel string
		record.ResponseJSON, responseModel, record.AssistantText, record.UsageJSON = audit.ParseResponse(rawResponseBody, record.ResponseContentType)
		if responseModel != "" {
			record.Model = responseModel
		}
	}

	if record.IsCapturePath {
		record.UserText = audit.RedactTextContent(record.UserText)
		record.AssistantText = audit.RedactTextContent(record.AssistantText)
		record.ResponseBody = audit.RedactCapturedBody(rawResponseBody, record.ResponseContentType)
		record.ResponseJSON = audit.RedactCapturedJSON(record.ResponseJSON)
	}

	record.PromptTokens, record.CompletionTokens, record.TotalTokens = audit.ExtractUsageTotals(record.UsageJSON)
}

func (h *Handler) prepareRequestBody(record *audit.Record, r *http.Request) (io.ReadCloser, func(), error) {
	if r.Body == nil || r.Body == http.NoBody {
		return nil, nil, nil
	}

	if !record.IsCapturePath {
		return r.Body, nil, nil
	}

	tmpFile, err := os.CreateTemp("", "newapi-audit-request-*")
	if err != nil {
		return nil, nil, err
	}

	requestBuffer := audit.NewLimitedBuffer(h.cfg.MaxBodyBytes)
	multiWriter := io.MultiWriter(tmpFile, requestBuffer)
	_, copyErr := io.Copy(multiWriter, r.Body)
	closeErr := r.Body.Close()
	if copyErr != nil {
		cleanupTempFile(tmpFile)
		record.RequestBody = audit.RedactCapturedBody(requestBuffer.Bytes(), record.RequestContentType)
		record.RequestBytes = requestBuffer.Total()
		record.RequestTruncated = requestBuffer.Truncated()
		return nil, nil, copyErr
	}
	if closeErr != nil {
		cleanupTempFile(tmpFile)
		record.RequestBody = audit.RedactCapturedBody(requestBuffer.Bytes(), record.RequestContentType)
		record.RequestBytes = requestBuffer.Total()
		record.RequestTruncated = requestBuffer.Truncated()
		return nil, nil, closeErr
	}

	rawRequestBody := requestBuffer.Bytes()
	record.RequestBytes = requestBuffer.Total()
	record.RequestTruncated = requestBuffer.Truncated()
	record.RequestJSON, record.Model, record.Stream, record.UserText = audit.ParseRequest(rawRequestBody, record.RequestContentType)
	record.RequestBody = audit.RedactCapturedBody(rawRequestBody, record.RequestContentType)
	record.RequestJSON = audit.RedactCapturedJSON(record.RequestJSON)

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		cleanupTempFile(tmpFile)
		return nil, nil, err
	}

	reader := &tempFileReadCloser{file: tmpFile, name: tmpFile.Name()}
	return reader, func() { _ = reader.Close() }, nil
}

func (h *Handler) buildUpstreamRequest(r *http.Request, body io.ReadCloser) (*http.Request, error) {
	targetURL := *h.upstream
	targetURL.Path = joinURLPath(h.upstream.Path, r.URL.Path)
	targetURL.RawPath = targetURL.Path
	targetURL.RawQuery = r.URL.RawQuery

	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header = r.Header.Clone()
	removeHopByHopHeaders(req.Header)
	req.Header.Set("Accept-Encoding", "identity")
	req.Host = h.upstream.Host
	req.ContentLength = r.ContentLength

	return req, nil
}

func (h *Handler) writeLocalError(w http.ResponseWriter, record *audit.Record, status int, message string) {
	responseBody := []byte(message + "\n")
	record.StatusCode = status
	if record.ErrorText == "" {
		record.ErrorText = message
	}
	record.ResponseContentType = "text/plain; charset=utf-8"
	record.ResponseHeaders = map[string][]string{
		"Content-Type": {"text/plain; charset=utf-8"},
	}
	record.ResponseBytes = int64(len(responseBody))
	if record.IsCapturePath {
		record.ResponseBody = bytes.Clone(responseBody)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(responseBody)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func streamResponse(
	w http.ResponseWriter,
	body io.ReadCloser,
	flusher http.Flusher,
	onChunk func([]byte),
) error {
	buf := make([]byte, 32*1024)

	for {
		n, readErr := body.Read(buf)
		if n > 0 {
			chunk := bytes.Clone(buf[:n])
			if onChunk != nil {
				onChunk(chunk)
			}

			written, writeErr := w.Write(chunk)
			if flusher != nil {
				flusher.Flush()
			}
			if writeErr != nil {
				return writeErr
			}
			if written != len(chunk) {
				return io.ErrShortWrite
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return nil
			}
			return readErr
		}
	}
}

func joinURLPath(basePath, reqPath string) string {
	switch {
	case basePath == "" || basePath == "/":
		return reqPath
	case strings.HasSuffix(basePath, "/") && strings.HasPrefix(reqPath, "/"):
		return basePath + strings.TrimPrefix(reqPath, "/")
	case !strings.HasSuffix(basePath, "/") && !strings.HasPrefix(reqPath, "/"):
		return basePath + "/" + reqPath
	default:
		return basePath + reqPath
	}
}

func copyResponseHeaders(dst, src http.Header) {
	for key := range dst {
		dst.Del(key)
	}

	headerCopy := src.Clone()
	removeHopByHopHeaders(headerCopy)
	for key, values := range headerCopy {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func removeHopByHopHeaders(headers http.Header) {
	for _, token := range connectionTokens(headers) {
		headers.Del(token)
	}
	for header := range hopByHopHeaders {
		headers.Del(header)
	}
}

func connectionTokens(headers http.Header) []string {
	raw := headers.Values("Connection")
	if len(raw) == 0 {
		return nil
	}

	var tokens []string
	for _, value := range raw {
		for _, token := range strings.Split(value, ",") {
			token = strings.TrimSpace(token)
			if token != "" {
				tokens = append(tokens, token)
			}
		}
	}
	return tokens
}

func isExpectedDisconnect(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}
	if ctx != nil && errors.Is(ctx.Err(), context.Canceled) {
		return true
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, io.ErrClosedPipe) {
		return true
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "context canceled") ||
		strings.Contains(message, "broken pipe") ||
		strings.Contains(message, "connection reset by peer") ||
		strings.Contains(message, "client disconnected")
}

func cleanupTempFile(file *os.File) {
	name := file.Name()
	_ = file.Close()
	_ = os.Remove(name)
}

type tempFileReadCloser struct {
	file *os.File
	name string
}

func (t *tempFileReadCloser) Read(p []byte) (int, error) {
	return t.file.Read(p)
}

func (t *tempFileReadCloser) Close() error {
	err := t.file.Close()
	_ = os.Remove(t.name)
	return err
}
