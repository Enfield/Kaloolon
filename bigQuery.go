package main

import (
	"cloud.google.com/go/bigquery"
	"context"
)

func loadCommentsToBigQuery(tableName string, comments []Comment, client *bigquery.Client, ctx context.Context) {
	if len(comments) > 0 {
		Info.Printf("Channel: [%v] Video: [%v] Saving data to BigQuery\n", tableName, comments[0].VideoId)
		u := client.Dataset("kolomo").Table("comments").Uploader()
		u.TableTemplateSuffix = "_" + tableName
		if err := u.Put(ctx, comments); err != nil {
			handleError(err, "Can't save comments to BigQuery")
		}
	}
}
