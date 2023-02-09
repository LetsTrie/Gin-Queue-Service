package fcm

import (
	"context"
	"encoding/base64"
	"fmt"
	"server/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	"google.golang.org/api/option"
)

var FcmClient *messaging.Client
var ctx = context.Background()

type MulticastMessage = messaging.MulticastMessage
type Notification = messaging.Notification

func FcmConfig() error {
	decodedKey, err := base64.StdEncoding.DecodeString(config.Env.FirebaseAuthKey)
	if err != nil {
		return fmt.Errorf("failed to get decoded Firebase Auth key: %v", err)
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase App: %v", err)
	}

	FcmClient, err = app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Messaging client: %v", err)
	}

	return nil
}
