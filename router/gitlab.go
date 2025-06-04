package router

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"gitlab.yogorobot.com/sre/lark-base-mapping/middlewares"
)

// FlexibleTime 自定义时间类型，能够解析多种时间格式
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON 自定义JSON解析，支持多种时间格式
func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	// 移除引号
	timeStr := strings.Trim(string(data), `"`)

	// 尝试多种时间格式
	formats := []string{
		time.RFC3339,              // 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano,          // 2006-01-02T15:04:05.999999999Z07:00
		"2006-01-02T15:04:05Z",    // 2006-01-02T15:04:05Z
		"2006-01-02T15:04:05",     // 2006-01-02T15:04:05
		"2006-01-02 15:04:05 UTC", // 2006-01-02 15:04:05 UTC (GitLab格式)
		"2006-01-02 15:04:05",     // 2006-01-02 15:04:05
	}

	var err error
	for _, format := range formats {
		ft.Time, err = time.Parse(format, timeStr)
		if err == nil {
			return nil
		}
	}

	// 如果所有格式都解析失败，返回最后一个错误
	return err
}

// GitLabMergeRequestEvent GitLab Merge Request事件数据结构
type GitLabMergeRequestEvent struct {
	ObjectKind       string                 `json:"object_kind"`
	EventType        string                 `json:"event_type"`
	User             User                   `json:"user"`
	Project          Project                `json:"project"`
	ObjectAttributes MergeRequestAttributes `json:"object_attributes"`
	Labels           []Label                `json:"labels"`
	Changes          Changes                `json:"changes"`
	Repository       Repository             `json:"repository"`
}

// SystemHookMergeRequestEvent System Hook格式的Merge Request事件数据结构
type SystemHookMergeRequestEvent struct {
	ObjectKind       string                           `json:"object_kind"`
	EventType        string                           `json:"event_type"`
	User             User                             `json:"user"`
	Project          Project                          `json:"project"`
	ObjectAttributes SystemHookMergeRequestAttributes `json:"object_attributes"`
	Labels           []Label                          `json:"labels"`
	Changes          SystemHookChanges                `json:"changes"`
	Repository       Repository                       `json:"repository"`
}

// SystemHookMergeRequestAttributes System Hook格式的MR属性
type SystemHookMergeRequestAttributes struct {
	ID                          int     `json:"id"`
	IID                         int     `json:"iid"`
	Title                       string  `json:"title"`
	Description                 string  `json:"description"`
	State                       string  `json:"state"`
	CreatedAt                   string  `json:"created_at"` // System Hook使用字符串格式
	UpdatedAt                   string  `json:"updated_at"` // System Hook使用字符串格式
	MergeStatus                 string  `json:"merge_status"`
	TargetBranch                string  `json:"target_branch"`
	SourceBranch                string  `json:"source_branch"`
	SourceProjectID             int     `json:"source_project_id"`
	TargetProjectID             int     `json:"target_project_id"`
	URL                         string  `json:"url"`
	Source                      Project `json:"source"`
	Target                      Project `json:"target"`
	LastCommit                  Commit  `json:"last_commit"`
	WorkInProgress              bool    `json:"work_in_progress"`
	Assignee                    User    `json:"assignee"`
	Author                      User    `json:"author"`
	MergeCommitSHA              string  `json:"merge_commit_sha"`
	BlockingDiscussionsResolved bool    `json:"blocking_discussions_resolved"`
	Action                      string  `json:"action"`
}

// SystemHookChanges System Hook格式的变更信息
type SystemHookChanges struct {
	Title       SystemHookChange      `json:"title"`
	Description SystemHookChange      `json:"description"`
	Labels      SystemHookLabelChange `json:"labels"`
	State       SystemHookChange      `json:"state"`
	UpdatedAt   SystemHookChange      `json:"updated_at"`
}

// SystemHookChange System Hook格式的单个变更
type SystemHookChange struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

// SystemHookLabelChange System Hook格式的标签变更
type SystemHookLabelChange struct {
	Previous []Label `json:"previous"`
	Current  []Label `json:"current"`
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type Project struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int    `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

type MergeRequestAttributes struct {
	ID                          int          `json:"id"`
	IID                         int          `json:"iid"`
	Title                       string       `json:"title"`
	Description                 string       `json:"description"`
	State                       string       `json:"state"`
	CreatedAt                   FlexibleTime `json:"created_at"`
	UpdatedAt                   FlexibleTime `json:"updated_at"`
	MergeStatus                 string       `json:"merge_status"`
	TargetBranch                string       `json:"target_branch"`
	SourceBranch                string       `json:"source_branch"`
	SourceProjectID             int          `json:"source_project_id"`
	TargetProjectID             int          `json:"target_project_id"`
	URL                         string       `json:"url"`
	Source                      Project      `json:"source"`
	Target                      Project      `json:"target"`
	LastCommit                  Commit       `json:"last_commit"`
	WorkInProgress              bool         `json:"work_in_progress"`
	Assignee                    User         `json:"assignee"`
	Author                      User         `json:"author"`
	MergeCommitSHA              string       `json:"merge_commit_sha"`
	BlockingDiscussionsResolved bool         `json:"blocking_discussions_resolved"`
	Action                      string       `json:"action"`
}

type Label struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Color       string       `json:"color"`
	ProjectID   int          `json:"project_id"`
	CreatedAt   FlexibleTime `json:"created_at"`
	UpdatedAt   FlexibleTime `json:"updated_at"`
	Template    bool         `json:"template"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	GroupID     int          `json:"group_id"`
}

