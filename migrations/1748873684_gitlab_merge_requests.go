package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 创建 gitlab_merge_requests 集合
		collection := core.NewBaseCollection("gitlab_merge_requests")

		// 配置集合基本信息
		collection.Name = "gitlab_merge_requests"
		collection.Type = core.CollectionTypeBase
		collection.System = false

		// 添加字段
		collection.Fields.Add(&core.NumberField{
			Name:     "mr_id",
			Required: true,
		})

		collection.Fields.Add(&core.NumberField{
			Name:     "mr_iid",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "title",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "description",
			Required: false,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "state",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "action",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "author_name",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "author_username",
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
			Name:     "source_branch",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "target_branch",
			Required: true,
		})

		collection.Fields.Add(&core.URLField{
			Name:     "url",
			Required: false,
		})

		collection.Fields.Add(&core.JSONField{
			Name:     "event_data",
			Required: false,
		})

		// 添加索引
		collection.Indexes = []string{
			"CREATE INDEX idx_gitlab_mr_id ON gitlab_merge_requests (mr_id)",
			"CREATE INDEX idx_gitlab_project_id ON gitlab_merge_requests (project_id)",
			"CREATE INDEX idx_gitlab_state ON gitlab_merge_requests (state)",
			"CREATE INDEX idx_gitlab_action ON gitlab_merge_requests (action)",
		}

		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚操作：删除 gitlab_merge_requests 集合
		collection, err := app.FindCollectionByNameOrId("gitlab_merge_requests")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
