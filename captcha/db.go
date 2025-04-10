package captcha

import (
	"image"
	"log"
	"sync"
	"time"
	"waffe/utils"

	"github.com/fatih/color"
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
		color.Yellow("Captcha verification failed: No record found for IP %s", clientIP)
		return false
	}

	if record.CorrectRegion.Min.X <= x && x <= record.CorrectRegion.Max.X &&
		record.CorrectRegion.Min.Y <= y && y <= record.CorrectRegion.Max.Y {
		record.IsVerified = true
		now := time.Now()
		record.VerifiedAt = &now
		log.Printf("Captcha solved for IP %s, coordinates: (%d, %d)", clientIP, x, y)
		return true
	}

	log.Printf("Captcha verification failed for IP %s, coordinates: (%d, %d)", clientIP, x, y)
	log.Printf("Correct region: (%d, %d) to (%d, %d)", record.CorrectRegion.Min.X, record.CorrectRegion.Min.Y, record.CorrectRegion.Max.X, record.CorrectRegion.Max.Y)
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
		captchaTasksCache[clientIP].CorrectRegion = correctRegion

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
