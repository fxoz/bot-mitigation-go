package antibot

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"waffe/utils"

	"gorm.io/gorm"
)

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func BotProtectionHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	token, err := generateRandomToken(16)
	if err != nil {
		http.Error(w, "Internal error generating token", http.StatusInternalServerError)
		return
	}

	AddClientToWhitelist(db, utils.GetIP(r), token)
	ShowVerificationPage(db, w, r, token, utils.GetFullPath(r))
}
