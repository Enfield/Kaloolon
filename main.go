package main

import (
	"flag"
	"google.golang.org/api/youtube/v3"
	"os"
	"sync"
	"strings"
	"net/http"
	"google.golang.org/api/googleapi/transport"
)

var (
	developerKey = flag.String("api-key", "", "Youtube API Developer Key")
	channels     = flag.String("channel", "", "Channel id(s) to process")
	videos       = flag.String("video", "", "Video id(s) to process")
	logFile      = flag.String("log-file", "", "Logfile name")
)

func processVideos(service *youtube.Service){
	Info.Println("Start processing videos")
	//Make the API call to YouTube.
	videosList := strings.Split(*videos, " ")
	videosMap := getVideosById(videosList, service)
	videos2csv(videosMap, "")
	wg := new(sync.WaitGroup)
	for _, videoId := range videosList {
		Info.Printf("Video: [%v] Processing video info\n", videoId)
		if videosMap[videoId].CommentCount > 0 {
			wg.Add(1)
			go func(videoId string) {
				defer wg.Done()
				comments2csv(commentsByVideo(service, videosMap[videoId]), "")
			}(videoId)
		}
	}
	wg.Wait()
}

func processChannels(service *youtube.Service) {
	Info.Println("Start processing channels")
	//Make the API call to YouTube.
	channelsList := strings.Split(*channels, " ")
	wg := new(sync.WaitGroup)
	for _, channelId := range channelsList {
		videosChannel := make(chan Video)
		wg.Add(1)
		go func(channelId string) {
			defer wg.Done()
			videos := getVideos(videosChannel, channelId, service)
			videos2csv(videos, channelId)
		}(channelId)
		for video := range videosChannel {
			if video.CommentCount > 0 {
				wg.Add(1)
				go func(video Video) {
					defer wg.Done()
					comments2csv(commentsByVideo(service, video), channelId)
				}(video)
			}
		}
	}
	wg.Wait()
}

func main() {
	flag.Parse()
	if len(*logFile) > 0 {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			Init(file)
		} else {
			handleError(err, "Initialize error")
		}
	} else {
		Init(os.Stdout)
	}
	if len(*channels) > 0 || len(*videos) > 0 {
		client := &http.Client{
			Transport: &transport.APIKey{Key: *developerKey},
		}
		service, err := youtube.New(client)
		if err != nil {
			handleError(err, "")
		}
		if len(*channels) > 0 {
			processChannels(service)
		}
		if len(*videos) > 0 {
			processVideos(service)
		}
	} else {
		Info.Println("Please provide channel or video to process.")
	}
}
//Kaloolon