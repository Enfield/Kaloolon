package main

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
)

var (
	developerKey = "AIzaSyArihz4MAJjQTVN7Qd73MX-LD8e8x9msXY"
)

func setupRouter() *gin.Engine {
	projectID := "youtube-analyzer-206211"
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	//r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, x-client-token",
		ExposedHeaders:  "",
		MaxAge:          60 * time.Minute,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	r.POST("/channels/:channelId", func(ctx *gin.Context) {
		channelId := ctx.Param("channelId")
		bigQueryClient, err := bigquery.NewClient(ctx, projectID)
		if err != nil {
			HandleFatalError(err, "Initialize error")
		}
		if !IsLoadedToBigQuery(ctx, channelId, bigQueryClient) {
			youTubeClient := &http.Client{
				Transport: &transport.APIKey{Key: developerKey},
			}
			youTubeService, err := youtube.New(youTubeClient)
			if err != nil {
				HandleFatalError(err, "Initialize error")
			}
			processor := NewProcessor(bigQueryClient, youTubeService, ctx)
			cCp := ctx.Copy()
			ch := make(chan int)
			go processor.ProcessChannel(cCp, channelId, ch)
			<-ch
		}
		ctx.JSON(http.StatusOK, gin.H{"channelId": channelId})
	})
	return r
}

func main() {
	InitLogger(os.Stdout)
	r := setupRouter()
	r.Run(":8081")
}
