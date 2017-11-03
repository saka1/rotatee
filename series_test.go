package main

import (
	"github.com/jehiah/go-strftime"
	"testing"
	"time"
)

const (
	FORMAT = "%Y-%m-%d"
)

func Test_PeriodSeries_test(t *testing.T) {
	tm := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	series := NewPeriodSeriesWithGuess("%Y-%m-%d", tm)
	if strftime.Format(FORMAT, series.Current()) != "2009-11-10" {
		t.Fatal(series.Current())
	}
	series = series.Next()
	if strftime.Format(FORMAT, series.Current()) != "2009-11-11" {
		t.Fatal(series, series.Current())
	}
}

func Test_ConstSeries_test(t *testing.T) {
	series := NewConstSeries("hoge.log")
	if strftime.Format("hoge.log", series.Current()) != "hoge.log" {
		t.Fatal()
	}
	series = series.Next()
	if strftime.Format("hoge.log", series.Current()) != "hoge.log" {
		t.Fatal()
	}
}
