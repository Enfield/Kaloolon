# busyglide/Dockerfile
FROM golang:1.11

COPY ./src /go/src/github.com/enfield/kaloolon/src
WORKDIR /go/src/github.com/enfield/kaloolon/src

RUN go get ./
RUN go build

ENTRYPOINT ["kaloolon"]
