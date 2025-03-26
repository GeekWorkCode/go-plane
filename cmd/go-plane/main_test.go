package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetIssueKeys(t *testing.T) {
	tests := []struct {
		name         string
		ref          string
		issuePattern string
		want         []string
	}{
		{
			name:         "default pattern",
			ref:          "This is a test ABC-1234 and another DEF-5678",
			issuePattern: "",
			want:         []string{"ABC-1234", "DEF-5678"},
		},
		{
			name:         "custom pattern",
			ref:          "This is a test ABC-1234 and another DEF-5678",
			issuePattern: `([A-Z]{3}-[0-9]{4})`,
			want:         []string{"ABC-1234", "DEF-5678"},
		},
		{
			name:         "no matches",
			ref:          "This is a test with no issues",
			issuePattern: "",
			want:         []string{},
		},
		{
			name:         "duplicate issues",
			ref:          "This is a test ABC-1234 and another ABC-1234",
			issuePattern: "",
			want:         []string{"ABC-1234"},
		},
		{
			name:         "multiple issues in commit message",
			ref:          "Fix issues PROJ-123, PROJ-456 and PROJ-789",
			issuePattern: "",
			want:         []string{"PROJ-123", "PROJ-456", "PROJ-789"},
		},
		{
			name:         "issue keys with various formats",
			ref:          "Working on PROJ-123 and PROJECT-456",
			issuePattern: "",
			want:         []string{"PROJ-123", "PROJECT-456"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getIssueKeys(tt.ref, tt.issuePattern)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIssueKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// 保存原始环境变量
	// Save original environment variables
	originalEnv := map[string]string{
		"PLANE_BASE_URL":       os.Getenv("PLANE_BASE_URL"),
		"PLANE_TOKEN":          os.Getenv("PLANE_TOKEN"),
		"PLANE_WORKSPACE_SLUG": os.Getenv("PLANE_WORKSPACE_SLUG"),
		"PLANE_REF":            os.Getenv("PLANE_REF"),
		"PLANE_DEBUG":          os.Getenv("PLANE_DEBUG"),
	}

	// 测试完成后恢复环境变量
	// Restore environment variables after test
	defer func() {
		for k, v := range originalEnv {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	// 设置测试环境变量
	// Set test environment variables
	os.Setenv("PLANE_BASE_URL", "https://test.plane.so/api/v1")
	os.Setenv("PLANE_TOKEN", "test-token")
	os.Setenv("PLANE_WORKSPACE_SLUG", "test-workspace")
	os.Setenv("PLANE_REF", "TEST-123 Test commit")
	os.Setenv("PLANE_DEBUG", "true")

	config := loadConfig()

	// 检查配置是否正确
	// Check if config is correct
	if config.baseURL != "https://test.plane.so/api/v1" {
		t.Errorf("Expected baseURL to be %s, got %s", "https://test.plane.so/api/v1", config.baseURL)
	}
	if config.token != "test-token" {
		t.Errorf("Expected token to be %s, got %s", "test-token", config.token)
	}
	if config.workspaceSlug != "test-workspace" {
		t.Errorf("Expected workspaceSlug to be %s, got %s", "test-workspace", config.workspaceSlug)
	}
	if config.ref != "TEST-123 Test commit" {
		t.Errorf("Expected ref to be %s, got %s", "TEST-123 Test commit", config.ref)
	}
	if !config.debug {
		t.Errorf("Expected debug to be %v, got %v", true, config.debug)
	}
}
