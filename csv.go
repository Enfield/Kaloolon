package main

import (
	"github.com/mohae/struct2csv"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

const videosFileName = "videos.csv"
const commentsFolderName = "comments"

func mkDir(name string) {
	path := filepath.Join(".", name)
	err := os.MkdirAll(path, os.ModePerm)
	handleError(err, "Can't create folder")
}

func videos2csv(videosMap *map[string]Video, path string) {
	videos := make([]Video, len(*videosMap))
	idx := 0
	for _, value := range *videosMap {
		videos[idx] = value
		idx++
	}
	buff := &bytes.Buffer{}
	w := struct2csv.NewWriter(buff)
	err := w.WriteStructs(videos)
	if err != nil {
		handleError(err, "Can't save videos to csv file")
	}
	var filePath string
	if len(path) > 0 {
		filePath = path + string(os.PathSeparator) + videosFileName
		mkDir(path)
		Info.Printf("Channel: [%v] Saving videos info to [%v]\n", path, filePath)
	} else {
		filePath = videosFileName
		Info.Printf("Saving videos info to [%v]\n", filePath)
	}
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.Write(buff.Bytes()); err != nil {
		panic(err)
	}
}

func comments2csv(comments []Comment, path string) {
	var filePath string
	if len(comments) > 0 {
		if len(path) > 0 {
			mkDir(path + string(os.PathSeparator) + commentsFolderName)
			filePath = path + string(os.PathSeparator) + commentsFolderName + string(os.PathSeparator) + comments[0].VideoId + ".csv"
		} else {
			mkDir(commentsFolderName)
			filePath = commentsFolderName + string(os.PathSeparator) + comments[0].VideoId + ".csv"
		}
		Info.Printf("Video:   [%v] Saving comments to file [%v]\n", comments[0].VideoId, filePath)
		buff := &bytes.Buffer{}
		w := struct2csv.NewWriter(buff)
		err := w.WriteStructs(comments)
		if err != nil {
			handleError(err, "Can't save comments to csv file")
		}
		ioutil.WriteFile(filePath, buff.Bytes(), 0644)
	}
}
