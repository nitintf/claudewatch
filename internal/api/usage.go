package api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const (
	usageURL          = "https://api.anthropic.com/api/oauth/usage"
	cacheTTLOK        = 60 * time.Second
	cacheTTLFail      = 15 * time.Second
	cacheTTLRateLimit = 5 * time.Minute
	maxBackoff        = 30 * time.Minute
	requestTimeout    = 5 * time.Second
)

// QuotaLimit is a single usage quota with utilization and reset time.
type QuotaLimit struct {
	Utilization float64 `json:"utilization"`
	ResetsAt    string  `json:"resets_at"`
}

// ExtraUsage is pay-as-you-go overage info.
type ExtraUsage struct {
	IsEnabled    bool     `json:"is_enabled"`
	MonthlyLimit *float64 `json:"monthly_limit"`
	UsedCredits  *float64 `json:"used_credits"`
}

// Usage is the response from the usage API.
type Usage struct {
	FiveHour         *QuotaLimit `json:"five_hour"`
	SevenDay         *QuotaLimit `json:"seven_day"`
	SevenDaySonnet   *QuotaLimit `json:"seven_day_sonnet"`
	SevenDayOpus     *QuotaLimit `json:"seven_day_opus"`
	SevenDayOAuthApp *QuotaLimit `json:"seven_day_oauth_apps"`
	SevenDayCowork   *QuotaLimit `json:"seven_day_cowork"`
	ExtraUsage       *ExtraUsage `json:"extra_usage"`
}

type cacheEntry struct {
	Data        json.RawMessage `json:"data"`
	Timestamp   int64           `json:"ts"`
	OK          bool            `json:"ok"`
	RateLimited bool            `json:"rl,omitempty"`
	RetryAfter  int64           `json:"ra,omitempty"`
}

// FetchUsage gets usage data from the API with caching.
func FetchUsage(token string) (*Usage, error) {
	if cached, err := readCache(); err == nil {
		return cached, nil
	}

	usage, retryAfter, err := fetchAPI(token)
	if err != nil {
		writeCache(nil, false, retryAfter)
		return nil, err
	}

	writeCache(usage, true, 0)
	return usage, nil
}

func fetchAPI(token string) (*Usage, time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, usageURL, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Anthropic-Beta", "oauth-2025-04-20")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests {
		ra := parseRetryAfter(resp.Header.Get("Retry-After"))
		return nil, ra, fmt.Errorf("rate limited (429)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var usage Usage
	if err := json.Unmarshal(body, &usage); err != nil {
		return nil, 0, err
	}
	return &usage, 0, nil
}

func parseRetryAfter(val string) time.Duration {
	if val == "" {
		return cacheTTLRateLimit
	}
	if secs, err := strconv.Atoi(val); err == nil && secs > 0 {
		d := time.Duration(secs) * time.Second
		if d > maxBackoff {
			return maxBackoff
		}
		return d
	}
	return cacheTTLRateLimit
}

func cacheFilePath() string {
	dir := "/tmp"
	if runtime.GOOS == "windows" {
		dir = os.TempDir()
	}
	suffix := ""
	if cd := os.Getenv("CLAUDE_CONFIG_DIR"); cd != "" {
		h := sha256.Sum256([]byte(cd))
		suffix = fmt.Sprintf("-%x", h[:4])
	}
	return filepath.Join(dir, "claudewatch-usage"+suffix+".json")
}

func readCache() (*Usage, error) {
	data, err := os.ReadFile(cacheFilePath())
	if err != nil {
		return nil, err
	}
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}

	age := time.Since(time.Unix(entry.Timestamp, 0))

	if entry.OK && age < cacheTTLOK {
		var usage Usage
		if err := json.Unmarshal(entry.Data, &usage); err != nil {
			return nil, err
		}
		return &usage, nil
	}
	if !entry.OK && entry.RateLimited {
		if entry.RetryAfter > 0 && time.Now().Unix() < entry.RetryAfter {
			return nil, fmt.Errorf("rate limited (cached)")
		}
		if entry.RetryAfter == 0 && age < cacheTTLRateLimit {
			return nil, fmt.Errorf("rate limited (cached)")
		}
	}
	if !entry.OK && age < cacheTTLFail {
		return nil, fmt.Errorf("cached failure")
	}
	return nil, fmt.Errorf("cache expired")
}

func writeCache(usage *Usage, ok bool, retryAfter time.Duration) {
	entry := cacheEntry{
		Timestamp:   time.Now().Unix(),
		OK:          ok,
		RateLimited: retryAfter > 0,
	}
	if retryAfter > 0 {
		entry.RetryAfter = time.Now().Add(retryAfter).Unix()
	}
	if usage != nil {
		d, err := json.Marshal(usage)
		if err != nil {
			return
		}
		entry.Data = d
	}
	d, err := json.Marshal(entry)
	if err != nil {
		return
	}
	_ = os.WriteFile(cacheFilePath(), d, 0o600)
}
