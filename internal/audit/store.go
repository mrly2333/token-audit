package audit

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const auditLogJoinSQL = `
	FROM audit_logs l
	LEFT JOIN token_aliases ta
		ON ta.token_fingerprint = l.token_fingerprint
`

const (
	asyncInsertTimeout  = 15 * time.Second
	logListPreviewChars = 160
	logTextPageChars    = 4000
)

const (
	DBMaintenanceVacuumAnalyze   = "vacuum_analyze"
	DBMaintenanceVacuumFull      = "vacuum_full"
	DBMaintenanceAnalyze         = "analyze"
	DBMaintenanceCompactPayloads = "compact_payloads"
)

var (
	ErrUnsupportedTextKind          = errors.New("unsupported text kind")
	ErrUnsupportedDBMaintenanceMode = errors.New("unsupported db maintenance mode")
	ErrAuditLogNotFound             = errors.New("audit log not found")
)

type Store struct {
	pool   *pgxpool.Pool
	logger *log.Logger
}

func NewStore(pool *pgxpool.Pool, logger *log.Logger) *Store {
	return &Store{pool: pool, logger: logger}
}

func (s *Store) InsertAsync(record Record) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), asyncInsertTimeout)
		defer cancel()

		if err := s.Insert(ctx, record); err != nil {
			s.logger.Printf("audit insert failed: %v", err)
		}
	}()
}

