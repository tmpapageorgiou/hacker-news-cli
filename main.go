package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

const (
	hostAPI = "hacker-news.firebaseio.com"
)

var options = struct {
	CSV    bool   `long:"csv" description:"output in csv format"`
	Output string `long:"output" description:"output filename (default: stdout)"`
	Limit  uint8  `long:"limit" description:"number of articles to get from server" default:"20"`
}{}

func main() {

	var fatalErr error
	defer func() {
		if fatalErr != nil {
			os.Exit(-1)
		}
		os.Exit(0)
	}()

	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			log.WithError(flagsErr).Error("Failed to parse options.")
			fatalErr = err
			return
		}
	}

	api := NewHackerNewsAPI(hostAPI)

	stories, err := api.TopStories(options.Limit)
	if err != nil {
		log.WithError(err).Error("Failed to get to stories.")
		fatalErr = err
		return
	}

	output := os.Stdout
	if options.Output != "" {
		output, err = os.Create(options.Output)
		if err != nil {
			log.WithError(err).Error("Failed to get open output file.")
			fatalErr = err
			return
		}
		defer output.Close()
	}

	var formatter Formatter
	switch {
	case options.CSV:
		formatter, fatalErr = NewCSVFormatter(output)
	default:
		formatter = NewPlainText(output)
	}

	if fatalErr != nil {
		return
	}

	for _, story := range stories {
		if err := formatter.Write(story); err != nil {
			fatalErr = err
			break
		}
	}
}
