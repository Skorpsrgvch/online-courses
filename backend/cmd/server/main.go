package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/auth"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/course"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/lesson"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/module"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/progress"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/review"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/user"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/postgres"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/login"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/register"
	courseCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/createwithmodules"
	courseDelete "github.com/Skorpsrgvch/online-courses/internal/usecase/course/delete"
	courseGet "github.com/Skorpsrgvch/online-courses/internal/usecase/course/get"
	courseGetFull "github.com/Skorpsrgvch/online-courses/internal/usecase/course/getfull"
	courseList "github.com/Skorpsrgvch/online-courses/internal/usecase/course/list"
	courseUpdate "github.com/Skorpsrgvch/online-courses/internal/usecase/course/update"

	lessonCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/create"
	lessonDelete "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/delete"
	lessonGet "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/get"
	lessonUpdate "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/update"

	moduleCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/module/create"
	moduleDelete "github.com/Skorpsrgvch/online-courses/internal/usecase/module/delete"
	moduleGet "github.com/Skorpsrgvch/online-courses/internal/usecase/module/get"
	moduleUpdate "github.com/Skorpsrgvch/online-courses/internal/usecase/module/update"

	forgotPassword "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/forgotpassword"
	resetPassword "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/resetpassword"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/progress/mark"
	reviewApprove "github.com/Skorpsrgvch/online-courses/internal/usecase/review/approve"
	reviewCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/review/create"
	reviewList "github.com/Skorpsrgvch/online-courses/internal/usecase/review/list"
	reviewPending "github.com/Skorpsrgvch/online-courses/internal/usecase/review/pending"
	reviewReject "github.com/Skorpsrgvch/online-courses/internal/usecase/review/reject"
	userCourses "github.com/Skorpsrgvch/online-courses/internal/usecase/user/courses"
	userProfile "github.com/Skorpsrgvch/online-courses/internal/usecase/user/profile"
	"github.com/Skorpsrgvch/online-courses/pkg/db"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://courses:securepass@localhost:5432/courses?sslmode=disable"
	}

	// Ожидание подключения к БД
	if err := waitForDB(dbURL); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Подключение к БД
	dbConn, err := db.NewPostgresDB(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	// === Репозитории ===
	userRepo := postgres.NewUserRepo(dbConn)
	courseRepo := postgres.NewCourseRepo(dbConn)
	courseTxRepo := postgres.NewCourseTxRepo(dbConn)
	moduleRepo := postgres.NewModuleRepo(dbConn)
	lessonRepo := postgres.NewLessonRepo(dbConn)
	progressRepo := postgres.NewProgressRepo(dbConn)
	reviewRepo := postgres.NewReviewRepo(dbConn)
	purchaseRepo := postgres.NewPurchaseRepo(dbConn)
	resetTokenRepo := postgres.NewResetTokenRepo(dbConn)

	// === Юзкейсы ===
	// Аутентификация
	registerUC, _ := register.NewUsecase(userRepo)
	loginUC, _ := login.NewUsecase(userRepo)

	// Курсы
	createCourseUC, _ := courseCreate.NewUsecase(courseRepo)
	createWithModulesUC, _ := createwithmodules.NewUsecase(courseTxRepo)
	getCourseUC, _ := courseGet.NewUsecase(courseRepo, purchaseRepo)
	getFullCourseUC, _ := courseGetFull.NewUsecase(courseRepo, moduleRepo, lessonRepo, purchaseRepo)
	updateCourseUC, _ := courseUpdate.NewUsecase(courseRepo, courseRepo)
	deleteCourseUC, _ := courseDelete.NewUsecase(courseRepo)
	listCourseUC, _ := courseList.NewUsecase(courseRepo)

	// Модули
	createModuleUC, _ := moduleCreate.NewUsecase(moduleRepo)
	getModuleUC, _ := moduleGet.NewUsecase(moduleRepo)
	updateModuleUC, _ := moduleUpdate.NewUsecase(moduleRepo, moduleRepo)
	deleteModuleUC, _ := moduleDelete.NewUsecase(moduleRepo)

	// Уроки
	createLessonUC, _ := lessonCreate.NewUsecase(lessonRepo)
	getLessonUC, _ := lessonGet.NewUsecase(lessonRepo)
	updateLessonUC, _ := lessonUpdate.NewUsecase(lessonRepo, lessonRepo)
	deleteLessonUC, _ := lessonDelete.NewUsecase(lessonRepo)

	// Прогресс
	markProgressUC, _ := mark.NewUsecase(progressRepo)

	// Отзывы
	createReviewUC, _ := reviewCreate.NewUsecase(reviewRepo, userRepo, courseRepo)
	approveReviewUC, _ := reviewApprove.NewUsecase(reviewRepo)
	listReviewUC, _ := reviewList.NewUsecase(reviewRepo)
	pendingReviewUC, _ := reviewPending.NewUsecase(reviewRepo)
	rejectReviewUC, _ := reviewReject.NewUsecase(reviewRepo)
	userProfileUC, _ := userProfile.NewUsecase(userRepo)
	userCoursesUC, _ := userCourses.NewUsecase(courseRepo, purchaseRepo, progressRepo)
	forgotPasswordUC, _ := forgotPassword.NewUsecase(userRepo, resetTokenRepo)
	resetPasswordUC, _ := resetPassword.NewUsecase(resetTokenRepo, userRepo)

	// === Хендлеры ===
	// Аутентификация
	registerHandler := auth.NewRegisterHandler(registerUC)
	loginHandler := auth.NewLoginHandler(loginUC)
	forgotPasswordHandler := auth.NewForgotPasswordHandler(forgotPasswordUC)
	resetPasswordHandler := auth.NewResetPasswordHandler(resetPasswordUC)

	// Курсы
	createCourseHandler := course.NewCreateHandler(createCourseUC)
	createWithModulesHandler := course.NewCreateWithModulesHandler(createWithModulesUC)
	getCourseHandler := course.NewGetHandler(getCourseUC)
	getFullCourseHandler := course.NewGetFullHandler(getFullCourseUC)
	updateCourseHandler := course.NewUpdateHandler(updateCourseUC)
	deleteCourseHandler := course.NewDeleteHandler(deleteCourseUC)
	listCourseHandler := course.NewListHandler(listCourseUC)

	// Модули
	createModuleHandler := module.NewCreateHandler(createModuleUC)
	getModuleHandler := module.NewGetHandler(getModuleUC)
	updateModuleHandler := module.NewUpdateHandler(updateModuleUC)
	deleteModuleHandler := module.NewDeleteHandler(deleteModuleUC)

	// Уроки
	createLessonHandler := lesson.NewCreateHandler(createLessonUC)
	getLessonHandler := lesson.NewGetHandler(getLessonUC)
	updateLessonHandler := lesson.NewUpdateHandler(updateLessonUC)
	deleteLessonHandler := lesson.NewDeleteHandler(deleteLessonUC)

	// Прогресс
	markProgressHandler := progress.NewMarkHandler(markProgressUC)

	// Отзывы
	createReviewHandler := review.NewCreateHandler(createReviewUC)
	approveReviewHandler := review.NewApproveHandler(approveReviewUC)
	listReviewHandler := review.NewListHandler(listReviewUC)
	pendingReviewHandler := review.NewPendingHandler(pendingReviewUC)
	rejectReviewHandler := review.NewRejectHandler(rejectReviewUC)

	// Пользователь
	userProfileHandler := user.NewProfileHandler(userProfileUC)
	userCoursesHandler := user.NewCoursesHandler(userCoursesUC)

	// === Роутер ===
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS для фронтенда
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Публичные эндпоинты
	r.POST("/api/auth/register", registerHandler.Handle)
	r.POST("/api/auth/login", loginHandler.Handle)
	r.POST("/api/auth/forgot-password", forgotPasswordHandler.Handle)
	r.POST("/api/auth/reset-password", resetPasswordHandler.Handle)

	// Защищённые эндпоинты
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Auth
		api.GET("/auth/me", auth.NewMeHandler().Handle)
		api.POST("/auth/refresh", auth.NewRefreshHandler().Handle)

		// Курсы
		api.GET("/courses", listCourseHandler.Handle)
		api.POST("/courses", createCourseHandler.Handle)
		api.POST("/courses/with-modules", createWithModulesHandler.Handle)
		api.GET("/courses/:id", getCourseHandler.Handle)
		api.GET("/courses/:id/full", getFullCourseHandler.Handle)
		api.PUT("/courses/:id", updateCourseHandler.Handle)
		api.DELETE("/courses/:id", deleteCourseHandler.Handle)

		// Модули
		api.POST("/modules", createModuleHandler.Handle)
		api.GET("/courses/:id/modules", getModuleHandler.Handle)
		api.PUT("/modules/:id", updateModuleHandler.Handle)
		api.DELETE("/modules/:id", deleteModuleHandler.Handle)

		// Уроки
		api.POST("/lessons", createLessonHandler.Handle)
		api.GET("/modules/:id/lessons", getLessonHandler.Handle)
		api.PUT("/lessons/:id", updateLessonHandler.Handle)
		api.DELETE("/lessons/:id", deleteLessonHandler.Handle)

		// Прогресс
		api.POST("/progress/lessons/:id/mark", markProgressHandler.Handle)

		// Отзывы
		api.POST("/reviews", createReviewHandler.Handle)
		api.POST("/reviews/:id/approve", approveReviewHandler.Handle)
		api.GET("/courses/:id/reviews", listReviewHandler.Handle)

		// Модерация отзывов (только admin)
		api.GET("/reviews/admin/pending", pendingReviewHandler.Handle)
		api.DELETE("/reviews/:id", rejectReviewHandler.Handle)

		// Профиль пользователя
		api.GET("/user/profile", userProfileHandler.Handle)
		api.GET("/user/courses", userCoursesHandler.Handle)
	}

	// Запуск сервера
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}

// waitForDB — ожидание подключения к БД
func waitForDB(dbURL string) error {
	for i := 0; i < 30; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		dbConn, err := db.NewPostgresDB(ctx, dbURL)
		cancel()
		if err != nil {
			log.Printf("Waiting for DB... attempt %d", i+1)
			time.Sleep(2 * time.Second)
			continue
		}
		dbConn.Close()
		return nil
	}
	return fmt.Errorf("timeout waiting for database")
}
