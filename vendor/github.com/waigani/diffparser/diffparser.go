// Copyright (c) 2015 Jesse Meek <https://github.com/waigani>
// This program is Free Software see LICENSE file for details.

package diffparser

import (
	"regexp"
	"strconv"
	"strings"

	"errors"
)

// FileMode represents the file status in a diff
type FileMode int

const (
	// DELETED if the file is deleted
	DELETED FileMode = iota
	// MODIFIED if the file is modified
	MODIFIED
	// NEW if the file is created and there is no diff
	NEW
)

// DiffRange contains the DiffLine's
type DiffRange struct {

	// starting line number
	Start int

	// the number of lines the change diffHunk applies to
	Length int

	// Each line of the hunk range.
	Lines []*DiffLine
}

// DiffLineMode tells the line if added, removed or unchanged
type DiffLineMode rune

const (
	// ADDED if the line is added (shown green in diff)
	ADDED DiffLineMode = iota
	// REMOVED if the line is deleted (shown red in diff)
	REMOVED
	// UNCHANGED if the line is unchanged (not colored in diff)
	UNCHANGED
)

// DiffLine is the least part of an actual diff
type DiffLine struct {
	Mode     DiffLineMode
	Number   int
	Content  string
	Position int // the line in the diff
}

// DiffHunk is a group of difflines
type DiffHunk struct {
	HunkHeader string
	OrigRange  DiffRange
	NewRange   DiffRange
	WholeRange DiffRange
}

// DiffFile is the sum of diffhunks and holds the changes of the file features
type DiffFile struct {
	DiffHeader string
	Mode       FileMode
	OrigName   string
	NewName    string
	Hunks      []*DiffHunk
}

// Diff is the collection of DiffFiles
type Diff struct {
	Files []*DiffFile
	Raw   string `sql:"type:text"`

	PullID uint `sql:"index"`
}

func (d *Diff) addFile(file *DiffFile) {
	d.Files = append(d.Files, file)
}

// Changed returns a map of filename to lines changed in that file. Deleted
// files are ignored.
func (d *Diff) Changed() map[string][]int {
	dFiles := make(map[string][]int)

	for _, f := range d.Files {
		if f.Mode == DELETED {
			continue
		}

		for _, h := range f.Hunks {
			for _, dl := range h.NewRange.Lines {
				if dl.Mode == ADDED { // TODO(waigani) return removed
					dFiles[f.NewName] = append(dFiles[f.NewName], dl.Number)
				}
			}
		}
	}

	return dFiles
}

func regFind(s string, reg string, group int) string {
	re := regexp.MustCompile(reg)
	return re.FindStringSubmatch(s)[group]
}

func lineMode(line string) (*DiffLineMode, error) {
	var m DiffLineMode
	switch line[:1] {
	case " ":
		m = UNCHANGED
	case "+":
		m = ADDED
	case "-":
		m = REMOVED
	default:
		return nil, errors.New("could not parse line mode for line: \"" + line + "\"")
	}
	return &m, nil
}

