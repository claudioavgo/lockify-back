package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"lockify-back/internal/models"
)

const (
	FeedbackCreated = "Feedback criado com sucesso"
	FeedbackUpdated = "Feedback atualizado com sucesso"
	FeedbackDeleted = "Feedback deletado com sucesso"
)

type FeedbackHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewFeedbackHandler(db *gorm.DB) *FeedbackHandler {
	return &FeedbackHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (h *FeedbackHandler) CreateFeedback(c *gin.Context) {
	var feedback models.Feedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	if err := h.validate.Struct(feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
		return
	}

	feedback.UserID = userID

	if err := h.db.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar feedback"})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

func (h *FeedbackHandler) GetFeedbacks(c *gin.Context) {
	var feedbacks []models.Feedback

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
		return
	}

	if err := h.db.Where("user_id = ?", userID).Find(&feedbacks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar feedbacks"})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}
