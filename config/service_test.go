package config_test

import (
	"os"
	"testing"

	"github.com/sparkymat/matrixmark/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigAPI(t *testing.T) {
	testCases := []struct {
		testCaseName          string
		env                   map[string]string
		expectedShioriURL     string
		expectedShioriUserMap map[string]config.Credentials
		panicExpected         bool
	}{
		{
			testCaseName:  "panics if required env missing",
			panicExpected: true,
		},
		{
			testCaseName: "panics if user map is incorrect",
			env: map[string]string{
				"SHIORI_URL":   "http://example.com",
				"SHIORI_USERS": "foo",
			},
			panicExpected: true,
		},
		{
			testCaseName: "no panic if user map is correct",
			env: map[string]string{
				"SHIORI_URL":   "http://example.com",
				"SHIORI_USERS": "foo:bar:secret",
			},
			expectedShioriURL: "http://example.com",
			expectedShioriUserMap: map[string]config.Credentials{
				"foo": {
					Username: "bar",
					Password: "secret",
				},
			},
			panicExpected: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			for envKey, envValue := range testCase.env {
				os.Setenv(envKey, envValue)
			}
			if testCase.panicExpected {
				require.Panics(t, func() {
					config.New()
				})
			} else {
				require.NotPanics(t, func() {
					c := config.New()
					assert.Equal(t, testCase.expectedShioriURL, c.ShioriURL())
					assert.Equal(t, testCase.expectedShioriUserMap, c.ShioriUsersMap())
				})
			}
			for envKey, _ := range testCase.env {
				os.Unsetenv(envKey)
			}
		})
	}
}
