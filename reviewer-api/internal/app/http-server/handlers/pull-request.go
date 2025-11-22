package handlers

import (
	"net/http"
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	pkg "reviewer-api/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type PKService interface {
	CreatePullRequest(pkDTO dto.PullRequestCreateDTO) (ds.PullRequest, error)
	ReassignReviewer(pk_id string, old_reviewer_id string) (ds.PullRequest, error)
	Merged(pk_id string) (ds.PullRequest, error)
}

type PKHandler struct {
	s PKService
}

func NewPKHandler(s PKService) *PKHandler {
	return &PKHandler{s: s}
}

func (h *PKHandler) CreateNewPullRequest(ctx *gin.Context) {
	var pkDTO dto.PullRequestCreateDTO
	err := ctx.BindJSON(&pkDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			pkg.BAD_REQUEST,
		)
		return
	}
	team, err := h.s.CreatePullRequest(pkDTO)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, dto.ToPullRequestDTO(team))
}

type ReassgnPRSchema struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldRevID      string `json:"old_reviewer_id" binding:"required"`
}

func (h *PKHandler) ReassignPullRequest(ctx *gin.Context) {
	var raw ReassgnPRSchema
	err := ctx.BindJSON(&raw)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			pkg.BAD_REQUEST,
		)
		return
	}
	pr, err := h.s.ReassignReviewer(raw.PullRequestID, raw.OldRevID)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	pkg.OkResponse(ctx, dto.ToPullRequestDTO(pr))
}

type MergedPRSchema struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

func (h *PKHandler) MergedPR(ctx *gin.Context) {
	var raw MergedPRSchema
	err := ctx.BindJSON(&raw)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			pkg.BAD_REQUEST,
		)
		return
	}
	pr, err := h.s.Merged(raw.PullRequestID)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	pkg.OkResponse(ctx, dto.ToPullRequestDTO(pr))
}
