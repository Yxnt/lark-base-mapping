package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"gitlab.yogorobot.com/sre/lark-base-mapping/middlewares"
	_ "gitlab.yogorobot.com/sre/lark-base-mapping/migrations"
	"gitlab.yogorobot.com/sre/lark-base-mapping/router"
)

func main() {
	app := pocketbase.New()

	// 加载配置
	config := LoadConfig()
	log.Printf("Loaded Lark config: AppID=%s, BaseURL=%s", config.LarkID, config.LarkBaseURL)

	// 创建飞书中间件配置
	larkConfig := &middlewares.LarkConfig{
		AppID:     config.LarkID,
		AppSecret: config.LarkSecret,
		BaseURL:   config.LarkBaseURL,
	}

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 注册路由并绑定飞书中间件
		se.Router.GET("/base/{baseID}/{tableID}/{recordID}", router.LarkBaseTable).BindFunc(
			middlewares.LarkAuth(larkConfig),
		)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
