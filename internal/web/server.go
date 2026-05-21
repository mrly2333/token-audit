package web

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"newapi-audit-proxy/internal/audit"
	"newapi-audit-proxy/internal/config"
)

const sessionCookieName = "newapi_audit_session"
const csrfFieldName = "csrf_token"
const csrfHeaderName = "X-CSRF-Token"
const maxLoginFailures = 5

const csrfTokenTTL = 24 * time.Hour
const loginAttemptWindow = 10 * time.Minute
const loginLockoutDuration = 10 * time.Minute

type Server struct {
	cfg       config.Config
	basePath  string
	logger    *log.Logger
	store     *audit.Store
	templates map[string]*template.Template
	loginMu   sync.Mutex
	loginFail map[string]loginAttempt
}

type loginViewData struct {
	Title     string
	Error     string
	CSRFToken string
}

type dashboardViewData struct {
	Title     string
	Now       time.Time
	Stats     audit.DashboardStats
	CSRFToken string
}

type logFiltersView struct {
	From             string
	To               string
	TokenFingerprint string
	TokenAlias       string
	Model            string
	StatusCode       string
	Keyword          string
}

type logsViewData struct {
	Title        string
	Filters      logFiltersView
	Page         int
	PageSize     int
	CurrentCount int
	TotalCount   int64
	TotalPages   int
	HasPrev      bool
	HasNext      bool
	PrevPage     int
	NextPage     int
	CSRFToken    string
}

type tokenFiltersView struct {
	From             string
	To               string
	TokenFingerprint string
	TokenAlias       string
	Model            string
}

type tokensViewData struct {
	Title      string
	Filters    tokenFiltersView
	Result     audit.TokenDirectoryResult
	PrevPage   int
	NextPage   int
	HasNext    bool
	Saved      bool
	Error      string
	CurrentURL string
	CSRFToken  string
}

type detailViewData struct {
	Title     string
	Log       audit.Record
	Images    []logImageAPIItem
	CSRFToken string
}

type loginAttempt struct {
	Failures    int
	FirstFailed time.Time
	LockedUntil time.Time
}

type logVersionAPIResponse struct {
	Version string `json:"version"`
}

type logListAPIItem struct {
	ID               int64  `json:"id"`
	StartedAt        string `json:"started_at"`
	TokenAlias       string `json:"token_alias"`
	TokenPreview     string `json:"token_preview"`
	TokenFingerprint string `json:"token_fingerprint"`
	Model            string `json:"model"`
	StatusCode       int    `json:"status_code"`
	TotalTokens      int64  `json:"total_tokens"`
	UserPreview      string `json:"user_preview"`
	AssistantPreview string `json:"assistant_preview"`
}

type logListAPIResponse struct {
	Version    string           `json:"version"`
	TotalCount int64            `json:"total_count"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
	HasPrev    bool             `json:"has_prev"`
	HasNext    bool             `json:"has_next"`
	Items      []logListAPIItem `json:"items"`
}

type logDetailAPIResponse struct {
	ID                int64                  `json:"id"`
	StartedAt         string                 `json:"started_at"`
	Method            string                 `json:"method"`
	PathWithQuery     string                 `json:"path_with_query"`
	StatusCode        int                    `json:"status_code"`
	Model             string                 `json:"model"`
	TokenAlias        string                 `json:"token_alias"`
	TokenPreview      string                 `json:"token_preview"`
	TokenFingerprint  string                 `json:"token_fingerprint"`
	PromptTokens      int64                  `json:"prompt_tokens"`
	CompletionTokens  int64                  `json:"completion_tokens"`
	TotalTokens       int64                  `json:"total_tokens"`
	DurationMS        int64                  `json:"duration_ms"`
	RequestBytes      int64                  `json:"request_bytes"`
	ResponseBytes     int64                  `json:"response_bytes"`
	RequestTruncated  bool                   `json:"request_truncated"`
	ResponseTruncated bool                   `json:"response_truncated"`
	Stream            string                 `json:"stream"`
	ResponseType      string                 `json:"response_type"`
	ErrorText         string                 `json:"error_text"`
	UserText          logTextPageAPIResponse `json:"user_text"`
	AssistantText     logTextPageAPIResponse `json:"assistant_text"`
}

type logDetailRawAPIResponse struct {
	ID              int64             `json:"id"`
	RequestBody     string            `json:"request_body"`
	ResponseBody    string            `json:"response_body"`
	RequestJSON     string            `json:"request_json"`
	ResponseJSON    string            `json:"response_json"`
	UsageJSON       string            `json:"usage_json"`
	RequestHeaders  string            `json:"request_headers"`
	ResponseHeaders string            `json:"response_headers"`
	Images          []logImageAPIItem `json:"images"`
}

type logTextPageAPIResponse struct {
	Kind       string `json:"kind"`
	Text       string `json:"text"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
	TotalChars int64  `json:"total_chars"`
}

type logImageAPIItem struct {
	Label   string `json:"label"`
	Mime    string `json:"mime"`
	DataURL string `json:"data_url"`
}

type dashboardAPIResponse struct {
	TotalRequests         int64                `json:"total_requests"`
	TodayRequests         int64                `json:"today_requests"`
	ErrorCount            int64                `json:"error_count"`
	DistinctTokens        int64                `json:"distinct_tokens"`
	TodayDistinctTokens   int64                `json:"today_distinct_tokens"`
	TotalPromptTokens     int64                `json:"total_prompt_tokens"`
	TotalCompletionTokens int64                `json:"total_completion_tokens"`
	TotalTokens           int64                `json:"total_tokens"`
	TodayTotalTokens      int64                `json:"today_total_tokens"`
	TokenGroups           []dashboardTokenItem `json:"token_groups"`
	ModelGroups           []dashboardModelItem `json:"model_groups"`
}

type dashboardTokenItem struct {
	TokenAlias       string `json:"token_alias"`
	TokenPreview     string `json:"token_preview"`
	TokenFingerprint string `json:"token_fingerprint"`
	RequestCount     int64  `json:"request_count"`
	TotalTokens      int64  `json:"total_tokens"`
	ErrorCount       int64  `json:"error_count"`
	LastSeen         string `json:"last_seen"`
}

type dashboardModelItem struct {
	Model        string `json:"model"`
	RequestCount int64  `json:"request_count"`
	TotalTokens  int64  `json:"total_tokens"`
	ErrorCount   int64  `json:"error_count"`
	LastSeen     string `json:"last_seen"`
}

