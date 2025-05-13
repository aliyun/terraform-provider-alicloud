package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/imports"
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
	functionName      = flag.String("f", "", "function name in service")
)

var specialServiceMap = map[string]string{
	"private_zone": "pvtz",
	"sls":          "log",
	"max_compute":  "maxcompute",
	"ack":          "cs",
}

var specialResourceMap = map[string]map[string]string{
	"vpc": {
		"vpc":                    "vpc",
		"vswitch":                "vswitch",
		"havip_attachment":       "havip_attachment",
		"network_acl":            "network_acl",
		"route_table":            "route_table",
		"route_table_attachment": "route_table_attachment",
	},
	"ecs": {
		"instance":                "instance",
		"security_group":          "security_group",
		"security_group_rule":     "security_group_rule",
		"auto_provisioning_group": "auto_provisioning_group",
		"image":                   "image",
		"image_copy":              "image_copy",
		"image_export":            "image_export",
		"image_import":            "image_import",
		"image_share_permission":  "image_share_permission",
		"ram_role_attachment":     "ram_role_attachment",
		"reserved_instance":       "reserved_instance",
	},
	"rds": {
		"instance": "db_instance",
	},
	"cbwp": {
		"common_bandwidth_package":            "common_bandwidth_package",
		"common_bandwidth_package_attachment": "common_bandwidth_package_attachment",
	},
	"private_zone": {
		"endpoint":               "pvtz_endpoint",
		"rule":                   "pvtz_rule",
		"rule_attachment":        "pvtz_rule_attachment",
		"user_vpc_authorization": "pvtz_user_vpc_authorization",
		"zone":                   "pvtz_zone",
		"zone_attachment":        "pvtz_zone_attachment",
		"zone_record":            "pvtz_zone_record",
	},
	"sls": {
		"project": "log_project",
		"store":   "log_store",
	},
	"nat_gateway": {
		"nat_gateway":   "nat_gateway",
		"forward_entry": "forward_entry",
		"snat_entry":    "snat_entry",
		"vpc_nat_ip":    "vpc_nat_ip",
		"nat_ip_cidr":   "vpc_nat_ip_cidr",
	},
	"max_compute": {
		"project": "maxcompute_project",
	},
	"ack": {
		"kubernetes":             "cs_kubernetes",
		"managed_kubernetes":     "cs_managed_kubernetes",
		"kubernetes_addon":       "cs_kubernetes_addon",
		"kubernetes_node_pool":   "cs_kubernetes_node_pool",
		"kubernetes_permissions": "cs_kubernetes_permissions",
	},
}

var specialDataSourceMap = map[string]map[string]string{
	"vpc": {
		"vpc":                    "vpcs",
		"vswitch":                "vswitches",
		"havip_attachment":       "havip_attachments",
		"network_acl":            "network_acls",
		"route_table":            "route_tables",
		"route_table_attachment": "route_table_attachments",
	},
	"ecs": {
		"instance":            "instances",
		"image":               "images",
		"security_group":      "security_groups",
		"security_group_rule": "security_group_rules",
	},
	"rds": {
		"instance": "db_instances",
	},
	"cbwp": {
		"common_bandwidth_package":            "common_bandwidth_packages",
		"common_bandwidth_package_attachment": "common_bandwidth_package_attachments",
	},
	"private_zone": {
		"endpoint":               "pvtz_endpoints",
		"rule":                   "pvtz_rules",
		"rule_attachment":        "pvtz_rule_attachments",
		"user_vpc_authorization": "pvtz_user_vpc_authorizations",
		"zone":                   "pvtz_zones",
		"zone_attachment":        "pvtz_zone_attachments",
		"zone_record":            "pvtz_zone_records",
	},
	"sls": {
		"project": "log_projects",
		"store":   "log_stores",
	},
	"nat_gateway": {
		"nat_gateway":   "nat_gateways",
		"forward_entry": "forward_entries",
		"snat_entry":    "snat_entries",
		"vpc_nat_ip":    "vpc_nat_ips",
		"nat_ip_cidr":   "vpc_nat_ip_cidrs",
	},
	"max_compute": {
		"project": "maxcompute_projects",
	},
	"ack": {
		"kubernetes":             "cs_kubernetes",
		"managed_kubernetes":     "cs_managed_kubernetes",
		"kubernetes_addon":       "cs_kubernetes_addons",
		"kubernetes_node_pool":   "cs_kubernetes_node_pools",
		"kubernetes_permissions": "cs_kubernetes_permissions",
	},
}

var irregularPlurals = map[string]string{
	"child":  "children",
	"person": "people",
	"sheep":  "sheep",
	"fish":   "fish",
}

func main() {
	flag.Parse()

	if *namespace == "" || *sourceProviderDir == "" || *destProviderDir == "" {
		log.Fatal("Parameters -n, -s, -t are required")
	}

	if functionName != nil && *functionName != "" {
		if err := migrateServiceFunction(namespace); err != nil {
			log.Printf("Error migrate function: %v", err)
		}
		return
	}

	var resources []string
	if *resource == "" {
		rs, err := listResources(*namespace, *sourceProviderDir)
		if err != nil {
			log.Fatal(err)
		}
		resources = rs
		log.Printf("Found %d resources to migrate: %v", len(resources), resources)
	} else {
		resources = []string{*resource}
	}

	for _, res := range resources {
		log.Printf("====== Migrating %s ======", res)
		migrateSingleResource(namespace, &res)
	}

	serviceName := *namespace
	if v, ok := specialServiceMap[*namespace]; ok {
		serviceName = v
	}
	serviceFileName := fmt.Sprintf("service_alicloud_%s.go", serviceName)
	if err := migrateService(serviceFileName, "v1"); err != nil {
		log.Printf("Error migrateService: %v", err)
	}

	serviceFileName = fmt.Sprintf("service_alicloud_%s_v2.go", *namespace)
	if err := migrateService(serviceFileName, "v2"); err != nil {
		log.Printf("Error migrateService: %v", err)
	}
	log.Println("All resources migrated!")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func listResources(namespace, sourceDir string) ([]string, error) {
	dir := filepath.Join(sourceDir, "alicloud")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir failed: %w", err)
	}

	pattern := regexp.MustCompile(fmt.Sprintf(`^resource_alicloud_%s_(.+?)\.go$`, regexp.QuoteMeta(namespace)))
	var resources []string

	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		if matches := pattern.FindStringSubmatch(entry.Name()); len(matches) == 2 {
			resources = append(resources, matches[1])
		}
	}

	if specialResources, ok := specialResourceMap[namespace]; ok {
		for res := range specialResources {
			if !contains(resources, res) {
				resources = append(resources, res)
			}
		}
	}

	return resources, nil
}

