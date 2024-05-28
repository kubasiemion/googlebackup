package photos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"net/http"

	"github.com/kubasiemion/googlebackup/service"
)

const rootphotosurl = "https://photoslibrary.googleapis.com/v1/mediaItems"

type MediaListChunk struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

/*
type MediaItem struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	BaseUrl  string `json:"baseUrl"`
}
*/

type MediaItem struct {
	ID            string        `json:"id"`
	Description   string        `json:"description"`
	ProductURL    string        `json:"productUrl"`
	BaseURL       string        `json:"baseUrl"`
	MimeType      string        `json:"mimeType"`
	MediaMetadata MediaMetadata `json:"mediaMetadata"`
	Filename      string        `json:"filename"`
}

type MediaMetadata struct {
	CreationTime string         `json:"creationTime"`
	Width        string         `json:"width"`
	Height       string         `json:"height"`
	Photo        *PhotoMetadata `json:"photo,omitempty"`
	Video        *VideoMetadata `json:"video,omitempty"`
}

type PhotoMetadata struct {
	CameraMake      string  `json:"cameraMake"`
	CameraModel     string  `json:"cameraModel"`
	FocalLength     float64 `json:"focalLength"`
	ApertureFNumber float64 `json:"apertureFNumber"`
	IsoEquivalent   int     `json:"isoEquivalent"`
	ExposureTime    string  `json:"exposureTime"`
}

type VideoMetadata struct {
	CameraMake  string  `json:"cameraMake"`
	CameraModel string  `json:"cameraModel"`
	FPS         float64 `json:"fps"`
	Status      string  `json:"status"`
}

func ListMediaItemsByMonth(client *http.Client, year int, month int, nextPageToken *string) (*MediaListChunk, error) {
	url := rootphotosurl + ":search"

	// Create the JSON body for the POST request
	body := []byte(fmt.Sprintf(`{
		"filters": {
			"dateFilter": {
				"dates": [
					{
						"year": %d,
						"month": %d
					}
				]
			}
		},
		"pageSize": 100,
		"pageToken": "%s"
	}`, year, month, *nextPageToken))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error fetching media items: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error (2) fetching media items: %v", resp.Status)
	}

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := new(MediaListChunk)
	err = json.Unmarshal(bts, res)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse response body: %v", err)
	}
	*nextPageToken = res.NextPageToken
	return res, err
}

// GetMediaItem retrieves a media item by its ID
// Returns filename and the media item data
func GetMediaItem(client *http.Client, mitem MediaItem) (content []byte, err error) {

	var resp *http.Response
	resp, err = http.Get(mitem.BaseURL + "=d")
	if err != nil {
		log.Fatalf("Unable to download image: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	content, err = io.ReadAll(resp.Body)

	return
}

func CountPhotos(year, month int, withsize bool) (count int, totalsize int64, err error) {
	client, err := service.GetClient(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	var nextPageToken *string
	for nextPageToken == nil || *nextPageToken != "" {
		if nextPageToken == nil {
			nextPageToken = new(string)
			*nextPageToken = ""
		}
		res, err := ListMediaItemsByMonth(client, year, month, nextPageToken)
		if err != nil {
			fmt.Println(err)
			break
		}
		count += len(res.MediaItems)
		if withsize {
			for i, item := range res.MediaItems {
				fmt.Printf("Counting %d/%d: %s\n", i+1, len(res.MediaItems), item.Filename)
				size, err := GetContentSize(item.BaseURL + "=d")
				if err != nil {
					fmt.Println(err)
					break
				}
				totalsize += size
			}
		}
	}

	return count, totalsize, nil
}

func DownloadPhotos(year, month int, dir string) (count int, totalsize int64, err error) {
	client, err := service.GetClient(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	var nextPageToken *string
	for nextPageToken == nil || *nextPageToken != "" {
		if nextPageToken == nil {
			nextPageToken = new(string)
			*nextPageToken = ""
		}
		res, err := ListMediaItemsByMonth(client, year, month, nextPageToken)
		if err != nil {
			fmt.Println(err)
			break
		}
		count += len(res.MediaItems)
		for i, item := range res.MediaItems {
			fmt.Printf("Downloading %d/%d: %s\n", i+1, len(res.MediaItems), item.Filename)
			content, err := GetMediaItem(client, item)
			if err != nil {
				fmt.Println(err)
				break
			}
			os.MkdirAll(dir, 0755)
			err = os.WriteFile(fmt.Sprintf("%s/%s", dir, item.Filename), content, 0644)
			if err != nil {
				log.Fatalf("Unable to save image: %v", err)
			}
			totalsize += int64(len(content))
		}
	}

	return count, totalsize, nil
}

func DeleteMediaItem(client *http.Client, mediaItemID string) error {
	url := fmt.Sprintf("%s/%s", rootphotosurl, mediaItemID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete media item: status %s", resp.Status)
	}

	return nil
}

func GetContentSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	size := resp.Header.Get("Content-Length")
	if size == "" {
		return 0, fmt.Errorf("Content-Length header not found")
	}

	var contentLength int
	//_, err = fmt.Sscanf(size, "%d", &contentLength)
	contentLength, err = strconv.Atoi(size)
	if err != nil {
		return 0, err
	}

	return int64(contentLength), nil
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
