package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	docsFileName := strings.TrimSpace(os.Args[1])
	docsFile, err := os.OpenFile(docsFileName, os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("open docs file %s failed. Error:%s", docsFileName, err)
		os.Exit(1)
	}
	defer docsFile.Close()
	scanner := bufio.NewScanner(docsFile)
	docsFileNameParts := strings.Split(docsFileName, "/")
	resourceName := "alicloud_" + strings.TrimSuffix(docsFileNameParts[len(docsFileNameParts)-1], ".html.markdown")
	exitCode := 0
	fmt.Printf("\n==> Checking docs content of %s ...", docsFileName)
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
				if text != "# "+resourceName {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "# "+resourceName, text)
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
	fmt.Println("\n--- Finished!\n")
	os.Exit(exitCode)
}
