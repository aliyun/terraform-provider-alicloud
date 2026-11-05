package vpc_test

// Shared acceptance-test helpers for the split VPC package.
//
// The generic reflection core (resourceCheck / resourceAttr / resourceAttrCheck
// / resourceConfig config-builders, the sentinels and the valueConvert engine)
// is NOT duplicated here: the migrated VPC tests call alicloud/acctest's exported
// forms (ResourceCheckInit, ResourceAttrInit, ResourceAttrCheckInit,
// ResourceTestAccConfig, DataSourceTestAccConfigFunc and the exported methods
// CheckResourceDestroy / ResourceAttrMapUpdateSet) directly. acctest already
// links the single terraform-plugin-testing harness the whole binary shares, so
// routing the VPC callers to it removes ~650 lines of byte-identical copy
// without changing any behaviour: the public API is the same reflection engine
// the root originals were, the service constructors / DescribeMethod strings
// stay at each caller, and the sentinel values (acctest.NOSET / CHECKSET /
// REMOVEKEY) are the same strings the local consts held.
//
// What stays local here is the VPC product surface acctest has no public
// equivalent for, or whose runner wiring differs: the SDKv2-shaped
// testAccProviderFactory (= acctest.Provider), the Provider() shim, the
// defaultRegionToTest / PreCheck / region-checkout helpers, the shared client +
// endpoint cache + sweep guards, and the dataSource runner. The dataSource runner
// (dataSourceTestCheck / dataSourceTestCheckWithPreCheck) keeps using this
// package's testAccPreCheck and testAccProviderFactory rather than acctest's
// PreCheck + SDKProviderFactories: testAccPreCheck writes this package's
// defaultRegionToTest (which the VPC fixtures read at entry time), so routing the
// runner to acctest's hardcoded PreCheck would detach that write. The
// dataSourceAttr / dataSourceTestAccConfig data types and their buildDataSourceSteps
// step builder therefore remain local; only checkDataSourceAttr delegates the
// attribute-existence/empty check construction to acctest.ResourceAttrInit +
// ResourceAttrMapUpdateSet.
//
// The migrated tests live in the EXTERNAL test package vpc_test. An internal
// package vpc test file importing acctest forms an import cycle
// (vpc test -> acctest -> alicloud -> provider/registry -> service/vpc), so the
// whole suite moved to vpc_test, which is outside package vpc's import graph and
// may import acctest freely. The runner is terraform-plugin-testing for every
// alicloud test binary; the baseline root VPC suite supplied providers via the
// SDKv2-shaped ProviderFactories: testAccProviderFactory local (mirrored below),
// never the ProtoV5 flip, so this file must not reintroduce the SDKv2
// helper/resource package (helper/schema alone is safe — the *schema.Provider
// reflection needs it and it registers no flags).

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ---------------------------------------------------------------------------
// Provider wiring (mirrors alicloud/provider_test.go init()). acctest.Provider
// is alicloud.Provider(); the SDKv2 *schema.Provider shape is preserved so the
// resourceCheck reflection machinery can still read ResourcesMap/Schema.
// Acceptance TestCases in the migrated suite supply providers via the legacy
// SDKv2-shaped ProviderFactories: testAccProviderFactory field (the baseline
// root VPC suite used this exact local factory, never the ProtoV5 flip; the
// multi-account tests keep their own per-test providerFactories map). The
// framework-flipped ProtoV5ProviderFactories field is therefore not used here.
// ---------------------------------------------------------------------------

var testAccProvider = acctest.Provider

// testAccProviderFactory mirrors alicloud/provider_test.go's baseline var
// (declared there at line 35, initialized at line 113). The migrated VPC
// acceptance tests reference it by the same name so their TestCase factory
// fields stay byte-identical to the baseline root suite.
var testAccProviderFactory = map[string]func() (*schema.Provider, error){
	"alicloud": func() (*schema.Provider, error) {
		return testAccProvider, nil
	},
}

