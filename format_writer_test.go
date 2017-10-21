package main

import (
	"bytes"
	"io"
	//"os"
	"testing"
)

func mockWriter() (FormatWriterCtx, *bytes.Buffer) {
	buff := new(bytes.Buffer)
	return FormatWriterCtx{writerGen: func(fmt string) io.Writer { return buff }}, buff
}

func Test_startFormatWriter(t *testing.T) {
	fmt := ""
	ch := make(chan Event, 100)
	ctx, buff := mockWriter()
	ch <- NewEvent([]byte("hello"))
	close(ch)
	startFormatWriter(fmt, ch, ctx)
	if buff.String() != "hello" {
		t.Errorf("failed: '%x'", buff.String())
	}
}

func Test_e2e(t *testing.T) {
	ch := make(chan Event, 100) // set enough capacity to avoid blocking
	ctx := newFormatWriterCtx()
	ch <- NewEvent([]byte("hello"))
	ch <- NewEvent([]byte("hello"))
	ch <- NewEvent([]byte("hello"))
	close(ch)
	startFormatWriter("/tmp/hoge.log", ch, ctx)
	//defer os.Remove("/tmp/hoge.log")
}

func Test_defaultFormatGen(t *testing.T) {

}
