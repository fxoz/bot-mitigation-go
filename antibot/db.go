package antibot

import (
	"log"
	"sync"
	"time"
	"waffe/utils"
)

type ClientRecord struct {
	IP           string
	VerifiedAt   *time.Time
	IsJsVerified bool
}

var (
	clientCache = make(map[string]*ClientRecord)
	cacheMutex  sync.RWMutex
	cfg         = utils.LoadConfig("config.yml")
)

func IsVerified(clientIP string) bool {
	cacheMutex.RLock()
	record, exists := clientCache[clientIP]
	cacheMutex.RUnlock()

	if !exists || record.VerifiedAt == nil {
		return false
	}

	if time.Since(*record.VerifiedAt) > time.Duration(cfg.AntiBot.VerificationValidForSeconds)*time.Second {
		return false
	}

	return true
}

func RegisterClient(clientIP string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if _, exists := clientCache[clientIP]; exists {
		return
	}

	clientCache[clientIP] = &ClientRecord{
		IP:           clientIP,
		IsJsVerified: false,
	}
	log.Printf("Registered new client with IP %s", clientIP)
}

func MarkClientVerified(clientIP string) {
	now := time.Now()
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	record, exists := clientCache[clientIP]
	if !exists {
		record = &ClientRecord{IP: clientIP}
		clientCache[clientIP] = record
	}

	record.IsJsVerified = true
	record.VerifiedAt = &now
	log.Printf("Verified client with IP %s", clientIP)
}
