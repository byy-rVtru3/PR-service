package interfaces

import (
	"AvitoTech/internal/domain/dto"
	"context"
)

type PRRepository interface {
	PRExists(ctx context.Context, prID string) (bool, error)

	CreatePR(ctx context.Context, prID, prName, authorID string) error

	GetPR(ctx context.Context, prID string) (*dto.PullRequestDTO, error)

	UpdatePRStatus(ctx context.Context, prID, status string) error

	SetMergedAt(ctx context.Context, prID string) error

	GetReviewers(ctx context.Context, prID string) ([]string, error)
	AssignReviewers(ctx context.Context, prID string, reviewerIDs []string) error
	RemoveReviewer(ctx context.Context, prID, reviewerID string) error
	AddReviewer(ctx context.Context, prID, reviewerID string) error
	IsReviewerAssigned(ctx context.Context, prID, reviewerID string) (bool, error)
}
