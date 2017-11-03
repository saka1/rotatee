package main

import (
	"regexp"
	"strconv"
)

var (
	formatHistorySpecRegexp   = regexp.MustCompile("[^%]%i")
	formatHistoryCaptureRegxp = regexp.MustCompile("([^%])%i")
)

type Format string

func (f Format) String() string {
	return string(f)
}

func (f Format) HasHistoryNumberSpec() bool {
	return formatHistorySpecRegexp.FindString(f.String()) != ""
}

func (f Format) evalHistory(history int) string {
	r := formatHistoryCaptureRegxp
	if history == 0 {
		return r.ReplaceAllString(f.String(), "${1}")
	}
	return r.ReplaceAllString(f.String(), "${1}"+strconv.FormatInt(int64(history), 10))
}
