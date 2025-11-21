package team

import "AvitoTech/internal/domain/dto"

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
