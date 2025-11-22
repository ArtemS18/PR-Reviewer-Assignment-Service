package handlers_test

import (
	"reviewer-api/internal/app/http-server/handlers"
	"reviewer-api/internal/app/repository/mocks"
	"reviewer-api/internal/app/service"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	teamService := service.NewTeamService(mocks.MockTeamRepo{})
	userService := service.NewUserService(mocks.MockUserRepo{})
	prService := service.NewPullRequestService(mocks.MockPRRepo{})

	h := &handlers.Handlers{
		Team: handlers.NewTeamHandler(teamService),
		User: handlers.NewUserHandler(userService),
		PR:   handlers.NewPKHandler(prService),
	}

	h.Register(r)

	return r
}