type Changes struct {
	Title       Change      `json:"title"`
	Description Change      `json:"description"`
	Labels      LabelChange `json:"labels"`
	State       Change      `json:"state"`
	UpdatedAt   Change      `json:"updated_at"`
}

type Change struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type LabelChange struct {
	Previous []Label `json:"previous"`
	Current  []Label `json:"current"`
}

type Repository struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

type Commit struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	Timestamp FlexibleTime `json:"timestamp"`
	URL       string       `json:"url"`
	Author    Author       `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// System Hook Event 数据结构
type SystemHookEvent struct {
	EventName  string `json:"event_name"`
	ObjectKind string `json:"object_kind,omitempty"` // 新格式事件使用
	Action     string `json:"action,omitempty"`      // 新格式事件使用
}

// Project System Hook Events
type ProjectSystemHookEvent struct {
	EventName            string       `json:"event_name"`
	CreatedAt            FlexibleTime `json:"created_at"`
	UpdatedAt            FlexibleTime `json:"updated_at"`
	Name                 string       `json:"name"`
	OwnerEmail           string       `json:"owner_email"`
	OwnerName            string       `json:"owner_name"`
	Owners               []Owner      `json:"owners"`
	Path                 string       `json:"path"`
	PathWithNamespace    string       `json:"path_with_namespace"`
	ProjectID            int          `json:"project_id"`
	ProjectNamespaceID   int          `json:"project_namespace_id"`
	ProjectVisibility    string       `json:"project_visibility"`
	OldPathWithNamespace string       `json:"old_path_with_namespace,omitempty"` // for rename/transfer
}

// User System Hook Events
type UserSystemHookEvent struct {
	EventName    string       `json:"event_name"`
	CreatedAt    FlexibleTime `json:"created_at"`
	UpdatedAt    FlexibleTime `json:"updated_at"`
	UserEmail    string       `json:"user_email"`
	UserName     string       `json:"user_name"`
	UserUsername string       `json:"user_username"`
	UserID       int          `json:"user_id"`
	OldUsername  string       `json:"old_username,omitempty"` // for user_rename
}

// Group System Hook Events
type GroupSystemHookEvent struct {
	EventName            string       `json:"event_name"`
	CreatedAt            FlexibleTime `json:"created_at"`
	UpdatedAt            FlexibleTime `json:"updated_at"`
	Name                 string       `json:"name"`
	Path                 string       `json:"path"`
	PathWithNamespace    string       `json:"path_with_namespace"`
	GroupID              int          `json:"group_id"`
	OwnerEmail           string       `json:"owner_email,omitempty"`
	OwnerName            string       `json:"owner_name,omitempty"`
	OldPath              string       `json:"old_path,omitempty"`
	OldPathWithNamespace string       `json:"old_path_with_namespace,omitempty"`
}

// Repository Update Event
type RepositoryUpdateEvent struct {
	EventName  string      `json:"event_name"`
	UserID     int         `json:"user_id"`
	UserName   string      `json:"user_name"`
	UserEmail  string      `json:"user_email"`
	UserAvatar string      `json:"user_avatar"`
	ProjectID  int         `json:"project_id"`
	Project    Project     `json:"project"`
	Changes    []RefChange `json:"changes"`
	Refs       []string    `json:"refs"`
}

type RefChange struct {
	Before string `json:"before"`
	After  string `json:"after"`
	Ref    string `json:"ref"`
}

// Access Request Events
type AccessRequestEvent struct {
	EventName     string       `json:"event_name"`
	CreatedAt     FlexibleTime `json:"created_at"`
	UpdatedAt     FlexibleTime `json:"updated_at"`
	GroupAccess   string       `json:"group_access,omitempty"`
	ProjectAccess string       `json:"project_access,omitempty"`
	GroupID       int          `json:"group_id,omitempty"`
	ProjectID     int          `json:"project_id,omitempty"`
	GroupName     string       `json:"group_name,omitempty"`
	ProjectName   string       `json:"project_name,omitempty"`
	GroupPath     string       `json:"group_path,omitempty"`
	ProjectPath   string       `json:"project_path,omitempty"`
	UserEmail     string       `json:"user_email"`
	UserName      string       `json:"user_name"`
	UserUsername  string       `json:"user_username"`
	UserID        int          `json:"user_id"`
}

// Key Events
type KeyEvent struct {
	EventName string       `json:"event_name"`
	CreatedAt FlexibleTime `json:"created_at"`
	UpdatedAt FlexibleTime `json:"updated_at"`
	UserName  string       `json:"user_name"`
	UserEmail string       `json:"user_email"`
	UserID    int          `json:"user_id"`
	KeyID     int          `json:"key_id"`
}

// Member Approval Events (新格式)
type MemberApprovalEvent struct {
	ObjectKind           string                   `json:"object_kind"`
	Action               string                   `json:"action"`
	ObjectAttributes     MemberApprovalAttributes `json:"object_attributes"`
	UserID               int                      `json:"user_id"`
	RequestedByUserID    int                      `json:"requested_by_user_id,omitempty"`
	ReviewedByUserID     int                      `json:"reviewed_by_user_id,omitempty"`
	PromotionNamespaceID int                      `json:"promotion_namespace_id,omitempty"`
	CreatedAt            FlexibleTime             `json:"created_at"`
	UpdatedAt            FlexibleTime             `json:"updated_at"`
}

type MemberApprovalAttributes struct {
	NewAccessLevel                       int    `json:"new_access_level,omitempty"`
	OldAccessLevel                       int    `json:"old_access_level,omitempty"`
	ExistingMemberID                     int    `json:"existing_member_id,omitempty"`
	PromotionRequestIDsThatFailedToApply []int  `json:"promotion_request_ids_that_failed_to_apply,omitempty"`
	Status                               string `json:"status,omitempty"`
}

type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GitLabWebhook 处理GitLab webhook事件
func GitLabWebhook(e *core.RequestEvent) error {
	app := e.App

	// 从中间件上下文中获取GitLab配置
	gitlabConfig, ok := middlewares.GetGitLabConfigFromContext(e.Request.Context())
	if !ok {
		return e.BadRequestError("GitLab config not found in context", nil)
	}

	// 获取事件类型
	eventType := e.Request.Header.Get("X-Gitlab-Event")
	if eventType == "" {
		return e.BadRequestError("Missing X-Gitlab-Event header", nil)
	}

	app.Logger().Info("Processing GitLab webhook",
		"eventType", eventType,
		"gitlabURL", gitlabConfig.BaseURL,
	)

	// 读取请求体
	body, err := io.ReadAll(e.Request.Body)
	if err != nil {
		app.Logger().Error("Failed to read request body", "error", err)
		return e.BadRequestError("Failed to read request body", err)
	}

	// 根据事件类型处理
	switch eventType {
	case "System Hook":
		return handleSystemHookEvent(e, body)
	case "Merge Request Hook":
		return handleMergeRequestEvent(e, body)
	case "Note Hook":
		return handleNoteEvent(e, body)
	case "Push Hook":
		return handlePushEvent(e, body)
	case "Tag Push Hook":
		return handleTagPushEvent(e, body)
	case "Issues Hook":
		return handleIssuesEvent(e, body)
	default:
		app.Logger().Info("Unsupported GitLab event type", "eventType", eventType)
		return e.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Event received but not processed",
			"event":   eventType,
		})
	}
}

// handleSystemHookEvent 处理System Hook事件
func handleSystemHookEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	// 先解析基本的事件信息来确定事件类型
	var baseEvent SystemHookEvent
	if err := json.Unmarshal(body, &baseEvent); err != nil {
		app.Logger().Error("Failed to parse system hook event", "error", err)
		return e.BadRequestError("Invalid system hook event format", err)
	}

	app.Logger().Info("Processing system hook event",
		"eventName", baseEvent.EventName,
		"objectKind", baseEvent.ObjectKind,
		"action", baseEvent.Action,
	)

	// 根据事件名称或对象类型分发处理
	if baseEvent.ObjectKind != "" {
		// 新格式事件
		return handleNewFormatSystemEvent(e, body, baseEvent)
	} else {
		// 传统格式事件
		return handleTraditionalSystemEvent(e, body, baseEvent.EventName)
	}
}

// handleNewFormatSystemEvent 处理新格式的系统事件
func handleNewFormatSystemEvent(e *core.RequestEvent, body []byte, baseEvent SystemHookEvent) error {
	app := e.App

	switch baseEvent.ObjectKind {
	case "gitlab_subscription_member_approval", "gitlab_subscription_member_approvals":
		return handleMemberApprovalEvent(e, body, baseEvent)
	case "merge_request":
		// System Hook格式的Merge Request事件，使用专门的处理器
		app.Logger().Info("Processing merge request system hook event",
			"objectKind", baseEvent.ObjectKind,
			"action", baseEvent.Action,
		)
		return handleSystemHookMergeRequestEvent(e, body)
	default:
		app.Logger().Info("Unsupported new format system event",
			"objectKind", baseEvent.ObjectKind,
			"action", baseEvent.Action,
		)
		return e.JSON(http.StatusOK, map[string]interface{}{
			"status":      "success",
			"message":     "New format system event received but not processed",
			"object_kind": baseEvent.ObjectKind,
			"action":      baseEvent.Action,
		})
	}
}

// handleTraditionalSystemEvent 处理传统格式的系统事件
func handleTraditionalSystemEvent(e *core.RequestEvent, body []byte, eventName string) error {
	app := e.App

	switch eventName {
	// 项目相关事件
	case "project_create", "project_destroy", "project_rename", "project_transfer", "project_update":
		return handleProjectSystemEvent(e, body, eventName)

	// 用户相关事件
	case "user_create", "user_destroy", "user_rename", "user_failed_login":
		return handleUserSystemEvent(e, body, eventName)

	// 组相关事件
	case "group_create", "group_destroy", "group_rename":
		return handleGroupSystemEvent(e, body, eventName)

	// 访问请求事件
	case "user_access_request_revoked_for_group", "user_access_request_revoked_for_project",
		"user_access_request_to_group", "user_access_request_to_project",
		"user_add_to_group", "user_add_to_team", "user_remove_from_group",
		"user_remove_from_team", "user_update_for_group", "user_update_for_team":
		return handleAccessRequestEvent(e, body, eventName)

	// 密钥事件
	case "key_create", "key_destroy":
		return handleKeyEvent(e, body, eventName)

	// 仓库更新事件
	case "repository_update":
		return handleRepositoryUpdateEvent(e, body)

	default:
		app.Logger().Info("Unsupported system hook event", "eventName", eventName)
		return e.JSON(http.StatusOK, map[string]interface{}{
			"status":     "success",
			"message":    "System hook event received but not processed",
			"event_name": eventName,
		})
	}
}

// handleProjectSystemEvent 处理项目系统事件
func handleProjectSystemEvent(e *core.RequestEvent, body []byte, eventName string) error {
	app := e.App

	var event ProjectSystemHookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse project system event", "error", err, "eventName", eventName)
		return e.BadRequestError("Invalid project system event format", err)
	}

	app.Logger().Info("Processing project system event",
		"eventName", eventName,
		"projectID", event.ProjectID,
		"projectName", event.Name,
		"path", event.PathWithNamespace,
		"visibility", event.ProjectVisibility,
	)

	// 保存项目系统事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_project_system_events")
	if err != nil {
		app.Logger().Warn("gitlab_project_system_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("event_name", event.EventName)
		record.Set("project_id", event.ProjectID)
		record.Set("project_name", event.Name)
		record.Set("path", event.Path)
		record.Set("path_with_namespace", event.PathWithNamespace)
		record.Set("project_visibility", event.ProjectVisibility)
		record.Set("owner_name", event.OwnerName)
		record.Set("owner_email", event.OwnerEmail)
		if event.OldPathWithNamespace != "" {
			record.Set("old_path_with_namespace", event.OldPathWithNamespace)
		}
		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save project system event", "error", err)
		} else {
			app.Logger().Info("Project system event saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Project system event processed",
		"event": map[string]interface{}{
			"event_name":   eventName,
			"project_id":   event.ProjectID,
			"project_name": event.Name,
			"path":         event.PathWithNamespace,
		},
	})
}

// handleUserSystemEvent 处理用户系统事件
func handleUserSystemEvent(e *core.RequestEvent, body []byte, eventName string) error {
	app := e.App

	var event UserSystemHookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse user system event", "error", err, "eventName", eventName)
		return e.BadRequestError("Invalid user system event format", err)
	}

	app.Logger().Info("Processing user system event",
		"eventName", eventName,
		"userID", event.UserID,
		"userName", event.UserName,
		"userEmail", event.UserEmail,
		"username", event.UserUsername,
	)

	// 保存用户系统事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_user_system_events")
	if err != nil {
		app.Logger().Warn("gitlab_user_system_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("event_name", event.EventName)
		record.Set("user_id", event.UserID)
		record.Set("user_name", event.UserName)
		record.Set("user_email", event.UserEmail)
		record.Set("user_username", event.UserUsername)
		if event.OldUsername != "" {
			record.Set("old_username", event.OldUsername)
		}
		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save user system event", "error", err)
		} else {
			app.Logger().Info("User system event saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "User system event processed",
		"event": map[string]interface{}{
			"event_name": eventName,
			"user_id":    event.UserID,
			"user_name":  event.UserName,
			"username":   event.UserUsername,
		},
	})
}

// handleGroupSystemEvent 处理组系统事件
func handleGroupSystemEvent(e *core.RequestEvent, body []byte, eventName string) error {
	app := e.App

	var event GroupSystemHookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse group system event", "error", err, "eventName", eventName)
		return e.BadRequestError("Invalid group system event format", err)
	}

	app.Logger().Info("Processing group system event",
		"eventName", eventName,
		"groupID", event.GroupID,
		"groupName", event.Name,
		"path", event.PathWithNamespace,
	)

	// 保存组系统事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_group_system_events")
	if err != nil {
		app.Logger().Warn("gitlab_group_system_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("event_name", event.EventName)
		record.Set("group_id", event.GroupID)
		record.Set("group_name", event.Name)
		record.Set("path", event.Path)
		record.Set("path_with_namespace", event.PathWithNamespace)
		if event.OldPath != "" {
			record.Set("old_path", event.OldPath)
		}
		if event.OldPathWithNamespace != "" {
			record.Set("old_path_with_namespace", event.OldPathWithNamespace)
		}
		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save group system event", "error", err)
		} else {
			app.Logger().Info("Group system event saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Group system event processed",
		"event": map[string]interface{}{
			"event_name": eventName,
			"group_id":   event.GroupID,
			"group_name": event.Name,
			"path":       event.PathWithNamespace,
		},
	})
}

// handleAccessRequestEvent 处理访问请求事件
func handleAccessRequestEvent(e *core.RequestEvent, body []byte, eventName string) error {
	app := e.App

	var event AccessRequestEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse access request event", "error", err, "eventName", eventName)
		return e.BadRequestError("Invalid access request event format", err)
	}

	app.Logger().Info("Processing access request event",
		"eventName", eventName,
		"userID", event.UserID,
		"userName", event.UserName,
		"groupID", event.GroupID,
		"projectID", event.ProjectID,
	)

	// 保存访问请求事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_access_request_events")
	if err != nil {
		app.Logger().Warn("gitlab_access_request_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("event_name", event.EventName)
		record.Set("user_id", event.UserID)
		record.Set("user_name", event.UserName)
		record.Set("user_email", event.UserEmail)
		record.Set("user_username", event.UserUsername)
		if event.GroupID > 0 {
			record.Set("group_id", event.GroupID)
			record.Set("group_name", event.GroupName)
			record.Set("group_path", event.GroupPath)
			record.Set("group_access", event.GroupAccess)
		}
		if event.ProjectID > 0 {
			record.Set("project_id", event.ProjectID)
			record.Set("project_name", event.ProjectName)
			record.Set("project_path", event.ProjectPath)
			record.Set("project_access", event.ProjectAccess)
		}
		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save access request event", "error", err)
		} else {
			app.Logger().Info("Access request event saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Access request event processed",
		"event": map[string]interface{}{
			"event_name": eventName,
			"user_id":    event.UserID,
			"user_name":  event.UserName,
		},
	})
}

// handleKeyEvent 处理密钥事件
func handleKeyEvent(e *core.RequestEvent, body []byte, eventName string) error {
	app := e.App

	var event KeyEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse key event", "error", err, "eventName", eventName)
		return e.BadRequestError("Invalid key event format", err)
	}

	app.Logger().Info("Processing key event",
		"eventName", eventName,
		"userID", event.UserID,
		"userName", event.UserName,
		"keyID", event.KeyID,
	)

	// 保存密钥事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_key_events")
	if err != nil {
		app.Logger().Warn("gitlab_key_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("event_name", event.EventName)
		record.Set("user_id", event.UserID)
		record.Set("user_name", event.UserName)
		record.Set("user_email", event.UserEmail)
		record.Set("key_id", event.KeyID)
		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save key event", "error", err)
		} else {
			app.Logger().Info("Key event saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Key event processed",
		"event": map[string]interface{}{
			"event_name": eventName,
			"user_id":    event.UserID,
			"key_id":     event.KeyID,
		},
	})
}

// handleRepositoryUpdateEvent 处理仓库更新事件
func handleRepositoryUpdateEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	var event RepositoryUpdateEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse repository update event", "error", err)
		return e.BadRequestError("Invalid repository update event format", err)
	}

	app.Logger().Info("Processing repository update event",
		"eventName", event.EventName,
		"userID", event.UserID,
		"userName", event.UserName,
		"projectID", event.ProjectID,
		"projectName", event.Project.Name,
		"refsCount", len(event.Refs),
	)

	// 保存仓库更新事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_repository_update_events")
	if err != nil {
		app.Logger().Warn("gitlab_repository_update_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("event_name", event.EventName)
		record.Set("user_id", event.UserID)
		record.Set("user_name", event.UserName)
		record.Set("user_email", event.UserEmail)
		record.Set("project_id", event.ProjectID)
		record.Set("project_name", event.Project.Name)
		record.Set("project_path", event.Project.PathWithNamespace)

		// 将 refs 和 changes 序列化为 JSON
		if refsJSON, err := json.Marshal(event.Refs); err == nil {
			record.Set("refs", string(refsJSON))
		}
		if changesJSON, err := json.Marshal(event.Changes); err == nil {
			record.Set("changes", string(changesJSON))
		}

		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save repository update event", "error", err)
		} else {
			app.Logger().Info("Repository update event saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Repository update event processed",
		"event": map[string]interface{}{
			"event_name":   event.EventName,
			"user_id":      event.UserID,
			"project_id":   event.ProjectID,
			"project_name": event.Project.Name,
			"refs_count":   len(event.Refs),
		},
	})
}

// handleMemberApprovalEvent 处理成员审批事件 (新格式)
func handleMemberApprovalEvent(e *core.RequestEvent, body []byte, baseEvent SystemHookEvent) error {
	app := e.App

	var event MemberApprovalEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse member approval event", "error", err)
		return e.BadRequestError("Invalid member approval event format", err)
	}

	app.Logger().Info("Processing member approval event",
		"objectKind", event.ObjectKind,
		"action", event.Action,
		"userID", event.UserID,
		"reviewedByUserID", event.ReviewedByUserID,
		"status", event.ObjectAttributes.Status,
	)

	// 保存成员审批事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_member_approval_events")
	if err != nil {
		app.Logger().Warn("gitlab_member_approval_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("object_kind", event.ObjectKind)
		record.Set("action", event.Action)
		record.Set("user_id", event.UserID)
		if event.RequestedByUserID > 0 {
			record.Set("requested_by_user_id", event.RequestedByUserID)
		}
		if event.ReviewedByUserID > 0 {
			record.Set("reviewed_by_user_id", event.ReviewedByUserID)
		}
		if event.PromotionNamespaceID > 0 {
			record.Set("promotion_namespace_id", event.PromotionNamespaceID)
		}
		if event.ObjectAttributes.Status != "" {
			record.Set("status", event.ObjectAttributes.Status)
		}
		if event.ObjectAttributes.NewAccessLevel > 0 {
			record.Set("new_access_level", event.ObjectAttributes.NewAccessLevel)
		}
		if event.ObjectAttributes.OldAccessLevel > 0 {
			record.Set("old_access_level", event.ObjectAttributes.OldAccessLevel)
		}
		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save member approval event", "error", err)
		} else {
			app.Logger().Info("Member approval event saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Member approval event processed",
		"event": map[string]interface{}{
			"object_kind": event.ObjectKind,
			"action":      event.Action,
			"user_id":     event.UserID,
			"status":      event.ObjectAttributes.Status,
		},
	})
}

// handleMergeRequestEvent 处理Merge Request事件
func handleMergeRequestEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	var event GitLabMergeRequestEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse merge request event", "error", err)
		return e.BadRequestError("Invalid merge request event format", err)
	}

	app.Logger().Info("Processing merge request event",
		"action", event.ObjectAttributes.Action,
		"mrID", event.ObjectAttributes.IID,
		"title", event.ObjectAttributes.Title,
		"state", event.ObjectAttributes.State,
		"author", event.ObjectAttributes.Author.Name,
		"targetBranch", event.ObjectAttributes.TargetBranch,
		"sourceBranch", event.ObjectAttributes.SourceBranch,
		"projectName", event.Project.Name,
	)

	// 这里可以添加具体的业务逻辑
	// 例如：
	// 1. 保存MR信息到数据库
	// 2. 触发自动化流程
	// 3. 发送通知到飞书
	// 4. 执行代码质量检查

	// 示例：保存MR事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_merge_requests")
	if err != nil {
		// 如果表不存在，先创建（这里简化处理，实际应该通过迁移创建）
		app.Logger().Warn("gitlab_merge_requests collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("mr_id", event.ObjectAttributes.ID)
		record.Set("mr_iid", event.ObjectAttributes.IID)
		record.Set("title", event.ObjectAttributes.Title)
		record.Set("description", event.ObjectAttributes.Description)
		record.Set("state", event.ObjectAttributes.State)
		record.Set("action", event.ObjectAttributes.Action)

		// 安全设置author信息，处理空值情况
		if event.ObjectAttributes.Author.ID > 0 && event.ObjectAttributes.Author.Name != "" {
			record.Set("author_name", event.ObjectAttributes.Author.Name)
		} else if event.User.ID > 0 && event.User.Name != "" {
			// 使用事件触发用户作为备用author
			app.Logger().Info("Using event user as author fallback",
				"eventUserID", event.User.ID,
				"eventUserName", event.User.Name,
				"mrID", event.ObjectAttributes.IID)
			record.Set("author_name", event.User.Name)
		} else {
			app.Logger().Warn("Both author and event user are invalid, using default",
				"authorID", event.ObjectAttributes.Author.ID,
				"authorName", event.ObjectAttributes.Author.Name,
				"eventUserID", event.User.ID,
				"eventUserName", event.User.Name,
				"mrID", event.ObjectAttributes.IID)
			record.Set("author_name", "Unknown Author")
		}

		if event.ObjectAttributes.Author.ID > 0 && event.ObjectAttributes.Author.Username != "" {
			record.Set("author_username", event.ObjectAttributes.Author.Username)
		} else if event.User.ID > 0 && event.User.Username != "" {
			// 使用事件触发用户作为备用author
			app.Logger().Info("Using event user username as author fallback",
				"eventUserID", event.User.ID,
				"eventUserUsername", event.User.Username,
				"mrID", event.ObjectAttributes.IID)
			record.Set("author_username", event.User.Username)
		} else {
			app.Logger().Warn("Both author and event user username are invalid, using default",
				"authorID", event.ObjectAttributes.Author.ID,
				"authorUsername", event.ObjectAttributes.Author.Username,
				"eventUserID", event.User.ID,
				"eventUserUsername", event.User.Username,
				"mrID", event.ObjectAttributes.IID)
			record.Set("author_username", "unknown")
		}

		record.Set("project_id", event.Project.ID)
		record.Set("project_name", event.Project.Name)
		record.Set("source_branch", event.ObjectAttributes.SourceBranch)
		record.Set("target_branch", event.ObjectAttributes.TargetBranch)
		record.Set("url", event.ObjectAttributes.URL)
		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save merge request record", "error", err)
		} else {
			app.Logger().Info("Merge request record saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Merge request event processed",
		"event": map[string]interface{}{
			"action":  event.ObjectAttributes.Action,
			"mr_id":   event.ObjectAttributes.IID,
			"title":   event.ObjectAttributes.Title,
			"state":   event.ObjectAttributes.State,
			"project": event.Project.Name,
		},
	})
}

// handleSystemHookMergeRequestEvent 处理System Hook格式的Merge Request事件
func handleSystemHookMergeRequestEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	var event SystemHookMergeRequestEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse system hook merge request event", "error", err)
		return e.BadRequestError("Invalid system hook merge request event format", err)
	}

	app.Logger().Info("Processing system hook merge request event",
		"action", event.ObjectAttributes.Action,
		"mrID", event.ObjectAttributes.IID,
		"title", event.ObjectAttributes.Title,
		"state", event.ObjectAttributes.State,
		"author", event.ObjectAttributes.Author.Name,
		"targetBranch", event.ObjectAttributes.TargetBranch,
		"sourceBranch", event.ObjectAttributes.SourceBranch,
		"projectName", event.Project.Name,
	)

	// 保存System Hook MR事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_merge_requests")
	if err != nil {
		app.Logger().Warn("gitlab_merge_requests collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("mr_id", event.ObjectAttributes.ID)
		record.Set("mr_iid", event.ObjectAttributes.IID)
		record.Set("title", event.ObjectAttributes.Title)
		record.Set("description", event.ObjectAttributes.Description)
		record.Set("state", event.ObjectAttributes.State)
		record.Set("action", event.ObjectAttributes.Action)

		// 安全设置author信息，处理空值情况
		if event.ObjectAttributes.Author.ID > 0 && event.ObjectAttributes.Author.Name != "" {
			record.Set("author_name", event.ObjectAttributes.Author.Name)
		} else if event.User.ID > 0 && event.User.Name != "" {
			// 使用事件触发用户作为备用author
			app.Logger().Info("Using event user as author fallback in system hook",
				"eventUserID", event.User.ID,
				"eventUserName", event.User.Name,
				"mrID", event.ObjectAttributes.IID)
			record.Set("author_name", event.User.Name)
		} else {
			app.Logger().Warn("Both author and event user are invalid in system hook, using default",
				"authorID", event.ObjectAttributes.Author.ID,
				"authorName", event.ObjectAttributes.Author.Name,
				"eventUserID", event.User.ID,
				"eventUserName", event.User.Name,
				"mrID", event.ObjectAttributes.IID)
			record.Set("author_name", "Unknown Author")
		}

		if event.ObjectAttributes.Author.ID > 0 && event.ObjectAttributes.Author.Username != "" {
			record.Set("author_username", event.ObjectAttributes.Author.Username)
		} else if event.User.ID > 0 && event.User.Username != "" {
			// 使用事件触发用户作为备用author
			app.Logger().Info("Using event user username as author fallback in system hook",
				"eventUserID", event.User.ID,
				"eventUserUsername", event.User.Username,
				"mrID", event.ObjectAttributes.IID)
			record.Set("author_username", event.User.Username)
		} else {
			app.Logger().Warn("Both author and event user username are invalid in system hook, using default",
				"authorID", event.ObjectAttributes.Author.ID,
				"authorUsername", event.ObjectAttributes.Author.Username,
				"eventUserID", event.User.ID,
				"eventUserUsername", event.User.Username,
				"mrID", event.ObjectAttributes.IID)
			record.Set("author_username", "unknown")
		}

		record.Set("project_id", event.Project.ID)
		record.Set("project_name", event.Project.Name)
		record.Set("source_branch", event.ObjectAttributes.SourceBranch)
		record.Set("target_branch", event.ObjectAttributes.TargetBranch)
		record.Set("url", event.ObjectAttributes.URL)
		record.Set("created_at", event.ObjectAttributes.CreatedAt) // 保存原始时间字符串
		record.Set("updated_at", event.ObjectAttributes.UpdatedAt) // 保存原始时间字符串
		record.Set("event_source", "system_hook")                  // 标记事件来源
		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save system hook merge request record", "error", err)
		} else {
			app.Logger().Info("System hook merge request record saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "System hook merge request event processed",
		"event": map[string]interface{}{
			"action":  event.ObjectAttributes.Action,
			"mr_id":   event.ObjectAttributes.IID,
			"title":   event.ObjectAttributes.Title,
			"state":   event.ObjectAttributes.State,
			"project": event.Project.Name,
			"source":  "system_hook",
		},
	})
}

// handlePushEvent 处理Push事件（占位符）
func handlePushEvent(e *core.RequestEvent, body []byte) error {
	e.App.Logger().Info("Push event received")
	// TODO: 实现Push事件处理逻辑
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Push event received",
	})
}

// handleTagPushEvent 处理Tag Push事件（占位符）
func handleTagPushEvent(e *core.RequestEvent, body []byte) error {
	e.App.Logger().Info("Tag push event received")
	// TODO: 实现Tag Push事件处理逻辑
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Tag push event received",
	})
}

// handleIssuesEvent 处理Issues事件（占位符）
func handleIssuesEvent(e *core.RequestEvent, body []byte) error {
	e.App.Logger().Info("Issues event received")
	// TODO: 实现Issues事件处理逻辑
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Issues event received",
	})
}

// handleNoteEvent 处理Note事件（评论事件）
func handleNoteEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	var event GitLabNoteEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse note event", "error", err)
		return e.BadRequestError("Invalid note event format", err)
	}

	app.Logger().Info("Processing note event",
		"action", event.ObjectAttributes.Action,
		"noteID", event.ObjectAttributes.ID,
		"noteableType", event.ObjectAttributes.NoteableType,
		"authorID", event.ObjectAttributes.AuthorID,
		"projectName", event.Project.Name,
		"noteContent", func() string {
			if len(event.ObjectAttributes.Note) > 50 {
				return event.ObjectAttributes.Note[:50] + "..."
			}
			return event.ObjectAttributes.Note
		}(),
	)

	// 根据评论类型添加额外的日志信息
	switch event.ObjectAttributes.NoteableType {
	case "MergeRequest":
		if event.MergeRequest != nil {
			app.Logger().Info("Note on merge request",
				"mrID", event.MergeRequest.IID,
				"mrTitle", event.MergeRequest.Title,
				"mrState", event.MergeRequest.State,
			)
		}
	case "Issue":
		if event.Issue != nil {
			app.Logger().Info("Note on issue",
				"issueID", event.Issue.IID,
				"issueTitle", event.Issue.Title,
				"issueState", event.Issue.State,
			)
		}
	case "Commit":
		if event.Commit != nil {
			app.Logger().Info("Note on commit",
				"commitID", event.Commit.ID,
				"commitMessage", func() string {
					if len(event.Commit.Message) > 50 {
						return event.Commit.Message[:50] + "..."
					}
					return event.Commit.Message
				}(),
			)
		}
	case "Snippet":
		if event.Snippet != nil {
			app.Logger().Info("Note on snippet",
				"snippetID", event.Snippet.ID,
				"snippetTitle", event.Snippet.Title,
			)
		}
	}

	// 保存Note事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_note_events")
	if err != nil {
		app.Logger().Warn("gitlab_note_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("note_id", event.ObjectAttributes.ID)
		record.Set("note_content", event.ObjectAttributes.Note)
		record.Set("noteable_type", event.ObjectAttributes.NoteableType)
		record.Set("author_id", event.ObjectAttributes.AuthorID)
		record.Set("project_id", event.Project.ID)
		record.Set("project_name", event.Project.Name)
		record.Set("action", event.ObjectAttributes.Action)
		record.Set("system", event.ObjectAttributes.System)
		record.Set("created_at", event.ObjectAttributes.CreatedAt)
		record.Set("updated_at", event.ObjectAttributes.UpdatedAt)
		record.Set("url", event.ObjectAttributes.URL)

		// 根据评论类型设置相关的ID和信息
		switch event.ObjectAttributes.NoteableType {
		case "MergeRequest":
			if event.MergeRequest != nil {
				record.Set("noteable_id", event.MergeRequest.IID)
				record.Set("noteable_title", event.MergeRequest.Title)
				record.Set("noteable_state", event.MergeRequest.State)
			}
		case "Issue":
			if event.Issue != nil {
				record.Set("noteable_id", event.Issue.IID)
				record.Set("noteable_title", event.Issue.Title)
				record.Set("noteable_state", event.Issue.State)
			}
		case "Commit":
			if event.Commit != nil {
				record.Set("noteable_id", event.Commit.ID)
				record.Set("noteable_title", event.Commit.Message)
				record.Set("commit_id", event.Commit.ID)
			}
		case "Snippet":
			if event.Snippet != nil {
				record.Set("noteable_id", event.Snippet.ID)
				record.Set("noteable_title", event.Snippet.Title)
			}
		}

		// 如果有代码行相关的评论信息
		if event.ObjectAttributes.LineCode != "" {
			record.Set("line_code", event.ObjectAttributes.LineCode)
		}
		if event.ObjectAttributes.CommitID != "" {
			record.Set("commit_id", event.ObjectAttributes.CommitID)
		}

		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save note event record", "error", err)
		} else {
			app.Logger().Info("Note event record saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Note event processed",
		"event": map[string]interface{}{
			"action":        event.ObjectAttributes.Action,
			"note_id":       event.ObjectAttributes.ID,
			"noteable_type": event.ObjectAttributes.NoteableType,
			"project":       event.Project.Name,
			"author_id":     event.ObjectAttributes.AuthorID,
		},
	})
}

// GitLabNoteEvent GitLab Note Hook事件数据结构
type GitLabNoteEvent struct {
	ObjectKind       string         `json:"object_kind"`
	EventType        string         `json:"event_type"`
	User             User           `json:"user"`
	ProjectID        int            `json:"project_id"`
	Project          Project        `json:"project"`
	Repository       Repository     `json:"repository"`
	ObjectAttributes NoteAttributes `json:"object_attributes"`
	MergeRequest     *MergeRequest  `json:"merge_request,omitempty"` // 当是 MR 评论时
	Issue            *Issue         `json:"issue,omitempty"`         // 当是 Issue 评论时
	Commit           *Commit        `json:"commit,omitempty"`        // 当是 Commit 评论时
	Snippet          *Snippet       `json:"snippet,omitempty"`       // 当是 Snippet 评论时
}

// NoteAttributes Note 评论属性
type NoteAttributes struct {
	ID           int         `json:"id"`
	Note         string      `json:"note"`
	NoteableType string      `json:"noteable_type"`
	AuthorID     int         `json:"author_id"`
	CreatedAt    string      `json:"created_at"` // Note Hook 使用字符串格式
	UpdatedAt    string      `json:"updated_at"` // Note Hook 使用字符串格式
	ProjectID    int         `json:"project_id"`
	Attachment   interface{} `json:"attachment"`
	LineCode     string      `json:"line_code"`
	CommitID     string      `json:"commit_id"`
	NoteableID   interface{} `json:"noteable_id"` // 可能是 int 或 null
	System       bool        `json:"system"`
	StDiff       *StDiff     `json:"st_diff"`
	Action       string      `json:"action"`
	URL          string      `json:"url"`
}

// StDiff 代码差异信息
type StDiff struct {
	Diff        string `json:"diff"`
	NewPath     string `json:"new_path"`
	OldPath     string `json:"old_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
}

