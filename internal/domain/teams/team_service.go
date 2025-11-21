package teams

import (
	"AvitoTech/internal/domain/dto"
	"AvitoTech/internal/domain/interfaces"
	"AvitoTech/internal/infrastructure/postgres"
	"context"
	"fmt"
)

type Service struct {
	teams *postgres.TeamRepo
	users interfaces.UserRepository
}

func NewService(t *postgres.TeamRepo, u interfaces.UserRepository) *Service {
	return &Service{teams: t, users: u}
}

func (s *Service) CreateTeam(ctx context.Context, req dto.TeamDTO) error {
	exists, err := s.teams.TeamExists(ctx, req.TeamName)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования команды: %v", err)
	}
	if exists {
		return ErrTeamExists
	}

	err = s.teams.CreateTeam(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка при создании команды: %v", err)
	}

	for _, member := range req.Members {
		err = s.users.CreateOrUpdateUser(ctx, member, req.TeamName)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении пользователя: %v", err)
		}
	}

	return nil
}

func (s *Service) GetTeam(ctx context.Context, teamName string) (dto.TeamDTO, error) {
	return s.teams.GetTeam(ctx, teamName)
}
