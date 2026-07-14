package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveEndpointDoesNotNeedOpenAPI(t *testing.T) {
	path, method := resolveEndpoint("google/search")
	if path != "/v1/google/search" {
		t.Fatalf("expected /v1/google/search, got %q", path)
	}
	if method != "" {
		t.Fatalf("expected empty inferred method, got %q", method)
	}
}

func TestResolveEndpointPreservesInlineQuery(t *testing.T) {
	path, method := resolveEndpoint("/v1/google/search?deal_id=a2c250a5&country=us")
	if path != "/v1/google/search?deal_id=a2c250a5&country=us" {
		t.Fatalf("expected inline query to be preserved, got %q", path)
	}
	if method != "" {
		t.Fatalf("expected empty inferred method, got %q", method)
	}
}

func TestPlatformFromEndpoint(t *testing.T) {
	tests := map[string]string{
		"google/search":                         "google",
		"/v1/google/search":                     "google",
		"https://api.keyapi.ai/v1/tiktok/video": "tiktok",
	}
	for input, want := range tests {
		if got := platformFromEndpoint(input); got != want {
			t.Fatalf("platformFromEndpoint(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestSplitKV(t *testing.T) {
	k, v, err := splitKV("q=openai")
	if err != nil {
		t.Fatal(err)
	}
	if k != "q" || v != "openai" {
		t.Fatalf("expected q/openai, got %q/%q", k, v)
	}
}

func TestNormalizeBodyJSON(t *testing.T) {
	tests := map[string]string{
		`{"page":1,"verified":true}`:          `{"page":1,"verified":true}`,
		`'{page:1,verified:true}'`:            `{"page":1,"verified":true}`,
		`{page:1,verified:true}`:              `{"page":1,"verified":true}`,
		`{q:openai,page:1,verified:true}`:     `{"q":"openai","page":1,"verified":true}`,
		`{filter:{countryCode:US},limit:100}`: `{"filter":{"countryCode":"US"},"limit":100}`,
	}

	for input, want := range tests {
		got, err := normalizeBodyJSON(input)
		if err != nil {
			t.Fatalf("normalizeBodyJSON(%q) returned error: %v", input, err)
		}
		if got != want {
			t.Fatalf("normalizeBodyJSON(%q) = %s, want %s", input, got, want)
		}
	}
}

func TestCallRejectsInvalidBodyJSON(t *testing.T) {
	t.Setenv("KEYAPI_API_KEY", "sk_test")

	err := runCall([]string{"google/search", "--body-json", `{"keyword":`})
	if err == nil {
		t.Fatal("expected invalid --body-json error")
	}
}

func TestCallBodyJSONAcceptsPowerShellMangledObject(t *testing.T) {
	server := newCallTestServer(t, func(r *http.Request, body string) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("expected application/json content type, got %q", got)
		}
		if body != `{"page":1,"verified":true}` {
			t.Fatalf("unexpected body %q", body)
		}
	})
	defer server.Close()

	writeCallTestConfig(t, server.URL)

	err := runCall([]string{"google/search", "--body-json", `{page:1,verified:true}`})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCallParamsAliasAppendsQuery(t *testing.T) {
	server := newCallTestServer(t, func(r *http.Request, body string) {
		if got := r.URL.Query().Get("gl"); got != "us" {
			t.Fatalf("expected gl=us, got %q", got)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Fatalf("expected page=1, got %q", got)
		}
	})
	defer server.Close()

	writeCallTestConfig(t, server.URL)

	err := runCall([]string{"/v1/google/search?gl=us", "--params", "page=1"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCallRejectsBodyRawAndBodyJSONTogether(t *testing.T) {
	t.Setenv("KEYAPI_API_KEY", "sk_test")

	err := runCall([]string{"google/search", "--body-raw", "plain text", "--body-json", `{"keyword":"openai"}`})
	if err == nil {
		t.Fatal("expected --body-raw and --body-json conflict error")
	}
}

func TestRunTopLevelConfigSetAndGet(t *testing.T) {
	home := t.TempDir()
	t.Setenv("KEYAPI_HOME", home)
	t.Setenv("KEYAPI_API_KEY", "")

	if err := Run([]string{"set", "token", "sk_test"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"get", "token"}); err != nil {
		t.Fatal(err)
	}
}

func TestCallBodyJSONUsesJSONContentType(t *testing.T) {
	server := newCallTestServer(t, func(r *http.Request, body string) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("expected application/json content type, got %q", got)
		}
		if body != `{"timeRange":7,"countryCode":"US","page":1,"limit":100}` {
			t.Fatalf("unexpected body %q", body)
		}
	})
	defer server.Close()

	writeCallTestConfig(t, server.URL)

	err := runCall([]string{"google/search", "--body-json", `{"timeRange":7,"countryCode":"US","page":1,"limit":100}`})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCallBodyRawAlwaysUsesTextContentType(t *testing.T) {
	server := newCallTestServer(t, func(r *http.Request, body string) {
		if got := r.Header.Get("Content-Type"); got != "text/plain; charset=utf-8" {
			t.Fatalf("expected text/plain content type, got %q", got)
		}
		if body != `{"q":"openai"}` {
			t.Fatalf("unexpected body %q", body)
		}
	})
	defer server.Close()

	writeCallTestConfig(t, server.URL)

	err := runCall([]string{"google/search", "--body-raw", `{"q":"openai"}`})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCallOutputWritesRelativePathFromWorkingDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"message":"success","data":{"items":[1,2]}}`))
	}))
	defer server.Close()

	writeCallTestConfig(t, server.URL)
	wd := t.TempDir()
	t.Chdir(wd)

	err := runCall([]string{"google/search", "--output", filepath.Join("result", "google.json")})
	if err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(filepath.Join(wd, "result", "google.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(raw) {
		t.Fatalf("expected valid JSON output, got %q", raw)
	}
	if !strings.Contains(string(raw), `"items"`) {
		t.Fatalf("expected response body in output file, got %q", raw)
	}
}

func TestCallOutputWritesNon2xxResponseBeforeReturningError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":401,"message":"apiKey is null","data":null}`))
	}))
	defer server.Close()

	writeCallTestConfig(t, server.URL)
	wd := t.TempDir()
	t.Chdir(wd)

	err := runCall([]string{"google/search", "--output", "error.json"})
	if err == nil {
		t.Fatal("expected non-2xx error")
	}
	raw, readErr := os.ReadFile(filepath.Join(wd, "error.json"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(raw), `"apiKey is null"`) {
		t.Fatalf("expected error response body in output file, got %q", raw)
	}
}

func TestCallReturnsErrorForNon2xxHTTPStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":401,"message":"apiKey is null","data":null}`))
	}))
	defer server.Close()

	writeCallTestConfig(t, server.URL)

	err := runCall([]string{"google/search"})
	if err == nil {
		t.Fatal("expected non-2xx HTTP status error")
	}
	statusErr, ok := err.(httpStatusError)
	if !ok {
		t.Fatalf("expected httpStatusError, got %T: %v", err, err)
	}
	if statusErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, statusErr.StatusCode)
	}
}

func newCallTestServer(t *testing.T, assert func(r *http.Request, body string)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert(r, string(raw))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"message":"success","data":{}}`))
	}))
}

func writeCallTestConfig(t *testing.T, baseURL string) {
	t.Helper()

	home := t.TempDir()
	t.Setenv("KEYAPI_HOME", home)
	t.Setenv("KEYAPI_API_KEY", "")
	raw := []byte(`{"api_key":"sk_test","base_url":` + strconvQuote(baseURL) + `,"auth_scheme":"Bearer","language":"en"}`)
	if err := os.WriteFile(filepath.Join(home, "config.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
}

func strconvQuote(value string) string {
	raw, _ := json.Marshal(value)
	return string(raw)
}
