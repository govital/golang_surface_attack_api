package logger

import "log"

type HandlerLogging struct{}

func (l *HandlerLogging) Init() {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
