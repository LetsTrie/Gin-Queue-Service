package fcm

import (
	"encoding/json"
	"fmt"
	"log"
	db "server/database"
)

func SendPushNotification(message *MulticastMessage) error {
	response, err := FcmClient.SendMulticast(ctx, message)

	if err != nil {
		return fmt.Errorf("failed to send push notification: %v", err)
	}

	log.Printf("Successfully sent push notification: %d tokens, %d successes, %d failures\n",
		len(message.Tokens), response.SuccessCount, response.FailureCount)

	if response.FailureCount > 0 {
		b, _ := json.Marshal(response)
		log.Printf("Failed tokens: %s\n", string(b))
	}

	return nil
}

func GetUserDeviceTokens(userId string) ([]string, error) {
	return db.RedisClient.SMembers(ctx, fmt.Sprintf("fcmToken:%s", userId)).Result()
}

func SendPushNotificationToUser(userId string, message *MulticastMessage) error {
	deviceTokens, err := GetUserDeviceTokens(userId)
	if err != nil {
		return fmt.Errorf("failed to get user device tokens: %v", err)
	}
	message.Tokens = deviceTokens
	return SendPushNotification(message)
}
