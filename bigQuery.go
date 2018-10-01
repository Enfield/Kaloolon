package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"strings"
	"time"
)

func (v *Video) LoadToBigQuery(comments []Comment) {
	if len(comments) > 0 {
		u := v.BigQueryClient.Dataset("comments").Table("cm").Uploader()
		u.SkipInvalidRows = true
		u.TableTemplateSuffix = "_" + strings.Replace(v.Plist.Channel.Id, "-", "__", -1)
		ctxWithTimeout, cancel := context.WithTimeout(v.Ctx, time.Minute)
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

func IsLoadedToBigQuery(ctx context.Context, channelId string, c *bigquery.Client) bool {
	iter, err :=  c.Query("select 1 from `channels.ch` where id = '" + channelId + "'").Read(ctx)
	if err != nil {
		Error.Print(err)
	}
	isAlreadyLoded := iter.TotalRows > 0
	if isAlreadyLoded {
		Info.Printf("Channel:[%v] Already processed", channelId)
	}
	return isAlreadyLoded
}
