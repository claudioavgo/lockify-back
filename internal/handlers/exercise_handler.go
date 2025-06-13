package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"lockify-back/internal/models"
)

const (
	ExerciseNotFound = "Exercício não encontrado"
	ExerciseCreated  = "Exercício criado com sucesso"
	ExerciseUpdated  = "Exercício atualizado com sucesso"
	ExerciseDeleted  = "Exercício deletado com sucesso"
	Unauthorized     = "Usuário não autenticado"
)

type ExerciseHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewExerciseHandler(db *gorm.DB) *ExerciseHandler {
	return &ExerciseHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
	var exercise models.Exercise
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	if err := h.validate.Struct(exercise); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	exercise.UserID = c.GetUint("user_id")

	if err := h.db.Create(&exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar exercício"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": ExerciseCreated})
}

func (h *ExerciseHandler) GetExercises(c *gin.Context) {
	var exercises []models.Exercise

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
		return
	}

	if err := h.db.Where("user_id = ?", userID).Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar exercícios"})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

func (h *ExerciseHandler) UpdateExercise(c *gin.Context) {
	var exercise models.Exercise

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
		return
	}

	if err := h.db.Model(&exercise).Where("user_id = ?", userID).Updates(exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar exercício"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": ExerciseUpdated})
}

func (h *ExerciseHandler) DeleteExercise(c *gin.Context) {
	var exercise models.Exercise

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	if err := h.db.Model(&exercise).Where("user_id = ?", userID).Delete(&exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar exercício"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": ExerciseDeleted})
}
