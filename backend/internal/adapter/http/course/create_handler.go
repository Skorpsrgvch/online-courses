package course

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	createUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createCourseRequest struct {
	Title             string             `json:"title" binding:"required"`
	Description       string             `json:"description"`
	IsPublic          bool               `json:"is_public"`
	Price             int                `json:"price"`
	CoverImageURL     string             `json:"cover_image_url"`
	Contraindications string             `json:"contraindications"`
	Recommendations   string             `json:"recommendations"`
	TargetAudience    string             `json:"target_audience"`
	CourseBasis       string             `json:"course_basis"`
	ClassBasis        string             `json:"class_basis"`
	Bonuses           []domain.BonusItem `json:"bonuses"`
}

type CreateHandler struct {
	usecase *createUC.Usecase
}

func NewCreateHandler(usecase *createUC.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req createCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	userID := middleware.GetUserID(c)
	zap.L().Info("Create course request", zap.Int("authorID", userID), zap.String("title", req.Title))

	input := createUC.Input{
		Title:             req.Title,
		Description:       req.Description,
		IsPublic:          req.IsPublic,
		Price:             req.Price,
		AuthorID:          userID,
		CoverImageURL:     req.CoverImageURL,
		Contraindications: req.Contraindications,
		Recommendations:   req.Recommendations,
		TargetAudience:    req.TargetAudience,
		CourseBasis:       req.CourseBasis,
		ClassBasis:        req.ClassBasis,
		Bonuses:           req.Bonuses,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Create course failed", zap.Int("authorID", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Course created successfully", zap.Int("authorID", userID), zap.String("title", req.Title))
	c.Status(http.StatusCreated)
}
