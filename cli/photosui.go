package cli

import (
	"fmt"
	"strconv"

	"github.com/kubasiemion/googlebackup/photos"
	"github.com/manifoldco/promptui"
)

var PeriodItem = &Item{Label: "Period", Details: "Select a period"}
var YearItem = &Item{Label: "Year", Details: "Select a year"}
var MonthItem = &Item{Label: "Month", Details: "Select a month"}
var CountPhotosItem = &Item{Label: "Count Photos", Details: "Count photos in the selected period"}
var CountPhotosWithSizeItem = &Item{Label: "Count Photos with size", Details: "Count photos in the selected period with size"}
var DownloadPhotosItem = &Item{Label: "Download media", Details: "Download madia from the selected period"}
var DownloadDirItem = &Item{Label: "Download Dir", Details: "Select a download directory", Value: "./tmpDownload"}

func PhotosUI() {
	sel := promptui.Select{}
	sel.Items = []*Item{PeriodItem, YearItem, MonthItem, CountPhotosItem, CountPhotosWithSizeItem, DownloadDirItem, DownloadPhotosItem, back}
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
		case CountPhotosItem.Label:
			c, _, err := photos.CountPhotos(year, month, false)
			if err != nil {
				fmt.Println("Error counting photos", err)
				return
			}
			fmt.Println("Found", c, "photos")
		case CountPhotosWithSizeItem.Label:
			c, s, err := photos.CountPhotos(year, month, true)
			if err != nil {
				fmt.Println("Error counting photos", err)
				return
			}
			fmt.Println("Found", c, "photos with total size", photos.HumanReadableSize(s))

		case DownloadDirItem.Label:
			pr := promptui.Prompt{Label: "Enter download directory", Default: DownloadDirItem.Value}
			dir, err := pr.Run()
			if err != nil {
				fmt.Println("Error entering download directory", err)
				return
			}
			DownloadDirItem.Value = dir
			DownloadDirItem.DisplayValue = dir
		case DownloadPhotosItem.Label:
			c, s, err := photos.DownloadPhotos(year, month, DownloadDirItem.Value)
			if err != nil {
				fmt.Println("Error downloading photos", err)
				return
			}
			fmt.Println("Downloaded", c, "photos with total size", photos.HumanReadableSize(s))
		case back.Label:
			return
		}
	}

}

func SetYearUI() {
	sel := promptui.Select{}
	Items := []*Item{}
	for i := 2024; i >= 2000; i-- {
		Items = append(Items, &Item{Label: fmt.Sprintf("%d", i)})
	}
	sel.Items = append(Items, back)

	sel.Label = "Select a year"
	sel.Templates = ItemTemplate
	sel.Size = 5

	_, o, _ := sel.Run()
	switch o {
	case back.Label:
		return
	default:
		YearItem.Value = o
		YearItem.DisplayValue = o
	}

}

func SetMonthUI() {
	sel := promptui.Select{}
	Items := []*Item{}
	for i := 1; i <= 12; i++ {
		Items = append(Items, &Item{Label: fmt.Sprintf("%d", i)})
	}
	sel.Items = append(Items, back)

	sel.Label = "Select a month"
	sel.Templates = ItemTemplate
	sel.Size = 5

	_, o, _ := sel.Run()
	switch o {
	case back.Label:
		return
	default:
		MonthItem.Value = o
		MonthItem.DisplayValue = o
	}

}
