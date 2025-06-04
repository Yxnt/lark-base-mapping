package router

import (
	"io"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"gitlab.yogorobot.com/sre/lark-base-mapping/handlers"
	"gitlab.yogorobot.com/sre/lark-base-mapping/middlewares"
)

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

	// 根据事件类型分发到对应的处理器
	switch eventType {
	case "System Hook":
		return handlers.HandleSystemHookEvent(e, body)
	case "Merge Request Hook":
		return handlers.HandleMergeRequestEvent(e, body)
	case "Note Hook":
		return handlers.HandleNoteEvent(e, body)
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
