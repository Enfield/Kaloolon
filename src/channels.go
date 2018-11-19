package main

import (
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
	"strconv"
)

//Channel struct represents info about youtube channel
type Channel struct {
	Id              string `json:"id,omitempty"`
	Title           string `json:"title,omitempty"`
	Description     string `json:"description,omitempty"`
	Thumbnail       string `json:"avatar,omitempty"`
	ViewCount       string `json:"viewCount,omitempty"`
	SubscriberCount string `json:"subscriberCount,omitempty"`
	CommentCount    string `json:"commentCount,omitempty"`
	VideoCount      string `json:"videoCount,omitempty"`
	YouTubeService  *youtube.Service
	BigQueryClient  *bigquery.Client
	Ctx             context.Context
	Plist           Playlist
}

// Save implements the ValueSaver interface.
func (c *Channel) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Id":          c.Id,
		"Title":       c.Title,
		"Description": c.Description,
		"Thumbnail":   c.Thumbnail,
		//bigquery not support uint64
		"ViewCount":       c.ViewCount,
		"SubscriberCount": c.SubscriberCount,
		//"CommentCount":    c.CommentCount,
		//"VideoCount":      c.VideoCount,
	}, c.Id, nil
}

// Playlist is function to create default empty channel playlist with provided context
func (c *Channel) Playlist(id string) *Playlist {
	c.Plist = Playlist{
		YouTubeService: c.YouTubeService,
		BigQueryClient: c.BigQueryClient,
		Ctx:            c.Ctx,
		PlaylistId:     id,
		Videos:         make([]string, 0),
		Channel:        c,
	}
	return &c.Plist
}

//LoadYouTubeData get information about youtube channel
func (c *Channel) LoadYouTubeData() {
	Info.Printf("Channel:[%v] Fetching playlists\n", c.Id)
	call := c.YouTubeService.Channels.List("snippet,contentDetails,statistics").
		Id(c.Id)
	response, err := call.Do()
	i := 0
	for !HandleApiError(err) {
		if i == 5 {
			Error.Fatalf(err.Error())
		}
		response, err = call.Do()
		i++
	}
	if len(response.Items) > 0 {
		channelResponse := response.Items[0]
		c.Description = channelResponse.Snippet.Description
		c.Title = channelResponse.Snippet.Title
		c.Thumbnail = channelResponse.Snippet.Thumbnails.Default.Url
		c.ViewCount = strconv.FormatUint(channelResponse.Statistics.ViewCount, 10)
		c.SubscriberCount = strconv.FormatUint(channelResponse.Statistics.SubscriberCount, 10)
		c.Playlist(channelResponse.ContentDetails.RelatedPlaylists.Uploads).LoadYouTubeData()
	}
}
