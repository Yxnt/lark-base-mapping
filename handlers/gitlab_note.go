package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"gitlab.yogorobot.com/sre/lark-base-mapping/types"
)

// HandleNoteEvent 处理Note事件（评论事件）
func HandleNoteEvent(e *core.RequestEvent, body []byte) error {
	app := e.App

	var event types.GitLabNoteEvent
	if err := json.Unmarshal(body, &event); err != nil {
		app.Logger().Error("Failed to parse note event", "error", err)
		return e.BadRequestError("Invalid note event format", err)
	}

	app.Logger().Info("Processing note event",
		"action", event.ObjectAttributes.Action,
		"noteID", event.ObjectAttributes.ID,
		"noteableType", event.ObjectAttributes.NoteableType,
		"authorID", event.ObjectAttributes.AuthorID,
		"projectName", event.Project.Name,
		"noteContent", func() string {
			if len(event.ObjectAttributes.Note) > 50 {
				return event.ObjectAttributes.Note[:50] + "..."
			}
			return event.ObjectAttributes.Note
		}(),
	)

	// 根据评论类型添加额外的日志信息
	switch event.ObjectAttributes.NoteableType {
	case "MergeRequest":
		if event.MergeRequest != nil {
			app.Logger().Info("Note on merge request",
				"mrID", event.MergeRequest.IID,
				"mrTitle", event.MergeRequest.Title,
				"mrState", event.MergeRequest.State,
			)
		}
	case "Issue":
		if event.Issue != nil {
			app.Logger().Info("Note on issue",
				"issueID", event.Issue.IID,
				"issueTitle", event.Issue.Title,
				"issueState", event.Issue.State,
			)
		}
	case "Commit":
		if event.Commit != nil {
			app.Logger().Info("Note on commit",
				"commitID", event.Commit.ID,
				"commitMessage", func() string {
					if len(event.Commit.Message) > 50 {
						return event.Commit.Message[:50] + "..."
					}
					return event.Commit.Message
				}(),
			)
		}
	case "Snippet":
		if event.Snippet != nil {
			app.Logger().Info("Note on snippet",
				"snippetID", event.Snippet.ID,
				"snippetTitle", event.Snippet.Title,
			)
		}
	}

	// 保存Note事件到数据库
	collection, err := app.FindCollectionByNameOrId("gitlab_note_events")
	if err != nil {
		app.Logger().Warn("gitlab_note_events collection not found", "error", err)
	} else {
		record := core.NewRecord(collection)
		record.Set("note_id", event.ObjectAttributes.ID)
		record.Set("note_content", event.ObjectAttributes.Note)
		record.Set("noteable_type", event.ObjectAttributes.NoteableType)
		record.Set("author_id", event.ObjectAttributes.AuthorID)
		record.Set("project_id", event.Project.ID)
		record.Set("project_name", event.Project.Name)
		record.Set("action", event.ObjectAttributes.Action)
		record.Set("system", event.ObjectAttributes.System)
		record.Set("created_at", event.ObjectAttributes.CreatedAt)
		record.Set("updated_at", event.ObjectAttributes.UpdatedAt)
		record.Set("url", event.ObjectAttributes.URL)

		// 根据评论类型设置相关的ID和信息
		switch event.ObjectAttributes.NoteableType {
		case "MergeRequest":
			if event.MergeRequest != nil {
				record.Set("noteable_id", event.MergeRequest.IID)
				record.Set("noteable_title", event.MergeRequest.Title)
				record.Set("noteable_state", event.MergeRequest.State)
			}
		case "Issue":
			if event.Issue != nil {
				record.Set("noteable_id", event.Issue.IID)
				record.Set("noteable_title", event.Issue.Title)
				record.Set("noteable_state", event.Issue.State)
			}
		case "Commit":
			if event.Commit != nil {
				record.Set("noteable_id", event.Commit.ID)
				record.Set("noteable_title", event.Commit.Message)
				record.Set("commit_id", event.Commit.ID)
			}
		case "Snippet":
			if event.Snippet != nil {
				record.Set("noteable_id", event.Snippet.ID)
				record.Set("noteable_title", event.Snippet.Title)
			}
		}

		// 如果有代码行相关的评论信息
		if event.ObjectAttributes.LineCode != "" {
			record.Set("line_code", event.ObjectAttributes.LineCode)
		}
		if event.ObjectAttributes.CommitID != "" {
			record.Set("commit_id", event.ObjectAttributes.CommitID)
		}

		record.Set("event_data", string(body))

		if err := app.Save(record); err != nil {
			app.Logger().Error("Failed to save note event record", "error", err)
		} else {
			app.Logger().Info("Note event record saved", "recordID", record.Id)
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Note event processed",
		"event": map[string]interface{}{
			"action":        event.ObjectAttributes.Action,
			"note_id":       event.ObjectAttributes.ID,
			"noteable_type": event.ObjectAttributes.NoteableType,
			"project":       event.Project.Name,
			"author_id":     event.ObjectAttributes.AuthorID,
		},
	})
}
