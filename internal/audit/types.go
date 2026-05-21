package audit

import "time"

type Record struct {
	ID                  int64
	StartedAt           time.Time
	FinishedAt          time.Time
	DurationMS          int64
	Method              string
	Path                string
	QueryString         string
	RemoteAddr          string
	RequestHost         string
	UpstreamBase        string
	StatusCode          int
	ErrorText           string
	IsCapturePath       bool
	ResponseIsSSE       bool
	TokenFingerprint    string
	TokenPreview        string
	TokenAlias          string
	Model               string
	Stream              *bool
	RequestHeaders      map[string][]string
	ResponseHeaders     map[string][]string
	RequestContentType  string
	ResponseContentType string
	RequestBody         []byte
	ResponseBody        []byte
	RequestJSON         []byte
	ResponseJSON        []byte
	UserText            string
	AssistantText       string
	UsageJSON           []byte
	PromptTokens        int64
	CompletionTokens    int64
	TotalTokens         int64
	RequestBytes        int64
	ResponseBytes       int64
	RequestTruncated    bool
	ResponseTruncated   bool
}

type DashboardStats struct {
	TotalRequests          int64
	TodayRequests          int64
	ErrorCount             int64
	DistinctTokens         int64
	TodayDistinctTokens    int64
	TotalPromptTokens      int64
	TotalCompletionTokens  int64
	TotalTokens            int64
	TodayPromptTokens      int64
	TodayCompletionTokens  int64
	TodayTotalTokens       int64
	TokenGroups            []TokenStats
	ModelGroups            []ModelStats
}

type TextPage struct {
	Kind       string
	Text       string
	Page       int
	PageSize   int
	TotalChars int64
	TotalPages int
}

type LogCore struct {
	Record         Record
	UserText       TextPage
	AssistantText  TextPage
}

type TokenStats struct {
	TokenFingerprint string
	TokenPreview     string
	TokenAlias       string
	RequestCount     int64
	ErrorCount       int64
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
	LastSeen         time.Time
}

type ModelStats struct {
	Model            string
	RequestCount     int64
	ErrorCount       int64
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
	LastSeen         time.Time
}

type ListFilters struct {
	From             time.Time
	To               time.Time
	TokenFingerprint string
	TokenAlias       string
	Model            string
	StatusCode       string
	Keyword          string
	Page             int
	PageSize         int
}

type LogListItem struct {
	ID               int64
	StartedAt        time.Time
	Method           string
	Path             string
	StatusCode       int
	Model            string
	TokenFingerprint string
	TokenPreview     string
	TokenAlias       string
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
	UserText         string
	AssistantText    string
	ErrorText        string
}

type ListResult struct {
	Items      []LogListItem
	TotalCount int64
	Page       int
	PageSize   int
}

type TokenDirectoryFilters struct {
	From             time.Time
	To               time.Time
	TokenFingerprint string
	TokenAlias       string
	Model            string
	Page             int
	PageSize         int
}

type TokenDirectoryItem struct {
	TokenFingerprint string
	TokenPreview     string
	TokenAlias       string
	RequestCount     int64
	ErrorCount       int64
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
	FirstSeen        time.Time
	LastSeen         time.Time
}

type TokenDirectoryResult struct {
	Items      []TokenDirectoryItem
	TotalCount int64
	Page       int
	PageSize   int
}

type FilterOptions struct {
	Models            []string
	TokenAliases      []string
	TokenFingerprints []string
	StatusCodes       []string
}

type DatabaseStats struct {
	TotalRows         int64
	TodayRows         int64
	DatabaseSize      int64
	DatabasePretty    string
	AuditTotalSize    int64
	AuditTotalPretty  string
	AuditTableSize    int64
	AuditTablePretty  string
	AuditIndexSize    int64
	AuditIndexPretty  string
	AuditToastSize    int64
	AuditToastPretty  string
	LiveTuples        int64
	DeadTuples        int64
	LastVacuum        time.Time
	LastAutovacuum    time.Time
	LastAnalyze       time.Time
	LastAutoanalyze   time.Time
}

type DeleteLogsFilters struct {
	From             time.Time
	To               time.Time
	TokenFingerprint string
	TokenAlias       string
	Model            string
}
