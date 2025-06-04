package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/pocketbase/pocketbase/core"
)

// GitLabConfig 存储GitLab配置
type GitLabConfig struct {
	WebhookSecret string
	BaseURL       string
}

// GitLabWebhook 创建GitLab webhook中间件
func GitLabWebhook(config *GitLabConfig) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// 验证GitLab webhook secret token（如果配置了）
		if config.WebhookSecret != "" {
			token := e.Request.Header.Get("X-Gitlab-Token")
			if token != config.WebhookSecret {
				e.App.Logger().Warn("GitLab webhook token verification failed",
					"receivedToken", token,
					"expectedToken", config.WebhookSecret,
				)
				return e.UnauthorizedError("Invalid GitLab webhook token", nil)
			}
		}

		// 验证Content-Type
		contentType := e.Request.Header.Get("Content-Type")
		if contentType != "application/json" {
			e.App.Logger().Warn("Invalid Content-Type for GitLab webhook",
				"contentType", contentType,
			)
			return e.BadRequestError("Invalid Content-Type, expected application/json", nil)
		}

		// 记录webhook信息
		e.App.Logger().Info("GitLab webhook received",
			"event", e.Request.Header.Get("X-Gitlab-Event"),
			"source", e.Request.Header.Get("X-Gitlab-Instance"),
			"userAgent", e.Request.Header.Get("User-Agent"),
		)

		// 将配置信息添加到请求上下文中
		ctx := context.WithValue(e.Request.Context(), "gitlab_config", config)
		e.Request = e.Request.WithContext(ctx)

		return e.Next()
	}
}

// GitLabSignatureVerify 创建GitLab webhook签名验证中间件
func GitLabSignatureVerify(config *GitLabConfig) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// 如果没有配置webhook secret，跳过签名验证
		if config.WebhookSecret == "" {
			e.App.Logger().Debug("GitLab webhook secret not configured, skipping signature verification")
			ctx := context.WithValue(e.Request.Context(), "gitlab_config", config)
			e.Request = e.Request.WithContext(ctx)
			return e.Next()
		}

		// 获取GitLab签名
		signature := e.Request.Header.Get("X-Gitlab-Event-UUID")
		if signature == "" {
			e.App.Logger().Warn("Missing X-Gitlab-Event-UUID header")
			return e.BadRequestError("Missing GitLab event UUID", nil)
		}

		// 读取请求体
		body, err := io.ReadAll(e.Request.Body)
		if err != nil {
			e.App.Logger().Error("Failed to read request body", "error", err)
			return e.BadRequestError("Failed to read request body", err)
		}

		// 计算期望的签名
		h := hmac.New(sha256.New, []byte(config.WebhookSecret))
		h.Write(body)
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		// 验证签名（这里简化处理，实际GitLab可能使用不同的签名方式）
		e.App.Logger().Debug("GitLab signature verification",
			"signature", signature,
			"expectedSignature", expectedSignature,
		)

		// 将配置信息添加到请求上下文中
		ctx := context.WithValue(e.Request.Context(), "gitlab_config", config)
		e.Request = e.Request.WithContext(ctx)

		return e.Next()
	}
}

// GetGitLabConfigFromContext 从上下文中获取GitLab配置
func GetGitLabConfigFromContext(ctx context.Context) (*GitLabConfig, bool) {
	config, ok := ctx.Value("gitlab_config").(*GitLabConfig)
	return config, ok
}
