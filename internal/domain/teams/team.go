package teams

import (
	"AvitoTech/internal/domain/dto"
	"errors"
)

type Team struct {
	Name    string
	Members []dto.TeamMemberDTO
}

func NewTeam(dto dto.TeamDTO) Team {
	return Team{
		Name:    dto.TeamName,
		Members: dto.Members,
	}
}

const TeamExistsCode = "TEAM_EXISTS"

var ErrTeamExists = errors.New("teams already exists")
