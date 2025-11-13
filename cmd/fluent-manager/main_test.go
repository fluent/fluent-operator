/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetContainerLogPath(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		setupEnvFile   bool
		envFileContent string
		expectedPath   string
	}{
		{
			name:         "returns path from CONTAINER_LOG_PATH environment variable",
			envValue:     "/custom/log/path",
			setupEnvFile: false,
			expectedPath: "/custom/log/path",
		},
		{
			name:           "returns path from fluent-bit.env file when env var not set",
			envValue:       "",
			setupEnvFile:   true,
			envFileContent: "CONTAINER_ROOT_DIR=/var/lib/docker",
			expectedPath:   "/var/lib/docker/containers",
		},
		{
			name:           "returns path with different CONTAINER_ROOT_DIR",
			envValue:       "",
			setupEnvFile:   true,
			envFileContent: "CONTAINER_ROOT_DIR=/custom/container/root",
			expectedPath:   "/custom/container/root/containers",
		},
		{
			name:         "returns default path when neither env var nor file is available",
			envValue:     "",
			setupEnvFile: false,
			expectedPath: "/var/log/containers",
		},
		{
			name:           "prefers env var over file when both are set",
			envValue:       "/priority/path",
			setupEnvFile:   true,
			envFileContent: "CONTAINER_ROOT_DIR=/var/lib/docker",
			expectedPath:   "/priority/path",
		},
		{
			name:           "handles env file with multiple variables",
			envValue:       "",
			setupEnvFile:   true,
			envFileContent: "OTHER_VAR=value\nCONTAINER_ROOT_DIR=/multi/var/test\nANOTHER_VAR=other",
			expectedPath:   "/multi/var/test/containers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment variable
			originalEnv := os.Getenv("CONTAINER_LOG_PATH")
			defer func() {
				if originalEnv != "" {
					os.Setenv("CONTAINER_LOG_PATH", originalEnv)
				} else {
					os.Unsetenv("CONTAINER_LOG_PATH")
				}
			}()

			// Set up environment variable for this test
			if tt.envValue != "" {
				os.Setenv("CONTAINER_LOG_PATH", tt.envValue)
			} else {
				os.Unsetenv("CONTAINER_LOG_PATH")
			}

			// Set up temporary fluent-bit.env file if needed
			var envFilePath string
			if tt.setupEnvFile {
				// Create a temporary directory for testing
				tempDir := t.TempDir()
				envFilePath = filepath.Join(tempDir, "fluent-bit.env")

				// Write the test env file
				err := os.WriteFile(envFilePath, []byte(tt.envFileContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test env file: %v", err)
				}
			} else {
				// Use a non-existent file path
				envFilePath = "/non/existent/path/fluent-bit.env"
			}

			// Call the function with the test file path
			result := getContainerLogPath(envFilePath)

			// Verify the result
			if result != tt.expectedPath {
				t.Errorf("getContainerLogPath() = %q, want %q", result, tt.expectedPath)
			}
		})
	}
}
