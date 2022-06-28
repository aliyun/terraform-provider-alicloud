package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSDDPInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sddp_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudSDDPInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SddpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSddpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSDDPInstanceBasicDependence0)
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
					"payment_type": "Subscription",
					"sddp_version": "version_company",
					"sdc":          "3",
					"sd_cbool":     "no",
					"period":       "1",
					"ud_cbool":     "yes",
					"udc":          "2000",
					"dataphin":     "yes",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"sddp_version": "version_company",
						"sd_cbool":     "no",
						"sdc":          "3",
						"ud_cbool":     "yes",
						"udc":          "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type": "Upgrade",
					"sdc":         "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sdc": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type": "Upgrade",
					"sd_cbool":    "no",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sd_cbool": "no",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type": "Upgrade",
					"udc":         "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"udc": "3000",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sd_cbool", "period", "ud_cbool", "udc", "dataphin", "modify_type"},
			},
		},
	})
}

var AlicloudSDDPInstanceMap0 = map[string]string{}

func AlicloudSDDPInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
