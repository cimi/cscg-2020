package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func setupLogging() {
	ts := time.Now()
	var path = fmt.Sprintf("logs/%s", "debug.log") // ts.Format("2006-01-02-15-04-05"))
	var logfile, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer logfile.Close()
	// wrt := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(logfile)
	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	// log.SetOutput(os.Stdout)
	// log.SetPrefix("mazerunner")
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Started at %s", ts)
}
