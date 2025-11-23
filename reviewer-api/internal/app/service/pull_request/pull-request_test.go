package pull_request_test

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
	"reviewer-api/internal/app/repository/mocks"
	"reviewer-api/internal/app/service/pull_request"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPullRequestService_CreatePullRequest_Success(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{})

	pr, err := svc.CreatePullRequest(dto.PullRequestCreateDTO{
		ID:       "pr-1",
		Name:     "Test PR",
		AuthorID: "user-author",
	})

	assert.NoError(t, err)
	assert.Equal(t, "pr-1", pr.ID)
}

func TestPullRequestService_CreatePullRequest_NoMembers(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{ReturnedMembers: []string{}})

	pr, err := svc.CreatePullRequest(dto.PullRequestCreateDTO{
		ID:       "pr-1",
		Name:     "Test PR",
		AuthorID: "user-author",
	})

	assert.NoError(t, err)
	assert.Equal(t, "pr-1", pr.ID)
	assert.Len(t, pr.AssignedReviewers, 0)
}

func TestPullRequestService_CreatePullRequest_MembersError(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{MemberErr: true})

	_, err := svc.CreatePullRequest(dto.PullRequestCreateDTO{
		ID:       "pr-1",
		Name:     "Test PR",
		AuthorID: "user-author",
	})

	assert.Error(t, err)
}

func TestPullRequestService_ReassignReviewer_Success(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{})

	pr, err := svc.ReassignReviewer("pr-1", "old-reviewer")

	assert.NoError(t, err)
	assert.Equal(t, "pr-1", pr.ID)
}

func TestPullRequestService_ReassignReviewer_NoMembers(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{EmptyMembers: true})

	_, err := svc.ReassignReviewer("pr-1", "old-reviewer")

	assert.Error(t, err)
	assert.Equal(t, err, repository.ErrNotFound)
}
func TestPullRequestService_ReassignReviewer_MergedPR(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{ForcedStatus: string(ds.MERGED)})

	_, err := svc.ReassignReviewer("pr-1", "old-reviewer")

	assert.Error(t, err)
	assert.Equal(t, err, repository.ErrReassign)
}

func TestPullRequestService_Merged_Success(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{})

	pr, err := svc.Merged("pr-1")

	assert.NoError(t, err)
	assert.Equal(t, string(ds.MERGED), pr.Status)
	assert.NotNil(t, pr.MergedAt)
	assert.WithinDuration(t, time.Now().UTC(), *pr.MergedAt, time.Second)
}

func TestPullRequestService_Merged_AlreadyMerged(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{ForcedStatus: string(ds.MERGED)})

	pr, err := svc.Merged("pr-1")

	assert.NoError(t, err)
	assert.Equal(t, string(ds.MERGED), pr.Status)
}

func TestPullRequestService_Merged_FindError(t *testing.T) {
	svc := pull_request.NewPullRequestService(mocks.MockPRRepo{FindErr: true})

	_, err := svc.Merged("pr-1")

	assert.Error(t, err)
	assert.Equal(t, err, repository.ErrNotFound)
}
