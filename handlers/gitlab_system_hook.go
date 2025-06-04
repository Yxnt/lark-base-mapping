package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"gitlab.yogorobot.com/sre/lark-base-mapping/types"
)

// HandleSystemHookEvent 处理System Hook事件
func HandleSystemHookEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	// 先解析基本的事件信息来确定事件类型
	var baseEvent types.SystemHookEvent
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
func handleNewFormatSystemEvent(e *core.RequestEvent, body []byte, baseEvent types.SystemHookEvent) error {
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
		return HandleSystemHookMergeRequestEvent(e, body)
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

	var event types.ProjectSystemHookEvent
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

	var event types.UserSystemHookEvent
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

	var event types.GroupSystemHookEvent
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

	var event types.AccessRequestEvent
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

	var event types.KeyEvent
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

	var event types.RepositoryUpdateEvent
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
func handleMemberApprovalEvent(e *core.RequestEvent, body []byte, baseEvent types.SystemHookEvent) error {
	app := e.App

	var event types.MemberApprovalEvent
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
