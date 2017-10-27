package main

import (
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
		rotatee := NewRotatee(RotateeSetting{
			args:    c.Args(),
			verbose: verbose,
		})
		rotatee.Start()
		return nil
	}
	app.Run(os.Args)
}
