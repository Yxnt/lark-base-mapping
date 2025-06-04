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
	log.Printf("Loaded Lark config: AppID=%s, BaseURL=%s, WebURL=%s", config.LarkID, config.LarkBaseURL, config.LarkWebURL)

	// 加载GitLab配置
	gitlabConfig := LoadGitLabConfig()
	log.Printf("Loaded GitLab config: BaseURL=%s, WebhookSecret configured=%t",
		gitlabConfig.BaseURL, gitlabConfig.WebhookSecret != "")

	// 创建飞书中间件配置，使用NewLarkConfig函数
	larkConfig := middlewares.NewLarkConfig(
		config.LarkID,
		config.LarkSecret,
		config.LarkBaseURL,
		config.LarkWebURL,
	)

	// 创建GitLab中间件配置
	gitlabMiddlewareConfig := &middlewares.GitLabConfig{
		WebhookSecret: gitlabConfig.WebhookSecret,
		BaseURL:       gitlabConfig.BaseURL,
	}

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 注册飞书路由并绑定飞书中间件
		se.Router.GET("/base/{baseID}/{tableID}/{recordID}", router.LarkBaseTable).BindFunc(
			middlewares.LarkAuth(larkConfig),
		)
		se.Router.GET("/base/{baseID}/{tableID}", router.LarkBaseTable).BindFunc(
			middlewares.LarkAuth(larkConfig),
		)

		// 注册GitLab webhook路由并绑定GitLab和飞书中间件
		se.Router.POST("/webhook/gitlab", router.GitLabWebhook).BindFunc(
			middlewares.GitLabWebhook(gitlabMiddlewareConfig),
			middlewares.LarkAuth(larkConfig),
		)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
