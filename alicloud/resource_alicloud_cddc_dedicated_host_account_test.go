package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCddcDedicatedHostAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host_account.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHostAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostAccountBasicDependence0)
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
					"dedicated_host_id": "${data.alicloud_cddc_dedicated_hosts.default.hosts.0.dedicated_host_id}",
					"account_type":      "Normal",
					"account_password":  "Test1234+!",
					"account_name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_id": CHECKSET,
						"account_name":      name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "Test1234+!" + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "Test1234+!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AlicloudCDDCDedicatedHostAccountMap0 = map[string]string{}

func AlicloudCDDCDedicatedHostAccountBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_cddc_dedicated_host_groups" "default" {}

data "alicloud_cddc_dedicated_hosts" "default" {
  dedicated_host_group_id = data.alicloud_cddc_dedicated_host_groups.default.ids.0
}
`, name)
}
