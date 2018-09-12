package main

import (
	"google.golang.org/api/youtube/v3"
	"cloud.google.com/go/bigquery"
)

type Channel struct {
	Id          string
	Title       string
	Description string
	Thumbnail   string
	Videos      map[string]Video
}

// Save implements the ValueSaver interface.
func (i *Channel) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Id":          i.Id,
		"Title":       i.Title,
		"Description": i.Description,
		"Thumbnail":   i.Thumbnail,
	}, i.Id, nil
}

func getChannel(channelId string, service *youtube.Service) Channel {
	Info.Printf("Channel: [%v] Fetching playlists\n", channelId)
	call := service.Channels.List("snippet,contentDetails").
		Id(channelId)
	response, err := call.Do()
	i := 0
	for !handleApiError(err) {
		if i == 5 {
			Error.Fatalf(err.Error())
		}
		response, err = call.Do()
		i++
	}
	channel := Channel{}
	if len(response.Items) > 0 {
		c := response.Items[0]
		channel.Id = c.Id
		channel.Description = c.Snippet.Description
		channel.Title = c.Snippet.Title
		channel.Thumbnail = c.Snippet.Thumbnails.Default.Url
		videos := make(map[string]Video)
		Info.Printf("Channel: [%v] Playlist [%v] found", channel.Id, c.ContentDetails.RelatedPlaylists.Uploads)
		getPlaylistVideos(c.ContentDetails.RelatedPlaylists.Uploads, service, &videos)
		channel.Videos = videos
		return channel
	}
	Error.Fatalf("Can't fetch info for channel: [%v]", channelId)
	return channel
}
