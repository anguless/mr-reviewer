package handlers

import (
	"encoding/json"
	"github.com/anguless/reviewer/internal/service"
	"net/http"
)

type StatisticsHandler struct {
	Service service.StatService
}

func (h *StatisticsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Service.GetReviewStats()
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
