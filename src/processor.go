package main

import (
	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
	"sync"
	"sync/atomic"
)

// Processor struct represents info about processor
type Processor struct {
	BigQueryClient *bigquery.Client
	YouTubeService *youtube.Service
	Ctx            context.Context
}

// NewProcessor create new processor with provided context
func NewProcessor(c context.Context, b *bigquery.Client, y *youtube.Service) *Processor {
	return &Processor{
		Ctx:            c,
		BigQueryClient: b,
		YouTubeService: y,
	}
}

// Channel struct represents info about youtube channel
func (p *Processor) Channel(id string) *Channel {
	return &Channel{
		YouTubeService: p.YouTubeService,
		Id:             id,
		BigQueryClient: p.BigQueryClient,
		Ctx:            p.Ctx,
	}
}

//ProcessChannel process information about youtube channel
//get information about playlists, videos and comments and save it to BigQuery
func (p *Processor) ProcessChannel(ctx *gin.Context, channelId string) {
	Info.Printf("Channel:[%v] Processing started\n", channelId)
	//max 250 goroutines per channel
	semaphore := make(chan struct{}, 200)
	wg := new(sync.WaitGroup)
	videosChannel := make(chan *Video)
	wg.Add(1)
	channel := p.Channel(channelId)
	go func(channelId string) {
		semaphore <- struct{}{}
		defer func() {
			<-semaphore
			wg.Done()
		}()
		channel.LoadYouTubeData()
		wg.Add(1)
		go func() {
			defer wg.Done()
			channel.LoadToBigQuery()
		}()
		channel.Plist.LoadVideos(wg, videosChannel)
	}(channelId)
	var commentCounter int64
	for video := range videosChannel {
		if video.CommentCount > 0 {
			atomic.AddInt64(&commentCounter, int64(video.CommentCount))
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

	Info.Printf("Channel:[%v] Processing finished", channelId)
}