// Provider is a test-only shim for the root alicloud.Provider() symbol the
// migrated VPC tests call directly (e.g. Provider().ResourcesMap). package vpc
// cannot import root alicloud, so it returns the cached acctest.Provider
// instance (= alicloud.Provider()); the schema is static, which is all the
// moved tests read from it. It must not be used to instantiate fresh
// per-test providers.
func Provider() *schema.Provider {
	return alicloud.Provider()
}

// defaultRegionToTest snapshots ALICLOUD_REGION at package init, falling back
// to cn-beijing, mirroring the root suite's defaultRegionToTest: test fixtures
// read it at test-function-entry time, before testAccPreCheck runs. PreCheck,
// testAccPreCheckWithRegions and checkoutSupportedRegions write it back when
// the region changes, but code that captured the old value before a region
// jump keeps seeing it — the root suite has the same entry-time gap.
var defaultRegionToTest = func() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	return "cn-beijing"
}()

// ---------------------------------------------------------------------------
// PreCheck helpers (mirrors alicloud/provider_test.go).
// ---------------------------------------------------------------------------

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
		defaultRegionToTest = "cn-beijing"
	} else {
		defaultRegionToTest = v
	}
}

// Skip automatically the testcases which does not support some known regions.
// If supported is true, the regions should a list of supporting the service regions.
// If supported is false, the regions should a list of unsupporting the service regions.
// If the region is unsupported and has backend region, the backend region will instead
func testAccPreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		t.Logf("[WARNING] The region is not set and using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
	}
	checkoutSupportedRegions(t, supported, regions)
	defaultRegionToTest = os.Getenv("ALICLOUD_REGION")
}

func checkoutSupportedRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	region := os.Getenv("ALICLOUD_REGION")
	find := false
	backupRegion := string(connectivity.APSouthEast1)
	if region == string(connectivity.APSouthEast1) {
		backupRegion = string(connectivity.EUCentral1)
	}

	checkoutRegion := os.Getenv("CHECKOUT_REGION")
	if checkoutRegion == "true" {
		if region == string(connectivity.Hangzhou) {
			region = string(connectivity.EUCentral1)
			os.Setenv("ALICLOUD_REGION", region)
		}
	}
	backupRegionFind := false
	hangzhouRegionFind := false
	for _, r := range regions {
		if region == string(r) {
			find = true
		}
		if string(r) == backupRegion {
			backupRegionFind = true
		}
		if string(connectivity.Hangzhou) == string(r) {
			hangzhouRegionFind = true
		}
	}

	if (find && !supported) || (!find && supported) {
		if supported {
			// Try to use backup region first
			if backupRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				defaultRegionToTest = backupRegion
				return
			}
			// Try to use Hangzhou region
			if hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				defaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			// If neither backup region nor Hangzhou region is available, use the first supported region
			if len(regions) > 0 {
				targetRegion := string(regions[0])
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, targetRegion)
				os.Setenv("ALICLOUD_REGION", targetRegion)
				defaultRegionToTest = targetRegion
				return
			}
			t.Skipf("Skipping unsupported region %s. Supported regions: %s.", region, regions)
		} else {
			if !backupRegionFind {
				t.Logf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				defaultRegionToTest = backupRegion
				return
			}
			if !hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				defaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			t.Skipf("Skipping unsupported region %s. Unsupported regions: %s.", region, regions)
		}
		t.Skipped()
	}
}

func testAccPreCheckWithEnvVariable(t *testing.T, envVariableName string) {
	if v := strings.TrimSpace(os.Getenv(envVariableName)); v == "" {
		t.Skipf("Skipping the test case with no env variable %s", envVariableName)
		t.Skipped()
	}
}

func testAccPreCheckWithMultipleAccount(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_2")); v == "" {
		t.Skipf("Skipping unsupported test with multiple account")
		t.Skipped()
	}
}

// ---------------------------------------------------------------------------
// Sweeper helpers (mirrors alicloud/alicloud_sweeper_test.go and the
// testSweepPreCheckWithRegions helper from alicloud/provider_test.go).
// ---------------------------------------------------------------------------

var endpoints sync.Map
var signVersion sync.Map
var regionalEndpoints sync.Map

