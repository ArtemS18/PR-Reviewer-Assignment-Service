package mocks

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/repository"
)

type MockUserRepo struct {
	NotFoundErr bool
}

func (m MockUserRepo) SetUserFlagDB(user_id string, is_active bool) (ds.User, error) {
	if m.NotFoundErr {
		return ds.User{}, repository.ErrNotFound
	}
	return ds.User{ID: user_id, IsActive: is_active, Name: "Test"}, nil
}

func (m MockUserRepo) GetReviewDB(user_id string) (ds.User, error) {
	if m.NotFoundErr {
		return ds.User{}, repository.ErrNotFound
	}
	return ds.User{ID: user_id, IsActive: true, Name: "Test"}, nil
}
