// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jxpath

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

// A DecimalFormat defines the symbols used in a FormatNumber
// picture string.
//
// See the XPath documentation for background.
//
// https://www.w3.org/TR/xpath-functions-31/#defining-decimal-format
type DecimalFormat struct {
	DecimalSeparator  rune
	GroupSeparator    rune
	ExponentSeparator rune
	MinusSign         rune
	Infinity          string
	NaN               string
	Percent           string
	PerMille          string
	ZeroDigit         rune
	OptionalDigit     rune
	PatternSeparator  rune
}

// NewDecimalFormat returns a new DecimalFormat object with
// the default number formatting settings.
func NewDecimalFormat() DecimalFormat {
	return DecimalFormat{
		DecimalSeparator:  '.',
		GroupSeparator:    ',',
		ExponentSeparator: 'e',
		MinusSign:         '-',
		Infinity:          "Infinity",
		NaN:               "NaN",
		Percent:           "%",
		PerMille:          "‰",
		ZeroDigit:         '0',
		OptionalDigit:     '#',
		PatternSeparator:  ';',
	}
}

// The following helper methods are designed for use with
// functions like strings.IndexFunc. Note that:
//
// 1. The methods have pointer receivers. There's significant
//    performance overhead in using value receivers as these
//    methods are potentially called once for every rune in
//    a string.
//
// 2. The methods are not passed directly to the string
//    functions as declaring method values causes the entire
//    DecimalFormat object to be allocated on the heap. This,
//    for example, is perfectly valid Go code:
//
//        df := NewDecimalFormat()
//        pos := strings.IndexFunc(s, df.isDigit)
//
//    But the method value df.isDigit causes an allocation
//    (which is more work for the garbage collector). We can
//    avoid that performance hit by using an anonymous wrapper
//    function instead:
//
//        df := NewDecimalFormat()
//        pos := strings.IndexFunc(s, func(r rune) bool {
//            return df.isDigit(r)
//        })
//
//    This code may not be as readable but it runs ~40% faster
//    in Go v1.10.3. See issue #27557 for updates:
//
//    https://github.com/golang/go/issues/27557

func (format *DecimalFormat) isZeroDigit(r rune) bool {
	return r == format.ZeroDigit
}

func (format *DecimalFormat) isDecimalDigit(r rune) bool {
	r -= format.ZeroDigit
	return r >= 0 && r <= 9
}

func (format *DecimalFormat) isDigit(r rune) bool {
	return r == format.OptionalDigit || format.isDecimalDigit(r)
}

func (format *DecimalFormat) isActive(r rune) bool {
	switch r {
	case
		format.DecimalSeparator,
		format.ExponentSeparator,
		format.GroupSeparator,
		format.PatternSeparator,
		format.OptionalDigit:
		return true
	default:
		return format.isDecimalDigit(r)
	}
}

// FormatNumber converts a number to a string, formatted according
// to the given picture string and decimal format.
//
// See the XPath function format-number for the syntax of the
// picture string.
//
// https://www.w3.org/TR/xpath-functions-31/#formatting-numbers
func FormatNumber(value float64, picture string, format DecimalFormat) (string, error) {
	if picture == "" {
		return "", fmt.Errorf("picture string cannot be empty")
	}

	vars, err := processPicture(picture, &format, value < 0)
	if err != nil {
		return "", err
	}

	if math.IsNaN(value) {
		return vars.Prefix + format.NaN + vars.Suffix, nil
	}
	if math.IsInf(value, 0) {
		return vars.Prefix + format.Infinity + vars.Suffix, nil
	}

	switch vars.NumberType {
	case typePercent:
		value *= 100
	case typePermille:
		value *= 1000
	}

	exponent := 0
	if vars.MinExponentSize != 0 {

		maxMantissa := math.Pow(10, float64(vars.ScalingFactor))
		minMantissa := math.Pow(10, float64(vars.ScalingFactor-1))

		for value < minMantissa {
			value *= 10
			exponent--
		}

		for value > maxMantissa {
			value /= 10
			exponent++
		}
	}

	var integerPart, fractionalPart, exponentPart string

	value = round(value, vars.MaxFractionalSize)
	s := makeNumberString(value, vars.MaxFractionalSize, &format)
	sint, sfrac := splitStringAtByte(s, '.')
	if sint != "" {
		integerPart = formatIntegerPart(sint, &vars, &format)
	}
	if sfrac != "" {
		fractionalPart = formatFractionalPart(sfrac, &vars, &format)
	}

	if vars.MinExponentSize != 0 {
		s := makeNumberString(float64(exponent), 0, &format)
		exponentPart = formatExponentPart(s, &vars, &format)
	}

	buf := make([]byte, 0, 128)

	buf = append(buf, vars.Prefix...)
	buf = append(buf, integerPart...)

	if len(fractionalPart) > 0 {
		buf = append(buf, string(format.DecimalSeparator)...)
		buf = append(buf, fractionalPart...)
	}

	if len(exponentPart) > 0 {
		buf = append(buf, string(format.ExponentSeparator)...)
		if exponent < 0 {
			buf = append(buf, string(format.MinusSign)...)
		}
		buf = append(buf, exponentPart...)
	}

	buf = append(buf, vars.Suffix...)

	return string(buf), nil
}