func migrateSingleResource(namespace, resource *string) {

	if err := migrateResource(namespace, resource); err != nil {
		log.Printf("Error migrateResource: %v", err)
	}

	if err := migrateDataSource(namespace, resource); err != nil {
		log.Printf("Error migrateDataSource: %v", err)
	}

	if err := migrateResourceTest(namespace, resource); err != nil {
		log.Printf("Error migrateResourceTest: %v", err)
	}

	if err := migrateDataSourceTest(namespace, resource); err != nil {
		log.Printf("Error migrateDataSource: %v", err)
	}

	if err := migrateResourceDocument(namespace, resource); err != nil {
		log.Printf("Error migrateResourceDocument: %v", err)
	}

	if err := migrateDataSourceDocument(namespace, resource); err != nil {
		log.Printf("Error migrateResourceDocument: %v", err)
	}

	log.Printf("Migrate %v finished! ^-^ ", *resource)
}

func migrateServiceFunction(namespace *string) error {
	serviceFiles := []struct {
		pattern string
		version string
	}{
		{fmt.Sprintf("service_alicloud_%s.go", *namespace), "v1"},
		// {fmt.Sprintf("service_alicloud_%s_v2.go", *namespace), "v2"},
	}

	for _, sf := range serviceFiles {
		sourcePath := filepath.Join(*sourceProviderDir, "alicloud", sf.pattern)
		destDir := filepath.Join(*destProviderDir, "internal", "service", *namespace)

		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("create service dir failed: %w", err)
		}

		destFile := "service_common.go"
		if sf.version == "v2" {
			destFile = "service_common_v2.go"
		}

		if err := copyAndModifyFunction(sourcePath, filepath.Join(destDir, destFile), *namespace, sf.version); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", sf.pattern, err)
		}
	}
	return nil
}

func functionExists(filePath string, funcName string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		log.Printf("Error parsing file %s: %v", filePath, err)
		return false
	}

	var exists bool
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.Name == funcName {
				exists = true
				return false
			}
		}
		return true
	})

	return exists
}

func copyAndModifyFunction(src, dest, namespace, version string) error {

	if functionExists(dest, *functionName) {
		log.Printf("Function %s already exists in %s, skip migration", *functionName, dest)
		return nil
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	var funcDecl *ast.FuncDecl
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == *functionName {
			funcDecl = fn
			return false
		}
		return true
	})

	if funcDecl == nil {
		return fmt.Errorf("function %s not found", *functionName)
	}

	var buf bytes.Buffer
	printer.Fprint(&buf, fset, funcDecl)

	code := applyFunctionModifications(buf.String(), version)

	f, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file failed: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + code + "\n"); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}

func applyFunctionModifications(code string, version string) string {
	lines := strings.Split(code, "\n")
	var modifiedLines []string
	for _, line := range lines {
		line = modifyServiceFunctionLine(line, version)
		modifiedLines = append(modifiedLines, line)
	}
	finalCode := strings.Join(modifiedLines, "\n")
	return finalCode
}

func migrateResource(namespace, resource *string) error {
	resourceName := getResourceName(*namespace, *resource)
	sourceFileName := fmt.Sprintf("resource_alicloud_%s.go", resourceName)
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

	resourceFunc, err := extractResourceFunction(destFilePath, *namespace, *resource, "Resource")
	resourceKey := getResourceName(*namespace, *resource)
	if err != nil {
		return fmt.Errorf("extractResourceFunction failed: %w", err)
	}

	if err = registerToProvider(resourceKey, resourceFunc, false); err != nil {
		log.Fatalf("Error registerToProvider: %v", err)
	}

	return err
}

func registerToProvider(resourceKey, resourceFunc string, datasource bool) error {
	providerPath := filepath.Join(*destProviderDir, "internal", "provider", "provider.go")
	return updateProviderMap(providerPath, *namespace, resourceKey, resourceFunc, datasource)
}

func extractResourceFunction(filePath, namespace, resource, prefix string) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	var resourceFunc string
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && strings.HasPrefix(fn.Name.Name, prefix) {
			resourceFunc = fmt.Sprintf("%s.%s()", namespace, fn.Name.Name)
			return false
		}
		return true
	})

	if resourceFunc == "" {
		return "", fmt.Errorf("no ResourceXXX function found in %s", filePath)
	}

	return resourceFunc, nil
}

func updateProviderMap(filePath, namespace, resourceKey, resourceFunc string, datasource bool) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read provider.go failed: %w", err)
	}

	re := regexp.MustCompile("\"" + resourceKey + "\"")
	containsKey := re.Match(content)
	if containsKey {
		return nil
	}

	targetComment := fmt.Sprintf("// %s", strings.ToUpper(namespace))
	insertLine := fmt.Sprintf("\t\t\"%s\":\t\t\t%s,", resourceKey, resourceFunc)

	var newContent []string
	scanner := bufio.NewScanner(bytes.NewReader(content))
	start := false
	foundSection := false
	inserted := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "import (") && !strings.Contains(string(content),
			fmt.Sprintf("service/%s\"", namespace)) {
			line += fmt.Sprintf("\n\t\"gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/service/%s\"",
				namespace)
		}

		if datasource {
			if strings.Contains(line, "func generateDataSourceMap()") {
				start = true
			}
			if strings.Contains(line, "func generateResourceMap()") {
				start = false
			}
		} else {
			if strings.Contains(line, "func generateResourceMap()") {
				start = true
			}
		}

		if start && strings.Contains(line, targetComment) {
			foundSection = true
		}

		if foundSection && !inserted {
			if strings.Contains(line, "//") || strings.TrimSpace(line) == "" {
				newContent = append(newContent, line)
			} else {
				newContent = append(newContent, insertLine)
				newContent = append(newContent, line)
				inserted = true
				foundSection = false
			}
		} else {
			newContent = append(newContent, line)
		}
	}

	output := strings.Join(newContent, "\n")
	formatted, err := format.Source([]byte(output))
	if err != nil {
		return fmt.Errorf("format failed: %w", err)
	}

	return os.WriteFile(filePath, formatted, 0644)
}

func migrateResourceDocument(namespace, resource *string) error {
	resourceName := getResourceName(*namespace, *resource)
	sourceFileName := fmt.Sprintf("%s.html.markdown", resourceName)
	sourceFilePath := fmt.Sprintf("%s/website/docs/r/%s", *sourceProviderDir, sourceFileName)

	targetDir := filepath.Join(*destProviderDir, "website", "docs", "r")

	destFileName := fmt.Sprintf("%s.html.markdown", resourceName)
	destFilePath := filepath.Join(targetDir, destFileName)

	err := copyFile(sourceFilePath, destFilePath)
	if err != nil {
		log.Fatalf("Error copying file: %v", err)
	}

	if err = modifyDocument(destFilePath, *namespace, *resource); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}
	return err
}

