// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jxpath

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type dateComponent rune

const (
	dateYear        dateComponent = 'Y'
	dateMonth       dateComponent = 'M'
	dateDay         dateComponent = 'D'
	dateDayOfYear   dateComponent = 'd'
	dateDayOfWeek   dateComponent = 'F'
	dateWeekOfYear  dateComponent = 'W'
	dateWeekOfMonth dateComponent = 'w'
	dateHour24      dateComponent = 'H'
	dateHour12      dateComponent = 'h'
	dateAMPM        dateComponent = 'P'
	dateMinute      dateComponent = 'm'
	dateSecond      dateComponent = 's'
	dateNanosecond  dateComponent = 'f'
	dateTZ          dateComponent = 'Z'
	dateTZPrefixed  dateComponent = 'z'
	dateCalendar    dateComponent = 'C'
	dateEra         dateComponent = 'E'
)

var defaultDateFormats = map[dateComponent]string{
	dateYear:        "1",
	dateMonth:       "1",
	dateDay:         "1",
	dateDayOfYear:   "1",
	dateDayOfWeek:   "n",
	dateWeekOfYear:  "1",
	dateWeekOfMonth: "1",
	dateHour24:      "1",
	dateHour12:      "1",
	dateAMPM:        "n",
	dateMinute:      "01",
	dateSecond:      "01",
	dateNanosecond:  "1",
	dateTZ:          "01:01",
	dateTZPrefixed:  "01:01",
	dateCalendar:    "n",
	dateEra:         "n",
}

type formatModifier uint8

const (
	_ formatModifier = iota
	modOrdinal
	modCardinal
	modAlphabetic
	modTraditional
)

type variableMarker struct {
	format   string
	modifier formatModifier
	minWidth int
	maxWidth int
}

var errUnsupported = errors.New("unsupported date format")

// FormatTime converts a time to a string, formatted according
// to the given picture string.
//
// See the XPath documentation for the syntax of the picture
// string.
//
// https://www.w3.org/TR/xpath-functions-31/#rules-for-datetime-formatting
func FormatTime(t time.Time, picture string) (string, error) {
	var start int
	var inMarker, doubleClosingBracket, expanded bool

	result := make([]byte, 0, 128)

	for current, r := range picture {
		if r == '[' {
			if inMarker {
				if current != start {
					return "", fmt.Errorf("open bracket inside variable marker")
				}
				inMarker = false
			} else {
				result = append(result, picture[start:current]...)
				start = current + 1
				inMarker = true
			}

			continue
		}

		if r == ']' {
			if inMarker {
				if current == start {
					return "", fmt.Errorf("empty variable marker")
				}
				s, err := expandVariableMarker(t, picture[start:current])
				if err != nil {
					return "", err
				}
				result = append(result, s...)
				start = current + 1
				inMarker = false
				expanded = true
			} else {
				if doubleClosingBracket {
					doubleClosingBracket = false
					continue
				}
				next := current + 1
				if next >= len(picture) || picture[next] != ']' {
					return "", fmt.Errorf("closing bracket outside variable marker")
				}
				doubleClosingBracket = true
				result = append(result, picture[start:current]...)
				start = next
			}

			continue
		}
	}

	if inMarker {
		return "", fmt.Errorf("unterminated variable marker")
	}

	if !expanded {
		return "", fmt.Errorf("no variable markers found")
	}

	result = append(result, picture[start:]...)
	return string(result), nil
}

func expandVariableMarker(t time.Time, s string) (string, error) {

	component, marker, err := parseVariableMarker(s)
	if err != nil {
		return "", err
	}

	var isDefaultFormat bool

	if marker.format == "" {
		marker.modifier = 0
		marker.format = defaultDateFormats[component]
		isDefaultFormat = true
	}

	repl, err := expandDateComponent(t, component, &marker)

	if err == errUnsupported && !isDefaultFormat {
		marker.modifier = 0
		marker.format = defaultDateFormats[component]
		repl, err = expandDateComponent(t, component, &marker)
	}

	return repl, err
}

