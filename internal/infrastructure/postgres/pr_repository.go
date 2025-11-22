package postgres

import (
	"AvitoTech/internal/domain/dto"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	prExistsQuery = `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`

	createPRQuery = `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	getPRQuery = `
		SELECT pull_request_id, pull_request_name, author_id, status
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	updatePRStatusQuery = `
		UPDATE pull_requests
		SET status = $2
		WHERE pull_request_id = $1
	`

	setMergedAtQuery = `
		UPDATE pull_requests
		SET merged_at = $2
		WHERE pull_request_id = $1
	`

	getReviewersQuery = `
		SELECT reviewer_id
		FROM pr_reviewers
		WHERE pull_request_id = $1
		ORDER BY assigned_at
	`

	assignReviewerQuery = `
		INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`

	removeReviewerQuery = `
		DELETE FROM pr_reviewers
		WHERE pull_request_id = $1 AND reviewer_id = $2
	`

	isReviewerAssignedQuery = `
		SELECT EXISTS(
			SELECT 1 FROM pr_reviewers
			WHERE pull_request_id = $1 AND reviewer_id = $2
		)
	`
)

type PRRepo struct {
	db *pgx.Conn
}

func NewPRRepo(db *Postgres) *PRRepo {
	return &PRRepo{db: db.conn}
}

func (r *PRRepo) PRExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, prExistsQuery, prID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке существования PR: %v", err)
	}
	return exists, nil
}

func (r *PRRepo) CreatePR(ctx context.Context, prID, prName, authorID string) error {
	_, err := r.db.Exec(ctx, createPRQuery, prID, prName, authorID, dto.StatusOpen, time.Now())
	if err != nil {
		return fmt.Errorf("ошибка при создании PR: %v", err)
	}
	return nil
}

func (r *PRRepo) GetPR(ctx context.Context, prID string) (*dto.PullRequestDTO, error) {
	var pr dto.PullRequestDTO
	err := r.db.QueryRow(ctx, getPRQuery, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("PR не найден")
		}
		return nil, fmt.Errorf("ошибка при получении PR: %v", err)
	}

	reviewers, err := r.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *PRRepo) UpdatePRStatus(ctx context.Context, prID, status string) error {
	_, err := r.db.Exec(ctx, updatePRStatusQuery, prID, status)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса PR: %v", err)
	}
	return nil
}

func (r *PRRepo) SetMergedAt(ctx context.Context, prID string) error {
	_, err := r.db.Exec(ctx, setMergedAtQuery, prID, time.Now())
	if err != nil {
		return fmt.Errorf("ошибка при установке времени мерджа: %v", err)
	}
	return nil
}

func (r *PRRepo) GetReviewers(ctx context.Context, prID string) ([]string, error) {
	rows, err := r.db.Query(ctx, getReviewersQuery, prID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении ревьюверов: %v", err)
	}
	defer rows.Close()

	reviewers := []string{}
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("ошибка при чтении ревьювера: %v", err)
		}
		reviewers = append(reviewers, reviewerID)
	}

	return reviewers, nil
}

func (r *PRRepo) AssignReviewers(ctx context.Context, prID string, reviewerIDs []string) error {
	for _, reviewerID := range reviewerIDs {
		if err := r.AddReviewer(ctx, prID, reviewerID); err != nil {
			return err
		}
	}
	return nil
}

func (r *PRRepo) RemoveReviewer(ctx context.Context, prID, reviewerID string) error {
	_, err := r.db.Exec(ctx, removeReviewerQuery, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении ревьювера: %v", err)
	}
	return nil
}

func (r *PRRepo) AddReviewer(ctx context.Context, prID, reviewerID string) error {
	_, err := r.db.Exec(ctx, assignReviewerQuery, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("ошибка при назначении ревьювера: %v", err)
	}
	return nil
}

func (r *PRRepo) IsReviewerAssigned(ctx context.Context, prID, reviewerID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, isReviewerAssignedQuery, prID, reviewerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке назначения ревьювера: %v", err)
	}
	return exists, nil
}
