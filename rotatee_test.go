package main

import (
	"bytes"
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
