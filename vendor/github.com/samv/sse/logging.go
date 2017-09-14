package sse

import (
	"io"
	"io/ioutil"
	"log"
)

// Logger is a way to enable logging for this module if you want it;
// approach borrowed with thanks from github.com/Shopify/sarama
var Logger StdLogger = log.New(ioutil.Discard, "[sse] ", log.LstdFlags)

// StdLogger is used to log error messages.
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	SetOutput(w io.Writer)
}