type tokenListAPIResponse struct {
	Items []tokenListAPIItem `json:"items"`
}

type tokenListAPIItem struct {
	TokenAlias       string `json:"token_alias"`
	TokenPreview     string `json:"token_preview"`
	TokenFingerprint string `json:"token_fingerprint"`
	RequestCount     int64  `json:"request_count"`
	TotalTokens      int64  `json:"total_tokens"`
	LastSeen         string `json:"last_seen"`
}

type filterOptionsAPIResponse struct {
	Models            []string `json:"models"`
	TokenAliases      []string `json:"token_aliases"`
	TokenFingerprints []string `json:"token_fingerprints"`
	StatusCodes       []string `json:"status_codes"`
}

type dbStatsAPIResponse struct {
	TotalRows        int64  `json:"total_rows"`
	TodayRows        int64  `json:"today_rows"`
	DatabaseSize     int64  `json:"database_size"`
	DatabasePretty   string `json:"database_pretty"`
	AuditTotalSize   int64  `json:"audit_total_size"`
	AuditTotalPretty string `json:"audit_total_pretty"`
	AuditTableSize   int64  `json:"audit_table_size"`
	AuditTablePretty string `json:"audit_table_pretty"`
	AuditIndexSize   int64  `json:"audit_index_size"`
	AuditIndexPretty string `json:"audit_index_pretty"`
	AuditToastSize   int64  `json:"audit_toast_size"`
	AuditToastPretty string `json:"audit_toast_pretty"`
	LiveTuples       int64  `json:"live_tuples"`
	DeadTuples       int64  `json:"dead_tuples"`
	LastVacuum       string `json:"last_vacuum"`
	LastAutovacuum   string `json:"last_autovacuum"`
	LastAnalyze      string `json:"last_analyze"`
	LastAutoanalyze  string `json:"last_autoanalyze"`
}

type dbCleanupResponse struct {
	DeletedRows int64  `json:"deleted_rows"`
	Message     string `json:"message"`
}

type dbMaintenanceResponse struct {
	Mode         string `json:"mode"`
	AffectedRows int64  `json:"affected_rows"`
	Message      string `json:"message"`
}

func New(cfg config.Config, store *audit.Store, logger *log.Logger) (*Server, error) {
	basePath := cfg.WebBasePath
	funcMap := template.FuncMap{
		"basePath": func() string {
			return basePath
		},
		"path": func(path string) string {
			return joinBasePath(basePath, path)
		},
		"formatTime":    formatTime,
		"statusClass":   statusClass,
		"truncate":      truncate,
		"prettyJSON":    prettyJSON,
		"prettyAny":     prettyAny,
		"renderBody":    renderBody,
		"streamValue":   streamValue,
		"pathWithQuery": pathWithQuery,
		"pageURL": func(filters logFiltersView, page int) string {
			return buildPageURL(basePath, filters, page)
		},
		"tokenPageURL": func(filters tokenFiltersView, page int) string {
			return buildTokenPageURL(basePath, filters, page)
		},
	}

	pageSources := map[string]string{
		"login":     baseTemplate + loginTemplate,
		"dashboard": baseTemplate + dashboardTemplate,
		"logs":      baseTemplate + logsTemplateV2,
		"tokens":    baseTemplate + tokensTemplate,
		"detail":    baseTemplate + detailTemplate,
	}

	templates := make(map[string]*template.Template, len(pageSources))
	for name, source := range pageSources {
		tmpl, err := template.New(name).Funcs(funcMap).Parse(source)
		if err != nil {
			return nil, fmt.Errorf("parse template %s: %w", name, err)
		}
		templates[name] = tmpl
	}

	return &Server{
		cfg:       cfg,
		basePath:  basePath,
		logger:    logger,
		store:     store,
		templates: templates,
		loginFail: make(map[string]loginAttempt),
	}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /login", s.handleLoginPage)
	mux.HandleFunc("POST /login", s.requireCSRF(s.handleLoginSubmit))
	mux.HandleFunc("POST /logout", s.requireAuth(s.requireCSRF(s.handleLogout)))
	mux.HandleFunc("GET /", s.requireAuth(s.handleDashboard))
	mux.HandleFunc("GET /api/dashboard", s.requireAuth(s.handleDashboardAPI))
	mux.HandleFunc("GET /logs", s.requireAuth(s.handleLogs))
	mux.HandleFunc("GET /api/logs/version", s.requireAuth(s.handleLogsVersionAPI))
	mux.HandleFunc("GET /api/logs", s.requireAuth(s.handleLogsListAPI))
	mux.HandleFunc("GET /api/filter-options", s.requireAuth(s.handleFilterOptionsAPI))
	mux.HandleFunc("GET /api/logs/{id}", s.requireAuth(s.handleLogDetailAPI))
	mux.HandleFunc("GET /api/logs/{id}/text", s.requireAuth(s.handleLogTextPageAPIv2))
	mux.HandleFunc("GET /api/logs/{id}/raw", s.requireAuth(s.handleLogDetailRawAPI))
	mux.HandleFunc("GET /api/tokens", s.requireAuth(s.handleTokensAPI))
	mux.HandleFunc("POST /api/tokens/alias", s.requireAuth(s.requireCSRF(s.handleTokenAliasAPI)))
	mux.HandleFunc("GET /api/db/stats", s.requireAuth(s.handleDBStatsAPI))
	mux.HandleFunc("POST /api/db/cleanup", s.requireAuth(s.requireCSRF(s.handleDBCleanupAPI)))
	mux.HandleFunc("POST /api/db/maintenance", s.requireAuth(s.requireCSRF(s.handleDBMaintenanceAPI)))
	mux.HandleFunc("GET /tokens", s.requireAuth(s.handleTokens))
	mux.HandleFunc("POST /tokens/alias", s.requireAuth(s.requireCSRF(s.handleTokenAliasSubmit)))
	mux.HandleFunc("GET /logs/{id}", s.requireAuth(s.handleLogDetail))
	return mux
}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if s.isAuthenticated(r) {
		http.Redirect(w, r, s.webPath("/logs"), http.StatusSeeOther)
		return
	}
	s.render(w, "login", loginViewData{Title: "登录 - newapi-audit-proxy"})
}

