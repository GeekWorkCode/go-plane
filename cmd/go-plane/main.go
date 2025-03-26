package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/GeekWorkCode/go-plane/pkg/markdown"
	"github.com/GeekWorkCode/go-plane/pkg/util"
	"github.com/GeekWorkCode/plane-api-go"
	"github.com/GeekWorkCode/plane-api-go/api"
	"github.com/GeekWorkCode/plane-api-go/models"
	"github.com/joho/godotenv"
)

// Version for command
var Version string

// Commit ID for command
var Commit string

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

func main() {
	// 加载环境变量(从.env文件)
	// Load environment variables (from .env file)
	if err := godotenv.Load(); err != nil {
		// 仅在开发环境中记录错误，生产环境中这不是问题
		// Only log this in development, not an error in production
		log.Printf("Info: .env file not found, using environment variables")
	}

	config := loadConfig()

	// 创建Plane客户端
	// Create Plane client
	planeClient := plane.NewClient(config.token)
	planeClient.SetDebug(config.debug)
	if config.baseURL != "" {
		planeClient.SetBaseURL(config.baseURL)
	}

	// 如果没有ref，则打印版本并退出
	// If no ref, print version and exit
	if config.ref == "" {
		fmt.Printf("go-plane version %s, commit %s\n", Version, Commit)
		os.Exit(0)
	}

	// 处理issue key
	// Process issue keys
	issueKeys := getIssueKeys(config.ref, config.issuePattern)
	if len(issueKeys) == 0 {
		log.Println("未找到issue keys")
		log.Println("No issue keys found")
		os.Exit(0)
	}

	// 获取当前用户信息 - 假设我们无法直接获取
	// Get current user information - assume we can't directly get it
	// In a real implementation, we would add error handling for this placeholder
	self := &User{
		ID:          "current-user",
		Email:       "user@example.com",
		Username:    "current-user",
		DisplayName: "Current User",
	}
	log.Printf("当前用户: %s\n", self.Email)
	log.Printf("Current user: %s\n", self.Email)

	// 获取要分配的用户
	// Get user to assign - this is a placeholder since we don't have User API access
	var assigneeUser *User
	if config.assignee != "" {
		assigneeUser = &User{
			ID:          config.assignee,
			Email:       config.assignee + "@example.com",
			Username:    config.assignee,
			DisplayName: config.assignee,
		}
	}

	// 处理issues
	// Process issues
	issues := processIssues(planeClient, config)
	if len(issues) == 0 {
		log.Println("未找到issues")
		log.Println("No issues found")
		os.Exit(0)
	}

	// 添加评论
	// Add comments
	if config.comment != "" {
		addComments(planeClient, config, issues, self)
	}

	// 更新状态
	// Update state
	if config.toState != "" {
		processState(planeClient, config, issues)
	}

	// 分配责任人
	// Assign issues
	if assigneeUser != nil {
		processAssignee(planeClient, config, issues, assigneeUser)
	}
}

// 处理issues
// Process issues
func processIssues(planeClient *plane.Plane, config Config) []models.Issue {
	var issues []models.Issue

	for _, key := range getIssueKeys(config.ref, config.issuePattern) {
		// 提取项目标识符和问题编号
		// Extract project identifier and issue number
		parts := strings.Split(key, "-")
		if len(parts) != 2 {
			log.Printf("警告: 无效的issue key: %s\n", key)
			log.Printf("Warning: Invalid issue key: %s\n", key)
			continue
		}

		project, err := findProjectByIdentifier(planeClient, config.workspaceSlug, parts[0])
		if err != nil {
			log.Printf("警告: 无法找到项目 '%s': %v\n", parts[0], err)
			log.Printf("Warning: Could not find project '%s': %v\n", parts[0], err)
			continue
		}

		issue, err := findIssueByNumber(planeClient, config.workspaceSlug, project.ID, parts[1])
		if err != nil {
			log.Printf("警告: 无法找到issue '%s': %v\n", key, err)
			log.Printf("Warning: Could not find issue '%s': %v\n", key, err)
			continue
		}

		issues = append(issues, issue)
		log.Printf("找到issue: %s (%s) - %s\n", key, issue.ID, issue.Name)
		log.Printf("Found issue: %s (%s) - %s\n", key, issue.ID, issue.Name)
	}

	return issues
}

// 处理issue分配
// Process issue assignment
func processAssignee(planeClient *plane.Plane, config Config, issues []models.Issue, assignee *User) {
	for _, issue := range issues {
		log.Printf("将issue %s 分配给 %s\n", issue.ID, assignee.DisplayName)
		log.Printf("Assigning issue %s to %s\n", issue.ID, assignee.DisplayName)

		updateReq := &api.IssueUpdateRequest{
			AssigneeID: assignee.ID,
		}

		_, err := planeClient.Issues.Update(config.workspaceSlug, issue.Project, issue.ID, updateReq)
		if err != nil {
			log.Printf("分配issue失败: %v\n", err)
			log.Printf("Failed to assign issue: %v\n", err)
		}
	}
}