var zeroVariableMarker variableMarker

func parseVariableMarker(s string) (dateComponent, variableMarker, error) {

	var format string
	var modifier formatModifier
	var minWidth, maxWidth int

	s = stripSpace(s)
	if s == "" {
		return 0, zeroVariableMarker, fmt.Errorf("empty variable marker")
	}

	if len(s) == 1 {
		return dateComponent(s[0]), zeroVariableMarker, nil
	}

	presentationModifiers, widthModifier, err := parseVariableMarkerModifiers(s[1:])
	if err != nil {
		return 0, zeroVariableMarker, err
	}

	if presentationModifiers != "" {
		format, modifier = parsePresentationModifiers(presentationModifiers)
	}

	if widthModifier != "" {
		minWidth, maxWidth, err = parseWidthModifier(widthModifier)
		if err != nil {
			return 0, zeroVariableMarker, err
		}
	}

	return dateComponent(s[0]), variableMarker{
		format:   format,
		modifier: modifier,
		minWidth: minWidth,
		maxWidth: maxWidth,
	}, nil
}

func parseVariableMarkerModifiers(s string) (string, string, error) {

	pos := strings.LastIndexByte(s, ',')
	if pos < 0 {
		return s, "", nil
	}

	presentationModifiers := s[:pos]
	widthModifier := s[pos+1:]

	if widthModifier == "" {
		return "", "", fmt.Errorf("empty width modifier")
	}

	return presentationModifiers, widthModifier, nil
}

func parsePresentationModifiers(s string) (string, formatModifier) {

	var format string
	var modifier formatModifier

	switch len := len(s); len {
	case 0:
	case 1:
		format = s
	default:
		last := len - 1
		format = s[:last]
		switch s[last] {
		case 'a':
			modifier = modAlphabetic
		case 't':
			modifier = modTraditional
		case 'c':
			modifier = modCardinal
		case 'o':
			modifier = modOrdinal
		default:
			format = s
		}
	}

	return format, modifier
}

func parseWidthModifier(s string) (int, int, error) {

	var err error
	var min, max int

	parts := strings.Split(s, "-")
	switch len(parts) {
	case 1:
		min, err = parseWidth(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid width %q: %s", parts[0], err)
		}
	case 2:
		min, err = parseWidth(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid minimum width %q: %s", parts[0], err)
		}
		max, err = parseWidth(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid maximum width %q: %s", parts[1], err)
		}
		if max < min {
			return 0, 0, fmt.Errorf("invalid width modifier %q: maximum width cannot be less than minimum width", s)
		}
	default:
		return 0, 0, fmt.Errorf("invalid width modifier %q", s)
	}

	return min, max, nil
}

func parseWidth(s string) (int, error) {

	if s == "*" {
		return 0, nil
	}

	if !isAllDigits(s) {
		return 0, fmt.Errorf("width contains illegal characters")
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("width is not an integer")
	}

	if n < 1 {
		return 0, fmt.Errorf("width cannot be less than 1")
	}

	return n, nil
}

func expandDateComponent(t time.Time, component dateComponent, marker *variableMarker) (string, error) {
	switch component {
	case dateYear:
		return formatYear(t, marker)
	case dateMonth:
		return formatMonth(t, marker)
	case dateDay:
		return formatDay(t, marker)
	case dateDayOfYear:
		return formatDayInYear(t, marker)
	case dateDayOfWeek:
		return formatDayOfWeek(t, marker)
	case dateWeekOfYear:
		return formatWeekInYear(t, marker)
	case dateWeekOfMonth:
		return formatWeekInMonth(t, marker)
	case dateHour24:
		return formatHour24(t, marker)
	case dateHour12:
		return formatHour12(t, marker)
	case dateAMPM:
		return formatAMPM(t, marker)
	case dateMinute:
		return formatMinute(t, marker)
	case dateSecond:
		return formatSecond(t, marker)
	case dateNanosecond:
		return formatNanosecond(t, marker)
	case dateTZ:
		return formatTimezoneUnprefixed(t, marker)
	case dateTZPrefixed:
		return formatTimezonePrefixed(t, marker)
	case dateCalendar:
		return formatCalendar(t, marker)
	case dateEra:
		return formatEra(t, marker)
	default:
		return "", fmt.Errorf("unknown component specifier %c", component)
	}
}