func processPicture(picture string, format *DecimalFormat, isNegative bool) (subpictureVariables, error) {

	pic1, pic2 := splitStringAtRune(picture, format.PatternSeparator)
	if pic1 == "" {
		return subpictureVariables{}, fmt.Errorf("picture string must contain 1 or 2 subpictures")
	}

	vars1, err := processSubpicture(pic1, format)
	if err != nil {
		return subpictureVariables{}, err
	}

	var vars2 subpictureVariables
	if pic2 != "" {
		vars2, err = processSubpicture(pic2, format)
		if err != nil {
			return subpictureVariables{}, err
		}
	}

	vars := vars1
	if isNegative {
		if pic2 != "" {
			vars = vars2
		} else {
			vars.Prefix = string(format.MinusSign) + vars.Prefix
		}
	}

	return vars, nil
}

func processSubpicture(subpicture string, format *DecimalFormat) (subpictureVariables, error) {

	parts := extractSubpictureParts(subpicture, format)
	err := validateSubpictureParts(parts, format)
	if err != nil {
		return subpictureVariables{}, err
	}

	return analyseSubpictureParts(parts, format), nil
}

type subpictureParts struct {
	Prefix     string
	Suffix     string
	Mantissa   string
	Exponent   string
	Integer    string
	Fractional string
	Picture    string
	Active     string
}

func extractSubpictureParts(subpicture string, format *DecimalFormat) subpictureParts {

	isActive := func(r rune) bool {
		return r != format.ExponentSeparator && format.isActive(r)
	}

	first := strings.IndexFunc(subpicture, isActive)
	if first < 0 {
		first = 0
	}

	last := strings.LastIndexFunc(subpicture, isActive)
	if last < 0 {
		last = len(subpicture)
	} else {
		_, w := utf8.DecodeRuneInString(subpicture[last:])
		last += w
	}

	prefix := subpicture[:first]
	suffix := subpicture[last:]
	activePart := subpicture[first:last]

	mantissaPart := activePart
	exponentPart := ""
	if pos := strings.IndexRune(activePart, format.ExponentSeparator); pos >= 0 {
		w := utf8.RuneLen(format.ExponentSeparator)
		mantissaPart = activePart[:pos]
		exponentPart = activePart[pos+w:]
	}

	integerPart := mantissaPart
	fractionalPart := suffix
	if pos := strings.IndexRune(mantissaPart, format.DecimalSeparator); pos >= 0 {
		w := utf8.RuneLen(format.DecimalSeparator)
		integerPart = mantissaPart[:pos]
		fractionalPart = mantissaPart[pos+w:]
	}

	return subpictureParts{
		Picture:    subpicture,
		Prefix:     prefix,
		Suffix:     suffix,
		Active:     activePart,
		Mantissa:   mantissaPart,
		Exponent:   exponentPart,
		Integer:    integerPart,
		Fractional: fractionalPart,
	}
}

