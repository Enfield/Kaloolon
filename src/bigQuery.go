package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
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

type ChannelVideosCountAndStatusResult struct {
	Status   string
	VideosCount uint64
}

func GetChannelVideosCountAndStatus(ctx context.Context, channelId string, c *bigquery.Client) (*ChannelVideosCountAndStatusResult, error) {
	data, err := c.Dataset("videos").Table("vi_" + strings.Replace(channelId, "-", "__", -1)).Metadata(ctx)
	if err != nil {
		HandleApiError(err)
	}
	iter, err := c.Query(fmt.Sprintf("select status from `channels.ch` where id = '%s'", channelId)).Read(ctx)
	if err != nil {
		HandleApiError(err)
	}
	var status string
	err = iter.Next(&status)
	if err != nil {
		return nil, err
	}
	return &ChannelVideosCountAndStatusResult{Status: status, VideosCount: data.NumRows}, nil
}
