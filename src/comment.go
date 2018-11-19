package main

import (
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/youtube/v3"
	"sync"
	"sync/atomic"
)

// Comment struct represents info about youtube comment
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

//LoadComments load all comments for provided video
func (v *Video) LoadComments(wg *sync.WaitGroup) {
	var commentCounter uint64
	//Info.Printf("Channel:[%v] Playlist:[%v] Video: [%v] Starting processing comments\n", v.Plist.Channel.Id, v.Plist.PlaylistId, v.Id)
	call := v.YouTubeService.CommentThreads.List("snippet").VideoId(v.Id).MaxResults(100)
	response, err := call.Do()
	i := 0
	for !HandleApiError(err) {
		if i == 5 {
			Error.Fatalf(err.Error())
		}
		response, err = call.Do()
		i++
	}
	comments := v.commentThreadsFromResponse(response, &commentCounter)
	wg.Add(1)
	go func() {
		defer wg.Done()
		v.LoadToBigQuery(comments)
	}()
	nextPageToken := response.NextPageToken
	for len(nextPageToken) > 0 {
		//Info.Printf("Channel:[%v] Playlist:[%v] Video: [%v] Downloaded %.2f%%\n",
		//	v.Plist.Channel.Id, v.Plist.PlaylistId, v.Id,
		//	float64(atomic.LoadUint64(&commentCounter))/float64(v.CommentCount)*100)
		call := v.YouTubeService.CommentThreads.List("snippet").VideoId(v.Id).MaxResults(100).PageToken(nextPageToken)
		response, err := call.Do()
		i := 0
		for !HandleApiError(err) {
			if i == 5 {
				Error.Fatalf(err.Error())
			}
			response, err = call.Do()
			i++
		}
		comments := v.commentThreadsFromResponse(response, &commentCounter)
		wg.Add(1)
		go func() {
			defer wg.Done()
			v.LoadToBigQuery(comments)
		}()
		nextPageToken = response.NextPageToken
	}
	//Info.Printf("Channel:[%v] Playlist:[%v] Video: [%v] Downloaded 100%% Total: %d\n",
	//	v.Plist.Channel.Id, v.Plist.PlaylistId, v.Id, atomic.LoadUint64(&commentCounter))
}

func (v *Video) commentThreadsFromResponse(response *youtube.CommentThreadListResponse, commentCounter *uint64) []Comment {
	comments := make([]Comment, 0)
	for _, item := range response.Items {
		atomic.AddUint64(commentCounter, 1)
		comments = append(comments, v.Comment(item))
		if item.Snippet.TotalReplyCount > 0 {
			call := v.YouTubeService.Comments.List("snippet").ParentId(item.Snippet.TopLevelComment.Id).MaxResults(100)
			response, err := call.Do()
			i := 0
			for !HandleApiError(err) {
				if i == 5 {
					Error.Fatalf(err.Error())
				}
				response, err = call.Do()
				i++
			}
			appendComments(response, commentCounter, comments, v)
			nextPageToken := response.NextPageToken
			for len(nextPageToken) > 0 {
				//Info.Printf("Channel:[%v] Playlist:[%v] Video: [%v] Downloaded %.2f%%\n",
				//	v.Plist.Channel.Id, v.Plist.PlaylistId, v.Id,
				//	float64(atomic.LoadUint64(commentCounter))/float64(v.CommentCount)*100)
				call := v.YouTubeService.Comments.List("snippet").ParentId(item.Snippet.TopLevelComment.Id).MaxResults(100).PageToken(nextPageToken)
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
				appendComments(response, commentCounter, comments, v)
			}
		}
	}
	return comments
}

func (v *Video) Comment(item *youtube.CommentThread) Comment {
	return Comment{
		Id:                    item.Id,
		AuthorDisplayName:     item.Snippet.TopLevelComment.Snippet.AuthorDisplayName,
		AuthorProfileImageUrl: item.Snippet.TopLevelComment.Snippet.AuthorProfileImageUrl,
		AuthorChannelUrl:      item.Snippet.TopLevelComment.Snippet.AuthorChannelUrl,
		VideoId:               v.Id,
		ChannelId:             v.Plist.Channel.Id,
		ParentId:              item.Snippet.TopLevelComment.Snippet.ParentId,
		CanRate:               item.Snippet.TopLevelComment.Snippet.CanRate,
		ViewerRating:          item.Snippet.TopLevelComment.Snippet.ViewerRating,
		LikeCount:             item.Snippet.TopLevelComment.Snippet.LikeCount,
		ModerationStatus:      item.Snippet.TopLevelComment.Snippet.ModerationStatus,
		PublishedAt:           item.Snippet.TopLevelComment.Snippet.PublishedAt,
		UpdatedAt:             item.Snippet.TopLevelComment.Snippet.UpdatedAt,
	}
}

func appendComments(response *youtube.CommentListResponse, commentCounter *uint64, comments []Comment, v *Video) {
	for _, item := range response.Items {
		atomic.AddUint64(commentCounter, 1)
		comments = append(comments, Comment{
			Id:                    item.Id,
			AuthorDisplayName:     item.Snippet.AuthorDisplayName,
			AuthorProfileImageUrl: item.Snippet.AuthorProfileImageUrl,
			AuthorChannelUrl:      item.Snippet.AuthorChannelUrl,
			VideoId:               v.Id,
			ChannelId:             v.Plist.Channel.Id,
			ParentId:              item.Snippet.ParentId,
			CanRate:               item.Snippet.CanRate,
			ViewerRating:          item.Snippet.ViewerRating,
			LikeCount:             item.Snippet.LikeCount,
			ModerationStatus:      item.Snippet.ModerationStatus,
			PublishedAt:           item.Snippet.PublishedAt,
			UpdatedAt:             item.Snippet.UpdatedAt,
		})
	}
}
