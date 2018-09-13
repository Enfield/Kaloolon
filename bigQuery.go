package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"time"
	"strings"
)

func loadCommentsToBigQuery(ctx context.Context, channelId string, videoId string, comments []Comment, client *bigquery.Client) {
	if len(comments) > 0 {
		u := client.Dataset("comments").Table("cm").Uploader()
		u.SkipInvalidRows = true
		u.TableTemplateSuffix = "_" + strings.Replace(channelId, "-", "__", -1)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		if err := u.Put(ctxWithTimeout, comments); err != nil {
			handleApiError(err)
		}
		ctxWithTimeout.Done()
	}
}

func loadVideosToBigQuery(ctx context.Context, videos []Video, channelId string, client *bigquery.Client) {
	if len(videos) > 0 {
		u := client.Dataset("videos").Table("vi").Uploader()
		u.SkipInvalidRows = true
		u.TableTemplateSuffix = "_" + strings.Replace(channelId, "-", "__", -1)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		if err := u.Put(ctxWithTimeout, videos); err != nil {
			handleApiError(err)
		}
	}
}

func loadChannelsToBigQuery(ctx context.Context, channel *Channel, client *bigquery.Client) {
	u := client.Dataset("channels").Table("ch").Uploader()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	if err := u.Put(ctxWithTimeout, channel); err != nil {
		handleApiError(err)
	}
}
