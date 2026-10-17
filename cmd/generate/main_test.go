package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func parseSpec(t *testing.T, source string) OpenAPISpec {
	t.Helper()
	var spec OpenAPISpec
	if err := json.Unmarshal([]byte(source), &spec); err != nil {
		t.Fatalf("decode test spec: %v", err)
	}
	return spec
}

func TestGenerateSourcePreservesSchemaSemantics(t *testing.T) {
	spec := parseSpec(t, `{
		"paths": {
			"/sample.getRecords": {
				"post": {
					"operationId": "sample.getRecords",
					"requestBody": {"content": {"application/json": {"schema": {
						"required": ["requiredId"],
						"properties": {
							"requiredId": {"type": "string"},
							"enabled": {"type": "boolean"},
							"kind": {"anyOf": [{"type": "string"}, {"type": "integer"}]},
							"tags": {"type": "array", "items": {"type": "string"}}
						}
					}}}}
				}
			}
		}
	}`)

	generated, err := generateSource(spec)
	if err != nil {
		t.Fatalf("generate source: %v", err)
	}
	source := string(generated)
	wantFields := map[string][2]string{
		"RequiredID": {"string", "`json:\"requiredId\"`"},
		"Enabled":    {"*bool", "`json:\"enabled,omitempty\"`"},
		"Kind":       {"any", "`json:\"kind,omitempty\"`"},
		"Tags":       {"[]string", "`json:\"tags,omitempty\"`"},
	}
	for field, want := range wantFields {
		var fieldLine string
		for _, line := range strings.Split(source, "\n") {
			if strings.Contains(line, field+" ") {
				fieldLine = line
				break
			}
		}
		if !strings.Contains(fieldLine, want[0]) || !strings.Contains(fieldLine, want[1]) {
			t.Errorf("generated field %s is incorrect: %q", field, fieldLine)
		}
	}

	again, err := generateSource(spec)
	if err != nil {
		t.Fatalf("generate source a second time: %v", err)
	}
	if string(again) != source {
		t.Fatal("generation was not deterministic")
	}
}

func TestGenerateSourceRejectsUnsupportedSchema(t *testing.T) {
	spec := parseSpec(t, `{
		"paths": {"/sample.getRecords": {"post": {
			"operationId": "sample.getRecords",
			"requestBody": {"content": {"application/json": {"schema": {
				"properties": {"cursor": {"type": "object"}}
			}}}}
		}}}
	}`)

	if _, err := generateSource(spec); err == nil || !strings.Contains(err.Error(), `unsupported OpenAPI type "object"`) {
		t.Fatalf("expected unsupported object schema error, got %v", err)
	}
}

func TestGenerateSourceRejectsFieldNameCollision(t *testing.T) {
	spec := parseSpec(t, `{
		"paths": {"/sample.getRecords": {"post": {
			"operationId": "sample.getRecords",
			"requestBody": {"content": {"application/json": {"schema": {
				"properties": {"userId": {"type": "string"}, "userID": {"type": "string"}}
			}}}}
		}}}
	}`)

	if _, err := generateSource(spec); err == nil || !strings.Contains(err.Error(), "both map to Go field UserID") {
		t.Fatalf("expected field name collision error, got %v", err)
	}
}
