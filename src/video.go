package main

import (
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
	"strconv"
	"strings"
	"sync"
)

// Video struct represents info about youtube video
type Video struct {
	Id                   string
	ChannelId            string
	CategoryId           string
	PublishedAt          string
	Title                string
	LiveBroadcastContent string
	DefaultLanguage      string
	DefaultAudioLanguage string
	Duration             string
	Dimension            string
	Definition           string
	Caption              string
	LicensedContent      bool
	Projection           string
	HasCustomThumbnail   bool
	ViewCount            uint64
	LikeCount            uint64
	DislikeCount         uint64
	FavoriteCount        uint64
	CommentCount         uint64
	BigQueryClient       *bigquery.Client
	YouTubeService       *youtube.Service
	Ctx                  context.Context
	Plist                *Playlist
}

// Save implements the ValueSaver interface.
func (i Video) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Id":                   i.Id,
		"ChannelId":            i.ChannelId,
		"CategoryId":           i.CategoryId,
		"PublishedAt":          i.PublishedAt,
		"Title":                i.Title,
		"LiveBroadcastContent": i.LiveBroadcastContent,
		"DefaultLanguage":      i.DefaultLanguage,
		"DefaultAudioLanguage": i.DefaultAudioLanguage,
		"Duration":             i.Duration,
		"Dimension":            i.Dimension,
		"Definition":           i.Definition,
		"Caption":              i.Caption,
		"LicensedContent":      i.LicensedContent,
		"Projection":           i.Projection,
		"HasCustomThumbnail":   i.HasCustomThumbnail,
		//BigQuery not support uint64
		"ViewCount":     strconv.FormatUint(i.ViewCount, 10),
		"LikeCount":     strconv.FormatUint(i.LikeCount, 10),
		"DislikeCount":  strconv.FormatUint(i.DislikeCount, 10),
		"FavoriteCount": strconv.FormatUint(i.FavoriteCount, 10),
		"CommentCount":  strconv.FormatUint(i.CommentCount, 10),
	}, i.Id, nil
}

func (p *Playlist) setVideosAdditionalParametersFromResponse(response *youtube.VideoListResponse, videosChannel chan *Video) []Video {
	v := make([]Video, 0)
	for _, item := range response.Items {
		v = append(v, Video{
			Id:                   item.Id,
			ChannelId:            p.Channel.Id,
			ViewCount:            item.Statistics.ViewCount,
			LikeCount:            item.Statistics.LikeCount,
			DislikeCount:         item.Statistics.DislikeCount,
			FavoriteCount:        item.Statistics.FavoriteCount,
			CommentCount:         item.Statistics.CommentCount,
			PublishedAt:          item.Snippet.PublishedAt,
			LiveBroadcastContent: item.Snippet.LiveBroadcastContent,
			DefaultLanguage:      item.Snippet.DefaultLanguage,
			DefaultAudioLanguage: item.Snippet.DefaultAudioLanguage,
			CategoryId:           item.Snippet.CategoryId,
			Duration:             item.ContentDetails.Duration,
			Dimension:            item.ContentDetails.Dimension,
			Definition:           item.ContentDetails.Definition,
			LicensedContent:      item.ContentDetails.LicensedContent,
			Caption:              item.ContentDetails.Caption,
			Projection:           item.ContentDetails.Projection,
			HasCustomThumbnail:   item.ContentDetails.HasCustomThumbnail,
			Title:                item.Snippet.Title,
		})
		videosChannel <- &Video{Id: item.Id, CommentCount: item.Statistics.CommentCount, Ctx: p.Ctx, BigQueryClient: p.BigQueryClient, YouTubeService: p.YouTubeService, Plist: p}
	}
	return v
}

func (p *Playlist) batchLoadVideosInfo(wg *sync.WaitGroup, videosChannel chan *Video) {
	call := p.YouTubeService.Videos.List("snippet,contentDetails,statistics")
	for len(p.Videos) > 0 {
		var ids string
		if len(p.Videos) >= 50 {
			ids = strings.Join(p.Videos[:50], ",")
		} else {
			ids = strings.Join(p.Videos, ",")
		}
		call = call.Id(ids)
		response, err := call.Do()
		i := 0
		for !HandleApiError(err) {
			if i == 5 {
				Error.Fatalf(err.Error())
			}
			response, err = call.Do()
			i++
		}
		if len(p.Videos) >= 50 {
			p.Videos = p.Videos[50:]
		} else {
			p.Videos = p.Videos[len(p.Videos):]
		}
		videos := p.setVideosAdditionalParametersFromResponse(response, videosChannel)
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.LoadToBigQuery(videos)
		}()
	}
	close(videosChannel)
}

//LoadVideos load videos info from playlist
func (p *Playlist) LoadVideos(wg *sync.WaitGroup, videosChannel chan *Video) {
	Info.Printf("Channel:[%v] Playlist:[%v] Processing videos\n", p.Channel.Id, p.PlaylistId)
	p.batchLoadVideosInfo(wg, videosChannel)
}
