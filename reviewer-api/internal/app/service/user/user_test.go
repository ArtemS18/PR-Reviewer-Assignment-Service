package user_test

import (
	"reviewer-api/internal/app/repository"
	"reviewer-api/internal/app/repository/mocks"
	service "reviewer-api/internal/app/service/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_SetUserFlag_Success(t *testing.T) {
	svc := service.NewUserService(mocks.MockUserRepo{})

	u, err := svc.SetUserFlag("user-1", true)

	assert.NoError(t, err)
	assert.Equal(t, "user-1", u.ID)
	assert.True(t, u.IsActive)
}

func TestUserService_SetUserFlag_Error(t *testing.T) {
	svc := service.NewUserService(mocks.MockUserRepo{NotFoundErr: true})

	_, err := svc.SetUserFlag("user-1", true)

	assert.Error(t, err)
	assert.Equal(t, err, repository.ErrNotFound)
}

func TestUserService_GetReview_Success(t *testing.T) {
	svc := service.NewUserService(mocks.MockUserRepo{})

	u, err := svc.GetReview("user-1")

	assert.NoError(t, err)
	assert.Equal(t, "user-1", u.ID)
}

func TestUserService_GetReview_Error(t *testing.T) {
	svc := service.NewUserService(mocks.MockUserRepo{NotFoundErr: true})

	_, err := svc.GetReview("user-1")

	assert.Error(t, err)
	assert.Equal(t, err, repository.ErrNotFound)
}