func migrateDataSourceDocument(namespace, resource *string) error {
	resourceName := getDataSourceName(*namespace, *resource)
	sourceFileName := fmt.Sprintf("%s.html.markdown", resourceName)
	sourceFilePath := fmt.Sprintf("%s/website/docs/d/%s", *sourceProviderDir, sourceFileName)

	targetDir := filepath.Join(*destProviderDir, "website", "docs", "d")

	destFileName := fmt.Sprintf("%s.html.markdown", resourceName)
	destFilePath := filepath.Join(targetDir, destFileName)

	err := copyFile(sourceFilePath, destFilePath)
	if err != nil {
		log.Printf("Error copying file: %v", err)
		log.Printf("Skipped copying file: %v", sourceFilePath)
		return nil
	}

	if err = modifyDocument(destFilePath, *namespace, *resource); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}
	return err
}

func migrateDataSource(namespace, resource *string) error {
	resourceName := getDataSourceName(*namespace, *resource)
	sourceFileName := fmt.Sprintf("data_source_alicloud_%s.go", resourceName)
	sourceFilePath := fmt.Sprintf("%s/alicloud/%s", *sourceProviderDir, sourceFileName)

	serviceDir := filepath.Join(*destProviderDir, "internal", "service", *namespace)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		log.Fatalf("Create Service Dir Failed: %v", err)
	}

	destFileName := fmt.Sprintf("%s.go", toPlural(*resource))
	destFilePath := filepath.Join(serviceDir, destFileName)

	err := copyFile(sourceFilePath, destFilePath)
	if err != nil {
		log.Printf("Error copying file: %v", err)
		log.Printf("Skipped copying file: %v", sourceFilePath)
		return nil
	}

	if err = modifyResourceFile(destFilePath, *namespace, *resource); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}

	if err = formatFile(destFilePath); err != nil {
		log.Fatalf("Error formatting file: %v", err)
	}

	resourceFunc, err := extractResourceFunction(destFilePath, *namespace, *resource, "DataSource")
	resourceKey := getDataSourceName(*namespace, *resource)
	if err != nil {
		return fmt.Errorf("extractResourceFunction failed: %w", err)
	}

	if err = registerToProvider(resourceKey, resourceFunc, true); err != nil {
		log.Fatalf("Error registerToProvider: %v", err)
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

	if _, err := os.Stat(destFilePath); err == nil {
		return nil
	}

	err := copyFile(sourceFilePath, destFilePath)
	if err != nil {
		log.Printf("Error copying file: %v", err)
		return nil
	}

	if err = modifyServiceFile(destFilePath, *namespace, version); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}

	if err = formatFile(destFilePath); err != nil {
		log.Fatalf("Error formatting file: %v", err)
	}
	return err
}

func migrateResourceTest(namespace, resource *string) error {
	resourceName := getResourceName(*namespace, *resource)
	sourceFileTestName := fmt.Sprintf("resource_alicloud_%s_test.go", resourceName)
	sourceFileTestPath := fmt.Sprintf("%s/alicloud/%s", *sourceProviderDir, sourceFileTestName)

	serviceDir := filepath.Join(*destProviderDir, "internal", "service", *namespace)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		log.Fatalf("Create Service Dir Failed: %v", err)
	}

	destFileTestName := fmt.Sprintf("%s_test.go", *resource)
	destFileTestPath := filepath.Join(serviceDir, destFileTestName)

	err := copyFile(sourceFileTestPath, destFileTestPath)
	if err != nil {
		log.Fatalf("Error copying file: %v", err)
	}

	if err = modifyResourceTestFile(destFileTestPath, *namespace, *resource); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}

	if err = formatFile(destFileTestPath); err != nil {
		log.Fatalf("Error formatting file: %v", err)
	}
	return err
}

