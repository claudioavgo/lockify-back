package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"lockify-back/internal/config"
	"lockify-back/internal/handlers"
	"lockify-back/internal/middleware"
)

func SetupRoutes(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// CORS
	allowedOrigin := cfg.AllowedOrigin

	fmt.Println("[ALLOWED_ORIGIN]", allowedOrigin)

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Handlers
	userHandler := handlers.NewUserHandler(db, cfg)
	habitHandler := handlers.NewHabitHandler(db)
	exerciseHandler := handlers.NewExerciseHandler(db)
	feedbackHandler := handlers.NewFeedbackHandler(db)

	// Rotas
	api := router.Group("/api")
	{
		// Rotas públicas
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.CreateUser)
		}

		// Rotas protegidas
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("", userHandler.GetUsers)
		}

		me := api.Group("/me")
		me.Use(middleware.AuthMiddleware())
		{
			me.GET("", userHandler.GetMe)
		}

		habits := api.Group("/habits")
		habits.Use(middleware.AuthMiddleware())
		{
			habits.GET("", habitHandler.GetHabits)
			habits.POST("", habitHandler.CreateHabit)
			habits.GET("/day", habitHandler.GetDayHabits)
			habits.POST("/:habit_id/check-in", habitHandler.CheckInHabit)
		}

		exercises := api.Group("/exercises")
		exercises.Use(middleware.AuthMiddleware())
		{
			exercises.GET("", exerciseHandler.GetExercises)
			exercises.POST("", exerciseHandler.CreateExercise)
			exercises.PUT("/:id", exerciseHandler.UpdateExercise)
			exercises.DELETE("/:id", exerciseHandler.DeleteExercise)
		}

		feedbacks := api.Group("/feedbacks")
		feedbacks.Use(middleware.AuthMiddleware())
		{
			feedbacks.GET("", feedbackHandler.GetFeedbacks)
			feedbacks.POST("", feedbackHandler.CreateFeedback)
		}
	}

	return router
}