func (s *Store) Insert(ctx context.Context, record Record) error {
	requestHeaders, err := json.Marshal(record.RequestHeaders)
	if err != nil {
		return fmt.Errorf("marshal request headers: %w", err)
	}
	responseHeaders, err := json.Marshal(record.ResponseHeaders)
	if err != nil {
		return fmt.Errorf("marshal response headers: %w", err)
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO audit_logs (
			started_at,
			finished_at,
			duration_ms,
			method,
			path,
			query_string,
			remote_addr,
			request_host,
			upstream_base,
			status_code,
			error_text,
			is_capture_path,
			response_is_sse,
			token_fingerprint,
			token_preview,
			model,
			stream,
			request_headers,
			response_headers,
			request_content_type,
			response_content_type,
			request_body,
			response_body,
			request_json,
			response_json,
			user_text,
			assistant_text,
			usage_json,
			prompt_tokens,
			completion_tokens,
			total_tokens,
			request_bytes,
			response_bytes,
			request_truncated,
			response_truncated
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18::jsonb, $19::jsonb, $20,
			$21, $22, $23, $24::jsonb, $25::jsonb, $26, $27, $28::jsonb, $29, $30,
			$31, $32, $33, $34, $35
		)
	`,
		record.StartedAt,
		nullTime(record.FinishedAt),
		record.DurationMS,
		record.Method,
		record.Path,
		record.QueryString,
		record.RemoteAddr,
		record.RequestHost,
		record.UpstreamBase,
		nullInt(record.StatusCode),
		nullString(record.ErrorText),
		record.IsCapturePath,
		record.ResponseIsSSE,
		nullString(record.TokenFingerprint),
		nullString(record.TokenPreview),
		nullString(record.Model),
		record.Stream,
		string(requestHeaders),
		string(responseHeaders),
		nullString(record.RequestContentType),
		nullString(record.ResponseContentType),
		nullBytes(record.RequestBody),
		nullBytes(record.ResponseBody),
		jsonParam(record.RequestJSON),
		jsonParam(record.ResponseJSON),
		nullString(record.UserText),
		nullString(record.AssistantText),
		jsonParam(record.UsageJSON),
		record.PromptTokens,
		record.CompletionTokens,
		record.TotalTokens,
		record.RequestBytes,
		record.ResponseBytes,
		record.RequestTruncated,
		record.ResponseTruncated,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

func (s *Store) DashboardStats(ctx context.Context, now time.Time) (DashboardStats, error) {
	stats := DashboardStats{}
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if err := s.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) AS total_requests,
			COUNT(*) FILTER (WHERE started_at >= $1) AS today_requests,
			COUNT(*) FILTER (WHERE status_code >= 400 OR error_text IS NOT NULL) AS error_count,
			COUNT(DISTINCT token_fingerprint) FILTER (
				WHERE COALESCE(token_fingerprint, '') <> ''
			) AS distinct_tokens,
			COUNT(DISTINCT token_fingerprint) FILTER (
				WHERE started_at >= $1 AND COALESCE(token_fingerprint, '') <> ''
			) AS today_distinct_tokens,
			COALESCE(SUM(prompt_tokens), 0) AS total_prompt_tokens,
			COALESCE(SUM(completion_tokens), 0) AS total_completion_tokens,
			COALESCE(SUM(total_tokens), 0) AS total_tokens,
			COALESCE(SUM(prompt_tokens) FILTER (WHERE started_at >= $1), 0) AS today_prompt_tokens,
			COALESCE(SUM(completion_tokens) FILTER (WHERE started_at >= $1), 0) AS today_completion_tokens,
			COALESCE(SUM(total_tokens) FILTER (WHERE started_at >= $1), 0) AS today_total_tokens
		FROM audit_logs
		WHERE is_capture_path = TRUE
	`, dayStart).Scan(
		&stats.TotalRequests,
		&stats.TodayRequests,
		&stats.ErrorCount,
		&stats.DistinctTokens,
		&stats.TodayDistinctTokens,
		&stats.TotalPromptTokens,
		&stats.TotalCompletionTokens,
		&stats.TotalTokens,
		&stats.TodayPromptTokens,
		&stats.TodayCompletionTokens,
		&stats.TodayTotalTokens,
	); err != nil {
		return DashboardStats{}, fmt.Errorf("query dashboard totals: %w", err)
	}

	rows, err := s.pool.Query(ctx, `
		SELECT
			COALESCE(l.token_fingerprint, '') AS token_fingerprint,
			COALESCE(MAX(l.token_preview), '') AS token_preview,
			COALESCE(MAX(ta.token_alias), '') AS token_alias,
			COUNT(*) AS request_count,
			COUNT(*) FILTER (WHERE l.status_code >= 400 OR l.error_text IS NOT NULL) AS error_count,
			COALESCE(SUM(l.prompt_tokens), 0) AS prompt_tokens,
			COALESCE(SUM(l.completion_tokens), 0) AS completion_tokens,
			COALESCE(SUM(l.total_tokens), 0) AS total_tokens,
			MAX(l.started_at) AS last_seen
		`+auditLogJoinSQL+`
		WHERE l.is_capture_path = TRUE
			AND COALESCE(l.token_fingerprint, '') <> ''
		GROUP BY COALESCE(l.token_fingerprint, '')
		ORDER BY total_tokens DESC, request_count DESC, last_seen DESC
		LIMIT 50
	`)
	if err != nil {
		return DashboardStats{}, fmt.Errorf("query token stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item TokenStats
		if err := rows.Scan(
			&item.TokenFingerprint,
			&item.TokenPreview,
			&item.TokenAlias,
			&item.RequestCount,
			&item.ErrorCount,
			&item.PromptTokens,
			&item.CompletionTokens,
			&item.TotalTokens,
			&item.LastSeen,
		); err != nil {
			return DashboardStats{}, fmt.Errorf("scan token stats: %w", err)
		}
		stats.TokenGroups = append(stats.TokenGroups, item)
	}
	if err := rows.Err(); err != nil {
		return DashboardStats{}, fmt.Errorf("iterate token stats: %w", err)
	}

	modelRows, err := s.pool.Query(ctx, `
		SELECT
			COALESCE(l.model, '') AS model,
			COUNT(*) AS request_count,
			COUNT(*) FILTER (WHERE l.status_code >= 400 OR l.error_text IS NOT NULL) AS error_count,
			COALESCE(SUM(l.prompt_tokens), 0) AS prompt_tokens,
			COALESCE(SUM(l.completion_tokens), 0) AS completion_tokens,
			COALESCE(SUM(l.total_tokens), 0) AS total_tokens,
			MAX(l.started_at) AS last_seen
		FROM audit_logs l
		WHERE l.is_capture_path = TRUE
			AND COALESCE(l.model, '') <> ''
		GROUP BY COALESCE(l.model, '')
		ORDER BY total_tokens DESC, request_count DESC, last_seen DESC
		LIMIT 20
	`)
	if err != nil {
		return DashboardStats{}, fmt.Errorf("query model stats: %w", err)
	}
	defer modelRows.Close()

	for modelRows.Next() {
		var item ModelStats
		if err := modelRows.Scan(
			&item.Model,
			&item.RequestCount,
			&item.ErrorCount,
			&item.PromptTokens,
			&item.CompletionTokens,
			&item.TotalTokens,
			&item.LastSeen,
		); err != nil {
			return DashboardStats{}, fmt.Errorf("scan model stats: %w", err)
		}
		stats.ModelGroups = append(stats.ModelGroups, item)
	}
	if err := modelRows.Err(); err != nil {
		return DashboardStats{}, fmt.Errorf("iterate model stats: %w", err)
	}

	return stats, nil
}

func (s *Store) ListLogs(ctx context.Context, filters ListFilters) (ListResult, error) {
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 || filters.PageSize > 100 {
		filters.PageSize = 50
	}

	whereSQL, args := buildListWhere(filters)

	var total int64
	if err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		`+auditLogJoinSQL+`
		WHERE `+whereSQL, args...).Scan(&total); err != nil {
		return ListResult{}, fmt.Errorf("count audit logs: %w", err)
	}

	if total > 0 {
		totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))
		if filters.Page > totalPages {
			filters.Page = totalPages
		}
	}

	offset := (filters.Page - 1) * filters.PageSize
	listArgs := append(append([]any(nil), args...), filters.PageSize, offset)
	limitIndex := len(args) + 1
	offsetIndex := len(args) + 2

	rows, err := s.pool.Query(ctx, fmt.Sprintf(`
		SELECT
			l.id,
			l.started_at,
			l.method,
			l.path,
			COALESCE(l.status_code, 0) AS status_code,
			COALESCE(l.model, '') AS model,
			COALESCE(l.token_fingerprint, '') AS token_fingerprint,
			COALESCE(l.token_preview, '') AS token_preview,
			COALESCE(ta.token_alias, '') AS token_alias,
			COALESCE(l.prompt_tokens, 0) AS prompt_tokens,
			COALESCE(l.completion_tokens, 0) AS completion_tokens,
			COALESCE(l.total_tokens, 0) AS total_tokens,
			COALESCE(substring(COALESCE(l.user_text, '') FROM 1 FOR %d), '') AS user_text,
			COALESCE(substring(COALESCE(l.assistant_text, '') FROM 1 FOR %d), '') AS assistant_text,
			COALESCE(l.error_text, '') AS error_text
		%s
		WHERE %s
		ORDER BY l.started_at DESC
		LIMIT $%d OFFSET $%d
	`, logListPreviewChars, logListPreviewChars, auditLogJoinSQL, whereSQL, limitIndex, offsetIndex), listArgs...)
	if err != nil {
		return ListResult{}, fmt.Errorf("query audit logs: %w", err)
	}
	defer rows.Close()

	result := ListResult{
		Items:      make([]LogListItem, 0, filters.PageSize),
		TotalCount: total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
	}

	for rows.Next() {
		var item LogListItem
		if err := rows.Scan(
			&item.ID,
			&item.StartedAt,
			&item.Method,
			&item.Path,
			&item.StatusCode,
			&item.Model,
			&item.TokenFingerprint,
			&item.TokenPreview,
			&item.TokenAlias,
			&item.PromptTokens,
			&item.CompletionTokens,
			&item.TotalTokens,
			&item.UserText,
			&item.AssistantText,
			&item.ErrorText,
		); err != nil {
			return ListResult{}, fmt.Errorf("scan audit log: %w", err)
		}
		result.Items = append(result.Items, item)
	}
	if err := rows.Err(); err != nil {
		return ListResult{}, fmt.Errorf("iterate audit logs: %w", err)
	}

	return result, nil
}

func (s *Store) LogsVersion(ctx context.Context, filters ListFilters) (string, error) {
	whereSQL, args := buildListWhere(filters)

	var total int64
	var maxID sql.NullInt64
	var aliasUpdatedAt sql.NullTime
	if err := s.pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			MAX(l.id),
			MAX(COALESCE(ta.updated_at, 'epoch'::timestamptz))
		`+auditLogJoinSQL+`
		WHERE `+whereSQL, args...).Scan(&total, &maxID, &aliasUpdatedAt); err != nil {
		return "", fmt.Errorf("query audit log version: %w", err)
	}

	var lastID int64
	if maxID.Valid {
		lastID = maxID.Int64
	}
	var aliasUnix int64
	if aliasUpdatedAt.Valid {
		aliasUnix = aliasUpdatedAt.Time.Unix()
	}
	return fmt.Sprintf("%d:%d:%d", total, lastID, aliasUnix), nil
}

