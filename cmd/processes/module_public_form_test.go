package processes

import (
	"encoding/json"
	"testing"
)

func TestParsePublicFormSubmissionPayloadWrapsPayload(t *testing.T) {
	body, err := parsePublicFormSubmissionPayload(`{"email":"laura@example.com"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("cannot decode body: %v", err)
	}

	payload, ok := decoded["payload"].(map[string]any)
	if !ok {
		t.Fatalf("expected payload object, got %T", decoded["payload"])
	}
	if payload["email"] != "laura@example.com" {
		t.Fatalf("unexpected email value: %v", payload["email"])
	}
}

func TestParsePublicFormSubmissionPayloadRejectsInvalidJSON(t *testing.T) {
	_, err := parsePublicFormSubmissionPayload(`{"email":`)
	if err == nil {
		t.Fatal("expected error for invalid payload_json")
	}
}

func TestParsePublicFormSubmissionPayloadRejectsEmptyInput(t *testing.T) {
	_, err := parsePublicFormSubmissionPayload("   ")
	if err == nil {
		t.Fatal("expected error for empty payload_json")
	}
}

func TestBuildPublicFormSubmissionPath(t *testing.T) {
	path := buildPublicFormSubmissionPath("form-generator", "demo-sl", "contact-form")
	expected := "/api/v1/external/public/modules/form-generator/forms/demo-sl/contact-form/submit"
	if path != expected {
		t.Fatalf("unexpected path\nwant: %s\ngot:  %s", expected, path)
	}
}
