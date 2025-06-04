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

		// 添加 record_id_field 字段，用于存储动态的记录ID字段名
		collection.Fields.Add(&core.TextField{
			Name:     "record_id_field",
			Required: false, // 暂时设为非必需，避免现有记录出错
		})

		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚：删除 record_id_field 字段
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		// 删除 record_id_field 字段
		field := collection.Fields.GetByName("record_id_field")
		if field != nil {
			collection.Fields.RemoveById(field.GetId())
		}

		return app.Save(collection)
	})
}
