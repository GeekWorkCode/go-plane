package main

import (
	"os"
	"testing"
)

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

func TestFindIssueBySequenceID(t *testing.T) {
	// Skip the test in short mode - this is a placeholder for when proper mocking is implemented
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// This test would require mocking the Plane client's Issues.GetBySequenceID method
	// For proper testing, we would need to:
	// 1. Create a mock Plane client
	// 2. Mock the GetBySequenceID method to return a predefined issue
	// 3. Call findIssueBySequenceID with the mock client
	// 4. Assert that the returned issue matches our expectations

	// This is a placeholder for future implementation
	t.Log("Test for findIssueBySequenceID with GetBySequenceID is a placeholder")
}
