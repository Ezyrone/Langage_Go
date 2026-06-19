package api

import (
	"log/slog"
	"net/http"
	"strings"
)

func NewRouter(handler *Handler, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", Healthz)

	mux.HandleFunc("/v1/checks", handler.PostChecks)

	mux.HandleFunc("/v1/checks/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/v1/checks/")
		if id == "" {
			handler.PostChecks(w, r)
			return
		}
		handler.GetCheck(w, r)
	})

	var h http.Handler = mux
	h = LoggingMiddleware(logger)(h)
	h = RecoveryMiddleware(logger)(h)

	return h
}
