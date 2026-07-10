package openapi

import "testing"

func TestParseSpecLinksSupportsMarkdownAndBareURLs(t *testing.T) {
	content := `## OpenAPI Specs
- [google.openapi](https://docs.keyapi.ai/en/google.openapi.json)
- https://docs.keyapi.ai/en/amazon.openapi.json
- [tiktok.zh.openapi](https://docs.keyapi.ai/zh/tiktok.zh.openapi.json)
- [openapi](https://docs.keyapi.ai/api-reference/openapi.json)
- [google.openapi](https://docs.keyapi.ai/en/google.openapi.json)
`

	specs := ParseOpenAPISpecLinks(content)
	if len(specs) != 4 {
		t.Fatalf("expected 4 specs, got %d", len(specs))
	}
	if specs[0].Platform != "google" || specs[0].Locale != "en" {
		t.Fatalf("unexpected first spec: %#v", specs[0])
	}
	foundZH := false
	for _, spec := range specs {
		if spec.Platform == "tiktok" && spec.Locale == "zh" {
			foundZH = true
		}
	}
	if !foundZH {
		t.Fatalf("expected zh tiktok spec in %#v", specs)
	}

	enSpecs := filterSpecsByLocale(specs, "en")
	if len(enSpecs) != 2 {
		t.Fatalf("expected 2 en specs, got %d", len(enSpecs))
	}
	for _, spec := range enSpecs {
		if spec.Platform == "openapi" || spec.URL == "https://docs.keyapi.ai/api-reference/openapi.json" {
			t.Fatalf("aggregate spec should be filtered: %#v", spec)
		}
	}
}

func TestParseDocLinks(t *testing.T) {
	content := `## Docs
- [Search](https://docs.keyapi.ai/en/google/search.md): Google search endpoint.
- [搜索](https://docs.keyapi.ai/zh/google/search.md): Google 搜索接口。

## OpenAPI Specs
- [google.openapi](https://docs.keyapi.ai/en/google.openapi.json)
`

	docs := parseDocLinks(content)
	if len(docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(docs))
	}
	if docs[0].Summary != "Search" || docs[0].Description != "Google search endpoint." {
		t.Fatalf("unexpected doc metadata: %#v", docs[0])
	}
	if docs[0].Lookup != "google/search" || docs[0].Locale != "en" || docs[0].Platform != "google" {
		t.Fatalf("unexpected doc lookup metadata: %#v", docs[0])
	}
}

func TestRegistryFindSupportsPathAliases(t *testing.T) {
	registry := Registry{
		Endpoints: []Endpoint{
			{
				Platform: "google",
				Method:   "GET",
				Path:     "/v1/google/search",
			},
		},
	}

	for _, input := range []string{"google/search", "/google/search", "/v1/google/search"} {
		endpoint, ok := registry.Find(input)
		if !ok {
			t.Fatalf("expected to find endpoint for %q", input)
		}
		if endpoint.Path != "/v1/google/search" {
			t.Fatalf("expected /v1/google/search, got %q", endpoint.Path)
		}
	}
}

func TestRegistryPlatformsAndSearch(t *testing.T) {
	registry := Registry{
		Endpoints: []Endpoint{
			{Platform: "google", Method: "GET", Path: "/v1/google/search", Summary: "search"},
			{Platform: "amazon", Method: "GET", Path: "/v1/amazon/product-details", Summary: "product details"},
		},
	}

	platforms := registry.Platforms()
	if len(platforms) != 2 || platforms[0] != "amazon" || platforms[1] != "google" {
		t.Fatalf("unexpected platforms: %#v", platforms)
	}

	matches := registry.Search("product details")
	if len(matches) != 1 || matches[0].Platform != "amazon" {
		t.Fatalf("unexpected matches: %#v", matches)
	}
}

func TestBodyFieldsExtractsJSONSchemaProperties(t *testing.T) {
	fields := bodyFields(rawRequestBody{
		Content: map[string]mediaType{
			"application/json": {
				Schema: schema{
					Type:     "object",
					Required: []string{"keyword"},
					Properties: map[string]schema{
						"keyword": {Type: "string", Description: "search keyword"},
						"offset":  {Type: "integer"},
					},
				},
			},
		},
	}, components{})

	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}
	if fields[0].Name != "keyword" || !fields[0].Required || fields[0].Type != "string" {
		t.Fatalf("unexpected first field: %#v", fields[0])
	}
}

func TestNormalizeParametersPreservesSchemaMetadata(t *testing.T) {
	params := normalizeParameters([]rawParameter{
		{
			Name:     "region",
			In:       "query",
			Required: true,
			Schema: schema{
				Type:    "string",
				Enum:    []interface{}{"US", "GB"},
				Default: "US",
				Example: "GB",
			},
		},
	})

	if len(params) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(params))
	}
	param := params[0]
	if param.Schema == nil {
		t.Fatal("expected schema to be preserved")
	}
	if param.Schema.Default != "US" || param.Schema.Example != "GB" {
		t.Fatalf("unexpected schema defaults/examples: %#v", param.Schema)
	}
	if len(param.Schema.Enum) != 2 || param.Schema.Enum[0] != "US" || param.Schema.Enum[1] != "GB" {
		t.Fatalf("unexpected schema enum: %#v", param.Schema.Enum)
	}
}