func (s *Store) GetLogCore(ctx context.Context, id int64) (LogCore, error) {
	var record Record
	var core LogCore
	var model sql.NullString
	var errorText sql.NullString
	var tokenFingerprint sql.NullString
	var tokenPreview sql.NullString
	var tokenAlias sql.NullString
	var userText sql.NullString
	var userTextChars int64
	var assistantText sql.NullString
	var assistantTextChars int64
	var statusCode sql.NullInt64
	var stream sql.NullBool

	err := s.pool.QueryRow(ctx, `
		SELECT
			l.id,
			l.started_at,
			l.duration_ms,
			l.method,
			l.path,
			l.query_string,
			l.status_code,
			l.error_text,
			l.response_is_sse,
			l.token_fingerprint,
			l.token_preview,
			ta.token_alias,
			l.model,
			l.stream,
			COALESCE(substring(COALESCE(l.user_text, '') FROM 1 FOR $2), ''),
			char_length(COALESCE(l.user_text, '')),
			COALESCE(substring(COALESCE(l.assistant_text, '') FROM 1 FOR $2), ''),
			char_length(COALESCE(l.assistant_text, '')),
			l.prompt_tokens,
			l.completion_tokens,
			l.total_tokens,
			l.request_bytes,
			l.response_bytes,
			l.request_truncated,
			l.response_truncated
		`+auditLogJoinSQL+`
		WHERE l.id = $1
	`, id, logTextPageChars).Scan(
		&record.ID,
		&record.StartedAt,
		&record.DurationMS,
		&record.Method,
		&record.Path,
		&record.QueryString,
		&statusCode,
		&errorText,
		&record.ResponseIsSSE,
		&tokenFingerprint,
		&tokenPreview,
		&tokenAlias,
		&model,
		&stream,
		&userText,
		&userTextChars,
		&assistantText,
		&assistantTextChars,
		&record.PromptTokens,
		&record.CompletionTokens,
		&record.TotalTokens,
		&record.RequestBytes,
		&record.ResponseBytes,
		&record.RequestTruncated,
		&record.ResponseTruncated,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return LogCore{}, err
		}
		return LogCore{}, fmt.Errorf("get audit log core: %w", err)
	}

	if statusCode.Valid {
		record.StatusCode = int(statusCode.Int64)
	}
	if stream.Valid {
		streamValue := stream.Bool
		record.Stream = &streamValue
	}
	record.ErrorText = nullStringValue(errorText)
	record.TokenFingerprint = nullStringValue(tokenFingerprint)
	record.TokenPreview = nullStringValue(tokenPreview)
	record.TokenAlias = nullStringValue(tokenAlias)
	record.Model = nullStringValue(model)
	record.UserText = nullStringValue(userText)
	record.AssistantText = nullStringValue(assistantText)

	core.Record = record
	core.UserText = buildTextPage("user", record.UserText, 1, logTextPageChars, userTextChars)
	core.AssistantText = buildTextPage("assistant", record.AssistantText, 1, logTextPageChars, assistantTextChars)

	return core, nil
}

func (s *Store) GetLogTextPage(ctx context.Context, id int64, kind string, page int) (TextPage, error) {
	column, err := textColumn(kind)
	if err != nil {
		return TextPage{}, err
	}

	textPage, err := s.getLogTextPageSQL(ctx, id, kind, column, page)
	if err == nil || errors.Is(err, ErrAuditLogNotFound) {
		return textPage, err
	}

	fallbackPage, fallbackErr := s.getLogTextPageFallback(ctx, id, kind, page)
	if fallbackErr == nil {
		return fallbackPage, nil
	}

	return TextPage{}, fmt.Errorf("%w; fallback full-text pagination failed: %v", err, fallbackErr)
}

func (s *Store) getLogTextPageSQL(ctx context.Context, id int64, kind, column string, page int) (TextPage, error) {
	if page <= 0 {
		page = 1
	}

	start := (page-1)*logTextPageChars + 1
	var text sql.NullString
	var totalChars int64
	if err := s.pool.QueryRow(ctx, `
		SELECT COALESCE(substring(COALESCE(`+column+`, '') FROM $2 FOR $3), '')
			, char_length(COALESCE(`+column+`, ''))
		FROM audit_logs l
		WHERE l.id = $1
	`, id, start, logTextPageChars).Scan(&text, &totalChars); err != nil {
		if err == pgx.ErrNoRows {
			return TextPage{}, ErrAuditLogNotFound
		}
		return TextPage{}, fmt.Errorf("get %s text page: %w", kind, err)
	}

	totalPages := totalTextPages(totalChars, logTextPageChars)
	if page > totalPages {
		page = totalPages
		start = (page-1)*logTextPageChars + 1
		if err := s.pool.QueryRow(ctx, `
			SELECT COALESCE(substring(COALESCE(`+column+`, '') FROM $2 FOR $3), '')
			FROM audit_logs l
			WHERE l.id = $1
		`, id, start, logTextPageChars).Scan(&text); err != nil {
			if err == pgx.ErrNoRows {
				return TextPage{}, ErrAuditLogNotFound
			}
			return TextPage{}, fmt.Errorf("reload %s text page after clamp: %w", kind, err)
		}
	}

	return buildTextPage(kind, nullStringValue(text), page, logTextPageChars, totalChars), nil
}

