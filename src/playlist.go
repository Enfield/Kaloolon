package main

import (
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
)

//Playlist struct represents info about youtube playlist
type Playlist struct {
	YouTubeService *youtube.Service
	BigQueryClient *bigquery.Client
	Videos         []string
	PlaylistId     string
	Channel        *Channel
	Ctx            context.Context
}

//LoadYouTubeData get information about youtube playlist
func (p *Playlist) LoadYouTubeData() {
	Info.Printf("Channel:[%v] Playlist:[%v] Fetching videos info\n", p.Channel.Id, p.PlaylistId)
	call := p.YouTubeService.PlaylistItems.List("contentDetails").
		PlaylistId(p.PlaylistId).MaxResults(50)
	response, err := call.Do()
	i := 0
	for !HandleApiError(err) {
		if i == 5 {
			Error.Fatalf(err.Error())
		}
		response, err = call.Do()
		i++
	}
	for _, i := range response.Items {
		p.Videos = append(p.Videos, i.ContentDetails.VideoId)
	}
	pageToken := response.NextPageToken
	for len(pageToken) > 0 {
		call = p.YouTubeService.PlaylistItems.List("contentDetails").
			PlaylistId(p.PlaylistId).MaxResults(50).PageToken(pageToken)
		response, err = call.Do()
		i = 0
		for !HandleApiError(err) {
			if i == 5 {
				Error.Fatalf(err.Error())
			}
			response, err = call.Do()
			i++
		}
		for _, i := range response.Items {
			p.Videos = append(p.Videos, i.ContentDetails.VideoId)
		}
		pageToken = response.NextPageToken
	}
}