func validateSubpictureParts(parts subpictureParts, format *DecimalFormat) error {

	if strings.Count(parts.Picture, string(format.DecimalSeparator)) > 1 {
		return fmt.Errorf("a subpicture cannot contain more than one decimal separator")
	}

	percents := strings.Count(parts.Picture, format.Percent)
	if percents > 1 {
		return fmt.Errorf("a subpicture cannot contain more than one percent character")
	}

	permilles := strings.Count(parts.Picture, format.PerMille)
	if permilles > 1 {
		return fmt.Errorf("a subpicture cannot contain more than one per-mille character")
	}

	if percents > 0 && permilles > 0 {
		return fmt.Errorf("a subpicture cannot contain both percent and per-mille characters")
	}

	// Passing an anonymous function to IndexFunc instead of
	// a method value prevents format escaping to the heap.
	if strings.IndexFunc(parts.Mantissa, func(r rune) bool {
		return format.isDigit(r)
	}) == -1 {
		return fmt.Errorf("a mantissa part must contain at least one decimal or optional digit")
	}

	isPassive := func(r rune) bool {
		return !format.isActive(r)
	}
	if strings.IndexFunc(parts.Active, isPassive) != -1 {
		return fmt.Errorf("a subpicture cannot contain a passive character that is both preceded by and followed by an active character")
	}

	if lastRuneInString(parts.Integer) == format.GroupSeparator ||
		firstRuneInString(parts.Fractional) == format.GroupSeparator {
		if strings.ContainsRune(parts.Picture, format.DecimalSeparator) {
			return fmt.Errorf("a group separator cannot be adjacent to a decimal separator")
		}
		return fmt.Errorf("an integer part cannot end with a group separator")
	}

	if strings.Contains(parts.Picture, doubleRune(format.GroupSeparator)) {
		return fmt.Errorf("a subpicture cannot contain adjacent group separators")
	}

	// Passing this wrapper function to IndexFunc instead of
	// a method value prevents format escaping to the heap.
	isDecimalDigit := func(r rune) bool {
		return format.isDecimalDigit(r)
	}

	pos := strings.IndexFunc(parts.Integer, isDecimalDigit)
	if pos != -1 {
		pos += utf8.RuneLen(format.ZeroDigit)
		if strings.ContainsRune(parts.Integer[pos:], format.OptionalDigit) {
			return fmt.Errorf("an integer part cannot contain a decimal digit followed by an optional digit")
		}
	}

	pos = strings.IndexRune(parts.Fractional, format.OptionalDigit)
	if pos != -1 {
		pos += utf8.RuneLen(format.OptionalDigit)
		if strings.IndexFunc(parts.Fractional[pos:], isDecimalDigit) != -1 {
			return fmt.Errorf("a fractional part cannot contain an optional digit followed by a decimal digit")
		}
	}

	exponents := strings.Count(parts.Picture, string(format.ExponentSeparator))
	if exponents > 1 {
		return fmt.Errorf("a subpicture cannot contain more than one exponent separator")
	}

	if exponents > 0 && (percents > 0 || permilles > 0) {
		return fmt.Errorf("a subpicture cannot contain a percent/per-mille character and an exponent separator")
	}

	if exponents > 0 {
		isNotDecimalDigit := func(r rune) bool {
			return !format.isDecimalDigit(r)
		}
		if strings.IndexFunc(parts.Exponent, isNotDecimalDigit) != -1 {
			return fmt.Errorf("an exponent part must consist solely of one or more decimal digits")
		}
	}

	return nil
}

type numberType uint8

const (
	_ numberType = iota
	typePercent
	typePermille
)

type subpictureVariables struct {
	NumberType               numberType
	IntegerGroupPositions    []int
	GroupSize                int
	MinIntegerSize           int
	ScalingFactor            int
	FractionalGroupPositions []int
	MinFractionalSize        int
	MaxFractionalSize        int
	MinExponentSize          int
	Prefix                   string
	Suffix                   string
}