// MergeRequest MR信息
type MergeRequest struct {
	ID                  int     `json:"id"`
	TargetBranch        string  `json:"target_branch"`
	SourceBranch        string  `json:"source_branch"`
	SourceProjectID     int     `json:"source_project_id"`
	AuthorID            int     `json:"author_id"`
	AssigneeID          int     `json:"assignee_id"`
	Title               string  `json:"title"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
	MilestoneID         int     `json:"milestone_id"`
	State               string  `json:"state"`
	MergeStatus         string  `json:"merge_status"`
	TargetProjectID     int     `json:"target_project_id"`
	IID                 int     `json:"iid"`
	Description         string  `json:"description"`
	Position            int     `json:"position"`
	Labels              []Label `json:"labels"`
	Source              Project `json:"source"`
	Target              Project `json:"target"`
	LastCommit          Commit  `json:"last_commit"`
	WorkInProgress      bool    `json:"work_in_progress"`
	Draft               bool    `json:"draft"`
	Assignee            User    `json:"assignee"`
	DetailedMergeStatus string  `json:"detailed_merge_status"`
}

// Issue Issue信息
type Issue struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	AssigneeIDs []int   `json:"assignee_ids"`
	AssigneeID  int     `json:"assignee_id"`
	AuthorID    int     `json:"author_id"`
	ProjectID   int     `json:"project_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Position    int     `json:"position"`
	BranchName  string  `json:"branch_name"`
	Description string  `json:"description"`
	MilestoneID int     `json:"milestone_id"`
	State       string  `json:"state"`
	IID         int     `json:"iid"`
	Labels      []Label `json:"labels"`
}

// Snippet Snippet信息
type Snippet struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Content         string `json:"content"`
	AuthorID        int    `json:"author_id"`
	ProjectID       int    `json:"project_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	FileName        string `json:"file_name"`
	Type            string `json:"type"`
	VisibilityLevel int    `json:"visibility_level"`
	URL             string `json:"url"`
}
