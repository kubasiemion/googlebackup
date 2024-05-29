package cli

import (
	"fmt"

	"github.com/kubasiemion/googlebackup/photos"
	"github.com/manifoldco/promptui"
)

var MailItem = &Item{Label: "Mail"}
var PhotosItem = &Item{Label: "Photos"}
var ServiceItem = &Item{Label: "Service"}

func TopUI() {
	sel := promptui.Select{}
	sel.Items = []*Item{MailItem, PhotosItem, ServiceItem, exit}
	sel.Label = "Select an option"
	sel.Templates = ItemTemplate
	for {
		ok := TestService()
		if ok {
			MailItem.DisplayValue = "ok"
			PhotosItem.DisplayValue = "ok"
		} else {
			MailItem.DisplayValue = "not logged in"
			PhotosItem.DisplayValue = "not logged in"
		}
		_, o, _ := sel.Run()
		switch o {
		case MailItem.Label:
			MailUI()
		case PhotosItem.Label:
			PhotosUI()
		case ServiceItem.Label:
			ServiceUI()
		case exit.Label:
			return
		}
	}
}

func TestService() bool {
	_, _, err := photos.CountPhotos(2024, 1, false)
	return err == nil
}

// humanReadableSize converts a size in bytes to a human-readable string.
func HumanReadableSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
