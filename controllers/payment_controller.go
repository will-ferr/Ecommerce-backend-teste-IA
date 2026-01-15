package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"smart-choice/database"
	"smart-choice/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PaymentWebhookPayload struct {
	OrderID   uint    `json:"order_id"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
	Signature string  `json:"signature"`
	Timestamp int64   `json:"timestamp"`
}

func PaymentWebhook(c *gin.Context) {
	var payload PaymentWebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Error().Err(err).Msg("Invalid webhook payload structure")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	webhookSecret := c.GetHeader("X-Webhook-Secret")
	if webhookSecret == "" {
		log.Warn().Msg("Missing webhook secret header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing webhook secret"})
		return
	}

	expectedSignature := generateSignature(payload, webhookSecret)
	if !hmac.Equal([]byte(payload.Signature), []byte(expectedSignature)) {
		log.Warn().Msg("Invalid webhook signature - potential attack")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Log only non-sensitive information
	log.Info().
		Uint64("order_id", uint64(payload.OrderID)).
		Str("status", payload.Status).
		Msg("Payment webhook received")

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.First(&order, payload.OrderID).Error; err != nil {
			return err
		}

		if order.Status != payload.Status {
			order.Status = payload.Status
			if err := tx.Save(&order).Error; err != nil {
				return err
			}

			activityLog := models.ActivityLog{
				UserID:    order.UserID,
				Action:    "Payment status updated to " + payload.Status,
				Timestamp: time.Now(),
			}
			if err := tx.Create(&activityLog).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Uint64("order_id", uint64(payload.OrderID)).Msg("Failed to process webhook")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func generateSignature(payload PaymentWebhookPayload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))

	data := strconv.FormatUint(uint64(payload.OrderID), 10) + payload.Status + strconv.FormatFloat(payload.Amount, 'f', 2, 64) + strconv.FormatInt(payload.Timestamp, 10)
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}
