package photos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func TestCount(t *testing.T) {

	fmt.Println(CountPhotos(2016, 8, false))
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/photoslibrary")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	nextPageToken := ""
	mr, err := ListMediaItemsByMonth(client, 2016, 8, &nextPageToken)
	if err != nil {
		log.Fatalf("Error fetching media items: %v", err)
	}
	fmt.Println("Found media items", mr)
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func TestIndy(t *testing.T) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/photoslibrary")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	// List media items
	listURL := "https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=10"
	resp, err := client.Get(listURL)
	if err != nil {
		log.Fatalf("Unable to retrieve media items: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error fetching media items: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read response body: %v", err)
	}

	var result MediaListChunk
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Unable to parse response body: %v", err)
	}

	log.Println("Found media items", string(body))
}

func TestScanAtoiTime(t *testing.T) {
	s := "20345679"
	var v int
	start := time.Now()
	iter := 10000
	for i := 0; i < iter; i++ {
		fmt.Sscanf(s, "%d", v)
	}
	fmt.Println(v)
	dur1 := time.Since(start)
	start = time.Now()
	for i := 0; i < iter; i++ {
		v, _ = strconv.Atoi(s)
	}
	dur2 := time.Since(start)
	fmt.Println(v)
	fmt.Println("Sscanf", dur1, "Atoi", dur2)
}
