package user

import (
	"AvitoTech/internal/domain/dto"
	"AvitoTech/internal/domain/interfaces"
	"context"
	"errors"
)

const (
	BadRequest    = "BAD_REQUEST"
	InternalError = "INTERNAL_ERROR"
	UserNotFound  = "NOT_FOUND"
)

var (
	ErrUserIDEmpty  = errors.New("user_id не может быть пустым")
	ErrUserNotFound = errors.New("пользователь не найден")
)

type Service struct {
	repo interfaces.UserRepository
}

func NewService(repo interfaces.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserReviews(ctx context.Context, userID string) ([]dto.PullRequestShortDTO, error) {
	reviews, err := s.repo.GetUserReviews(ctx, userID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *Service) SetUserActive(ctx context.Context, userID string, isActive bool) (*dto.UserDTO, error) {
	if userID == "" {
		return nil, ErrUserIDEmpty
	}

	user, err := s.repo.SetUserActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}
