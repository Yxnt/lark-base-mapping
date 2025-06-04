package types

import (
	"strings"
	"time"
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

// 基础用户结构
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// 基础项目结构
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

// 基础仓库结构
type Repository struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

// 基础提交结构
type Commit struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	Timestamp FlexibleTime `json:"timestamp"`
	URL       string       `json:"url"`
	Author    Author       `json:"author"`
}

// 基础作者结构
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// 基础标签结构
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

// 基础变更结构
type Change struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

// 标签变更结构
type LabelChange struct {
	Previous []Label `json:"previous"`
	Current  []Label `json:"current"`
}

// 变更集合结构
type Changes struct {
	Title       Change      `json:"title"`
	Description Change      `json:"description"`
	Labels      LabelChange `json:"labels"`
	State       Change      `json:"state"`
	UpdatedAt   Change      `json:"updated_at"`
}

// GitLab Merge Request事件数据结构
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

// MR属性结构
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

// System Hook Merge Request事件数据结构
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

// System Hook MR属性结构
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

// System Hook变更结构
type SystemHookChanges struct {
	Title       SystemHookChange      `json:"title"`
	Description SystemHookChange      `json:"description"`
	Labels      SystemHookLabelChange `json:"labels"`
	State       SystemHookChange      `json:"state"`
	UpdatedAt   SystemHookChange      `json:"updated_at"`
}

// System Hook单个变更结构
type SystemHookChange struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

// System Hook标签变更结构
type SystemHookLabelChange struct {
	Previous []Label `json:"previous"`
	Current  []Label `json:"current"`
}

// System Hook基础事件结构
type SystemHookEvent struct {
	EventName  string `json:"event_name"`
	ObjectKind string `json:"object_kind,omitempty"` // 新格式事件使用
	Action     string `json:"action,omitempty"`      // 新格式事件使用
}

// 项目系统钩子事件
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

// 用户系统钩子事件
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

// 组系统钩子事件
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

// 仓库更新事件
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

// 引用变更结构
type RefChange struct {
	Before string `json:"before"`
	After  string `json:"after"`
	Ref    string `json:"ref"`
}

// 访问请求事件
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

// 密钥事件
type KeyEvent struct {
	EventName string       `json:"event_name"`
	CreatedAt FlexibleTime `json:"created_at"`
	UpdatedAt FlexibleTime `json:"updated_at"`
	UserName  string       `json:"user_name"`
	UserEmail string       `json:"user_email"`
	UserID    int          `json:"user_id"`
	KeyID     int          `json:"key_id"`
}

// 成员审批事件 (新格式)
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

// 成员审批属性
type MemberApprovalAttributes struct {
	NewAccessLevel                       int    `json:"new_access_level,omitempty"`
	OldAccessLevel                       int    `json:"old_access_level,omitempty"`
	ExistingMemberID                     int    `json:"existing_member_id,omitempty"`
	PromotionRequestIDsThatFailedToApply []int  `json:"promotion_request_ids_that_failed_to_apply,omitempty"`
	Status                               string `json:"status,omitempty"`
}

// 所有者结构
type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GitLab Note Hook事件数据结构
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

// Note 评论属性
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

// 代码差异信息
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

// MR信息
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

// Issue信息
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

// Snippet信息
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
