package sddp_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sddp"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSDDPInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sddp_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSDDPInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sddp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSddpInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSDDPInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithTime(t, []int{30})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dataphin_count": "CreateInstance",
					"logistics":      "yes",
					"oss_size":       "3000",
					"renew_period":   "1",
					"renewal_status": "CreateInstance",
					"payment_type":   "Subscription",
					"sddp_version":   "version_company",
					"sdc":            "3",
					"sd_cbool":       "no",
					"period":         "1",
					"ud_cbool":       "yes",
					"udc":            "2000",
					"dataphin":       "yes",
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
					"dataphin_count": "ModifyInstance",
					"logistics":      "no",
					"oss_size":       "3001",
					"renew_period":   "2",
					"renewal_status": "AutoRenewal",
					"sddp_version":   "version_audit",
					"modify_type":    "Upgrade",
					"sdc":            "8",
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
