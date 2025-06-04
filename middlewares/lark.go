package middlewares

import (
	"context"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/pocketbase/pocketbase/core"
)

// LarkConfig 存储飞书配置
type LarkConfig struct {
	AppID     string
	AppSecret string
	BaseURL   string
	WebURL    string
	Client    *lark.Client
}

// NewLarkConfig 创建新的飞书配置并初始化客户端
func NewLarkConfig(appID, appSecret, baseURL, webURL string) *LarkConfig {
	config := &LarkConfig{
		AppID:     appID,
		AppSecret: appSecret,
		BaseURL:   baseURL,
		WebURL:    webURL,
	}
	config.Client = lark.NewClient(appID, appSecret)
	return config
}

// LarkAuth 创建飞书认证中间件
func LarkAuth(config *LarkConfig) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// 在这里可以进行飞书相关的认证或预处理
		// 例如：验证访问令牌、检查权限等

		// 确保客户端已初始化
		if config.Client == nil {
			config.Client = lark.NewClient(config.AppID, config.AppSecret)
		}

		// 将配置信息和客户端添加到请求上下文中，供后续使用
		ctx := context.WithValue(e.Request.Context(), "lark_config", config)
		ctx = context.WithValue(ctx, "lark_client", config.Client)
		e.Request = e.Request.WithContext(ctx)

		// 记录请求信息
		e.App.Logger().Info("Lark middleware processing request",
			"baseID", e.Request.PathValue("baseID"),
			"tableID", e.Request.PathValue("tableID"),
			"recordID", e.Request.PathValue("recordID"),
		)

		// 可以在这里添加更多的飞书 API 认证逻辑
		// 例如获取访问令牌等

		return e.Next()
	}
}

// LarkAuthRequired 创建需要飞书认证的中间件
func LarkAuthRequired(config *LarkConfig) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// 检查飞书配置是否完整
		if config.AppID == "" || config.AppSecret == "" {
			return e.BadRequestError("Missing Lark configuration", nil)
		}

		// 确保客户端已初始化
		if config.Client == nil {
			config.Client = lark.NewClient(config.AppID, config.AppSecret)
		}

		// 可以在这里添加飞书 token 验证逻辑
		// 例如：从请求头中获取 token，验证其有效性等

		// 将配置和客户端添加到上下文
		ctx := context.WithValue(e.Request.Context(), "lark_config", config)
		ctx = context.WithValue(ctx, "lark_client", config.Client)
		e.Request = e.Request.WithContext(ctx)

		return e.Next()
	}
}

// GetLarkConfigFromContext 从上下文中获取飞书配置
func GetLarkConfigFromContext(ctx context.Context) (*LarkConfig, bool) {
	config, ok := ctx.Value("lark_config").(*LarkConfig)
	return config, ok
}

// GetLarkClientFromContext 从上下文中获取飞书客户端
func GetLarkClientFromContext(ctx context.Context) (*lark.Client, bool) {
	client, ok := ctx.Value("lark_client").(*lark.Client)
	return client, ok
}