func (s *Server) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderLoginError(w, http.StatusBadRequest, "解析登录表单失败")
		return
	}

	username := strings.TrimSpace(r.Form.Get("username"))
	password := r.Form.Get("password")
	if retryAfter, limited := s.isLoginRateLimited(r); limited {
		s.renderLoginError(w, http.StatusTooManyRequests, fmt.Sprintf("登录失败次数过多，请 %s 后再试", formatDurationForHuman(retryAfter)))
		return
	}
	passwordOK := bcrypt.CompareHashAndPassword([]byte(s.cfg.AdminPasswordHash), []byte(password)) == nil
	if username != s.cfg.AdminUsername || !passwordOK {
		s.recordLoginFailure(r)
		s.renderLoginError(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	s.clearLoginFailures(r)
	s.issueSession(w, r)
	http.Redirect(w, r, s.webPath("/logs"), http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     s.cookiePath(),
		HttpOnly: true,
		Secure:   s.secureCookie(r),
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, s.webPath("/login"), http.StatusSeeOther)
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.webPath("/logs"), http.StatusSeeOther)
}

func (s *Server) handleDashboardAPI(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.DashboardStats(r.Context(), time.Now())
	if err != nil {
		s.logger.Printf("dashboard api failed: %v", err)
		http.Error(w, "查询统计失败", http.StatusInternalServerError)
		return
	}

	resp := dashboardAPIResponse{
		TotalRequests:         stats.TotalRequests,
		TodayRequests:         stats.TodayRequests,
		ErrorCount:            stats.ErrorCount,
		DistinctTokens:        stats.DistinctTokens,
		TodayDistinctTokens:   stats.TodayDistinctTokens,
		TotalPromptTokens:     stats.TotalPromptTokens,
		TotalCompletionTokens: stats.TotalCompletionTokens,
		TotalTokens:           stats.TotalTokens,
		TodayTotalTokens:      stats.TodayTotalTokens,
		TokenGroups:           make([]dashboardTokenItem, 0, len(stats.TokenGroups)),
		ModelGroups:           make([]dashboardModelItem, 0, len(stats.ModelGroups)),
	}
	for _, item := range stats.TokenGroups {
		resp.TokenGroups = append(resp.TokenGroups, dashboardTokenItem{
			TokenAlias:       item.TokenAlias,
			TokenPreview:     item.TokenPreview,
			TokenFingerprint: item.TokenFingerprint,
			RequestCount:     item.RequestCount,
			TotalTokens:      item.TotalTokens,
			ErrorCount:       item.ErrorCount,
			LastSeen:         formatTime(item.LastSeen),
		})
	}
	for _, item := range stats.ModelGroups {
		resp.ModelGroups = append(resp.ModelGroups, dashboardModelItem{
			Model:        item.Model,
			RequestCount: item.RequestCount,
			TotalTokens:  item.TotalTokens,
			ErrorCount:   item.ErrorCount,
			LastSeen:     formatTime(item.LastSeen),
		})
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filters := parseLogFilters(query)
	result, err := s.store.ListLogs(r.Context(), filters)
	if err != nil {
		s.logger.Printf("list logs for page failed: %v", err)
		http.Error(w, "查询日志失败", http.StatusInternalServerError)
		return
	}

	totalPages := 1
	if result.TotalCount > 0 && result.PageSize > 0 {
		totalPages = int((result.TotalCount + int64(result.PageSize) - 1) / int64(result.PageSize))
	}
	prevPage := result.Page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := result.Page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	s.render(w, "logs", logsViewData{
		Title: "日志 - newapi-audit-proxy",
		Filters: logFiltersView{
			From:             query.Get("from"),
			To:               query.Get("to"),
			TokenFingerprint: strings.TrimSpace(query.Get("token")),
			TokenAlias:       strings.TrimSpace(query.Get("alias")),
			Model:            strings.TrimSpace(query.Get("model")),
			StatusCode:       query.Get("status"),
			Keyword:          strings.TrimSpace(query.Get("q")),
		},
		Page:         result.Page,
		PageSize:     result.PageSize,
		CurrentCount: len(result.Items),
		TotalCount:   result.TotalCount,
		TotalPages:   totalPages,
		HasPrev:      result.Page > 1,
		HasNext:      result.Page < totalPages,
		PrevPage:     prevPage,
		NextPage:     nextPage,
	})
}

func (s *Server) handleLogsVersionAPI(w http.ResponseWriter, r *http.Request) {
	version, err := s.store.LogsVersion(r.Context(), parseLogFilters(r.URL.Query()))
	if err != nil {
		s.logger.Printf("query log version failed: %v", err)
		http.Error(w, "查询日志版本失败", http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, logVersionAPIResponse{Version: version})
}

func (s *Server) handleLogsListAPI(w http.ResponseWriter, r *http.Request) {
	filters := parseLogFilters(r.URL.Query())
	result, err := s.store.ListLogs(r.Context(), filters)
	if err != nil {
		s.logger.Printf("list logs api failed: %v", err)
		http.Error(w, "查询日志失败", http.StatusInternalServerError)
		return
	}

	version, err := s.store.LogsVersion(r.Context(), filters)
	if err != nil {
		s.logger.Printf("query log version failed: %v", err)
		http.Error(w, "查询日志版本失败", http.StatusInternalServerError)
		return
	}

	resp := logListAPIResponse{
		Version:    version,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		Items:      make([]logListAPIItem, 0, len(result.Items)),
	}
	if result.TotalCount > 0 && result.PageSize > 0 {
		resp.TotalPages = int((result.TotalCount + int64(result.PageSize) - 1) / int64(result.PageSize))
	} else {
		resp.TotalPages = 1
	}
	resp.HasPrev = resp.Page > 1
	resp.HasNext = resp.Page < resp.TotalPages
	for _, item := range result.Items {
		resp.Items = append(resp.Items, logListAPIItem{
			ID:               item.ID,
			StartedAt:        formatTime(item.StartedAt),
			TokenAlias:       item.TokenAlias,
			TokenPreview:     item.TokenPreview,
			TokenFingerprint: item.TokenFingerprint,
			Model:            item.Model,
			StatusCode:       item.StatusCode,
			TotalTokens:      item.TotalTokens,
			UserPreview:      previewText(audit.RedactTextContent(item.UserText), 56),
			AssistantPreview: previewText(audit.RedactTextContent(item.AssistantText), 56),
		})
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleFilterOptionsAPI(w http.ResponseWriter, r *http.Request) {
	options, err := s.store.FilterOptions(r.Context())
	if err != nil {
		s.logger.Printf("filter options api failed: %v", err)
		http.Error(w, "查询筛选候选项失败", http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, filterOptionsAPIResponse{
		Models:            options.Models,
		TokenAliases:      options.TokenAliases,
		TokenFingerprints: options.TokenFingerprints,
		StatusCodes:       options.StatusCodes,
	})
}

func (s *Server) handleLogDetailAPI(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	record, err := s.store.GetLogCore(r.Context(), id)
	if err != nil {
		http.Error(w, "日志不存在", http.StatusNotFound)
		return
	}

	responseType := "普通 JSON"
	if record.Record.ResponseIsSSE {
		responseType = "SSE"
	}

	s.writeJSON(w, http.StatusOK, logDetailAPIResponse{
		ID:                record.Record.ID,
		StartedAt:         formatTime(record.Record.StartedAt),
		Method:            record.Record.Method,
		PathWithQuery:     pathWithQuery(record.Record.Path, record.Record.QueryString),
		StatusCode:        record.Record.StatusCode,
		Model:             record.Record.Model,
		TokenAlias:        record.Record.TokenAlias,
		TokenPreview:      record.Record.TokenPreview,
		TokenFingerprint:  record.Record.TokenFingerprint,
		PromptTokens:      record.Record.PromptTokens,
		CompletionTokens:  record.Record.CompletionTokens,
		TotalTokens:       record.Record.TotalTokens,
		DurationMS:        record.Record.DurationMS,
		RequestBytes:      record.Record.RequestBytes,
		ResponseBytes:     record.Record.ResponseBytes,
		RequestTruncated:  record.Record.RequestTruncated,
		ResponseTruncated: record.Record.ResponseTruncated,
		Stream:            streamValue(record.Record.Stream),
		ResponseType:      responseType,
		ErrorText:         record.Record.ErrorText,
		UserText:          toLogTextPageAPIResponse(record.UserText),
		AssistantText:     toLogTextPageAPIResponse(record.AssistantText),
	})
}

func (s *Server) handleLogTextPageAPI(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	kind := strings.TrimSpace(r.URL.Query().Get("kind"))
	page := max(1, parseInt(r.URL.Query().Get("page")))
	textPage, err := s.store.GetLogTextPage(r.Context(), id, kind, page)
	if err != nil {
		if err == audit.ErrUnsupportedTextKind {
			http.Error(w, "不支持的文本类型", http.StatusBadRequest)
			return
		}
		http.Error(w, "日志不存在", http.StatusNotFound)
		return
	}

	s.writeJSON(w, http.StatusOK, toLogTextPageAPIResponse(textPage))
}

func (s *Server) handleLogTextPageAPIv2(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	kind := strings.TrimSpace(r.URL.Query().Get("kind"))
	page := max(1, parseInt(r.URL.Query().Get("page")))
	textPage, err := s.store.GetLogTextPage(r.Context(), id, kind, page)
	if err != nil {
		if err == audit.ErrUnsupportedTextKind {
			http.Error(w, "不支持的文本类型", http.StatusBadRequest)
			return
		}
		if err == audit.ErrAuditLogNotFound {
			http.Error(w, "日志不存在或文本分页已失效", http.StatusNotFound)
			return
		}
		s.logger.Printf("load log text page failed: id=%d kind=%s page=%d err=%v", id, kind, page, err)
		http.Error(w, "加载文本分页失败", http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, toLogTextPageAPIResponse(textPage))
}

func (s *Server) handleLogDetailRawAPI(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	record, err := s.store.GetLog(r.Context(), id)
	if err != nil {
		http.Error(w, "日志不存在", http.StatusNotFound)
		return
	}

	record.UserText = audit.RedactTextContent(record.UserText)
	record.AssistantText = audit.RedactTextContent(record.AssistantText)
	record.RequestBody = audit.RedactCapturedBody(record.RequestBody, record.RequestContentType)
	record.ResponseBody = audit.RedactCapturedBody(record.ResponseBody, record.ResponseContentType)
	record.RequestJSON = audit.RedactCapturedJSON(record.RequestJSON)
	record.ResponseJSON = audit.RedactCapturedJSON(record.ResponseJSON)

	s.writeJSON(w, http.StatusOK, logDetailRawAPIResponse{
		ID:              record.ID,
		RequestBody:     renderBody(record.RequestBody),
		ResponseBody:    renderBody(record.ResponseBody),
		RequestJSON:     prettyJSON(record.RequestJSON),
		ResponseJSON:    prettyJSON(record.ResponseJSON),
		UsageJSON:       prettyJSON(record.UsageJSON),
		RequestHeaders:  prettyAny(record.RequestHeaders),
		ResponseHeaders: prettyAny(record.ResponseHeaders),
		Images:          extractImagesFromRecord(record),
	})
}

func (s *Server) handleTokensAPI(w http.ResponseWriter, r *http.Request) {
	filters := parseTokenFilters(r.URL.Query())
	result, err := s.store.ListTokens(r.Context(), filters)
	if err != nil {
		s.logger.Printf("list tokens api failed: %v", err)
		http.Error(w, "查询 token 数据失败", http.StatusInternalServerError)
		return
	}

	resp := tokenListAPIResponse{
		Items: make([]tokenListAPIItem, 0, len(result.Items)),
	}
	for _, item := range result.Items {
		resp.Items = append(resp.Items, tokenListAPIItem{
			TokenAlias:       item.TokenAlias,
			TokenPreview:     item.TokenPreview,
			TokenFingerprint: item.TokenFingerprint,
			RequestCount:     item.RequestCount,
			TotalTokens:      item.TotalTokens,
			LastSeen:         formatTime(item.LastSeen),
		})
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleTokenAliasAPI(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "解析 token 别名表单失败", http.StatusBadRequest)
		return
	}

	tokenFingerprint, tokenAlias, err := s.resolveTokenAliasInput(r.Form.Get("token_fingerprint"), r.Form.Get("token_value"), r.Form.Get("token_alias"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.store.UpsertTokenAlias(r.Context(), tokenFingerprint, tokenAlias); err != nil {
		s.logger.Printf("save token alias api failed: %v", err)
		http.Error(w, friendlyAliasError(err), http.StatusBadRequest)
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]string{
		"message": "令牌代号已保存",
	})
}

func (s *Server) handleDBStatsAPI(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.DatabaseStats(r.Context(), time.Now())
	if err != nil {
		s.logger.Printf("db stats api failed: %v", err)
		http.Error(w, "查询数据库统计失败", http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, dbStatsAPIResponse{
		TotalRows:        stats.TotalRows,
		TodayRows:        stats.TodayRows,
		DatabaseSize:     stats.DatabaseSize,
		DatabasePretty:   stats.DatabasePretty,
		AuditTotalSize:   stats.AuditTotalSize,
		AuditTotalPretty: stats.AuditTotalPretty,
		AuditTableSize:   stats.AuditTableSize,
		AuditTablePretty: stats.AuditTablePretty,
		AuditIndexSize:   stats.AuditIndexSize,
		AuditIndexPretty: stats.AuditIndexPretty,
		AuditToastSize:   stats.AuditToastSize,
		AuditToastPretty: stats.AuditToastPretty,
		LiveTuples:       stats.LiveTuples,
		DeadTuples:       stats.DeadTuples,
		LastVacuum:       formatTime(stats.LastVacuum),
		LastAutovacuum:   formatTime(stats.LastAutovacuum),
		LastAnalyze:      formatTime(stats.LastAnalyze),
		LastAutoanalyze:  formatTime(stats.LastAutoanalyze),
	})
}

func (s *Server) handleDBCleanupAPI(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "解析数据清理表单失败", http.StatusBadRequest)
		return
	}

	filters := audit.DeleteLogsFilters{
		From:             parseLocalDateTime(r.Form.Get("from")),
		To:               parseLocalDateTime(r.Form.Get("to")),
		TokenFingerprint: strings.TrimSpace(r.Form.Get("token")),
		TokenAlias:       strings.TrimSpace(r.Form.Get("alias")),
		Model:            strings.TrimSpace(r.Form.Get("model")),
	}

	deletedRows, err := s.store.DeleteLogs(r.Context(), filters)
	if err != nil {
		if strings.Contains(err.Error(), "at least one delete filter is required") {
			http.Error(w, "请至少提供一个清理条件", http.StatusBadRequest)
			return
		}
		s.logger.Printf("db cleanup api failed: %v", err)
		http.Error(w, "清理数据失败", http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, dbCleanupResponse{
		DeletedRows: deletedRows,
		Message:     fmt.Sprintf("已删除 %d 条记录。PostgreSQL 删除后通常不会立刻缩小磁盘文件，如需回收磁盘空间，请继续执行“整理空间”或“强制缩盘”。", deletedRows),
	})
}

func (s *Server) handleDBMaintenanceAPI(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "解析数据库维护表单失败", http.StatusBadRequest)
		return
	}

	mode := strings.TrimSpace(r.Form.Get("mode"))
	timeout := 3 * time.Minute
	switch mode {
	case audit.DBMaintenanceVacuumFull:
		timeout = 30 * time.Minute
	case audit.DBMaintenanceCompactPayloads:
		timeout = 15 * time.Minute
	case "", "vacuum", audit.DBMaintenanceVacuumAnalyze, audit.DBMaintenanceAnalyze:
		timeout = 5 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var affectedRows int64
	var err error
	if normalizeMaintenanceMode(mode) == audit.DBMaintenanceCompactPayloads {
		affectedRows, err = s.store.CompactCapturedPayloads(ctx)
	} else {
		err = s.store.RunDBMaintenance(ctx, mode)
	}
	if err != nil {
		if errors.Is(err, audit.ErrUnsupportedDBMaintenanceMode) {
			http.Error(w, "不支持的数据库维护模式", http.StatusBadRequest)
			return
		}
		s.logger.Printf("db maintenance api failed: %v", err)
		http.Error(w, "数据库维护失败", http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, dbMaintenanceResponse{
		Mode:         normalizeMaintenanceMode(mode),
		AffectedRows: affectedRows,
		Message:      maintenanceMessage(mode, affectedRows),
	})
}

func (s *Server) handleTokens(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filters := audit.TokenDirectoryFilters{
		From:             parseLocalDateTime(query.Get("from")),
		To:               parseLocalDateTime(query.Get("to")),
		TokenFingerprint: strings.TrimSpace(query.Get("token")),
		TokenAlias:       strings.TrimSpace(query.Get("alias")),
		Model:            strings.TrimSpace(query.Get("model")),
		Page:             max(1, parseInt(query.Get("page"))),
		PageSize:         50,
	}

	result, err := s.store.ListTokens(r.Context(), filters)
	if err != nil {
		s.logger.Printf("list tokens failed: %v", err)
		http.Error(w, "查询令牌统计失败", http.StatusInternalServerError)
		return
	}

	hasNext := int64(result.Page*result.PageSize) < result.TotalCount
	s.render(w, "tokens", tokensViewData{
		Title: "令牌映射 - newapi-audit-proxy",
		Filters: tokenFiltersView{
			From:             query.Get("from"),
			To:               query.Get("to"),
			TokenFingerprint: filters.TokenFingerprint,
			TokenAlias:       filters.TokenAlias,
			Model:            filters.Model,
		},
		Result:     result,
		PrevPage:   result.Page - 1,
		NextPage:   result.Page + 1,
		HasNext:    hasNext,
		Saved:      query.Get("saved") == "1",
		Error:      query.Get("err"),
		CurrentURL: currentRequestURL(s.basePath, r, "saved", "err"),
	})
}

func (s *Server) handleTokenAliasSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "解析令牌别名表单失败", http.StatusBadRequest)
		return
	}

	redirectTo := normalizeRedirectTarget(s.basePath, r.Form.Get("redirect_to"))
	tokenFingerprint, tokenAlias, err := s.resolveTokenAliasInput(r.Form.Get("token_fingerprint"), r.Form.Get("token_value"), r.Form.Get("token_alias"))
	if err != nil {
		http.Redirect(w, r, withQueryValue(s.webPath("/tokens"), redirectTo, "err", err.Error()), http.StatusSeeOther)
		return
	}

	if err := s.store.UpsertTokenAlias(r.Context(), tokenFingerprint, tokenAlias); err != nil {
		s.logger.Printf("save token alias failed: %v", err)
		http.Redirect(w, r, withQueryValue(s.webPath("/tokens"), redirectTo, "err", friendlyAliasError(err)), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, withQueryValue(s.webPath("/tokens"), redirectTo, "saved", "1"), http.StatusSeeOther)
}

func (s *Server) handleLogDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	record, err := s.store.GetLog(r.Context(), id)
	if err != nil {
		http.Error(w, "日志不存在", http.StatusNotFound)
		return
	}

	record.UserText = audit.RedactTextContent(record.UserText)
	record.AssistantText = audit.RedactTextContent(record.AssistantText)
	record.RequestBody = audit.RedactCapturedBody(record.RequestBody, record.RequestContentType)
	record.ResponseBody = audit.RedactCapturedBody(record.ResponseBody, record.ResponseContentType)
	record.RequestJSON = audit.RedactCapturedJSON(record.RequestJSON)
	record.ResponseJSON = audit.RedactCapturedJSON(record.ResponseJSON)

	s.render(w, "detail", detailViewData{
		Title:  fmt.Sprintf("日志详情 #%d - newapi-audit-proxy", id),
		Log:    record,
		Images: extractImagesFromRecord(record),
	})
}

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.isAuthenticated(r) {
			http.Redirect(w, r, s.webPath("/login"), http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (s *Server) requireCSRF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.validCSRFToken(r) {
			http.Error(w, "CSRF 校验失败，请刷新页面后重试", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func (s *Server) validCSRFToken(r *http.Request) bool {
	token := strings.TrimSpace(r.Header.Get(csrfHeaderName))
	if token == "" {
		if err := r.ParseForm(); err != nil {
			return false
		}
		token = strings.TrimSpace(r.Form.Get(csrfFieldName))
	}
	return verifyCSRFToken(token, s.cfg.HMACSecret)
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return false
	}

	decoded, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return false
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 3 {
		return false
	}
	if parts[0] != s.cfg.AdminUsername {
		return false
	}

	expiresAt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > expiresAt {
		return false
	}

	signingInput := parts[0] + "|" + parts[1]
	expectedSig := signCookie(signingInput, s.cfg.HMACSecret)
	return subtle.ConstantTimeCompare([]byte(expectedSig), []byte(parts[2])) == 1
}

func (s *Server) issueSession(w http.ResponseWriter, r *http.Request) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	payload := fmt.Sprintf("%s|%d", s.cfg.AdminUsername, expiresAt)
	sig := signCookie(payload, s.cfg.HMACSecret)
	value := base64.RawURLEncoding.EncodeToString([]byte(payload + "|" + sig))

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     s.cookiePath(),
		HttpOnly: true,
		Secure:   s.secureCookie(r),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60,
	})
}

func (s *Server) secureCookie(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	if strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		return true
	}
	if strings.EqualFold(r.Header.Get("X-Forwarded-Ssl"), "on") {
		return true
	}
	return strings.Contains(strings.ToLower(r.Header.Get("Forwarded")), "proto=https")
}

func signCookie(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func newCSRFToken(secret string) (string, error) {
	var nonce [32]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", fmt.Errorf("generate csrf nonce: %w", err)
	}

	expiresAt := time.Now().Add(csrfTokenTTL).Unix()
	payload := fmt.Sprintf("%d|%s", expiresAt, base64.RawURLEncoding.EncodeToString(nonce[:]))
	sig := signCookie("csrf|"+payload, secret)
	return base64.RawURLEncoding.EncodeToString([]byte(payload + "|" + sig)), nil
}

func verifyCSRFToken(token, secret string) bool {
	if token == "" {
		return false
	}
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 3 {
		return false
	}
	expiresAt, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || time.Now().Unix() > expiresAt {
		return false
	}
	payload := parts[0] + "|" + parts[1]
	expectedSig := signCookie("csrf|"+payload, secret)
	return subtle.ConstantTimeCompare([]byte(expectedSig), []byte(parts[2])) == 1
}

func (s *Server) isLoginRateLimited(r *http.Request) (time.Duration, bool) {
	key := loginRateKey(r)
	now := time.Now()

	s.loginMu.Lock()
	defer s.loginMu.Unlock()

	attempt, ok := s.loginFail[key]
	if !ok {
		return 0, false
	}
	if !attempt.LockedUntil.IsZero() && now.Before(attempt.LockedUntil) {
		return attempt.LockedUntil.Sub(now), true
	}
	if !attempt.FirstFailed.IsZero() && now.Sub(attempt.FirstFailed) > loginAttemptWindow {
		delete(s.loginFail, key)
	}
	return 0, false
}

func (s *Server) recordLoginFailure(r *http.Request) {
	key := loginRateKey(r)
	now := time.Now()

	s.loginMu.Lock()
	defer s.loginMu.Unlock()

	attempt := s.loginFail[key]
	if attempt.FirstFailed.IsZero() || now.Sub(attempt.FirstFailed) > loginAttemptWindow {
		attempt = loginAttempt{FirstFailed: now}
	}
	attempt.Failures++
	if attempt.Failures >= maxLoginFailures {
		attempt.LockedUntil = now.Add(loginLockoutDuration)
		attempt.Failures = 0
		attempt.FirstFailed = now
	}
	s.loginFail[key] = attempt
}

func (s *Server) clearLoginFailures(r *http.Request) {
	s.loginMu.Lock()
	defer s.loginMu.Unlock()
	delete(s.loginFail, loginRateKey(r))
}

func loginRateKey(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

func formatDurationForHuman(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d 秒", max(1, int(d.Seconds())))
	}
	minutes := int(d.Round(time.Minute) / time.Minute)
	if minutes < 1 {
		minutes = 1
	}
	return fmt.Sprintf("%d 分钟", minutes)
}

func (s *Server) renderLoginError(w http.ResponseWriter, status int, message string) {
	s.renderWithStatus(w, status, "login", loginViewData{
		Title: "登录 - newapi-audit-proxy",
		Error: message,
	})
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	s.renderWithStatus(w, http.StatusOK, name, data)
}

func (s *Server) renderWithStatus(w http.ResponseWriter, status int, name string, data any) {
	tmpl, ok := s.templates[name]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	var err error
	data, err = s.withCSRFToken(data)
	if err != nil {
		s.logger.Printf("prepare csrf token failed: %v", err)
		http.Error(w, "生成安全令牌失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		s.logger.Printf("render template %s failed: %v", name, err)
	}
}

func (s *Server) withCSRFToken(data any) (any, error) {
	token, err := newCSRFToken(s.cfg.HMACSecret)
	if err != nil {
		return data, err
	}
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return data, nil
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return data, nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return data, nil
	}
	copyValue := reflect.New(value.Type()).Elem()
	copyValue.Set(value)
	field := copyValue.FieldByName("CSRFToken")
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(token)
		return copyValue.Interface(), nil
	}
	return data, nil
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return "-"
	}
	return value.Local().Format("2006-01-02 15:04:05")
}

func statusClass(status int) string {
	if status >= 400 {
		return "status-err"
	}
	return "status-ok"
}

func truncate(value string, maxLen int) string {
	runes := []rune(value)
	if maxLen <= 0 || len(runes) <= maxLen {
		if value == "" {
			return "-"
		}
		return value
	}
	return string(runes[:maxLen]) + "..."
}

func prettyJSON(value []byte) string {
	if len(value) == 0 {
		return "（空）"
	}
	var out bytes.Buffer
	if err := json.Indent(&out, value, "", "  "); err != nil {
		return string(value)
	}
	return out.String()
}

func prettyAny(value any) string {
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(raw)
}

func renderBody(value []byte) string {
	if len(value) == 0 {
		return "（空）"
	}
	if utf8.Valid(value) {
		return string(value)
	}
	return hex.EncodeToString(value)
}

func pathWithQuery(path, rawQuery string) string {
	if rawQuery == "" {
		return path
	}
	return path + "?" + rawQuery
}

func streamValue(value *bool) string {
	if value == nil {
		return "-"
	}
	if *value {
		return "是"
	}
	return "否"
}

func previewText(value string, maxLen int) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if maxLen <= 0 || len(runes) <= maxLen {
		return value
	}
	return string(runes[:maxLen]) + "..."
}

