package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var log = logrus.New()

type Rotatee struct {
	setting RotateeSetting
}

func newRotatee(setting RotateeSetting) *Rotatee {
	return &Rotatee{setting: setting}
}

type RotateeSetting struct {
	verbose bool
}

func (r *Rotatee) start() {
	log.WithFields(logrus.Fields{"Rotatee": r}).Debug("Start rotatee")
	writeCh := make(chan []byte)
	// writer loop
	go func() {
		writer := os.Stdin
		for {
			chunk := <-writeCh
			_, err := writer.Write(chunk)
			if err != nil {
				panic("Reader goroutine IO failed")
			}
		}
	}()

	// reader loop
	reader := os.Stdin
	readBuf := make([]byte, 1024)
	for {
		len, err := reader.Read(readBuf)
		if err != nil {
			log.Panic("Writer goroutine IO failed")
		}
		// do copy because 'content' is shared among goroutine(s)
		content := make([]byte, len)
		copy(content, readBuf[:len])
		writeCh <- content
	}
}

func main() {
	// setup logger
	log.Level = logrus.DebugLevel //TODO
	log.Out = os.Stdout
	log.Debug("At main")

	app := cli.NewApp()
	app.Name = "rotatee"
	app.Usage = "advanced tee, advanced input rotation"
	app.Action = func(c *cli.Context) error {
		log.WithFields(logrus.Fields{"Args": c.Args()}).Debug("Parsed input arguments")
		rotatee := newRotatee(RotateeSetting{
			verbose: false,
		})
		rotatee.start()
		return nil
	}
	app.Run(os.Args)
}
