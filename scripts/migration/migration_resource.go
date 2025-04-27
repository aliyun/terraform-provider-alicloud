package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	namespace       = flag.String("n", "", "namespace")
	resource        = flag.String("r", "", "resource")
	destProviderDir = flag.String("t", "", "target dir")
)

func main() {

	flag.Parse()

	if *namespace == "" || *resource == "" || *destProviderDir == "" {
		log.Fatal("All parameters ( -n, -r, -t) are required.")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	sourceDir := fmt.Sprintf("%s/alicloud", filepath.Dir(filepath.Dir(currentDir)))
	sourceFileName := fmt.Sprintf("resource_alicloud_%s_%s.go", *namespace, *resource)
	if sourceFileName == "resource_alicloud_vpc_vswitch.go" {
		sourceFileName = "resource_alicloud_vswitch.go"
	}

	sourceFile := fmt.Sprintf("%s/%s", sourceDir, sourceFileName)

	destFileName := fmt.Sprintf("%s.go", *resource)
	destFile := filepath.Join(*destProviderDir, "internal", "service", *namespace, destFileName)

	err = copyFile(sourceFile, destFile)
	if err != nil {
		log.Fatalf("Error copying file: %v", err)
	}

	if err = modifyFile(destFile, *namespace, *resource); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}

	if err = formatFile(destFile); err != nil {
		log.Fatalf("Error formatting file: %v", err)
	}
}

func copyFile(src, dest string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, content, 0644)
}

func modifyFile(filePath, namespace, resource string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for modification: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	headers := "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	headers = headers + "\"\n\"" + "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/service"
	headers = headers + "\"\n" + "tferr \"gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err"

	imports := "import ("
	imports = imports + "\n\"" + "context\""

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.ReplaceAll(line, "package alicloud", "package "+namespace)
		line = strings.ReplaceAll(line, "import (", imports)
		line = strings.ReplaceAll(line, "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity", "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/connectivity")
		line = strings.ReplaceAll(line, "github.com/hashicorp/terraform-plugin-sdk/helper/resource", headers)
		line = strings.ReplaceAll(line, "github.com/hashicorp/terraform-plugin-sdk/helper/schema", "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")

		if strings.Contains(line, "Create:") {
			if !strings.Contains(line, "Create: schema") {
				line = strings.ReplaceAll(line, "Create:", "CreateContext:")
			}
		}
		if strings.Contains(line, "Read:") {
			if !strings.Contains(line, "Read: schema") {
				line = strings.ReplaceAll(line, "Read:", "ReadContext:")
			}
		}
		if strings.Contains(line, "Update:") {
			if !strings.Contains(line, "Update: schema") {
				line = strings.ReplaceAll(line, "Update:", "UpdateContext:")
			}
		}
		if strings.Contains(line, "Delete:") {
			if !strings.Contains(line, "Delete: schema") {
				line = strings.ReplaceAll(line, "Delete:", "DeleteContext:")
			}
		}

		line = strings.ReplaceAll(line, "tagsSchema()", "service.TagsSchema()")

		line = strings.ReplaceAll(line, "AliCloud", "ApsaraCloud")
		line = strings.ReplaceAll(line, "connectivity.AliyunClient", "connectivity.Client")
		line = strings.ReplaceAll(line, "(d *schema.ResourceData, meta interface{}) error {", "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {")

		line = strings.ReplaceAll(line, "State: schema.ImportStatePassthrough,", "StateContext: schema.ImportStatePassthroughContext,")
		line = strings.ReplaceAll(line, "buildClientToken", "helper.BuildClientToken")
		line = strings.ReplaceAll(line, "incrementalWait", "helper.IncrementalWait")
		line = strings.ReplaceAll(line, "resource.", "retry.")
		line = strings.ReplaceAll(line, "NeedRetry(err)", "tferr.NeedRetry(err)")
		line = strings.ReplaceAll(line, "addDebug", "helper.AddDebug")

		line = strings.ReplaceAll(line, "IdMsg", "tferr.IdMsg")
		line = strings.ReplaceAll(line, "WrapErrorf", "tferr.WrapErrorf")
		line = strings.ReplaceAll(line, "DefaultErrorMsg", "tferr.DefaultErrorMsg")
		line = strings.ReplaceAll(line, "AlibabaCloudSdkGoERROR", "tferr.SdkGoERROR")
		line = strings.ReplaceAll(line, "BuildStateConf", "helper.BuildStateConf")

		if strings.Contains(line, "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {") {
			line = line + "\nvar diags diag.Diagnostics\n"
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading file: %w", err)
	}

	fileOut, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file for writing: %w", err)
	}
	defer fileOut.Close()

	writer := bufio.NewWriter(fileOut)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error while writing to file: %w", err)
		}
	}
	writer.Flush()

	return nil
}

func formatFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	formattedContent, err := format.Source(content)
	if err != nil {
		return fmt.Errorf("格式化失败: %v", err)
	}

	return os.WriteFile(filePath, formattedContent, 0644)
}