func emptyText(value string) string {
	if strings.TrimSpace(value) == "" {
		return "（空）"
	}
	return value
}

func toLogTextPageAPIResponse(page audit.TextPage) logTextPageAPIResponse {
	return logTextPageAPIResponse{
		Kind:       page.Kind,
		Text:       emptyText(audit.RedactTextContent(page.Text)),
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: page.TotalPages,
		TotalChars: page.TotalChars,
	}
}

func parseLogFilters(values url.Values) audit.ListFilters {
	return audit.ListFilters{
		From:             parseLocalDateTime(values.Get("from")),
		To:               parseLocalDateTime(values.Get("to")),
		TokenFingerprint: strings.TrimSpace(values.Get("token")),
		TokenAlias:       strings.TrimSpace(values.Get("alias")),
		Model:            strings.TrimSpace(values.Get("model")),
		StatusCode:       strings.TrimSpace(values.Get("status")),
		Keyword:          strings.TrimSpace(values.Get("q")),
		Page:             max(1, parseInt(values.Get("page"))),
		PageSize:         100,
	}
}

func parseTokenFilters(values url.Values) audit.TokenDirectoryFilters {
	return audit.TokenDirectoryFilters{
		From:             parseLocalDateTime(values.Get("from")),
		To:               parseLocalDateTime(values.Get("to")),
		TokenFingerprint: strings.TrimSpace(values.Get("token")),
		TokenAlias:       strings.TrimSpace(values.Get("alias")),
		Model:            strings.TrimSpace(values.Get("model")),
		Page:             1,
		PageSize:         50,
	}
}

