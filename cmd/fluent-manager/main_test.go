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
			originalEnv := os.Getenv("CONTAINER_LOG_PATH")
			defer func() {
				if originalEnv != "" {
					if err := os.Setenv("CONTAINER_LOG_PATH", originalEnv); err != nil {
						t.Errorf("Failed to restore CONTAINER_LOG_PATH: %v", err)
					}
				} else {
					if err := os.Unsetenv("CONTAINER_LOG_PATH"); err != nil {
						t.Errorf("Failed to unset CONTAINER_LOG_PATH: %v", err)
					}
				}
			}()

			if tt.envValue != "" {
				if err := os.Setenv("CONTAINER_LOG_PATH", tt.envValue); err != nil {
					t.Fatalf("Failed to set CONTAINER_LOG_PATH: %v", err)
				}
			} else {
				if err := os.Unsetenv("CONTAINER_LOG_PATH"); err != nil {
					t.Fatalf("Failed to unset CONTAINER_LOG_PATH: %v", err)
				}
			}

			var envFilePath string
			if tt.setupEnvFile {
				tempDir := t.TempDir()
				envFilePath = filepath.Join(tempDir, "fluent-bit.env")

				err := os.WriteFile(envFilePath, []byte(tt.envFileContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test env file: %v", err)
				}
			} else {
				envFilePath = "/non/existent/path/fluent-bit.env"
			}

			result := getContainerLogPath(envFilePath)
			if result != tt.expectedPath {
				t.Errorf("getContainerLogPath() = %q, want %q", result, tt.expectedPath)
			}
		})
	}
}
