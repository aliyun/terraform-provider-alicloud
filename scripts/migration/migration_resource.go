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
	namespace = flag.String("n", "", "namespace")
	resource  = flag.String("r", "", "resource")
	targetDir = flag.String("t", "", "target dir")
)

func main() {

	flag.Parse()

	if *namespace == "" || *resource == "" || *targetDir == "" {
		log.Fatal("All parameters ( -n, -r, -t) are required.")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	sourceDir := fmt.Sprintf("%s/alicloud", filepath.Dir(filepath.Dir(currentDir)))
	sourceFileName := fmt.Sprintf("alicloud_%s_%s.go", *namespace, *resource)
	if sourceFileName == "alicloud_vpc_vswitch.go" {
		sourceFileName = "alicloud_vswitch.go"
	}

	sourceFile := fmt.Sprintf("%s/%s", sourceDir, sourceFileName)

	destFileName := fmt.Sprintf("%s.go", *resource)
	destFile := filepath.Join(*targetDir, destFileName)

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

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.ReplaceAll(line, "Create:", "CreateContext:")
		line = strings.ReplaceAll(line, "Read:", "ReadContext:")
		line = strings.ReplaceAll(line, "Update:", "UpdateContext:")
		line = strings.ReplaceAll(line, "Delete:", "DeleteContext:")
		line = strings.ReplaceAll(line, "AliCloud", "ApsaraCloud")
		line = strings.ReplaceAll(line, "State: schema.ImportStatePassthrough,", "StateContext: schema.ImportStatePassthroughContext,")
		line = strings.ReplaceAll(line, "buildClientToken", "helper.BuildClientToken")
		line = strings.ReplaceAll(line, "incrementalWait", "helper.IncrementalWait")
		line = strings.ReplaceAll(line, "resource.", "retry.")
		line = strings.ReplaceAll(line, "NeedRetry(err)", "tferr.NeedRetry(err)")
		line = strings.ReplaceAll(line, "addDebug", "helper.AddDebug")

		if strings.Contains(line, "client := meta.(*connectivity.AliyunClient)") {
			lines = append(lines, "var diags diag.Diagnostics")
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
