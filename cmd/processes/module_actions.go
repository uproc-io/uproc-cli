package processes

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func runModuleAction(cmd *cobra.Command, moduleSlug, action string, payload map[string]any) error {
	client, err := mustClient()
	if err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot encode %s payload: %w", action, err)
	}

	path := fmt.Sprintf("/api/v1/external/modules/%s/actions/%s", moduleSlug, action)
	respBody, status, reqErr := client.Do("POST", path, body)
	return printResponse(cmd, respBody, status, reqErr)
}
