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

		collection.Fields.Add(&core.TextField{
			Name:     "description",
			Required: false,
		})

		return app.Save(collection)
	}, nil)
}