func (s *Store) getLogTextPageFallback(ctx context.Context, id int64, kind string, page int) (TextPage, error) {
	record, err := s.GetLog(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return TextPage{}, ErrAuditLogNotFound
		}
		return TextPage{}, fmt.Errorf("load full audit log for %s text fallback: %w", kind, err)
	}

	switch kind {
	case "user":
		return buildTextPageFromFullText(kind, record.UserText, page, logTextPageChars), nil
	case "assistant":
		return buildTextPageFromFullText(kind, record.AssistantText, page, logTextPageChars), nil
	default:
		return TextPage{}, ErrUnsupportedTextKind
	}
}

func (s *Store) GetLog(ctx context.Context, id int64) (Record, error) {
	var record Record
	var requestHeadersJSON []byte
	var responseHeadersJSON []byte
	var model sql.NullString
	var errorText sql.NullString
	var tokenFingerprint sql.NullString
	var tokenPreview sql.NullString
	var tokenAlias sql.NullString
	var requestContentType sql.NullString
	var responseContentType sql.NullString
	var userText sql.NullString
	var assistantText sql.NullString
	var statusCode sql.NullInt64
	var finishedAt sql.NullTime
	var stream sql.NullBool

	err := s.pool.QueryRow(ctx, `
		SELECT
			l.id,
			l.started_at,
			l.finished_at,
			l.duration_ms,
			l.method,
			l.path,
			l.query_string,
			l.remote_addr,
			l.request_host,
			l.upstream_base,
			l.status_code,
			l.error_text,
			l.is_capture_path,
			l.response_is_sse,
			l.token_fingerprint,
			l.token_preview,
			ta.token_alias,
			l.model,
			l.stream,
			l.request_headers,
			l.response_headers,
			l.request_content_type,
			l.response_content_type,
			l.request_body,
			l.response_body,
			l.request_json,
			l.response_json,
			l.user_text,
			l.assistant_text,
			l.usage_json,
			l.prompt_tokens,
			l.completion_tokens,
			l.total_tokens,
			l.request_bytes,
			l.response_bytes,
			l.request_truncated,
			l.response_truncated
		`+auditLogJoinSQL+`
		WHERE l.id = $1
	`, id).Scan(
		&record.ID,
		&record.StartedAt,
		&finishedAt,
		&record.DurationMS,
		&record.Method,
		&record.Path,
		&record.QueryString,
		&record.RemoteAddr,
		&record.RequestHost,
		&record.UpstreamBase,
		&statusCode,
		&errorText,
		&record.IsCapturePath,
		&record.ResponseIsSSE,
		&tokenFingerprint,
		&tokenPreview,
		&tokenAlias,
		&model,
		&stream,
		&requestHeadersJSON,
		&responseHeadersJSON,
		&requestContentType,
		&responseContentType,
		&record.RequestBody,
		&record.ResponseBody,
		&record.RequestJSON,
		&record.ResponseJSON,
		&userText,
		&assistantText,
		&record.UsageJSON,
		&record.PromptTokens,
		&record.CompletionTokens,
		&record.TotalTokens,
		&record.RequestBytes,
		&record.ResponseBytes,
		&record.RequestTruncated,
		&record.ResponseTruncated,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Record{}, err
		}
		return Record{}, fmt.Errorf("get audit log: %w", err)
	}

	if finishedAt.Valid {
		record.FinishedAt = finishedAt.Time
	}
	if statusCode.Valid {
		record.StatusCode = int(statusCode.Int64)
	}
	if stream.Valid {
		streamValue := stream.Bool
		record.Stream = &streamValue
	}
	record.ErrorText = nullStringValue(errorText)
	record.TokenFingerprint = nullStringValue(tokenFingerprint)
	record.TokenPreview = nullStringValue(tokenPreview)
	record.TokenAlias = nullStringValue(tokenAlias)
	record.Model = nullStringValue(model)
	record.RequestContentType = nullStringValue(requestContentType)
	record.ResponseContentType = nullStringValue(responseContentType)
	record.UserText = nullStringValue(userText)
	record.AssistantText = nullStringValue(assistantText)

	if err := json.Unmarshal(requestHeadersJSON, &record.RequestHeaders); err != nil {
		return Record{}, fmt.Errorf("decode request headers: %w", err)
	}
	if err := json.Unmarshal(responseHeadersJSON, &record.ResponseHeaders); err != nil {
		return Record{}, fmt.Errorf("decode response headers: %w", err)
	}

	return record, nil
}

