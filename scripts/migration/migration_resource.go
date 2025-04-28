package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	namespace         = flag.String("n", "", "namespace")
	resource          = flag.String("r", "", "resource")
	sourceProviderDir = flag.String("s", "", "source provider dir path")
	destProviderDir   = flag.String("t", "", "target provider dir path")
)

var specialResourceMap = map[string]map[string]string{
	"vpc": {
		"vpc":     "alicloud_vpc",
		"vswitch": "alicloud_vswitch",
	},
	"ecs": {
		"instance": "alicloud_instance",
	},
	"rds": {
		"instance": "alicloud_db_instance",
	},
}

func main() {

	flag.Parse()

	if *namespace == "" || *resource == "" || *destProviderDir == "" {
		log.Fatal("All parameters ( -n, -r, -t) are required.")
	}

	if err := migrateResource(namespace, resource); err != nil {
		log.Fatalf("Error migrateResource: %v", err)
	}

	serviceFileName := fmt.Sprintf("service_alicloud_%s.go", *namespace)
	if err := migrateService(serviceFileName, "v1"); err != nil {
		log.Fatalf("Error migrateService: %v", err)
	}

	serviceFileName = fmt.Sprintf("service_alicloud_%s_v2.go", *namespace)
	if err := migrateService(serviceFileName, "v2"); err != nil {
		log.Fatalf("Error migrateService: %v", err)
	}

	//if err = migrateResourceTest(sourceFile, destFile, namespace, resource); err != nil {
	//	log.Fatalf("Error modifying file: %v", err)
	//}
}

func migrateResource(namespace, resource *string) error {
	resourceName := getResourceName(*namespace, *resource)
	sourceFileName := fmt.Sprintf("resource_%s.go", resourceName)
	sourceFilePath := fmt.Sprintf("%s/alicloud/%s", *sourceProviderDir, sourceFileName)

	serviceDir := filepath.Join(*destProviderDir, "internal", "service", *namespace)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		log.Fatalf("Create Service Dir Failed: %v", err)
	}

	destFileName := fmt.Sprintf("%s.go", *resource)
	destFilePath := filepath.Join(serviceDir, destFileName)

	err := copyFile(sourceFilePath, destFilePath)
	if err != nil {
		log.Fatalf("Error copying file: %v", err)
	}

	if err = modifyResourceFile(destFilePath, *namespace, *resource); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}

	if err = formatFile(destFilePath); err != nil {
		log.Fatalf("Error formatting file: %v", err)
	}
	return err
}

func migrateService(serviceFileName, version string) error {

	sourceFilePath := filepath.Join(*sourceProviderDir, "alicloud", serviceFileName)

	if _, err := os.Stat(sourceFilePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Cannot find file: %s", sourceFilePath)
			return nil
		}
	}

	serviceDir := filepath.Join(*destProviderDir, "internal", "service", *namespace)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		log.Fatalf("Create Service Dir Failed: %v", err)
	}

	targetFileName := "service_common.go"
	if version == "v2" {
		targetFileName = "service_common_v2.go"
	}
	destFilePath := filepath.Join(serviceDir, targetFileName)

	err := copyFile(sourceFilePath, destFilePath)
	if err != nil {
		log.Fatalf("Error copying file: %v", err)
	}

	if err = modifyServiceFile(destFilePath, *namespace, version); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}

	if err = formatFile(destFilePath); err != nil {
		log.Fatalf("Error formatting file: %v", err)
	}
	return err
}

func getResourceName(namespace, resource string) string {
	if productMap, ok := specialResourceMap[namespace]; ok {
		if mappedName, ok := productMap[resource]; ok {
			return mappedName
		}
	}
	return fmt.Sprintf("alicloud_%s_%s", namespace, resource)
}

func migrateResourceTest(sourceFile, destFile string, namespace, resource *string) error {
	err := copyFile(sourceFile, destFile)
	if err != nil {
		log.Fatalf("Error copying file: %v", err)
	}

	if err = modifyResourceTestFile(destFile, *namespace, *resource); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}

	if err = formatFile(destFile); err != nil {
		log.Fatalf("Error formatting file: %v", err)
	}
	return err
}

