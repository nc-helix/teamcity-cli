package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/JetBrains/teamcity-cli/api"
)

func parseHeaderFlags(headerFlags []string) (http.Header, error) {
	headers := make(http.Header, len(headerFlags))
	for _, h := range headerFlags {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format %q (expected 'Key: Value')", h)
		}

		key := http.CanonicalHeaderKey(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		headers.Set(key, value)
	}
	return headers, nil
}

func requestHeaderOption() (api.ClientOption, error) {
	headers, err := parseHeaderFlags(RequestHeaders)
	if err != nil {
		return nil, err
	}
	return api.WithDefaultHeaders(headers), nil
}
