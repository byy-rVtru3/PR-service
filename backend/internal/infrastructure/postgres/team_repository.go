package postgres

import (
	"AvitoTech/internal/domain/dto"
	"AvitoTech/internal/domain/interfaces"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type TeamRepo struct {
	db *pgx.Conn
}

func NewTeamRepo(db *pgx.Conn) interfaces.TeamRepository {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) TeamExist(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке существования команды: %v", err)
	}
	return exists, nil
}

func (r *TeamRepo) CreateTeam(ctx context.Context, team dto.TeamDTO) error {
	_, err := r.db.Exec(ctx, `INSERT INTO teams(team_name) VALUES ($1)`, team.TeamName)
	if err != nil {
		return fmt.Errorf("ошибка при создании команды: %v", err)
	}

	for _, member := range team.Members {
		_, err := r.db.Exec(ctx, `INSERT INTO users(user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4)`,
			member.UserID, member.Username, team.TeamName, member.IsActive)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении участника команды: %v", err)
		}
	}

	return nil
}

func (r *TeamRepo) GetTeam(ctx context.Context, name string) (dto.TeamDTO, error) {
	row := r.db.QueryRow(ctx, `SELECT team_name FROM teams WHERE team_name = $1`, name)
	var teamName string
	err := row.Scan(&teamName)
	if err != nil {
		return dto.TeamDTO{}, fmt.Errorf("команда не найдена: %v", err)
	}

	rows, err := r.db.Query(ctx, `SELECT user_id, username, is_active FROM users WHERE team_name = $1`, teamName)
	if err != nil {
		return dto.TeamDTO{}, fmt.Errorf("ошибка при загрузке участников команды: %v", err)
	}
	defer rows.Close()

	var members []dto.TeamMemberDTO
	for rows.Next() {
		var member dto.TeamMemberDTO
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return dto.TeamDTO{}, fmt.Errorf("ошибка при чтении участников команды: %v", err)
		}
		members = append(members, member)
	}

	return dto.TeamDTO{
		TeamName: teamName,
		Members:  members,
	}, nil
}
