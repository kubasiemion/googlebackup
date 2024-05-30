package mail

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/kubasiemion/googlebackup/service"
	"google.golang.org/api/gmail/v1"
)

func CountMailWithAttachments(year, month int) (count int, err error) {
	query := "has:attachment"

	count, _, err = CountMail(year, month, false, query)
	return
}

func DownloadAttachments(year, month int, dir string) (count int, err error) {
	query := "has:attachment"
	srv, err := service.GetService()
	if err != nil {
		return
	}
	user := "me"
	pageToken := ""
	messages, err := ListEmailsByMonth(srv, user, year, month, &pageToken, false, query)
	if err != nil {
		return 0, err
	}
	count += len(messages)
	for _, m := range messages {
		msg, err := srv.Users.Messages.Get(user, m.Id).Do()
		if err != nil {
			return 0, fmt.Errorf("error retrieving message %s: %w", m.Id, err)
		}
		names, contents, err := GetAttachments(srv, user, msg)
		if err != nil {
			return 0, fmt.Errorf("error retrieving attachments for message %s: %w", m.Id, err)
		}
		for i, name := range names {
			err = SaveAttachment(fmt.Sprintf("%s/%s", dir, name), contents[i])
			if err != nil {
				return 0, fmt.Errorf("error saving attachment %s: %w", name, err)
			}
		}
	}
	return
}

func SaveAttachment(name string, content []byte) error {
	os.WriteFile(name, content, 0644)
	return nil
}

// Retrieves and returns the attachments of a message.
func GetAttachments(srv *gmail.Service, user string, message *gmail.Message) (names []string, contents [][]byte, err error) {
	for _, part := range message.Payload.Parts {
		if part.Filename != "" && part.Body != nil && part.Body.AttachmentId != "" {
			var attachment *gmail.MessagePartBody
			attachment, err = srv.Users.Messages.Attachments.Get(user, message.Id, part.Body.AttachmentId).Do()
			if err != nil {
				return
			}

			// Decode the base64 encoded data
			var data []byte
			data, err = base64.URLEncoding.DecodeString(attachment.Data)
			if err != nil {
				err = fmt.Errorf("unable to decode attachment data: %v", err)
				return
			}

			names = append(names, part.Filename)
			contents = append(contents, data)

		}
	}
	return
}

// Retrieves and returns the attachments of a message.
func SizeMessage(srv *gmail.Service, user string, message *gmail.Message) (body, attachments int64, err error) {
	for _, part := range message.Payload.Parts {
		if part.Filename != "" && part.Body != nil && part.Body.AttachmentId != "" {
			attachments += part.Body.Size
		} else {
			body += part.Body.Size
		}
	}
	bod, err := GetMessageBody(message)
	if err != nil {
		return
	}
	body += int64(len(bod))
	return
}
