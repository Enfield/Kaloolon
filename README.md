[![Build Status](https://travis-ci.com/Enfield/Kaloolon.svg?branch=master)](https://travis-ci.com/Enfield/Kaloolon)
[![Go Report Card](https://goreportcard.com/badge/github.com/Enfield/Kaloolon)](https://goreportcard.com/report/github.com/Enfield/Kaloolon)
[![Maintainability](https://api.codeclimate.com/v1/badges/fb2f399b5f5391538ee2/maintainability)](https://codeclimate.com/github/Enfield/Kaloolon/maintainability)
# Kaloolon
It's my pet project that i want to open source.
A simple web service to collect all comments on all videos from YouTube channel.
Very fast. Carefully spending YouTube Data API quota.
Kaloolon written entirely in Go and intensively uses the features of Google Cloud Platform:
* [BigQuery](https://cloud.google.com/bigquery/) for data storage
* [Compute Engine](https://cloud.google.com/compute/) for hosts

Easily can be packaged to Docker container and loaded as POD to [Kubernetes Engine](https://cloud.google.com/kubernetes-engine/)

Kaloolon effectively uses memory and CPU resources and can work without problems on cheapest and slowest GCP host: [f1-micro (0.2 vCPU, 0.6 GB memory](https://cloud.google.com/compute/docs/machine-types)
On f1-micro host Kaloolon can completely extract data for channel like [BroScienceLife](https://www.youtube.com/channel/UCduKuJToxWPizJ7I2E6n1kA) with:
* 2,200,267 subscribers
* 130 videos
* ~ 200000 comments

and save it to BigQuery less then 2 minutes. This time can be dropped dramatically (to seconds) if you use hosts faster than f1-micro.

All data stored in BigQuery with provided schemas:

**Channels**
```plsql
Id:STRING,
Title:STRING,
Description:STRING,
Thumbnail:STRING,
Status:STRING,
ViewCount:STRING,
SubscriberCount:STRING
```
**Videos**
```plsql
Id:STRING,
ChannelId:STRING,
CategoryId:STRING,
PublishedAt:TIMESTAMP,
Title:STRING,
LiveBroadcastContent:STRING,
DefaultLanguage:STRING,
DefaultAudioLanguage:STRING,
Duration:STRING,
Dimension:STRING,
Definition:STRING,
Caption:BOOLEAN,
LicensedContent:BOOLEAN,
Projection:STRING,
HasCustomThumbnail:BOOLEAN,
ViewCount:STRING,
LikeCount:STRING,
DislikeCount:STRING,
FavoriteCount:STRING,
CommentCount:STRING
```
**Comments**
```plsql
Id:STRING,
AuthorDisplayName:STRING,
AuthorProfileImageUrl:STRING,
AuthorChannelUrl:STRING,
AuthorChannelId:STRING,
ChannelId:STRING,
VideoId:STRING,
ParentId:STRING,
CanRate:BOOLEAN,
ViewerRating:STRING,
LikeCount:INTEGER,
ModerationStatus:STRING,
PublishedAt:TIMESTAMP,
UpdatedAt:TIMESTAMP
```
## Installation
```sh
git clone https://github.com/Enfield/Kaloolon
cd Kaloolon/src/
go get ./
go build -o kaloolon
./kaloolon
```
## Usage example
Api is very simple and provide only one REST method.
### Extract data for channel
```http
POST /channels/UCPDXXXJj9nax0fr0Wfc048g HTTP/1.1
Host: 127.0.0.1:8080
Accept: */*
Content-Length: 0

HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Vary: Origin
Date: Mon, 19 Nov 2018 14:55:27 GMT
Content-Length: 0
```
---
Kirill Suslov – kirill@suslov.pro  
[https://github.com/Enfield](https://github.com/Enfield/)  
Distributed under the MIT license. See ``LICENSE`` for more information.  
