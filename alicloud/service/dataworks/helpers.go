package dataworks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

// dataResourceIdHash mirrors the SDKv2 dataResourceIdHash: a stable hash of the
// returned ids so the data source instance id is deterministic for a given
// result set.
func dataResourceIdHash(ids []string) string {
	var buf bytes.Buffer
	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}
	return fmt.Sprintf("%d", helper.Hashcode(buf.String()))
}

// formatInt coerces a JSON-decoded numeric to int without panicking. The SDK
// decode path may surface json.Number, float64, int64 or a string; unknown
// types fall back to 0 instead of crashing the data source Read.
func formatInt(src interface{}) int {
	if src == nil {
		return 0
	}
	switch v := src.(type) {
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int:
		return v
	case json.Number:
		n, _ := strconv.Atoi(v.String())
		return n
	case string:
		n, _ := strconv.Atoi(v)
		return n
	default:
		n, _ := strconv.Atoi(fmt.Sprint(src))
		return n
	}
}

// asString coerces a JSON-decoded value to string, tolerating json.Number and
// numeric types so TableGuid/TableName/DatabaseName do not silently drop.
func asString(src interface{}) string {
	if src == nil {
		return ""
	}
	switch v := src.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	default:
		return fmt.Sprint(src)
	}
}

// writeToFile writes the lineages as indented JSON to output, expanding a
// leading "~" to the user's home dir, mirroring the SDKv2 writeToFile.
func writeToFile(output string, entries []tableLineageEntry) error {
	data := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		data = append(data, map[string]interface{}{
			"id":               e.Id.ValueString(),
			"table_guid":       e.TableGuid.ValueString(),
			"table_name":       e.TableName.ValueString(),
			"database_name":    e.DatabaseName.ValueString(),
			"create_timestamp": e.CreateTimestamp.ValueInt64(),
		})
	}
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal lineage output: %w", err)
	}

	path := output
	if len(path) > 0 && path[0] == '~' {
		home, herr := os.UserHomeDir()
		if herr == nil && home != "" {
			path = filepath.Join(home, path[1:])
		}
	}
	// Replace any existing file so re-plans do not accumulate stale results.
	if _, statErr := os.Stat(path); statErr == nil {
		_ = os.Remove(path)
	}
	return os.WriteFile(path, out, 0o644)
}
