package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_rotatee_simpleTeeTest(t *testing.T) {
	var stdin, stdout bytes.Buffer
	tmpdir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpdir)
	rotatee := NewRotatee(RotateeSetting{
		stdin:  &stdin,
		stdout: &stdout,
		args:   []string{filepath.Join(tmpdir, "foo.log")},
	})
	stdin.Write([]byte("abc"))
	rotatee.Start()
	if str, err := stdout.ReadString('c'); str != "abc" || err != nil {
		t.Fatalf("invalid stdout: %v, %v", str, err)
	}
}

func Test_rotatee_sizeBasedRotation(t *testing.T) {
	// Suppress logger
	origLevel := log.Level
	log.SetLevel(logrus.ErrorLevel)
	defer log.SetLevel(origLevel)

	var stdin, stdout bytes.Buffer
	tmpdir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpdir)
	rotatee := NewRotatee(RotateeSetting{
		stdin:       &stdin,
		stdout:      &stdout,
		args:        []string{filepath.Join(tmpdir, "foo%i.log")},
		maxFileSize: 3,
		historySize: 2,
	})
	stdin.Write([]byte("abcdefghi"))
	rotatee.Start()
	if str, err := stdout.ReadString('f'); str != "abcdef" || err != nil {
		t.Fatalf("invalid stdout: %v, %v", str, err)
	}
	if b, err := ioutil.ReadFile(filepath.Join(tmpdir, "foo.log")); err != nil || string(b) != "ghi" {
		t.Fatalf("invalid stdout: %v, %v", b, err)
	}
	if b, err := ioutil.ReadFile(filepath.Join(tmpdir, "foo1.log")); err != nil || string(b) != "def" {
		t.Fatalf("invalid stdout: %v, %v", b, err)
	}
	_, err := os.Stat(filepath.Join(tmpdir, "foo2.log"))
	if err == nil {
		t.Fatalf("'foo2.log' must not exist")
	}
}

func Test_rotatee_appendMode(t *testing.T) {
	// Suppress logger
	origLevel := log.Level
	log.SetLevel(logrus.ErrorLevel)
	defer log.SetLevel(origLevel)

	var stdin, stdout bytes.Buffer
	tmpdir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpdir)
	tmpfile := filepath.Join(tmpdir, "foo.log")
	err := ioutil.WriteFile(tmpfile, []byte("abc"), 0644)
	if err != nil {
		t.Fatal()
	}
	rotatee := NewRotatee(RotateeSetting{
		stdin:      &stdin,
		stdout:     &stdout,
		args:       []string{filepath.Join(tmpdir, "foo.log")},
		appendMode: true,
	})
	stdin.Write([]byte("def"))
	rotatee.Start()
	if output, err := ioutil.ReadFile(tmpfile); string(output) != "abcdef" || err != nil {
		t.Fatalf("invalid state: output = %v, err = %v", output, err)
	}
}
