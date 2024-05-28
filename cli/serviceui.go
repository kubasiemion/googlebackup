package cli

import (
	"context"
	"fmt"

	"github.com/kubasiemion/googlebackup/service"
	"github.com/manifoldco/promptui"
)

const ConfigItem = "OAuth Config"
const TokenItem = "OAuth Token"
const ServiceTestItem = "Service Test"
const UpItem = "Up"

func ServiceUI() {
	sel := promptui.Select{}
	sel.Items = []string{ConfigItem, TokenItem, ServiceTestItem, UpItem}
	sel.Label = "Select an option"
	for {
		_, o, _ := sel.Run()
		switch o {
		case ConfigItem:
			ConfigUI()
		case TokenItem:
			TokenUI()
		case ServiceTestItem:
			_, err := service.GetService()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Service OK")
			}

		case UpItem:
			return
		}
	}

}

func ConfigUI() {
	fmt.Print(service.Config())
}

const ShowTokenItem = "Show Token"
const GetTokenItem = "Get Token"
const ForceNewTokenItem = "Force New Token"

func TokenUI() {
	sel := promptui.Select{}
	sel.Items = []string{ShowTokenItem, GetTokenItem, ForceNewTokenItem, UpItem}
	sel.Label = "Select an option"
	for {
		_, o, _ := sel.Run()
		switch o {
		case ShowTokenItem:
			if service.HaveToken() {
				fmt.Println("service.Token()")
			} else {
				fmt.Println("No token")
			}

		case GetTokenItem:
			_, err := service.TokenFromFile(service.TokenFile)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Token OK")
			}

		case ForceNewTokenItem:
			service.GetTokenFromWeb(context.Background())

		case UpItem:
			return
		}
	}
}
