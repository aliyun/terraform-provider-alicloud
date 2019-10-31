package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCloudConnectNetworkGrant_sameAccount(t *testing.T) {
	var grantRule smartag.GrantRule
	resourceId := "alicloud_cloud_connect_network_grant.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &grantRule, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccCloudConnectNetworkGrant"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCcnGrantRuleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheckWithCloudConnectNetworkGrantSetting(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ccn_id":  "${alicloud_cloud_connect_network.default.id}",
					"cen_id":  "${alicloud_cen_instance.default.id}",
					"cen_uid": "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ccn_id":  CHECKSET,
						"cen_id":  CHECKSET,
						"cen_uid": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"cen_uid"},
			},
		},
	})
}

func TestAccAlicloudCloudConnectNetworkGrant_differentAccount(t *testing.T) {
	var grantRule smartag.GrantRule
	resourceId := "alicloud_cloud_connect_network_grant.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &grantRule, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccCloudConnectNetworkGrantCen"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCcnGrantRuleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheckWithCloudConnectNetworkGrantSetting(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ccn_id":  "${alicloud_cloud_connect_network.default.id}",
					"cen_id":  os.Getenv("GRANT_CEN_ID"),
					"cen_uid": os.Getenv("GRANT_CEN_UID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ccn_id":  CHECKSET,
						"cen_id":  os.Getenv("GRANT_CEN_ID"),
						"cen_uid": os.Getenv("GRANT_CEN_UID"),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"cen_uid"},
			},
		},
	})
}

func resourceCcnGrantRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default"{
}

resource "alicloud_cloud_connect_network" "default" {
  name = "${var.name}"
  is_default = "true"
}

resource "alicloud_cen_instance" "default" {
  name = "${var.name}"
}
`, name)
}
