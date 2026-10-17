package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/scripts/collection-order/internal/coverage"
	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

func main() {
	base := flag.String("base", "", "full merge-base SHA from resolve-change-range.sh")
	head := flag.String("head", "", "full checked-out head SHA")
	resourceName := flag.String("resource", "", "check one registered resource locally")
	flag.Parse()
	if err := run(*base, *head, *resourceName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(base, head, resourceName string) error {
	current, err := os.ReadFile("alicloud/provider.go")
	if err != nil {
		return err
	}
	registers, err := currentRegisters()
	if err != nil {
		return err
	}
	registry, err := coverage.CombinedRegistry(current, registers)
	if err != nil {
		return err
	}
	sources := map[string][]byte{}
	for _, pattern := range []string{"alicloud/resource_alicloud_*.go", "alicloud/service/*/resource_*.go"} {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		for _, path := range files {
			if strings.HasSuffix(path, "_test.go") {
				continue
			}
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			sources[path] = src
		}
	}
	var changed []string
	before := registry
	if resourceName != "" {
		if base != "" || head != "" {
			return fmt.Errorf("use -resource or -base/-head, not both")
		}
		if registry[resourceName] == "" {
			return fmt.Errorf("resource is not registered: %s", resourceName)
		}
		before = map[string]string{}
		for name, constructor := range registry {
			if name != resourceName {
				before[name] = constructor
			}
		}
		changed = []string{"alicloud/provider.go"}
	} else {
		valid := regexp.MustCompile(`^[0-9a-f]{40}$`)
		if !valid.MatchString(base) || !valid.MatchString(head) {
			return fmt.Errorf("-base and -head must be full commit SHAs")
		}
		actual, err := git("rev-parse", "HEAD")
		if err != nil {
			return err
		}
		if strings.TrimSpace(string(actual)) != head {
			return fmt.Errorf("-head must equal the checked-out HEAD")
		}
		output, err := git("diff", "--name-only", "--no-renames", "-z", base, head, "--", "alicloud/")
		if err != nil {
			return err
		}
		changed = strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00")
		src, err := git("show", base+":alicloud/provider.go")
		if err != nil {
			return err
		}
		// Products that did not exist at the base revision simply contribute
		// no registrations there: their registrations are all new.
		baseRegisters := map[string][]byte{}
		for _, path := range sortedRegisterPaths(registers) {
			registerSrc, err := git("show", base+":"+path)
			if err != nil {
				continue
			}
			baseRegisters[path] = registerSrc
		}
		before, err = coverage.CombinedRegistry(src, baseRegisters)
		if err != nil {
			return err
		}
	}
	selected, err := coverage.Select(changed, before, registry, sources, func(path string) ([]byte, error) { return git("show", base+":"+path) })
	if err != nil {
		return err
	}
	if resourceName != "" {
		selected = map[string]string{resourceName: selected[resourceName]}
	}
	if len(selected) == 0 {
		fmt.Println("TypeList Order Coverage: no directly changed resource implementations, tests or registrations.")
		return nil
	}
	provider := alicloud.Provider()
	failed := false
	checked := map[string]bool{}
	for _, name := range coverage.Names(selected) {
		constructor := registry[name]
		if checked[constructor] {
			continue
		}
		checked[constructor] = true
		var aliases []string
		for alias, fn := range registry {
			if fn == constructor {
				aliases = append(aliases, alias)
			}
		}
		r := provider.ResourcesMap[name]
		if r == nil {
			// Framework-native resources are not in the SDKv2 runtime schema.
			// Instead of skipping them as out of scope, synthesize their schema
			// from the real Schema() AST so the collection-order check applies
			// with the same fail-closed semantics as SDKv2: a framework resource
			// with a configurable ordered collection (ListNestedBlock) is checked,
			// and one without reaches an empty candidate set by reading its real
			// schema, not by a generic skip. A parse failure is reported, not
			// swallowed.
			resolved := splitlayout.FrameworkFindByTypeName(name)
			if resolved == nil {
				return fmt.Errorf("resource missing from runtime schema: %s", name)
			}
			synth, err := coverage.FrameworkResource(resolved.ImplPath)
			if err != nil {
				return fmt.Errorf("framework resource %s: %w", name, err)
			}
			r = synth
		}
		if len(coverage.Candidates(r)) == 0 {
			fmt.Printf("%s: no configurable multi-member TypeList attributes\n", name)
			continue
		}
		path := selected[name]
		src, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
			failed = true
			continue
		}
		issues := coverage.Check(path, src, name, r, aliases...)
		for _, issue := range issues {
			fmt.Fprintf(os.Stderr, "%s (%s): %s\n", name, path, issue)
		}
		if len(issues) > 0 {
			failed = true
		} else {
			fmt.Printf("%s: TypeList Order Coverage passed\n", name)
		}
	}
	if failed {
		return fmt.Errorf("TypeList Order Coverage failed; see scripts/collection-order/README.md")
	}
	return nil
}

func git(args ...string) ([]byte, error) {
	output, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, output)
	}
	return output, nil
}

// currentRegisters reads every product register.go in the working tree.
func currentRegisters() (map[string][]byte, error) {
	files, err := filepath.Glob("alicloud/service/*/register.go")
	if err != nil {
		return nil, err
	}
	registers := map[string][]byte{}
	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		registers[path] = src
	}
	return registers, nil
}

// sortedRegisterPaths gives the base-revision reads a deterministic order.
func sortedRegisterPaths(registers map[string][]byte) []string {
	paths := make([]string, 0, len(registers))
	for path := range registers {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}
