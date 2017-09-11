package utils

import (
	"encoding/base64"

	"github.com/BKellogg/iuga-timecapsule/tc-server/models"
	gmail "google.golang.org/api/gmail/v1"
)

// SendSuccessEmail sends a time capsule success email with the given capsule
func SendSuccessEmail(gmailService *gmail.Service, capsule *models.NewCapsule) error {
	email := capsule.NetID + "@uw.edu"
	messageStr := []byte(
		"From: IUGA <iuga@uw.edu>\r\n" +
			"To: " + email + "\r\n" +
			"Subject: Your IUGA TimeCapsule\r\n\r\n" +
			"You sent this message to the IUGA TimeCapsule:\r\n" +
			capsule.Message)
	message := gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(messageStr),
	}
	_, err := gmailService.Users.Messages.Send("me", &message).Do()
	return err
}
