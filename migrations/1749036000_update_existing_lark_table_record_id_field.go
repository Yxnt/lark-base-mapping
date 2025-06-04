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

		// 获取所有现有的 lark_table 记录
		records, err := app.FindRecordsByFilter(
			collection.Id,
			"record_id_field = '' OR record_id_field IS NULL",
			"",
			100,
			0,
		)
		if err != nil {
			return err
		}

		// 为每个记录设置默认的 record_id_field 值
		for _, record := range records {
			record.Set("record_id_field", "编号")
			if err := app.Save(record); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		// 获取 lark_table collection
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		// 回滚操作：将所有记录的 record_id_field 设为空
		records, err := app.FindRecordsByFilter(
			collection.Id,
			"record_id_field = '编号'",
			"",
			100,
			0,
		)
		if err != nil {
			return err
		}

		for _, record := range records {
			record.Set("record_id_field", "")
			if err := app.Save(record); err != nil {
				return err
			}
		}

		return nil
	})
}
