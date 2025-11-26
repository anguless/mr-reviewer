package handlers

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/anguless/reviewer/internal/model"
	"github.com/anguless/reviewer/internal/service"

	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TeamHandler struct {
	TeamService service.TeamService
	PRService   service.PrService
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team model.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	createdTeam, err := h.TeamService.CreateTeam(r.Context(), &team)
	if err != nil {
		if errors.Is(err, model.ErrTeamExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			log.Println(err)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("TeamHandler error: createdTeam error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"team": createdTeam,
	})
	if err != nil {
		return
	}
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["team_id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	team, err := h.TeamService.GetTeamByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(team)
	if err != nil {
		return
	}
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["team_id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	var team model.Team
	if err = json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	team.ID = id

	updatedTeam, err := h.TeamService.UpdateTeam(r.Context(), &team)
	if err != nil {
		if err.Error() == "team not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "team with this name already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedTeam)
	if err != nil {
		return
	}
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["team_id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	err = h.TeamService.DeleteTeam(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	_, err = w.Write([]byte("Команда удалена"))
	if err != nil {
		return
	}
}
