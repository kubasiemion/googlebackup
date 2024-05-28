package cli

import "github.com/manifoldco/promptui"

var MailItem = &Item{Label: "Mail"}
var PhotosItem = &Item{Label: "Photos"}
var ServiceItem = &Item{Label: "Service"}

func TopUI() {
	sel := promptui.Select{}
	sel.Items = []*Item{MailItem, PhotosItem, ServiceItem, exit}
	sel.Label = "Select an option"
	sel.Templates = ItemTemplate
	for {
		_, o, _ := sel.Run()
		switch o {
		case MailItem.Label:
		case PhotosItem.Label:
			PhotosUI()
		case ServiceItem.Label:
			ServiceUI()
		case exit.Label:
			return
		}
	}
}
