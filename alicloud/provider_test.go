package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/features"
	ecssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	fcsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/fc"
	osssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Marked setup copies remain local because importing alicloud/acctest would create an import cycle.

// Migrated to alicloud/acctest.SDKProviders.
var testAccProviders map[string]*schema.Provider

// Migrated to alicloud/acctest.SDKProviderFactories.
var testAccProviderFactory map[string]func() (*schema.Provider, error)

// Migrated to alicloud/acctest.Provider.
var testAccProvider *schema.Provider

// Migrated to alicloud/acctest.DefaultRegionToTest.
var defaultRegionToTest = func() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	// Fallback to cn-beijing when ALICLOUD_REGION is unset so test fixtures
	// that read defaultRegionToTest at test-function-entry time (before
	// testAccPreCheck runs) don't produce malformed strings like
	// "oss--internal.aliyuncs.com" or "...fcv3.<acct>..fc.devsapp.net".
	return "cn-beijing"
}()

// Migrated to alicloud/acctest/setup.go (same-named constants).
const (
	ProviderNameAlicloudAlt  = "alicloudalt"
	ProviderNameAlicloudAlt2 = "alicloudalt2"
	ProviderNameAlicloudAlt3 = "alicloudalt3"
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"alicloud": testAccProvider,
	}

	testAccProviderFactory = map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
	// TODO disable setting STS from oss
	//setStsCredential()
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProviderSkipRegionValidation(t *testing.T) {
	original := os.Getenv("ALICLOUD_SKIP_REGION_VALIDATION")
	defer os.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", original)

	testCases := []struct {
		name     string
		env      string
		raw      map[string]interface{}
		expected bool
	}{
		{"unset env defaults to false", "", map[string]interface{}{}, false},
		{"env true enables skip", "true", map[string]interface{}{}, true},
		{"env false disables skip", "false", map[string]interface{}{}, false},
		{"config overrides env", "true", map[string]interface{}{"skip_region_validation": false}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			os.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", testCase.env)

			d := schema.TestResourceDataRaw(t, Provider().Schema, testCase.raw)
			if got := d.Get("skip_region_validation").(bool); got != testCase.expected {
				t.Fatalf("skip_region_validation: expected %v, got %v", testCase.expected, got)
			}
		})
	}
}

// TestProviderFeatures covers expandFeatures over the shapes the features block can take: every
// level of it is optional, so a missing block or field has to fall back to features.Default.
func TestProviderFeatures(t *testing.T) {
	testCases := []struct {
		name     string
		features []interface{}
		expected features.Features
	}{
		{
			name:     "no features block at all",
			features: []interface{}{},
			expected: features.Default(),
		},
		{
			name:     "a features block that is null",
			features: []interface{}{nil},
			expected: features.Default(),
		},
		{
			name:     "a features block without the ecs_instance block",
			features: []interface{}{map[string]interface{}{}},
			expected: features.Default(),
		},
		{
			name:     "an empty ecs_instance block",
			features: []interface{}{map[string]interface{}{"ecs_instance": []interface{}{}}},
			expected: features.Default(),
		},
		{
			name: "replace_on_image_update enabled",
			features: []interface{}{map[string]interface{}{
				"ecs_instance": []interface{}{map[string]interface{}{"replace_on_image_update": true}},
			}},
			expected: features.Features{
				EcsInstance: features.EcsInstance{ReplaceOnImageUpdate: true},
			},
		},
		{
			name: "replace_on_image_update disabled explicitly",
			features: []interface{}{map[string]interface{}{
				"ecs_instance": []interface{}{map[string]interface{}{"replace_on_image_update": false}},
			}},
			expected: features.Default(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := expandFeatures(testCase.features); got != testCase.expected {
				t.Fatalf("expandFeatures: expected %+v, got %+v", testCase.expected, got)
			}
		})
	}
}

// Migrated to alicloud/acctest.PreCheck.
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

