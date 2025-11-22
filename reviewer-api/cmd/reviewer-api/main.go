package main

import (
	"log"
	"reviewer-api/internal/app/config"
	http "reviewer-api/internal/app/http-server/handlers"
	"reviewer-api/internal/app/repository/postgres"
	"reviewer-api/internal/app/service"
	pkg "reviewer-api/internal/pkg/app"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	pg, err := postgres.NewPostgers(cfg.GetDSN(), true)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	teamService := service.NewTeamService(pg)
	userService := service.NewUserService(pg)
	prService := service.NewPullRequestService(pg)

	handl := &http.Handlers{
		Team: http.NewTeamHandler(teamService),
		User: http.NewUserHandler(userService),
		PR:   http.NewPKHandler(prService),
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	app := pkg.NewApplication(cfg, router, handl)
	app.RunApplication()

}
