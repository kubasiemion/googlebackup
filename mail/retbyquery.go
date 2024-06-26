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
func QueryMessages(ctx context.Context, srv *gmail.Service, user, query string, pageSize int64, nextPageToken *string, paging bool) (messages []*gmail.Message, err error) {

	req := srv.Users.Messages.List(user).Context(ctx).Q(query)
	var resp *gmail.ListMessagesResponse
	if nextPageToken != nil {
		req.PageToken(*nextPageToken)
	}
	for {
		if pageSize > 0 {
			req.MaxResults(pageSize) // Set the page size
		}
		resp, err = req.Do()
		if err != nil {
			return
		}
		messages = append(messages, resp.Messages...)
		*nextPageToken = resp.NextPageToken
		req.PageToken(*nextPageToken)
		if paging || *nextPageToken == "" {
			break
		}
	}

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

func ListEmailsByMonth(service *gmail.Service, user string, year int, month int, pageToken *string, paging bool, aquery ...string) ([]*gmail.Message, error) {
	startDate := fmt.Sprintf("%d-%02d-01", year, month)
	endDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	query := fmt.Sprintf("after:%s before:%s", startDate, endDate)
	for _, q := range aquery {
		query += " " + q
	}
	return QueryMessages(context.Background(), service, user, query, 100, pageToken, paging)

}

func CountMail(year, month int, withSize bool, query ...string) (count int, size int64, err error) {
	srv, err := service.GetService()
	if err != nil {
		return
	}
	user := "me"
	pageToken := ""
	messages, err := ListEmailsByMonth(srv, user, year, month, &pageToken, false, query...)
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

	return
}
