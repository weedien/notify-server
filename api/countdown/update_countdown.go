package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/weedien/countdown-server/api"
	"github.com/weedien/countdown-server/store"
)

// 更新倒计时
type UpdateCountdownRequest struct {
	Increment string `json:"increment"`
	Remark    string `json:"remark"`
	Message   string `json:"message"`
}

// 更新倒计时
func UpdateCountdown(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	input := parseUpdateArgs(w, r)

	var expireAt time.Time
	err := store.DB.QueryRow("SELECT expire_at FROM countdowns WHERE id = ?", id).Scan(&expireAt)
	if err != nil {
		if err == sql.ErrNoRows {
			api.ErrorResponse(w, err, http.StatusNotFound)
			return
		}
		api.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// 对已经结束的倒计时进行修改没有意义
	if expireAt.Before(time.Now()) {
		api.ErrorResponse(w, fmt.Errorf("countdown %s has expired", id), http.StatusBadRequest)
		return
	}

	if input.Increment != "" {
		increment, err := time.ParseDuration(input.Increment)
		if err != nil {
			api.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// 忽略：扣减后早于当前时间，因为查询的时候会显示已经结束
		expireAt = expireAt.Add(increment)
	}

	query := `
        UPDATE countdowns
        SET expire_at = CASE WHEN ? = '' THEN expire_at ELSE ? END,
						remark = CASE WHEN ? = '' THEN remark ELSE ? END,
						message = CASE WHEN ? = '' THEN message ELSE ? END
						updated_at = NOW()
        WHERE id = ?
    `
	_, err = store.DB.Exec(query,
		input.Increment, expireAt,
		input.Remark, input.Remark,
		input.Message, input.Message,
		id,
	)

	if err != nil {
		api.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":         id,
		"updated_at": time.Now().Format(time.DateTime),
		"expire_at":  expireAt.Format(time.DateTime),
	}

	api.SuccessResponse(w, response)
}

// 解析更新倒计时的参数
func parseUpdateArgs(w http.ResponseWriter, r *http.Request) UpdateCountdownRequest {
	var input UpdateCountdownRequest
	if r.Header.Get("Content-Type") == "application/json" {

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			api.ErrorResponse(w, err, http.StatusBadRequest)
			return UpdateCountdownRequest{}
		}
	} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			api.ErrorResponse(w, err, http.StatusBadRequest)
			return UpdateCountdownRequest{}
		}
		input.Increment = r.Form.Get("increment")
		input.Remark = r.Form.Get("remark")
		input.Message = r.Form.Get("message")
	}
	// 如果input中的任何字段为空，尝试从请求路径中获取
	if input.Increment == "" {
		input.Increment = r.URL.Query().Get("increment")
	}
	if input.Remark == "" {
		input.Remark = r.URL.Query().Get("remark")
	}
	if input.Message == "" {
		input.Message = r.URL.Query().Get("message")
	}
	return input
}