func (s *Store) ListTokens(ctx context.Context, filters TokenDirectoryFilters) (TokenDirectoryResult, error) {
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 || filters.PageSize > 100 {
		filters.PageSize = 50
	}

	whereSQL, args := buildTokenWhere(filters)

	var total int64
	if err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM (
			SELECT COALESCE(l.token_fingerprint, '')
			`+auditLogJoinSQL+`
			WHERE `+whereSQL+`
			GROUP BY COALESCE(l.token_fingerprint, '')
		) grouped
	`, args...).Scan(&total); err != nil {
		return TokenDirectoryResult{}, fmt.Errorf("count token groups: %w", err)
	}

	offset := (filters.Page - 1) * filters.PageSize
	listArgs := append(append([]any(nil), args...), filters.PageSize, offset)
	limitIndex := len(args) + 1
	offsetIndex := len(args) + 2

	rows, err := s.pool.Query(ctx, fmt.Sprintf(`
		SELECT
			COALESCE(l.token_fingerprint, '') AS token_fingerprint,
			COALESCE(MAX(l.token_preview), '') AS token_preview,
			COALESCE(MAX(ta.token_alias), '') AS token_alias,
			COUNT(*) AS request_count,
			COUNT(*) FILTER (WHERE l.status_code >= 400 OR l.error_text IS NOT NULL) AS error_count,
			COALESCE(SUM(l.prompt_tokens), 0) AS prompt_tokens,
			COALESCE(SUM(l.completion_tokens), 0) AS completion_tokens,
			COALESCE(SUM(l.total_tokens), 0) AS total_tokens,
			MIN(l.started_at) AS first_seen,
			MAX(l.started_at) AS last_seen
		%s
		WHERE %s
		GROUP BY COALESCE(l.token_fingerprint, '')
		ORDER BY total_tokens DESC, request_count DESC, last_seen DESC
		LIMIT $%d OFFSET $%d
	`, auditLogJoinSQL, whereSQL, limitIndex, offsetIndex), listArgs...)
	if err != nil {
		return TokenDirectoryResult{}, fmt.Errorf("query token groups: %w", err)
	}
	defer rows.Close()

	result := TokenDirectoryResult{
		Items:      make([]TokenDirectoryItem, 0, filters.PageSize),
		TotalCount: total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
	}

	for rows.Next() {
		var item TokenDirectoryItem
		if err := rows.Scan(
			&item.TokenFingerprint,
			&item.TokenPreview,
			&item.TokenAlias,
			&item.RequestCount,
			&item.ErrorCount,
			&item.PromptTokens,
			&item.CompletionTokens,
			&item.TotalTokens,
			&item.FirstSeen,
			&item.LastSeen,
		); err != nil {
			return TokenDirectoryResult{}, fmt.Errorf("scan token group: %w", err)
		}
		result.Items = append(result.Items, item)
	}
	if err := rows.Err(); err != nil {
		return TokenDirectoryResult{}, fmt.Errorf("iterate token groups: %w", err)
	}

	return result, nil
}

func (s *Store) FilterOptions(ctx context.Context) (FilterOptions, error) {
	options := FilterOptions{}

	modelRows, err := s.pool.Query(ctx, `
		SELECT COALESCE(l.model, '') AS model
		FROM audit_logs l
		WHERE l.is_capture_path = TRUE
			AND COALESCE(l.model, '') <> ''
		GROUP BY COALESCE(l.model, '')
		ORDER BY COUNT(*) DESC, MAX(l.started_at) DESC
		LIMIT 80
	`)
	if err != nil {
		return FilterOptions{}, fmt.Errorf("query model filter options: %w", err)
	}
	defer modelRows.Close()

	for modelRows.Next() {
		var value string
		if err := modelRows.Scan(&value); err != nil {
			return FilterOptions{}, fmt.Errorf("scan model filter option: %w", err)
		}
		options.Models = append(options.Models, value)
	}
	if err := modelRows.Err(); err != nil {
		return FilterOptions{}, fmt.Errorf("iterate model filter options: %w", err)
	}

	aliasRows, err := s.pool.Query(ctx, `
		SELECT ta.token_alias
		FROM token_aliases ta
		WHERE COALESCE(ta.token_alias, '') <> ''
		ORDER BY ta.updated_at DESC, ta.token_alias ASC
		LIMIT 80
	`)
	if err != nil {
		return FilterOptions{}, fmt.Errorf("query token alias filter options: %w", err)
	}
	defer aliasRows.Close()

	for aliasRows.Next() {
		var value string
		if err := aliasRows.Scan(&value); err != nil {
			return FilterOptions{}, fmt.Errorf("scan token alias filter option: %w", err)
		}
		options.TokenAliases = append(options.TokenAliases, value)
	}
	if err := aliasRows.Err(); err != nil {
		return FilterOptions{}, fmt.Errorf("iterate token alias filter options: %w", err)
	}

	tokenRows, err := s.pool.Query(ctx, `
		SELECT COALESCE(l.token_fingerprint, '') AS token_fingerprint
		FROM audit_logs l
		WHERE l.is_capture_path = TRUE
			AND COALESCE(l.token_fingerprint, '') <> ''
		GROUP BY COALESCE(l.token_fingerprint, '')
		ORDER BY COUNT(*) DESC, MAX(l.started_at) DESC
		LIMIT 80
	`)
	if err != nil {
		return FilterOptions{}, fmt.Errorf("query token fingerprint filter options: %w", err)
	}
	defer tokenRows.Close()

	for tokenRows.Next() {
		var value string
		if err := tokenRows.Scan(&value); err != nil {
			return FilterOptions{}, fmt.Errorf("scan token fingerprint filter option: %w", err)
		}
		options.TokenFingerprints = append(options.TokenFingerprints, value)
	}
	if err := tokenRows.Err(); err != nil {
		return FilterOptions{}, fmt.Errorf("iterate token fingerprint filter options: %w", err)
	}

	statusRows, err := s.pool.Query(ctx, `
		SELECT DISTINCT CAST(l.status_code AS TEXT) AS status_code
		FROM audit_logs l
		WHERE l.is_capture_path = TRUE
			AND l.status_code IS NOT NULL
		ORDER BY status_code ASC
	`)
	if err != nil {
		return FilterOptions{}, fmt.Errorf("query status filter options: %w", err)
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var value string
		if err := statusRows.Scan(&value); err != nil {
			return FilterOptions{}, fmt.Errorf("scan status filter option: %w", err)
		}
		options.StatusCodes = append(options.StatusCodes, value)
	}
	if err := statusRows.Err(); err != nil {
		return FilterOptions{}, fmt.Errorf("iterate status filter options: %w", err)
	}

	return options, nil
}

func (s *Store) DatabaseStats(ctx context.Context, now time.Time) (DatabaseStats, error) {
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var stats DatabaseStats
	var lastVacuum sql.NullTime
	var lastAutovacuum sql.NullTime
	var lastAnalyze sql.NullTime
	var lastAutoanalyze sql.NullTime
	if err := s.pool.QueryRow(ctx, `
		WITH rel AS (
			SELECT COALESCE(to_regclass('public.audit_logs'), to_regclass('audit_logs')) AS relid
		),
		counts AS (
			SELECT
				COUNT(*) FILTER (WHERE is_capture_path = TRUE) AS total_rows,
				COUNT(*) FILTER (WHERE is_capture_path = TRUE AND started_at >= $1) AS today_rows
			FROM audit_logs
		),
		sizes AS (
			SELECT
				pg_database_size(current_database()) AS database_size,
				pg_size_pretty(pg_database_size(current_database())) AS database_pretty,
				COALESCE(pg_total_relation_size(rel.relid), 0) AS audit_total_size,
				pg_size_pretty(COALESCE(pg_total_relation_size(rel.relid), 0)) AS audit_total_pretty,
				COALESCE(pg_relation_size(rel.relid), 0) AS audit_table_size,
				pg_size_pretty(COALESCE(pg_relation_size(rel.relid), 0)) AS audit_table_pretty,
				COALESCE(pg_indexes_size(rel.relid), 0) AS audit_index_size,
				pg_size_pretty(COALESCE(pg_indexes_size(rel.relid), 0)) AS audit_index_pretty,
				COALESCE(pg_total_relation_size(c.reltoastrelid), 0) AS audit_toast_size,
				pg_size_pretty(COALESCE(pg_total_relation_size(c.reltoastrelid), 0)) AS audit_toast_pretty
			FROM rel
			LEFT JOIN pg_class c
				ON c.oid = rel.relid
		),
		vacuum_stats AS (
			SELECT
				COALESCE(st.n_live_tup, 0) AS live_tuples,
				COALESCE(st.n_dead_tup, 0) AS dead_tuples,
				st.last_vacuum,
				st.last_autovacuum,
				st.last_analyze,
				st.last_autoanalyze
			FROM rel
			LEFT JOIN pg_stat_user_tables st
				ON st.relid = rel.relid
		)
		SELECT
			counts.total_rows,
			counts.today_rows,
			sizes.database_size,
			sizes.database_pretty,
			sizes.audit_total_size,
			sizes.audit_total_pretty,
			sizes.audit_table_size,
			sizes.audit_table_pretty,
			sizes.audit_index_size,
			sizes.audit_index_pretty,
			sizes.audit_toast_size,
			sizes.audit_toast_pretty,
			vacuum_stats.live_tuples,
			vacuum_stats.dead_tuples,
			vacuum_stats.last_vacuum,
			vacuum_stats.last_autovacuum,
			vacuum_stats.last_analyze,
			vacuum_stats.last_autoanalyze
		FROM counts
		CROSS JOIN sizes
		CROSS JOIN vacuum_stats
	`, dayStart).Scan(
		&stats.TotalRows,
		&stats.TodayRows,
		&stats.DatabaseSize,
		&stats.DatabasePretty,
		&stats.AuditTotalSize,
		&stats.AuditTotalPretty,
		&stats.AuditTableSize,
		&stats.AuditTablePretty,
		&stats.AuditIndexSize,
		&stats.AuditIndexPretty,
		&stats.AuditToastSize,
		&stats.AuditToastPretty,
		&stats.LiveTuples,
		&stats.DeadTuples,
		&lastVacuum,
		&lastAutovacuum,
		&lastAnalyze,
		&lastAutoanalyze,
	); err != nil {
		return DatabaseStats{}, fmt.Errorf("query database stats: %w", err)
	}
	if lastVacuum.Valid {
		stats.LastVacuum = lastVacuum.Time
	}
	if lastAutovacuum.Valid {
		stats.LastAutovacuum = lastAutovacuum.Time
	}
	if lastAnalyze.Valid {
		stats.LastAnalyze = lastAnalyze.Time
	}
	if lastAutoanalyze.Valid {
		stats.LastAutoanalyze = lastAutoanalyze.Time
	}

	return stats, nil
}

func (s *Store) RunDBMaintenance(ctx context.Context, mode string) error {
	mode = strings.TrimSpace(strings.ToLower(mode))

	var statement string
	switch mode {
	case "", DBMaintenanceVacuumAnalyze, "vacuum":
		statement = "VACUUM (ANALYZE) audit_logs"
	case DBMaintenanceVacuumFull:
		statement = "VACUUM (FULL, ANALYZE) audit_logs"
	case DBMaintenanceAnalyze:
		statement = "ANALYZE audit_logs"
	case DBMaintenanceCompactPayloads:
		_, err := s.CompactCapturedPayloads(ctx)
		return err
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedDBMaintenanceMode, mode)
	}

	if _, err := s.pool.Exec(ctx, statement); err != nil {
		return fmt.Errorf("run %s on audit_logs: %w", mode, err)
	}

	return nil
}

func (s *Store) CompactCapturedPayloads(ctx context.Context) (int64, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT
			id,
			request_content_type,
			response_content_type,
			request_body,
			response_body,
			request_json,
			response_json,
			user_text,
			assistant_text
		FROM audit_logs
		WHERE is_capture_path = TRUE
	`)
	if err != nil {
		return 0, fmt.Errorf("query compact payload candidates: %w", err)
	}
	defer rows.Close()

	var updatedRows int64
	for rows.Next() {
		var id int64
		var requestContentType sql.NullString
		var responseContentType sql.NullString
		var requestBody []byte
		var responseBody []byte
		var requestJSON []byte
		var responseJSON []byte
		var userText sql.NullString
		var assistantText sql.NullString

		if err := rows.Scan(
			&id,
			&requestContentType,
			&responseContentType,
			&requestBody,
			&responseBody,
			&requestJSON,
			&responseJSON,
			&userText,
			&assistantText,
		); err != nil {
			return updatedRows, fmt.Errorf("scan compact payload candidate: %w", err)
		}

		redactedRequestBody := RedactCapturedBody(requestBody, requestContentType.String)
		redactedResponseBody := RedactCapturedBody(responseBody, responseContentType.String)
		redactedRequestJSON := RedactCapturedJSON(requestJSON)
		redactedResponseJSON := RedactCapturedJSON(responseJSON)
		redactedUserText := RedactTextContent(userText.String)
		redactedAssistantText := RedactTextContent(assistantText.String)

		if bytes.Equal(requestBody, redactedRequestBody) &&
			bytes.Equal(responseBody, redactedResponseBody) &&
			bytes.Equal(requestJSON, redactedRequestJSON) &&
			bytes.Equal(responseJSON, redactedResponseJSON) &&
			userText.String == redactedUserText &&
			assistantText.String == redactedAssistantText {
			continue
		}

		if _, err := s.pool.Exec(ctx, `
			UPDATE audit_logs
			SET
				request_body = $2,
				response_body = $3,
				request_json = $4::jsonb,
				response_json = $5::jsonb,
				user_text = $6,
				assistant_text = $7
			WHERE id = $1
		`,
			id,
			nullBytes(redactedRequestBody),
			nullBytes(redactedResponseBody),
			jsonParam(redactedRequestJSON),
			jsonParam(redactedResponseJSON),
			nullString(redactedUserText),
			nullString(redactedAssistantText),
		); err != nil {
			return updatedRows, fmt.Errorf("update compact payload candidate %d: %w", id, err)
		}

		updatedRows++
	}
	if err := rows.Err(); err != nil {
		return updatedRows, fmt.Errorf("iterate compact payload candidates: %w", err)
	}

	return updatedRows, nil
}