func formatYear(t time.Time, marker *variableMarker) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}

	size := marker.maxWidth
	if size <= 0 {
		if n := countDigits(marker.format); n >= 2 {
			size = n
		}
	}

	y := t.Year()
	if size > 0 {
		y = y % pow10(size)
	}

	return formatIntegerComponent(y, marker)
}

func formatMonth(t time.Time, marker *variableMarker) (string, error) {

	month := t.Month()

	if isNameFormat(marker.format) {
		names := defaultLanguage.months[month]
		return formatNameComponent(names, marker)
	}

	if isDecimalFormat(marker.format) {
		return formatIntegerComponent(int(month), marker)
	}

	return "", errUnsupported
}

func formatDay(t time.Time, marker *variableMarker) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}
	return formatIntegerComponent(t.Day(), marker)
}

func formatDayInYear(t time.Time, marker *variableMarker) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}
	return formatIntegerComponent(t.YearDay(), marker)
}

func formatDayOfWeek(t time.Time, marker *variableMarker) (string, error) {

	day := t.Weekday()

	if isNameFormat(marker.format) {
		names := defaultLanguage.days[day]
		return formatNameComponent(names, marker)
	}

	if isDecimalFormat(marker.format) {
		return formatIntegerComponent(int(day)+1, marker)
	}

	return "", errUnsupported
}

func formatWeekInYear(t time.Time, marker *variableMarker) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}

	_, w := t.ISOWeek()
	return formatIntegerComponent(w, marker)
}

func formatWeekInMonth(t time.Time, marker *variableMarker) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}
	return formatIntegerComponent(daysToWeeks(t.Day()), marker)
}

func formatHour24(t time.Time, marker *variableMarker) (string, error) {
	return formatHour(t, marker, false)
}

func formatHour12(t time.Time, marker *variableMarker) (string, error) {
	return formatHour(t, marker, true)
}

func formatHour(t time.Time, marker *variableMarker, hour12 bool) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}

	h := t.Hour()
	if hour12 && h > 12 {
		h -= 12
	}
	return formatIntegerComponent(h, marker)
}

func formatAMPM(t time.Time, marker *variableMarker) (string, error) {

	if !isNameFormat(marker.format) {
		return "", errUnsupported
	}

	names := defaultLanguage.am
	if t.Hour() >= 12 {
		names = defaultLanguage.pm
	}

	return formatNameComponent(names, marker)
}

func formatMinute(t time.Time, marker *variableMarker) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}
	return formatIntegerComponent(t.Minute(), marker)
}

func formatSecond(t time.Time, marker *variableMarker) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}
	return formatIntegerComponent(t.Second(), marker)
}

func formatNanosecond(t time.Time, marker *variableMarker) (string, error) {

	if !isDecimalFormat(marker.format) {
		return "", errUnsupported
	}

	l := utf8.RuneCountInString(marker.format)

	if l == 1 || !isAllDigits(marker.format) {
		return formatNano(t.Nanosecond(), 9), nil
	}

	return formatNano(t.Nanosecond(), l), nil
}

func formatNano(n, maxlen int) string {

	var buf [9]byte
	for start := len(buf); start > 0; {
		start--
		buf[start] = byte(n%10 + '0')
		n /= 10
	}

	if maxlen > 9 {
		maxlen = 9
	}

	return string(buf[:maxlen])
}

type tzStyle uint

const (
	_          tzStyle = iota
	tzShort            // -07
	tzLong             // -0700
	tzSplit            // -07:00
	tzMilitary         // T
	tzName             // MST
)

type tzSplitLayout struct {
	hours     string
	minutes   string
	separator string
}

