package service

import (
	"time"

	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
	"reviewer-api/internal/pkg/utils"
)

type PullRequestRepository interface {
	CreatePullRequestDB(pk ds.PullRequest) (ds.PullRequest, error)
	UpdatePullRequestDB(pk ds.PullRequest) (ds.PullRequest, error)
	FindPullRequestByID(id string) (ds.PullRequest, error)
	AssignReviewersDB(rs []ds.Reviewer) error
	GetMemberIDsDB(excludedID, pkID string) ([]string, error)
	UpdateReviewersDB(rs ds.Reviewer) error
}

type PullRequestService struct {
	repo PullRequestRepository
}

func NewPullRequestService(repo PullRequestRepository) *PullRequestService {
	return &PullRequestService{repo: repo}
}

func (s *PullRequestService) assignReviewers(membersIDs []string, pkID string) ([]ds.Reviewer, error) {
	assignedIDs := utils.GetRandomSlice(membersIDs)
	if len(assignedIDs) == 0 {
		return nil, nil
	}
	revs := make([]ds.Reviewer, 0, len(assignedIDs))
	for _, id := range assignedIDs {
		revs = append(revs, ds.Reviewer{UserID: id, PullRequestID: pkID})
	}
	if err := s.repo.AssignReviewersDB(revs); err != nil {
		return nil, err
	}
	return revs, nil
}

func (s *PullRequestService) CreatePullRequest(pkDTO dto.PullRequestCreateDTO) (ds.PullRequest, error) {
	membersIDs, err := s.repo.GetMemberIDsDB(pkDTO.AuthorID, pkDTO.ID)
	if err != nil {
		return ds.PullRequest{}, repository.ErrNotFound
	}

	pk := ds.PullRequest{
		ID:       pkDTO.ID,
		Name:     pkDTO.Name,
		AuthorID: pkDTO.AuthorID,
	}

	pk, err = s.repo.CreatePullRequestDB(pk)
	if err != nil {
		return ds.PullRequest{}, err
	}

	revs, err := s.assignReviewers(membersIDs, pk.ID)
	if err != nil {
		return ds.PullRequest{}, err
	}
	pk.AssignedReviewers = revs

	return pk, nil
}

func (s *PullRequestService) ReassignReviewer(pkID, oldReviewerID string) (ds.PullRequest, error) {
	membersIDs, err := s.repo.GetMemberIDsDB(oldReviewerID, pkID)
	if err != nil {
		return ds.PullRequest{}, repository.ErrNotFound
	}

	newReviewerID := utils.GetRandomNumber(membersIDs)
	if newReviewerID == "" {
		return ds.PullRequest{}, repository.ErrNotFound
	}

	pk, err := s.repo.FindPullRequestByID(pkID)
	if err != nil {
		return ds.PullRequest{}, err
	}
	if pk.Status == string(ds.MERGED) {
		return ds.PullRequest{}, repository.ErrReassign
	}

	err = s.repo.UpdateReviewersDB(
		ds.Reviewer{
			UserID:        newReviewerID,
			PullRequestID: pkID,
		},
	)
	if err != nil {
		return ds.PullRequest{}, err
	}
	return s.repo.FindPullRequestByID(pkID)
}

func (s *PullRequestService) Merged(pkID string) (ds.PullRequest, error) {
	pk, err := s.repo.FindPullRequestByID(pkID)
	if err != nil {
		return ds.PullRequest{}, err
	}
	if pk.Status == string(ds.MERGED) {
		return pk, nil
	}
	now := time.Now().UTC()
	pk.Status = string(ds.MERGED)
	pk.MergedAt = &now

	pk, err = s.repo.UpdatePullRequestDB(pk)
	if err != nil {
		return ds.PullRequest{}, err
	}
	return pk, nil
}