func copyFile(src, dest string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, content, 0644)
}

func modifyResourceFile(filePath, namespace, resource string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for modification: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	headers := "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	headers = headers + "\"\n\"" + "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/names"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err/sdkdiag"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/service"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/helper"
	headers = headers + "\"\n" + "tferr \"gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err"

	imports := "import ("
	imports = imports + "\n\"" + "context\""

	clientRe := regexp.MustCompile(`client\.Rpc([A-Za-z]+)\(\s*"([^"]+)",\s*"([^"]+)",\s*([^,]+),\s*([^,]+),\s*([^,]+),\s*([^)]+)\)`)
	serviceRe := regexp.MustCompile(`([A-Z]\w*?)Service(V2)?\b`)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "SetPartial") {
			continue
		}

		if strings.Contains(line, "ValidateFunc") {
			continue
		}

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
		line = strings.ReplaceAll(line, "(d, meta)", "(ctx, d, meta)")

		line = strings.ReplaceAll(line, "tagsSchema()", "service.TagsSchema()")
		line = strings.ReplaceAll(line, "tagsToMap", "service.TagsToMap")
		line = strings.ReplaceAll(line, "AliCloud", "ApsaraCloud")
		line = strings.ReplaceAll(line, "connectivity.AliyunClient", "connectivity.Client")
		line = strings.ReplaceAll(line, "(d *schema.ResourceData, meta interface{}) error {", "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {")

		line = strings.ReplaceAll(line, "State: schema.ImportStatePassthrough,", "StateContext: schema.ImportStatePassthroughContext,")
		line = strings.ReplaceAll(line, "buildClientToken", "helper.BuildClientToken")
		line = strings.ReplaceAll(line, "incrementalWait", "helper.IncrementalWait")
		line = strings.ReplaceAll(line, "resource.", "retry.")
		line = strings.ReplaceAll(line, "NeedRetry(err)", "tferr.NeedRetry(err)")
		line = strings.ReplaceAll(line, "addDebug", "helper.AddDebug")
		line = strings.ReplaceAll(line, "retry.Retry(", "retry.RetryContext(ctx, ")
		line = strings.ReplaceAll(line, "StateRefreshFunc(", "StateRefreshFunc(ctx, ")
		line = strings.ReplaceAll(line, "WaitForState()", "WaitForStateContext(ctx)")

		line = strings.ReplaceAll(line, "convertListToCommaSeparate", "helper.ConvertListToCommaSeparate")
		line = strings.ReplaceAll(line, "ConvertTags", "service.ConvertTags")
		line = strings.ReplaceAll(line, "expandTagsToMap", "service.ExpandTagsToMap")
		line = strings.ReplaceAll(line, "InArray", "helper.InArray")

		line = strings.ReplaceAll(line, "IdMsg", "tferr.IdMsg")
		line = strings.ReplaceAll(line, "WrapError(", "tferr.WrapError(")
		line = strings.ReplaceAll(line, "WrapErrorf(", "tferr.WrapErrorf(")

		line = strings.ReplaceAll(line, "DefaultErrorMsg", "tferr.DefaultErrorMsg")
		line = strings.ReplaceAll(line, "AlibabaCloudSdkGoERROR", "tferr.SdkGoERROR")
		line = strings.ReplaceAll(line, "IsExpectedErrors", "tferr.IsExpectedErrors")
		line = strings.ReplaceAll(line, "NotFoundError", "tferr.NotFoundError")
		line = strings.ReplaceAll(line, "BuildStateConf", "helper.BuildStateConf")

		line = strings.ReplaceAll(line, "PostPaidDiffSuppressFunc", "helper.PostPaidDiffSuppressFunc")
		line = strings.ReplaceAll(line, "PostPaidAndRenewDiffSuppressFunc", "helper.PostPaidAndRenewDiffSuppressFunc")
		line = strings.ReplaceAll(line, "securityIpsDiffSuppressFunc", "helper.SecurityIpsDiffSuppressFunc")

		if strings.Contains(line, "return tferr.") {
			line = strings.ReplaceAll(line, "WrapError(", "tferr.WrapError(")
			line = strings.ReplaceAll(line, "WrapErrorf(", "tferr.WrapErrorf(")
			line = strings.ReplaceAll(line, "return tferr.", "return sdkdiag.AppendFromErr(diags,")
			line += ")"
		}

		if strings.TrimSpace(line) == "return err" {
			line = "sdkdiag.AppendFromErr(diags, tferr.WrapError(err))"
		}

		if strings.Contains(line, "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {") {
			line = line + "\nvar diags diag.Diagnostics\n"
		}

		if strings.Contains(line, "client.Rpc") {
			// 使用正则表达式动态提取方法类型和参数
			matches := clientRe.FindStringSubmatch(line)

			if len(matches) == 8 {
				httpMethod := strings.ToUpper(matches[1]) // 提取POST/GET等动词
				service := matches[2]
				version := matches[3]
				action := matches[4]
				pathParams := matches[5]
				request := matches[6]
				async := matches[7]

				// 重构参数结构
				newLine := fmt.Sprintf(
					`client.Do("%s", client.NewRpcParam("%s", "%s", %s), %s, %s, nil, nil, %s)`,
					service,
					httpMethod,
					version,
					action,
					pathParams,
					request,
					async,
				)

				line = strings.Replace(line, matches[0], newLine, 1)
			}
		}
		line = serviceRe.ReplaceAllString(line, "Service")

		line = strings.ReplaceAll(line, "alicloud_", "apsara_")

		line = strings.ReplaceAll(line, "PrePaid,", "names.PrePaid,")
		line = strings.ReplaceAll(line, "PostPaid,", "names.PostPaid,")
		line = strings.ReplaceAll(line, "Prepaid,", "names.Prepaid,")
		line = strings.ReplaceAll(line, "Postpaid,", "names.Postpaid,")
		line = strings.ReplaceAll(line, "Serverless,", "names.Serverless,")

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

