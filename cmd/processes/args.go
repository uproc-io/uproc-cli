package processes

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func parsePositiveIntArg(name, raw string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return value, nil
}

func parseJSONObjectArg(name, raw string) (map[string]any, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, fmt.Errorf("%s is required", name)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", name, err)
	}
	if payload == nil {
		return nil, fmt.Errorf("%s must be a JSON object", name)
	}
	return payload, nil
}
