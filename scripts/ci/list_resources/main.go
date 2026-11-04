// list_resources resolves changed Terraform Provider source files to the
// registered SDKv2 resource and data-source type names they implement, so
// local-ci-check.sh can list and integration-test split-layout files the way
// it lists root-layout files. A split-layout path does not yield the resource
// name on its own (alicloud_cen_* lives in service/cbn), so the product's
// register.go is the source of truth, consulted here via splitlayout.Resolve.
//
// Input: one file path per line on stdin. Output: one line per resolved
// entity, "<type>\t<doc>" where <type> is the full Terraform type name and
// <doc> is the website/docs/{r,d}/<name>.html.markdown path. Test files,
// framework-native registrations and non-resource files produce no line:
// integration tests and doc association target implementation files, and the
// SDKv2 diff checkers' framework out-of-scope boundary is preserved.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if path == "" {
			continue
		}
		for _, e := range resolveEntity(path) {
			fmt.Printf("%s\t%s\n", e.typeName, e.docPath)
		}
	}
}

type entity struct {
	typeName string
	docPath  string
}

// resolveEntity maps one source path to the SDKv2 entities it implements. It
// routes split-layout files through register.go (splitlayout.Resolve) so the
// product slug in the path never masquerades as the resource name, and falls
// back to the root-layout filename convention for alicloud/resource_alicloud_*
// and alicloud/data_source_alicloud_* files.
func resolveEntity(path string) []entity {
	if strings.HasSuffix(path, "_test.go") {
		return nil
	}
	if product, stem, isTest, ok := splitlayout.IsSplitSource(path); ok {
		if isTest {
			return nil
		}
		resolved, err := splitlayout.Resolve(product, stem)
		if err != nil {
			return nil
		}
		var out []entity
		for _, r := range resolved {
			if r.Kind != splitlayout.KindResource && r.Kind != splitlayout.KindDataSource {
				continue
			}
			out = append(out, entity{r.TypeName, docPath(r.TypeName, r.Kind)})
		}
		if len(out) == 0 {
			// A split-layout file that registers no SDKv2 entity may still be
			// framework-native (managed resource, data source, ephemeral, list
			// or action): those are outside the SDKv2 integration-test scope
			// and are logged by the runtime coverage checker. Surface the
			// boundary here so the skip reason is explicit, not silent.
			if fw, err := splitlayout.FrameworkResolve(product, stem); err == nil && len(fw) > 0 {
				fmt.Fprintf(os.Stderr, "[list_resources] %s: framework-native, outside SDKv2 integration-test scope\n", path)
			}
		}
		return out
	}
	if tfType, kind, ok := resolveRootEntity(path); ok {
		return []entity{{tfType, docPath(tfType, kind)}}
	}
	return nil
}

// resolveRootEntity derives the Terraform type name from a root-layout
// implementation file. Root files follow the resource_alicloud_<name> /
// data_source_alicloud_<name> convention, so the type name is the stem with
// the resource_/data_source_ prefix stripped. It is pure string logic (no
// register.go) so it is unit-testable without the repository root as the
// working directory that splitlayout.Resolve requires.
func resolveRootEntity(path string) (tfType string, kind splitlayout.Kind, ok bool) {
	if strings.HasSuffix(path, "_test.go") {
		return "", "", false
	}
	base := filepath.Base(path)
	if !strings.HasSuffix(base, ".go") {
		return "", "", false
	}
	stem := strings.TrimSuffix(base, ".go")
	switch {
	case strings.HasPrefix(stem, "resource_alicloud_"):
		return strings.TrimPrefix(stem, "resource_"), splitlayout.KindResource, true
	case strings.HasPrefix(stem, "data_source_alicloud_"):
		return strings.TrimPrefix(stem, "data_source_"), splitlayout.KindDataSource, true
	}
	return "", "", false
}

// docPath returns the website/docs path for a Terraform type name. Resources
// document under r/, data sources under d/; the name is the type without the
// alicloud_ prefix, matching the root-layout doc convention.
func docPath(tfType string, kind splitlayout.Kind) string {
	name := strings.TrimPrefix(tfType, "alicloud_")
	dir := "r"
	if kind == splitlayout.KindDataSource {
		dir = "d"
	}
	return fmt.Sprintf("website/docs/%s/%s.html.markdown", dir, name)
}
