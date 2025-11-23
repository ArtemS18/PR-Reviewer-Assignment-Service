package team

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
	"reviewer-api/internal/pkg/utils"
)

type TeamRepository interface {
	GetTeam(teamName string) (ds.Team, error)
	CreateTeam(team ds.Team) (ds.Team, error)
	CreateOrUpdateMembers(teamID string, users []dto.UserDTO) ([]ds.User, error)
	ReassgnUsersDB(assignedMap map[string]string, deactivatedIDs []string) error
	DeactivateUsersDB(teamID string) ([]string, error)
	GetNewAssigned(deactivated_ids []string) ([]string, error)

	WithTeamTransaction(fn TxFunc) error
}

type TeamService struct {
	repo TeamRepository
}
type TxFunc func(r TeamRepository) error

func NewTeamService(repo TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) GetTeam(teamName string) (ds.Team, error) {
	return s.repo.GetTeam(teamName)
}

func (s *TeamService) AddTeam(teamData dto.TeamDTO) (ds.Team, error) {
	var result ds.Team
	err := s.repo.WithTeamTransaction(func(r TeamRepository) error {
		if teamData.Name == "" {
			return repository.ErrBadRequest
		}

		team := ds.Team{Name: teamData.Name}
		team, err := r.CreateTeam(team)
		if err != nil {
			return err
		}

		members, err := r.CreateOrUpdateMembers(team.ID, teamData.Members)
		if err != nil {
			return err
		}
		team.Members = members
		result = team
		return nil
	})

	if err != nil {
		return ds.Team{}, err
	}
	return result, nil
}

func (s *TeamService) DeactivateTeam(teamName string) (ds.Team, error) {
	var result ds.Team
	err := s.repo.WithTeamTransaction(func(r TeamRepository) error {
		t, err := r.GetTeam(teamName)
		if err != nil {
			return err
		}

		deactivatedIDs, err := r.DeactivateUsersDB(t.ID)
		if err != nil {
			return err
		}

		assignedIDs, err := r.GetNewAssigned(deactivatedIDs)
		if err != nil {
			return err
		}

		assignedMap := make(map[string]string, len(deactivatedIDs))
		for _, u := range deactivatedIDs {
			assignedMap[u] = utils.GetRandomNumber(assignedIDs)
		}

		if err := r.ReassgnUsersDB(assignedMap, deactivatedIDs); err != nil {
			return err
		}

		team, err := r.GetTeam(teamName)
		if err != nil {
			return err
		}

		result = team
		return nil
	})
	if err != nil {
		return ds.Team{}, err
	}
	return result, nil
}
