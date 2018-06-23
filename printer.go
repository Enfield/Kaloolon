package main

import "fmt"

func printAllVideos(videos map[string]Video) {
	for _, video := range videos {
		fmt.Printf("[%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v][%v]\n",
			video.Id,
			video.ChannelId,
			video.CategoryId,
			video.PublishedAt,
			video.Title,
			video.Description,
			video.LiveBroadcastContent,
			video.DefaultLanguage,
			video.DefaultAudioLanguage,
			video.Duration,
			video.Dimension,
			video.Definition,
			video.Caption,
			video.LicensedContent,
			video.Projection,
			video.HasCustomThumbnail,
			video.ViewCount,
			video.LikeCount,
			video.DislikeCount,
			video.FavoriteCount,
			video.CommentCount)
	}
	fmt.Printf("\n\n")
}


func printAllComments(comments []Comment) {
	for _, comment := range comments {
		fmt.Printf("CommentID [%v] ParentId [%v]\n", comment.Id, comment.ParentId)
	}
	fmt.Printf("Count: %v", len(comments))
}
