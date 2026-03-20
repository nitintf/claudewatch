package statusline

import (
	"strings"
	"testing"

	"github.com/nitintf/claudewatch/internal/api"
	"github.com/nitintf/claudewatch/internal/theme"
)

func ptrFloat(f float64) *float64 { return &f }

func testTheme() *theme.Theme {
	return &theme.Theme{
		Name: "test",
		Colors: theme.Colors{
			Bg: "#282a36", Fg: "#f8f8f2", Accent: "#bd93f9",
			Success: "#50fa7b", Warning: "#f1fa8c", Error: "#ff5555",
			Info: "#8be9fd", Muted: "#6272a4",
		},
	}
}

func TestParse(t *testing.T) {
	input := `{"session_id":"abc","model":{"display_name":"Opus","id":"claude-opus-4-6"},"context_window":{"used_percentage":42},"cost":{"total_cost_usd":1.23}}`

	s, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if s.SessionID != "abc" {
		t.Errorf("got session %q", s.SessionID)
	}
	if s.Model.DisplayName != "Opus" {
		t.Errorf("got model %q", s.Model.DisplayName)
	}
	if s.ContextWindow.UsedPercentage == nil || *s.ContextWindow.UsedPercentage != 42 {
		t.Errorf("got context %v", s.ContextWindow.UsedPercentage)
	}
	if s.Cost.TotalCostUSD != 1.23 {
		t.Errorf("got cost %v", s.Cost.TotalCostUSD)
	}
}

func TestParseInvalid(t *testing.T) {
	_, err := Parse([]byte("nope"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRenderWithUsage(t *testing.T) {
	s := ClaudeStatus{
		Model:         modelInfo{DisplayName: "Opus"},
		ContextWindow: contextWindow{UsedPercentage: ptrFloat(22)},
	}
	usage := &api.Usage{
		FiveHour: &api.QuotaLimit{Utilization: 23},
		SevenDay: &api.QuotaLimit{Utilization: 2},
	}
	out := Render(s, testTheme(), "Pro", usage)
	if !strings.Contains(out, "Opus") {
		t.Error("expected Opus in output")
	}
	if !strings.Contains(out, "Pro") {
		t.Error("expected Pro in output")
	}
	if !strings.Contains(out, "22%") {
		t.Error("expected 22% in output")
	}
	if !strings.Contains(out, "23%") {
		t.Error("expected 23% for 5h in output")
	}
}

func TestRenderWithoutUsage(t *testing.T) {
	s := ClaudeStatus{
		Model:         modelInfo{DisplayName: "Opus"},
		ContextWindow: contextWindow{UsedPercentage: ptrFloat(50)},
	}
	out := Render(s, testTheme(), "", nil)
	if !strings.Contains(out, "Opus") {
		t.Error("expected Opus in output")
	}
	if !strings.Contains(out, "50%") {
		t.Error("expected 50% in output")
	}
}

func TestRenderContainsANSI(t *testing.T) {
	s := ClaudeStatus{
		Model:         modelInfo{DisplayName: "Opus"},
		ContextWindow: contextWindow{UsedPercentage: ptrFloat(50)},
	}
	out := Render(s, testTheme(), "Pro", nil)
	if !strings.Contains(out, "\x1b[") {
		t.Error("expected ANSI escape codes in output")
	}
}

func TestProgressBar(t *testing.T) {
	b := progressBar(50, 3, "#8be9fd", "#6272a4")
	if !strings.Contains(b, "\u2588") {
		t.Error("expected filled blocks")
	}
	if !strings.Contains(b, "\u2591") {
		t.Error("expected empty blocks")
	}
}

func TestBarColor(t *testing.T) {
	th := testTheme()
	if barColor(30, th) != th.Colors.Info {
		t.Error("30% should be info")
	}
	if barColor(55, th) != th.Colors.Warning {
		t.Error("55% should be warning")
	}
	if barColor(90, th) != th.Colors.Error {
		t.Error("90% should be error")
	}
}

func TestHexToRGB(t *testing.T) {
	r, g, b := hexToRGB("#ff8800")
	if r != 255 || g != 136 || b != 0 {
		t.Errorf("got %d,%d,%d", r, g, b)
	}
}