// Parse takes a diff, such as produced by "git diff", and parses it into a
// Diff struct.
func Parse(diffString string) (*Diff, error) {
	var diff Diff
	diff.Raw = diffString
	lines := strings.Split(diffString, "\n")

	var file *DiffFile
	var hunk *DiffHunk
	var ADDEDCount int
	var REMOVEDCount int
	var inHunk bool
	oldFilePrefix := "--- a/"
	newFilePrefix := "+++ b/"

	var diffPosCount int
	var firstHunkInFile bool
	// Parse each line of diff.
	for idx, l := range lines {
		diffPosCount++
		switch {
		case strings.HasPrefix(l, "diff "):
			inHunk = false

			// Start a new file.
			file = &DiffFile{}
			header := l
			if len(lines) > idx+3 {
				rein := regexp.MustCompile(`^index .+$`)
				remp := regexp.MustCompile(`^(-|\+){3} .+$`)
				index := lines[idx+1]
				if rein.MatchString(index) {
					header = header + "\n" + index
				}
				mp1 := lines[idx+2]
				mp2 := lines[idx+3]
				if remp.MatchString(mp1) && remp.MatchString(mp2) {
					header = header + "\n" + mp1 + "\n" + mp2
				}
			}
			file.DiffHeader = header
			diff.Files = append(diff.Files, file)
			firstHunkInFile = true

			// File mode.
			file.Mode = MODIFIED
		case l == "+++ /dev/null":
			file.Mode = DELETED
		case l == "--- /dev/null":
			file.Mode = NEW
		case strings.HasPrefix(l, oldFilePrefix):
			file.OrigName = strings.TrimPrefix(l, oldFilePrefix)
		case strings.HasPrefix(l, newFilePrefix):
			file.NewName = strings.TrimPrefix(l, newFilePrefix)
		case strings.HasPrefix(l, "@@ "):
			if firstHunkInFile {
				diffPosCount = 0
				firstHunkInFile = false
			}

			inHunk = true
			// Start new hunk.
			hunk = &DiffHunk{}
			file.Hunks = append(file.Hunks, hunk)

			// Parse hunk heading for ranges
			re := regexp.MustCompile(`@@ \-(\d+),?(\d+)? \+(\d+),?(\d+)? @@ ?(.+)?`)
			m := re.FindStringSubmatch(l)
			if len(m) < 5 {
				return nil, errors.New("Error parsing line: " + l)
			}
			a, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			b := a
			if len(m[2]) > 0 {
				b, err = strconv.Atoi(m[2])
				if err != nil {
					return nil, err
				}
			}
			c, err := strconv.Atoi(m[3])
			if err != nil {
				return nil, err
			}
			d := c
			if len(m[4]) > 0 {
				d, err = strconv.Atoi(m[4])
				if err != nil {
					return nil, err
				}
			}
			if len(m[5]) > 0 {
				hunk.HunkHeader = m[5]
			}

			// hunk orig range.
			hunk.OrigRange = DiffRange{
				Start:  a,
				Length: b,
			}

			// hunk new range.
			hunk.NewRange = DiffRange{
				Start:  c,
				Length: d,
			}

			// (re)set line counts
			ADDEDCount = hunk.NewRange.Start
			REMOVEDCount = hunk.OrigRange.Start
		case inHunk && isSourceLine(l):
			m, err := lineMode(l)
			if err != nil {
				return nil, err
			}
			line := DiffLine{
				Mode:     *m,
				Content:  l[1:],
				Position: diffPosCount,
			}
			newLine := line
			origLine := line

			// add lines to ranges
			switch *m {
			case ADDED:
				newLine.Number = ADDEDCount
				hunk.NewRange.Lines = append(hunk.NewRange.Lines, &newLine)
				hunk.WholeRange.Lines = append(hunk.WholeRange.Lines, &newLine)
				ADDEDCount++

			case REMOVED:
				origLine.Number = REMOVEDCount
				hunk.OrigRange.Lines = append(hunk.OrigRange.Lines, &origLine)
				hunk.WholeRange.Lines = append(hunk.WholeRange.Lines, &origLine)
				REMOVEDCount++

			case UNCHANGED:
				newLine.Number = ADDEDCount
				hunk.NewRange.Lines = append(hunk.NewRange.Lines, &newLine)
				hunk.WholeRange.Lines = append(hunk.WholeRange.Lines, &newLine)
				origLine.Number = REMOVEDCount
				hunk.OrigRange.Lines = append(hunk.OrigRange.Lines, &origLine)
				ADDEDCount++
				REMOVEDCount++
			}
		}
	}

	return &diff, nil
}

func isSourceLine(line string) bool {
	if line == `\ No newline at end of file` {
		return false
	}
	if l := len(line); l == 0 || (l >= 3 && (line[:3] == "---" || line[:3] == "+++")) {
		return false
	}
	return true
}

// Length returns the hunks line length
func (hunk *DiffHunk) Length() int {
	return len(hunk.WholeRange.Lines) + 1
}
