package handlers

import (
	"AvitoTech/internal/domain/dto"
	"AvitoTech/internal/domain/teams"
	"encoding/json"
	"errors"
	"net/http"
)

type TeamHandler struct {
	service *teams.Service
}

func NewTeamHandler(service *teams.Service) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    teams.BadRequest,
				Message: "invalid request body",
			},
		})
		return
	}

	ctx := r.Context()
	err := h.service.CreateTeam(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, teams.ErrTeamExists) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error: dto.Error{
					Code:    teams.TeamExistsCode,
					Message: "team_name already exists",
				},
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    teams.InternalError,
				Message: "internal server error",
			},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	if err := json.NewEncoder(w).Encode(dto.TeamResponse{Team: req}); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    teams.BadRequest,
				Message: "team_name is required",
			},
		})
		return
	}

	ctx := r.Context()
	t, err := h.service.GetTeam(ctx, teamName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    teams.TeamNotFound,
				Message: "resource not found",
			},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}
