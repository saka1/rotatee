package main

import (
	"github.com/alecthomas/units"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

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
			Usage: "Verbose logging to stderr",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug mode (very verbose logging to stderr)",
		},
		cli.StringFlag{
			Name:  "s, size",
			Usage: "Max file size",
		},
		cli.IntFlag{
			Name: "history",
			Usage: "Limit number of files to keep. " +
				"After file rotation, rotatee remove the oldest file if the count are exceeded",
		},
		cli.BoolFlag{
			Name:  "a, appendMode",
			Usage: "Open output files with appendMode mode",
		},
	}
	app.Version = APP_VERSION
	app.Action = func(c *cli.Context) error {
		log.WithFields(logrus.Fields{"Args": c.Args()}).Debug("Parsed input arguments")
		verbose, debug := c.Bool("verbose"), c.Bool("debug")
		setupLogger(verbose, debug)
		maxFileSize, err := units.ParseBase2Bytes("0B")
		if c.GlobalIsSet("size") {
			maxFileSize, err = units.ParseBase2Bytes(c.String("size"))
			if err != nil {
				log.Error("invalid size format")
				cli.ShowAppHelp(c)
				os.Exit(1)
			}
		}
		historySize := c.GlobalInt("history")
		appendMode := c.GlobalBool("appendMode")
		rotatee := NewRotatee(RotateeSetting{
			args:        c.Args(),
			verbose:     verbose,
			maxFileSize: int64(maxFileSize),
			historySize: historySize,
			appendMode:  appendMode,
		})
		rotatee.Start()
		return nil
	}
	app.Run(os.Args)
}
