package team

import (
	"AvitoTech/internal/domain/dto"
	"AvitoTech/internal/domain/interfaces"
	"context"
	"fmt"
)

type Service struct {
	teams interfaces.TeamRepository
	users interfaces.UserRepository
}

func NewService(t interfaces.TeamRepository, u interfaces.UserRepository) *Service {
	return &Service{teams: t, users: u}
}

func (s *Service) CreateTeam(ctx context.Context, req dto.TeamDTO) error {
	exists, err := s.teams.TeamExist(ctx, req.TeamName)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования команды: %v", err)
	}
	if exists {
		return fmt.Errorf("команда с таким именем уже существует")
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
