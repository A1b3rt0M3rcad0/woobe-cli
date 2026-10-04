package controlplane

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	Base         string
	Token        string
	HTTP         *http.Client
	Headers      http.Header
	ResponseHook func(*http.Response) error
}

func New(base, token string, timeout time.Duration) (*Client, error) {
	u, e := url.Parse(base)
	if e != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, output.New(2, "api URL must be an HTTP(S) origin without credentials, query or fragment")
	}
	h := &http.Client{Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return &Client{Base: strings.TrimRight(base, "/"), Token: token, HTTP: h, Headers: make(http.Header)}, nil
}
func (c *Client) Request(ctx context.Context, method, path string, q url.Values, body []byte) (any, http.Header, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	return c.RequestReader(ctx, method, path, q, reader, "application/json")
}
func (c *Client) RequestReader(ctx context.Context, method, path string, q url.Values, body io.Reader, contentType string) (any, http.Header, error) {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.ContainsAny(path, "?#\\") {
		return nil, nil, output.New(2, "request path must be an origin-relative path; use --query for query parameters")
	}
	raw := c.Base + path
	if len(q) > 0 {
		raw += "?" + q.Encode()
	}
	req, e := http.NewRequestWithContext(ctx, method, raw, body)
	if e != nil {
		return nil, nil, output.New(2, "invalid request")
	}
	req.Header = c.Headers.Clone()
	req.Header.Set("Accept", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	resp, e := c.HTTP.Do(req)
	if e != nil {
		err := output.Normalize(e)
		if err.Code == 1 {
			err.Code = 7
			err.Message = "HTTP connection failed"
		}
		if method != "GET" && method != "HEAD" {
			err.Outcome = "unknown"
		}
		return nil, nil, err
	}
	defer resp.Body.Close()
	if c.ResponseHook != nil {
		if e = c.ResponseHook(resp); e != nil {
			return nil, resp.Header, &output.Error{Code: 10, Message: "server responded but session persistence failed", Outcome: "unknown"}
		}
	}
	b, e := io.ReadAll(io.LimitReader(resp.Body, 32<<20+1))
	if e != nil {
		return nil, resp.Header, responseError(method, resp, 7, "response could not be read")
	}
	if len(b) > 32<<20 {
		return nil, resp.Header, responseError(method, resp, 9, "response exceeds limit")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		code := 7
		switch resp.StatusCode {
		case 400, 422:
			code = 2
		case 401:
			code = 3
		case 403:
			code = 4
		case 404:
			code = 5
		case 409, 412:
			code = 6
		case 408, 504:
			code = 8
		case 405, 501:
			code = 9
		}
		outcome := "rejected"
		if resp.StatusCode >= 500 {
			outcome = "unknown"
		}
		return nil, resp.Header, &output.Error{Code: code, Message: "server rejected request", Status: resp.StatusCode, RequestID: resp.Header.Get("X-Request-ID"), Outcome: outcome}
	}
	if len(bytes.TrimSpace(b)) == 0 {
		return nil, resp.Header, nil
	}
	var v any
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	if e = d.Decode(&v); e != nil {
		return nil, resp.Header, responseError(method, resp, 9, "server returned non-JSON response")
	}
	if env, ok := v.(map[string]any); ok {
		if success, ok := env["success"].(bool); ok && !success {
			return nil, resp.Header, &output.Error{Code: 7, Message: "server returned unsuccessful response", Outcome: "unknown"}
		}
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return nil, resp.Header, responseError(method, resp, 9, "server returned multiple JSON values")
	}
	return v, resp.Header, nil
}

func responseError(method string, resp *http.Response, code int, message string) *output.Error {
	outcome := ""
	if method != "GET" && method != "HEAD" {
		outcome = "unknown"
	}
	return &output.Error{Code: code, Message: message, Status: resp.StatusCode, RequestID: resp.Header.Get("X-Request-ID"), Outcome: outcome}
}
