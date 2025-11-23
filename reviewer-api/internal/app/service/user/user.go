package user

import (
	"reviewer-api/internal/app/ds"
)

type UserRepository interface {
	SetUserFlagDB(user_id string, is_active bool) (ds.User, error)
	GetReviewDB(user_id string) (ds.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}

}

func (p *UserService) SetUserFlag(user_id string, is_active bool) (ds.User, error) {
	user, err := p.repo.SetUserFlagDB(user_id, is_active)
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}

func (p *UserService) GetReview(user_id string) (ds.User, error) {
	user, err := p.repo.GetReviewDB(user_id)
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}
