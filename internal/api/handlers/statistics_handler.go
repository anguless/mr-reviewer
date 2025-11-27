package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anguless/mr-reviewer/internal/service"
)

type StatisticsHandler struct {
	StatService service.StatService
}

func (h *StatisticsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.StatService.GetReviewStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		return
	}
}
