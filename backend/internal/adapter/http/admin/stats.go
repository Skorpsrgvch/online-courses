package admin

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	getDetailsUC "github.com/Skorpsrgvch/online-courses/internal/usecase/admin/get_student_details"
	getTopUC "github.com/Skorpsrgvch/online-courses/internal/usecase/admin/get_top_students"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StatsHandler struct {
	getTopUsecase     *getTopUC.Usecase
	getDetailsUsecase *getDetailsUC.Usecase
}

func NewStatsHandler(getTopUsecase *getTopUC.Usecase, getDetailsUsecase *getDetailsUC.Usecase) *StatsHandler {
	return &StatsHandler{
		getTopUsecase:     getTopUsecase,
		getDetailsUsecase: getDetailsUsecase,
	}
}

// GetTopStudents godoc
// @Summary Получить топ студентов
// @Description Возвращает список студентов с наибольшим прогрессом
// @Tags admin
// @Security BearerAuth
// @Param limit query int false "Лимит записей (по умолчанию 20)"
// @Success 200 {object} object{students=[]domain.StudentStat}
// @Failure 403 {object} common.ErrorResponse
// @Router /api/admin/students [get]
func (h *StatsHandler) GetTopStudents(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	input := getTopUC.Input{Limit: limit}
	output, err := h.getTopUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Failed to execute get top students usecase", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetStudentDetails godoc
// @Summary Получить детали студента
// @Description Возвращает информацию о студенте и его прогресс по курсам
// @Tags admin
// @Security BearerAuth
// @Param id path int true "ID студента"
// @Success 200 {object} object{student=domain.StudentStat,courses=[]domain.StudentCourseDetail}
// @Failure 403 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /api/admin/students/:id [get]
func (h *StatsHandler) GetStudentDetails(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		common.HandleError(c, common.HttpError("invalid user ID", http.StatusBadRequest))
		return
	}

	input := getDetailsUC.Input{UserID: userID}
	output, err := h.getDetailsUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Failed to execute get student details usecase", zap.Int("userID", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
