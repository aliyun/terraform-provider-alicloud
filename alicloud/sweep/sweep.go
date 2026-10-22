// Package sweep provides shared sweeper client setup and filtering helpers.
// Test runners, registrations, and resource-specific cleanup remain with callers.
package sweep

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// SweepAll reports whether sweeping all resources is explicitly enabled.
func SweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// SweepOlderThan returns the minimum age a resource must reach before the
// sweeper is allowed to delete it. It is driven by the ALICLOUD_SWEEP_OLDER_THAN
// environment variable, a Go time.Duration string (for example "12h" or "24h").
// A missing, empty or invalid value returns 0, which disables the age filter
// and preserves the legacy sweep behavior.
func SweepOlderThan() time.Duration {
	v := strings.TrimSpace(os.Getenv("ALICLOUD_SWEEP_OLDER_THAN"))
	if v == "" {
		return 0
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("[WARN] Invalid ALICLOUD_SWEEP_OLDER_THAN %q (expected a Go duration such as 12h): %s; age filter disabled.", v, err)
		return 0
	}
	return d
}

// SweepRecentCheck inspects a list-API item and decides whether the sweeper
// should spare it because it is too new. It returns skip=true with a
// human-readable reason when the resource was created within the
// SweepOlderThan() window, or when its creation time cannot be established
// (the conservative, protective default). It returns skip=false (proceed to
// delete) when the age filter is disabled or the resource is older than the
// threshold. The caller logs the reason when skipping.
func SweepRecentCheck(item map[string]interface{}) (skip bool, reason string) {
	threshold := SweepOlderThan()
	if threshold <= 0 {
		return false, ""
	}
	raw := ""
	for _, k := range []string{"CreationTime", "CreateTime", "creation_time", "create_time"} {
		if v, ok := item[k]; ok && v != nil {
			if s := fmt.Sprint(v); s != "" {
				raw = s
				break
			}
		}
	}
	if raw == "" {
		return true, fmt.Sprintf("no creation time field found; ALICLOUD_SWEEP_OLDER_THAN=%s set, protecting", threshold)
	}
	created, ok := ParseSweepCreationTime(raw)
	if !ok {
		return true, fmt.Sprintf("cannot parse creation time %q; ALICLOUD_SWEEP_OLDER_THAN=%s set, protecting", raw, threshold)
	}
	if age := time.Since(created); age < threshold {
		return true, fmt.Sprintf("created %s ago, newer than ALICLOUD_SWEEP_OLDER_THAN=%s", age.Truncate(time.Second), threshold)
	}
	return false, ""
}

// ParseSweepCreationTime parses the timestamp formats Alibaba Cloud list APIs
// commonly return. The VPC DescribeVpcs API, for example, returns CreationTime
// as RFC3339 (e.g. 2021-04-18T15:02:37Z). The remaining layouts cover the
// no-seconds and space-separated variants other services use.
func ParseSweepCreationTime(s string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339,             // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04Z07:00", // no seconds (provider's UnStandardRFC3339)
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05", // ECS-style wall-clock, assumed UTC
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// SweepPreCheckWithRegions reports whether to skip the region. When supported
// is true, regions is an allowlist; otherwise, it is a list of excluded regions.
func SweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}
