package main

import (
	"fmt"
	"github.com/josuehennemann/logger"
	"os"
	"time"
)

var (
	Logger *logger.Logger
	rules  *RuleList
)

const (
	HTTP_READ_TIMEOUT  = time.Second * 5
	HTTP_WRITE_TIMEOUT = time.Second * 5
)

func CheckErrorAndKillMe(e error) {
	if e != nil {
		if Logger != nil {
			Logger.Printf(logger.ERROR, "CheckErrorAndKillMe - Error [%s]", e.Error())
		} else {
			fmt.Printf("CheckErrorAndKillMe - Error [%s]", e.Error())
		}
		os.Exit(1)
	}
}

func parseUrl(url string) string{
		if string(url[0]) == "/" {
			url = url[1:]
		}
		return url
}