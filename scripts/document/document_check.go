package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	os.Exit(checkDocument(strings.TrimSpace(os.Args[1])))
}

func checkDocument(docsFileName string) int {
	docsFile, err := os.OpenFile(docsFileName, os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("open docs file %s failed. Error:%s", docsFileName, err)
		return 1
	}
	defer docsFile.Close()
	scanner := bufio.NewScanner(docsFile)
	docsFileNameParts := strings.Split(docsFileName, "/")
	resourceName := "alicloud_" + strings.TrimSuffix(docsFileNameParts[len(docsFileNameParts)-1], ".html.markdown")
	expectedTitle := "# " + resourceName
	if len(docsFileNameParts) > 1 &&
		(docsFileNameParts[len(docsFileNameParts)-2] == "e" || docsFileNameParts[len(docsFileNameParts)-2] == "ephemeral-resources") {
		expectedTitle = "# Ephemeral: " + resourceName
	}
	exitCode := 0
	fmt.Printf("\n==> Checking docs content of %s ...", docsFileName)
	// Build a docs index once so every line can be validated against the
	// actual resource/data-source/function/guide/action/ephemeral inventory.
	// When the file is not under website/docs (e.g. unit-test fixtures) the
	// root is empty and the cross-reference check is skipped.
	linkIndex := buildDocIndex(repoRootFromDocPath(docsFileName))
	titleCheck := false
	titleDescriptionCheck := false
	descriptionCheck := false
	versionChecked := false
	exampleCheck := false
	exampleBlockOpen := false
	argumentCheck := false
	attributesCheck := false
	importCheck := false
	importBlockOpen := false
	timeoutCheck := false
	line := 0
	for scanner.Scan() {
		line += 1
		text := scanner.Text()
		// Cross-reference links may appear in any prose section, so they
		// are validated on every line regardless of the structure state.
		for _, msg := range checkCrossReferenceLinks(text, line, linkIndex) {
			fmt.Println(msg)
			exitCode = 1
		}
		if line == 1 {
			titleCheck = true
			if text != "---" {
				fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "---", text)
				exitCode = 1
			}
			continue
		}
		if titleCheck {
			if strings.HasPrefix(text, "layout:") && text != "layout: \"alicloud\"" {
				fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "layout: \"alicloud\"", text)
				exitCode = 1
				continue
			}
			if strings.HasPrefix(text, "page_title:") && text != "page_title: \"Alicloud: "+resourceName+"\"" {
				fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "page_title: \"Alicloud: "+resourceName, text)
				exitCode = 1
				continue
			}
			if strings.HasPrefix(text, "description:") {
				titleDescriptionCheck = true
				if text != "description: |-" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "description: |- ", text)
					exitCode = 1
				}
				continue
			}
			if titleDescriptionCheck {
				titleDescriptionCheck = false
				if !strings.HasPrefix(text, "  ") {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "  "+text, text)
					exitCode = 1
				}
				continue
			}
			if text == "---" {
				titleCheck = false
				descriptionCheck = true
				continue
			}
		}
		if descriptionCheck {
			if strings.HasPrefix(text, "# alicloud") || strings.HasSuffix(text, resourceName) {
				if titleCheck {
					fmt.Printf("\n[Error] line %d: docs title has not been closed.", line)
					exitCode = 1
				}
				if strings.Contains(text, "\\") {
					fmt.Printf("\n[WARNING] line %d: please remove the \\", line)
					text = strings.Replace(text, "\\", "", -1)
				}
				if text != expectedTitle {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, expectedTitle, text)
					exitCode = 1
				}
				continue
			}
			if strings.Contains(text, "https://help.aliyun.com/") {
				fmt.Printf("\n[Error] line %d: Expected: an international site link: %s. Got: %s.", line, "https://www.alibabacloud.com/help", "https://help.aliyun.com/")
				exitCode = 1
				continue
			}
			if strings.Contains(text, "Available since v") {
				versionChecked = true
				v := "Available since v"
				if strings.Contains(text, v) && !strings.HasPrefix(text, "-> **NOTE:** "+v) {
					parts := strings.Split(text, v)
					fmt.Printf("\n[Error] line %d: Expected: %s Got: %s", line, "-> **NOTE:** "+v+strings.Replace(parts[1], "+", "", -1), text)
					exitCode = 1
				}
				continue
			}
			if strings.Contains(text, "Available in v") {
				versionChecked = true
				v := "Available in v"
				if strings.Contains(text, v) {
					parts := strings.Split(text, v)
					v = strings.Replace(strings.Replace(v, "in", "since", -1), "from", "since", -1)
					fmt.Printf("\n[Error] line %d: Expected: %s Got: %s", line, parts[0]+v+strings.Replace(parts[1], "+", "", -1), text)
					exitCode = 1
				}
				continue
			}
			if strings.Contains(text, "-> **DEPRECATED:**") || strings.Contains(text, "-> **NOTE:** deprecated since") ||
				strings.Contains(text, "-> **NOTE:** Deprecated since") ||
				(strings.Contains(text, "This resource has been") && (strings.Contains(text, "deprecated") || strings.Contains(text, "DEPRECATED"))) {
				versionChecked = true
				if !strings.HasPrefix(text, "-> **DEPRECATED:**") {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "-> **DEPRECATED:** This [resource|data source] has been deprecated from vxxx", text)
					exitCode = 1
				}
				continue
			}

			if strings.HasPrefix(text, "## Example Usage") || strings.HasSuffix(text, "Example Usage") {
				descriptionCheck = false
				exampleCheck = true
				if !versionChecked {
					fmt.Printf("\n[Error] line %d: missing available or deprecated verison info.", line)
					exitCode = 1
				}
				if text != "## Example Usage" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Example Usage", text)
					exitCode = 1
				}
				continue
			}
		}
		if exampleCheck {
			if strings.HasPrefix(text, "```") {
				if !exampleBlockOpen {
					exampleBlockOpen = true
				} else {
					exampleBlockOpen = false
					if text != "```" {
						fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "```", text)
						exitCode = 1
					}
				}
				if exampleBlockOpen && text != "```terraform" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "```terraform", text)
					exitCode = 1
				}
				continue
			}
			if exampleBlockOpen && strings.Contains(text, "test") {
				text = strings.TrimSpace(text)
				parts := strings.Split(text, "=")
				if key := strings.TrimSpace(parts[0]); key == "default" || strings.HasSuffix(key, "name") {
					if value := strings.TrimSpace(parts[1]); strings.Contains(value, "test") {
						fmt.Printf("\n[Error] line %d: avoid using 'test' as soon as possable in the example", line)
						exitCode = 1
					}
				}
			}
			if exampleBlockOpen && strings.Contains(text, "depends_on") {
				text = strings.TrimSpace(text)
				parts := strings.Split(text, "=")
				if key := strings.TrimSpace(parts[0]); key == "depends_on" {
					fmt.Printf("\n[Error] line %d: avoid using 'depends_on' as soon as possable in the example", line)
					exitCode = 1
				}
			}
			if strings.HasPrefix(text, "## Argument Reference") || strings.HasSuffix(text, "Argument Reference") {
				exampleCheck = false
				argumentCheck = true
				if exampleBlockOpen {
					fmt.Printf("\n[Error] line %d: example block has not been closed.", line)
					exitCode = 1
				}
				if text != "## Argument Reference" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Argument Reference", text)
					exitCode = 1
				}
				continue
			}
			if !argumentCheck && strings.HasPrefix(text, "## Attributes Reference") || strings.HasSuffix(text, "Attributes Reference") {
				exampleCheck = false
				attributesCheck = true
				if text != "## Attributes Reference" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Attributes Reference", text)
					exitCode = 1
				}
				continue
			}
		}
		if argumentCheck {
			if strings.HasPrefix(text, "* `") {
				parts := strings.Split(strings.Split(text, ")")[0], "(")
				for i, tag := range strings.Split(parts[len(parts)-1], ",") {
					if strings.HasPrefix(tag, " ") {
						if i == 0 {
							fmt.Printf("\n[Error] line %d: please remove redundant space prefix for %s. ", line, tag)
							exitCode = 1
						}
					} else if i > 0 {
						fmt.Printf("\n[Error] line %d: missing space prefix for %s. ", line, tag)
						exitCode = 1
					}
					if strings.HasSuffix(tag, " ") {
						fmt.Printf("\n[Error] line %d: please remove redundant space suffix for %s.", line, tag)
						exitCode = 1
					}
				}
				continue
			}
			if strings.HasPrefix(text, "###") {
				parts := strings.Split(text, " ")
				block := strings.ToLower(parts[len(parts)-1])
				if len(parts) > 2 ||
					(strings.HasPrefix(block, "`") && !strings.HasSuffix(block, "`")) ||
					(!strings.HasPrefix(block, "`") && strings.HasSuffix(block, "`")) {
					block = strings.Trim(block, "`")
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "### `"+block+"`", text)
					exitCode = 1
				}
				block = strings.Trim(block, "`")
				blockParts := strings.Split(block, "-")
				blockLink := fmt.Sprintf("[`%s`](#%s)", blockParts[len(blockParts)-1], block)
				docsContent, err := os.ReadFile(docsFileName)
				if err != nil {
					fmt.Printf("\n[Error] reading docs file %s failed. Error: %s", docsFileName, err)
					exitCode = 1
				} else if !strings.Contains(string(docsContent), blockLink) {
					fmt.Printf("\n[Error] line %d: missing link for block `%s`. Expected link like: See %s below.", line, block, blockLink)
					exitCode = 1
				}
				continue
			}
			if strings.HasPrefix(text, "## Argument Reference") || strings.HasSuffix(text, "Argument Reference") {
				argumentCheck = false
				attributesCheck = true
				if text != "## Attributes Reference" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Attributes Reference", text)
					exitCode = 1
				}
				continue
			}
			if strings.HasPrefix(text, "## Attributes Reference") || strings.HasSuffix(text, "Attributes Reference") {
				argumentCheck = false
				attributesCheck = true
				if text != "## Attributes Reference" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Attributes Reference", text)
					exitCode = 1
				}
				continue
			}
		}

		if attributesCheck {
			if strings.HasPrefix(text, "## Timeout") || strings.HasSuffix(text, "Timeouts") || strings.HasSuffix(text, "Timeout") {
				attributesCheck = false
				timeoutCheck = true
				if text != "## Timeouts" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Timeouts", text)
					exitCode = 1
				}
				continue
			}
			if strings.HasPrefix(text, "## Import") || strings.HasSuffix(text, "Import") {
				attributesCheck = false
				importCheck = true
				if text != "## Import" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Import", text)
					exitCode = 1
				}
				continue
			}
		}

		if timeoutCheck {
			if strings.HasPrefix(text, "## Import") || strings.HasSuffix(text, "Import") {
				timeoutCheck = false
				importCheck = true
				if text != "## Import" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Import", text)
					exitCode = 1
				}
				continue
			}
		}

		if importCheck {
			if !importBlockOpen && strings.Contains(text, "terraform import") {
				fmt.Printf("\n[Error] line %d: Expected: %s.", line-1, "```shell")
				exitCode = 1
				break
			}

			if strings.HasPrefix(text, "```") {
				if !importBlockOpen {
					importBlockOpen = true
				} else {
					importBlockOpen = false
					if text != "```" {
						fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "```", text)
						exitCode = 1
					}
				}
				if importBlockOpen && text != "```shell" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "```shell", text)
					exitCode = 1
				}
				continue
			}
		}
	}
	if importBlockOpen {
		fmt.Println("\n[Error] Error: import block has not been closed.")
		exitCode = 1
	}
	fmt.Println("\n--- Finished!")
	return exitCode
}

