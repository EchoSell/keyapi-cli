package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	defaultBaseURL    = "https://api.keyapi.ai"
	defaultAuthScheme = "Bearer"
)

type Config struct {
	APIKey     string `json:"api_key"`
	BaseURL    string `json:"base_url"`
	AuthScheme string `json:"auth_scheme"`
	Language   string `json:"language"`
}

func Load() (Config, error) {
	cfg := Config{
		BaseURL:    defaultBaseURL,
		AuthScheme: defaultAuthScheme,
		Language:   "en",
	}

	path, err := FilePath()
	if err != nil {
		return cfg, err
	}
	raw, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return cfg, err
		}
	} else if !os.IsNotExist(err) {
		return cfg, err
	}

	if value := os.Getenv("KEYAPI_API_KEY"); value != "" {
		cfg.APIKey = value
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.AuthScheme == "" {
		cfg.AuthScheme = defaultAuthScheme
	}
	if cfg.Language == "" {
		cfg.Language = "en"
	}
	return cfg, nil
}

func Save(cfg Config) error {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.AuthScheme == "" {
		cfg.AuthScheme = defaultAuthScheme
	}
	if cfg.Language == "" {
		cfg.Language = "en"
	}
	path, err := FilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o600)
}

func Dir() (string, error) {
	if value := os.Getenv("KEYAPI_HOME"); value != "" {
		return value, nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "keyapi"), nil
}

func FilePath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}
