package main

import (
	"flag"
	"google.golang.org/api/youtube/v3"
	"os"
	"net/http"
	"google.golang.org/api/googleapi/transport"
	// Imports the Google Cloud BigQuery client package.
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
)

var (
	developerKey = flag.String("api-key", "", "Youtube API Developer Key")
	channels     = flag.String("channel", "", "Channel id(s) to process")
	videos       = flag.String("video", "", "Video id(s) to process")
	logFile      = flag.String("log-file", "", "Logfile name")
)

func main() {
	flag.Parse()
	if len(*logFile) > 0 {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			InitLogger(file)
		} else {
			HandleFatalError(err, "Initialize error")
		}
	} else {
		InitLogger(os.Stdout)
	}
	if len(*channels) > 0 || len(*videos) > 0 {
		youTubeClient := &http.Client{
			Transport: &transport.APIKey{Key: *developerKey},
		}
		youTubeService, err := youtube.New(youTubeClient)
		if err != nil {
			HandleFatalError(err, "Initialize error")
		}
		// Sets your Google Cloud Platform project ID.
		projectID := "youtube-analyzer-206211"

		// Creates a client.
		ctx := context.Background()
		bigQueryClient, err := bigquery.NewClient(ctx, projectID)
		if err != nil {
			HandleFatalError(err, "Initialize error")
		}
		processor :=  NewProcessor(bigQueryClient, youTubeService, ctx)
		processor.ProcessChannels(channels)
	} else {
		Error.Println("Please provide channel or video to process.")
	}
}
