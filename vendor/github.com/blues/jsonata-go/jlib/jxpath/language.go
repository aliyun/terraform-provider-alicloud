// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jxpath

import (
	"time"
)

type dateLanguage struct {
	days     [7][]string
	months   [13][]string
	am       []string
	pm       []string
	tzPrefix string
}

var dateLanguages = map[string]dateLanguage{
	"en": {
		days: [...][]string{
			time.Sunday: {
				"Sunday",
				"Sun",
				"Su",
			},
			time.Monday: {
				"Monday",
				"Mon",
				"Mo",
			},
			time.Tuesday: {
				"Tuesday",
				"Tues",
				"Tue",
				"Tu",
			},
			time.Wednesday: {
				"Wednesday",
				"Weds",
				"Wed",
				"We",
			},
			time.Thursday: {
				"Thursday",
				"Thurs",
				"Thur",
				"Thu",
				"Th",
			},
			time.Friday: {
				"Friday",
				"Fri",
				"Fr",
			},
			time.Saturday: {
				"Saturday",
				"Sat",
				"Sa",
			},
		},
		months: [...][]string{
			time.January: {
				"January",
				"Jan",
				"Ja",
			},
			time.February: {
				"February",
				"Feb",
				"Fe",
			},
			time.March: {
				"March",
				"Mar",
				"Mr",
			},
			time.April: {
				"April",
				"Apr",
				"Ap",
			},
			time.May: {
				"May",
				"My",
			},
			time.June: {
				"June",
				"Jun",
				"Jn",
			},
			time.July: {
				"July",
				"Jul",
				"Jl",
			},
			time.August: {
				"August",
				"Aug",
				"Au",
			},
			time.September: {
				"September",
				"Sept",
				"Sep",
				"Se",
			},
			time.October: {
				"October",
				"Oct",
				"Oc",
			},
			time.November: {
				"November",
				"Nov",
				"No",
			},
			time.December: {
				"December",
				"Dec",
				"De",
			},
		},
		am: []string{
			"am",
			"a",
		},
		pm: []string{
			"pm",
			"p",
		},
		tzPrefix: "GMT",
	},
}

var defaultLanguage = dateLanguages["en"]
