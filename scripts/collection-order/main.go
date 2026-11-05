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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func main() {
	base := flag.String("base", "", "full merge-base SHA from resolve-change-range.sh")
	head := flag.String("head", "", "full checked-out head SHA")
	resourceName := flag.String("resource", "", "check one registered resource locally")
	dataSourceName := flag.String("data-source", "", "check one registered data source locally")
	flag.Parse()
	if err := run(*base, *head, *resourceName, *dataSourceName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(base, head, resourceName, dataSourceName string) error {
	current, err := os.ReadFile("alicloud/provider.go")
	if err != nil {
		return err
	}
	registers, err := currentRegisters()
	if err != nil {
		return err
	}
	registry, err := coverage.RegistryUnion(current, registers)
	if err != nil {
		return err
	}
	dsRegistry, err := coverage.DataSourcesRegistryUnion(current, registers)
	if err != nil {
		return err
	}
	sources := map[string][]byte{}
	for _, pattern := range []string{"alicloud/resource_alicloud_*.go", "alicloud/service/*/resource_*.go", "alicloud/provider/core/resource_*.go"} {
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
	dsSources := map[string][]byte{}
	for _, pattern := range []string{"alicloud/data_source_alicloud_*.go", "alicloud/service/*/data_source_*.go", "alicloud/provider/core/data_source_*.go"} {
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
			dsSources[path] = src
		}
	}
	var changed []string
	before := registry
	dsBefore := dsRegistry
	if resourceName != "" || dataSourceName != "" {
		if base != "" || head != "" {
			return fmt.Errorf("use -resource/-data-source or -base/-head, not both")
		}
		if resourceName != "" && dataSourceName != "" {
			return fmt.Errorf("use -resource or -data-source, not both")
		}
		changed = []string{"alicloud/provider.go"}
		if resourceName != "" {
			if registry[resourceName] == "" {
				return fmt.Errorf("resource is not registered: %s", resourceName)
			}
			before = withoutTarget(registry, resourceName)
			// dsBefore stays dsRegistry: no data source is selected in
			// resource single-entity mode.
		} else {
			if dsRegistry[dataSourceName] == "" {
				return fmt.Errorf("data source is not registered: %s", dataSourceName)
			}
			dsBefore = withoutTarget(dsRegistry, dataSourceName)
			// before stays registry: no resource is selected in data-source
			// single-entity mode.
		}
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
		before, err = coverage.RegistryUnion(src, baseRegisters)
		if err != nil {
			return err
		}
		dsBefore, err = coverage.DataSourcesRegistryUnion(src, baseRegisters)
		if err != nil {
			return err
		}
	}
	oldSource := func(path string) ([]byte, error) { return git("show", base+":"+path) }
	selected, err := coverage.Select(changed, before, registry, sources, oldSource)
	if err != nil {
		return err
	}
	dsSelected, err := coverage.SelectScope(changed, dsBefore, dsRegistry, dsSources, oldSource, coverage.DataSourceScope)
	if err != nil {
		return err
	}
	if resourceName != "" {
		selected = map[string]string{resourceName: selected[resourceName]}
		dsSelected = map[string]string{}
	}
	if dataSourceName != "" {
		dsSelected = map[string]string{dataSourceName: dsSelected[dataSourceName]}
		selected = map[string]string{}
	}
	if len(selected) == 0 && len(dsSelected) == 0 {
		fmt.Println("TypeList Order Coverage: no directly changed resource or data source implementations, tests or registrations.")
		return nil
	}
	provider := alicloud.Provider()
	failed := false
	failed = checkSelected(selected, registry, provider.ResourcesMap, coverage.ResourceScope) || failed
	failed = checkSelected(dsSelected, dsRegistry, provider.DataSourcesMap, coverage.DataSourceScope) || failed
	if failed {
		return fmt.Errorf("TypeList Order Coverage failed; see scripts/collection-order/README.md")
	}
	return nil
}

// checkSelected runs the order-coverage Check over a scope's selected
// implementations and reports per-entity results. The schema map and scope
// pair so resources run against ResourcesMap with ResourceScope and data
// sources against DataSourcesMap with DataSourceScope; the Check logic and
// diagnostics are identical apart from the read-only waiver a data source's
// order-insensitive Read lifecycle earns.
func checkSelected(selected, registry map[string]string, schemas map[string]*schema.Resource, scope coverage.Scope) bool {
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
		r := schemas[name]
		if r == nil {
			fmt.Fprintf(os.Stderr, "%s missing from runtime schema: %s\n", scope.Label(), name)
			failed = true
			continue
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
		// The suppressor unit net for a split-package entity may live in a
		// sibling *_unit_test.go file in the same package directory rather
		// than the main *_test.go file that carries the TestAcc steps. Gather
		// every Go test file in the same directory so resolveSuppressNet can
		// span them; the main file itself is parsed by Check and is excluded.
		testSources := map[string][]byte{}
		testFiles, terr := filepath.Glob(filepath.Join(filepath.Dir(path), "*_test.go"))
		if terr == nil {
			for _, tf := range testFiles {
				if tf == path {
					continue
				}
				ts, terr := os.ReadFile(tf)
				if terr != nil {
					continue
				}
				testSources[tf] = ts
			}
		}
		issues := coverage.CheckScope(path, src, testSources, name, r, scope, aliases...)
		for _, issue := range issues {
			fmt.Fprintf(os.Stderr, "%s (%s): %s\n", name, path, issue)
		}
		if len(issues) > 0 {
			failed = true
		} else {
			fmt.Printf("%s: TypeList Order Coverage passed\n", name)
		}
	}
	return failed
}

// withoutTarget returns a copy of registry without the target entry, so a
// single-entity check selects only that entity the way a diff would.
func withoutTarget(registry map[string]string, target string) map[string]string {
	out := make(map[string]string, len(registry))
	for name, constructor := range registry {
		if name != target {
			out[name] = constructor
		}
	}
	return out
}

func git(args ...string) ([]byte, error) {
	output, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, output)
	}
	return output, nil
}

// currentRegisters reads supported product and core registrations.
func currentRegisters() (map[string][]byte, error) {
	files, err := splitlayout.RegisterPaths()
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
