package db

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"
)

func (store *SQLStore) SendChildNotification(ct context.Context, userId int32, notificationType, body string) error {

	token, err := store.GetChildToken(ct, userId)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := store.fApp.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("error getting Messaging client: %v", err)
	}

	message := &messaging.Message{
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Title:                 "Sanduk",
				Body:                  body,
				Icon:                  "",
				VibrateTimingMillis:   []int64{1500, 500},
				ClickAction:           "FLUTTER_NOTIFICATION_CLICK",
				Visibility:            messaging.VisibilityPublic,
				Tag:                   "",
				Priority:              messaging.PriorityMax,
				Sticky:                true,
				DefaultSound:          true,
				DefaultVibrateTimings: true,
			},
		},

		Token: token.Token,

		Data: map[string]string{
			"key":    "NOTIFICATION_TYPE",
			"value":  notificationType,
			"status": "status",
		},
	}
	_, err = client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("could not send message: %v", err)
	}
	return nil

}

func (store *SQLStore) SendParentsNotification(ct context.Context, userId int32, notificationType, body string) error {

	token, err := store.GetParentsToken(ct, userId)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := store.fApp.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("error getting Messaging client: %v", err)
	}

	message := &messaging.Message{
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Title:                 "Sanduk",
				Body:                  body,
				Icon:                  "",
				VibrateTimingMillis:   []int64{1500, 500},
				ClickAction:           "FLUTTER_NOTIFICATION_CLICK",
				Visibility:            messaging.VisibilityPublic,
				Tag:                   "",
				Priority:              messaging.PriorityMax,
				Sticky:                true,
				DefaultSound:          true,
				DefaultVibrateTimings: true,
			},
		},

		Token: token.Token,

		Data: map[string]string{
			"key":    "NOTIFICATION_TYPE",
			"value":  notificationType,
			"status": "status",
		},
	}
	_, err = client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("could not send message: %v", err)
	}
	return nil

}
