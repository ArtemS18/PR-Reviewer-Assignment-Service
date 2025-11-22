package service

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
)

type TeamRepository interface {
	GetTeam(teamName string) (ds.Team, error)
	CreateTeam(team ds.Team) (ds.Team, error)
	CreateOrUpdateMembers(teamID string, users []dto.UserDTO) ([]ds.User, error)
}

type TeamService struct {
	repo TeamRepository
}

func NewTeamService(repo TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) GetTeam(teamName string) (ds.Team, error) {
	return s.repo.GetTeam(teamName)
}

func (s *TeamService) AddTeam(teamData dto.TeamDTO) (ds.Team, error) {
	if teamData.Name == "" {
		return ds.Team{}, repository.ErrBadRequest
	}

	team := ds.Team{Name: teamData.Name}
	team, err := s.repo.CreateTeam(team)
	if err != nil {
		return ds.Team{}, err
	}

	members, err := s.repo.CreateOrUpdateMembers(team.ID, teamData.Members)
	if err != nil {
		return ds.Team{}, err
	}
	team.Members = members

	return team, nil
}
