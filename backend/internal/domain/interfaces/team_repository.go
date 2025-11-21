package interfaces

import (
	"AvitoTech/internal/domain/dto"
	"context"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team dto.TeamDTO) error
	TeamExist(ctx context.Context, name string) (bool, error)
	GetTeam(ctx context.Context, name string) (dto.TeamDTO, error)
}