// sectionToDir maps Terraform Registry docs path segments to the local
// website/docs subdirectory that holds the matching document files.
var sectionToDir = map[string]string{
	"resources":           "r",
	"data-sources":        "d",
	"functions":           "functions",
	"guides":              "guides",
	"actions":             "actions",
	"ephemeral-resources": "ephemeral-resources",
}

// docIndex records the document basenames available under each website/docs
// subdirectory. It is used to validate that cross-reference links in a
// document point at resources, data sources, functions, guides, actions or
// ephemeral resources that actually exist in this provider.
type docIndex struct {
	byDir map[string]map[string]bool
}

// registryLinkRe captures the docs path of an alicloud Terraform Registry
// link, e.g. https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/foo
// The path-safe character class avoids trailing punctuation in bare URLs and
// stops at markdown delimiters such as ')' in [label](url).
var registryLinkRe = regexp.MustCompile(
	`https?://registry\.terraform\.io/providers/aliyun/alicloud/[^/]+/docs/([a-zA-Z0-9/._#-]+)`)

// oldLinkRe captures the path of a deprecated www.terraform.io alicloud link,
// e.g. https://www.terraform.io/docs/providers/alicloud/r/foo
var oldLinkRe = regexp.MustCompile(
	`https?://www\.terraform\.io/docs/providers/alicloud/([a-zA-Z0-9/._#-]+)`)

