package alicloud

import (
	"fmt"
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
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
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
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
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
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
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
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
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
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
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

// Test CloudFirewall Instance. <<< Resource test cases, automatically generated.
