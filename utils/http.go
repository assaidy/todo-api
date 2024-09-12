package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

func Make(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("New request", "from", r.RemoteAddr, "path", r.URL.Path, "method", r.Method)
		if err := f(w, r); err != nil {
			if apiErr, ok := err.(ApiError); ok {
				slog.Error("Api error", "err", err.Error(), "path", r.URL.Path)
				WriteJSON(w, apiErr.StatusCode, apiErr)
			} else {
				resp := map[string]string{
					"statusCode": "500",
					"msg":        "internal server error",
				}
				slog.Error("Internal error", "err", err.Error(), "path", r.URL.Path)
				WriteJSON(w, http.StatusInternalServerError, resp)
			}
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Failed to encode JSON response", "err", err.Error())
		return err
	}

	return nil
}
