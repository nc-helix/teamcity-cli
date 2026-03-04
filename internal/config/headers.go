package config

import (
	"fmt"
	"os"
	"strings"
)

// GetExtraHeaders returns additional HTTP headers that should be sent with every request.
// Priority: TEAMCITY_HEADERS env var > per-server config headers.
func GetExtraHeaders() (map[string]string, error) {
	if raw := os.Getenv(EnvHeaders); raw != "" {
		headers, err := ParseHeaders(splitHeaderEnv(raw))
		if err != nil {
			return nil, fmt.Errorf("invalid %s: %w", EnvHeaders, err)
		}
		return headers, nil
	}

	serverURL := GetServerURL()
	if serverURL == "" || cfg == nil {
		return nil, nil
	}

	server, ok := cfg.Servers[serverURL]
	if !ok || len(server.Headers) == 0 {
		return nil, nil
	}

	headers := make(map[string]string, len(server.Headers))
	for k, v := range server.Headers {
		headers[k] = v
	}

	return headers, nil
}

// ParseHeaders parses headers in "Key: Value" format.
func ParseHeaders(entries []string) (map[string]string, error) {
	headers := make(map[string]string)
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format %q (expected 'Key: Value')", entry)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("invalid header format %q (header name is empty)", entry)
		}
		headers[key] = value
	}

	return headers, nil
}

func splitHeaderEnv(raw string) []string {
	return strings.FieldsFunc(raw, func(r rune) bool {
		return r == ';' || r == '\n'
	})
}
