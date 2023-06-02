// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/blues/jsonata-go/jlib/jxpath"
	"github.com/blues/jsonata-go/jtypes"
)

// 2006-01-02T15:04:05.000Z07:00
const defaultFormatTimeLayout = "[Y]-[M01]-[D01]T[H01]:[m]:[s].[f001][Z01:01t]"

var defaultParseTimeLayouts = []string{
	"[Y]-[M01]-[D01]T[H01]:[m]:[s][Z01:01t]",
	"[Y]-[M01]-[D01]T[H01]:[m]:[s][Z0100t]",
	"[Y]-[M01]-[D01]T[H01]:[m]:[s]",
	"[Y]-[M01]-[D01]",
	"[Y]",
}

// FromMillis (golint)
func FromMillis(ms int64, picture jtypes.OptionalString, tz jtypes.OptionalString) (string, error) {

	t := msToTime(ms).UTC()

	if tz.String != "" {
		loc, err := parseTimeZone(tz.String)
		if err != nil {
			return "", err
		}

		t = t.In(loc)
	}

	layout := picture.String
	if layout == "" {
		layout = defaultFormatTimeLayout
	}

	return jxpath.FormatTime(t, layout)
}

// parseTimeZone parses a JSONata timezone.
//
// The format is a "+" or "-" character, followed by four digits, the first two
// denoting the hour offset, and the last two denoting the minute offset.
func parseTimeZone(tz string) (*time.Location, error) {
	// must be exactly 5 characters
	if len(tz) != 5 {
		return nil, fmt.Errorf("invalid timezone")
	}

	plusOrMinus := string(tz[0])

	// the first character must be a literal "+" or "-" character.
	// Any other character will error.
	var offsetMultiplier int
	switch plusOrMinus {
	case "-":
		offsetMultiplier = -1
	case "+":
		offsetMultiplier = 1
	default:
		return nil, fmt.Errorf("invalid timezone")
	}

	// take the first two digits as "HH"
	hours, err := strconv.Atoi(tz[1:3])
	if err != nil {
		return nil, fmt.Errorf("invalid timezone")
	}

	// take the last two digits as "MM"
	minutes, err := strconv.Atoi(tz[3:5])
	if err != nil {
		return nil, fmt.Errorf("invalid timezone")
	}

	// convert to seconds
	offsetSeconds := offsetMultiplier * (60 * ((60 * hours) + minutes))

	// construct a time.Location based on the tz string and the offset in seconds.
	loc := time.FixedZone(tz, offsetSeconds)

	return loc, nil
}

// ToMillis (golint)
func ToMillis(s string, picture jtypes.OptionalString, tz jtypes.OptionalString) (int64, error) {
	layouts := defaultParseTimeLayouts
	if picture.String != "" {
		layouts = []string{picture.String}
	}

	// TODO: How are timezones used for parsing?

	for _, l := range layouts {
		if t, err := parseTime(s, l); err == nil {
			return timeToMS(t), nil
		}
	}

	return 0, fmt.Errorf("could not parse time %q", s)
}

var reMinus7 = regexp.MustCompile("-(0*7)")

func parseTime(s string, picture string) (time.Time, error) {
	// Go's reference time: Mon Jan 2 15:04:05 MST 2006
	refTime := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*60*60))

	layout, err := jxpath.FormatTime(refTime, picture)
	if err != nil {
		return time.Time{}, fmt.Errorf("the second argument of the toMillis function must be a valid date format")
	}

	// Replace -07:00 with Z07:00
	layout = reMinus7.ReplaceAllString(layout, "Z$1")

	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse time %q", s)
	}

	return t, nil
}

func msToTime(ms int64) time.Time {
	return time.Unix(ms/1000, (ms%1000)*int64(time.Millisecond))
}

func timeToMS(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