// 处理issue状态更新
// Process issue state update
func processState(planeClient *plane.Plane, config Config, issues []models.Issue) {
	for _, issue := range issues {
		states, err := planeClient.States.List(config.workspaceSlug, issue.Project)
		if err != nil {
			log.Printf("获取状态列表失败: %v\n", err)
			log.Printf("Failed to get states: %v\n", err)
			continue
		}

		// 根据名称查找状态ID
		// Find state ID by name
		var stateID string
		for _, state := range states {
			if strings.EqualFold(state.Name, config.toState) {
				stateID = state.ID
				break
			}
		}

		if stateID == "" {
			log.Printf("警告: 无法找到状态 '%s'\n", config.toState)
			log.Printf("Warning: Could not find state '%s'\n", config.toState)
			continue
		}

		log.Printf("将issue %s 状态更新为 %s\n", issue.ID, config.toState)
		log.Printf("Updating issue %s state to %s\n", issue.ID, config.toState)

		updateReq := &api.IssueUpdateRequest{
			State: stateID,
		}

		_, err = planeClient.Issues.Update(config.workspaceSlug, issue.Project, issue.ID, updateReq)
		if err != nil {
			log.Printf("更新状态失败: %v\n", err)
			log.Printf("Failed to update state: %v\n", err)
		}
	}
}

// 获取issue keys
// Get issue keys
func getIssueKeys(ref, issuePattern string) []string {
	if ref == "" {
		return []string{}
	}

	if issuePattern == "" {
		issuePattern = `([A-Z][A-Z0-9]+-[0-9]+)`
	}

	// Remove commas and normalize spaces
	cleanRef := strings.ReplaceAll(ref, ",", " ")
	cleanRef = regexp.MustCompile(`\s+`).ReplaceAllString(cleanRef, " ")

	re := regexp.MustCompile(issuePattern)
	matches := re.FindAllString(cleanRef, -1)

	if matches == nil {
		return []string{}
	}

	// 去重
	// Remove duplicates
	uniqueKeys := make(map[string]struct{})
	for _, match := range matches {
		uniqueKeys[match] = struct{}{}
	}

	var keys []string
	for key := range uniqueKeys {
		keys = append(keys, key)
	}

	return keys
}

// 配置结构体
// Configuration struct
type Config struct {
	baseURL       string
	insecure      string
	token         string
	workspaceSlug string
	ref           string
	issuePattern  string
	toState       string
	comment       string
	assignee      string
	markdown      bool
	debug         bool
}

// 加载配置
// Load configuration
func loadConfig() Config {
	return Config{
		baseURL:       util.GetGlobalValue("PLANE_BASE_URL"),
		insecure:      util.GetGlobalValue("PLANE_INSECURE"),
		token:         util.GetGlobalValue("PLANE_TOKEN"),
		workspaceSlug: util.GetGlobalValue("PLANE_WORKSPACE_SLUG"),
		ref:           util.GetGlobalValue("PLANE_REF"),
		issuePattern:  util.GetGlobalValue("PLANE_ISSUE_PATTERN"),
		toState:       util.GetGlobalValue("PLANE_TO_STATE"),
		comment:       util.GetGlobalValue("PLANE_COMMENT"),
		assignee:      util.GetGlobalValue("PLANE_ASSIGNEE"),
		markdown:      util.ToBool(util.GetGlobalValue("PLANE_MARKDOWN")),
		debug:         util.ToBool(util.GetGlobalValue("PLANE_DEBUG")),
	}
}

// 创建HTTP客户端
// Create HTTP client
func createHTTPClient(config Config) *http.Client {
	transport := &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DisableKeepAlives: false,
	}

	// 如果insecure为true，跳过TLS验证
	// Skip TLS verification if insecure is true
	if util.ToBool(config.insecure) {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// 根据项目标识符查找项目
// Find project by identifier
func findProjectByIdentifier(planeClient *plane.Plane, workspaceSlug, identifier string) (models.Project, error) {
	projects, err := planeClient.Projects.List(workspaceSlug)
	if err != nil {
		return models.Project{}, fmt.Errorf("获取项目列表失败: %w", err)
	}

	for _, project := range projects {
		if strings.EqualFold(project.Identifier, identifier) {
			return project, nil
		}
	}

	return models.Project{}, fmt.Errorf("未找到项目: %s", identifier)
}

// 根据编号查找issue
// Find issue by number
func findIssueByNumber(planeClient *plane.Plane, workspaceSlug, projectID, number string) (models.Issue, error) {
	issues, err := planeClient.Issues.List(workspaceSlug, projectID)
	if err != nil {
		return models.Issue{}, fmt.Errorf("获取issue列表失败: %w", err)
	}

	// 由于models.Issue没有SequenceID字段，我们需要在名称中查找编号
	// Since models.Issue doesn't have a SequenceID field, we need to look for the number in the name
	for _, issue := range issues {
		// 在issue名称中查找编号
		// Look for the number in the issue name
		if strings.Contains(issue.Name, number) {
			return issue, nil
		}
	}

	return models.Issue{}, fmt.Errorf("未找到issue: %s", number)
}

// 添加评论
// Add comments
func addComments(planeClient *plane.Plane, config Config, issues []models.Issue, user *User) {
	for _, issue := range issues {
		var commentText string
		if config.markdown {
			// 将Markdown转换为HTML
			// Convert Markdown to HTML
			commentText = markdown.ToHTML(config.comment)
		} else {
			commentText = config.comment
		}

		log.Printf("为issue %s 添加评论\n", issue.ID)
		log.Printf("Adding comment to issue %s\n", issue.ID)

		_, err := planeClient.Comments.Create(config.workspaceSlug, issue.Project, issue.ID, &api.CommentCreateRequest{
			CommentHTML: commentText,
		})
		if err != nil {
			log.Printf("添加评论失败: %v\n", err)
			log.Printf("Failed to add comment: %v\n", err)
		}
	}
}