// SDK endpoint caches are keyed by product. Sharing one across regions can route
// a Hangzhou request to a previously cached Beijing endpoint.
func sharedEndpointCacheForRegion(region string) *sync.Map {
	cache, _ := regionalEndpoints.LoadOrStore(region, new(sync.Map))
	return cache.(*sync.Map)
}

// sharedClientForRegion returns a common AlicloudClient setup needed for the sweeper
// functions for a given region
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
		Endpoints:   sharedEndpointCacheForRegion(region),
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

// Skip automatically the sweep testcases which does not support some known regions.
// If supported is true, the regions should a list of supporting the service regions.
// If supported is false, the regions should a list of unsupporting the service regions.
func testSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}

// ---------------------------------------------------------------------------
// DataSource runner. acctest's DataSourceTestCheck(WithPreCheck) hardcode
// SDKProviderFactories + acctest.PreCheck, but the VPC suite's dataSourceTestCheck
// uses this package's testAccProviderFactory (= acctest.Provider) and
// testAccPreCheck (which writes this package's defaultRegionToTest that the VPC
// fixtures read at entry time). Routing the runner to acctest's exported form
// would detach that region write, so the dataSourceAttr / dataSourceTestAccConfig
// data types, buildDataSourceSteps step builder and the dataSourceTestCheck /
// dataSourceTestCheckWithPreCheck runners stay local. Only checkDataSourceAttr
// delegates the exist/empty attribute-check construction to acctest's exported
// ResourceAttrInit + ResourceAttrMapUpdateSet.
// ---------------------------------------------------------------------------

// dataSourceAttr has 3 field ,incloud resourceId  existMapFunc fakeMapFunc, every dataSource test can use only one
type dataSourceAttr struct {
	// IDRefreshName, like "data.alicloud_dns_records.record"
	resourceId string

	// get existMap function
	existMapFunc func(rand int) map[string]string

	// get fakeMap function
	fakeMapFunc func(rand int) map[string]string
}

// get exist and empty resourceAttrMapUpdate function
func (dsa *dataSourceAttr) checkDataSourceAttr(rand int) (exist, empty acctest.ResourceAttrMapUpdate) {
	exist = acctest.ResourceAttrInit(dsa.resourceId, dsa.existMapFunc(rand)).ResourceAttrMapUpdateSet()
	empty = acctest.ResourceAttrInit(dsa.resourceId, dsa.fakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	return
}

// according to configs generate step list and execute the test
// lintignore: AT001
func (dsa *dataSourceAttr) dataSourceTestCheck(t *testing.T, rand int, configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps:             steps,
	})
}

// according to configs generate step list and execute the test with preCheck
// lintignore: AT001
func (dsa *dataSourceAttr) dataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          preCheck,
		ProviderFactories: testAccProviderFactory,
		Steps:             steps,
	})
}

// per schema attribute test config
type dataSourceTestAccConfig struct {
	// be equal to testCase config string,but the result has only one record
	existConfig string

	// if the dataSourceAttr.existMapFunc returned map value not match we want, existChangMap can alter checkMap for existConfig
	existChangMap map[string]string

	// be equal to testCase config string,but the result is empty
	fakeConfig string

	// if the dataSourceAttr.fakeMapFunc returned map value not match we want, fakeChangMap can alter checkMap for fakeConfig
	fakeChangMap map[string]string
}

// build test cases for each attribute
func (conf *dataSourceTestAccConfig) buildDataSourceSteps(t *testing.T, info *dataSourceAttr, rand int) []resource.TestStep {
	testAccCheckExist, testAccCheckEmpty := info.checkDataSourceAttr(rand)
	var steps []resource.TestStep
	if conf.existConfig != "" {
		step := resource.TestStep{
			Config: conf.existConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckExist(conf.existChangMap),
			),
		}
		steps = append(steps, step)
	}
	if conf.fakeConfig != "" {
		step := resource.TestStep{
			Config: conf.fakeConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckEmpty(conf.fakeChangMap),
			),
		}
		steps = append(steps, step)
	}
	return steps
}
