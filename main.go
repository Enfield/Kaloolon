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
	videos2csv(&videosMap, "")
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
	semaphore := make(chan struct{}, 200)
	wg := new(sync.WaitGroup)
	for _, channelId := range channelsList {
		videosChannel := make(chan Video)
		var channel Channel
		wg.Add(1)
		go func(channelId string) {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
				wg.Done()
			}()
			channel = getChannel(channelId, service)
			getVideosByChannel(&channel, videosChannel, service)
			//videos2csv(&channel.Videos, channel.Title+"_"+channel.Id)
		}(channelId)
		for video := range videosChannel {
			if video.CommentCount > 0 {
				wg.Add(1)
				go func(video Video) {
					semaphore <- struct{}{}
					defer func() {
						<-semaphore
						wg.Done()
					}()
					comments2csv(commentsByVideo(service, video), channel.Title+"_"+channel.Id)
				}(video)
			}
		}
	}
	wg.Wait()
}
func main() {
	flag.Parse()
	if len(*logFile) > 0 {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY, 0644)
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
			handleError(err, "Initialize error")
		}
		if len(*channels) > 0 {
			processChannels(service)
		}
		if len(*videos) > 0 {
			processVideos(service)
		}
	} else {
		Error.Println("Please provide channel or video to process.")
	}
}
//Kaloolon