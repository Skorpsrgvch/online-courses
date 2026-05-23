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

	"github.com/Skorpsrgvch/online-courses/internal/adapter/email"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/admin"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/auth"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/course"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/lesson"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/module"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/payment"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/payment/yookassa"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/progress"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/review"
	serviceHttp "github.com/Skorpsrgvch/online-courses/internal/adapter/http/service"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/user"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/postgres"
	"github.com/Skorpsrgvch/online-courses/pkg/db"

	// Auth UseCases
	forgotPassword "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/forgotpassword"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/login"
	logout "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/logout"
	refresh "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/refresh"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/register"
	resetPassword "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/resetpassword"

	// Admin UseCases
	access "github.com/Skorpsrgvch/online-courses/internal/usecase/admin/access"
	search "github.com/Skorpsrgvch/online-courses/internal/usecase/admin/search"

	// Course UseCases
	courseCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/createwithmodules"
	courseDelete "github.com/Skorpsrgvch/online-courses/internal/usecase/course/delete"
	enroll_free "github.com/Skorpsrgvch/online-courses/internal/usecase/course/enroll_free"
	courseGetFull "github.com/Skorpsrgvch/online-courses/internal/usecase/course/getfull"
	courseList "github.com/Skorpsrgvch/online-courses/internal/usecase/course/list"
	courseUpdate "github.com/Skorpsrgvch/online-courses/internal/usecase/course/update"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/updatefullcourse"

	// Lesson UseCases
	lessonCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/create"
	lessonDelete "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/delete"
	lessonGet "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/get"
	lessonUpdate "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/update"

	// Module UseCases
	moduleCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/module/create"
	moduleDelete "github.com/Skorpsrgvch/online-courses/internal/usecase/module/delete"
	moduleGet "github.com/Skorpsrgvch/online-courses/internal/usecase/module/get"
	moduleUpdate "github.com/Skorpsrgvch/online-courses/internal/usecase/module/update"

	// Payment UseCases
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/callback"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/confirm"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/create"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/list"

	// Progress UseCase
	"github.com/Skorpsrgvch/online-courses/internal/usecase/progress/mark"

	// Review UseCases
	admin_delete "github.com/Skorpsrgvch/online-courses/internal/usecase/review/admin_delete"
	reviewApprove "github.com/Skorpsrgvch/online-courses/internal/usecase/review/approve"
	reviewCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/review/create"
	reviewDelete "github.com/Skorpsrgvch/online-courses/internal/usecase/review/delete"
	reviewList "github.com/Skorpsrgvch/online-courses/internal/usecase/review/list"
	myReviews "github.com/Skorpsrgvch/online-courses/internal/usecase/review/myreviews"
	reviewPending "github.com/Skorpsrgvch/online-courses/internal/usecase/review/pending"
	reviewReject "github.com/Skorpsrgvch/online-courses/internal/usecase/review/reject"

	// Service UseCases
	serviceCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/service/create"
	serviceDelete "github.com/Skorpsrgvch/online-courses/internal/usecase/service/delete"
	serviceGet "github.com/Skorpsrgvch/online-courses/internal/usecase/service/get"
	serviceList "github.com/Skorpsrgvch/online-courses/internal/usecase/service/list"
	serviceUpdate "github.com/Skorpsrgvch/online-courses/internal/usecase/service/update"

	// User UseCases
	changePassword "github.com/Skorpsrgvch/online-courses/internal/usecase/user/change_password"
	userCourses "github.com/Skorpsrgvch/online-courses/internal/usecase/user/courses"
	userProfile "github.com/Skorpsrgvch/online-courses/internal/usecase/user/profile"
	updateProfile "github.com/Skorpsrgvch/online-courses/internal/usecase/user/update_profile"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func main() {
	// Инициализация глобального логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	zap.L().Info("Starting application...")

	gin.SetMode(gin.ReleaseMode)

	// === Конфигурация ===
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		zap.L().Fatal("DB_URL must be set")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpHost == "" || smtpPort == "" {
		zap.L().Fatal("SMTP_HOST and SMTP_PORT must be set")
	}

	emailSender := email.NewSMTPSender()

	yookassaConfig := yookassa.Config{
		ShopID:    os.Getenv("YOOKASSA_SHOP_ID"),
		SecretKey: os.Getenv("YOOKASSA_SECRET_KEY"),
		BaseURL:   os.Getenv("YOOKASSA_BASE_URL"),
	}
	if yookassaConfig.BaseURL == "" {
		yookassaConfig.BaseURL = "https://api.yookassa.ru/v3 "
	}

	if yookassaConfig.ShopID == "" || yookassaConfig.SecretKey == "" {
		zap.L().Fatal("YOOKASSA_SHOP_ID and YOOKASSA_SECRET_KEY must be set")
	}

	paymentGateway := yookassa.NewGateway(yookassaConfig)

	// === Подключение к БД ===
	if err := waitForDB(dbURL); err != nil {
		zap.L().Fatal("Failed to connect to DB", zap.Error(err))
	}

	dbConn, err := db.NewPostgresDB(context.Background(), dbURL)
	if err != nil {
		zap.L().Fatal("Failed to create DB connection", zap.Error(err))
	}
	defer dbConn.Close()

	// === Репозитории ===
	userRepo := postgres.NewUserRepo(dbConn)
	courseRepo := postgres.NewCourseRepo(dbConn)
	courseTxRepo := postgres.NewCourseTxRepo(dbConn)
	courseFullRepo := postgres.NewCourseFullRepo(dbConn)
	moduleRepo := postgres.NewModuleRepo(dbConn)
	lessonRepo := postgres.NewLessonRepo(dbConn)
	progressRepo := postgres.NewProgressRepo(dbConn)
	reviewRepo := postgres.NewReviewRepo(dbConn)
	serviceRepo := postgres.NewServiceRepo(dbConn)
	purchaseRepo := postgres.NewPurchaseRepo(dbConn)
	resetTokenRepo := postgres.NewResetTokenRepo(dbConn)
	paymentRepo := postgres.NewPaymentRepository(dbConn)
	refreshTokenRepo := postgres.NewRefreshTokenRepo(dbConn)

	// === Инициализация UseCases ===
	// Теперь мы НЕ передаем логгер в конструкторы. Логгер берется внутри через zap.L().
	// Проверка на ошибку остается только там, где конструктор может вернуть ошибку (валидация зависимостей).

	// Auth
	registerUC, err := register.NewUsecase(userRepo)
	if err != nil {
		zap.L().Fatal("Failed to create register usecase", zap.Error(err))
	}

	loginUC, err := login.NewUsecase(userRepo)
	if err != nil {
		zap.L().Fatal("Failed to create login usecase", zap.Error(err))
	}

	refreshUC, err := refresh.NewUsecase(refreshTokenRepo, userRepo)
	if err != nil {
		zap.L().Fatal("Failed to create refresh usecase", zap.Error(err))
	}

	logoutUC, err := logout.NewUsecase(refreshTokenRepo)
	if err != nil {
		zap.L().Fatal("Failed to create logout usecase", zap.Error(err))
	}

	forgotPasswordUC, err := forgotPassword.NewUsecase(userRepo, resetTokenRepo, emailSender)
	if err != nil {
		zap.L().Fatal("Failed to create forgot password usecase", zap.Error(err))
	}

	resetPasswordUC, err := resetPassword.NewUsecase(resetTokenRepo, userRepo)
	if err != nil {
		zap.L().Fatal("Failed to create reset password usecase", zap.Error(err))
	}

	// Courses
	createCourseUC, err := courseCreate.NewUsecase(courseRepo)
	if err != nil {
		zap.L().Fatal("Failed to create course usecase", zap.Error(err))
	}

	createWithModulesUC, err := createwithmodules.NewUsecase(courseTxRepo)
	if err != nil {
		zap.L().Fatal("Failed to create course with modules usecase", zap.Error(err))
	}

	getFullCourseUC, err := courseGetFull.NewUsecase(courseRepo, moduleRepo, lessonRepo, purchaseRepo, progressRepo)
	if err != nil {
		zap.L().Fatal("Failed to create get full course usecase", zap.Error(err))
	}

	updateCourseUC, err := courseUpdate.NewUsecase(courseRepo)
	if err != nil {
		zap.L().Fatal("Failed to create update course usecase", zap.Error(err))
	}

	deleteCourseUC, err := courseDelete.NewUsecase(courseRepo)
	if err != nil {
		zap.L().Fatal("Failed to create delete course usecase", zap.Error(err))
	}

	listCourseUC, err := courseList.NewUsecase(courseRepo)
	if err != nil {
		zap.L().Fatal("Failed to create list course usecase", zap.Error(err))
	}

	// UpdateFullCourse может не возвращать ошибку, если конструктор простой, но лучше проверить сигнатуру.
	// Если конструктор возвращает (*Usecase, error), используем этот паттерн.
	// Если просто *Usecase, то присваиваем напрямую.
	// Предполагаем, что мы привели все конструкторы к виду (deps) (*Usecase, error) для единообразия.
	updateFullCourseUC, err := updatefullcourse.NewUsecase(courseFullRepo)
	if err != nil {
		zap.L().Fatal("Failed to create update full course usecase", zap.Error(err))
	}

	// Modules
	createModuleUC, err := moduleCreate.NewUsecase(moduleRepo)
	if err != nil {
		zap.L().Fatal("Failed to create module usecase", zap.Error(err))
	}

	getModuleUC, err := moduleGet.NewUsecase(moduleRepo)
	if err != nil {
		zap.L().Fatal("Failed to create get module usecase", zap.Error(err))
	}

	updateModuleUC, err := moduleUpdate.NewUsecase(moduleRepo, moduleRepo)
	if err != nil {
		zap.L().Fatal("Failed to create update module usecase", zap.Error(err))
	}

	deleteModuleUC, err := moduleDelete.NewUsecase(moduleRepo)
	if err != nil {
		zap.L().Fatal("Failed to create delete module usecase", zap.Error(err))
	}

	// Lessons
	createLessonUC, err := lessonCreate.NewUsecase(lessonRepo)
	if err != nil {
		zap.L().Fatal("Failed to create lesson usecase", zap.Error(err))
	}

	getLessonUC, err := lessonGet.NewUsecase(lessonRepo)
	if err != nil {
		zap.L().Fatal("Failed to create get lesson usecase", zap.Error(err))
	}

	updateLessonUC, err := lessonUpdate.NewUsecase(lessonRepo, lessonRepo)
	if err != nil {
		zap.L().Fatal("Failed to create update lesson usecase", zap.Error(err))
	}

	deleteLessonUC, err := lessonDelete.NewUsecase(lessonRepo)
	if err != nil {
		zap.L().Fatal("Failed to create delete lesson usecase", zap.Error(err))
	}

	// Services
	getServiceUC, err := serviceGet.NewUsecase(serviceRepo)
	if err != nil {
		zap.L().Fatal("Failed to create get service usecase", zap.Error(err))
	}

	listServiceUC, err := serviceList.NewUsecase(serviceRepo)
	if err != nil {
		zap.L().Fatal("Failed to create list service usecase", zap.Error(err))
	}

	createServiceUC, err := serviceCreate.NewUsecase(serviceRepo)
	if err != nil {
		zap.L().Fatal("Failed to create service usecase", zap.Error(err))
	}

	updateServiceUC, err := serviceUpdate.NewUsecase(serviceRepo)
	if err != nil {
		zap.L().Fatal("Failed to create update service usecase", zap.Error(err))
	}

	deleteServiceUC, err := serviceDelete.NewUsecase(serviceRepo)
	if err != nil {
		zap.L().Fatal("Failed to create delete service usecase", zap.Error(err))
	}

	// Progress
	markProgressUC, err := mark.NewUsecase(progressRepo)
	if err != nil {
		zap.L().Fatal("Failed to create progress usecase", zap.Error(err))
	}

	// Reviews
	createReviewUC, err := reviewCreate.NewUsecase(reviewRepo, userRepo, courseRepo)
	if err != nil {
		zap.L().Fatal("Failed to create review usecase", zap.Error(err))
	}

	approveReviewUC, err := reviewApprove.NewUsecase(reviewRepo)
	if err != nil {
		zap.L().Fatal("Failed to create approve review usecase", zap.Error(err))
	}

	listReviewUC, err := reviewList.NewUsecase(reviewRepo)
	if err != nil {
		zap.L().Fatal("Failed to create list review usecase", zap.Error(err))
	}

	pendingReviewUC, err := reviewPending.NewUsecase(reviewRepo)
	if err != nil {
		zap.L().Fatal("Failed to create pending review usecase", zap.Error(err))
	}

	rejectReviewUC, err := reviewReject.NewUsecase(reviewRepo)
	if err != nil {
		zap.L().Fatal("Failed to create reject review usecase", zap.Error(err))
	}

	myReviewsUC, err := myReviews.NewUsecase(reviewRepo)
	if err != nil {
		zap.L().Fatal("Failed to create my reviews usecase", zap.Error(err))
	}

	deleteReviewsUC, err := reviewDelete.NewUsecase(reviewRepo)
	if err != nil {
		zap.L().Fatal("Failed to create delete review usecase", zap.Error(err))
	}

	adminDeleteUC, err := admin_delete.NewUsecase(reviewRepo)
	if err != nil {
		zap.L().Fatal("Failed to create admin delete review usecase", zap.Error(err))
	}

	// Users
	userProfileUC, err := userProfile.NewUsecase(userRepo)
	if err != nil {
		zap.L().Fatal("Failed to create user profile usecase", zap.Error(err))
	}

	userCoursesUC, err := userCourses.NewUsecase(courseRepo, purchaseRepo, progressRepo)
	if err != nil {
		zap.L().Fatal("Failed to create user courses usecase", zap.Error(err))
	}

	changePasswordUC, err := changePassword.NewUsecase(userRepo)
	if err != nil {
		zap.L().Fatal("Failed to create change password usecase", zap.Error(err))
	}

	updateProfileUC, err := updateProfile.NewUsecase(userRepo)
	if err != nil {
		zap.L().Fatal("Failed to create update profile usecase", zap.Error(err))
	}

	// Admin
	searchUsersUC, err := search.NewUsecase(userRepo)
	if err != nil {
		zap.L().Fatal("Failed to create search users usecase", zap.Error(err))
	}

	accessUC, err := access.NewUsecase(userRepo, courseRepo, purchaseRepo)
	if err != nil {
		zap.L().Fatal("Failed to create access usecase", zap.Error(err))
	}

	// Payments
	confirmPaymentUC := confirm.NewUseCase(paymentRepo, purchaseRepo)
	callbackPaymentUC := callback.NewUseCase(paymentRepo, purchaseRepo)
	listPaymentUC := list.NewUseCase(paymentRepo)
	createPaymentUC := create.NewUseCase(paymentRepo, courseRepo, purchaseRepo, paymentGateway)

	enrollFreeUC, err := enroll_free.NewUsecase(purchaseRepo)
	if err != nil {
		zap.L().Fatal("Failed to create enroll free usecase", zap.Error(err))
	}

	// === Хендлеры ===
	registerHandler := auth.NewRegisterHandler(registerUC)
	loginHandler := auth.NewLoginHandler(loginUC)
	forgotPasswordHandler := auth.NewForgotPasswordHandler(forgotPasswordUC)
	resetPasswordHandler := auth.NewResetPasswordHandler(resetPasswordUC)
	refreshHandler := auth.NewRefreshHandler(refreshUC)
	logoutHandler := auth.NewLogoutHandler(logoutUC)

	createCourseHandler := course.NewCreateHandler(createCourseUC)
	createWithModulesHandler := course.NewCreateWithModulesHandler(createWithModulesUC)
	getFullCourseHandler := course.NewGetFullHandler(getFullCourseUC)
	updateCourseHandler := course.NewUpdateHandler(updateCourseUC)
	deleteCourseHandler := course.NewDeleteHandler(deleteCourseUC)
	listCourseHandler := course.NewListHandler(listCourseUC)
	updateFullCourseHandler := course.NewUpdateFullCourseHandler(updateFullCourseUC)
	reorderHandler := course.NewReorderHandler(updateFullCourseUC)

	getServiceHandler := serviceHttp.NewGetByIDHandler(getServiceUC)
	listServiceHandler := serviceHttp.NewListHandler(listServiceUC)
	createServiceHandler := serviceHttp.NewCreateHandler(createServiceUC)
	updateServiceHandler := serviceHttp.NewUpdateHandler(updateServiceUC)
	deleteServiceHandler := serviceHttp.NewDeleteHandler(deleteServiceUC)

	createModuleHandler := module.NewCreateHandler(createModuleUC)
	getModuleHandler := module.NewGetHandler(getModuleUC)
	updateModuleHandler := module.NewUpdateHandler(updateModuleUC)
	deleteModuleHandler := module.NewDeleteHandler(deleteModuleUC)

	createLessonHandler := lesson.NewCreateHandler(createLessonUC)
	getLessonHandler := lesson.NewGetHandler(getLessonUC)
	updateLessonHandler := lesson.NewUpdateHandler(updateLessonUC)
	deleteLessonHandler := lesson.NewDeleteHandler(deleteLessonUC)

	markProgressHandler := progress.NewMarkHandler(markProgressUC)

	createReviewHandler := review.NewCreateHandler(createReviewUC)
	approveReviewHandler := review.NewApproveHandler(approveReviewUC)
	listReviewHandler := review.NewListHandler(listReviewUC)
	pendingReviewHandler := review.NewPendingHandler(pendingReviewUC)
	rejectReviewHandler := review.NewRejectHandler(rejectReviewUC)
	myReviewsHandler := review.NewMyReviewsHandler(myReviewsUC)
	deleteReviewHandler := review.NewDeleteHandler(deleteReviewsUC)
	adminDeleteHandler := review.NewAdminDeleteHandler(adminDeleteUC)

	userProfileHandler := user.NewProfileHandler(userProfileUC)
	userCoursesHandler := user.NewCoursesHandler(userCoursesUC)
	changePasswordHandler := user.NewChangePasswordHandler(changePasswordUC)
	updateProfileHandler := user.NewUpdateProfileHandler(updateProfileUC)

	createPaymentHandler := payment.NewCreateHandler(createPaymentUC)
	confirmPaymentHandler := payment.NewConfirmHandler(confirmPaymentUC)
	callbackPaymentHandler := payment.NewCallbackHandler(callbackPaymentUC)
	listPaymentHandler := payment.NewListHandler(listPaymentUC)
	enrollFreeHandler := course.NewEnrollFreeHandler(enrollFreeUC)

	searchUsersHandler := admin.NewSearchUsersHandler(searchUsersUC)
	accessHandler := admin.NewGrantAccessHandler(accessUC)

	// === Роутер ===
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", frontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public Auth
	authGroup := r.Group("/api/auth")
	authGroup.Use(middleware.RateLimitMiddleware(rate.Every(1*time.Second), 5))
	{
		authGroup.POST("/register", registerHandler.Handle)
		authGroup.POST("/login", loginHandler.Handle)
		authGroup.POST("/forgot-password", forgotPasswordHandler.Handle)
		authGroup.POST("/reset-password", resetPasswordHandler.Handle)
		authGroup.POST("/refresh", refreshHandler.Handle)
	}

	// Public Content
	publicApi := r.Group("/api")
	{
		publicApi.GET("/courses", listCourseHandler.Handle)
		publicApi.GET("/services", listServiceHandler.Handle)
		publicApi.GET("/payments/:payment_id/confirm", confirmPaymentHandler.HandleGet)
		publicApi.POST("/payments/callback", callbackPaymentHandler.Handle)
		publicApi.GET("/courses/:id/modules", getModuleHandler.Handle)
		publicApi.GET("/courses/:id/reviews", listReviewHandler.Handle)
	}

	// Protected API
	protectedApi := r.Group("/api")
	protectedApi.Use(middleware.AuthMiddleware())
	{
		protectedApi.GET("/auth/me", auth.NewMeHandler().Handle)
		protectedApi.POST("/auth/logout", logoutHandler.Handle)

		protectedApi.GET("/user/profile", userProfileHandler.Handle)
		protectedApi.PUT("/user/profile", updateProfileHandler.Handle)
		protectedApi.PUT("/user/password", changePasswordHandler.Handle)
		protectedApi.GET("/user/courses", userCoursesHandler.Handle)

		protectedApi.POST("/users/search", searchUsersHandler.Handle)
		protectedApi.POST("/access", accessHandler.Handle)

		protectedApi.GET("/courses/:id/full", getFullCourseHandler.Handle)

		protectedApi.POST("/services", createServiceHandler.Handle)
		protectedApi.PUT("/services/:id", updateServiceHandler.Handle)
		protectedApi.DELETE("/services/:id", deleteServiceHandler.Handle)
		protectedApi.GET("/services/:id", getServiceHandler.Handle)

		protectedApi.POST("/courses", createCourseHandler.Handle)
		protectedApi.POST("/courses/with-modules", createWithModulesHandler.Handle)
		protectedApi.PUT("/courses/:id", updateCourseHandler.Handle)
		protectedApi.PUT("/courses/:id/full-update", updateFullCourseHandler.Handle)
		protectedApi.PATCH("/courses/:id/status", updateCourseHandler.HandleStatusPatch)
		protectedApi.DELETE("/courses/:id", deleteCourseHandler.Handle)
		protectedApi.POST("/courses/:id/enroll-free", enrollFreeHandler.Handle)
		protectedApi.GET("/admin/courses/all", listCourseHandler.HandleAdmin)

		protectedApi.POST("/modules", createModuleHandler.Handle)
		protectedApi.GET("/modules/:id/lessons", getLessonHandler.Handle)
		protectedApi.PUT("/modules/:id", updateModuleHandler.Handle)
		protectedApi.DELETE("/modules/:id", deleteModuleHandler.Handle)

		protectedApi.PUT("/courses/:id/modules/reorder", reorderHandler.HandleModules)
		protectedApi.PUT("/modules/:id/lessons/reorder", reorderHandler.HandleLessons)

		protectedApi.POST("/lessons", createLessonHandler.Handle)
		protectedApi.PUT("/lessons/:id", updateLessonHandler.Handle)
		protectedApi.DELETE("/lessons/:id", deleteLessonHandler.Handle)

		protectedApi.POST("/progress/lessons/:id/mark", markProgressHandler.Handle)
		protectedApi.POST("/reviews", createReviewHandler.Handle)
		protectedApi.POST("/reviews/:id/approve", approveReviewHandler.Handle)
		protectedApi.GET("/reviews/admin/pending", pendingReviewHandler.Handle)
		protectedApi.PUT("/reviews/:id", rejectReviewHandler.Handle)
		protectedApi.DELETE("/reviews/:id", deleteReviewHandler.Handle)
		protectedApi.DELETE("/admin/reviews/:id", adminDeleteHandler.Handle)
		protectedApi.GET("/reviews/my", myReviewsHandler.Handle)

		protectedApi.POST("/payments", createPaymentHandler.Handle)
		protectedApi.GET("/payments", listPaymentHandler.Handle)
		protectedApi.POST("/payments/:payment_id/confirm", confirmPaymentHandler.Handle)
	}

	// Start Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		zap.L().Info("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}
	zap.L().Info("Server exited gracefully")
}

func waitForDB(dbURL string) error {
	for i := 0; i < 30; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		dbConn, err := db.NewPostgresDB(ctx, dbURL)
		cancel()
		if err != nil {
			zap.L().Warn("Waiting for DB...", zap.Int("attempt", i+1), zap.Error(err))
			time.Sleep(2 * time.Second)
			continue
		}
		dbConn.Close()
		return nil
	}
	return fmt.Errorf("timeout waiting for database")
}
