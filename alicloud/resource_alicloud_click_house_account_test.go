package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudClickHouseAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_account.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	pwd := fmt.Sprintf("Tf-test%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseAccountBasicDependence0)
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
					"db_cluster_id":    "${alicloud_click_house_db_cluster.default.id}",
					"account_name":     name,
					"account_password": pwd,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":    CHECKSET,
						"account_name":     name,
						"account_password": pwd,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": pwd + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": pwd + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "AccountDescription_all",
					"account_password":    pwd + "updateall",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "AccountDescription_all",
						"account_password":    pwd + "updateall",
					}),
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

var AlicloudClickHouseAccountMap0 = map[string]string{
	"account_name":  CHECKSET,
	"db_cluster_id": CHECKSET,
	"type":          "Normal",
}

func AlicloudClickHouseAccountBasicDependence0(name string) string {
	return fmt.Sprintf(` 
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = "${data.alicloud_vswitches.default.vswitches.0.id}"
}
variable "name" {
  default = "%s"
}
`, name)
}
