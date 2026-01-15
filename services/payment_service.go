package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"

	"smart-choice/database"
	"smart-choice/repository"

	"gorm.io/gorm"
)

func ProcessPaymentWebhook(payload, signature string, orderID uint) error {
	if !validatePayload(payload, signature) {
		return errors.New("invalid signature")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := repository.UpdateOrderStatus(orderID, "paid"); err != nil {
			return err
		}

		// Adicionar lógica adicional da transação aqui, se necessário

		return nil
	})
}

func validatePayload(payload, signature string) bool {
	secret := os.Getenv("WEBHOOK_SECRET")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
