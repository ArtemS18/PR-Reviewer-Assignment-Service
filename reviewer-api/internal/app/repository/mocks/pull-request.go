package mocks

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/repository"
)

type MockPRRepo struct {
	MemberErr       bool
	CreateErr       bool
	AssignErr       bool
	FindErr         bool
	UpdateErr       bool
	UpdateRevErr    bool
	ForcedStatus    string
	ReturnedMembers []string
	EmptyMembers    bool
}

func (m MockPRRepo) CreatePullRequestDB(pk ds.PullRequest) (ds.PullRequest, error) {
	if m.CreateErr {
		return ds.PullRequest{}, repository.ErrPRAlreadyExists
	}
	return pk, nil
}

func (m MockPRRepo) UpdatePullRequestDB(pk ds.PullRequest) (ds.PullRequest, error) {
	if m.UpdateErr {
		return ds.PullRequest{}, repository.ErrNotFound
	}
	return pk, nil
}

func (m MockPRRepo) FindPullRequestByID(id string) (ds.PullRequest, error) {
	if m.FindErr {
		return ds.PullRequest{}, repository.ErrNotFound
	}
	status := m.ForcedStatus
	if status == "" {
		status = string(ds.OPEN)
	}
	return ds.PullRequest{
		ID:     id,
		Name:   "PR",
		Status: status,
	}, nil
}

func (m MockPRRepo) AssignReviewersDB(rs []ds.Reviewer) error {
	if m.AssignErr {
		return repository.ErrReassign
	}
	return nil
}

func (m MockPRRepo) GetMemberIDsDB(excludedID, pkID string) ([]string, error) {
	if m.MemberErr {
		return nil, repository.ErrNotFound
	}
	if m.ReturnedMembers != nil {
		return m.ReturnedMembers, nil
	}
	if m.EmptyMembers {
		return []string{}, nil
	}
	return []string{"u1", "u2", "u3"}, nil

}

func (m MockPRRepo) UpdateReviewersDB(r ds.Reviewer) error {
	if m.UpdateRevErr {
		return repository.ErrNotFound
	}
	return nil
}
