package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudFirewallInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BssOpenApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "QueryAvailableInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{30})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":    "Subscription",
					"spec":            "premium_version",
					"ip_number":       "20",
					"band_width":      "10",
					"cfw_log":         "false",
					"cfw_log_storage": "1000",
					"cfw_service":     "false",
					"period":          "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":    "Subscription",
						"spec":            "premium_version",
						"ip_number":       "20",
						"band_width":      "10",
						"cfw_log":         "false",
						"cfw_log_storage": "1000",
						"cfw_service":     "false",
						"period":          "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_service": "true",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_service": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fw_vpc_number": "3",
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fw_vpc_number": "3",
						"modify_type": "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"band_width":  "20",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"band_width":  "20",
						"modify_type": "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log_storage": "1200",
					"modify_type":     "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log_storage": "1200",
						"modify_type":     "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":     "true",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log":     "true",
						"modify_type": "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_period": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_period": "6",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"band_width", "cfw_log", "cfw_log_storage", "cfw_service", "ip_number", "payment_type", "period", "modify_type", "spec"},
			},
		},
	})
}

var AlicloudCloudFirewallInstanceMap0 = map[string]string{}

func AlicloudCloudFirewallInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
