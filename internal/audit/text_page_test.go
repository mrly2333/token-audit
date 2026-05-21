package audit

import (
	"strings"
	"testing"
)

func TestBuildTextPageFromFullText_ClampsAndSlicesRunes(t *testing.T) {
	text := strings.Repeat("你", 3999) + "A" + strings.Repeat("界", 10)

	page1 := buildTextPageFromFullText("user", text, 1, 4000)
	if page1.Page != 1 {
		t.Fatalf("expected page 1, got %d", page1.Page)
	}
	if page1.TotalPages != 2 {
		t.Fatalf("expected 2 total pages, got %d", page1.TotalPages)
	}
	if page1.TotalChars != 4010 {
		t.Fatalf("expected 4010 total chars, got %d", page1.TotalChars)
	}
	if got := len([]rune(page1.Text)); got != 4000 {
		t.Fatalf("expected 4000 chars on page 1, got %d", got)
	}
	if !strings.HasSuffix(page1.Text, "A") {
		t.Fatalf("expected page 1 to end with A, got %q", string([]rune(page1.Text)[len([]rune(page1.Text))-1]))
	}

	page2 := buildTextPageFromFullText("user", text, 2, 4000)
	if page2.Page != 2 {
		t.Fatalf("expected page 2, got %d", page2.Page)
	}
	if got := len([]rune(page2.Text)); got != 10 {
		t.Fatalf("expected 10 chars on page 2, got %d", got)
	}
	if page2.Text != strings.Repeat("界", 10) {
		t.Fatalf("unexpected page 2 content: %q", page2.Text)
	}
}

func TestBuildTextPageFromFullText_OutOfRangeFallsBackToLastPage(t *testing.T) {
	text := strings.Repeat("x", 4500)

	page := buildTextPageFromFullText("assistant", text, 99, 4000)
	if page.Page != 2 {
		t.Fatalf("expected clamped page 2, got %d", page.Page)
	}
	if page.TotalPages != 2 {
		t.Fatalf("expected 2 total pages, got %d", page.TotalPages)
	}
	if got := len([]rune(page.Text)); got != 500 {
		t.Fatalf("expected 500 chars on last page, got %d", got)
	}
}
