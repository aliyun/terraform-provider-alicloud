package alicloud

import (
	"fmt"
	"testing"

	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudWafInstance_basic(t *testing.T) {
	var v waf_openapi.DescribeInstanceInfoResponse
	resourceId := "alicloud_waf_instance.default"
	ra := resourceAttrInit(resourceId, WafInstanceMap)
	serviceFunc := func() interface{} {
		return &Waf_openapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeWafInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccWafInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, WafInstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"big_screen":           "0",
					"exclusive_ip_package": "1",
					"ext_bandwidth":        "50",
					"ext_domain_package":   "1",
					"package_code":         "version_3",
					"prefessional_service": "false",
					"subscription_type":    "Subscription",
					"period":               "1",
					"waf_log":              "false",
					"log_storage":          "3",
					"log_time":             "180",
					"resource_group_id":    `${"data.alicloud_resource_manager_resource_groups.this.groups.0.id"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"big_screen":           "0",
						"exclusive_ip_package": "1",
						"ext_bandwidth":        "50",
						"ext_domain_package":   "1",
						"package_code":         "version_3",
						"prefessional_service": "false",
						"subscription_type":    "Subscription",
						"period":               "1",
						"waf_log":              "false",
						"log_storage":          "3",
						"log_time":             "180",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"big_screen", "exclusive_ip_package", "ext_bandwidth", "ext_domain_package", "log_storage", "log_time", "modify_type", "package_code", "period", "prefessional_service", "renew_period", "renewal_status", "resource_group_id", "waf_log"},
			},
			//  modify_type does not support updating alone
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"modify_type": "Upgrade",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"modify_type": "Upgrade",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"big_screen": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"big_screen": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclusive_ip_package": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclusive_ip_package": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ext_bandwidth": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ext_bandwidth": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ext_domain_package": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ext_domain_package": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_code": "version_4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_code": "version_4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prefessional_service": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefessional_service": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"waf_log":     "true",
					"log_storage": "3",
					"log_time":    "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"waf_log":     "true",
						"log_storage": "3",
						"log_time":    "180",
					}),
				),
			},
		},
	})
}

var WafInstanceMap = map[string]string{
	"status": CHECKSET,
}

func WafInstanceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_resource_manager_resource_groups" "this" {
  name_regex = "^default$"
}
`, name)
}
