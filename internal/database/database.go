package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"lockify-back/internal/config"
	"lockify-back/internal/models"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Exercise{}, &models.Habit{}, &models.HabitDayCheckIn{}, &models.Feedback{})
	if err != nil {
		return nil, fmt.Errorf("erro ao executar migrações: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")
	return db, nil
}
