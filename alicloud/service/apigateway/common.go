package apigateway

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// Page sizes used by the CloudAPI DescribeSignatures paginated read. They mirror the
// PageSizeSmall/PageSizeLarge constants the flat alicloud package uses, kept local so
// the service package does not import the legacy package.
const (
	pageSizeSmall = 10
	pageSizeLarge = 50
)

// nextPageNumber increments the CloudAPI page number the way the flat package's
// getNextpageNumber does, so the DescribeSignatures list walk matches the SDKv2 data
// source's paging exactly.
func nextPageNumber(number requests.Integer) (requests.Integer, error) {
	page, err := strconv.Atoi(string(number))
	if err != nil {
		return "", err
	}
	return requests.NewInteger(page + 1), nil
}

// userHomeDir returns the current user's home directory. os.UserHomeDir is the standard
// resolver; the indirection keeps writeToFile testable if a test overrides it.
func userHomeDir() (string, error) {
	return os.UserHomeDir()
}

// writeToFile writes data as indented JSON to filePath, expanding a leading ~ and
// replacing any existing file, matching the flat package's writeToFile for the
// alicloud_api_gateway_signatures output_file argument.
func writeToFile(filePath string, data interface{}) error {
	if data == nil {
		return nil
	}
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if strings.HasPrefix(filePath, "~") {
		home, err := userHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, out, 0o644)
}
