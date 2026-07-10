package openapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/keyapi-ai/keyapi-cli/internal/config"
)

const DefaultIndexURL = "https://docs.keyapi.ai/llms.txt"

type Manifest struct {
	IndexURL string     `json:"index_url"`
	SyncedAt time.Time  `json:"synced_at"`
	Language string     `json:"language"`
	Specs    []SpecMeta `json:"specs"`
	Docs     []DocMeta  `json:"docs,omitempty"`
}

type SpecMeta struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	Platform string `json:"platform"`
	Locale   string `json:"locale,omitempty"`
}

type DocMeta struct {
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
	OpenAPIURL  string `json:"openapi_url"`
	Platform    string `json:"platform"`
	Locale      string `json:"locale"`
	Lookup      string `json:"lookup"`
}

type Registry struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type PlatformEndpoints map[string][]Endpoint

type Endpoint struct {
	Platform    string      `json:"platform"`
	SpecName    string      `json:"spec_name"`
	Method      string      `json:"method"`
	Path        string      `json:"path"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty"`
	BodyFields  []BodyField `json:"body_fields,omitempty"`
}

type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Type        string  `json:"type,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Description string  `json:"description,omitempty"`
	Schema      *schema `json:"schema,omitempty"`
}

type BodyField struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
}

type document struct {
	Paths      map[string]map[string]operation `json:"paths"`
	Components components                      `json:"components"`
}

type operation struct {
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Parameters  []rawParameter         `json:"parameters"`
	RequestBody rawRequestBody         `json:"requestBody"`
	Responses   map[string]interface{} `json:"responses"`
}

type components struct {
	Schemas map[string]schema `json:"schemas"`
}

type rawParameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Schema      schema `json:"schema"`
}

type rawRequestBody struct {
	Required bool                 `json:"required"`
	Content  map[string]mediaType `json:"content"`
}

type mediaType struct {
	Schema schema `json:"schema"`
}

type schema struct {
	Ref         string            `json:"$ref,omitempty"`
	Type        string            `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Description string            `json:"description,omitempty"`
	Enum        []interface{}     `json:"enum,omitempty"`
	Default     interface{}       `json:"default,omitempty"`
	Example     interface{}       `json:"example,omitempty"`
	Properties  map[string]schema `json:"properties,omitempty"`
	Required    []string          `json:"required,omitempty"`
	Items       *schema           `json:"items,omitempty"`
}

func LoadPlatforms() ([]string, error) {
	manifest, _, err := loadManifest(DefaultIndexURL)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var platforms []string
	for _, spec := range manifest.Specs {
		if spec.Platform == "" || seen[spec.Platform] {
			continue
		}
		seen[spec.Platform] = true
		platforms = append(platforms, spec.Platform)
	}
	sort.Strings(platforms)
	return platforms, nil
}

func LoadEndpointsForPlatforms(platforms []string) (PlatformEndpoints, error) {
	normalized, includeAll := normalizePlatformArgs(platforms)
	manifest, client, err := loadManifest(DefaultIndexURL)
	if err != nil {
		return nil, err
	}
	if !includeAll {
		manifest.Specs = filterSpecsByPlatforms(manifest.Specs, normalized)
		manifest.Docs = filterDocsByPlatforms(manifest.Docs, normalized)
		if len(manifest.Specs) == 0 {
			return nil, fmt.Errorf("no platforms found for %q", strings.Join(normalized, ", "))
		}
	}
	registry, err := loadRegistryFromManifest(manifest, client)
	if err != nil {
		return nil, err
	}
	grouped := PlatformEndpoints{}
	for _, endpoint := range registry.Endpoints {
		if includeAll || containsPlatform(normalized, endpoint.Platform) {
			grouped[endpoint.Platform] = append(grouped[endpoint.Platform], endpoint)
		}
	}
	for platform := range grouped {
		sort.Slice(grouped[platform], func(i, j int) bool {
			return grouped[platform][i].Path < grouped[platform][j].Path
		})
	}
	return grouped, nil
}

func LoadRegistryForPlatform(platform string) (Registry, error) {
	platform = strings.ToLower(strings.TrimSpace(platform))
	if platform == "" {
		return LoadRegistry()
	}
	manifest, client, err := loadManifest(DefaultIndexURL)
	if err != nil {
		return Registry{}, err
	}
	manifest.Specs = filterSpecsByPlatform(manifest.Specs, platform)
	manifest.Docs = filterDocsByPlatform(manifest.Docs, platform)
	if len(manifest.Specs) == 0 {
		return Registry{}, fmt.Errorf("platform %q not found", platform)
	}
	return loadRegistryFromManifest(manifest, client)
}

