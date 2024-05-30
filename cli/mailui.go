package cli

import (
	"fmt"
	"strconv"

	"github.com/kubasiemion/googlebackup/mail"
	"github.com/manifoldco/promptui"
)

var CountMailItem = &Item{Label: "Count Mail", Details: "Count mail in the selected period"}
var CountMailWithSizeItem = &Item{Label: "Count Mail with size", Details: "Count mail in the selected period with size"}
var CountMailWithAttachmentsItem = &Item{Label: "Count Mail with attachments", Details: "Count mail with attachments in the selected period"}
var DownloadAttachmentsItem = &Item{Label: "Download attachments", Details: "Download attachments from the selected period"}

func MailUI() {
	sel := promptui.Select{}
	sel.Items = []*Item{YearItem, MonthItem, CountMailItem, CountMailWithAttachmentsItem, CountMailWithSizeItem, DownloadAttachmentsItem, back}
	sel.Label = "Select an option"
	sel.Templates = ItemTemplate
	var year, month int
	var err error
	for {
		_, o, _ := sel.Run()

		switch o {
		case PeriodItem.Label:
		case YearItem.Label:
			SetYearUI()
			year, err = strconv.Atoi(YearItem.Value)
			if err != nil {
				fmt.Printf("Invalid year: %s", err)

			} else {
				fmt.Println("Year", year)
			}
		case MonthItem.Label:
			SetMonthUI()
			month, err = strconv.Atoi(MonthItem.Value)
			if err != nil {
				fmt.Printf("Invalid month: %s", err)
			} else {
				fmt.Println("Month", month)
			}
		case CountMailItem.Label:
			c, _, err := mail.CountMail(year, month, false)
			if err != nil {
				fmt.Println("Error counting mail", err)
				return
			}
			fmt.Println("Found", c, "mail")
		case CountMailWithSizeItem.Label:
			c, s, err := mail.CountMail(year, month, true)
			if err != nil {
				fmt.Println("Error counting mail", err)
				return
			}
			fmt.Println("Found", c, "mail with size", HumanReadableSize(s))
		case CountMailWithAttachmentsItem.Label:
			c, err := mail.CountMailWithAttachments(year, month)
			if err != nil {
				fmt.Println("Error counting mail", err)
				return
			}
			fmt.Println("Found", c, "mail with attachments")
		case DownloadAttachmentsItem.Label:
			_, err := mail.DownloadAttachments(year, month)
			if err != nil {
				fmt.Println("Error downloading attachments", err)
				return
			}
			fmt.Println("Attachments downloaded")

		case back.Label:
			return
		}

	}
}
