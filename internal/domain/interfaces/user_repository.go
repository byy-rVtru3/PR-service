package interfaces

import (
	"AvitoTech/internal/domain/dto"
	"context"
)

type UserRepository interface {
	GetUserReviews(ctx context.Context, userID string) ([]dto.PullRequestShortDTO, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) (*dto.UserDTO, error)
	GetUser(ctx context.Context, userID string) (*dto.UserDTO, error)
	CreateOrUpdateUser(ctx context.Context, member dto.TeamMemberDTO, teamName string) error
}
