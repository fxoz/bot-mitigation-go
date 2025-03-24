package antibot

import (
	"fmt"
	"log"
	"os"
	"time"
	"waffe/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type CheckedClient struct {
	IP           string `gorm:"primaryKey"`
	PassedAtUnix *int64
	IsVerified   bool
}

func getClient(db *gorm.DB, clientIP string) (*CheckedClient, error) {
	var client CheckedClient
	err := db.Where("ip = ?", clientIP).First(&client).Error
	return &client, err
}

func NeedsVerification(db *gorm.DB, clientIP string) bool {
	cfg := utils.LoadConfig("config.yml")
	client, err := getClient(db, clientIP)
	if err != nil {
		return true
	}

	if client.PassedAtUnix == nil {
		return true
	}

	verifiedAt := time.Unix(*client.PassedAtUnix, 0)
	return time.Since(verifiedAt).Seconds() > float64(cfg.AntiBot.VerificationValidForSeconds)
}

func IsClientVerified(db *gorm.DB, clientIP string) bool {
	cfg := utils.LoadConfig("config.yml")
	client, err := getClient(db, clientIP)
	if err != nil || !client.IsVerified || client.PassedAtUnix == nil {
		return false
	}

	verifiedAt := time.Unix(*client.PassedAtUnix, 0)
	return time.Since(verifiedAt).Seconds() <= float64(cfg.AntiBot.VerificationValidForSeconds)
}

func AddClient(db *gorm.DB, clientIP string) {
	if _, err := getClient(db, clientIP); err == nil {
		return
	}

	client := CheckedClient{IP: clientIP, IsVerified: false}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ip"}},
		DoNothing: true,
	}).Create(&client).Error; err != nil {
		log.Fatalf("failed to add client %s: %v", clientIP, err)
	}
	log.Printf("added client with IP %s", clientIP)
}

func SetClientVerified(db *gorm.DB, clientIP string) {
	client, err := getClient(db, clientIP)
	if err != nil {
		log.Fatalf("failed to find client %s: %v", clientIP, err)
	}
	now := time.Now().Unix()
	client.IsVerified = true
	client.PassedAtUnix = &now

	if err := db.Save(client).Error; err != nil {
		log.Fatalf("failed to verify client %s: %v", clientIP, err)
	}
	log.Printf("verified client with IP %s", clientIP)
}

func InitDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{SlowThreshold: time.Second, LogLevel: logger.Silent, Colorful: false},
	)

	const folder = "db"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.Mkdir(folder, 0755); err != nil {
			fmt.Println("Error creating folder:", err)
		} else {
			fmt.Println("Folder created successfully.")
		}
	}

	db, err := gorm.Open(sqlite.Open("db/antibot-verification.db"), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&CheckedClient{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return db
}
