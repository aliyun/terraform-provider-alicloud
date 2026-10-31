package ecd_test

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain delegates to the testing framework's TestMain so resource.Test
// sweepers and the SDK provider factories are wired up for the ecd test suite.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

var ecdTestEndpoints sync.Map
var ecdTestSignVersion sync.Map

// ecdTestClientForRegion rebuilds the sweeper/raw client on the exported
// connectivity API: the flat-suite sharedClientForRegion is not reachable from
// this external _test package.
func ecdTestClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &ecdTestEndpoints,
		SignVersion: &ecdTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
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

	return conf.Client()
}

// ecdTestRegion returns the region under test, honoring ALICLOUD_REGION.
func ecdTestRegion() string {
	if region := os.Getenv("ALICLOUD_REGION"); region != "" {
		return region
	}
	return "cn-hangzhou"
}

// ecdString is a local *string convenience mirroring the flat-suite helper of
// the same name, which is unreachable from this external _test package.
func ecdString(s string) *string {
	return &s
}

// ecdMapMerge merges source into target in place, mirroring the flat-suite
// mapMerge helper that is unreachable from this external _test package.
func ecdMapMerge(target, source map[string]interface{}) {
	for k, v := range source {
		target[k] = v
	}
}

// testAccPreCheckWithTime mirrors the flat-suite gate of the same name
// (alicloud/provider_test.go), which is unreachable from this external
// _test package. It skips the test unless the current day is in days.
func testAccPreCheckWithTime(t *testing.T, days []int) {
	skipped := true
	for _, d := range days {
		if time.Now().Day() == d {
			skipped = false
			break
		}
	}
	if skipped {
		t.Skipf("Skipping the test case when not in specified days %#v of every month", days)
		t.Skipped()
	}
}
