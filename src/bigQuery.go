package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"strings"
	"time"
)

// LoadToBigQuery load video info to BigQuery
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

// LoadToBigQuery load playlist info to BigQuery
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

// LoadToBigQuery load channel info to BigQuery
func (c *Channel) LoadToBigQuery() {
	u := c.BigQueryClient.Dataset("channels").Table("ch").Uploader()
	ctxWithTimeout, cancel := context.WithTimeout(c.Ctx, time.Minute)
	defer cancel()
	if err := u.Put(ctxWithTimeout, c); err != nil {
		HandleApiError(err)
	}
}

// IsLoadedToBigQuery checks is info about youtube channel is already processed
func IsLoadedToBigQuery(ctx context.Context, channelId string, c *bigquery.Client) bool {
	iter, err := c.Query(fmt.Sprintf("select 1 from `channels.ch` where id = '%s' limit 1", channelId)).Read(ctx)
	if err != nil {
		HandleApiError(err)
	}
	isAlreadyLoded := iter.TotalRows > 0
	if isAlreadyLoded {
		Info.Printf("Channel:[%v] Already processed", channelId)
	}
	return isAlreadyLoded
}