var militaryOffsets = map[int]string{
	1: "A",
	2: "B",
	3: "C",
	4: "D",
	5: "E",
	6: "F",
	7: "G",
	8: "H",
	9: "I",
	// "J" is reserved for times with no timezone (i.e. local time).
	10:  "K",
	11:  "L",
	12:  "M",
	-1:  "N",
	-2:  "O",
	-3:  "P",
	-4:  "Q",
	-5:  "R",
	-6:  "S",
	-7:  "T",
	-8:  "U",
	-9:  "V",
	-10: "W",
	-11: "X",
	-12: "Y",
	0:   "Z",
}

var reTZSplit = regexp.MustCompile("^([0-9]+)([^0-9A-Za-z])([0-9]+)$")

func getTimezoneStyle(s string) (tzStyle, *tzSplitLayout) {

	if s == "Z" {
		return tzMilitary, nil
	}

	if isNameFormat(s) {
		return tzName, nil
	}

	if isAllDigits(s) {
		switch len(s) {
		case 1, 2:
			return tzShort, nil
		case 3, 4:
			return tzLong, nil
		default:
			return 0, nil
		}
	}

	matches := reTZSplit.FindAllStringSubmatch(s, -1)
	if len(matches) == 1 && len(matches[0]) == 4 {
		return tzSplit, &tzSplitLayout{
			hours:     matches[0][1],
			minutes:   matches[0][3],
			separator: matches[0][2],
		}
	}

	return 0, nil
}

func formatTimezoneUnprefixed(t time.Time, marker *variableMarker) (string, error) {
	return formatTimezone(t, marker, false)
}

func formatTimezonePrefixed(t time.Time, marker *variableMarker) (string, error) {
	return formatTimezone(t, marker, true)
}

func formatTimezone(t time.Time, marker *variableMarker, prefixed bool) (string, error) {

	var tz string
	var err error

	style, split := getTimezoneStyle(marker.format)
	isNumeric := style == tzShort || style == tzLong || style == tzSplit

	name, hours, minutes := getTimezoneInfo(t)

	switch {
	case marker.modifier == modTraditional && isNumeric && hours == 0 && minutes == 0:
		tz = "Z"
		isNumeric = false

	case style == tzShort:
		tz, err = formatTimezoneShort(hours, minutes, marker.format)

	case style == tzLong:
		tz, err = formatTimezoneLong(hours, minutes, marker.format)

	case style == tzSplit:
		tz, err = formatTimezoneSplit(hours, split.hours, minutes, split.minutes, split.separator)

	case style == tzName && name != "":
		tz, err = formatNameComponent([]string{name}, &variableMarker{
			format: marker.format,
		})

	case style == tzMilitary && minutes == 0 && hours >= -12 && hours <= 12:
		tz = militaryOffsets[hours]

	default:
		return "", errUnsupported
	}

	if err != nil {
		return "", err
	}

	if prefixed && isNumeric {
		tz = defaultLanguage.tzPrefix + tz
	}

	if marker.minWidth > 0 {
		padding := marker.minWidth - utf8.RuneCountInString(tz)
		if padding > 0 {
			tz += strings.Repeat(" ", padding)
		}
	}

	return tz, nil
}

func formatTimezoneShort(h int, m int, layout string) (string, error) {

	tz, err := formatInteger(h, layout)
	if err != nil {
		return "", err
	}

	if h >= 0 {
		tz = "+" + tz
	}

	if m != 0 {
		tz += fmt.Sprintf(":%02d", abs(m))
	}

	return tz, nil
}

func formatTimezoneLong(h int, m int, layout string) (string, error) {

	tz, err := formatInteger(h*100+m, layout)
	if err != nil {
		return "", err
	}

	if h >= 0 {
		tz = "+" + tz
	}

	return tz, nil
}

func formatTimezoneSplit(h int, layoutH string, m int, layoutM string, separator string) (string, error) {

	hh, err := formatInteger(h, layoutH)
	if err != nil {
		return "", err
	}

	mm, err := formatInteger(abs(m), layoutM)
	if err != nil {
		return "", err
	}

	tz := hh + separator + mm

	if h >= 0 {
		tz = "+" + tz
	}

	return tz, nil
}