// slugRe extracts the leading document slug (resource/data-source name) from a
// path fragment, stopping at any character that cannot be part of a name.
var slugRe = regexp.MustCompile(`^[a-z0-9_-]+`)

// cleanName strips an anchor fragment and trailing non-slug characters from a
// path segment so a residual file suffix or sentence punctuation does not
// masquerade as part of the target name.
func cleanName(raw string) string {
	if i := strings.Index(raw, "#"); i >= 0 {
		raw = raw[:i]
	}
	return slugRe.FindString(raw)
}

// hasResidualSuffix reports whether a path segment still carries a document
// file extension that should not appear in a Registry link.
func hasResidualSuffix(seg string) bool {
	seg = strings.ToLower(seg)
	return strings.HasSuffix(seg, ".html.markdown") ||
		strings.HasSuffix(seg, ".html.md") ||
		strings.HasSuffix(seg, ".html") ||
		strings.HasSuffix(seg, ".md")
}

// repoRootFromDocPath derives the repository root from a documentation file
// path by locating the "website/docs/" segment. It returns "" for paths that
// are not under website/docs (for example test fixtures in a temp directory),
// in which case the link check is skipped.
func repoRootFromDocPath(docPath string) string {
	p := filepath.ToSlash(docPath)
	idx := strings.Index(p, "website/docs/")
	if idx < 0 {
		return ""
	}
	return p[:idx]
}

