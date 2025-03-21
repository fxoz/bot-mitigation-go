package antibot

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	"waffe/utils"
)

type WhitelistedClient struct {
	IP                string `gorm:"primaryKey"`
	TokenIssuedAtUnix int64
	VerifiedAtUnix    *int64
	Token             string
	IsVerified        bool
}

func AddClientToWhitelist(db *gorm.DB, clientIP, token string) {
	client := WhitelistedClient{
		IP:                clientIP,
		TokenIssuedAtUnix: time.Now().Unix(),
		Token:             token,
		IsVerified:        false,
	}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ip"}},
		DoUpdates: clause.AssignmentColumns([]string{"token", "token_issued_at_unix", "is_verified"}),
	}).Create(&client).Error; err != nil {
		log.Fatalf("failed to upsert client: %v", err)
	}
}

func IsValidTokenForIP(db *gorm.DB, clientIP, token string) bool {
	var client WhitelistedClient
	if err := db.Where("ip = ?", clientIP).First(&client).Error; err != nil {
		log.Printf("client with IP %s not found: %v", clientIP, err)
		return false
	}

	if client.Token != token {
		log.Printf("token mismatch for IP %s", clientIP)
		return false
	}

	cfg := utils.LoadConfig("config.yml")
	if time.Now().Unix()-client.TokenIssuedAtUnix > int64(cfg.AntiBot.TokenValidForSeconds) {
		log.Printf("token for IP %s expired", clientIP)
		return false
	}

	return true
}

func IsClientVerified(db *gorm.DB, clientIP string) bool {
	var client WhitelistedClient
	if err := db.Where("ip = ?", clientIP).First(&client).Error; err != nil {
		log.Printf("client with IP %s not found: %v", clientIP, err)
		return false
	}
	return client.IsVerified
}

func SetClientVerified(db *gorm.DB, clientIP string) {
	var client WhitelistedClient
	if err := db.Where("ip = ?", clientIP).First(&client).Error; err != nil {
		log.Fatalf("failed to find client %s: %v", clientIP, err)
	}

	now := time.Now().Unix()
	client.IsVerified = true
	client.VerifiedAtUnix = &now

	if err := db.Save(&client).Error; err != nil {
		log.Fatalf("failed to verify client %s: %v", clientIP, err)
	}

	log.Printf("verified client with IP %s", clientIP)
}

func IsClientWhitelisted(db *gorm.DB, clientIP string) bool {
	var client WhitelistedClient
	if err := db.Where("ip = ?", clientIP).First(&client).Error; err != nil ||
		!client.IsVerified || client.VerifiedAtUnix == nil {
		log.Printf("client with IP %s not verified or not found: %v", clientIP, err)
		return false
	}

	cfg := utils.LoadConfig("config.yml")
	if time.Now().Unix()-*client.VerifiedAtUnix > int64(cfg.AntiBot.WhitelistValidForSeconds) {
		log.Printf("client with IP %s verification expired", clientIP)
		return false
	}
	return true
}

func InitDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      false,
		},
	)

	const folder = "db"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.Mkdir(folder, 0755); err != nil {
			fmt.Println("Error creating folder:", err)
		} else {
			fmt.Println("Folder created successfully.")
		}
	}

	db, err := gorm.Open(sqlite.Open("db/whitelisted_clients.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&WhitelistedClient{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return db
}
