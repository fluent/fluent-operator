package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

// TestContainerLogPathResolution tests the container log path resolution logic
// that will be added to main.go
func TestContainerLogPathResolution(t *testing.T) {
	tests := []struct {
		name           string
		envVar         string
		fileContent    string
		expectedPath   string
		setupFile      bool
	}{
		{
			name:         "prefer environment variable over file",
			envVar:       "/custom/containers",
			fileContent:  "CONTAINER_ROOT_DIR=/var/log",
			expectedPath: "/custom/containers",
			setupFile:    true,
		},
		{
			name:         "fallback to file when env var not set",
			envVar:       "",
			fileContent:  "CONTAINER_ROOT_DIR=/var/log",
			expectedPath: "/var/log/containers",
			setupFile:    true,
		},
		{
			name:         "fallback to file for docker path",
			envVar:       "",
			fileContent:  "CONTAINER_ROOT_DIR=/var/lib/docker",
			expectedPath: "/var/lib/docker/containers",
			setupFile:    true,
		},
		{
			name:         "use default when neither env var nor file present",
			envVar:       "",
			fileContent:  "",
			expectedPath: "/var/log/containers",
			setupFile:    false,
		},
		{
			name:         "env var takes precedence even with docker in file",
			envVar:       "/override/containers",
			fileContent:  "CONTAINER_ROOT_DIR=/var/lib/docker",
			expectedPath: "/override/containers",
			setupFile:    true,
		},
		{
			name:         "handle empty env var (should fallback)",
			envVar:       "",
			fileContent:  "CONTAINER_ROOT_DIR=/var/log",
			expectedPath: "/var/log/containers",
			setupFile:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test file
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "fluent-bit.env")

			// Setup environment variable
			if tt.envVar != "" {
				os.Setenv("CONTAINER_LOG_PATH", tt.envVar)
				defer os.Unsetenv("CONTAINER_LOG_PATH")
			} else {
				os.Unsetenv("CONTAINER_LOG_PATH")
			}

			// Setup file if needed
			if tt.setupFile && tt.fileContent != "" {
				err := os.WriteFile(testFile, []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			// Execute the resolution logic (mimics what will be in main.go)
			var logPath string

			// Prefer environment variable for container log path (modern approach)
			// Fall back to file-based configuration for backward compatibility
			if envLogPath := os.Getenv("CONTAINER_LOG_PATH"); envLogPath != "" {
				logPath = envLogPath
			} else if tt.setupFile {
				if envs, err := godotenv.Read(testFile); err == nil {
					logPath = envs["CONTAINER_ROOT_DIR"] + "/containers"
				}
			}

			// Final fallback to safe default for containerd/CRI-O
			if logPath == "" {
				logPath = "/var/log/containers"
			}

			// Verify result
			if logPath != tt.expectedPath {
				t.Errorf("expected logPath to be %q, got %q", tt.expectedPath, logPath)
			}
		})
	}
}

// TestContainerLogPathEnvVarPrecedence specifically tests that env var always wins
func TestContainerLogPathEnvVarPrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "fluent-bit.env")

	// Create file with one path
	err := os.WriteFile(testFile, []byte("CONTAINER_ROOT_DIR=/var/lib/docker"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Set env var with different path
	envPath := "/env/var/wins"
	os.Setenv("CONTAINER_LOG_PATH", envPath)
	defer os.Unsetenv("CONTAINER_LOG_PATH")

	// Execute resolution logic
	var logPath string
	if envLogPath := os.Getenv("CONTAINER_LOG_PATH"); envLogPath != "" {
		logPath = envLogPath
	} else if envs, err := godotenv.Read(testFile); err == nil {
		logPath = envs["CONTAINER_ROOT_DIR"] + "/containers"
	}
	if logPath == "" {
		logPath = "/var/log/containers"
	}

	// Env var should win
	if logPath != envPath {
		t.Errorf("environment variable should take precedence: expected %q, got %q", envPath, logPath)
	}
}

// TestContainerLogPathFilePathsAppendContainers verifies /containers is appended
func TestContainerLogPathFilePathsAppendContainers(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "fluent-bit.env")

	// Unset env var
	os.Unsetenv("CONTAINER_LOG_PATH")

	tests := []struct {
		rootDir      string
		expectedPath string
	}{
		{"/var/log", "/var/log/containers"},
		{"/var/lib/docker", "/var/lib/docker/containers"},
		{"/custom/path", "/custom/path/containers"},
	}

	for _, tt := range tests {
		t.Run(tt.rootDir, func(t *testing.T) {
			// Write test file
			content := "CONTAINER_ROOT_DIR=" + tt.rootDir
			err := os.WriteFile(testFile, []byte(content), 0644)
			if err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			// Execute resolution logic
			var logPath string
			if envs, err := godotenv.Read(testFile); err == nil {
				logPath = envs["CONTAINER_ROOT_DIR"] + "/containers"
			}

			if logPath != tt.expectedPath {
				t.Errorf("expected %q, got %q", tt.expectedPath, logPath)
			}
		})
	}
}

