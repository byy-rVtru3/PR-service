package interfaces

import (
	"AvitoTech/internal/domain/dto"
	"context"
)

type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, member dto.TeamMemberDTO, team string) error
}
