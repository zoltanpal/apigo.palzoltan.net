package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	url   string
	token string
	http  *http.Client
}

func New(url, token string, timeout time.Duration) *Client {
	return &Client{
		url:   url, // e.g. https://api.palzoltan.net/sentiment_analyzer/analyze_text
		token: token,
		http:  &http.Client{Timeout: timeout},
	}
}

type analyzeReq struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

// The API returns a flat JSON map (keys -> float values)
type AnalyzeResp map[string]float64

func (c *Client) Analyze(ctx context.Context, lang, text string) (AnalyzeResp, error) {
	body, _ := json.Marshal(analyzeReq{Lang: lang, Text: text})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.http.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// surface auth errors clearly
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("sentiment auth failed: HTTP %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("sentiment HTTP %d", resp.StatusCode)
	}

	var out AnalyzeResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
