package team_test

import (
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
	"reviewer-api/internal/app/repository/mocks"
	"reviewer-api/internal/app/service/team"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamService_GetTeam_Success(t *testing.T) {
	svc := team.NewTeamService(mocks.MockTeamRepo{})

	team, err := svc.GetTeam("team-1")

	assert.NoError(t, err)
	assert.Equal(t, "team-1", team.Name)
}

func TestTeamService_GetTeam_Error(t *testing.T) {
	svc := team.NewTeamService(mocks.MockTeamRepo{GetErr: true})

	_, err := svc.GetTeam("team-1")

	assert.Error(t, err)
	assert.Equal(t, err, repository.ErrNotFound)
}

func TestTeamService_AddTeam_EmptyName(t *testing.T) {
	svc := team.NewTeamService(mocks.MockTeamRepo{})

	_, err := svc.AddTeam(dto.TeamDTO{
		Name:    "",
		Members: nil,
	})

	assert.Error(t, err)
}

func TestTeamService_AddTeam_Success(t *testing.T) {
	svc := team.NewTeamService(mocks.MockTeamRepo{})

	team, err := svc.AddTeam(dto.TeamDTO{
		Name: "team-1",
		Members: []dto.UserDTO{
			{ID: "user-1", Name: "User 1", IsActive: true},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "team-1", team.Name)
	assert.Len(t, team.Members, 1)
	assert.Equal(t, "user-1", team.Members[0].ID)
}

func TestTeamService_AddTeam_CreateTeamError(t *testing.T) {
	svc := team.NewTeamService(mocks.MockTeamRepo{CreateErr: true})

	_, err := svc.AddTeam(dto.TeamDTO{
		Name: "team-1",
	})

	assert.Error(t, err)
	assert.Equal(t, err, repository.ErrTeamAlreadyExists)
}

func TestTeamService_AddTeam_CreateMembersError(t *testing.T) {
	svc := team.NewTeamService(mocks.MockTeamRepo{MembersErr: true})

	_, err := svc.AddTeam(dto.TeamDTO{
		Name: "team-1",
		Members: []dto.UserDTO{
			{ID: "user-1", Name: "User 1", IsActive: true},
		},
	})

	assert.Error(t, err)
}
