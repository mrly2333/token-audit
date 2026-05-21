package web

import (
	"bytes"
	"io"
	"log"
	"strings"
	"testing"

	"newapi-audit-proxy/internal/config"
)

func TestBuildPageURL(t *testing.T) {
	url := buildPageURL("/audit", logFiltersView{
		From:       "2026-05-18T10:00",
		To:         "2026-05-18T11:00",
		TokenAlias: "demo",
		Model:      "gpt-5.5",
		StatusCode: "200",
		Keyword:    "hello",
	}, 2)

	expectedParts := []string{
		"/audit/logs?",
		"from=2026-05-18T10%3A00",
		"to=2026-05-18T11%3A00",
		"alias=demo",
		"model=gpt-5.5",
		"status=200",
		"q=hello",
		"page=2",
	}
	for _, part := range expectedParts {
		if !strings.Contains(url, part) {
			t.Fatalf("expected %q to contain %q", url, part)
		}
	}
}

func TestLogsTemplateRendersNativePaginationLinks(t *testing.T) {
	server, err := New(config.Config{
		WebBasePath: "/audit",
		HMACSecret:  "test-secret",
	}, nil, log.New(io.Discard, "", 0))
	if err != nil {
		t.Fatalf("create server: %v", err)
	}

	var buf bytes.Buffer
	data := logsViewData{
		Title: "logs",
		Filters: logFiltersView{
			Model: "gpt-5.5",
		},
		Page:       1,
		PageSize:   100,
		TotalCount: 250,
		TotalPages: 3,
		HasPrev:    false,
		HasNext:    true,
		PrevPage:   1,
		NextPage:   2,
	}
	renderData, err := server.withCSRFToken(data)
	if err != nil {
		t.Fatalf("csrf token: %v", err)
	}
	if err := server.templates["logs"].ExecuteTemplate(&buf, "logs", renderData); err != nil {
		t.Fatalf("render logs template: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `id="next-page" href="/audit/logs?model=gpt-5.5&amp;page=2"`) {
		t.Fatalf("next page link not rendered as expected: %s", html)
	}
	if !strings.Contains(html, `id="prev-page" href="/audit/logs?model=gpt-5.5&amp;page=1"`) {
		t.Fatalf("prev page link not rendered as expected: %s", html)
	}
	if strings.Contains(html, `href="#"`) {
		t.Fatalf("unexpected placeholder href found in rendered html")
	}
	if strings.Contains(html, `window.location.assign(refs.nextPage.href)`) {
		t.Fatalf("unexpected JS pagination interception still present")
	}
}
