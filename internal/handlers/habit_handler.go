package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"lockify-back/internal/models"
)

const (
	HabitCreated   = "Hábito criado com sucesso"
	HabitUpdated   = "Hábito atualizado com sucesso"
	HabitDeleted   = "Hábito deletado com sucesso"
	HabitCompleted = "Hábito completado com sucesso"
)

type HabitHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewHabitHandler(db *gorm.DB) *HabitHandler {
	return &HabitHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (h *HabitHandler) CreateHabit(c *gin.Context) {
	var habit models.Habit
	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	if err := h.validate.Struct(habit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	habit.UserID = c.GetUint("user_id")

	if err := h.db.Create(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar hábito"})
		return
	}

	c.JSON(http.StatusCreated, habit)
}

func (h *HabitHandler) GetHabits(c *gin.Context) {
	var habits []models.Habit = make([]models.Habit, 0)

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
		return
	}

	if err := h.db.Where("user_id = ?", userID).Find(&habits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar hábitos"})
		return
	}

	c.JSON(http.StatusOK, habits)
}

func (h *HabitHandler) CheckInHabit(c *gin.Context) {
	var habit models.Habit
	var existingCheckIn models.HabitDayCheckIn

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
		return
	}

	habitID, err := strconv.ParseUint(c.Param("habit_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do hábito inválido"})
		return
	}

	if err := h.db.Where("user_id = ? AND id = ?", userID, habitID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hábito não encontrado"})
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	if err := h.db.Where("user_id = ? AND habit_id = ? AND date >= ?",
		userID, habitID, startOfDay, endOfDay).First(&existingCheckIn).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hábito já foi marcado como completado hoje"})
		return
	}

	habitDayCheckIn := models.HabitDayCheckIn{
		HabitID: uint(habitID),
		UserID:  userID,
		Date:    time.Now(),
	}

	if err := h.db.Create(&habitDayCheckIn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao completar hábito"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": HabitCompleted})
}

func (h *HabitHandler) GetDayHabits(c *gin.Context) {
	var habits []models.Habit
	var checkIns []models.HabitDayCheckIn

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Unauthorized})
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	date := now.Format("2006-01-02")
	dayOfWeek := now.Format("Mon") // Retorna a abreviação do dia (Mon, Tue, etc)

	if err := h.db.Where("user_id = ? AND is_active = ? AND starts_at <= ? AND days_of_week LIKE ?",
		userID, true, date, "%"+dayOfWeek+"%").Find(&habits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar hábitos"})
		return
	}

	if err := h.db.Where("user_id = ? AND date >= ? AND date <= ?",
		userID, startOfDay, endOfDay).Find(&checkIns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar check-ins"})
		return
	}

	completedHabits := make(map[uint]bool)
	for _, checkIn := range checkIns {
		completedHabits[checkIn.HabitID] = true
	}

	var incompleteHabits []models.Habit = make([]models.Habit, 0)
	for _, habit := range habits {
		if !completedHabits[habit.ID] {
			incompleteHabits = append(incompleteHabits, habit)
		}
	}

	c.JSON(http.StatusOK, incompleteHabits)
}