func (s *Store) DeleteLogs(ctx context.Context, filters DeleteLogsFilters) (int64, error) {
	if !hasDeleteCriteria(filters) {
		return 0, fmt.Errorf("at least one delete filter is required")
	}

	whereSQL, args := buildDeleteWhere(filters)
	commandTag, err := s.pool.Exec(ctx, `
		DELETE FROM audit_logs
		WHERE id IN (
			SELECT l.id
			`+auditLogJoinSQL+`
			WHERE `+whereSQL+`
		)
	`, args...)
	if err != nil {
		return 0, fmt.Errorf("delete audit logs: %w", err)
	}

	return commandTag.RowsAffected(), nil
}

func (s *Store) UpsertTokenAlias(ctx context.Context, tokenFingerprint, tokenAlias string) error {
	tokenFingerprint = strings.TrimSpace(tokenFingerprint)
	tokenAlias = strings.TrimSpace(tokenAlias)

	if tokenFingerprint == "" {
		return fmt.Errorf("token fingerprint is required")
	}

	if tokenAlias == "" {
		if _, err := s.pool.Exec(ctx, `
			DELETE FROM token_aliases
			WHERE token_fingerprint = $1
		`, tokenFingerprint); err != nil {
			return fmt.Errorf("delete token alias: %w", err)
		}
		return nil
	}

	if _, err := s.pool.Exec(ctx, `
		INSERT INTO token_aliases (
			token_fingerprint,
			token_alias,
			updated_at
		)
		VALUES ($1, $2, NOW())
		ON CONFLICT (token_fingerprint)
		DO UPDATE SET
			token_alias = EXCLUDED.token_alias,
			updated_at = NOW()
	`, tokenFingerprint, tokenAlias); err != nil {
		return fmt.Errorf("upsert token alias: %w", err)
	}

	return nil
}

