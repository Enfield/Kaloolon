package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"time"
	"strings"
)

func loadCommentsToBigQuery(channelId string, comments *[]Comment, client *bigquery.Client, ctx context.Context) {
	c := *comments
	if len(c) > 0 {
		Info.Printf("Channel: [%v] Video: [%v] Saving comments to BigQuery\n", channelId, c[0].VideoId)
		u := client.Dataset("comments").Table("cm").Uploader()
		u.TableTemplateSuffix = "_" + strings.Replace(channelId, "-", "__", -1)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		if err := u.Put(ctxWithTimeout, *comments); err != nil {
			handleError(err, "Can't save comments to BigQuery")
		}
	}
}

func loadVideosToBigQuery(channelId string, videosMap *map[string]Video, client *bigquery.Client, ctx context.Context) {
	if len(*videosMap) > 0 {
		Info.Printf("Channel: [%v] Saving videos info to BigQuery\n", channelId)
		videos := make([]Video, len(*videosMap))
		idx := 0
		for _, value := range *videosMap {
			videos[idx] = value
			idx++
		}
		u := client.Dataset("videos").Table("vi").Uploader()
		u.TableTemplateSuffix = "_" + strings.Replace(channelId, "-", "__", -1)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		if err := u.Put(ctxWithTimeout, videos); err != nil {
			handleError(err, "Can't save videos to BigQuery")
		}
	}
}

func loadChannelsToBigQuery(channel *Channel, client *bigquery.Client, ctx context.Context) {
	Info.Printf("Channel: [%v] Saving data to BigQuery\n", channel.Id)
	u := client.Dataset("channels").Table("ch").Uploader()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	if err := u.Put(ctxWithTimeout, channel); err != nil {
		handleError(err, "Can't save comments to BigQuery")
	}
}
