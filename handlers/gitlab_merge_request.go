package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"gitlab.yogorobot.com/sre/lark-base-mapping/types"
)

// HandleMergeRequestEvent 处理Merge Request事件
func HandleMergeRequestEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	var event types.GitLabMergeRequestEvent
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

// HandleSystemHookMergeRequestEvent 处理System Hook格式的Merge Request事件
func HandleSystemHookMergeRequestEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	var event types.SystemHookMergeRequestEvent
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
