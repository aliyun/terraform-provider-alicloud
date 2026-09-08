package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall Instance. >>> Resource test cases, automatically generated.
// Case 国内版按量付费2.0 11709
func TestAccAliCloudCloudFirewallInstanceV2_basic11709(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":          "PayAsYouGo",
					"product_code":          "cfw",
					"product_type":          "cfw_elasticity_public_cn",
					"spec":                  "payg_version",
					"auto_asset_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "PayAsYouGo",
						"product_code":          "cfw",
						"product_type":          "cfw_elasticity_public_cn",
						"spec":                  "payg_version",
						"auto_asset_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sdl":                   "true",
					"auto_asset_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sdl":                   "true",
						"auto_asset_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":     "false",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11709_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
					"sdl":          "true",
					"cfw_log":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
						"spec":         "payg_version",
						"sdl":          "true",
						"cfw_log":      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number"},
			},
		},
	})
}

var AliCloudCloudFirewallInstanceV2Map11709 = map[string]string{
	"cfw_log":     CHECKSET,
	"create_time": CHECKSET,
	"end_time":    CHECKSET,
	"user_status": CHECKSET,
	"status":      CHECKSET,
}

func AliCloudCloudFirewallInstanceV2BasicDependence11709(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 国际版按量付费2.0 11710
func TestAccAliCloudCloudFirewallInstanceV2_basic11710(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_intl",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_intl",
						"spec":         "payg_version",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sdl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sdl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":     "false",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11710_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_intl",
					"spec":         "payg_version",
					"sdl":          "true",
					"cfw_log":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_intl",
						"spec":         "payg_version",
						"sdl":          "true",
						"cfw_log":      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number"},
			},
		},
	})
}

// Case 国内版预付费2.0 11711
func TestAccAliCloudCloudFirewallInstanceV2_basic11711(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11711)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"product_code": "cfw",
					"product_type": "cfw_sub_public_cn",
					"spec":         "premium_version",
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"product_code": "cfw",
						"product_type": "cfw_sub_public_cn",
						"spec":         "premium_version",
						"period":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11711_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11711)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":          "Subscription",
					"product_code":          "cfw",
					"product_type":          "cfw_sub_public_cn",
					"spec":                  "enterprise_version",
					"sdl":                   "true",
					"cfw_log":               "false",
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
					"period":                "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "Subscription",
						"product_code":          "cfw",
						"product_type":          "cfw_sub_public_cn",
						"spec":                  "enterprise_version",
						"sdl":                   "true",
						"cfw_log":               "false",
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number"},
			},
		},
	})
}

var AliCloudCloudFirewallInstanceV2Map11711 = map[string]string{
	"cfw_log":               CHECKSET,
	"renewal_duration_unit": CHECKSET,
	"renewal_status":        CHECKSET,
	"create_time":           CHECKSET,
	"end_time":              CHECKSET,
	"user_status":           CHECKSET,
	"status":                CHECKSET,
}

// Case 国际版预付费2.0 11712
func TestAccAliCloudCloudFirewallInstanceV2_basic11712(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11711)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"product_code": "cfw",
					"product_type": "cfw_sub_public_intl",
					"spec":         "premium_version",
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"product_code": "cfw",
						"product_type": "cfw_sub_public_intl",
						"spec":         "premium_version",
						"period":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11712_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11711)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":          "Subscription",
					"product_code":          "cfw",
					"product_type":          "cfw_sub_public_intl",
					"spec":                  "ultimate_version",
					"sdl":                   "true",
					"cfw_log":               "false",
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
					"period":                "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "Subscription",
						"product_code":          "cfw",
						"product_type":          "cfw_sub_public_intl",
						"spec":                  "ultimate_version",
						"sdl":                   "true",
						"cfw_log":               "false",
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number"},
			},
		},
	})
}

// cfw_log_storage and fw_vpc_number are intentionally absent from this map.
// DescribeUserBuyVersion reports LogStorage/VpcNumber as 0 for Subscription
// instances even after the instance is normal, and the Read function
// deliberately skips writing a 0 value back to state (it would otherwise
// clear the user-configured count and produce a spurious diff). On step 0
// the test config does not set these attributes, so the state legitimately
// has no value for them, and a CHECKSET would always fail. They are still
// validated by the explicit testAccCheck overrides in the steps that
// actually set them (step 2 sets cfw_log_storage=5000, step 1/4 set
// fw_vpc_number=5/6).
var AliCloudCloudFirewallInstanceV2Map11713 = map[string]string{
	"cfw_log":               CHECKSET,
	"ip_number":             CHECKSET,
	"auto_asset_protection": CHECKSET,
	"renewal_duration_unit": CHECKSET,
	"renewal_status":        CHECKSET,
	"create_time":           CHECKSET,
	"end_time":              CHECKSET,
	"user_status":           CHECKSET,
	"status":                CHECKSET,
}

// Case 国内版预付费2.0 11713
func TestAccAliCloudCloudFirewallInstanceV2_basic11713(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11713)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			// testAccPreCheckWithTime(t, []int{1}) // monthly gate disabled to let remote ACC exercise this new case
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"product_code": "cfw",
					"product_type": "cfw_sub_public_cn",
					"spec":         "enterprise_version",
					"ip_number":    "50",
					"band_width":   "50",
					"cfw_log":      "false",
					"logistics":    cloudFirewallInstanceLogisticsConfig,
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"product_code": "cfw",
						"product_type": "cfw_sub_public_cn",
						"spec":         "enterprise_version",
						"ip_number":    "50",
						"cfw_log":      "false",
						"logistics":    cloudFirewallInstanceLogisticsState,
						"period":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec":          "ultimate_version",
					"ip_number":     "400",
					"band_width":    "200",
					"fw_vpc_number": "5",
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":          "ultimate_version",
						"ip_number":     "400",
						"fw_vpc_number": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":         "true",
					"cfw_log_storage": "5000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log":         "true",
						"cfw_log_storage": "5000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_number": "405",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_number": "405",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fw_vpc_number": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fw_vpc_number": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// ip_number, cfw_log_storage and fw_vpc_number are ignored on Subscription:
				// DescribeUserBuyVersion reports these counts as 0 and the Read guard skips
				// writing them, so the imported state defaults to 0 while the pre-import state
				// held the user-configured values.
				ImportStateVerifyIgnore: []string{"modify_type", "period", "band_width", "logistics", "instance_count", "cfw_account", "account_number", "ip_number", "cfw_log_storage", "fw_vpc_number"},
			},
		},
	})
}

// Case unsupported order parameters 11714
func TestAccAliCloudCloudFirewallInstanceV2_unsupportedOrderParameters(t *testing.T) {
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			// testAccPreCheckWithTime(t, []int{1}) // monthly gate disabled to let remote ACC exercise this new case
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"product_code":   "cfw",
					"product_type":   "cfw_sub_public_cn",
					"spec":           "enterprise_version",
					"ip_number":      "50",
					"band_width":     "50",
					"cfw_log":        "false",
					"cfw_account":    "true",
					"account_number": "1",
					"instance_count": "5",
					"logistics":      cloudFirewallInstanceLogisticsConfig,
					"period":         "1",
				}),
				ExpectError: regexp.MustCompile("InvalidParameter"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_account":    "false",
					"account_number": "2",
					"instance_count": "10",
				}),
				ExpectError: regexp.MustCompile("InvalidParameter"),
			},
		},
	})
}

// Test CloudFirewall Instance. <<< Resource test cases, automatically generated.
