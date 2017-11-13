package main

import (
	"errors"
	"math"
	"regexp"
	"time"
)

var (
	unitRegexp   = regexp.MustCompile("%[SMHdmy]")
	unitScoreMap = map[string]int{
		"S": 1,
		"M": 60,
		"H": 60 * 60,
		"d": 60 * 60 * 24,
		"m": 60 * 60 * 24 * 30,
		"y": 60 * 60 * 24 * 365,
	}
)

//
// Detect which series to use
//
func DetectSeries(format string, t0 time.Time) Series {
	if unitRegexp.MatchString(format) {
		return NewPeriodSeriesWithGuess(format, t0)
	}
	return NewConstSeries(format)
}

/*
 * Series is a infinite time.Time list that generated from strftime's pattern sucn as:
 * "hoge-%S" => ["hoge-13", "hoge-14", "hoge-15", ...]
 * "%y-%m-%d" => ["2016-2-1", "2016-2-2", ...]
 * Iteration starts at `t0`.
 */
type Series interface {
	Current() time.Time
	Next() Series
	Sub(t time.Time) time.Duration
}

type PeriodSeries struct {
	time   time.Time
	format string
	period time.Duration
}

func NewPeriodSeriesWithGuess(format string, t0 time.Time) Series {
	period, err := guessPeriod(format)
	if err != nil {
		panic("todo impl") //TODO
	}
	return PeriodSeries{t0, format, period}
}

func (s PeriodSeries) Current() time.Time {
	return s.time
}

func (s PeriodSeries) Next() Series {
	return PeriodSeries{
		time:   s.time.Add(s.period),
		format: s.format,
		period: s.period,
	}
}

func (s PeriodSeries) Sub(t time.Time) time.Duration {
	return s.Current().Sub(t)
}

type ConstSeries struct {
	format string
}

func NewConstSeries(format string) Series {
	return ConstSeries{format}
}

func (s ConstSeries) Next() Series {
	return NewConstSeries(s.format)
}

func (s ConstSeries) Current() time.Time {
	return time.Now()
}

func (s ConstSeries) Sub(t time.Time) time.Duration {
	return time.Duration(math.MaxInt64)
}

/*
 * Guess period in second from the argument.
 * In general, the most smallest pattern letter is used.
 */
func guessPeriod(format string) (time.Duration, error) {
	re := unitRegexp
	units := re.FindAllString(format, -1)
	score := unitScoreMap
	// find minimum unit
	minScore := math.MaxInt32
	for _, u := range units {
		// u[1:] is truncation '%'
		value, ok := score[u[1:]]
		if ok && value < minScore {
			minScore = value
		}
	}
	if minScore == math.MaxInt32 {
		return time.Duration(0), errors.New("guessPeriod: fail to period detection from '" + format + "'")
	}
	// convert sec to nanosec
	return time.Duration(minScore * 1000 * 1000 * 1000), nil
}
