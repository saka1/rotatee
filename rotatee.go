package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"time"
)

var log = logrus.New()

type Rotatee struct {
	setting RotateeSetting
}

func newRotatee(setting RotateeSetting) *Rotatee {
	return &Rotatee{setting: setting}
}

type RotateeSetting struct {
	args    []string
	verbose bool
}

func setupEventPipe(arg string) EventPipe {
	pipe := NewEventPipe()
	pipe.Add(NewTimer(DetectSeries(arg, time.Now())))
	pipe.Add(NewFormatEval(arg))
	pipe.Add(NewRoller())
	pipe.Add(NewWriter())
	return pipe
}

func (r *Rotatee) start() {
	log.WithFields(logrus.Fields{"Rotatee": r}).Debug("Start rotatee")
	// init pipe
	pipeGroup := NewEventPipeGroup()
	for _, arg := range r.setting.args {
		pipeGroup.Add(setupEventPipe(arg))
	}
	pipeGroup.Start()
	// init first destination
	pipeGroup.Broadcast(NewWriteTarget())
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
		pipeGroup.Broadcast(NewPayload(content))
		// Also, write to stdout
		os.Stdout.Write(content)
	}
}

func setupLogger(verbose bool, debug bool) {
	log.Out = os.Stderr
	if verbose {
		log.Level = logrus.InfoLevel
	} else {
		log.Level = logrus.WarnLevel
	}
	if debug {
		log.Level = logrus.DebugLevel
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "rotatee"
	app.Usage = "advanced tee, advanced input rotation"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose logging to stderr",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug mode (very verbose logging to stderr)",
		},
	}
	app.Action = func(c *cli.Context) error {
		verbose, debug := c.Bool("verbose"), c.Bool("debug")
		setupLogger(verbose, debug)
		log.WithFields(logrus.Fields{"Args": c.Args()}).Debug("Parsed input arguments")
		rotatee := newRotatee(RotateeSetting{
			args:    c.Args(),
			verbose: verbose,
		})
		rotatee.start()
		return nil
	}
	app.Run(os.Args)
}
