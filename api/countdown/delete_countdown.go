package api

import (
	"net/http"

	"github.com/weedien/countdown-server/api"
	"github.com/weedien/countdown-server/store"
)

// 删除倒计时
func DeleteCountdown(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := store.DB.Exec("DELETE FROM countdowns WHERE id = ?", id)
	if err != nil {
		api.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	api.SuccessResponse(w, map[string]interface{}{"id": id, "status": "success"})
}
