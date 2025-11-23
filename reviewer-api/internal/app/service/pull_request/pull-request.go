package pull_request

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

	WithPRTransaction(fn TxFunc) error
}

type PullRequestService struct {
	repo PullRequestRepository
}
type TxFunc func(r PullRequestRepository) error

func NewPullRequestService(repo PullRequestRepository) *PullRequestService {
	return &PullRequestService{repo: repo}
}

func (s *PullRequestService) CreatePullRequest(pkDTO dto.PullRequestCreateDTO) (ds.PullRequest, error) {
	var result ds.PullRequest
	err := s.repo.WithPRTransaction(func(r PullRequestRepository) error {
		membersIDs, err := r.GetMemberIDsDB(pkDTO.AuthorID, pkDTO.ID)
		if err != nil {
			return repository.ErrNotFound
		}

		pk := ds.PullRequest{
			ID:       pkDTO.ID,
			Name:     pkDTO.Name,
			AuthorID: pkDTO.AuthorID,
		}

		pk, err = r.CreatePullRequestDB(pk)
		if err != nil {
			return err
		}

		assignedIDs := utils.GetRandomSlice(membersIDs)
		if len(assignedIDs) == 0 {
			return nil
		}
		revs := make([]ds.Reviewer, 0, len(assignedIDs))
		for _, id := range assignedIDs {
			revs = append(revs, ds.Reviewer{UserID: id, PullRequestID: pk.ID})
		}
		if err := r.AssignReviewersDB(revs); err != nil {
			return err
		}
		pk.AssignedReviewers = revs
		result = pk
		return nil
	})
	if err != nil {
		return ds.PullRequest{}, err
	}

	return result, nil
}

func (s *PullRequestService) ReassignReviewer(pkID, oldReviewerID string) (ds.PullRequest, error) {
	var result ds.PullRequest
	err := s.repo.WithPRTransaction(func(r PullRequestRepository) error {
		membersIDs, err := r.GetMemberIDsDB(oldReviewerID, pkID)
		if err != nil {
			return repository.ErrNotFound
		}

		newReviewerID := utils.GetRandomNumber(membersIDs)
		if newReviewerID == "" {
			return repository.ErrNotFound
		}

		pk, err := r.FindPullRequestByID(pkID)
		if err != nil {
			return err
		}
		if pk.Status == string(ds.MERGED) {
			return repository.ErrReassign
		}

		err = r.UpdateReviewersDB(
			ds.Reviewer{
				UserID:        newReviewerID,
				PullRequestID: pkID,
			},
		)
		if err != nil {
			return err
		}
		pk, err = r.FindPullRequestByID(pkID)
		if err != nil {
			return err
		}
		result = pk
		return nil
	})
	if err != nil {
		return ds.PullRequest{}, err
	}
	return result, nil
}

func (s *PullRequestService) Merged(pkID string) (ds.PullRequest, error) {
	var result ds.PullRequest
	err := s.repo.WithPRTransaction(func(r PullRequestRepository) error {
		pk, err := r.FindPullRequestByID(pkID)
		if err != nil {
			return err
		}
		if pk.Status == string(ds.MERGED) {
			result = pk
			return nil
		}

		now := time.Now().UTC()
		pk.Status = string(ds.MERGED)
		pk.MergedAt = &now

		pk, err = r.UpdatePullRequestDB(pk)
		if err != nil {
			return err
		}
		result = pk
		return nil
	})
	if err != nil {
		return ds.PullRequest{}, err
	}
	return result, nil
}
