package main

import "regexp"

var REGEXP_FORMAT_HISTORY = regexp.MustCompile("[^%]%i")

type Format string

func (f Format) String() string {
	return string(f)
}

func (f Format) HasHistoryNumberSpec() bool {
	return REGEXP_FORMAT_HISTORY.FindString(f.String()) != ""
}
