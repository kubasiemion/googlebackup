package cli

import (
	"github.com/manifoldco/promptui"
)

const (
	EXIT = "EXIT"
	BACK = "BACK"
)

// Create a struct to keep cipherLabel and cipherName
type Item struct {
	Label        string
	Details      string
	Value        string
	DisplayValue string
	Color        string
}

// Rewrite the above as variable instantiation witt OK=&Item{Label: "OK"} pattern
var OK = &Item{Label: "OK", Details: "Confirm and proceed"}
var back = &Item{Label: BACK, Details: "Go back to previous menu", Color: "yellow"}
var exit = &Item{Label: EXIT, Details: "Exit the program", Color: "red"}

/*
var vanity = &Item{Label: "Vanity", Name: "", Details: "Vanity pattern for the address as regex,\n e.g. '^F00D'"}
var vanCS = &Item{Label: "Vanity case sensitive?", Name: "no", Details: "Is vanity pattern case sensitive? (yes/no)"}
var setPassword = &Item{Label: "Set password"}
var enterPassword = &Item{Label: "Enter password"}
var kdf = &Item{Label: "Key derivation function", Name: "scrypt", Details: "scrypt and pbkdf2 are supported"}
*/
func (i Item) String() string {
	return i.Label
}

var ItemTemplate = &promptui.SelectTemplates{
	Label:    "{{ . | bold | cyan}}",
	Inactive: `{{if eq .Label "BACK"}}{{.Label | yellow}}{{else if eq .Label "EXIT"}}{{.Label | red}}{{else}}{{ .Label }}{{with .DisplayValue}}: {{.}}{{end}}{{end}}`,
	Active:   `{{if eq .Label "BACK"}}{{.Label | yellow | bold | underline}}{{else if eq .Label "EXIT"}}{{.Label | red | bold | underline}}{{else}}{{ .Label | bold | underline }}{{with .DisplayValue}}: {{. | bold}}{{end}}{{end}}`,
	Selected: "{{ .Label | red }}",
	Details:  "{{ .Details | faint }}",
}

var ItemAltTemplate = &promptui.SelectTemplates{
	Label:    "{{ . | bold | cyan}}",
	Inactive: `{{ . | bold | .Color}}`,
	Active:   `{{if not eq .Color ""}}{{.Label | .Color | bold | underline}}{{else}}{{ .Label | bold | underline }}{{with .DisplayValue}}: {{. | bold}}{{end}}{{end}}`,
	Selected: "{{ .Label | red }}",
	Details:  "{{ .Details | faint }}",
}
