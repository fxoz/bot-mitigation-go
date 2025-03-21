package antibot

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	"waffe/utils"
)

type WhitelistedClients struct {
	IP                string `gorm:"primaryKey"`
	WhitelistedAtUnix int64
}

func IsClientWhitelisted(db *gorm.DB, clientIP string) bool {
	cfg := utils.LoadConfig("config.yml")

	var client WhitelistedClients
	if err := db.First(&client, "ip = ?", clientIP).Error; err != nil {
		return false
	}

	return time.Now().Unix()-client.WhitelistedAtUnix < int64(cfg.AntiBot.WhitelistDurationSeconds)
}

func AddClientToWhitelist(db *gorm.DB, clientIP string) {
	newClient := WhitelistedClients{IP: clientIP, WhitelistedAtUnix: time.Now().Unix()}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ip"}},
		DoUpdates: clause.AssignmentColumns([]string{"whitelisted_at_unix"}),
	}).Create(&newClient).Error; err != nil {
		log.Fatalf("failed to upsert client: %v", err)
	}
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

	db, err := gorm.Open(sqlite.Open("db/whitelisted_clients.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&WhitelistedClients{})

	return db
}