func migrateDataSourceTest(namespace, resource *string) error {
	resourceName := getDataSourceName(*namespace, *resource)
	sourceFileTestName := fmt.Sprintf("data_source_alicloud_%s_test.go", resourceName)
	sourceFileTestPath := fmt.Sprintf("%s/alicloud/%s", *sourceProviderDir, sourceFileTestName)

	serviceDir := filepath.Join(*destProviderDir, "internal", "service", *namespace)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		log.Fatalf("Create Service Dir Failed: %v", err)
	}

	destFileTestName := fmt.Sprintf("%s_test.go", toPlural(*resource))
	destFileTestPath := filepath.Join(serviceDir, destFileTestName)

	err := copyFile(sourceFileTestPath, destFileTestPath)
	if err != nil {
		log.Printf("Error copying file: %v", err)
		log.Printf("Skipped copying file: %v", sourceFileTestPath)
		return nil
	}

	if err = modifyResourceTestFile(destFileTestPath, *namespace, *resource); err != nil {
		log.Fatalf("Error modifying file: %v", err)
	}

	if err = formatFile(destFileTestPath); err != nil {
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

func modifyDocument(filePath, namespace, resource string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for modification: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "alicloud", "apsara")
		line = strings.ReplaceAll(line, "Alibaba", "Apsara")
		line = strings.ReplaceAll(line, "Alicloud", "Apsara")
		line = strings.ReplaceAll(line, "AliCloud", "Apsara")
		line = strings.ReplaceAll(line, "alibabacloud", "apsaracloud")
		line = strings.ReplaceAll(line, "aliyun", "apsara")
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

func commonReplaces(line string) string {

	clientRe := regexp.MustCompile(`client\.Rpc([A-Za-z]+)\(\s*"([^"]+)",\s*"([^"]+)",\s*([^,]+),\s*([^,]+)(?:,\s*([^,]+))?(?:,\s*([^,]+))?(?:,\s*([^,]+))?(?:,\s*([^)]+))?\)`)
	if strings.Contains(line, "client.Rpc") {
		matches := clientRe.FindStringSubmatch(line)
		if len(matches) >= 6 {
			httpMethod := strings.ToUpper(matches[1])
			service := matches[2]
			version := matches[3]
			action := matches[4]
			pathParams := matches[5]
			request := "nil"
			options := "nil"
			async := "true"

			if len(matches) >= 7 && matches[6] != "" {
				request = matches[6]
			}
			if len(matches) >= 8 && matches[7] != "" {
				options = matches[7]
			}
			if len(matches) >= 9 && matches[8] != "" {
				async = matches[8]
			}

			newLine := fmt.Sprintf(
				`client.Do("%s", client.NewRpcParam("%s", "%s", %s), %s, %s, %s, nil, %s)`,
				service,
				httpMethod,
				version,
				action,
				pathParams,
				request,
				options,
				async,
			)
			line = strings.Replace(line, matches[0], newLine, 1)
		}
	}

	roaRe := regexp.MustCompile(`client\.Roa([A-Za-z]+)\(\s*"([^"]+)",\s*"([^"]+)",\s*([^,]+),\s*([^,]+)(?:,\s*([^,]+))?(?:,\s*([^,]+))?(?:,\s*([^)]+))?\)`)
	if strings.Contains(line, "client.Roa") {
		matches := roaRe.FindStringSubmatch(line)
		if len(matches) >= 7 {
			httpMethod := strings.ToUpper(matches[1])
			service := matches[2]
			version := matches[3]
			action := matches[4]
			query := matches[5]
			request := "nil"
			body := "nil"
			async := "true"

			if len(matches) >= 7 && matches[6] != "" {
				request = matches[6]
			}
			if len(matches) >= 8 && matches[7] != "" {
				body = matches[7]
			}
			if len(matches) >= 9 && matches[8] != "" {
				async = matches[8]
			}

			newLine := fmt.Sprintf(
				`client.Do("%s", client.NewRoaParam("%s", "%s", %s), %s, %s, %s, nil, %s)`,
				service,
				httpMethod,
				version,
				action,
				query,
				body,
				request,
				async,
			)
			line = strings.Replace(line, matches[0], newLine, 1)
		}
	}

	line = strings.ReplaceAll(line, "string(PostgreSQL)", "\"PostgreSQL\"")
	line = strings.ReplaceAll(line, "string(MySQL)", "\"MySQL\"")
	line = strings.ReplaceAll(line, "string(MongoDB)", "\"MongoDB\"")
	line = strings.ReplaceAll(line, "string(SQLServer)", "\"SQLServer\"")
	line = strings.ReplaceAll(line, "NormalMode", "\"normal\"")

	line = strings.ReplaceAll(line, "string(Month)", "string(names.Month)")
	line = strings.ReplaceAll(line, "string(Year)", "string(names.Year)")

	line = strings.ReplaceAll(line, "string(Postpaid)", "\"PostPaid\"")
	line = strings.ReplaceAll(line, "string(Prepaid)", "\"PrePaid\"")
	line = strings.ReplaceAll(line, "string(PostPaid)", "\"PostPaid\"")
	line = strings.ReplaceAll(line, "string(PrePaid)", "\"PrePaid\"")
	line = strings.ReplaceAll(line, "string(Serverless)", "\"Serverless\"")

	line = strings.ReplaceAll(line, "PrePaid,", "names.PrePaid,")
	line = strings.ReplaceAll(line, "PostPaid,", "names.PostPaid,")
	line = strings.ReplaceAll(line, "Prepaid,", "names.Prepaid,")
	line = strings.ReplaceAll(line, "Postpaid,", "names.Postpaid,")
	line = strings.ReplaceAll(line, "Serverless,", "names.Serverless,")

	line = strings.ReplaceAll(line, "PageNumSmall", "names.PageNumSmall")
	line = strings.ReplaceAll(line, "PageSizeSmall", "names.PageSizeSmall")
	line = strings.ReplaceAll(line, "PageSizeMedium", "names.PageSizeMedium")
	line = strings.ReplaceAll(line, "PageSizeLarge", "names.PageSizeLarge")
	line = strings.ReplaceAll(line, "PageSizeXLarge", "names.PageSizeXLarge")

	line = strings.ReplaceAll(line, "convertListToJsonString", "helper.ConvertListToJsonString")
	line = strings.ReplaceAll(line, "convertObjectToJsonString", "helper.ConvertObjectToJsonString")
	line = strings.ReplaceAll(line, "expandStringList", "helper.ExpandStringList")
	line = strings.ReplaceAll(line, "ParseResourceId", "helper.ParseResourceId")
	line = strings.ReplaceAll(line, "convertListMapToJsonString", "helper.ConvertListMapToJsonString")
	line = strings.ReplaceAll(line, "convertMaptoJsonString", "helper.ConvertMaptoJsonString")
	line = strings.ReplaceAll(line, "convertListStringToListInterface", "helper.ConvertListStringToListInterface")
	line = strings.ReplaceAll(line, "convertArrayObjectToJsonString", "helper.ConvertArrayObjectToJsonString")
	line = strings.ReplaceAll(line, "convertMapToJsonStringIgnoreError", "helper.ConvertMapToJsonStringIgnoreError")
	line = strings.ReplaceAll(line, "convertArrayToString", "helper.ConvertArrayToString")
	line = strings.ReplaceAll(line, "convertJsonStringToStringList", "helper.ConvertJsonStringToStringList")
	line = strings.ReplaceAll(line, "convertJsonStringToMap", "helper.ConvertJsonStringToMap")
	line = strings.ReplaceAll(line, "convertStringToBool", "helper.ConvertStringToBool")
	line = strings.ReplaceAll(line, "convertMapFloat64ToJsonString", "helper.ConvertMapFloat64ToJsonString")
	line = strings.ReplaceAll(line, "convertJsonStringToList", "helper.ConvertJsonStringToList")
	line = strings.ReplaceAll(line, "normalizeYamlString", "helper.NormalizeYamlString")
	line = strings.ReplaceAll(line, "compareYamlTemplateAreEquivalent", "helper.CompareYamlTemplateAreEquivalent")
	line = strings.ReplaceAll(line, "compareArrayJsonTemplateAreEquivalent", "helper.CompareArrayJsonTemplateAreEquivalent")
	line = strings.ReplaceAll(line, "expandSingletonToList", "helper.ExpandSingletonToList")
	line = strings.ReplaceAll(line, "filterEmptyStrings", "helper.FilterEmptyStrings")
	line = strings.ReplaceAll(line, "GetFunc", "helper.GetFunc")
	line = strings.ReplaceAll(line, "xmlParam", "client.NewXmlParam")
	line = strings.ReplaceAll(line, "SplitSlice", "helper.SplitSlice")
	line = strings.ReplaceAll(line, "decodeFromBase64String", "helper.DecodeFromBase64String")
	line = strings.ReplaceAll(line, "WaitTimeoutMsg", "tferr.WaitTimeoutMsg")
	line = strings.ReplaceAll(line, "COMMA_SEPARATED", "names.COMMA_SEPARATED")
	line = strings.ReplaceAll(line, "COLON_SEPARATED", "names.COLON_SEPARATED")
	line = strings.ReplaceAll(line, "LOCAL_HOST_IP", "names.LOCAL_HOST_IP")

	line = strings.ReplaceAll(line, "IsNil", "helper.IsNil")
	roaParamRe := regexp.MustCompile(`roaParam\(\s*("[^"]+"|\w+)\s*,\s*("[^"]+"|\w+)\s*,\s*("[^"]+"|\w+)\s*,\s*([^)]+)\s*\)`)
	line = roaParamRe.ReplaceAllString(line, `client.NewRpcParam($1, $2, $3)`)

	if !strings.Contains(line, "schema.DefaultTimeout") {
		line = strings.ReplaceAll(line, "DefaultTimeout", "names.DefaultTimeout")
	}

	if isVariable(line, "Throttling") {
		line = strings.ReplaceAll(line, "Throttling", "\"Throttling\"")
	}

	if isVariable(line, "RenewAutoRenewal") {
		line = strings.ReplaceAll(line, "RenewAutoRenewal", "\"AutoRenewal\"")
	}

	if isVariable(line, "RenewNormal") {
		line = strings.ReplaceAll(line, "RenewNormal", "\"Normal\"")
	}

	if isVariable(line, "RenewNotRenewal") {
		line = strings.ReplaceAll(line, "RenewNotRenewal", "\"NotRenewal\"")
	}

	if isVariable(line, "PayByTraffic") {
		line = strings.ReplaceAll(line, "PayByTraffic", "\"PayByTraffic\"")
	}

	if isVariable(line, "Deleted") {
		line = strings.ReplaceAll(line, "Deleted", "names.Deleted")
	}

	if isVariable(line, "Running") {
		line = strings.ReplaceAll(line, "Running", "names.Running")
	}

	if isVariable(line, "Month") {
		line = strings.ReplaceAll(line, "Month", "names.Month")
	}

	if isVariable(line, "Year") {
		line = strings.ReplaceAll(line, "Year", "names.Year")
	}

	if isVariable(line, "Status") {
		line = strings.ReplaceAll(line, "Status", "names.Status")
	}

	rdkPointerRe := regexp.MustCompile(`rdk\.StringPointer\(\s*d\.Get\("([^"]+)"\)(?:\.\(string\))?\s*\)`)
	line = rdkPointerRe.ReplaceAllString(line, `rdk.StringPointer(d.Get("$1").(string))`)

	if strings.Contains(line, " StringPointer") {
		line = strings.ReplaceAll(line, " StringPointer", " rdk.StringPointer")
	}

	rdkNestedRe := regexp.MustCompile(`rdk\.StringPointer\(\s*StringPointer\(([^)]+)\)\s*\)`)
	line = rdkNestedRe.ReplaceAllString(line, `rdk.StringPointer($1)`)

	createRequestRe := regexp.MustCompile(`request := (\w+)\.Create(\w+)Request\(\)`)
	line = createRequestRe.ReplaceAllString(line,
		"action := \"$2\"\n\trequest := make(map[string]interface{})")

	fieldAccessRe := regexp.MustCompile(`\b(request|req)\.(\w+)\b`)
	line = fieldAccessRe.ReplaceAllString(line, `${1}["${2}"]`)

	return line
}

func isVariable(line, code string) bool {
	codeRe := regexp.MustCompile(`^` + code + `$`)

	return codeRe.MatchString(line)
}

func skipUpdate(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return true
	}
	re := regexp.MustCompile(` Cannot update`)
	return re.Match(content)
}

