package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/cobra"

	"bizzmod-cli/internal/api"
	"bizzmod-cli/internal/config"
)

func testClientWithRole(t *testing.T, role string) *api.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/external/profile" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"role":"` + role + `"}}`))
	}))
	t.Cleanup(server.Close)
	return api.NewClient(config.Config{APIURL: server.URL, UserAPIKey: "test-key"})
}

func TestEnsureRoleAllowsAdmin(t *testing.T) {
	cmd := &cobra.Command{}
	err := ensureRole(cmd, testClientWithRole(t, "admin"), "license status", "admin", "superadmin")
	if err != nil {
		t.Fatalf("expected admin to be allowed, got %v", err)
	}
}

func TestEnsureRoleAllowsSuperadmin(t *testing.T) {
	cmd := &cobra.Command{}
	err := ensureRole(cmd, testClientWithRole(t, "superadmin"), "license management", "superadmin")
	if err != nil {
		t.Fatalf("expected superadmin to be allowed, got %v", err)
	}
}

func TestEnsureRoleDeniesAdminForSuperadminScope(t *testing.T) {
	cmd := &cobra.Command{}
	err := ensureRole(cmd, testClientWithRole(t, "admin"), "license management", "superadmin")
	if err == nil {
		t.Fatal("expected permission denied for admin on a superadmin-only scope")
	}
}

func TestEnsureRoleDeniesRegularUser(t *testing.T) {
	cmd := &cobra.Command{}
	err := ensureRole(cmd, testClientWithRole(t, "user"), "license status", "admin", "superadmin")
	if err == nil {
		t.Fatal("expected permission denied for a regular user")
	}
}
