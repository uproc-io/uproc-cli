package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bizzmod-cli/internal/config"
)

func TestClientUsesUserAPIKeyBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer user-key" {
			t.Fatalf("expected bearer user API key, got %q", got)
		}
		if r.Header.Get("x-api-key") != "" || r.Header.Get("x-customer-domain") != "" || r.Header.Get("x-user-email") != "" {
			t.Fatalf("customer authentication headers must not be sent")
		}
		_, _ = io.WriteString(w, `{"success":true}`)
	}))
	defer server.Close()

	client := NewClient(config.Config{APIURL: server.URL, UserAPIKey: "user-key"})
	if _, _, err := client.Do(http.MethodGet, "/api/v1/external/modules", nil); err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}
}
