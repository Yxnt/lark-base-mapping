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

		// 获取 record_id_field 字段并设置为必需字段
		field := collection.Fields.GetByName("record_id_field")
		if field != nil {
			if textField, ok := field.(*core.TextField); ok {
				textField.Required = true
			}
		}

		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚：将 record_id_field 字段设置为非必需
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		field := collection.Fields.GetByName("record_id_field")
		if field != nil {
			if textField, ok := field.(*core.TextField); ok {
				textField.Required = false
			}
		}

		return app.Save(collection)
	})
}