// buildDocIndex globs each website/docs subdirectory and records the available
// document basenames. A nil index is returned when repoRoot is empty.
func buildDocIndex(repoRoot string) *docIndex {
	if repoRoot == "" {
		return nil
	}
	idx := &docIndex{byDir: map[string]map[string]bool{}}
	seen := map[string]bool{}
	for _, dir := range sectionToDir {
		if seen[dir] {
			continue
		}
		seen[dir] = true
		set := map[string]bool{}
		matches, _ := filepath.Glob(filepath.Join(repoRoot, "website", "docs", dir, "*"))
		for _, m := range matches {
			base := filepath.Base(m)
			for _, ext := range []string{".html.markdown", ".html.md", ".md"} {
				if strings.HasSuffix(base, ext) {
					set[strings.TrimSuffix(base, ext)] = true
					break
				}
			}
		}
		idx.byDir[dir] = set
	}
	return idx
}

// resolveTarget reports whether a cross-reference target exists in the section
// directory. For resources and data sources it also suggests the alternate
// section when the name is filed under the wrong one.
func (idx *docIndex) resolveTarget(section, name string) (exists bool, suggestion string) {
	dir, ok := sectionToDir[section]
	if !ok {
		return true, "" // unknown section: provider-level page, do not flag.
	}
	if idx.byDir[dir][name] {
		return true, ""
	}
	var other string
	switch section {
	case "resources":
		other = "data-sources"
	case "data-sources":
		other = "resources"
	}
	if other != "" && idx.byDir[sectionToDir[other]][name] {
		return false, other + "/" + name
	}
	return false, ""
}

