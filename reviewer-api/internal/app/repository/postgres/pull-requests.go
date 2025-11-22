package postgres

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/repository"
)

func (p *Postgres) CreatePullRequestDB(pk ds.PullRequest) (ds.PullRequest, error) {
	err := p.db.Create(&pk).Error
	if err != nil {
		return ds.PullRequest{}, repository.HandelPgError(err, "pr")
	}
	return pk, nil
}

func (p *Postgres) UpdatePullRequestDB(pk ds.PullRequest) (ds.PullRequest, error) {
	err := p.db.Save(&pk).Error
	if err != nil {
		return ds.PullRequest{}, repository.HandelPgError(err, "pr")
	}
	return pk, nil
}

func (p *Postgres) FindPullRequestByID(id string) (ds.PullRequest, error) {
	var pk ds.PullRequest
	err := p.db.Preload("AssignedReviewers").Where("id = ?", id).First(&pk).Error
	if err != nil {
		return ds.PullRequest{}, repository.HandelPgError(err, "pk")
	}
	return pk, nil
}

func (p *Postgres) AssignReviewersDB(rs []ds.Reviewer) error {
	if len(rs) == 0 {
		return nil
	}
	if err := p.db.Create(&rs).Error; err != nil {
		return repository.HandelPgError(err, "pk")
	}
	return nil
}

func (p *Postgres) UpdateReviewersDB(rs ds.Reviewer) error {
	if err := p.db.Updates(&rs).Error; err != nil {
		return repository.HandelPgError(err, "pk")
	}
	return nil
}

func (p *Postgres) GetMemberIDsDB(excludedID, pkID string) ([]string, error) {
	sub := p.db.Model(&ds.Team{}).
		Select("teams.id").
		Joins("JOIN users ON teams.id = users.team_id").
		Where("users.id = ?", excludedID)

	reviewerSub := p.db.Model(&ds.Reviewer{}).
		Distinct("user_id").
		Where("pull_request_id = ?", pkID)

	var memberIDs []string
	err := p.db.Model(&ds.User{}).
		Select("users.id").
		Where("users.team_id = (?) AND users.is_active = true AND users.id != ?", sub, excludedID).
		Where("users.id NOT IN (?)", reviewerSub).
		Scan(&memberIDs).Error

	if err != nil {
		return nil, repository.HandelPgError(err, "users")
	}
	return memberIDs, nil
}
