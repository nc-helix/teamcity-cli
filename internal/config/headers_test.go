package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(T *testing.T) {
	headers, err := ParseHeaders([]string{"X-Test: value", "X-Other:  second"})
	require.NoError(T, err)
	assert.Equal(T, "value", headers["X-Test"])
	assert.Equal(T, "second", headers["X-Other"])
}

func TestParseHeadersInvalid(T *testing.T) {
	_, err := ParseHeaders([]string{"X-Test"})
	require.Error(T, err)
	assert.Contains(T, err.Error(), "invalid header format")
}

func TestGetExtraHeadersFromEnv(T *testing.T) {
	saveCfgState(T)
	cfg = &Config{
		DefaultServer: "https://tc.example.com",
		Servers: map[string]ServerConfig{
			"https://tc.example.com": {
				Headers: map[string]string{"X-From-Config": "nope"},
			},
		},
	}

	T.Setenv(EnvHeaders, "X-From-Env: yes;X-Other: value")

	headers, err := GetExtraHeaders()
	require.NoError(T, err)
	assert.Equal(T, map[string]string{"X-From-Env": "yes", "X-Other": "value"}, headers)
}

func TestGetExtraHeadersFromConfig(T *testing.T) {
	saveCfgState(T)
	cfg = &Config{
		DefaultServer: "https://tc.example.com",
		Servers: map[string]ServerConfig{
			"https://tc.example.com": {
				Headers: map[string]string{"X-IAP": "token"},
			},
		},
	}
	T.Setenv(EnvHeaders, "")

	headers, err := GetExtraHeaders()
	require.NoError(T, err)
	assert.Equal(T, map[string]string{"X-IAP": "token"}, headers)
}

func TestGetExtraHeadersInvalidEnv(T *testing.T) {
	saveCfgState(T)
	T.Setenv(EnvHeaders, "bad-header")

	_, err := GetExtraHeaders()
	require.Error(T, err)
	assert.Contains(T, err.Error(), EnvHeaders)
}
