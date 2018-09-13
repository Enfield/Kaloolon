package main

import (
	"google.golang.org/api/youtube/v3"
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"sync/atomic"
)

type Comment struct {
	Id                    string
	AuthorDisplayName     string
	AuthorProfileImageUrl string
	AuthorChannelUrl      string
	AuthorChannelId       string
	ChannelId             string
	VideoId               string
	ParentId              string
	CanRate               bool
	ViewerRating          string
	LikeCount             int64
	ModerationStatus      string
	PublishedAt           string
	UpdatedAt             string
}

// Save implements the ValueSaver interface.
func (i Comment) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Id":                    i.Id,
		"AuthorDisplayName":     i.AuthorDisplayName,
		"AuthorProfileImageUrl": i.AuthorProfileImageUrl,
		"AuthorChannelUrl":      i.AuthorChannelUrl,
		"AuthorChannelId":       i.AuthorChannelId,
		"ChannelId":             i.ChannelId,
		"VideoId":               i.VideoId,
		"ParentId":              i.ParentId,
		"CanRate":               i.CanRate,
		"ViewerRating":          i.ViewerRating,
		"LikeCount":             i.LikeCount,
		"ModerationStatus":      i.ModerationStatus,
		"PublishedAt":           i.PublishedAt,
		"UpdatedAt":             i.UpdatedAt,
	}, i.Id, nil
}

func (v *Video) LoadComments() {
	var commentCounter uint64 = 0
	Info.Printf("Video:   [%v] Starting processing comments\n", video.Id)
	call := v.service.CommentThreads.List("snippet").VideoId(video.Id).MaxResults(100)
	response, err := call.Do()
	i := 0
	for !HandleApiError(err) {
		if i == 5 {
			Error.Fatalf(err.Error())
		}
		response, err = call.Do()
		i++
	}
	comments := commentThreadsFromResponse(response, service, video, &commentCounter)
	go LoadCommentsToBigQuery(ctx, video.ChannelId, video.Id, comments, client)
	nextPageToken := response.NextPageToken
	for len(nextPageToken) > 0 {
		Info.Printf("Video:   [%v] Downloaded %.2f%%\n", video.Id, float64(atomic.LoadUint64(&commentCounter))/float64(video.CommentCount)*100)
		call := service.CommentThreads.List("snippet").VideoId(video.Id).MaxResults(100).PageToken(nextPageToken)
		response, err := call.Do()
		i := 0
		for !HandleApiError(err) {
			if i == 5 {
				Error.Fatalf(err.Error())
			}
			response, err = call.Do()
			i++
		}
		comments := commentThreadsFromResponse(response, service, video, &commentCounter)
		go LoadCommentsToBigQuery(ctx, video.ChannelId, video.Id, comments, client)
		nextPageToken = response.NextPageToken
	}
	Info.Printf("Video:   [%v] Downloaded 100%% Total: %d\n", video.Id, atomic.LoadUint64(&commentCounter))
}

func commentThreadsFromResponse(response *youtube.CommentThreadListResponse, service *youtube.Service, video Video, commentCounter *uint64) [] Comment {
	comments := make([]Comment, 0)
	for _, item := range response.Items {
		atomic.AddUint64(commentCounter, 1)
		comments = append(comments, comment(item.Snippet.TopLevelComment, video.Id, video.ChannelId))
		if item.Snippet.TotalReplyCount > 0 {
			call := service.Comments.List("snippet").ParentId(item.Snippet.TopLevelComment.Id).MaxResults(100)
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
				atomic.AddUint64(commentCounter, 1)
				comments = append(comments, comment(i, video.Id, video.ChannelId))
			}
			nextPageToken := response.NextPageToken
			for len(nextPageToken) > 0 {
				Info.Printf("Video:   [%v] Downloaded %.2f%%\n", video.Id, float64(atomic.LoadUint64(commentCounter))/float64(video.CommentCount)*100)
				call := service.Comments.List("snippet").ParentId(item.Snippet.TopLevelComment.Id).MaxResults(100).PageToken(nextPageToken)
				response, err := call.Do()
				i := 0
				for !HandleApiError(err) {
					if i == 5 {
						Error.Fatalf(err.Error())
					}
					response, err = call.Do()
					i++
				}
				nextPageToken = response.NextPageToken
				for _, item := range response.Items {
					atomic.AddUint64(commentCounter, 1)
					comments = append(comments, comment(item, video.Id, video.ChannelId))
				}
			}
		}
	}
	return comments
}

func comment(item *youtube.Comment, videoId string, channelId string) Comment {
	comment := Comment{}
	comment.Id = item.Id
	comment.AuthorDisplayName = item.Snippet.AuthorDisplayName
	comment.AuthorProfileImageUrl = item.Snippet.AuthorProfileImageUrl
	comment.AuthorChannelUrl = item.Snippet.AuthorChannelUrl
	comment.VideoId = videoId
	comment.ChannelId = channelId
	comment.ParentId = item.Snippet.ParentId
	comment.CanRate = item.Snippet.CanRate
	comment.ViewerRating = item.Snippet.ViewerRating
	comment.LikeCount = item.Snippet.LikeCount
	comment.ModerationStatus = item.Snippet.ModerationStatus
	comment.PublishedAt = item.Snippet.PublishedAt
	comment.UpdatedAt = item.Snippet.UpdatedAt
	return comment
}
