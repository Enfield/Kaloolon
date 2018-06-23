package main

import (
	"google.golang.org/api/youtube/v3"
	"strings"
)

type Comment struct {
	Id                    string
	AuthorDisplayName     string
	AuthorProfileImageUrl string
	AuthorChannelUrl      string
	AuthorChannelId       string
	ChannelId             string
	VideoId               string
	TextDisplay           string
	TextOriginal          string
	ParentId              string
	CanRate               bool
	ViewerRating          string
	LikeCount             int64
	ModerationStatus      string
	PublishedAt           string
	UpdatedAt             string
}

func commentsByVideo(service *youtube.Service, video Video) []Comment {
	comments := make([]Comment,0)
	Info.Printf("Video: [%v] Starting processing comments\n", video.Id)
	call := service.CommentThreads.List("snippet").VideoId(video.Id).MaxResults(100)
	response, err := call.Do()
	if response.HTTPStatusCode != 400 {
		handleError(err, "")
	}
	commentThreadsFromResponse(response, service, video, &comments)
	nextPageToken := response.NextPageToken
	for len(nextPageToken) > 0 {
		Info.Printf("Video: [%v] Downloaded %.2f%%\n", video.Id, float64(len(comments))/float64(video.CommentCount)*100)
		call := service.CommentThreads.List("snippet").VideoId(video.Id).MaxResults(100).PageToken(nextPageToken)
		response, err := call.Do()
		if response.HTTPStatusCode != 400 {
			handleError(err, "")
		}
		commentThreadsFromResponse(response, service, video, &comments)
		nextPageToken = response.NextPageToken
	}
	Info.Printf("Video: [%v] Downloaded 100%%. Total: %d\n", video.Id, len(comments))
	return comments
}

func commentThreadsFromResponse(response *youtube.CommentThreadListResponse, service *youtube.Service, video Video, commentsPtr *[]Comment) {
	for _, item := range response.Items {
		*commentsPtr = append(*commentsPtr, comment(item.Snippet.TopLevelComment, video.Id))
		if item.Snippet.TotalReplyCount > 0 {
			call := service.Comments.List("snippet").ParentId(item.Snippet.TopLevelComment.Id).MaxResults(100)
			response, err := call.Do()
			if response.HTTPStatusCode != 400 {
				handleError(err, "")
			}
			for _, i := range response.Items {
				*commentsPtr = append(*commentsPtr, comment(i, video.Id))
			}
			nextPageToken := response.NextPageToken
			for len(nextPageToken) > 0 {
				Info.Printf("Video: [%v] Downloaded %.2f%%\n", video.Id, float64(len(*commentsPtr))/float64(video.CommentCount)*100)
				call := service.Comments.List("snippet").ParentId(item.Snippet.TopLevelComment.Id).MaxResults(100).PageToken(nextPageToken)
				response, err := call.Do()
				if response.HTTPStatusCode != 400 {
					handleError(err, "")
				}
				nextPageToken = response.NextPageToken
				for _, item := range response.Items {
					*commentsPtr = append(*commentsPtr, comment(item, video.Id))
				}
			}
		}
	}
}

func comment(item *youtube.Comment, videoId string) Comment {
	comment := Comment{}
	comment.Id = item.Id
	comment.AuthorDisplayName = item.Snippet.AuthorDisplayName
	comment.AuthorProfileImageUrl = item.Snippet.AuthorProfileImageUrl
	comment.AuthorChannelUrl = item.Snippet.AuthorChannelUrl
	comment.VideoId = videoId
	comment.TextDisplay = strings.Replace(item.Snippet.TextDisplay, "\n","",-1)
	comment.TextOriginal = strings.Replace(item.Snippet.TextOriginal, "\n","",-1)
	comment.ParentId = item.Snippet.ParentId
	comment.CanRate = item.Snippet.CanRate
	comment.ViewerRating = item.Snippet.ViewerRating
	comment.LikeCount = item.Snippet.LikeCount
	comment.ModerationStatus = item.Snippet.ModerationStatus
	comment.PublishedAt = item.Snippet.PublishedAt
	comment.UpdatedAt = item.Snippet.UpdatedAt
	return comment
}
