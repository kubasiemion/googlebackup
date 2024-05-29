package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/kubasiemion/googlebackup/service"
	"google.golang.org/api/gmail/v1"
)

// NextPageToken will be (most likely) modified as a result of this call
func QueryMessages(ctx context.Context, srv *gmail.Service, user, query string, pageSize int64, nextPageToken *string) (messages []*gmail.Message, err error) {

	req := srv.Users.Messages.List(user).Context(ctx).Q(query)
	if nextPageToken != nil {
		req.PageToken(*nextPageToken)
	}
	if pageSize > 0 {
		req.MaxResults(pageSize) // Set the page size
	}
	r, err := req.Do()
	if err != nil {
		return
	}
	messages = r.Messages
	*nextPageToken = r.NextPageToken

	return
}

// getMessageBody retrieves and decodes the body of a message
func GetMessageBody(message *gmail.Message) (string, error) {
	var body string

	var getMessagePartBody func(part *gmail.MessagePart) string
	getMessagePartBody = func(part *gmail.MessagePart) string {
		if part.MimeType == "text/plain" || part.MimeType == "text/html" {
			data, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err != nil {
				log.Printf("Error decoding message part body: %v", err)
				return ""
			}
			return string(data)
		}
		if part.Parts != nil {
			for _, subPart := range part.Parts {
				subBody := getMessagePartBody(subPart)
				if subBody != "" {
					return subBody
				}
			}
		}
		return ""
	}

	// Retrieve the body from the main payload or its parts
	body = getMessagePartBody(message.Payload)

	if body == "" {
		return "", fmt.Errorf("no body found in the message")
	}

	return body, nil
}

func ListEmailsByMonth(service *gmail.Service, user string, year int, month int, pageToken *string) ([]*gmail.Message, error) {
	startDate := fmt.Sprintf("%d-%02d-01", year, month)
	endDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	query := fmt.Sprintf("after:%s before:%s", startDate, endDate)

	return QueryMessages(context.Background(), service, user, query, 100, pageToken)

}

func CountMail(year, month int, withSize bool) (count int, size int64, err error) {
	srv, err := service.GetService()
	if err != nil {
		return
	}
	user := "me"
	pageToken := ""
	for {
		messages, err := ListEmailsByMonth(srv, user, year, month, &pageToken)
		if err != nil {
			return 0, 0, err
		}
		count += len(messages)
		if withSize {
			for _, m := range messages {
				msg, err := srv.Users.Messages.Get(user, m.Id).Do()
				if err != nil {
					return 0, 0, fmt.Errorf("error retrieving message %s: %w", m.Id, err)
				}
				size += msg.SizeEstimate
			}
		}
		if pageToken == "" {
			break
		}
	}
	return
}
