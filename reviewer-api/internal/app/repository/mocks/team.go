package mocks

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
)

type MockTeamRepo struct {
	CreateErr  bool
	MembersErr bool
	GetErr     bool
}

func (m MockTeamRepo) GetTeam(teamName string) (ds.Team, error) {
	if m.GetErr {
		return ds.Team{}, repository.ErrNotFound
	}
	return ds.Team{Name: teamName}, nil
}

func (m MockTeamRepo) CreateTeam(team ds.Team) (ds.Team, error) {
	if m.CreateErr {
		return ds.Team{}, repository.ErrTeamAlreadyExists
	}
	team.ID = "team-1"
	return team, nil
}

func (m MockTeamRepo) CreateOrUpdateMembers(teamID string, users []dto.UserDTO) ([]ds.User, error) {
	if m.MembersErr {
		return nil, repository.ErrUnexpect
	}
	res := make([]ds.User, 0, len(users))
	for _, u := range users {
		res = append(res, ds.User{
			ID:       u.ID,
			Name:     u.Name,
			IsActive: u.IsActive,
			TeamID:   teamID,
		})
	}
	return res, nil
}
