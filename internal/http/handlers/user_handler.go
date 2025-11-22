package handlers

import (
	"AvitoTech/internal/domain/dto"
	"AvitoTech/internal/domain/user"
	"AvitoTech/pkg/logger"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
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
		logger.Log.Warn("Запрос получения PR без user_id")
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

	logger.Log.Info("Получение PR пользователя", zap.String("user_id", userID))

	reviews, err := h.service.GetUserReviews(r.Context(), userID)
	if err != nil {
		logger.Log.Error("Ошибка получения PR пользователя",
			zap.String("user_id", userID),
			zap.Error(err),
		)
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

	logger.Log.Info("PR пользователя успешно получены",
		zap.String("user_id", userID),
		zap.Int("pr_count", len(reviews)),
	)

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
		logger.Log.Warn("Неверный формат запроса изменения активности", zap.Error(err))
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

	logger.Log.Info("Изменение активности пользователя",
		zap.String("user_id", req.UserID),
		zap.Bool("is_active", req.IsActive),
	)

	u, err := h.service.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if errors.Is(err, user.ErrUserNotFound) {
			logger.Log.Warn("Пользователь не найден", zap.String("user_id", req.UserID))
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error: dto.Error{
					Code:    user.UserNotFound,
					Message: "user not found",
				},
			})
			return
		}

		logger.Log.Error("Ошибка изменения активности пользователя",
			zap.String("user_id", req.UserID),
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: dto.Error{
				Code:    user.InternalError,
				Message: "internal server error",
			},
		})
		return
	}

	logger.Log.Info("Активность пользователя успешно изменена",
		zap.String("user_id", req.UserID),
		zap.Bool("is_active", req.IsActive),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.UserResponse{User: *u})
}
