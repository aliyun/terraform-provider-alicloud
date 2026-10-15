package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedClientForRegion returns a common AlicloudClient setup needed for the sweeper
// functions for a given region
func sharedClientForRegionWithBackendRegions(region string, supported bool, regions []connectivity.Region) (interface{}, error) {
	find := false
	backupRegion := string(connectivity.APSouthEast1)
	backupRegionFind := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
		if string(r) == backupRegion {
			backupRegionFind = true
		}
	}

	if (find && !supported) || (!find && supported) {
		if supported {
			// Try to use backup region first
			if backupRegionFind {
				log.Printf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, backupRegion)
				region = backupRegion
			} else if len(regions) > 0 {
				// If backup region is not available, use the first supported region
				targetRegion := string(regions[0])
				log.Printf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, targetRegion)
				region = targetRegion
			}
		} else {
			if !backupRegionFind {
				log.Printf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
				region = backupRegion
			}
		}
	}
	return sharedClientForRegion(region)
}

var endpoints sync.Map
var signVersion sync.Map

// sharedClientForRegion returns a common AlicloudClient setup needed for the sweeper
// functions for a give n region
func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, securityToken string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}

	securityToken = os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &endpoints,
		SignVersion: &signVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	accountId := os.Getenv("ALICLOUD_ACCOUNT_ID")
	if accountId == "" {
		accountId = os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID")
	}
	if accountId != "" {
		conf.AccountId = accountId
	}
	accountType := os.Getenv("ALICLOUD_ACCOUNT_TYPE")
	if accountType == "" {
		accountType = os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")
	}
	if accountType != "" {
		conf.AccountType = accountType
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func sweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// sweepOlderThan returns the minimum age a resource must reach before the
// sweeper is allowed to delete it. It is driven by the ALICLOUD_SWEEP_OLDER_THAN
// environment variable, a Go time.Duration string (for example "12h" or "24h").
// A missing, empty or invalid value returns 0, which disables the age filter
// and preserves the legacy sweep behavior.
func sweepOlderThan() time.Duration {
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

// sweepRecentCheck inspects a list-API item and decides whether the sweeper
// should spare it because it is too new. It returns skip=true with a
// human-readable reason when the resource was created within the
// sweepOlderThan() window, or when its creation time cannot be established
// (the conservative, protective default). It returns skip=false (proceed to
// delete) when the age filter is disabled or the resource is older than the
// threshold. The caller logs the reason when skipping.
func sweepRecentCheck(item map[string]interface{}) (skip bool, reason string) {
	threshold := sweepOlderThan()
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
	created, ok := parseSweepCreationTime(raw)
	if !ok {
		return true, fmt.Sprintf("cannot parse creation time %q; ALICLOUD_SWEEP_OLDER_THAN=%s set, protecting", raw, threshold)
	}
	if age := time.Since(created); age < threshold {
		return true, fmt.Sprintf("created %s ago, newer than ALICLOUD_SWEEP_OLDER_THAN=%s", age.Truncate(time.Second), threshold)
	}
	return false, ""
}

// parseSweepCreationTime parses the timestamp formats Alibaba Cloud list APIs
// commonly return. The VPC DescribeVpcs API, for example, returns CreationTime
// as RFC3339 (e.g. 2021-04-18T15:02:37Z). The remaining layouts cover the
// no-seconds and space-separated variants other services use.
func parseSweepCreationTime(s string) (time.Time, bool) {
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
