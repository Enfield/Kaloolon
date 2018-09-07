package main

import (
	"bytes"
	"github.com/mohae/struct2csv"
	"io/ioutil"
)

func loadCommentsToBigQuery(comments []Comment) {
	var filePath string
	if len(comments) > 0 {
		Info.Printf("Video:   [%v] Saving comments to biqQuery\n", comments[0].VideoId)
		if err != nil {
			handleError(err, "Can't save comments to csv file")
		}
		ioutil.WriteFile(filePath, buff.Bytes(), 0644)
	}
}
