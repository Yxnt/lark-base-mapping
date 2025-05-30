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

		// 添加 view_id 字段
		collection.Fields.Add(&core.TextField{
			Name:     "view_id",
			Required: true,
		})

		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚：删除 view_id 字段
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		// 删除 view_id 字段
		collection.Fields.RemoveById("view_id")

		return app.Save(collection)
	})
}
