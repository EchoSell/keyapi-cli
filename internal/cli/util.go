package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func printJSON(v any) error {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(raw))
	return nil
}

func splitKV(input string) (string, string, error) {
	key, value, ok := strings.Cut(input, "=")
	if !ok || key == "" {
		return "", "", errors.New("expected k=v")
	}
	return key, value, nil
}

func mask(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "********"
	}
	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}

func isHelp(value string) bool {
	return value == "-h" || value == "--help" || value == "help"
}

func parseTargetAndFormat(args []string, usage string) (string, string, error) {
	return parseTargetAndFormatDefault(args, usage, "pretty")
}

func parseTargetAndFormatDefault(args []string, usage string, defaultFormat string) (string, string, error) {
	target := ""
	format := defaultFormat
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format":
			if i+1 >= len(args) {
				return "", "", errors.New("missing value for --format")
			}
			format = args[i+1]
			i++
		default:
			if strings.HasPrefix(args[i], "--format=") {
				format = strings.TrimPrefix(args[i], "--format=")
				continue
			}
			if strings.HasPrefix(args[i], "-") {
				return "", "", fmt.Errorf("unknown flag %q", args[i])
			}
			if target != "" {
				return "", "", fmt.Errorf("usage: %s", usage)
			}
			target = args[i]
		}
	}
	if target == "" {
		return "", "", fmt.Errorf("usage: %s", usage)
	}
	return target, format, nil
}

func resolveEndpoint(input string) (string, string) {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		u, err := url.Parse(input)
		if err == nil {
			return u.RequestURI(), ""
		}
	}
	path := input
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(path, "/v") {
		path = "/v1" + path
	}
	return path, ""
}

func platformFromEndpoint(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	if parsed, err := url.Parse(input); err == nil && parsed.Path != "" && parsed.Scheme != "" {
		input = parsed.Path
	}
	input = strings.Trim(input, "/")
	input = strings.TrimPrefix(input, "v1/")
	parts := strings.Split(input, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func maskedAuthorization(value string) string {
	scheme, token, ok := strings.Cut(value, " ")
	if !ok {
		return mask(value)
	}
	return scheme + " " + mask(token)
}

type repeatedFlag []string

func (r *repeatedFlag) String() string {
	return strings.Join(*r, ",")
}

func (r *repeatedFlag) Set(value string) error {
	*r = append(*r, value)
	return nil
}
