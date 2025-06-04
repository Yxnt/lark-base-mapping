package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 创建 gitlab_note_events 集合
		collection := core.NewBaseCollection("gitlab_note_events")

		// 配置集合基本信息
		collection.Name = "gitlab_note_events"
		collection.Type = core.CollectionTypeBase
		collection.System = false

		// 添加字段
		collection.Fields.Add(&core.NumberField{
			Name:     "note_id",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "note_content",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "noteable_type",
			Required: true,
		})

		collection.Fields.Add(&core.NumberField{
			Name:     "author_id",
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
			Name:     "action",
			Required: false,
		})

		collection.Fields.Add(&core.BoolField{
			Name:     "system",
			Required: false,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "created_at",
			Required: false,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "updated_at",
			Required: false,
		})

		collection.Fields.Add(&core.URLField{
			Name:     "url",
			Required: false,
		})

		// 评论目标相关字段
		collection.Fields.Add(&core.TextField{
			Name:     "noteable_id",
			Required: false,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "noteable_title",
			Required: false,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "noteable_state",
			Required: false,
		})

		// 代码行相关字段
		collection.Fields.Add(&core.TextField{
			Name:     "line_code",
			Required: false,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "commit_id",
			Required: false,
		})

		// 原始事件数据
		collection.Fields.Add(&core.JSONField{
			Name:     "event_data",
			Required: false,
		})

		// 添加索引
		collection.Indexes = []string{
			"CREATE INDEX idx_gitlab_note_id ON gitlab_note_events (note_id)",
			"CREATE INDEX idx_gitlab_note_project_id ON gitlab_note_events (project_id)",
			"CREATE INDEX idx_gitlab_note_noteable_type ON gitlab_note_events (noteable_type)",
			"CREATE INDEX idx_gitlab_note_author_id ON gitlab_note_events (author_id)",
			"CREATE INDEX idx_gitlab_note_action ON gitlab_note_events (action)",
			"CREATE INDEX idx_gitlab_note_system ON gitlab_note_events (system)",
		}

		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚操作：删除 gitlab_note_events 集合
		collection, err := app.FindCollectionByNameOrId("gitlab_note_events")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