func modifyServiceFile(filePath, namespace, version string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for modification: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	headers := "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	headers = headers + "\"\n\"" + "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/names"
	headers = headers + "\"\n" + "tferr \"gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err"

	imports := "import ("
	imports = imports + "\n\"" + "context\""

	clientRe := regexp.MustCompile(`client\.Rpc([A-Za-z]+)\(\s*"([^"]+)",\s*"([^"]+)",\s*([^,]+),\s*([^,]+),\s*([^,]+),\s*([^)]+)\)`)
	serviceRe := regexp.MustCompile(`([A-Z]\w*?)Service\b`)
	if version == "v2" {
		serviceRe = regexp.MustCompile(`([A-Z]\w*?)ServiceV2\b`)
	}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "SetPartial") {
			continue
		}

		if strings.Contains(line, "ValidateFunc") {
			continue
		}

		line = strings.ReplaceAll(line, "package alicloud", "package "+namespace)
		line = strings.ReplaceAll(line, "import (", imports)
		line = strings.ReplaceAll(line, "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity", "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/connectivity")
		line = strings.ReplaceAll(line, "github.com/hashicorp/terraform-plugin-sdk/helper/resource", headers)
		line = strings.ReplaceAll(line, "github.com/hashicorp/terraform-plugin-sdk/helper/schema", "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")

		line = strings.ReplaceAll(line, "(d, meta)", "(ctx, d, meta)")

		line = strings.ReplaceAll(line, "tagsSchema()", "service.TagsSchema()")
		line = strings.ReplaceAll(line, "tagsToMap", "service.TagsToMap")
		line = strings.ReplaceAll(line, "AliCloud", "ApsaraCloud")
		line = strings.ReplaceAll(line, "connectivity.AliyunClient", "connectivity.Client")
		line = strings.ReplaceAll(line, "(d *schema.ResourceData, meta interface{}) error {", "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {")

		line = strings.ReplaceAll(line, "State: schema.ImportStatePassthrough,", "StateContext: schema.ImportStatePassthroughContext,")
		line = strings.ReplaceAll(line, "buildClientToken", "helper.BuildClientToken")
		line = strings.ReplaceAll(line, "incrementalWait", "helper.IncrementalWait")
		line = strings.ReplaceAll(line, "resource.", "retry.")
		line = strings.ReplaceAll(line, "NeedRetry(err)", "tferr.NeedRetry(err)")
		line = strings.ReplaceAll(line, "addDebug", "helper.AddDebug")
		line = strings.ReplaceAll(line, "retry.Retry(", "retry.RetryContext(ctx, ")
		line = strings.ReplaceAll(line, "StateRefreshFunc(", "StateRefreshFunc(ctx, ")
		line = strings.ReplaceAll(line, "WaitForState()", "WaitForStateContext(ctx)")

		line = strings.ReplaceAll(line, "convertListToCommaSeparate", "helper.ConvertListToCommaSeparate")
		line = strings.ReplaceAll(line, "ConvertTags", "service.ConvertTags")
		line = strings.ReplaceAll(line, "expandTagsToMap", "service.ExpandTagsToMap")
		line = strings.ReplaceAll(line, "InArray", "helper.InArray")

		line = strings.ReplaceAll(line, "IdMsg", "tferr.IdMsg")
		line = strings.ReplaceAll(line, "WrapError(", "tferr.WrapError(")
		line = strings.ReplaceAll(line, "WrapErrorf(", "tferr.WrapErrorf(")

		line = strings.ReplaceAll(line, "DefaultErrorMsg", "tferr.DefaultErrorMsg")
		line = strings.ReplaceAll(line, "AlibabaCloudSdkGoERROR", "tferr.SdkGoERROR")
		line = strings.ReplaceAll(line, "IsExpectedErrors", "tferr.IsExpectedErrors")
		line = strings.ReplaceAll(line, "NotFoundError", "tferr.NotFoundError")
		line = strings.ReplaceAll(line, "BuildStateConf", "helper.BuildStateConf")

		line = strings.ReplaceAll(line, "PostPaidDiffSuppressFunc", "helper.PostPaidDiffSuppressFunc")
		line = strings.ReplaceAll(line, "PostPaidAndRenewDiffSuppressFunc", "helper.PostPaidAndRenewDiffSuppressFunc")
		line = strings.ReplaceAll(line, "securityIpsDiffSuppressFunc", "helper.SecurityIpsDiffSuppressFunc")

		if strings.Contains(line, "return tferr.") {
			line = strings.ReplaceAll(line, "WrapError(", "tferr.WrapError(")
			line = strings.ReplaceAll(line, "WrapErrorf(", "tferr.WrapErrorf(")
			line = strings.ReplaceAll(line, "return tferr.", "return sdkdiag.AppendFromErr(diags,")
			line += ")"
		}

		if strings.TrimSpace(line) == "return err" {
			line = "sdkdiag.AppendFromErr(diags, tferr.WrapError(err))"
		}

		if strings.Contains(line, "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {") {
			line = line + "\nvar diags diag.Diagnostics\n"
		}

		if strings.Contains(line, "client.Rpc") {
			// 使用正则表达式动态提取方法类型和参数
			matches := clientRe.FindStringSubmatch(line)

			if len(matches) == 8 {
				httpMethod := strings.ToUpper(matches[1]) // 提取POST/GET等动词
				service := matches[2]
				version := matches[3]
				action := matches[4]
				pathParams := matches[5]
				request := matches[6]
				async := matches[7]

				// 重构参数结构
				newLine := fmt.Sprintf(
					`client.Do("%s", client.NewRpcParam("%s", "%s", %s), %s, %s, nil, nil, %s)`,
					service,
					httpMethod,
					version,
					action,
					pathParams,
					request,
					async,
				)

				line = strings.Replace(line, matches[0], newLine, 1)
			}
		}
		if version == "v2" {
			line = serviceRe.ReplaceAllString(line, "ServiceV2")
		} else {
			line = serviceRe.ReplaceAllString(line, "Service")
		}

		line = strings.ReplaceAll(line, "alicloud_", "apsara_")

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

