package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection := core.NewBaseCollection("lark_base")

		collection.Fields = core.NewFieldsList(
			&core.TextField{
				Name:     "base_id",
				Required: true,
			},
		)

		return app.Save(collection)
	}, nil)
}
