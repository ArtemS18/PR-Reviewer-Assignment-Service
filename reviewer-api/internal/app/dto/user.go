package dto

import "reviewer-api/internal/app/ds"

type UserDTO struct {
	ID       string `json:"user_id" binding:"required"`
	Name     string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active"`
}
type UserReviewDTO struct {
	UserDTO
	PullRequests []ds.PullRequest `json:"pull_requests"`
}

type UserWithTeamDTO struct {
	UserDTO
	TeamName string `json:"team_name"`
}

func ToUserWithTeamDTO(userORM ds.User) UserWithTeamDTO {
	return UserWithTeamDTO{
		UserDTO: UserDTO{
			ID:       userORM.ID,
			Name:     userORM.Name,
			IsActive: userORM.IsActive,
		},
		TeamName: userORM.Team.Name,
	}

}

func ToUserReviewDTO(userORM ds.User) UserReviewDTO {
	return UserReviewDTO{
		UserDTO: UserDTO{
			ID:       userORM.ID,
			Name:     userORM.Name,
			IsActive: userORM.IsActive,
		},
		PullRequests: userORM.Assigned,
	}

}

func ToUserDTO(userORM ds.User) UserDTO {
	return UserDTO{
		ID:       userORM.ID,
		Name:     userORM.Name,
		IsActive: userORM.IsActive,
	}
}

func ToUserORM(teamID string, userDTO UserDTO) ds.User {
	return ds.User{
		ID:       userDTO.ID,
		Name:     userDTO.Name,
		IsActive: userDTO.IsActive,
		TeamID:   teamID,
	}
}

func ToUserORMList(teamID string, usersDTO []UserDTO) []ds.User {
	usersORM := make([]ds.User, 0, len(usersDTO))
	for _, user := range usersDTO {
		userORM := ToUserORM(teamID, user)
		usersORM = append(usersORM, userORM)
	}
	return usersORM
}