func (s *Server) resolveTokenAliasInput(tokenFingerprint, tokenValue, tokenAlias string) (string, string, error) {
	tokenFingerprint = strings.TrimSpace(tokenFingerprint)
	tokenValue = strings.TrimSpace(tokenValue)
	tokenAlias = strings.TrimSpace(tokenAlias)

	if tokenFingerprint == "" && tokenValue != "" {
		tokenFingerprint, _ = audit.TokenMetadata(tokenValue, s.cfg.HMACSecret)
	}
	if tokenFingerprint == "" {
		return "", "", fmt.Errorf("请输入 token 原文或已记录的 token 指纹")
	}
	return tokenFingerprint, tokenAlias, nil
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.logger.Printf("write json failed: %v", err)
	}
}

func extractImages(raw []byte) []logImageAPIItem {
	if len(raw) == 0 {
		return nil
	}

	var payload any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}

	items := make([]logImageAPIItem, 0, 4)
	seen := make(map[string]struct{}, 4)
	walkImages(payload, "", &items, seen)
	for i := range items {
		if items[i].Label == "" {
			items[i].Label = fmt.Sprintf("生成图片 %d", i+1)
		}
	}
	return items
}

func extractImagesFromRecord(record audit.Record) []logImageAPIItem {
	images := extractImages(record.ResponseJSON)
	if len(images) > 0 {
		return images
	}
	return extractImages(record.ResponseBody)
}

