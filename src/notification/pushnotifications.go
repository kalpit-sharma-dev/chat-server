package notification

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func SendPushNotification(deviceTokens []string) error {
	decodedKey, err := getDecodedFireBaseKey()

	if err != nil {
		return err
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	app, err := firebase.NewApp(context.Background(), nil, opts...)

	if err != nil {
		log.Debug("Error in initializing firebase : %s", err)
		return err
	}

	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
		return err
	}

	// response, err := fcmClient.Send(context.Background(), &messaging.Message{

	// 	Notification: &messaging.Notification{
	// 	  Title: "Congratulations!!",
	// 	  Body: "You have just implement push notification",
	// 	},
	// 	  Token: "sample-device-token", // it's a single device token
	//   })

	//   if err != nil {
	// 	   return err
	//   }

	response, err := fcmClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: "Congratulations!!",
			Body:  "You have just implement push notification",
		},
		Tokens: deviceTokens,
	})

	if err != nil {
		return err
	}

	log.Println("Response success count : ", response.SuccessCount)
	log.Println("Response failure count : ", response.FailureCount)

	return nil
}

func getDecodedFireBaseKey() ([]byte, error) {

	fireBaseAuthKey := os.Getenv("FIREBASE_AUTH_KEY")

	decodedKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}
