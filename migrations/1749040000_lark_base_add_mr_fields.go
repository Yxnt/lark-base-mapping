package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("lark_base")
		if err != nil {
			return err
		}

		// 获取 lark_base 集合，用于设置关联
		baseCollection, err := app.FindCollectionByNameOrId("lark_base")
		if err != nil {
			return err
		}

		// 添加 merge request 表的 table_id 字段
		collection.Fields.Add(&core.TextField{
			Name:     "mr_table_id",
			Required: false,
		})

		// 添加 merge request 表的 view_id 字段
		collection.Fields.Add(&core.TextField{
			Name:     "mr_view_id",
			Required: false,
		})

		// 添加 merge request 表的 base_id 字段，外键关联到 lark_base 表
		collection.Fields.Add(&core.RelationField{
			Name:         "mr_base_id",
			Required:     false,
			CollectionId: baseCollection.Id,
			MaxSelect:    1, // 只能关联一个 base
		})

		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚操作：删除添加的字段
		collection, err := app.FindCollectionByNameOrId("lark_base")
		if err != nil {
			return err
		}

		// 删除 mr_table_id 字段
		mrTableIdField := collection.Fields.GetByName("mr_table_id")
		if mrTableIdField != nil {
			collection.Fields.RemoveById(mrTableIdField.GetId())
		}

		// 删除 mr_view_id 字段
		mrViewIdField := collection.Fields.GetByName("mr_view_id")
		if mrViewIdField != nil {
			collection.Fields.RemoveById(mrViewIdField.GetId())
		}

		// 删除 mr_base_id 字段
		mrBaseIdField := collection.Fields.GetByName("mr_base_id")
		if mrBaseIdField != nil {
			collection.Fields.RemoveById(mrBaseIdField.GetId())
		}

		return app.Save(collection)
	})
}