func analyseSubpictureParts(parts subpictureParts, format *DecimalFormat) subpictureVariables {

	var typ numberType
	switch {
	case strings.Contains(parts.Picture, format.Percent):
		typ = typePercent
	case strings.Contains(parts.Picture, format.PerMille):
		typ = typePermille
	}

	// Defining these wrapper functions instead of using method
	// values prevents format escaping to the heap.
	isDigit := func(r rune) bool {
		return format.isDigit(r)
	}
	isDecimalDigit := func(r rune) bool {
		return format.isDecimalDigit(r)
	}

	integerGroupPositions := getGroupPositions(parts.Integer, format.GroupSeparator, isDigit, false)
	fractionalGroupPositions := getGroupPositions(parts.Fractional, format.GroupSeparator, isDigit, true)
	groupSize := getGroupSize(integerGroupPositions)

	minIntegerSize := runeCountInStringFunc(parts.Integer, isDecimalDigit)
	scalingFactor := minIntegerSize

	minFractionalSize := runeCountInStringFunc(parts.Fractional, isDecimalDigit)
	maxFractionalSize := runeCountInStringFunc(parts.Fractional, isDigit)

	if minIntegerSize == 0 && maxFractionalSize == 0 {
		if parts.Exponent != "" {
			minFractionalSize = 1
			maxFractionalSize = 1
		} else {
			minIntegerSize = 1
		}
	}

	if parts.Exponent != "" && minIntegerSize == 0 &&
		strings.ContainsRune(parts.Integer, format.OptionalDigit) {
		minIntegerSize = 1
	}

	if minIntegerSize == 0 && minFractionalSize == 0 {
		minFractionalSize = 1
	}

	minExponentSize := runeCountInStringFunc(parts.Exponent, isDecimalDigit)

	return subpictureVariables{
		Prefix:                   parts.Prefix,
		Suffix:                   parts.Suffix,
		NumberType:               typ,
		IntegerGroupPositions:    integerGroupPositions,
		GroupSize:                groupSize,
		MinIntegerSize:           minIntegerSize,
		ScalingFactor:            scalingFactor,
		FractionalGroupPositions: fractionalGroupPositions,
		MinFractionalSize:        minFractionalSize,
		MaxFractionalSize:        maxFractionalSize,
		MinExponentSize:          minExponentSize,
	}
}

func getGroupPositions(s string, sep rune, fn func(rune) bool, lookLeft bool) []int {

	var rest string
	var positions []int

	length := utf8.RuneLen(sep)

	for {
		pos := strings.IndexRune(s, sep)
		if pos == -1 {
			break
		}

		if lookLeft {
			rest = s[:pos]
		} else {
			rest = s[pos+length:]
		}

		positions = append(positions, runeCountInStringFunc(rest, fn))

		if lookLeft {
			if l := len(positions); l > 1 {
				positions[l-1] += positions[l-2]
			}
		}

		s = s[pos+length:]
	}

	return positions
}

func getGroupSize(positions []int) int {

	if len(positions) == 0 {
		return 0
	}

	factor := gcdOf(positions)
	for i := 0; i < len(positions); i++ {
		if indexInt(positions, factor*(i+1)) == -1 {
			return 0
		}
	}

	return factor
}

func formatIntegerPart(integer string, vars *subpictureVariables, format *DecimalFormat) string {

	// Passing an anonymous function to TrimLeftFunc instead
	// of a method value prevents format escaping to the heap.
	integer = strings.TrimLeftFunc(integer, func(r rune) bool {
		return format.isZeroDigit(r)
	})

	padding := vars.MinIntegerSize - utf8.RuneCountInString(integer)

	switch {
	case padding == 1:
		integer = string(format.ZeroDigit) + integer
	case padding > 1:
		integer = strings.Repeat(string(format.ZeroDigit), padding) + integer
	}

	if vars.GroupSize > 0 {
		return insertSeparatorsEvery(integer, format.GroupSeparator, vars.GroupSize)
	}

	if len(vars.IntegerGroupPositions) > 0 {
		return insertSeparatorsAt(integer, format.GroupSeparator, vars.IntegerGroupPositions, true)
	}

	return integer
}

func formatFractionalPart(fractional string, vars *subpictureVariables, format *DecimalFormat) string {

	// Passing an anonymous function to TrimRightFunc instead
	// of a method value prevents format escaping to the heap.
	fractional = strings.TrimRightFunc(fractional, func(r rune) bool {
		return format.isZeroDigit(r)
	})

	padding := vars.MinFractionalSize - utf8.RuneCountInString(fractional)

	switch {
	case padding == 1:
		fractional += string(format.ZeroDigit)
	case padding > 1:
		fractional += strings.Repeat(string(format.ZeroDigit), padding)
	}

	if len(vars.FractionalGroupPositions) > 0 {
		return insertSeparatorsAt(fractional, format.GroupSeparator, vars.FractionalGroupPositions, false)
	}

	return fractional
}

func formatExponentPart(exponent string, vars *subpictureVariables, format *DecimalFormat) string {

	padding := vars.MinExponentSize - utf8.RuneCountInString(exponent)

	switch {
	case padding == 1:
		exponent = string(format.ZeroDigit) + exponent
	case padding > 1:
		exponent = strings.Repeat(string(format.ZeroDigit), padding) + exponent
	}

	return exponent
}

