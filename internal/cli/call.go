package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/EchoSell/keyapi-cli/internal/config"
)

func runCall(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		printCallHelp()
		return nil
	}

	endpoint := args[0]
	fs := flag.NewFlagSet("call", flag.ContinueOnError)
	method := fs.String("method", "", "HTTP method; inferred from body flags when omitted")
	apiKey := fs.String("api-key", "", "API key")
	token := fs.String("token", "", "API token; alias for --api-key")
	bodyJSON := fs.String("body-json", "", "raw JSON request body")
	bodyRaw := fs.String("body-raw", "", "raw text request body")
	format := fs.String("format", "json", "output format: json, pretty, or raw")
	dryRun := fs.Bool("dry-run", false, "print the request without sending it")
	var queries repeatedFlag
	fs.Var(&queries, "query", "query parameter as \"k=v\"; may be repeated")
	fs.Var(&queries, "qparam", "alias for --query")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if *bodyJSON != "" && *bodyRaw != "" {
		return errors.New("use either --body-raw or --body-json, not both")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if *apiKey != "" {
		cfg.APIKey = *apiKey
	}
	if *token != "" {
		cfg.APIKey = *token
	}
	if cfg.APIKey == "" {
		return errors.New("missing API key; set KEYAPI_API_KEY or run keyapi config set token <value>")
	}

	path, inferredMethod := resolveEndpoint(endpoint)
	if *method == "" {
		if *bodyJSON != "" || *bodyRaw != "" {
			*method = http.MethodPost
		} else if inferredMethod != "" {
			*method = inferredMethod
		} else {
			*method = http.MethodGet
		}
	}

	target, err := url.Parse(cfg.BaseURL + path)
	if err != nil {
		return err
	}
	values := target.Query()
	for _, item := range queries {
		k, v, err := splitKV(item)
		if err != nil {
			return fmt.Errorf("invalid --query %q: %w", item, err)
		}
		values.Add(k, v)
	}
	target.RawQuery = values.Encode()

	var body io.Reader
	var bodyValue interface{}
	contentType := ""
	if *bodyJSON != "" {
		if err := json.Unmarshal([]byte(*bodyJSON), &bodyValue); err != nil {
			return fmt.Errorf("invalid --body-json: %w", err)
		}
		body = strings.NewReader(*bodyJSON)
		contentType = "application/json"
	} else if *bodyRaw != "" {
		body = strings.NewReader(*bodyRaw)
		bodyValue = *bodyRaw
		contentType = "text/plain; charset=utf-8"
	}

	req, err := http.NewRequest(strings.ToUpper(*method), target.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	if strings.EqualFold(cfg.AuthScheme, "bearer") || cfg.AuthScheme == "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	} else {
		req.Header.Set("Authorization", cfg.AuthScheme+" "+cfg.APIKey)
	}

	if *dryRun {
		return printJSON(dryRunRequest{
			Method: req.Method,
			URL:    req.URL.String(),
			Headers: map[string]string{
				"Accept":        req.Header.Get("Accept"),
				"Content-Type":  req.Header.Get("Content-Type"),
				"Authorization": maskedAuthorization(req.Header.Get("Authorization")),
			},
			Body: bodyValue,
		})
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := printResponse(raw, *format); err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return httpStatusError{StatusCode: resp.StatusCode, Status: resp.Status}
	}
	return nil
}

func printResponse(raw []byte, format string) error {
	if format == "raw" {
		fmt.Println(string(raw))
		return nil
	}
	if format != "json" && format != "pretty" {
		return fmt.Errorf("unsupported format %q; use pretty, raw, or json", format)
	}

	var pretty bytes.Buffer
	if json.Valid(raw) && json.Indent(&pretty, raw, "", "  ") == nil {
		fmt.Println(pretty.String())
		return nil
	}
	fmt.Println(string(raw))
	return nil
}

type httpStatusError struct {
	StatusCode int
	Status     string
}

func (e httpStatusError) Error() string {
	return fmt.Sprintf("request failed with HTTP status %s", e.Status)
}