func LoadRegistry() (Registry, error) {
	manifest, client, err := loadManifest(DefaultIndexURL)
	if err != nil {
		return Registry{}, err
	}
	return loadRegistryFromManifest(manifest, client)
}

func loadManifest(indexURL string) (Manifest, *http.Client, error) {
	if indexURL == "" {
		indexURL = DefaultIndexURL
	}
	cfg, err := config.Load()
	if err != nil {
		return Manifest{}, nil, err
	}
	locale := normalizeLocale(cfg.Language)
	client := &http.Client{Timeout: 60 * time.Second}
	raw, err := getBytes(client, indexURL)
	if err != nil {
		return Manifest{}, nil, err
	}
	specs := filterSpecsByLocale(ParseOpenAPISpecLinks(string(raw)), locale)
	if len(specs) == 0 {
		return Manifest{}, nil, fmt.Errorf("no %q OpenAPI specs found in index", locale)
	}
	docs := filterDocsByLocale(parseDocLinks(string(raw)), locale)

	manifest := Manifest{
		IndexURL: indexURL,
		SyncedAt: time.Now().UTC(),
		Language: locale,
		Specs:    specs,
		Docs:     docs,
	}
	return manifest, client, nil
}

func loadRegistryFromManifest(manifest Manifest, client *http.Client) (Registry, error) {
	docByLookup := map[string]DocMeta{}
	for _, doc := range manifest.Docs {
		docByLookup[normalizeLookup(doc.Lookup)] = doc
	}

	var endpoints []Endpoint
	for _, spec := range manifest.Specs {
		if spec.Platform == "openapi" {
			continue
		}
		raw, err := getBytes(client, spec.URL)
		if err != nil {
			return Registry{}, fmt.Errorf("download %s: %w", spec.URL, err)
		}
		var doc document
		if err := json.Unmarshal(raw, &doc); err != nil {
			return Registry{}, fmt.Errorf("parse %s: %w", spec.URL, err)
		}
		for path, methods := range doc.Paths {
			for method, op := range methods {
				if !isHTTPMethod(method) {
					continue
				}
				summary := op.Summary
				description := op.Description
				if doc, ok := docByLookup[normalizeLookup(path)]; ok {
					if doc.Summary != "" {
						summary = doc.Summary
					}
					if doc.Description != "" {
						description = doc.Description
					}
				}
				endpoints = append(endpoints, Endpoint{
					Platform:    spec.Platform,
					SpecName:    spec.Name,
					Method:      strings.ToUpper(method),
					Path:        path,
					Summary:     summary,
					Description: description,
					Parameters:  normalizeParameters(op.Parameters),
					BodyFields:  bodyFields(op.RequestBody, doc.Components),
				})
			}
		}
	}

	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].Platform == endpoints[j].Platform {
			return endpoints[i].Path < endpoints[j].Path
		}
		return endpoints[i].Platform < endpoints[j].Platform
	})
	return Registry{Endpoints: endpoints}, nil
}

func (r Registry) Find(input string) (Endpoint, bool) {
	normalized := normalizeLookup(input)
	for _, endpoint := range r.Endpoints {
		candidates := []string{
			normalizeLookup(endpoint.Path),
			normalizeLookup(strings.TrimPrefix(endpoint.Path, "/v1/")),
			normalizeLookup(strings.TrimPrefix(endpoint.Path, "/")),
			normalizeLookup(endpoint.Platform + "/" + strings.TrimPrefix(strings.TrimPrefix(endpoint.Path, "/v1/"), endpoint.Platform+"/")),
		}
		for _, candidate := range candidates {
			if normalized == candidate {
				return endpoint, true
			}
		}
	}
	return Endpoint{}, false
}

func (r Registry) Platforms() []string {
	seen := map[string]bool{}
	var platforms []string
	for _, endpoint := range r.Endpoints {
		if endpoint.Platform == "" || seen[endpoint.Platform] {
			continue
		}
		seen[endpoint.Platform] = true
		platforms = append(platforms, endpoint.Platform)
	}
	sort.Strings(platforms)
	return platforms
}

