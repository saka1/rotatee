package main

import (
	"github.com/jehiah/go-strftime"
	"regexp"
	"strconv"
	"time"
)

var (
	formatHistorySpecRegexp    = regexp.MustCompile("[^%]%i")
	formatHistoryCaptureRegexp = regexp.MustCompile("([^%])%i")
)

type Format string

func (f Format) String() string {
	return string(f)
}

func (f Format) HasHistoryNumberSpec() bool {
	return formatHistorySpecRegexp.FindString(f.String()) != ""
}

func (f Format) evalHistory(t time.Time, history int) string {
	r := formatHistoryCaptureRegexp
	if history == 0 {
		str := r.ReplaceAllString(f.String(), "${1}")
		return strftime.Format(str, t)
	}
	str := r.ReplaceAllString(f.String(), "${1}"+strconv.FormatInt(int64(history), 10))
	return strftime.Format(str, t)
}
