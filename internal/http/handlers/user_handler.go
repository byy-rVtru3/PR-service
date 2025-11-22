package handlers

import (
	"AvitoTech/internal/domain/dto"
	"AvitoTech/internal/domain/user"
	"encoding/json"
	"errors"
	"net/http"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    user.BadRequest,
				Message: "user_id is required",
			},
		})
		return
	}

	reviews, err := h.service.GetUserReviews(r.Context(), userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    user.InternalError,
				Message: "internal server error",
			},
		})
		return
	}

	response := dto.GetUserReviewsResponse{
		UserID:       userID,
		PullRequests: reviews,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) SetUserActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    user.BadRequest,
				Message: "invalid request body",
			},
		})
		return
	}

	u, err := h.service.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if errors.Is(err, user.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error: dto.Error{
					Code:    user.UserNotFound,
					Message: "user not found",
				},
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    user.InternalError,
				Message: "internal server error",
			},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.UserResponse{User: *u})
}
