package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const layoutGuide = "scripts/ci/provider-layout.md"

// Stable codes and explicit destinations make failures useful to both humans
// and agents. Line 1 identifies a file-level organization rule, not a code defect.
type layoutIssue struct {
	Code   string `json:"code"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	Reason string `json:"reason"`
	Fix    string `json:"fix"`
	Guide  string `json:"guide"`
}

func (i layoutIssue) String() string {
	return fmt.Sprintf("[%s] %s:%d: %s Fix: %s Guide: %s", i.Code, i.File, i.Line, i.Reason, i.Fix, i.Guide)
}

// Read registration identity, never invent it by prepending a directory name.
// Multiple aliases may legitimately resolve to the same implementation.
func factoryFiles(s *snapshot) map[string]bool {
	result := map[string]bool{}
	index := map[string][]string{}
	for _, d := range s.decls {
		if !d.test && d.receiver == "" && d.result.name != "" {
			key := path.Dir(d.file) + ":" + d.name
			index[key] = append(index[key], d.file)
		}
	}
	for file, mappings := range s.registrations {
		for _, factory := range mappings {
			candidates := index[path.Dir(file)+":"+factory]
			if len(candidates) == 1 {
				result[candidates[0]] = true
			}
		}
	}
	return result
}

func entryStem(file string) bool {
	for _, prefix := range []string{"resource_", "data_source_", "ephemeral_", "action_", "list_"} {
		if strings.HasPrefix(path.Base(file), prefix) {
			return true
		}
	}
	return false
}

func helperOwners(s *snapshot, file string, entries map[string]bool, graph map[string][]edge) []string {
	owners, seen := map[string]bool{}, map[string]bool{}
	var queue []string
	for id, d := range s.decls {
		if d.file == file {
			queue = append(queue, id)
		}
	}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		if seen[id] {
			continue
		}
		seen[id] = true
		if d := s.decls[id]; d != nil && entries[d.file] {
			owners[d.file] = true
			continue
		}
		for _, e := range graph[id] {
			if e.ambiguity == "" {
				queue = append(queue, e.caller)
			}
		}
	}
	return sortedKeys(owners)
}

func layoutIssues(before, after *snapshot, changed []string) []layoutIssue {
	issues := []layoutIssue{}
	entries := factoryFiles(after)
	var graph map[string][]edge
	add := func(code, file, reason, fix string) {
		issues = append(issues, layoutIssue{code, file, 1, reason, fix, layoutGuide})
	}
	for _, file := range changed {
		// Grandfather existing layout. Deleting it is allowed; new destinations
		// of renames are checked. Do not demand unrelated repository migration.
		if before.files[file] || !after.files[file] || !strings.HasPrefix(file, "alicloud/service/") ||
			len(strings.Split(file, "/")) != 4 || !strings.HasSuffix(file, ".go") {
			continue
		}
		name := path.Base(file)
		if strings.HasSuffix(name, "_test.go") && entries[strings.TrimSuffix(file, "_test.go")+".go"] {
			continue // a real entity may itself have a local name ending in _unit
		}
		testOwner := ""
		if strings.HasSuffix(name, "_test.go") {
			for entry := range entries {
				if strings.HasPrefix(file, strings.TrimSuffix(entry, ".go")+"_") && len(entry) > len(testOwner) {
					testOwner = entry
				}
			}
		}
		if strings.HasSuffix(name, "_unit_test.go") || testOwner != "" {
			target := strings.TrimSuffix(file, "_unit_test.go") + "_test.go"
			impl := strings.TrimSuffix(file, "_unit_test.go") + ".go"
			if testOwner != "" {
				impl = testOwner
				target = strings.TrimSuffix(testOwner, ".go") + "_test.go"
			}
			if !entries[impl] {
				owner := ""
				for entry := range entries {
					base := strings.TrimSuffix(entry, ".go")
					if strings.HasPrefix(impl, base+"_") && len(entry) > len(owner) {
						owner = entry
					}
				}
				if owner != "" {
					target = strings.TrimSuffix(owner, ".go") + "_test.go"
					impl = owner
				}
			}
			// Keep an internal unit file when the companion is an external test
			// package. Do not tell users to export internals to merge two packages.
			pkg := after.packages[file]
			ownerPkg, companionPkg := after.packages[impl], after.packages[target]
			if ownerPkg != "" && ((pkg == ownerPkg && companionPkg == ownerPkg+"_test") || (pkg == ownerPkg+"_test" && companionPkg == ownerPkg)) {
				continue
			}
			add("LAYOUT002", file, "New split test file duplicates the default companion-test layout.",
				"Move test functions and their mocks into "+target+" (create it if needed); keep TestAcc/TestUnit names so CI runs them separately. Keep a separate file only for complementary internal/external (<package>/<package>_test) packages; never export production symbols just to merge tests.")
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue // mixed tests and existing supplemental files are legal
		}
		if entries[file] {
			for _, prefix := range []string{"resource_", "data_source_", "ephemeral_", "action_", "list_"} {
				if strings.HasPrefix(name, prefix+"alicloud_") {
					target := path.Join(path.Dir(file), prefix+strings.TrimPrefix(name, prefix+"alicloud_"))
					add("LAYOUT003", file, "New service implementation repeats the provider prefix alicloud_ in its filename.",
						"Use a package-local filename such as "+target+" and rename its companion tests consistently; check for destination collisions. Preserve every registered TypeName, Factory and compatibility alias. Service namespace differences must not be guessed from the package directory.")
				}
			}
			continue // a real resource named *_helpers is not an auxiliary file
		}
		stem := strings.TrimSuffix(file, ".go")
		owner := ""
		for entry := range entries {
			base := strings.TrimSuffix(entry, ".go")
			if strings.HasPrefix(stem, base+"_") && len(base) > len(strings.TrimSuffix(owner, ".go")) {
				owner = entry
			}
		}
		auxiliary := false
		for _, suffix := range []string{"_helper", "_helpers", "_decode", "_pagination", "_validation", "_request", "_response"} {
			auxiliary = auxiliary || strings.HasSuffix(stem, suffix)
		}
		if name == "service.go" || strings.HasPrefix(name, "service_") || strings.HasSuffix(name, "_service.go") {
			continue
		}
		// Catch private splits named decode.go/query.go as well, rather than
		// making a helper suffix the only way to recognize the same code.
		if name == "register.go" {
			continue
		}
		if graph == nil {
			graph = reverseGraph(after)
		}
		consumers := helperOwners(after, file, entries, graph)
		if len(consumers) == 1 {
			owner = consumers[0]
		}
		if len(consumers) > 1 {
			owner = ""
			auxiliary = true
		}
		auxiliary = auxiliary || name == "helpers.go" || name == "helper.go"
		if owner != "" || auxiliary {
			destination := owner
			if len(consumers) > 1 {
				add("LAYOUT001", file, "New helper file is shared by multiple registered implementations: "+strings.Join(consumers, ", "),
					"Move shared declarations into service.go or service_<purpose>.go and their shared tests into service_test.go. Keep every consuming resource/data-source test; do not merge shared logic into only one resource.")
				continue
			}
			if destination == "" {
				destination = "the resource_*.go or data_source_*.go that actually calls these helpers (resolve its Factory in register.go)"
			}
			add("LAYOUT001", file, "New split implementation is not a registered entrypoint; resource-specific helpers belong with their owner.",
				"Move the declarations into "+destination+" and move their tests into its companion *_test.go; remove this split file. If multiple entities actually share the logic, use service.go or service_<purpose>.go and retain all consumer tests. Preserve build constraints; incompatible build contexts require maintainer review instead of dropping tags. Renaming a private helper to service.go without checking consumers is not a coverage fix.")
		} else if entryStem(file) {
			add("LAYOUT004", file, "New entry-shaped file cannot be joined to a unique Factory in register.go.",
				"For a new entity, add its literal TypeName/Factory registration in register.go and define the Factory in this file. For private helpers, merge into the owning implementation; for shared service logic, use service.go or service_<purpose>.go. If registration is dynamic, simplify it or request maintainer scope review; do not add a fake registration to silence this error.")
		}
	}
	sort.Slice(issues, func(i, j int) bool { return issues[i].File+issues[i].Code < issues[j].File+issues[j].Code })
	return issues
}

// Local preflight includes tracked edits, deletions and untracked Go files.
// Never compile/load Provider dependencies and never follow a source symlink.
func readWorkingSnapshot() (*snapshot, error) {
	s := &snapshot{decls: map[string]*declaration{}, files: map[string]bool{}, assets: map[string]string{}}
	listed, err := command("ls-files", "-z", "--cached", "--others", "--exclude-standard", "--", "alicloud/")
	if err != nil {
		return nil, err
	}
	for _, file := range strings.Split(string(listed), "\x00") {
		if file == "" || s.files[file] {
			continue
		}
		info, err := os.Lstat(filepath.FromSlash(file))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("%s: expected a regular source file", file)
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		s.files[file] = true
		if strings.HasSuffix(file, ".go") {
			if err = parseSource(s, file, data); err != nil {
				return nil, err
			}
		} else {
			s.assets[file] = fmt.Sprintf("%x", sha256.Sum256(data))
		}
	}
	return s, nil
}
