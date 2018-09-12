package main

import (
	"google.golang.org/api/youtube/v3"
	"strings"
	"cloud.google.com/go/bigquery"
	"strconv"
)

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
		"ViewCount":            strconv.FormatUint(i.ViewCount, 10),
		"LikeCount":            strconv.FormatUint(i.LikeCount, 10),
		"DislikeCount":         strconv.FormatUint(i.DislikeCount, 10),
		"FavoriteCount":        strconv.FormatUint(i.FavoriteCount, 10),
		"CommentCount":         strconv.FormatUint(i.CommentCount, 10),
	}, i.Id, nil
}

func setVideosAdditionalParametersFromResponse(result *map[string]Video, videosChannel chan Video, response *youtube.VideoListResponse, channelId string) {
	for _, item := range response.Items {
		r := *result
		video := r[item.Id]
		video.ChannelId = channelId
		video.ViewCount = item.Statistics.ViewCount
		video.LikeCount = item.Statistics.LikeCount
		video.DislikeCount = item.Statistics.DislikeCount
		video.FavoriteCount = item.Statistics.FavoriteCount
		video.CommentCount = item.Statistics.CommentCount
		video.PublishedAt = item.Snippet.PublishedAt
		video.LiveBroadcastContent = item.Snippet.LiveBroadcastContent
		video.DefaultLanguage = item.Snippet.DefaultLanguage
		video.DefaultAudioLanguage = item.Snippet.DefaultAudioLanguage
		video.CategoryId = item.Snippet.CategoryId
		video.Duration = item.ContentDetails.Duration
		video.Dimension = item.ContentDetails.Dimension
		video.Definition = item.ContentDetails.Definition
		video.LicensedContent = item.ContentDetails.LicensedContent
		video.Caption = item.ContentDetails.Caption
		video.Projection = item.ContentDetails.Projection
		video.HasCustomThumbnail = item.ContentDetails.HasCustomThumbnail
		video.Title = item.Snippet.Title
		r[video.Id] = video
		videosChannel <- video
	}

}

func batchLoadVideosInfo(service *youtube.Service, videosChannel chan Video, videosMap *map[string]Video, channelId string) {
	keys := make([]string, 0, len(*videosMap))
	for k := range *videosMap {
		keys = append(keys, k)
	}
	call := service.Videos.List("snippet,contentDetails,statistics")
	for len(keys) > 0 {
		var ids string
		if len(keys) >= 50 {
			ids = strings.Join(keys[:50], ",")
		} else {
			ids = strings.Join(keys, ",")
		}
		call = call.Id(ids)
		response, err := call.Do()
		i := 0
		for !handleApiError(err) {
			if i == 5 {
				Error.Fatalf(err.Error())
			}
			response, err = call.Do()
			i++
		}
		if len(keys) >= 50 {
			keys = keys[50:]
		} else {
			keys = keys[len(keys):]
		}
		setVideosAdditionalParametersFromResponse(videosMap, videosChannel, response, channelId)
	}
	close(videosChannel)
}

func addVideosFromVideoListResponseToMap(result map[string]Video, response *youtube.VideoListResponse) map[string]Video {
	for _, item := range response.Items {
		video := Video{}
		video.Id = item.Id
		video.Title = item.Snippet.Title
		video.ChannelId = item.Snippet.ChannelId
		video.ViewCount = item.Statistics.ViewCount
		video.LikeCount = item.Statistics.LikeCount
		video.DislikeCount = item.Statistics.DislikeCount
		video.FavoriteCount = item.Statistics.FavoriteCount
		video.CommentCount = item.Statistics.CommentCount
		video.PublishedAt = item.Snippet.PublishedAt
		video.LiveBroadcastContent = item.Snippet.LiveBroadcastContent
		video.DefaultLanguage = item.Snippet.DefaultLanguage
		video.DefaultAudioLanguage = item.Snippet.DefaultAudioLanguage
		video.CategoryId = item.Snippet.CategoryId
		video.Duration = item.ContentDetails.Duration
		video.Dimension = item.ContentDetails.Dimension
		video.Definition = item.ContentDetails.Definition
		video.LicensedContent = item.ContentDetails.LicensedContent
		video.Caption = item.ContentDetails.Caption
		video.Projection = item.ContentDetails.Projection
		video.HasCustomThumbnail = item.ContentDetails.HasCustomThumbnail
		result[video.Id] = video
	}
	return result
}

func getVideosByChannel(channel *Channel, videosChannel chan Video, service *youtube.Service) {
	Info.Printf("Channel: [%v] Processing videos\n", channel.Id)
	batchLoadVideosInfo(service, videosChannel, &channel.Videos, channel.Id)
}

func getVideosById(videoIds []string, service *youtube.Service) map[string]Video {
	ids := strings.Join(videoIds[:], ",")
	videosMap := make(map[string]Video)
	call := service.Videos.List("snippet,contentDetails,statistics").Id(ids)
	response, err := call.Do()
	i := 0
	for !handleApiError(err) {
		if i == 5 {
			Error.Fatalf(err.Error())
		}
		response, err = call.Do()
		i++
	}
	videosMap = addVideosFromVideoListResponseToMap(videosMap, response)
	return videosMap
}
