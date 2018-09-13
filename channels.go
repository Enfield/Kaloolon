package main

import (
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/youtube/v3"
	"golang.org/x/net/context"
)

type Channel struct {
	Id             string
	Title          string
	Description    string
	Thumbnail      string
	YouTubeService *youtube.Service
	BigQueryClient *bigquery.Client
	Ctx            context.Context
	Plist Playlist
}

// Save implements the ValueSaver interface.
func (c *Channel) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Id":          c.Id,
		"Title":       c.Title,
		"Description": c.Description,
		"Thumbnail":   c.Thumbnail,
	}, c.Id, nil
}

func (c *Channel) Playlist(id string) *Playlist {
	c.Plist = Playlist{
		YouTubeService: c.YouTubeService,
		PlaylistId:     id,
		Videos:         make([]string, 0),
		Channel: c,
	}
	return &c.Plist
}

func (c *Channel) LoadYouTubeData(videosChannel chan string) {
	Info.Printf("Channel: [%v] Fetching playlists\n", c.Id)
	call := c.YouTubeService.Channels.List("snippet,contentDetails").
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
		Info.Printf("Channel:[%v] Playlist:[%v] found", c.Id, channelResponse.ContentDetails.RelatedPlaylists.Uploads)
		c.Playlist(channelResponse.ContentDetails.RelatedPlaylists.Uploads).LoadYouTubeData(videosChannel)
	}
	Error.Fatalf("Can't fetch info for channel: [%v]", c.Id)
}
