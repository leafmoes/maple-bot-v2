package api

import (
	"maple-bot/router"
	"net/http"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()
	router.RegisterRouter(mux)
	mux.ServeHTTP(w, req)
}