func skipDelete(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return true
	}
	re := regexp.MustCompile(`(\[WARN\] Cannot destroy|Cannot delete)`)
	return re.Match(content)
}

func modifyResourceFile(filePath, namespace, resource string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for modification: %w", err)
	}
	defer file.Close()

	skippedUpdate := skipUpdate(filePath)
	skippedDelete := skipDelete(filePath)

	scanner := bufio.NewScanner(file)
	var lines []string

	headers := "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	headers = headers + "\"\n\"" + "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	headers = headers + "\"\n\"" + "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/rdk"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/names"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err/sdkdiag"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/service"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/helper"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/rdk"
	headers = headers + "\"\n" + "tferr \"gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err"

	imports := "import ("
	imports = imports + "\n\"" + "context\""

	resourceFuncRe := regexp.MustCompile(`func\s+resourceAli[Cc]loud(\w+)\s*\(\s*\)\s*\*schema\.Resource\s*\{`)
	dataSourceFuncRe := regexp.MustCompile(`func\s+dataSourceAli[Cc]loud(\w+)\s*\(\s*\)\s*\*schema\.Resource\s*\{`)
	clientRe := regexp.MustCompile(`client\.Rpc([A-Za-z]+)\(\s*"([^"]+)",\s*"([^"]+)",\s*([^,]+),\s*([^,]+),\s*([^,]+),\s*([^)]+)\)`)
	serviceRe := regexp.MustCompile(`([A-Z]\w*?)Service(V2)?\b`)

	queryAssignRe := regexp.MustCompile(`query\["([^"]+)"\]\s*=\s*([^;\n]+)`)
	diffSuppressRe := regexp.MustCompile(`(\w+)DiffSuppressFunc`)

	var (
		inValidateFunc bool
		parenDepth     int
	)

	for scanner.Scan() {
		line := scanner.Text()

		if inValidateFunc {
			parenDepth += strings.Count(line, "(")
			parenDepth -= strings.Count(line, ")")

			trimmed := strings.TrimSpace(line)
			if parenDepth <= 0 && strings.HasSuffix(trimmed, "),") {
				inValidateFunc = false
				parenDepth = 0
			}
			continue
		}

		if strings.Contains(line, "ValidateFunc") {
			inValidateFunc = true
			parenDepth += strings.Count(line, "(")
			parenDepth -= strings.Count(line, ")")

			trimmed := strings.TrimSpace(line)
			if parenDepth <= 0 && (strings.HasSuffix(trimmed, ",") || strings.HasSuffix(trimmed, ")")) {
				inValidateFunc = false
				parenDepth = 0
			}
			continue
		}

		if strings.Contains(line, "// Package alicloud.") {
			continue
		}

		if strings.Contains(line, "SetPartial") {
			continue
		}

		if strings.Contains(line, "helper/validation") {
			continue
		}

		if strings.Contains(line, "helper/encryption") {
			continue
		}

		if strings.Contains(line, "ConflictsWith") {
			continue
		}

		if strings.Contains(line, "Removed:") {
			continue
		}

		//if strings.Contains(line, "Deprecated:") {
		//	continue
		//}

		line = strings.ReplaceAll(line, "package alicloud", "package "+namespace)
		line = strings.ReplaceAll(line, "import (", imports)
		line = strings.ReplaceAll(line, "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity", "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/connectivity")
		line = strings.ReplaceAll(line, "github.com/hashicorp/terraform-plugin-sdk/helper/resource", "")
		line = strings.ReplaceAll(line, "github.com/hashicorp/terraform-plugin-sdk/helper/schema", headers)

		line = resourceFuncRe.ReplaceAllString(line, "func Resource$1() *schema.Resource {")
		line = dataSourceFuncRe.ReplaceAllString(line, "func DataSource$1() *schema.Resource {")

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
		line = strings.ReplaceAll(line, "tagsSchemaForceNew()", "service.TagsSchemaForceNew()")
		line = strings.ReplaceAll(line, "tagsSchemaComputed()", "service.TagsSchemaComputed()")

		if !strings.Contains(line, "func") {
			line = strings.ReplaceAll(line, "tagsToMap", "service.TagsToMap")
		}
		line = strings.ReplaceAll(line, "AliCloud", "")
		line = strings.ReplaceAll(line, "Alicloud", "")
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
		line = strings.ReplaceAll(line, "compareJsonTemplateAreEquivalent", "helper.CompareJsonTemplateAreEquivalent")
		line = strings.ReplaceAll(line, "ConvertTags", "service.ConvertTags")
		line = strings.ReplaceAll(line, "expandTagsToMap", "service.ExpandTagsToMap")
		line = strings.ReplaceAll(line, "InArray", "helper.InArray")

		line = strings.ReplaceAll(line, "parsingTags", "service.ParsingTags")
		line = strings.ReplaceAll(line, "ignoredTags", "service.IgnoredTags")
		if !strings.Contains(line, "strings.Trim") {
			line = strings.ReplaceAll(line, "Trim(", "helper.Trim(")
		}
		line = strings.ReplaceAll(line, "isPagingRequest", "helper.IsPagingRequest")
		line = strings.ReplaceAll(line, "formatInt", "helper.FormatInt")
		line = strings.ReplaceAll(line, "formatBool", "helper.FormatBool")

		line = strings.ReplaceAll(line, "IdMsg", "tferr.IdMsg")
		line = strings.ReplaceAll(line, "WrapError(", "tferr.WrapError(")
		line = strings.ReplaceAll(line, "WrapErrorf(", "tferr.WrapErrorf(")

		line = strings.ReplaceAll(line, " DataDefaultErrorMsg", " tferr.DefaultErrorMsg")
		line = strings.ReplaceAll(line, " DefaultErrorMsg", " tferr.DefaultErrorMsg")
		line = strings.ReplaceAll(line, " FailedGetAttributeMsg", " tferr.FailedGetAttributeMsg")
		line = strings.ReplaceAll(line, "AlibabaCloudSdkGoERROR", "tferr.SdkGoERROR")
		line = strings.ReplaceAll(line, "IsExpectedErrors", "tferr.IsExpectedErrors")
		line = strings.ReplaceAll(line, "NotFoundError", "tferr.NotFoundError")
		line = strings.ReplaceAll(line, "FailedToReachTargetStatus", "tferr.FailedToReachTargetStatus")
		line = strings.ReplaceAll(line, "BuildStateConf", "helper.BuildStateConf")

		line = diffSuppressRe.ReplaceAllStringFunc(line, func(m string) string {
			parts := diffSuppressRe.FindStringSubmatch(m)
			if len(parts) < 2 {
				return m
			}
			prefix := strings.ToUpper(string(parts[1][0])) + parts[1][1:]
			return "helper." + prefix + "DiffSuppressFunc"
		})

		line = strings.ReplaceAll(line, "hashcode.String", "helper.HashString")

		line = strings.ReplaceAll(line, "query := make(map[string]interface{})", "query := make(map[string]*string)")
		line = strings.ReplaceAll(line, "var query map[string]interface{}", "var query map[string]*string")
		line = strings.ReplaceAll(line, "query = make(map[string]interface{})", "query = make(map[string]*string)")
		line = queryAssignRe.ReplaceAllString(line, `query["$1"] = rdk.StringPointer($2)`)

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
			if (!strings.Contains(line, "Update") || !skippedUpdate) && (!strings.Contains(line, "Delete") || !skippedDelete) {
				line = line + "\nvar diags diag.Diagnostics\n"
			}
		}

		if strings.Contains(line, "client.Rpc") {
			matches := clientRe.FindStringSubmatch(line)

			if len(matches) == 8 {
				httpMethod := strings.ToUpper(matches[1])
				service := matches[2]
				version := matches[3]
				action := matches[4]
				pathParams := matches[5]
				request := matches[6]
				async := matches[7]

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
		line = serviceRe.ReplaceAllString(line, "Service$2")

		line = strings.ReplaceAll(line, "alicloud_", "apsara_")

		line = commonReplaces(line)

		line = strings.ReplaceAll(line, "dataResourceIdHash", "helper.DataResourceIdHash")
		line = strings.ReplaceAll(line, "writeToFile", "helper.WriteToFile")
		line = strings.ReplaceAll(line, "setPagingRequest", "helper.SetPagingRequest")

		line = strings.ReplaceAll(line, "(Error(", "(tferr.Error(")

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
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/service"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/names"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err/sdkdiag"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/helper"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/rdk"
	headers = headers + "\"\n" + "tferr \"gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/err"

	imports := "import ("
	imports = imports + "\n\"" + "context\""

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "// Package alicloud.") {
			continue
		}

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
		line = modifyServiceFunctionLine(line, version)

		if strings.Contains(line, "type ServiceV2 struct {") {
			lines = append(lines, "func NewServiceV2(client *connectivity.Client) *ServiceV2 {")
			lines = append(lines, "return &ServiceV2{client}")
			lines = append(lines, "}")
		} else if strings.Contains(line, "type Service struct {") {
			lines = append(lines, "func NewService(client *connectivity.Client) *Service {")
			lines = append(lines, "return &Service{client}")
			lines = append(lines, "}")
		}

		line = commonReplaces(line)

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

func modifyServiceFunctionLine(line, version string) string {

	clientRe := regexp.MustCompile(`client\.Rpc([A-Za-z]+)\(\s*"([^"]+)",\s*"([^"]+)",\s*([^,]+),\s*([^,]+),\s*([^,]+),\s*([^)]+)\)`)
	serviceRe := regexp.MustCompile(`([A-Z]\w*?)Service\b`)
	if version == "v2" {
		serviceRe = regexp.MustCompile(`([A-Z]\w*?)ServiceV2\b`)
	}

	queryAssignRe := regexp.MustCompile(`query\["([^"]+)"\]\s*=\s*([^;\n]+)`)

	stateRefreshDeclRe := regexp.MustCompile(`(func\s*\(.*?\)\s*\w+StateRefreshFunc)\s*\(`)
	stateRefreshCallRe := regexp.MustCompile(`(\w+)\.(\w+StateRefreshFunc)\s*\(`)
	diffSuppressRe := regexp.MustCompile(`(\w+)DiffSuppressFunc`)

	line = strings.ReplaceAll(line, "(d, meta)", "(ctx, d, meta)")

	line = strings.ReplaceAll(line, "tagsSchema()", "service.TagsSchema()")
	if !strings.Contains(line, "func") {
		line = strings.ReplaceAll(line, "tagsToMap", "service.TagsToMap")
	}
	line = strings.ReplaceAll(line, "AliCloud", "")
	line = strings.ReplaceAll(line, "connectivity.AliyunClient", "connectivity.Client")
	line = strings.ReplaceAll(line, "(d *schema.ResourceData, meta interface{}) error {", "(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {")

	line = strings.ReplaceAll(line, "State: schema.ImportStatePassthrough,", "StateContext: schema.ImportStatePassthroughContext,")
	line = strings.ReplaceAll(line, "buildClientToken", "helper.BuildClientToken")
	line = strings.ReplaceAll(line, "incrementalWait", "helper.IncrementalWait")
	line = strings.ReplaceAll(line, "resource.", "retry.")
	line = strings.ReplaceAll(line, "NeedRetry(err)", "tferr.NeedRetry(err)")
	line = strings.ReplaceAll(line, "addDebug", "helper.AddDebug")
	line = strings.ReplaceAll(line, "retry.Retry(", "retry.RetryContext(ctx, ")
	line = strings.ReplaceAll(line, "WaitForState()", "WaitForStateContext(ctx)")

	if strings.Contains(line, "StateRefreshFunc") {
		if matches := stateRefreshDeclRe.FindStringSubmatch(line); len(matches) > 0 {
			line = strings.Replace(line, matches[0], matches[1]+"(ctx context.Context, ", 1)
		}

		if matches := stateRefreshCallRe.FindStringSubmatch(line); len(matches) > 0 {
			line = strings.Replace(line, matches[0], matches[1]+"."+matches[2]+"(ctx, ", 1)
		}
	}

	line = strings.ReplaceAll(line, "convertListToCommaSeparate", "helper.ConvertListToCommaSeparate")
	line = strings.ReplaceAll(line, "ConvertTags", "service.ConvertTags")
	line = strings.ReplaceAll(line, "expandTagsToMap", "service.ExpandTagsToMap")
	line = strings.ReplaceAll(line, "InArray", "helper.InArray")
	line = strings.ReplaceAll(line, "parsingTags", "service.ParsingTags")
	line = strings.ReplaceAll(line, "ignoredTags", "service.IgnoredTags")
	line = strings.ReplaceAll(line, "Trim(", "helper.Trim(")
	line = strings.ReplaceAll(line, "isPagingRequest", "helper.IsPagingRequest")
	line = strings.ReplaceAll(line, "formatInt", "helper.FormatInt")
	line = strings.ReplaceAll(line, "formatBool", "helper.FormatBool")

	line = strings.ReplaceAll(line, "IdMsg", "tferr.IdMsg")
	line = strings.ReplaceAll(line, "WrapError(", "tferr.WrapError(")
	line = strings.ReplaceAll(line, "WrapErrorf(", "tferr.WrapErrorf(")

	line = strings.ReplaceAll(line, "DefaultErrorMsg", "tferr.DefaultErrorMsg")
	line = strings.ReplaceAll(line, "AlibabaCloudSdkGoERROR", "tferr.SdkGoERROR")
	line = strings.ReplaceAll(line, "IsExpectedErrors", "tferr.IsExpectedErrors")
	line = strings.ReplaceAll(line, "NotFoundError", "tferr.NotFoundError")
	line = strings.ReplaceAll(line, "NotFoundErr(", "tferr.NotFoundErr(")
	line = strings.ReplaceAll(line, "NotFoundMsg", "tferr.NotFoundMsg")
	line = strings.ReplaceAll(line, "ProviderERROR", "tferr.ProviderERROR")
	line = strings.ReplaceAll(line, "FailedGetAttributeMsg", "tferr.FailedGetAttributeMsg")
	line = strings.ReplaceAll(line, "NotFoundWithResponse", "tferr.NotFoundWithResponse")
	line = strings.ReplaceAll(line, "FailedToReachTargetStatus", "tferr.FailedToReachTargetStatus")
	line = strings.ReplaceAll(line, "BuildStateConf", "helper.BuildStateConf")
	line = strings.ReplaceAll(line, "(Error(", "(tferr.Error(")
	line = strings.ReplaceAll(line, "tferr.NotFoundMsg, tferr.ProviderERROR,", "tferr.NotFoundMsg,")

	line = diffSuppressRe.ReplaceAllStringFunc(line, func(m string) string {
		parts := diffSuppressRe.FindStringSubmatch(m)
		if len(parts) < 2 {
			return m
		}
		prefix := strings.ToUpper(string(parts[1][0])) + parts[1][1:]
		return "helper." + prefix + "DiffSuppressFunc"
	})

	line = strings.ReplaceAll(line, "hashcode.String", "helper.HashString")

	if strings.Contains(line, "client.Rpc") {
		matches := clientRe.FindStringSubmatch(line)

		if len(matches) == 8 {
			httpMethod := strings.ToUpper(matches[1])
			service := matches[2]
			version := matches[3]
			action := matches[4]
			pathParams := matches[5]
			request := matches[6]
			async := matches[7]

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
		line = serviceRe.ReplaceAllString(line, "Service$2")
	}

	line = strings.ReplaceAll(line, "query := make(map[string]interface{})", "query := make(map[string]*string)")
	line = strings.ReplaceAll(line, "var query map[string]interface{}", "var query map[string]*string")
	line = strings.ReplaceAll(line, "query = make(map[string]interface{})", "query = make(map[string]*string)")
	line = queryAssignRe.ReplaceAllString(line, `query["$1"] = rdk.StringPointer($2)`)

	line = strings.ReplaceAll(line, "RetryContext(ctx,", "RetryContext(context.Background(),")
	line = strings.ReplaceAll(line, "alicloud_", "apsara_")
	return line
}

func modifyResourceTestFile(filePath, namespace, resource string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for modification: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	headers := "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/connectivity"
	headers = headers + "\"\n\"" + "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/provider"
	headers = headers + "\"\n" + "tftest \"gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/acctest"
	headers = headers + "\"\n\"" + "gitlab.alibaba-inc.com/opensource-tools/terraform-provider-atlanta/internal/service/" + namespace

	serviceConstructRe := regexp.MustCompile(`&\b([A-Za-z]+)(Service)(V2)?\b`)

	skipCheckDestroyFunc := false
	testAccCheckDestroyRe := regexp.MustCompile(`testAccCheck(\w+)Destroy\b`)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "// Package alicloud.") {
			continue
		}

		if strings.Contains(line, "testAccPreCheckWithEnvVariable") {
			continue
		}

		if strings.Contains(line, "testAccPreCheckWithAccountSiteType") {
			continue
		}

		if strings.Contains(line, "testAccPreCheckWithTime") {
			continue
		}

		if strings.Contains(line, "testAccPreCheckEnterpriseAccountEnabled") {
			continue
		}

		if strings.Contains(line, "func testAccCheck") && strings.Contains(line, "(s *terraform.State) error {") {
			skipCheckDestroyFunc = true
			continue
		}

		if skipCheckDestroyFunc {
			if line == "}" {
				skipCheckDestroyFunc = false
			}
			continue
		}

		line = testAccCheckDestroyRe.ReplaceAllStringFunc(line, func(m string) string {
			parts := testAccCheckDestroyRe.FindStringSubmatch(m)
			if len(parts) < 2 {
				return m
			}
			return fmt.Sprintf("rac.CheckResourceDestroy()")
		})

		line = strings.ReplaceAll(line, "package alicloud", "package "+namespace+"_test")
		line = strings.ReplaceAll(line, "github.com/hashicorp/terraform-plugin-sdk/helper/acctest", headers)
		line = strings.ReplaceAll(line, "github.com/hashicorp/terraform-plugin-sdk/helper/resource", "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource")
		line = strings.ReplaceAll(line, "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity", "")
		line = strings.ReplaceAll(line, "AliCloud", "")
		line = strings.ReplaceAll(line, "Alicloud", "")
		line = strings.ReplaceAll(line, "alicloud_", "apsara_")
		line = strings.ReplaceAll(line, "resourceAttrInit", "tftest.ResourceAttrInit")
		line = strings.ReplaceAll(line, "resourceCheckInitWithDescribeMethod", "tftest.ResourceCheckInitWithDescribeMethod")
		line = serviceConstructRe.ReplaceAllString(line, namespace+".New$2$3")
		line = strings.ReplaceAll(line, "{testAccProvider.Meta().(*connectivity.AliyunClient)}", "(tftest.TestAccProvider.Meta().(*connectivity.Client))")
		line = strings.ReplaceAll(line, "resourceAttrCheckInit", "tftest.ResourceAttrCheckInit")
		line = strings.ReplaceAll(line, "resourceAttrMapUpdateSet", "ResourceAttrMapUpdateSet")
		line = strings.ReplaceAll(line, "defaultRegionToTest", "provider.DefaultRegionToTest")
		line = strings.ReplaceAll(line, "resourceTestAccConfigFunc", "tftest.ResourceTestAccConfig")
		line = strings.ReplaceAll(line, "resourceCheckInit", "tftest.ResourceCheckInit")
		line = strings.ReplaceAll(line, "testAccPreCheck(t)", "tftest.PreCheck(nil, t)")
		line = strings.ReplaceAll(line, "Providers:     testAccProviders", "ProtoV5ProviderFactories: tftest.ProtoV5ProviderFactories")
		line = strings.ReplaceAll(line, "Providers:    testAccProviders", "ProtoV5ProviderFactories: tftest.ProtoV5ProviderFactories")
		line = strings.ReplaceAll(line, "rac.checkResourceDestroy()", "rac.CheckResourceDestroy()")
		line = strings.ReplaceAll(line, "resourceAttrMapUpdateSet", "ResourceAttrMapUpdateSet")
		line = strings.ReplaceAll(line, "testAccPreCheckWithRegions", "tftest.TestAccPreCheckWithRegions")
		line = strings.ReplaceAll(line, "checkoutSupportedRegions", "tftest.TestAccPreCheckWithRegions")
		line = strings.ReplaceAll(line, "testAccPreCheckEnterpriseAccountEnabled", "tftest.TestAccPreCheckEnterpriseAccountEnabled")
		line = strings.ReplaceAll(line, "testAccPreCheckPrePaidResources", "tftest.TestAccPreCheckPrePaidResources")
		line = strings.ReplaceAll(line, "testAccPreCheckWithResourceManagerAccountsSetting", "tftest.TestAccPreCheckWithResourceManagerAccountsSetting")
		line = strings.ReplaceAll(line, "testAccPreCheckWithResourceManagerHandshakesSetting", "tftest.TestAccPreCheckWithResourceManagerHandshakesSetting")
		line = strings.ReplaceAll(line, "CHECKSET", "tftest.CHECKSET")
		line = strings.ReplaceAll(line, "REMOVEKEY", "tftest.REMOVEKEY")
		line = strings.ReplaceAll(line, "NOSET", "tftest.NOSET")

		line = strings.ReplaceAll(line, "dataSourceAttr", "tftest.DataSourceAttr")
		line = strings.ReplaceAll(line, "dataSourceTestAccConfig", "tftest.DataSourceTestAccConfig")
		line = strings.ReplaceAll(line, "existConfig:", "ExistConfig:")
		line = strings.ReplaceAll(line, "fakeConfig:", "FakeConfig:")
		line = strings.ReplaceAll(line, "resourceId:", "ResourceId:")
		line = strings.ReplaceAll(line, "existMapFunc:", "ExistMapFunc:")
		line = strings.ReplaceAll(line, "fakeMapFunc:", "FakeMapFunc:")
		line = strings.ReplaceAll(line, "dataSourceTestCheck", "DataSourceTestCheck")
		line = strings.ReplaceAll(line, "dataSourceTestAccConfig", "tftest.DataSourceTestAccConfig")

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

	if strings.Contains(filePath, "_test.go") {
		funcNames := []string{"init", "testSweep", "TestUnit"}
		content, err = removeFunctionsWithFuncNames(filePath, content, funcNames)
		if err != nil {
			return err
		}
	}

	formattedContent, err := imports.Process(filePath, content, nil)
	if err != nil {
		log.Printf("error while optimizing %s imports: %v", filePath, err.Error())
		formattedContent, err = format.Source(content)
		if err != nil {
			log.Printf("error while formatting to file: %v", err.Error())
			formattedContent = content
		}
	}

	return os.WriteFile(filePath, formattedContent, 0644)
}

func getResourceName(namespace, resource string) string {
	if productMap, ok := specialResourceMap[namespace]; ok {
		if mappedName, ok := productMap[resource]; ok {
			return mappedName
		}
	}

	return fmt.Sprintf("%s_%s", namespace, resource)
}

func getDataSourceName(namespace, resource string) string {

	if productMap, ok := specialDataSourceMap[namespace]; ok {
		if mappedName, ok := productMap[resource]; ok {
			return mappedName
		}
	}

	return fmt.Sprintf("%s_%s", namespace, toPlural(resource))
}

func removeFunctionsWithFuncNames(filename string, content interface{}, funcNames []string) ([]byte, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error while parsing AST: %w", err)
	}

	filteredDecls := make([]ast.Decl, 0, len(file.Decls))
	for _, decl := range file.Decls {
		if fnDecl, ok := decl.(*ast.FuncDecl); ok {
			if containsFuncNames(fnDecl.Name.Name, funcNames) {
				continue
			}
		}
		filteredDecls = append(filteredDecls, decl)
	}
	file.Decls = filteredDecls

	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, file); err != nil {
		return nil, fmt.Errorf("error while formatting to file: %w", err)
	}

	return buf.Bytes(), nil
}

