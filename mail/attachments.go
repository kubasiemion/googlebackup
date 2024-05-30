package mail

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/api/gmail/v1"
)

func CountMailWithAttachments(year, month int) (count int, err error) {
	query := "has:attachment"

	count, _, err = CountMail(year, month, false, query)
	return
}

func DownloadAttachments(year, month int) (count int, err error) {
	err = fmt.Errorf("downoad of attachments not implemented")
	return
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
