package handlers

import (
	"AvitoTech/internal/domain/dto"
	"AvitoTech/internal/domain/team"
	"encoding/json"
	"fmt"
	"net/http"
)

type TeamHandler struct {
	service *team.Service
}

func NewTeamHandler(service *team.Service) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидные данные", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err := h.service.CreateTeam(ctx, req)
	if err != nil {
		if err.Error() == "команда с таким именем уже существует" {
			http.Error(w, "Команда с таким именем уже существует", http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("Ошибка при создании команды: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "team created"})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	ctx := r.Context()

	team, err := h.service.GetTeam(ctx, teamName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при получении команды: %v", err), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}
