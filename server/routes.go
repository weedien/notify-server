package server

import (
	"net/http"

	api "github.com/weedien/countdown-server/api/countdown"
)

func InitRoutes() {
	http.HandleFunc("POST /countdowns", api.CreateCountdown)
	http.HandleFunc("PUT /countdowns/{id}", api.UpdateCountdown)
	http.HandleFunc("DELETE /countdowns/{id}", api.DeleteCountdown)
	http.HandleFunc("GET /countdowns/{id}", api.GetCountdown)
	http.HandleFunc("GET /countdowns", api.GetAllCountdown)
	http.HandleFunc("GET /countdowns/code/{query_code}", api.GetCountdownByQueryCode)
}
