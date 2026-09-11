package processes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"bizzmod-cli/internal/config"
)

func TestHasConfigChanged(t *testing.T) {
	base := config.Config{
		APIURL:     "http://localhost:5000",
		UserAPIKey: "k",
	}

	if hasConfigChanged(base, base) {
		t.Fatalf("expected unchanged config")
	}

	changed := base
	changed.UserAPIKey = "other"
	if !hasConfigChanged(base, changed) {
		t.Fatalf("expected changed config")
	}
}

func TestPromptReviewFieldsKeepsDefaultsOnEnter(t *testing.T) {
	in := bytes.NewBufferString("\n\n")
	out := bytes.NewBuffer(nil)
	cmd := newLoginCmd()
	cmd.SetIn(in)
	cmd.SetOut(out)

	initial := config.Config{
		APIURL:     "http://localhost:5000",
		UserAPIKey: "abc",
	}

	got, err := promptReviewFields(initial, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != initial {
		t.Fatalf("expected defaults kept, got %+v", got)
	}
}

func TestPromptReviewFieldsUpdatesUserAPIKey(t *testing.T) {
	in := bytes.NewBufferString("\nnew-key\n")
	out := bytes.NewBuffer(nil)
	cmd := newLoginCmd()
	cmd.SetIn(in)
	cmd.SetOut(out)

	initial := config.Config{
		APIURL:     "http://localhost:5000",
		UserAPIKey: "abc",
	}

	got, err := promptReviewFields(initial, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.UserAPIKey != "new-key" {
		t.Fatalf("expected updated user API key")
	}
}

func TestValidateExternalCredentialsFailsOnBadResponse(t *testing.T) {
	cfg := config.Config{
		APIURL:     "http://127.0.0.1:1",
		UserAPIKey: "k",
	}

	err := validateExternalCredentials(cfg)
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestValidateExternalCredentialsWithStubServer(t *testing.T) {
	server := httpTestServer(http.StatusOK, `{"success":true,"data":{"rows":[]}}`)
	defer server.Close()

	cfg := config.Config{
		APIURL:     server.URL,
		UserAPIKey: "k",
	}

	if err := validateExternalCredentials(cfg); err != nil {
		t.Fatalf("expected successful validation, got %v", err)
	}
}

func httpTestServer(status int, body string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/external/modules" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"detail":"not found"}`))
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})
	return httptest.NewServer(handler)
}
