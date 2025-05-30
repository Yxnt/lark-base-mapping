package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 获取 lark_table 集合
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		// 获取 lark_base 集合，用于设置关联
		baseCollection, err := app.FindCollectionByNameOrId("lark_base")
		if err != nil {
			return err
		}

		// 添加 base_id 关联字段，指向 lark_base 表
		collection.Fields.Add(&core.RelationField{
			Name:         "base_id",
			Required:     true,
			CollectionId: baseCollection.Id,
			MaxSelect:    1, // 每个 table 只能关联一个 base
		})

		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚操作：移除 base_id 字段
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		// 移除 base_id 字段
		collection.Fields.RemoveById(collection.Fields.GetByName("base_id").GetId())

		return app.Save(collection)
	})
}
