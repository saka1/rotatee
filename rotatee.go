package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

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
	fmt.Println("started")
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
			panic("Writer goroutine IO failed")
		}
		// do copy because 'content' is shared among goroutine(s)
		content := make([]byte, len)
		copy(content, readBuf[:len])
		writeCh <- content
	}
	return
}

func main() {
	app := cli.NewApp()
	app.Name = "rotatee"
	app.Usage = "advanced tee, advanced input rotation"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("Args: %q\n", c.Args())
		rotatee := newRotatee(RotateeSetting{
			verbose: false,
		})
		rotatee.start()
		return nil
	}
	app.Run(os.Args)
}