func (idx *docIndex) validateRegistryPath(path string, line int) []string {
	rest := path
	if i := strings.Index(rest, "#"); i >= 0 {
		rest = rest[:i]
	}
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return nil // bare docs root, provider index.
	}
	section := parts[0]
	if _, ok := sectionToDir[section]; !ok {
		return nil // provider-level or unknown section page; do not flag.
	}
	if len(parts) < 2 || parts[1] == "" {
		return nil // section index page, e.g. docs/resources.
	}
	nameSeg := parts[1]
	name := cleanName(nameSeg)
	if name == "" {
		return nil
	}
	var msgs []string
	if hasResidualSuffix(nameSeg) {
		msgs = append(msgs, fmt.Sprintf(
			"\n[Error] line %d: cross-reference link has a residual file suffix in path: %s",
			line, path))
	}
	exists, suggestion := idx.resolveTarget(section, name)
	if exists {
		return msgs
	}
	if suggestion != "" {
		msgs = append(msgs, fmt.Sprintf(
			"\n[Error] line %d: cross-reference section mismatch: %s is filed under %s, not %s; use %s instead. Link: .../docs/%s",
			line, name, suggestion, section, suggestion, path))
	} else {
		msgs = append(msgs, fmt.Sprintf(
			"\n[Error] line %d: cross-reference target does not exist: %s/%s. No matching document at website/docs/%s/%s.html.markdown. Link: .../docs/%s",
			line, section, name, sectionToDir[section], name, path))
	}
	return msgs
}

func (idx *docIndex) validateOldFormat(path string, line int) []string {
	msgs := []string{fmt.Sprintf(
		"\n[Error] line %d: Terraform link uses deprecated www.terraform.io format; use https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/<section>/<name> instead. Got: https://www.terraform.io/docs/providers/alicloud/%s",
		line, path)}
	rest := path
	if i := strings.Index(rest, "#"); i >= 0 {
		rest = rest[:i]
	}
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) < 2 || parts[1] == "" {
		return msgs // provider index or section root; the format error is enough.
	}
	var section string
	switch parts[0] {
	case "r":
		section = "resources"
	case "d":
		section = "data-sources"
	default:
		return msgs // unknown old-format section; the format error is enough.
	}
	name := cleanName(parts[1])
	if name == "" {
		return msgs
	}
	exists, suggestion := idx.resolveTarget(section, name)
	if exists {
		return msgs
	}
	if suggestion != "" {
		msgs = append(msgs, fmt.Sprintf(
			"\n[Error] line %d: cross-reference section mismatch: %s is filed under %s, not %s; use https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/%s instead. Link: .../alicloud/%s",
			line, name, suggestion, section, suggestion, path))
	} else {
		msgs = append(msgs, fmt.Sprintf(
			"\n[Error] line %d: cross-reference target does not exist: %s/%s. No matching document at website/docs/%s/%s.html.markdown. Link: .../alicloud/%s",
			line, section, name, sectionToDir[section], name, path))
	}
	return msgs
}

// checkCrossReferenceLinks scans a documentation line for Terraform
// cross-reference links and returns descriptive error messages for any that
// are broken, mismatched, suffixed or in the deprecated www.terraform.io form.
// Third-party provider links (e.g. hashicorp/vault) are ignored. It returns
// nil when the index is unavailable, so test fixtures outside website/docs are
// unaffected.
func checkCrossReferenceLinks(text string, line int, idx *docIndex) []string {
	if idx == nil {
		return nil
	}
	var msgs []string
	for _, m := range oldLinkRe.FindAllStringSubmatch(text, -1) {
		msgs = append(msgs, idx.validateOldFormat(m[1], line)...)
	}
	for _, m := range registryLinkRe.FindAllStringSubmatch(text, -1) {
		msgs = append(msgs, idx.validateRegistryPath(m[1], line)...)
	}
	return msgs
}