func walkImages(node any, label string, items *[]logImageAPIItem, seen map[string]struct{}) {
	switch typed := node.(type) {
	case []any:
		for _, item := range typed {
			walkImages(item, label, items, seen)
		}
	case map[string]any:
		mimeHint := imageMimeHint(typed)
		for key, value := range typed {
			text, ok := value.(string)
			if !ok {
				continue
			}
			lowerKey := strings.ToLower(key)
			if lowerKey == "b64_json" ||
				lowerKey == "image_base64" ||
				lowerKey == "image_data" ||
				strings.Contains(lowerKey, "base64") ||
				strings.Contains(lowerKey, "b64") ||
				strings.Contains(lowerKey, "image") ||
				strings.HasPrefix(strings.ToLower(text), "data:image/") {
				if item, ok := imageItemFromValue(text, mimeHint, imageLabel(label, typed)); ok {
					if _, exists := seen[item.DataURL]; !exists {
						seen[item.DataURL] = struct{}{}
						*items = append(*items, item)
					}
				}
			}
		}
		for key, value := range typed {
			childLabel := label
			if childLabel == "" {
				childLabel = key
			}
			walkImages(value, childLabel, items, seen)
		}
	}
}

func imageItemFromValue(value, mimeHint, label string) (logImageAPIItem, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return logImageAPIItem{}, false
	}

	if strings.HasPrefix(strings.ToLower(value), "data:image/") {
		comma := strings.Index(value, ",")
		if comma > 0 {
			prefix := value[:comma]
			mime := strings.TrimPrefix(strings.SplitN(prefix, ";", 2)[0], "data:")
			return logImageAPIItem{
				Label:   label,
				Mime:    mime,
				DataURL: value,
			}, true
		}
	}

	clean := sanitizeBase64(value)
	decoded, ok := decodeBase64(clean)
	if !ok || len(decoded) == 0 {
		return logImageAPIItem{}, false
	}

	mime := http.DetectContentType(decoded)
	if !strings.HasPrefix(mime, "image/") {
		if strings.HasPrefix(mimeHint, "image/") {
			mime = mimeHint
		} else {
			return logImageAPIItem{}, false
		}
	}

	return logImageAPIItem{
		Label:   label,
		Mime:    mime,
		DataURL: "data:" + mime + ";base64," + clean,
	}, true
}

