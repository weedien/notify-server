package api

import (
	"encoding/json"
	"net/http"
)

// 错误响应
func ErrorResponse(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  code,
		"error": err.Error(),
	})
}

// 成功响应
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusOK,
		"data": data,
	})
}
func SuccessMsg(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusOK,
		"msg":  msg,
	})
}
