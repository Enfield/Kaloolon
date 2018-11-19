package main

import (
	"io"
	"log"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// InitLogger init logger
func InitLogger(logHandle io.Writer) {

	Debug = log.New(logHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(logHandle,
		"INFO: ",
		log.Ldate|log.Ltime)

	Warning = log.New(logHandle,
		"WARNING: ",
		log.Ldate|log.Ltime)

	Error = log.New(logHandle,
		"ERROR: ",
		log.Ldate|log.Ltime)
}
