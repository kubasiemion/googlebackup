package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var token *oauth2.Token
var config *oauth2.Config
var selfaddress = "http://localhost:18080/token"

func init() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err = google.ConfigFromJSON(b, gmail.GmailModifyScope, "https://www.googleapis.com/auth/photoslibrary")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	config.RedirectURL = selfaddress
	//---------------------
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		r.FormValue("code")
		if r.FormValue("code") != "" {
			tokenreceived <- r.FormValue("code")
			fmt.Fprintf(w, "Token received. You can close this tab now.")
		}
	})
	//---------------------
	token, err = TokenFromFile(TokenFile)
	if err != nil {
		log.Println(err)
	}
}

func HaveToken() bool {
	return token != nil
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func GetClient(ctx context.Context) (*http.Client, error) {
	if !HaveToken() {
		return nil, fmt.Errorf("do not have the OAuth token")
	}

	return config.Client(ctx, token), nil
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func GetTokenFromWeb(ctx context.Context) error {

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	//fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
	go srv.ListenAndServe()
	err := OpenInBrowser(authURL)
	if err != nil {
		return fmt.Errorf("unable to open browser: %w", err)

	}
	fmt.Println("Waiting for token...")
	var authCode string
	select {
	case authCode = <-tokenreceived:
		fmt.Println("Token received:", authCode)
		break
	case <-time.After(10 * time.Minute):
		log.Fatalf("Timeout waiting for token")
		return nil
	}
	srv.Shutdown(context.Background())

	//if _, err := fmt.Scan(&authCode); err != nil {
	//	log.Fatalf("Unable to read authorization code %v", err)
	//}

	token, err = config.Exchange(ctx, authCode)
	if err != nil {
		return fmt.Errorf("unable to retrieve token from web %w", err)
	}
	SaveToken(TokenFile, token)
	return nil
}

var TokenFile = "ntoken.json"

// TokenFromFile retrieves a Token from a given file path.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	token = tok
	return tok, err
}
func setByCode(code string) (err error) {
	token, err = config.Exchange(context.Background(), code)
	return
	//saveToken("token.json", token)
}

const code = "4/0AdLIrYeFfYYwomuACaHomQa3XH6DTdeEJLOT4xCt3yjR9NLck0hrok7irza72wFUEkCsjg"

// saveToken uses a file path to create a file and store the token in it.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetService() (*gmail.Service, error) {
	// Use a context with a timeout
	//ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	ctx := context.Background()

	client, err := GetClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting service: %w", err)
	}

	return gmail.NewService(ctx, option.WithHTTPClient(client))
}

func Config() string {
	return fmt.Sprintf("OAuth Config: %v", config)
}

func OpenInBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = "open"
		args = append(args, url)
	case "linux":
		cmd = "xdg-open"
		args = append(args, url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}

var tokenreceived = make(chan string)

var srv = http.Server{
	Addr: ":18080",
}