func modifyResourceTestFile(filePath, namespace, resource string) error {
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
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/helper"
	headers = headers + "\"\n" + "tferr \"gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err"

	imports := "import ("
	imports = imports + "\n\"" + "context\""

	clientRe := regexp.MustCompile(`client\.Rpc([A-Za-z]+)\(\s*"([^"]+)",\s*"([^"]+)",\s*([^,]+),\s*([^,]+),\s*([^,]+),\s*([^)]+)\)`)
	serviceRe := regexp.MustCompile(`([A-Z]\w*?)Service(V2)?\b`)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "SetPartial") {
			continue
		}

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
		line = strings.ReplaceAll(line, "(d, meta)", "(ctx, d, meta)")

		line = strings.ReplaceAll(line, "tagsSchema()", "service.TagsSchema()")
		line = strings.ReplaceAll(line, "tagsToMap", "service.TagsToMap")
		line = strings.ReplaceAll(line, "AliCloud", "ApsaraCloud")
		line = strings.ReplaceAll(line, "connectivity.AliyunClient", "connectivity.Client")
		line = strings.ReplaceAll(line, "(d *schema.ResourceData, meta interface{}) error {", "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {")

		line = strings.ReplaceAll(line, "State: schema.ImportStatePassthrough,", "StateContext: schema.ImportStatePassthroughContext,")
		line = strings.ReplaceAll(line, "buildClientToken", "helper.BuildClientToken")
		line = strings.ReplaceAll(line, "incrementalWait", "helper.IncrementalWait")
		line = strings.ReplaceAll(line, "resource.", "retry.")
		line = strings.ReplaceAll(line, "NeedRetry(err)", "tferr.NeedRetry(err)")
		line = strings.ReplaceAll(line, "addDebug", "helper.AddDebug")
		line = strings.ReplaceAll(line, "retry.Retry(", "retry.RetryContext(ctx, ")
		line = strings.ReplaceAll(line, "StateRefreshFunc(", "StateRefreshFunc(ctx, ")
		line = strings.ReplaceAll(line, "WaitForState()", "WaitForStateContext(ctx)")

		line = strings.ReplaceAll(line, "IdMsg", "tferr.IdMsg")
		line = strings.ReplaceAll(line, "WrapError(", "tferr.WrapError(")
		line = strings.ReplaceAll(line, "WrapErrorf(", "tferr.WrapErrorf(")

		line = strings.ReplaceAll(line, "DefaultErrorMsg", "tferr.DefaultErrorMsg")
		line = strings.ReplaceAll(line, "AlibabaCloudSdkGoERROR", "tferr.SdkGoERROR")
		line = strings.ReplaceAll(line, "IsExpectedErrors", "tferr.IsExpectedErrors")
		line = strings.ReplaceAll(line, "NotFoundError", "tferr.NotFoundError")
		line = strings.ReplaceAll(line, "BuildStateConf", "helper.BuildStateConf")

		if strings.Contains(line, "return tferr.") {
			line = strings.ReplaceAll(line, "WrapError(", "tferr.WrapError(")
			line = strings.ReplaceAll(line, "WrapErrorf(", "tferr.WrapErrorf(")
			line = strings.ReplaceAll(line, "return tferr.", "return sdkdiag.AppendFromErr(diags,")
			line += ")"
		}

		if strings.Contains(line, "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {") {
			line = line + "\nvar diags diag.Diagnostics\n"
		}

		if strings.Contains(line, "client.Rpc") {
			// 使用正则表达式动态提取方法类型和参数
			matches := clientRe.FindStringSubmatch(line)

			if len(matches) == 8 {
				httpMethod := strings.ToUpper(matches[1]) // 提取POST/GET等动词
				service := matches[2]
				version := matches[3]
				action := matches[4]
				pathParams := matches[5]
				request := matches[6]
				async := matches[7]

				// 重构参数结构
				newLine := fmt.Sprintf(
					`client.Do("%s", client.NewRpcParam("%s", "%s", %s), %s, %s, nil, nil, %s)`,
					service,
					httpMethod,
					version,
					action,
					pathParams,
					request,
					async,
				)

				line = strings.Replace(line, matches[0], newLine, 1)
			}
		}
		line = serviceRe.ReplaceAllString(line, "Service")

		line = strings.ReplaceAll(line, "alicloud_", "apsara_")

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