func sanitizeBase64(value string) string {
	replacer := strings.NewReplacer("\n", "", "\r", "", "\t", "", " ", "")
	return replacer.Replace(value)
}

func decodeBase64(value string) ([]byte, bool) {
	decoders := []func(string) ([]byte, error){
		base64.StdEncoding.DecodeString,
		base64.RawStdEncoding.DecodeString,
		base64.URLEncoding.DecodeString,
		base64.RawURLEncoding.DecodeString,
	}
	for _, decode := range decoders {
		decoded, err := decode(value)
		if err == nil {
			return decoded, true
		}
	}
	return nil, false
}

func imageMimeHint(payload map[string]any) string {
	for _, key := range []string{"mime_type", "mime", "content_type"} {
		if value, _ := payload[key].(string); strings.HasPrefix(strings.ToLower(value), "image/") {
			return value
		}
	}
	return ""
}

func imageLabel(parent string, payload map[string]any) string {
	for _, key := range []string{"type", "name", "role"} {
		if value, _ := payload[key].(string); value != "" {
			return value
		}
	}
	return parent
}

func buildPageURL(basePath string, filters logFiltersView, page int) string {
	values := url.Values{}
	if filters.From != "" {
		values.Set("from", filters.From)
	}
	if filters.To != "" {
		values.Set("to", filters.To)
	}
	if filters.TokenFingerprint != "" {
		values.Set("token", filters.TokenFingerprint)
	}
	if filters.TokenAlias != "" {
		values.Set("alias", filters.TokenAlias)
	}
	if filters.Model != "" {
		values.Set("model", filters.Model)
	}
	if filters.StatusCode != "" {
		values.Set("status", filters.StatusCode)
	}
	if filters.Keyword != "" {
		values.Set("q", filters.Keyword)
	}
	values.Set("page", strconv.Itoa(page))
	return joinBasePath(basePath, "/logs") + "?" + values.Encode()
}

