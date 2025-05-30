package router

import (
	"context"
	"fmt"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/pocketbase/pocketbase/core"
	"gitlab.yogorobot.com/sre/lark-base-mapping/middlewares"
)

func LarkBaseTable(e *core.RequestEvent) error {
	app := e.App
	baseID := e.Request.PathValue("baseID")
	tableID := e.Request.PathValue("tableID")
	recordID := e.Request.PathValue("recordID")

	// 从中间件上下文中获取飞书配置
	larkConfig, ok := middlewares.GetLarkConfigFromContext(e.Request.Context())
	if !ok {
		return e.BadRequestError("Lark config not found in context", nil)
	}

	app.Logger().Info("Using Lark config",
		"appID", larkConfig.AppID,
		"baseURL", larkConfig.BaseURL,
	)

	// 先查询 table_id 是否存在
	table, err := app.FindFirstRecordByData("lark_table", "table_id", tableID)
	if err != nil {
		return e.NotFoundError("Table not found", err)
	}

	app.Logger().Info("Found table", "id", table.Id, "tableID", tableID)

	// 检查 table 关联的 base_id 是否与请求的 baseID 匹配
	tableBaseID := table.GetString("base_id")
	if tableBaseID == "" {
		return e.NotFoundError("Table is not associated with any base", nil)
	}

	// 查询关联的 base 记录
	base, err := app.FindRecordById("lark_base", tableBaseID)
	if err != nil {
		return e.NotFoundError("Associated base not found", err)
	}

	// 验证 base 的 base_id 是否与请求的 baseID 匹配
	if base.GetString("base_id") != baseID {
		app.Logger().Warn("Base ID mismatch",
			"requestedBaseID", baseID,
			"tableAssociatedBaseID", base.GetString("base_id"))
		return e.NotFoundError("Table is not associated with the requested base", nil)
	}

	app.Logger().Info("Verified base-table association",
		"baseID", baseID,
		"tableID", tableID,
		"baseRecordID", base.Id)

	// 如果没有提供 recordID，直接重定向到飞书页面
	if recordID == "" {
		// 获取 table 记录中的 view_id
		viewID := table.GetString("view_id")
		if viewID == "" {
			return e.NotFoundError("View ID not found for this table", nil)
		}

		// 构建重定向 URL
		redirectURL := fmt.Sprintf("%s/base/%s?table=%s&view=%s", larkConfig.WebURL, baseID, tableID, viewID)

		app.Logger().Info("Redirecting to Feishu page",
			"baseID", baseID,
			"tableID", tableID,
			"viewID", viewID,
			"redirectURL", redirectURL)

		return e.Redirect(http.StatusFound, redirectURL)
	}

	app.Logger().Info("Processing record", "recordID", recordID)

	// 初始化飞书客户端
	client := lark.NewClient(larkConfig.AppID, larkConfig.AppSecret)

	// 使用搜索记录的方式获取记录
	searchReq := larkbitable.NewSearchAppTableRecordReqBuilder().
		AppToken(baseID).
		TableId(tableID).
		PageSize(20).
		Body(larkbitable.NewSearchAppTableRecordReqBodyBuilder().
			Filter(larkbitable.NewFilterInfoBuilder().
				Conjunction(`and`).
				Conditions([]*larkbitable.Condition{
					larkbitable.NewConditionBuilder().
						FieldName(`编号`).
						Operator(`is`).
						Value([]string{recordID}).
						Build(),
				}).
				Build()).
			AutomaticFields(false).
			Build()).
		Build()

	ctx := context.Background()
	searchResp, err := client.Bitable.V1.AppTableRecord.Search(ctx, searchReq)

	if err != nil {
		app.Logger().Error("Failed to search record", "error", err)
		return e.InternalServerError("Failed to search record", err)
	}

	if !searchResp.Success() {
		app.Logger().Error("Lark API search request failed", "code", searchResp.Code, "msg", searchResp.Msg, "requestId", searchResp.RequestId())
		return e.NotFoundError("Record not found", fmt.Errorf("code: %d, msg: %s", searchResp.Code, searchResp.Msg))
	}

	// 检查是否找到记录
	if searchResp.Data == nil || len(searchResp.Data.Items) == 0 {
		app.Logger().Warn("No record found with the given ID", "recordID", recordID)
		return e.NotFoundError("Record not found", nil)
	}

	// 获取第一个匹配的记录
	record := searchResp.Data.Items[0]

	app.Logger().Info("Successfully retrieved record", "recordID", *record.RecordId)

	// 使用BatchGet方法获取记录的详细信息，包括shared_url
	batchGetReq := larkbitable.NewBatchGetAppTableRecordReqBuilder().
		AppToken(baseID).
		TableId(tableID).
		Body(larkbitable.NewBatchGetAppTableRecordReqBodyBuilder().
			RecordIds([]string{*record.RecordId}).
			WithSharedUrl(true).
			AutomaticFields(true).
			Build()).
		Build()

	batchGetResp, err := client.Bitable.V1.AppTableRecord.BatchGet(ctx, batchGetReq)

	if err != nil {
		app.Logger().Error("Failed to batch get record", "error", err)
		return e.InternalServerError("Failed to get record details", err)
	}

	if !batchGetResp.Success() {
		app.Logger().Error("Lark API batch get request failed", "code", batchGetResp.Code, "msg", batchGetResp.Msg, "requestId", batchGetResp.RequestId())
		return e.InternalServerError("Failed to get record details", fmt.Errorf("code: %d, msg: %s", batchGetResp.Code, batchGetResp.Msg))
	}

	// 检查是否获取到记录详情
	if batchGetResp.Data == nil || len(batchGetResp.Data.Records) == 0 {
		app.Logger().Warn("No record details found", "recordID", *record.RecordId)
		return e.NotFoundError("Record details not found", nil)
	}

	// 获取记录的shared_url
	recordDetail := batchGetResp.Data.Records[0]
	var sharedURL string
	if recordDetail.SharedUrl != nil && *recordDetail.SharedUrl != "" {
		sharedURL = *recordDetail.SharedUrl
	} else {
		app.Logger().Warn("Shared URL not found in record details", "recordID", *record.RecordId)
		return e.NotFoundError("Shared URL not available", nil)
	}

	app.Logger().Info("Retrieved shared URL", "recordID", *record.RecordId, "sharedURL", sharedURL)

	return e.Redirect(http.StatusFound, sharedURL)
}
