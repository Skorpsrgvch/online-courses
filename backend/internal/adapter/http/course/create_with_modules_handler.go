package course

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	createUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
	createWithModulesUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/createwithmodules"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type moduleRequest struct {
	Title   string          `json:"title" binding:"required"`
	Order   int             `json:"order"`
	Lessons []lessonRequest `json:"lessons"`
}

type lessonRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type createCourseWithModulesRequest struct {
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
	Modules           []moduleRequest    `json:"modules"`
}

type CreateWithModulesHandler struct {
	usecase *createWithModulesUC.Usecase
}

func NewCreateWithModulesHandler(usecase *createWithModulesUC.Usecase) *CreateWithModulesHandler {
	return &CreateWithModulesHandler{usecase: usecase}
}

func (h *CreateWithModulesHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req createCourseWithModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	userID := middleware.GetUserID(c)
	zap.L().Info("Create course with modules request", zap.Int("authorID", userID), zap.String("title", req.Title), zap.Int("modulesCount", len(req.Modules)))

	modules := make([]createUC.ModuleInput, len(req.Modules))
	for i, mod := range req.Modules {
		lessons := make([]createUC.LessonInput, len(mod.Lessons))
		for j, l := range mod.Lessons {
			lessons[j] = createUC.LessonInput{
				Title:        l.Title,
				Description:  l.Description,
				VideoEmbedID: l.VideoEmbedID,
				PrivateKey:   l.PrivateKey,
				Order:        l.Order,
			}
		}
		modules[i] = createUC.ModuleInput{
			Title:   mod.Title,
			Order:   mod.Order,
			Lessons: lessons,
		}
	}

	input := createWithModulesUC.Input{
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
		Modules:           modules,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Create course with modules failed", zap.Int("authorID", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Course with modules created successfully", zap.Int("authorID", userID), zap.String("title", req.Title))
	c.Status(http.StatusCreated)
}