func buildTokenPageURL(basePath string, filters tokenFiltersView, page int) string {
	values := url.Values{}
	if filters.From != "" {
		values.Set("from", filters.From)
	}
	if filters.To != "" {
		values.Set("to", filters.To)
	}
	if filters.TokenFingerprint != "" {
		values.Set("token", filters.TokenFingerprint)
	}
	if filters.TokenAlias != "" {
		values.Set("alias", filters.TokenAlias)
	}
	if filters.Model != "" {
		values.Set("model", filters.Model)
	}
	values.Set("page", strconv.Itoa(page))
	return joinBasePath(basePath, "/tokens") + "?" + values.Encode()
}

func parseLocalDateTime(value string) time.Time {
	if strings.TrimSpace(value) == "" {
		return time.Time{}
	}
	parsed, err := time.ParseInLocation("2006-01-02T15:04", value, time.Local)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func parseInt(value string) int {
	if strings.TrimSpace(value) == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func currentRequestURL(basePath string, r *http.Request, omitKeys ...string) string {
	values := r.URL.Query()
	for _, key := range omitKeys {
		values.Del(key)
	}
	encoded := values.Encode()
	if encoded == "" {
		return joinBasePath(basePath, r.URL.Path)
	}
	return joinBasePath(basePath, r.URL.Path) + "?" + encoded
}

func normalizeRedirectTarget(basePath, target string) string {
	fallback := joinBasePath(basePath, "/tokens")
	if target == "" {
		return fallback
	}
	if !strings.HasPrefix(target, "/") || strings.HasPrefix(target, "//") {
		return fallback
	}
	parsed, err := url.Parse(target)
	if err != nil || parsed.IsAbs() {
		return fallback
	}
	if parsed.Path != basePath && parsed.Path != basePath+"/" && !strings.HasPrefix(parsed.Path, basePath+"/") {
		return fallback
	}
	return parsed.String()
}

func withQueryValue(fallback, target, key, value string) string {
	parsed, err := url.Parse(target)
	if err != nil {
		return fallback
	}
	query := parsed.Query()
	query.Del("saved")
	query.Del("err")
	query.Set(key, value)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func (s *Server) webPath(path string) string {
	return joinBasePath(s.basePath, path)
}

func (s *Server) cookiePath() string {
	if s.basePath == "" || s.basePath == "/" {
		return "/"
	}
	return s.basePath
}

func joinBasePath(basePath, path string) string {
	switch {
	case path == "" || path == "/":
		if basePath == "" || basePath == "/" {
			return "/"
		}
		return basePath + "/"
	case strings.HasPrefix(path, "/"):
		if basePath == "" || basePath == "/" {
			return path
		}
		return basePath + path
	default:
		if basePath == "" || basePath == "/" {
			return "/" + path
		}
		return basePath + "/" + path
	}
}

func friendlyAliasError(err error) string {
	message := err.Error()
	switch {
	case strings.Contains(message, "duplicate key value"):
		return "该令牌代号已被其他 token 使用"
	case strings.Contains(message, "token fingerprint is required"):
		return "令牌指纹不能为空"
	default:
		return "保存令牌代号失败"
	}
}

func normalizeMaintenanceMode(mode string) string {
	mode = strings.TrimSpace(strings.ToLower(mode))
	switch mode {
	case "", "vacuum":
		return audit.DBMaintenanceVacuumAnalyze
	default:
		return mode
	}
}

func maintenanceMessage(mode string, affectedRows int64) string {
	switch normalizeMaintenanceMode(mode) {
	case audit.DBMaintenanceVacuumFull:
		return "VACUUM FULL 已完成，audit_logs 表文件已尝试回收给操作系统。执行期间该表会被锁定。"
	case audit.DBMaintenanceAnalyze:
		return "ANALYZE 已完成，统计信息已刷新。"
	case audit.DBMaintenanceCompactPayloads:
		return fmt.Sprintf("历史记录瘦身已完成，共重写 %d 条记录。建议紧接着再执行一次“强制缩盘”，把旧的大字段空间真正回收掉。", affectedRows)
	default:
		return "VACUUM ANALYZE 已完成。已整理可复用空间并刷新统计信息，但磁盘文件未必会立刻缩小。"
	}
}
