package mail

import (
	"fmt"
	"testing"
)

func TestAttachmentCount(t *testing.T) {
	// Test the count of attachments
	year := 2024
	month := 1
	count, err := CountMailWithAttachments(year, month)
	if err != nil {
		t.Errorf("Error counting mail with attachments: %v", err)
	}
	fmt.Println("Found", count, "mail with attachments")
}

func TestMailCount(t *testing.T) {
	// Test the count of mail
	year := 2024
	month := 1
	count, _, err := CountMail(year, month, false)
	if err != nil {
		t.Errorf("Error counting mail: %v", err)
	}
	fmt.Println("Found", count, "mail")
}
