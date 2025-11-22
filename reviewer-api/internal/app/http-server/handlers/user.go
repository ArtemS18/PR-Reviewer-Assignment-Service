package handlers

import (
	"log"
	"net/http"
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	pkg "reviewer-api/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	SetUserFlag(user_id string, is_active bool) (ds.User, error)
	GetReview(user_id string) (ds.User, error)
}

type UserHandler struct {
	s UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{s: s}
}

type UserUpdateSchema struct {
	UserId   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandler) UpdateUserActivity(ctx *gin.Context) {
	var userData UserUpdateSchema
	err := ctx.BindJSON(&userData)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			pkg.BAD_REQUEST,
		)
		return
	}
	user, err := h.s.SetUserFlag(userData.UserId, userData.IsActive)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	pkg.OkResponse(ctx, dto.ToUserWithTeamDTO(user))
}

func (h *UserHandler) GetUserReview(ctx *gin.Context) {
	userName := ctx.Query("user_id")
	if userName == "" {
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			pkg.NOT_FOUND,
		)
		return
	}
	user, err := h.s.GetReview(userName)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	pkg.OkResponse(ctx, dto.ToUserReviewDTO(user))
}
