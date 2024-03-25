package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/weedien/countdown-server/api"
	"github.com/weedien/countdown-server/api/models"
	"github.com/weedien/countdown-server/store"
	"github.com/weedien/countdown-server/util"
)

// 获取倒计时
type GetCountdownRequest struct {
	ID        string `json:"id,omitempty"`
	QueryCode string `json:"query_code,omitempty"`
	Passed    string `json:"passed,omitempty"`
	Remaining string `json:"remaining,omitempty"`
	ExpireAt  string `json:"expire_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Remark    string `json:"remark,omitempty"`
	Message   string `json:"message,omitempty"`
}

// 获取倒计时
func GetCountdown(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var countdown models.Countdown
	err := store.DB.QueryRow("SELECT id, query_code, expire_at, created_at, updated_at, remark, message FROM countdowns WHERE id = ?", id).Scan(&countdown.ID, &countdown.QueryCode, &countdown.ExpireAt, &countdown.CreatedAt, &countdown.UpdatedAt, &countdown.Remark, &countdown.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			api.ErrorResponse(w, err, http.StatusNotFound)
			return
		}
		api.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	remaining := time.Until(countdown.ExpireAt)
	if remaining < 0 {
		remaining = time.Duration(0)
	}
	passed := time.Since(countdown.CreatedAt)

	response := GetCountdownRequest{
		ID:        countdown.ID,
		QueryCode: countdown.QueryCode,
		Passed:    util.FormatDuration(passed),
		Remaining: util.FormatDuration(remaining),
		ExpireAt:  countdown.ExpireAt.Format(time.DateTime),
		CreatedAt: countdown.CreatedAt.Format(time.DateTime),
		Remark:    countdown.Remark.String,
	}
	if countdown.UpdatedAt.Valid {
		response.UpdatedAt = countdown.UpdatedAt.Time.Format(time.DateTime)
	}
	if countdown.Message.Valid && remaining == 0 {
		response.Message = countdown.Message.String
	}

	api.SuccessResponse(w, response)
}

// 查询所有倒计时时间
func GetAllCountdown(w http.ResponseWriter, r *http.Request) {
	history, err := strconv.ParseBool(r.URL.Query().Get("history"))
	if err != nil {
		history = false
	}

	var rows *sql.Rows
	if history {
		rows, err = store.DB.Query("SELECT id, query_code, expire_at, created_at, remark, message FROM countdowns ORDER BY expire_at DESC")
	} else {
		rows, err = store.DB.Query("SELECT id, query_code, expire_at, created_at, remark, message FROM countdowns WHERE expire_at > NOW()")
	}

	if err != nil {
		api.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var countdowns []GetCountdownRequest = make([]GetCountdownRequest, 0)
	for rows.Next() {
		var countdown models.Countdown
		err = rows.Scan(&countdown.ID, &countdown.QueryCode, &countdown.ExpireAt, &countdown.CreatedAt, &countdown.Remark, &countdown.Message)
		if err != nil {
			api.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		remaining := time.Until(countdown.ExpireAt)
		if remaining < 0 {
			remaining = time.Duration(0)
		}
		passed := time.Since(countdown.CreatedAt)

		element := GetCountdownRequest{
			ID:        countdown.ID,
			QueryCode: countdown.QueryCode,
			Passed:    util.FormatDuration(passed),
			Remaining: util.FormatDuration(remaining),
			ExpireAt:  countdown.ExpireAt.Format(time.DateTime),
			CreatedAt: countdown.CreatedAt.Format(time.DateTime),
			Remark:    countdown.Remark.String,
		}
		if countdown.UpdatedAt.Valid {
			element.UpdatedAt = countdown.UpdatedAt.Time.Format(time.DateTime)
		}
		if countdown.Message.Valid && remaining == 0 {
			element.Message = countdown.Message.String
		}

		countdowns = append(countdowns, element)
	}

	api.SuccessResponse(w, countdowns)
}

// 根据QueryCode查询倒计时
func GetCountdownByQueryCode(w http.ResponseWriter, r *http.Request) {
	queryCode := r.PathValue("query_code")
	detail, err := strconv.ParseBool(r.URL.Query().Get("detail"))
	if err != nil {
		detail = false
	}

	var countdown models.Countdown
	err = store.DB.QueryRow("SELECT id, expire_at, created_at, updated_at, remark, message FROM countdowns WHERE query_code = ? and expire_at > NOW()", queryCode).Scan(&countdown.ID, &countdown.ExpireAt, &countdown.CreatedAt, &countdown.UpdatedAt, &countdown.Remark, &countdown.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			api.SuccessMsg(w, "the countdown not exists or is over")
			return
		}
		api.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	remaining := time.Until(countdown.ExpireAt)
	if remaining < 0 {
		remaining = time.Duration(0)
	}
	passed := time.Since(countdown.CreatedAt)

	var response GetCountdownRequest
	if detail {
		response = GetCountdownRequest{
			ID:        countdown.ID,
			Passed:    util.FormatDuration(passed),
			Remaining: util.FormatDuration(remaining),
			ExpireAt:  countdown.ExpireAt.Format(time.DateTime),
			CreatedAt: countdown.CreatedAt.Format(time.DateTime),
			Remark:    countdown.Remark.String,
		}
		if countdown.UpdatedAt.Valid {
			response.UpdatedAt = countdown.UpdatedAt.Time.Format(time.DateTime)
		}
	} else {
		response = GetCountdownRequest{
			Passed:    util.FormatDuration(passed),
			Remaining: util.FormatDuration(remaining),
			ExpireAt:  countdown.ExpireAt.Format(time.DateTime),
			Remark:    countdown.Remark.String,
		}
	}

	if countdown.Message.Valid && remaining == 0 {
		response.Message = countdown.Message.String
	}

	api.SuccessResponse(w, response)
}
