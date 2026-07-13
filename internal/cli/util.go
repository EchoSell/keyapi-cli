package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
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

func normalizeBodyJSON(input string) (string, error) {
	input = strings.TrimSpace(input)
	if len(input) >= 2 && input[0] == '\'' && input[len(input)-1] == '\'' {
		input = strings.TrimSpace(input[1 : len(input)-1])
	}
	if json.Valid([]byte(input)) {
		return input, nil
	}

	parser := looseJSONParser{input: input}
	repaired, err := parser.parseValue()
	if err != nil {
		return "", err
	}
	parser.skipSpaces()
	if parser.pos != len(parser.input) {
		return "", fmt.Errorf("unexpected trailing content at byte %d", parser.pos)
	}
	if !json.Valid([]byte(repaired)) {
		return "", errors.New("value is not valid JSON")
	}
	return repaired, nil
}

type looseJSONParser struct {
	input string
	pos   int
}

func (p *looseJSONParser) parseValue() (string, error) {
	p.skipSpaces()
	if p.pos >= len(p.input) {
		return "", errors.New("missing value")
	}

	switch p.input[p.pos] {
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '"':
		return p.parseDoubleQuotedString()
	case '\'':
		value, err := p.parseSingleQuotedString()
		if err != nil {
			return "", err
		}
		return quoteJSONString(value), nil
	default:
		token := p.readBareToken()
		if token == "" {
			return "", fmt.Errorf("unexpected character %q at byte %d", p.input[p.pos], p.pos)
		}
		if isJSONLiteral(token) || isJSONNumber(token) {
			return token, nil
		}
		return quoteJSONString(token), nil
	}
}

func (p *looseJSONParser) parseObject() (string, error) {
	p.pos++
	var b strings.Builder
	b.WriteByte('{')
	first := true

	for {
		p.skipSpaces()
		if p.pos >= len(p.input) {
			return "", errors.New("unterminated object")
		}
		if p.input[p.pos] == '}' {
			p.pos++
			b.WriteByte('}')
			return b.String(), nil
		}
		if !first {
			if p.input[p.pos] != ',' {
				return "", fmt.Errorf("expected comma at byte %d", p.pos)
			}
			p.pos++
			p.skipSpaces()
		}

		key, err := p.parseObjectKey()
		if err != nil {
			return "", err
		}
		p.skipSpaces()
		if p.pos >= len(p.input) || p.input[p.pos] != ':' {
			return "", fmt.Errorf("expected colon after object key at byte %d", p.pos)
		}
		p.pos++
		value, err := p.parseValue()
		if err != nil {
			return "", err
		}
		if !first {
			b.WriteByte(',')
		}
		b.WriteString(quoteJSONString(key))
		b.WriteByte(':')
		b.WriteString(value)
		first = false
	}
}

func (p *looseJSONParser) parseArray() (string, error) {
	p.pos++
	var b strings.Builder
	b.WriteByte('[')
	first := true

	for {
		p.skipSpaces()
		if p.pos >= len(p.input) {
			return "", errors.New("unterminated array")
		}
		if p.input[p.pos] == ']' {
			p.pos++
			b.WriteByte(']')
			return b.String(), nil
		}
		if !first {
			if p.input[p.pos] != ',' {
				return "", fmt.Errorf("expected comma at byte %d", p.pos)
			}
			p.pos++
		}
		value, err := p.parseValue()
		if err != nil {
			return "", err
		}
		if !first {
			b.WriteByte(',')
		}
		b.WriteString(value)
		first = false
	}
}

func (p *looseJSONParser) parseObjectKey() (string, error) {
	p.skipSpaces()
	if p.pos >= len(p.input) {
		return "", errors.New("missing object key")
	}
	switch p.input[p.pos] {
	case '"':
		raw, err := p.parseDoubleQuotedString()
		if err != nil {
			return "", err
		}
		var key string
		if err := json.Unmarshal([]byte(raw), &key); err != nil {
			return "", err
		}
		return key, nil
	case '\'':
		return p.parseSingleQuotedString()
	default:
		start := p.pos
		for p.pos < len(p.input) && p.input[p.pos] != ':' {
			p.pos++
		}
		key := strings.TrimSpace(p.input[start:p.pos])
		if key == "" {
			return "", errors.New("empty object key")
		}
		return key, nil
	}
}

func (p *looseJSONParser) parseDoubleQuotedString() (string, error) {
	start := p.pos
	p.pos++
	escaped := false
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		p.pos++
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == '"' {
			raw := p.input[start:p.pos]
			if !json.Valid([]byte(raw)) {
				return "", fmt.Errorf("invalid string at byte %d", start)
			}
			return raw, nil
		}
	}
	return "", errors.New("unterminated string")
}

func (p *looseJSONParser) parseSingleQuotedString() (string, error) {
	p.pos++
	start := p.pos
	for p.pos < len(p.input) {
		if p.input[p.pos] == '\'' {
			value := p.input[start:p.pos]
			p.pos++
			return value, nil
		}
		p.pos++
	}
	return "", errors.New("unterminated single-quoted string")
}

func (p *looseJSONParser) readBareToken() string {
	start := p.pos
	for p.pos < len(p.input) {
		switch p.input[p.pos] {
		case ',', '}', ']':
			return strings.TrimSpace(p.input[start:p.pos])
		default:
			p.pos++
		}
	}
	return strings.TrimSpace(p.input[start:p.pos])
}

func (p *looseJSONParser) skipSpaces() {
	for p.pos < len(p.input) {
		switch p.input[p.pos] {
		case ' ', '\t', '\r', '\n':
			p.pos++
		default:
			return
		}
	}
}

func isJSONLiteral(token string) bool {
	return token == "true" || token == "false" || token == "null"
}

func isJSONNumber(token string) bool {
	if token == "" {
		return false
	}
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func quoteJSONString(value string) string {
	raw, _ := json.Marshal(value)
	return string(raw)
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
