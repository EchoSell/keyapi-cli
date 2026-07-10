package cli

import "testing"

func TestResolveEndpointDoesNotNeedOpenAPI(t *testing.T) {
	path, method := resolveEndpoint("google/search")
	if path != "/v1/google/search" {
		t.Fatalf("expected /v1/google/search, got %q", path)
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
