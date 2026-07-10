package config

import "testing"

func TestLoadAppliesEnvironmentOverrides(t *testing.T) {
	t.Setenv("KEYAPI_HOME", t.TempDir())
	t.Setenv("KEYAPI_API_KEY", "env-key")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "env-key" {
		t.Fatalf("expected env API key, got %q", cfg.APIKey)
	}
	if cfg.BaseURL != defaultBaseURL {
		t.Fatalf("expected default base URL, got %q", cfg.BaseURL)
	}
	if cfg.Language != "en" {
		t.Fatalf("expected default language en, got %q", cfg.Language)
	}
}

func TestSaveAndLoad(t *testing.T) {
	t.Setenv("KEYAPI_HOME", t.TempDir())
	t.Setenv("KEYAPI_API_KEY", "")

	want := Config{
		APIKey:     "stored-key",
		BaseURL:    "https://api.example.com",
		AuthScheme: "Bearer",
		Language:   "zh",
	}
	if err := Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
}
