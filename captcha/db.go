package captcha

import (
	"image"
	"log"
	"sync"
	"time"
	"waffe/utils"
)

type CaptchaTask struct {
	IP            string
	VerifiedAt    *time.Time
	IsVerified    bool
	CorrectRegion image.Rectangle
}

var (
	captchaTasksCache = make(map[string]*CaptchaTask)
	cacheMutex        sync.RWMutex
	cfg               = utils.LoadConfig("config.yml")
)

func IsCaptchaCorrect(clientIP string, x int, y int) bool {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	record, exists := captchaTasksCache[clientIP]
	if !exists {
		return false
	}

	if record.CorrectRegion.Min.X <= x && x <= record.CorrectRegion.Max.X &&
		record.CorrectRegion.Min.Y <= y && y <= record.CorrectRegion.Max.Y {
		record.IsVerified = true
		now := time.Now()
		record.VerifiedAt = &now
		log.Printf("Captcha solved, IP %s", clientIP)
		return true
	}

	log.Printf("Captcha incorrect, IP %s", clientIP)
	return false
}

func RequiresVerification(clientIP string) bool {
	cacheMutex.RLock()
	record, exists := captchaTasksCache[clientIP]
	cacheMutex.RUnlock()

	if !exists || record.VerifiedAt == nil {
		return true
	}

	if time.Since(*record.VerifiedAt) > time.Duration(cfg.Captcha.VerificationValidForSeconds)*time.Second {
		return true
	}

	return false
}

func RegisterCaptcha(clientIP string, correctRegion image.Rectangle) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if _, exists := captchaTasksCache[clientIP]; exists {
		return
	}

	captchaTasksCache[clientIP] = &CaptchaTask{
		IP:            clientIP,
		IsVerified:    false,
		CorrectRegion: correctRegion,
	}
	log.Printf("Registered new client with IP %s", clientIP)
}

func MarkCaptchaSolved(clientIP string) {
	now := time.Now()
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	record, exists := captchaTasksCache[clientIP]
	if !exists {
		record = &CaptchaTask{IP: clientIP}
		captchaTasksCache[clientIP] = record
	}

	record.IsVerified = true
	record.VerifiedAt = &now
	log.Printf("Verified client with IP %s", clientIP)
}