func buildListWhere(filters ListFilters) (string, []any) {
	clauses := []string{"l.is_capture_path = TRUE"}
	args := make([]any, 0, 8)

	if !filters.From.IsZero() {
		args = append(args, filters.From)
		clauses = append(clauses, fmt.Sprintf("l.started_at >= $%d", len(args)))
	}
	if !filters.To.IsZero() {
		args = append(args, filters.To)
		clauses = append(clauses, fmt.Sprintf("l.started_at <= $%d", len(args)))
	}
	if filters.TokenFingerprint != "" {
		args = append(args, "%"+filters.TokenFingerprint+"%")
		clauses = append(clauses, fmt.Sprintf("COALESCE(l.token_fingerprint, '') ILIKE $%d", len(args)))
	}
	if filters.TokenAlias != "" {
		args = append(args, "%"+filters.TokenAlias+"%")
		clauses = append(clauses, fmt.Sprintf("COALESCE(ta.token_alias, '') ILIKE $%d", len(args)))
	}
	if filters.Model != "" {
		clauses, args = appendModelLikeClause(clauses, args, filters.Model)
	}
	if filters.StatusCode != "" {
		args = append(args, "%"+filters.StatusCode+"%")
		clauses = append(clauses, fmt.Sprintf("CAST(COALESCE(l.status_code, 0) AS TEXT) ILIKE $%d", len(args)))
	}
	if filters.Keyword != "" {
		args = append(args, "%"+filters.Keyword+"%")
		idx := len(args)
		clauses = append(clauses, fmt.Sprintf(`
			(
				l.path ILIKE $%d OR
				COALESCE(l.user_text, '') ILIKE $%d OR
				COALESCE(l.assistant_text, '') ILIKE $%d OR
				COALESCE(l.error_text, '') ILIKE $%d OR
				COALESCE(l.model, '') ILIKE $%d OR
				COALESCE(l.token_fingerprint, '') ILIKE $%d OR
				COALESCE(ta.token_alias, '') ILIKE $%d OR
				COALESCE(l.token_preview, '') ILIKE $%d OR
				CAST(COALESCE(l.status_code, 0) AS TEXT) ILIKE $%d OR
				COALESCE(l.request_json::text, '') ILIKE $%d OR
				COALESCE(l.response_json::text, '') ILIKE $%d OR
				COALESCE(l.usage_json::text, '') ILIKE $%d
			)
		`, idx, idx, idx, idx, idx, idx, idx, idx, idx, idx, idx))
	}

	return strings.Join(clauses, " AND "), args
}

