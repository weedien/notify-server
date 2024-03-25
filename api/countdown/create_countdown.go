package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/weedien/countdown-server/api"
	"github.com/weedien/countdown-server/api/models"
	"github.com/weedien/countdown-server/store"
	"github.com/weedien/countdown-server/util"
)

// 创建倒计时
type CreateCountdownRequest struct {
	Duration string `json:"duration" validate:"required"`
	Remark   string `json:"remark"`
	Message  string `json:"message"`
	// NotifyConfig `json:"notify_config"`
}

// 创建倒计时
func CreateCountdown(w http.ResponseWriter, r *http.Request) {
	input := parseCreateArgs(w, r)

	duration, err := time.ParseDuration(input.Duration)
	if err != nil {
		api.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// 倒计时结束，提示消息的默认值
	if input.Message == "" {
		input.Message = "Time's up!"
	}

	id := store.SnowFlake.Generate().String()
	countdown := models.Countdown{
		ID:        id,
		QueryCode: util.GenQueryCode(),
		ExpireAt:  time.Now().Add(duration),
		CreatedAt: time.Now(),
		Remark:    sql.NullString{String: input.Remark, Valid: input.Remark != ""},
		Message:   sql.NullString{String: input.Message, Valid: input.Message != ""},
	}

	_, err = store.DB.Exec("INSERT INTO countdowns (id, query_code, expire_at, created_at, remark, message) VALUES (?, ?, ?, ?, ?, ?)",
		countdown.ID, countdown.QueryCode, countdown.ExpireAt, countdown.CreatedAt, countdown.Remark, countdown.Message)
	if err != nil {
		api.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":         countdown.ID,
		"query_code": countdown.QueryCode,
		"created_at": countdown.CreatedAt.Format(time.DateTime),
		"expire_at":  countdown.ExpireAt.Format(time.DateTime),
	}

	api.SuccessResponse(w, response)
}

// 解析创建倒计时的参数
func parseCreateArgs(w http.ResponseWriter, r *http.Request) CreateCountdownRequest {
	var input CreateCountdownRequest
	if r.Header.Get("Content-Type") == "application/json" {

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			api.ErrorResponse(w, err, http.StatusBadRequest)
			return CreateCountdownRequest{}
		}
	} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			api.ErrorResponse(w, err, http.StatusBadRequest)
			return CreateCountdownRequest{}
		}
		input.Duration = r.Form.Get("duration")
		input.Remark = r.Form.Get("remark")
		input.Message = r.Form.Get("message")
	}
	// 如果input中的任何字段为空，尝试从请求路径中获取
	if input.Duration == "" {
		input.Duration = r.URL.Query().Get("duration")
	}
	if input.Remark == "" {
		input.Remark = r.URL.Query().Get("remark")
	}
	if input.Message == "" {
		input.Message = r.URL.Query().Get("message")
	}
	return input
}
