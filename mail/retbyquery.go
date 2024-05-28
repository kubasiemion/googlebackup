package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

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

/*
srv, err := service.GetService()
	user := "me"

	//query := fmt.Sprintf("STATEMENT OF INCOME 2018 IN SPAIN")
	query := fmt.Sprintf("rfc822msgid:%s", "1503754224.99672026@emsgrid.com1503754224")
	var NextPageToken = new(string)
	msgs, err := mail.QueryMessages(context.Background(), srv, user, query, 100, NextPageToken)
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	for _, msg := range msgs {
		m, err := srv.Users.Messages.Get(user, msg.Id).Format("full").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve message: %v", err)
		}
		fmt.Printf("Message <%s> snippet: %s\n", msg.Id, msg.Snippet)
		fmt.Println(m.ForceSendFields)
		tim := time.Unix(m.InternalDate/1000, 0)
		fmt.Println(tim)
		fmt.Println(m.Payload.Filename)
		fmt.Println(mail.SizeMessage(srv, user, m))
		names, bts, err := mail.GetAttachments(srv, user, m)
		if err != nil {
			log.Fatalf("Unable to retrieve attachments: %v", err)
		}
		for i, name := range names {
			fmt.Printf("Attachment %s: %d bytes\n", name, len(bts[i]))
		}
		_, err = mail.GetMessageBody(m)
		if err != nil {
			log.Fatalf("Unable to retrieve body: %v", err)
		}
*/
