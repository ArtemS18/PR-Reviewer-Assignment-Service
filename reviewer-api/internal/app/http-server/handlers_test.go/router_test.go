package handlers_test

import (
	"reviewer-api/internal/app/http-server/handlers"
	"reviewer-api/internal/app/repository/mocks"
	"reviewer-api/internal/app/service/pull_request"
	"reviewer-api/internal/app/service/team"
	"reviewer-api/internal/app/service/user"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	teamService := team.NewTeamService(mocks.MockTeamRepo{})
	userService := user.NewUserService(mocks.MockUserRepo{})
	prService := pull_request.NewPullRequestService(mocks.MockPRRepo{})

	h := &handlers.Handlers{
		Team: handlers.NewTeamHandler(teamService),
		User: handlers.NewUserHandler(userService),
		PR:   handlers.NewPKHandler(prService),
	}

	h.Register(r)

	return r
}
