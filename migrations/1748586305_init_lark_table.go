package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection := core.NewBaseCollection("lark_table")

		collection.Fields = core.NewFieldsList(
			&core.TextField{
				Name:     "table_id",
				Required: true,
			},
			&core.TextField{
				Name:     "table_name",
				Required: true,
			},
		)

		return app.Save(collection)
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}