// Test_testAccPreChecksSyncDefaultRegion verifies that prechecks keep the
// package-level default region in sync with the final
// ALICLOUD_REGION value. It does not contact any cloud API.
func Test_testAccPreChecksSyncDefaultRegion(t *testing.T) {
	testCases := []struct {
		name           string
		initialRegion  string
		checkoutRegion string
		wantRegion     string
		preCheck       func(*testing.T)
	}{
		{
			name:          "default precheck fallback",
			wantRegion:    string(connectivity.Beijing),
			initialRegion: "",
			preCheck:      testAccPreCheck,
		},
		{
			name:          "regional precheck fallback",
			wantRegion:    string(connectivity.Beijing),
			initialRegion: "",
			preCheck: func(t *testing.T) {
				testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
			},
		},
		{
			name:           "regional precheck checkout override",
			initialRegion:  string(connectivity.Hangzhou),
			checkoutRegion: "true",
			wantRegion:     string(connectivity.EUCentral1),
			preCheck: func(t *testing.T) {
				testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.EUCentral1})
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("ALICLOUD_ACCESS_KEY", "test-access-key")
			t.Setenv("ALICLOUD_SECRET_KEY", "test-secret-key")
			t.Setenv("CHECKOUT_REGION", testCase.checkoutRegion)
			t.Setenv("ALICLOUD_REGION", testCase.initialRegion)
			if testCase.initialRegion == "" {
				if err := os.Unsetenv("ALICLOUD_REGION"); err != nil {
					t.Fatalf("unset ALICLOUD_REGION: %s", err)
				}
			}

			originalDefaultRegion := defaultRegionToTest
			t.Cleanup(func() { defaultRegionToTest = originalDefaultRegion })
			defaultRegionToTest = testCase.initialRegion

			testCase.preCheck(t)

			if got := defaultRegionToTest; got != testCase.wantRegion {
				t.Fatalf("defaultRegionToTest = %q, want %q", got, testCase.wantRegion)
			}
			if got := os.Getenv("ALICLOUD_REGION"); got != testCase.wantRegion {
				t.Fatalf("ALICLOUD_REGION = %q, want %q", got, testCase.wantRegion)
			}
		})
	}
}

// currently not all account site type support create PostPaid resources, PayByBandwidth and other limits.
// The setting of account site type can skip some unsupported cases automatically.

// Skip automatically the testcases which does not support some known regions.
// If supported is true, the regions should a list of supporting the service regions.
// If supported is false, the regions should a list of unsupporting the service regions.
// If the region is unsupported and has backend region, the backend region will instead
// Migrated to alicloud/acctest.PreCheckWithRegions.
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

// Migrated to alicloud/acctest.CheckoutSupportedRegions.
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

type DingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

var providerCommon = `
provider "alicloud" {
	assume_role {}
}
`

func TestAccAlicloudProviderEcs(t *testing.T) {
	var v ecs.Instance

	resourceId := "alicloud_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return ecssvc.NewService(testAccProvider.Meta().(*connectivity.AliyunClient))
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceInstanceVpcConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":        "${data.alicloud_images.default.images.0.id}",
					"security_groups": []string{"${alicloud_security_group.default.0.id}"},
					"instance_type":   "${data.alicloud_instance_types.default.instance_types.0.id}",

					"availability_zone":             "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category":          "cloud_efficiency",
					"instance_name":                 "${var.name}",
					"key_name":                      "${alicloud_key_pair.default.key_name}",
					"spot_strategy":                 "NoSpot",
					"spot_price_limit":              "0",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${alicloud_vswitch.default.id}",
					"role_name":  "${alicloud_ram_role.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudProviderFC(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%salicloudfcfunction-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"service":     CHECKSET,
		"name":        name,
		"runtime":     "python3.9",
		"description": "tf",
		"handler":     "hello.handler",
		"oss_bucket":  CHECKSET,
		"oss_key":     CHECKSET,
	}
	resourceId := "alicloud_fc_function.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return fcsvc.NewService(testAccProvider.Meta().(*connectivity.AliyunClient))
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceFCFunctionConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service":     "${alicloud_fc_service.default.name}",
					"name":        "${var.name}",
					"runtime":     "python3.9",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  "${alicloud_oss_bucket.default.id}",
					"oss_key":     "${alicloud_oss_bucket_object.default.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudProviderOss(t *testing.T) {
	var v oss.GetBucketInfoResult

	resourceId := "alicloud_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return osssvc.NewOssService(testAccProvider.Meta().(*connectivity.AliyunClient))
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sbucket-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceOssBucketConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
		},
	})
}
