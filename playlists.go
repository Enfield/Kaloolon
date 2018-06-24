package main

import "google.golang.org/api/youtube/v3"

func getPlaylistVideos(playlistId string, service *youtube.Service, videos *map[string]Video) {
	v := *videos
	Info.Printf("Playlist:[%v] Fetching videos info\n", playlistId)
	call := service.PlaylistItems.List("contentDetails").
		PlaylistId(playlistId).MaxResults(50)
	response, err := call.Do()
	i := 0
	for !handleApiError(err) {
		if i == 5 {
			Error.Fatalf(err.Error())
		}
		response, err = call.Do()
		i++
	}
	for _, p := range response.Items {
		video := Video{}
		video.Id = p.ContentDetails.VideoId
		v[video.Id] = video
	}
	pageToken := response.NextPageToken
	for len(pageToken) > 0 {
		call = service.PlaylistItems.List("contentDetails").
			PlaylistId(playlistId).MaxResults(50).PageToken(pageToken)
		response, err = call.Do()
		i = 0
		for !handleApiError(err) {
			if i == 5 {
				Error.Fatalf(err.Error())
			}
			response, err = call.Do()
			i++
		}
		for _, p := range response.Items {
			video := Video{}
			video.Id = p.ContentDetails.VideoId
			v[video.Id] = video
		}
		pageToken = response.NextPageToken
	}
}

