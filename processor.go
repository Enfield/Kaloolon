package main

import (
	"strings"
	"sync"
	"google.golang.org/api/youtube/v3"
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
)

type Processor struct {
	BigQueryClient *bigquery.Client
	YouTubeService *youtube.Service
	Ctx            context.Context
}

func NewProcessor(b *bigquery.Client, y *youtube.Service, c context.Context) *Processor {
	return &Processor{
		BigQueryClient: b,
		YouTubeService: y,
		Ctx:            c,
	}
}

func (p *Processor) Channel(id string) *Channel {
	return &Channel{
		YouTubeService: p.YouTubeService,
		Id: id,
	}
}

func (p *Processor) ProcessVideos(service *youtube.Service) {
	Info.Println("Start processing videos")
	//Make the API call to YouTube.
	videosList := strings.Split(*videos, " ")
	videosMap := getVideosById(videosList, service)
	//videos2csv(&videosMap, "")
	wg := new(sync.WaitGroup)
	for _, videoId := range videosList {
		Info.Printf("Video: [%v] Processing video info\n", videoId)
		if videosMap[videoId].CommentCount > 0 {
			wg.Add(1)
			go func(videoId string) {
				defer wg.Done()
				//comments2csv(LoadYouTubeData(service, videosMap[videoId]), "")
			}(videoId)
		}
	}
	wg.Wait()
}

func (p *Processor) ProcessChannels(channels *[]string) {
	Info.Println("Start processing channels")
	//max 250 requests per channel
	semaphore := make(chan struct{}, 250)
	wg := new(sync.WaitGroup)
	for _, channelId := range *channels {
		videosChannel := make(chan string)
		wg.Add(1)
		go func(channelId string) {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
				wg.Done()
			}()
			channel := p.Channel(channelId)
			channel.LoadYouTubeData(videosChannel)
			go channel.LoadToBigQuery()
			channel.Plist.LoadVideos()
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
					video.LoadComments()
				}(video)
			}
		}
	}
	wg.Wait()
}
