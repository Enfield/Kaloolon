package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"time"
	"strings"
)

func LoadCommentsToBigQuery(ctx context.Context, channelId string, videoId string, comments []Comment, client *bigquery.Client) {
	if len(comments) > 0 {
		u := client.Dataset("comments").Table("cm").Uploader()
		u.SkipInvalidRows = true
		u.TableTemplateSuffix = "_" + strings.Replace(channelId, "-", "__", -1)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		if err := u.Put(ctxWithTimeout, comments); err != nil {
			HandleApiError(err)
		}
		ctxWithTimeout.Done()
	}
}

func (p *Playlist) LoadToBigQuery(videos []Video) {
	if len(videos) > 0 {
		u := p.BigQueryClient.Dataset("videos").Table("vi").Uploader()
		u.SkipInvalidRows = true
		u.TableTemplateSuffix = "_" + strings.Replace(p.Channel.Id, "-", "__", -1)
		ctxWithTimeout, cancel := context.WithTimeout(p.Ctx, time.Minute)
		defer cancel()
		if err := u.Put(ctxWithTimeout, videos); err != nil {
			HandleApiError(err)
		}
	}
}

func (c *Channel) LoadToBigQuery() {
	u := c.BigQueryClient.Dataset("channels").Table("ch").Uploader()
	ctxWithTimeout, cancel := context.WithTimeout(c.Ctx, time.Minute)
	defer cancel()
	if err := u.Put(ctxWithTimeout, c); err != nil {
		HandleApiError(err)
	}
}
