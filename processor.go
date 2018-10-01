package main

import (
	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
	"sync"
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
		Id:             id,
		BigQueryClient: p.BigQueryClient,
		Ctx:            p.Ctx,
	}
}

func (p *Processor) ProcessChannel(ctx *gin.Context, channelId string, ch chan int) {
	Info.Printf("Channel:[%v] Start processing\n", channelId)
	//max 250 requests per channel
	semaphore := make(chan struct{}, 250)
	wg := new(sync.WaitGroup)
	videosChannel := make(chan *Video)
	wg.Add(1)
	go func(channelId string) {
		semaphore <- struct{}{}
		defer func() {
			<-semaphore
			wg.Done()
		}()
		channel := p.Channel(channelId)
		channel.LoadYouTubeData()
		ch <- len(channel.Plist.Videos)
		wg.Add(1)
		go func() {
			defer wg.Done()
			channel.LoadToBigQuery()
		}()
		channel.Plist.LoadVideos(wg, videosChannel)
	}(channelId)
	for video := range videosChannel {
		if video.CommentCount > 0 {
			wg.Add(1)
			go func(video *Video) {
				semaphore <- struct{}{}
				defer func() {
					<-semaphore
					wg.Done()
				}()
				video.LoadComments(wg)
			}(video)
		}
	}
	wg.Wait()
	Info.Printf("Channel:[%v] Done", channelId)
}