var calendars = []string{"AD"}

func formatCalendar(t time.Time, marker *variableMarker) (string, error) {

	if !isNameFormat(marker.format) {
		return "", errUnsupported
	}
	return formatNameComponent(calendars, marker)
}

var eras = []string{"CE"}

func formatEra(t time.Time, marker *variableMarker) (string, error) {

	if !isNameFormat(marker.format) {
		return "", errUnsupported
	}
	return formatNameComponent(eras, marker)
}

func formatIntegerComponent(n int, marker *variableMarker) (string, error) {

	s, err := formatInteger(n, marker.format)
	if err != nil {
		return "", err
	}

	switch marker.modifier {
	case modOrdinal:
		s += ordinalSuffix(n)
	}

	return s, nil
}

var defaultDecimalFormat = NewDecimalFormat()

func formatInteger(n int, layout string) (string, error) {
	return FormatNumber(float64(n), layout, defaultDecimalFormat)
}

func formatNameComponent(names []string, marker *variableMarker) (string, error) {

	s := bestFittingString(names, marker.maxWidth)
	if s == "" {
		return "", fmt.Errorf("no name exists for max length %d", marker.maxWidth)
	}

	switch marker.format {
	case "N":
		s = strings.ToUpper(s)
	case "n":
		s = strings.ToLower(s)
	case "Nn":
		s = toTitle(s)
	}

	if marker.minWidth > 0 {
		padding := marker.minWidth - utf8.RuneCountInString(s)
		if padding > 0 {
			s += strings.Repeat(" ", padding)
		}
	}

	return s, nil
}

func bestFittingString(values []string, maxlen int) string {

	if len(values) == 0 {
		return ""
	}

	if maxlen <= 0 {
		return values[0]
	}

	for _, s := range values {
		if utf8.RuneCountInString(s) <= maxlen {
			return s
		}
	}

	pos := positionOfNthRune(values[0], maxlen)
	return values[0][:pos]
}

const secondsPerMinute = 60
const secondsPerHour = secondsPerMinute * 60

func getTimezoneInfo(t time.Time) (string, int, int) {
	name, secs := t.Zone()
	hours := secs / secondsPerHour
	secs = secs % secondsPerHour
	minutes := secs / secondsPerMinute
	return name, hours, minutes
}

func daysToWeeks(days int) int {
	return (days / 7) + 1
}

func ordinalSuffix(n int) string {

	mod10 := n % 10
	mod100 := n % 100

	switch {
	case mod10 == 1 && mod100 != 11:
		return "st"
	case mod10 == 2 && mod100 != 12:
		return "nd"
	case mod10 == 3 && mod100 != 13:
		return "rd"
	default:
		return "th"
	}
}

func stripSpace(s string) string {
	return strings.Map(func(r rune) rune {
		if isWhitespace(r) {
			return -1
		}
		return r
	}, s)
}

func toTitle(s string) string {
	toUpper := true
	return strings.Map(func(r rune) rune {
		if isWhitespace(r) {
			toUpper = true
			return r
		}
		if toUpper {
			toUpper = false
			return unicode.ToUpper(r)
		}
		return unicode.ToLower(r)
	}, s)
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', '\v':
		return true
	default:
		return false
	}
}

func isDecimalFormat(s string) bool {
	return strings.ContainsAny(s, "0123456789")
}

func isNameFormat(s string) bool {
	return s == "N" || s == "n" || s == "Nn"
}

func isAllDigits(s string) bool {

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return len(s) > 0
}

func countDigits(s string) int {

	n := 0
	for _, r := range s {
		if strings.ContainsRune("0123456789#", r) {
			n++
		}
	}

	return n
}

func pow10(n int) int {
	val := 1
	for i := 0; i < n; i++ {
		val *= 10
	}
	return val
}

func positionOfNthRune(s string, n int) int {

	i := 0
	for pos := range s {
		if i == n {
			return pos
		}
		i++
	}

	return -1
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
