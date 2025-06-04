package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 获取 lark_table collection
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		// 添加 gitlab_project_id 字段
		collection.Fields.Add(&core.TextField{
			Name:     "gitlab_project_id",
			Required: false,
		})

		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚：删除 gitlab_project_id 字段
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		// 删除 gitlab_project_id 字段
		field := collection.Fields.GetByName("gitlab_project_id")
		if field != nil {
			collection.Fields.RemoveById(field.GetId())
		}

		return app.Save(collection)
	})
}