func buildTokenWhere(filters TokenDirectoryFilters) (string, []any) {
	clauses := []string{
		"l.is_capture_path = TRUE",
		"COALESCE(l.token_fingerprint, '') <> ''",
	}
	args := make([]any, 0, 6)

	if !filters.From.IsZero() {
		args = append(args, filters.From)
		clauses = append(clauses, fmt.Sprintf("l.started_at >= $%d", len(args)))
	}
	if !filters.To.IsZero() {
		args = append(args, filters.To)
		clauses = append(clauses, fmt.Sprintf("l.started_at <= $%d", len(args)))
	}
	if filters.TokenFingerprint != "" {
		args = append(args, "%"+filters.TokenFingerprint+"%")
		clauses = append(clauses, fmt.Sprintf("COALESCE(l.token_fingerprint, '') ILIKE $%d", len(args)))
	}
	if filters.TokenAlias != "" {
		args = append(args, "%"+filters.TokenAlias+"%")
		clauses = append(clauses, fmt.Sprintf("COALESCE(ta.token_alias, '') ILIKE $%d", len(args)))
	}
	if filters.Model != "" {
		clauses, args = appendModelLikeClause(clauses, args, filters.Model)
	}

	return strings.Join(clauses, " AND "), args
}

func buildDeleteWhere(filters DeleteLogsFilters) (string, []any) {
	clauses := []string{"l.is_capture_path = TRUE"}
	args := make([]any, 0, 5)

	if !filters.From.IsZero() {
		args = append(args, filters.From)
		clauses = append(clauses, fmt.Sprintf("l.started_at >= $%d", len(args)))
	}
	if !filters.To.IsZero() {
		args = append(args, filters.To)
		clauses = append(clauses, fmt.Sprintf("l.started_at <= $%d", len(args)))
	}
	if filters.TokenFingerprint != "" {
		args = append(args, "%"+filters.TokenFingerprint+"%")
		clauses = append(clauses, fmt.Sprintf("COALESCE(l.token_fingerprint, '') ILIKE $%d", len(args)))
	}
	if filters.TokenAlias != "" {
		args = append(args, "%"+filters.TokenAlias+"%")
		clauses = append(clauses, fmt.Sprintf("COALESCE(ta.token_alias, '') ILIKE $%d", len(args)))
	}
	if filters.Model != "" {
		clauses, args = appendModelLikeClause(clauses, args, filters.Model)
	}

	return strings.Join(clauses, " AND "), args
}

func hasDeleteCriteria(filters DeleteLogsFilters) bool {
	return !filters.From.IsZero() ||
		!filters.To.IsZero() ||
		filters.TokenFingerprint != "" ||
		filters.TokenAlias != "" ||
		filters.Model != ""
}

func appendModelLikeClause(clauses []string, args []any, model string) ([]string, []any) {
	args = append(args, "%"+model+"%")
	idx := len(args)
	clauses = append(clauses, fmt.Sprintf(`
		(
			COALESCE(l.model, '') ILIKE $%d OR
			COALESCE(l.user_text, '') ILIKE $%d OR
			COALESCE(l.assistant_text, '') ILIKE $%d OR
			COALESCE(l.request_json::text, '') ILIKE $%d OR
			COALESCE(l.response_json::text, '') ILIKE $%d OR
			COALESCE(encode(l.request_body, 'escape'), '') ILIKE $%d OR
			COALESCE(encode(l.response_body, 'escape'), '') ILIKE $%d
		)
	`, idx, idx, idx, idx, idx, idx, idx))
	return clauses, args
}

func buildTextPage(kind, text string, page, pageSize int, totalChars int64) TextPage {
	if page <= 0 {
		page = 1
	}
	return TextPage{
		Kind:       kind,
		Text:       text,
		Page:       page,
		PageSize:   pageSize,
		TotalChars: totalChars,
		TotalPages: totalTextPages(totalChars, pageSize),
	}
}

func buildTextPageFromFullText(kind, text string, page, pageSize int) TextPage {
	if pageSize <= 0 {
		pageSize = logTextPageChars
	}
	if page <= 0 {
		page = 1
	}

	runes := []rune(text)
	totalChars := int64(len(runes))
	totalPages := totalTextPages(totalChars, pageSize)
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	if start < 0 {
		start = 0
	}
	if start > len(runes) {
		start = len(runes)
	}

	end := start + pageSize
	if end > len(runes) {
		end = len(runes)
	}

	return buildTextPage(kind, string(runes[start:end]), page, pageSize, totalChars)
}

func totalTextPages(totalChars int64, pageSize int) int {
	if pageSize <= 0 {
		pageSize = logTextPageChars
	}
	if totalChars <= 0 {
		return 1
	}
	return int((totalChars + int64(pageSize) - 1) / int64(pageSize))
}

func textColumn(kind string) (string, error) {
	switch strings.TrimSpace(strings.ToLower(kind)) {
	case "user":
		return "l.user_text", nil
	case "assistant":
		return "l.assistant_text", nil
	default:
		return "", ErrUnsupportedTextKind
	}
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func nullBytes(value []byte) any {
	if len(value) == 0 {
		return nil
	}
	return value
}

func jsonParam(value []byte) any {
	if len(value) == 0 {
		return nil
	}
	return string(value)
}

func nullInt(value int) any {
	if value == 0 {
		return nil
	}
	return value
}

func nullTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func nullStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
