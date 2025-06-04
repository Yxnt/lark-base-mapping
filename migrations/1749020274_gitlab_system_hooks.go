package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 创建项目系统事件表
		if err := createProjectSystemEventsCollection(app); err != nil {
			return err
		}

		// 创建用户系统事件表
		if err := createUserSystemEventsCollection(app); err != nil {
			return err
		}

		// 创建组系统事件表
		if err := createGroupSystemEventsCollection(app); err != nil {
			return err
		}

		// 创建访问请求事件表
		if err := createAccessRequestEventsCollection(app); err != nil {
			return err
		}

		// 创建密钥事件表
		if err := createKeyEventsCollection(app); err != nil {
			return err
		}

		// 创建仓库更新事件表
		if err := createRepositoryUpdateEventsCollection(app); err != nil {
			return err
		}

		// 创建成员审批事件表
		if err := createMemberApprovalEventsCollection(app); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		// 回滚操作：删除所有相关集合
		collections := []string{
			"gitlab_project_system_events",
			"gitlab_user_system_events",
			"gitlab_group_system_events",
			"gitlab_access_request_events",
			"gitlab_key_events",
			"gitlab_repository_update_events",
			"gitlab_member_approval_events",
		}

		for _, name := range collections {
			collection, err := app.FindCollectionByNameOrId(name)
			if err == nil {
				if err := app.Delete(collection); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// createProjectSystemEventsCollection 创建项目系统事件集合
func createProjectSystemEventsCollection(app core.App) error {
	collection := core.NewBaseCollection("gitlab_project_system_events")
	collection.Name = "gitlab_project_system_events"
	collection.Type = core.CollectionTypeBase
	collection.System = false

	// 添加字段
	collection.Fields.Add(&core.TextField{
		Name:     "event_name",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "project_id",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "project_name",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "path",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "path_with_namespace",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "project_visibility",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "owner_name",
		Required: false,
	})

	collection.Fields.Add(&core.EmailField{
		Name:     "owner_email",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "old_path_with_namespace",
		Required: false,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "event_data",
		Required: false,
	})

	// 添加索引
	collection.Indexes = []string{
		"CREATE INDEX idx_project_system_event_name ON gitlab_project_system_events (event_name)",
		"CREATE INDEX idx_project_system_project_id ON gitlab_project_system_events (project_id)",
	}

	return app.Save(collection)
}

// createUserSystemEventsCollection 创建用户系统事件集合
func createUserSystemEventsCollection(app core.App) error {
	collection := core.NewBaseCollection("gitlab_user_system_events")
	collection.Name = "gitlab_user_system_events"
	collection.Type = core.CollectionTypeBase
	collection.System = false

	// 添加字段
	collection.Fields.Add(&core.TextField{
		Name:     "event_name",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "user_id",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "user_name",
		Required: true,
	})

	collection.Fields.Add(&core.EmailField{
		Name:     "user_email",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "user_username",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "old_username",
		Required: false,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "event_data",
		Required: false,
	})

	// 添加索引
	collection.Indexes = []string{
		"CREATE INDEX idx_user_system_event_name ON gitlab_user_system_events (event_name)",
		"CREATE INDEX idx_user_system_user_id ON gitlab_user_system_events (user_id)",
	}

	return app.Save(collection)
}

// createGroupSystemEventsCollection 创建组系统事件集合
func createGroupSystemEventsCollection(app core.App) error {
	collection := core.NewBaseCollection("gitlab_group_system_events")
	collection.Name = "gitlab_group_system_events"
	collection.Type = core.CollectionTypeBase
	collection.System = false

	// 添加字段
	collection.Fields.Add(&core.TextField{
		Name:     "event_name",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "group_id",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "group_name",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "path",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "path_with_namespace",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "old_path",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "old_path_with_namespace",
		Required: false,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "event_data",
		Required: false,
	})

	// 添加索引
	collection.Indexes = []string{
		"CREATE INDEX idx_group_system_event_name ON gitlab_group_system_events (event_name)",
		"CREATE INDEX idx_group_system_group_id ON gitlab_group_system_events (group_id)",
	}

	return app.Save(collection)
}

// createAccessRequestEventsCollection 创建访问请求事件集合
func createAccessRequestEventsCollection(app core.App) error {
	collection := core.NewBaseCollection("gitlab_access_request_events")
	collection.Name = "gitlab_access_request_events"
	collection.Type = core.CollectionTypeBase
	collection.System = false

	// 添加字段
	collection.Fields.Add(&core.TextField{
		Name:     "event_name",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "user_id",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "user_name",
		Required: true,
	})

	collection.Fields.Add(&core.EmailField{
		Name:     "user_email",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "user_username",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "group_id",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "group_name",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "group_path",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "group_access",
		Required: false,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "project_id",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "project_name",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "project_path",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "project_access",
		Required: false,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "event_data",
		Required: false,
	})

	// 添加索引
	collection.Indexes = []string{
		"CREATE INDEX idx_access_request_event_name ON gitlab_access_request_events (event_name)",
		"CREATE INDEX idx_access_request_user_id ON gitlab_access_request_events (user_id)",
	}

	return app.Save(collection)
}

// createKeyEventsCollection 创建密钥事件集合
func createKeyEventsCollection(app core.App) error {
	collection := core.NewBaseCollection("gitlab_key_events")
	collection.Name = "gitlab_key_events"
	collection.Type = core.CollectionTypeBase
	collection.System = false

	// 添加字段
	collection.Fields.Add(&core.TextField{
		Name:     "event_name",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "user_id",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "user_name",
		Required: true,
	})

	collection.Fields.Add(&core.EmailField{
		Name:     "user_email",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "key_id",
		Required: true,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "event_data",
		Required: false,
	})

	// 添加索引
	collection.Indexes = []string{
		"CREATE INDEX idx_key_event_name ON gitlab_key_events (event_name)",
		"CREATE INDEX idx_key_user_id ON gitlab_key_events (user_id)",
	}

	return app.Save(collection)
}

// createRepositoryUpdateEventsCollection 创建仓库更新事件集合
func createRepositoryUpdateEventsCollection(app core.App) error {
	collection := core.NewBaseCollection("gitlab_repository_update_events")
	collection.Name = "gitlab_repository_update_events"
	collection.Type = core.CollectionTypeBase
	collection.System = false

	// 添加字段
	collection.Fields.Add(&core.TextField{
		Name:     "event_name",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "user_id",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "user_name",
		Required: true,
	})

	collection.Fields.Add(&core.EmailField{
		Name:     "user_email",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "project_id",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "project_name",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "project_path",
		Required: true,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "refs",
		Required: false,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "changes",
		Required: false,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "event_data",
		Required: false,
	})

	// 添加索引
	collection.Indexes = []string{
		"CREATE INDEX idx_repo_update_event_name ON gitlab_repository_update_events (event_name)",
		"CREATE INDEX idx_repo_update_project_id ON gitlab_repository_update_events (project_id)",
	}

	return app.Save(collection)
}

// createMemberApprovalEventsCollection 创建成员审批事件集合
func createMemberApprovalEventsCollection(app core.App) error {
	collection := core.NewBaseCollection("gitlab_member_approval_events")
	collection.Name = "gitlab_member_approval_events"
	collection.Type = core.CollectionTypeBase
	collection.System = false

	// 添加字段
	collection.Fields.Add(&core.TextField{
		Name:     "object_kind",
		Required: true,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "action",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "user_id",
		Required: true,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "requested_by_user_id",
		Required: false,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "reviewed_by_user_id",
		Required: false,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "promotion_namespace_id",
		Required: false,
	})

	collection.Fields.Add(&core.TextField{
		Name:     "status",
		Required: false,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "new_access_level",
		Required: false,
	})

	collection.Fields.Add(&core.NumberField{
		Name:     "old_access_level",
		Required: false,
	})

	collection.Fields.Add(&core.JSONField{
		Name:     "event_data",
		Required: false,
	})

	// 添加索引
	collection.Indexes = []string{
		"CREATE INDEX idx_member_approval_object_kind ON gitlab_member_approval_events (object_kind)",
		"CREATE INDEX idx_member_approval_action ON gitlab_member_approval_events (action)",
		"CREATE INDEX idx_member_approval_user_id ON gitlab_member_approval_events (user_id)",
	}

	return app.Save(collection)
}
