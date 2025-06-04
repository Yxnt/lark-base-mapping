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

		// 分页处理所有现有的 lark_table 记录
		limit := 100
		offset := 0

		for {
			records, err := app.FindRecordsByFilter(
				collection.Id,
				"record_id_field = '' OR record_id_field IS NULL",
				"",
				limit,
				offset,
			)
			if err != nil {
				return err
			}

			// 如果没有更多记录，退出循环
			if len(records) == 0 {
				break
			}

			// 为每个记录设置默认的 record_id_field 值
			for _, record := range records {
				record.Set("record_id_field", "编号")
				if err := app.Save(record); err != nil {
					return err
				}
			}

			// 如果返回的记录数小于限制数，说明已经处理完所有记录
			if len(records) < limit {
				break
			}

			// 更新偏移量以获取下一批记录
			offset += limit
		}

		return nil
	}, func(app core.App) error {
		// 获取 lark_table collection
		collection, err := app.FindCollectionByNameOrId("lark_table")
		if err != nil {
			return err
		}

		// 回滚操作：分页处理所有 record_id_field 为 '编号' 的记录
		limit := 100
		offset := 0

		for {
			records, err := app.FindRecordsByFilter(
				collection.Id,
				"record_id_field = '编号'",
				"",
				limit,
				offset,
			)
			if err != nil {
				return err
			}

			// 如果没有更多记录，退出循环
			if len(records) == 0 {
				break
			}

			// 将所有记录的 record_id_field 设为空
			for _, record := range records {
				record.Set("record_id_field", "")
				if err := app.Save(record); err != nil {
					return err
				}
			}

			// 如果返回的记录数小于限制数，说明已经处理完所有记录
			if len(records) < limit {
				break
			}

			// 更新偏移量以获取下一批记录
			offset += limit
		}

		return nil
	})
}
