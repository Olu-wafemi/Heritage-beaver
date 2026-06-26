package handlers

import (
	"database/sql"
	"net/http"
	"time"
)

type HealthHandler struct {
	db *sql.DB
}

type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

func NewHealthHandler(db *sql.DB) HealthHandler {
	return HealthHandler{db: db}
}

func (h HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	databaseStatus := "up"
	if err := h.db.PingContext(r.Context()); err != nil {
		databaseStatus = "down"
		writeJSON(w, http.StatusServiceUnavailable, healthResponse{
			Status:    "degraded",
			Service:   "heritage-weaver-api",
			Database:  databaseStatus,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	writeJSON(w, http.StatusOK, healthResponse{
		Status:    "ok",
		Service:   "heritage-weaver-api",
		Database:  databaseStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