func (r Registry) FilterPlatform(platform string) []Endpoint {
	var endpoints []Endpoint
	for _, endpoint := range r.Endpoints {
		if strings.EqualFold(endpoint.Platform, platform) {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

func (r Registry) Search(query string) []Endpoint {
	terms := strings.Fields(strings.ToLower(query))
	if len(terms) == 0 {
		return nil
	}
	var matches []Endpoint
	for _, endpoint := range r.Endpoints {
		haystack := strings.ToLower(strings.Join([]string{
			endpoint.Platform,
			endpoint.Method,
			endpoint.Path,
			endpoint.Summary,
			endpoint.Description,
		}, " "))
		ok := true
		for _, term := range terms {
			if !strings.Contains(haystack, term) {
				ok = false
				break
			}
		}
		if ok {
			matches = append(matches, endpoint)
		}
	}
	return matches
}

func PrintEndpointList(w io.Writer, registry Registry, platform string) error {
	platform = strings.ToLower(platform)
	for _, endpoint := range registry.Endpoints {
		if platform != "" && !strings.EqualFold(endpoint.Platform, platform) {
			continue
		}
		summary := endpoint.Summary
		if summary == "" {
			summary = endpoint.Description
		}
		if summary != "" {
			summary = " - " + oneLine(summary)
		}
		if _, err := fmt.Fprintf(w, "%-8s %-12s %s%s\n", endpoint.Method, endpoint.Platform, endpoint.Path, summary); err != nil {
			return err
		}
	}
	return nil
}

func PrintEndpointDetails(w io.Writer, endpoints []Endpoint) error {
	for i, endpoint := range endpoints {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if err := PrintEndpointDetail(w, endpoint); err != nil {
			return err
		}
	}
	return nil
}

func PrintEndpointDetail(w io.Writer, endpoint Endpoint) error {
	if _, err := fmt.Fprintf(w, "%s %s\n", endpoint.Method, endpoint.Path); err != nil {
		return err
	}
	if endpoint.Summary != "" {
		if _, err := fmt.Fprintf(w, "Summary: %s\n", endpoint.Summary); err != nil {
			return err
		}
	}
	if endpoint.Description != "" {
		if _, err := fmt.Fprintf(w, "Description: %s\n", oneLine(endpoint.Description)); err != nil {
			return err
		}
	}
	if len(endpoint.Parameters) > 0 {
		if _, err := fmt.Fprintln(w, "Parameters:"); err != nil {
			return err
		}
		for _, param := range endpoint.Parameters {
			required := "optional"
			if param.Required {
				required = "required"
			}
			typeName := param.Type
			if typeName == "" {
				typeName = "any"
			}
			if _, err := fmt.Fprintf(w, "  %-18s %-8s %-10s %s", param.Name, param.In, typeName, required); err != nil {
				return err
			}
			if param.Description != "" {
				if _, err := fmt.Fprintf(w, " - %s", oneLine(param.Description)); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
	}
	if len(endpoint.BodyFields) > 0 {
		if _, err := fmt.Fprintln(w, "Body Fields:"); err != nil {
			return err
		}
		for _, field := range endpoint.BodyFields {
			required := "optional"
			if field.Required {
				required = "required"
			}
			typeName := field.Type
			if typeName == "" {
				typeName = "any"
			}
			if _, err := fmt.Fprintf(w, "  %-18s %-10s %s", field.Name, typeName, required); err != nil {
				return err
			}
			if field.Description != "" {
				if _, err := fmt.Fprintf(w, " - %s", oneLine(field.Description)); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func ParseOpenAPISpecLinks(content string) []SpecMeta {
	start := strings.Index(content, "## OpenAPI Specs")
	if start >= 0 {
		content = content[start:]
	}

	seen := map[string]bool{}
	specs := []SpecMeta{}

	re := regexp.MustCompile(`\[([^\]]+)\]\((https://docs\.keyapi\.ai/[^\)]+\.json)\)`)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		name := match[1]
		specURL := match[2]
		if seen[specURL] {
			continue
		}
		seen[specURL] = true
		platform, locale := splitSpecURL(specURL)
		specs = append(specs, SpecMeta{
			Name:     name,
			URL:      specURL,
			Platform: platform,
			Locale:   locale,
		})
	}

	urlRe := regexp.MustCompile(`https://docs\.keyapi\.ai/[^\s\)]+\.json`)
	for _, specURL := range urlRe.FindAllString(content, -1) {
		if seen[specURL] {
			continue
		}
		seen[specURL] = true
		name := nameFromSpecURL(specURL)
		platform, locale := splitSpecURL(specURL)
		specs = append(specs, SpecMeta{
			Name:     name,
			URL:      specURL,
			Platform: platform,
			Locale:   locale,
		})
	}
	return specs
}

func parseDocLinks(content string) []DocMeta {
	start := strings.Index(content, "## Docs")
	if start < 0 {
		return nil
	}
	content = content[start:]
	if end := strings.Index(content, "## OpenAPI Specs"); end >= 0 {
		content = content[:end]
	}

	re := regexp.MustCompile(`(?m)^-\s+\[([^\]]+)\]\((https://docs\.keyapi\.ai/[^\)]+\.md)\)(?::\s*(.*))?$`)
	matches := re.FindAllStringSubmatch(content, -1)
	docs := make([]DocMeta, 0, len(matches))
	for _, match := range matches {
		docURL := match[2]
		platform, locale, lookup := splitDocURL(docURL)
		docs = append(docs, DocMeta{
			Summary:     strings.TrimSpace(match[1]),
			Description: strings.TrimSpace(match[3]),
			URL:         docURL,
			OpenAPIURL:  strings.TrimSuffix(docURL, ".md"),
			Platform:    platform,
			Locale:      locale,
			Lookup:      lookup,
		})
	}
	return docs
}

func filterSpecsByLocale(specs []SpecMeta, locale string) []SpecMeta {
	var filtered []SpecMeta
	for _, spec := range specs {
		if spec.Locale == locale && !isAggregateSpec(spec) {
			filtered = append(filtered, spec)
		}
	}
	return filtered
}

func isAggregateSpec(spec SpecMeta) bool {
	if spec.Platform == "openapi" {
		return true
	}
	parsed, err := url.Parse(spec.URL)
	if err != nil {
		return false
	}
	return strings.Contains(parsed.Path, "/api-reference/")
}

func filterDocsByLocale(docs []DocMeta, locale string) []DocMeta {
	var filtered []DocMeta
	for _, doc := range docs {
		if doc.Locale == locale {
			filtered = append(filtered, doc)
		}
	}
	return filtered
}

func filterSpecsByPlatform(specs []SpecMeta, platform string) []SpecMeta {
	var filtered []SpecMeta
	for _, spec := range specs {
		if strings.EqualFold(spec.Platform, platform) {
			filtered = append(filtered, spec)
		}
	}
	return filtered
}

func filterSpecsByPlatforms(specs []SpecMeta, platforms []string) []SpecMeta {
	var filtered []SpecMeta
	for _, spec := range specs {
		if containsPlatform(platforms, spec.Platform) {
			filtered = append(filtered, spec)
		}
	}
	return filtered
}

func filterDocsByPlatform(docs []DocMeta, platform string) []DocMeta {
	var filtered []DocMeta
	for _, doc := range docs {
		if strings.EqualFold(doc.Platform, platform) {
			filtered = append(filtered, doc)
		}
	}
	return filtered
}

func filterDocsByPlatforms(docs []DocMeta, platforms []string) []DocMeta {
	var filtered []DocMeta
	for _, doc := range docs {
		if containsPlatform(platforms, doc.Platform) {
			filtered = append(filtered, doc)
		}
	}
	return filtered
}

func normalizePlatformArgs(platforms []string) ([]string, bool) {
	if len(platforms) == 0 {
		return nil, true
	}
	seen := map[string]bool{}
	var normalized []string
	for _, platform := range platforms {
		for _, part := range strings.Split(platform, ",") {
			part = strings.ToLower(strings.TrimSpace(part))
			if part == "" {
				continue
			}
			if part == "all" {
				return nil, true
			}
			if !seen[part] {
				seen[part] = true
				normalized = append(normalized, part)
			}
		}
	}
	sort.Strings(normalized)
	return normalized, false
}

func containsPlatform(platforms []string, platform string) bool {
	for _, item := range platforms {
		if strings.EqualFold(item, platform) {
			return true
		}
	}
	return false
}

func nameFromSpecURL(specURL string) string {
	parsed, err := url.Parse(specURL)
	if err != nil {
		return "openapi"
	}
	parts := pathParts(parsed.Path)
	if len(parts) == 0 {
		return "openapi"
	}
	base := parts[len(parts)-1]
	return strings.TrimSuffix(base, ".json")
}

func splitSpecURL(rawURL string) (string, string) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", ""
	}
	parts := pathParts(parsed.Path)
	locale := ""
	platform := ""
	if len(parts) >= 2 {
		locale = normalizeLocale(parts[0])
		base := strings.TrimSuffix(parts[1], ".json")
		base = strings.TrimSuffix(base, ".openapi")
		base = strings.TrimSuffix(base, ".zh")
		platform = strings.ToLower(base)
	}
	return platform, locale
}

func splitDocURL(rawURL string) (string, string, string) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", "", ""
	}
	parts := pathParts(parsed.Path)
	if len(parts) < 2 {
		return "", "", ""
	}
	locale := normalizeLocale(parts[0])
	platform := strings.ToLower(parts[1])
	endpointParts := append([]string{}, parts[1:]...)
	last := endpointParts[len(endpointParts)-1]
	endpointParts[len(endpointParts)-1] = strings.TrimSuffix(last, ".md")
	return platform, locale, strings.Join(endpointParts, "/")
}

func pathParts(path string) []string {
	rawParts := strings.Split(strings.Trim(path, "/"), "/")
	parts := make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func normalizeLocale(locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if locale == "zh" {
		return "zh"
	}
	return "en"
}

func getBytes(client *http.Client, rawURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json,text/plain;q=0.9,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET %s failed with %s: %s", rawURL, resp.Status, strings.TrimSpace(string(body)))
	}
	return body, nil
}

func normalizeParameters(params []rawParameter) []Parameter {
	normalized := make([]Parameter, 0, len(params))
	for _, param := range params {
		normalized = append(normalized, Parameter{
			Name:        param.Name,
			In:          param.In,
			Type:        schemaType(param.Schema),
			Required:    param.Required,
			Description: param.Description,
			Schema:      schemaPointer(param.Schema),
		})
	}
	return normalized
}

func schemaPointer(value schema) *schema {
	if isEmptySchema(value) {
		return nil
	}
	return &value
}

func isEmptySchema(value schema) bool {
	return value.Ref == "" &&
		value.Type == "" &&
		value.Format == "" &&
		value.Description == "" &&
		len(value.Enum) == 0 &&
		value.Default == nil &&
		value.Example == nil &&
		len(value.Properties) == 0 &&
		len(value.Required) == 0 &&
		value.Items == nil
}

func bodyFields(body rawRequestBody, components components) []BodyField {
	var bodySchema schema
	for contentType, media := range body.Content {
		if strings.Contains(contentType, "json") {
			bodySchema = media.Schema
			break
		}
	}
	if bodySchema.Ref == "" && bodySchema.Type == "" && len(bodySchema.Properties) == 0 {
		return nil
	}
	bodySchema = resolveSchema(bodySchema, components)
	required := map[string]bool{}
	for _, name := range bodySchema.Required {
		required[name] = true
	}
	fields := make([]BodyField, 0, len(bodySchema.Properties))
	for name, property := range bodySchema.Properties {
		property = resolveSchema(property, components)
		fields = append(fields, BodyField{
			Name:        name,
			Type:        schemaType(property),
			Required:    required[name],
			Description: property.Description,
		})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

func resolveSchema(value schema, components components) schema {
	if value.Ref == "" {
		return value
	}
	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(value.Ref, prefix) {
		return value
	}
	name := strings.TrimPrefix(value.Ref, prefix)
	resolved, ok := components.Schemas[name]
	if !ok {
		return value
	}
	return resolved
}

func schemaType(value schema) string {
	if value.Type == "array" && value.Items != nil {
		itemType := schemaType(*value.Items)
		if itemType == "" {
			itemType = "any"
		}
		return "array<" + itemType + ">"
	}
	if value.Type != "" {
		return value.Type
	}
	if value.Ref != "" {
		return "object"
	}
	if len(value.Properties) > 0 {
		return "object"
	}
	return ""
}

func normalizeLookup(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	if parsed, err := url.Parse(input); err == nil && parsed.Path != "" && parsed.Scheme != "" {
		input = parsed.Path
	}
	input = strings.TrimPrefix(input, "/")
	input = strings.TrimPrefix(input, "v1/")
	return strings.Trim(input, "/")
}

func isHTTPMethod(method string) bool {
	switch strings.ToLower(method) {
	case "get", "post", "put", "patch", "delete":
		return true
	default:
		return false
	}
}

func oneLine(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if len(value) > 100 {
		return value[:97] + "..."
	}
	return value
}