func containsFuncNames(name string, funcNames []string) bool {
	for _, funcName := range funcNames {
		if strings.Contains(name, funcName) {
			return true
		}
	}

	return false
}

var pluralRules = []struct {
	suffix      string
	replacement string
}{
	{"ss", "sses"},
	{"s", "ses"},
	{"sh", "shes"},
	{"ch", "ches"},
	{"x", "xes"},
	{"z", "zes"},
	{"o", "oes"},
	{"f", "ves"},
	{"fe", "ves"},
	{"us", "i"},
	{"y", "ies"},
}

func toPlural(s string) string {
	if plural, exists := irregularPlurals[strings.ToLower(s)]; exists {
		return applyCase(s, plural)
	}

	lowerWord := strings.ToLower(s)
	for _, rule := range pluralRules {
		if strings.HasSuffix(lowerWord, rule.suffix) {
			if rule.suffix == "y" {
				if len(s) > 1 && !isVowel(rune(s[len(s)-2])) {
					return applyCase(s, s[:len(s)-1]+"ies")
				}
				continue
			}

			if rule.suffix == "o" {
				switch lowerWord {
				case "photo", "piano", "halo":
					return s + "s"
				}
			}

			return applyCase(s, s[:len(s)-len(rule.suffix)]+rule.replacement)
		}
	}

	return s + "s"
}

func applyCase(original, transformed string) string {
	if len(original) == 0 {
		return transformed
	}

	if original[0] >= 'A' && original[0] <= 'Z' {
		if len(transformed) > 0 {
			return strings.ToUpper(string(transformed[0])) + transformed[1:]
		}
	}
	return transformed
}

func isVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	}
	return false
}