func makeNumberString(value float64, dp int, format *DecimalFormat) string {

	s := strconv.AppendFloat(make([]byte, 0, 24), math.Abs(value), 'f', dp, 64)

	if format.ZeroDigit != '0' {
		s = bytes.Map(func(r rune) rune {
			offset := r - '0'
			if offset < 0 || offset > 9 {
				return r
			}
			return format.ZeroDigit + offset
		}, s)
	}

	return string(s)
}

func insertSeparatorsEvery(s string, sep rune, interval int) string {

	l := utf8.RuneCountInString(s)
	if interval <= 0 || l <= interval {
		return s
	}

	end := len(s)
	n := (l - 1) / interval
	chunks := make([]string, n+1)

	for n > 0 {
		pos := 0
		for i := 0; i < interval; i++ {
			_, w := utf8.DecodeLastRuneInString(s[:end])
			pos += w
		}
		chunks[n] = s[end-pos : end]
		end -= pos
		n--
	}

	chunks[n] = s[:end]
	return strings.Join(chunks, string(sep))
}

func insertSeparatorsAt(integer string, sep rune, positions []int, fromRight bool) string {

	s := integer
	chunks := make([]string, 0, len(positions)+1)

	for i := range positions {

		n := positions[i]
		if fromRight {
			n = utf8.RuneCountInString(s) - n
		}

		pos := 0
		for n > 0 {
			_, w := utf8.DecodeRuneInString(s[pos:])
			pos += w
			n--
		}

		chunks = append(chunks, s[:pos])
		s = s[pos:]
	}

	chunks = append(chunks, s)
	return strings.Join(chunks, string(sep))
}

func splitStringAtRune(s string, r rune) (string, string) {

	pos := strings.IndexRune(s, r)
	if pos == -1 {
		return s, ""
	}

	if s2 := s[pos+utf8.RuneLen(r):]; !strings.ContainsRune(s2, r) {
		return s[:pos], s2
	}

	return "", ""
}

func splitStringAtByte(s string, b byte) (string, string) {

	pos := strings.IndexByte(s, b)
	if pos == -1 {
		return s, ""
	}

	if s2 := s[pos+1:]; strings.IndexByte(s2, b) == -1 {
		return s[:pos], s2
	}

	return "", ""
}

func firstRuneInString(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

func lastRuneInString(s string) rune {
	r, _ := utf8.DecodeLastRuneInString(s)
	return r
}

func runeCountInStringFunc(s string, f func(rune) bool) int {

	var count int
	for _, c := range s {
		if f(c) {
			count++
		}
	}

	return count
}

func doubleRune(r rune) string {
	return string([]rune{r, r})
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func gcdOf(values []int) int {
	res := 0
	for _, n := range values {
		res = gcd(res, n)
	}
	return res
}

func indexInt(values []int, want int) int {
	for i, n := range values {
		if n == want {
			return i
		}
	}
	return -1
}

func round(x float64, prec int) float64 {
	// From gonum's floats.RoundEven.
	// https://github.com/gonum/gonum/tree/master/floats
	if x == 0 {
		// Make sure zero is returned
		// without the negative bit set.
		return 0
	}
	// Fast path for positive precision on integers.
	if prec >= 0 && x == math.Trunc(x) {
		return x
	}
	pow := math.Pow10(prec)
	intermed := x * pow
	if math.IsInf(intermed, 0) {
		return x
	}
	if isHalfway(intermed) {
		correction, _ := math.Modf(math.Mod(intermed, 2))
		intermed += correction
		if intermed > 0 {
			x = math.Floor(intermed)
		} else {
			x = math.Ceil(intermed)
		}
	} else {
		if x < 0 {
			x = math.Ceil(intermed - 0.5)
		} else {
			x = math.Floor(intermed + 0.5)
		}
	}

	if x == 0 {
		return 0
	}

	return x / pow
}

func isHalfway(x float64) bool {
	_, frac := math.Modf(x)
	frac = math.Abs(frac)
	return frac == 0.5 || (math.Nextafter(frac, math.Inf(-1)) < 0.5 && math.Nextafter(frac, math.Inf(1)) > 0.5)
}
