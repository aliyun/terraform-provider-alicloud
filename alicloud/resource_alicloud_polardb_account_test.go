package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPolarDBAccount_update(t *testing.T) {
	var v *polardb.DBAccount
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbaccount-%d", rand)
	var basicMap = map[string]string{
		"cluster_id": CHECKSET,
		"name":       "tftestnormal",
		"password":   "YourPassword_123",
		"type":       string(PolarDBAccountNormal),
	}
	resourceId := "alicloud_polardb_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBAccountConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${alicloud_polardb_cluster.cluster.id}",
					"name":       "tftestnormal",
					"password":   "YourPassword_123",
					"type":       string(PolarDBAccountNormal),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "from terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YourPassword_1234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "YourPassword_1234",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf test",
					"password":    "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf test",
						"password":    "YourPassword_123",
					}),
				),
			},
		},
	})

}

func resourcePolarDBAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

	variable "instancechargetype" {
		default = "PostPaid"
	}

	variable "engine" {
		default = "MySQL"
	}

	variable "engineversion" {
		default = "8.0"
	}

	variable "instanceclass" {
		default = "polar.mysql.x4.large"
	}

	resource "alicloud_polardb_cluster" "cluster" {
		db_type = "${var.engine}"
		db_version = "${var.engineversion}"
		pay_type = "${var.instancechargetype}"
		db_node_class = "${var.instanceclass}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		description = "${var.name}"
	}`, PolarDBCommonTestCase, name)
}
