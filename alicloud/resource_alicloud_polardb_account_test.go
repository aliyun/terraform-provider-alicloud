package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPolarDBAccount_update(t *testing.T) {
	var v *polardb.DBAccount
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbaccount-%d", rand)
	var basicMap = map[string]string{
		"db_cluster_id":    CHECKSET,
		"account_name":     "tftestnormal",
		"account_password": "YourPassword_123",
		"account_type":     "Normal",
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
					"db_cluster_id":    "${alicloud_polardb_cluster.cluster.id}",
					"account_name":     "tftestnormal",
					"account_password": "YourPassword_123",
					"account_type":     "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "from terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_1234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_1234",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "tf test",
					"account_password":    "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "tf test",
						"account_password":    "YourPassword_123",
					}),
				),
			},
		},
	})

}
func TestAccAlicloudPolarDBAccount_update_forSuper(t *testing.T) {
	var v *polardb.DBAccount
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbaccount-%d", rand)
	var basicMap = map[string]string{
		"db_cluster_id":    CHECKSET,
		"account_name":     "tftestsuper",
		"account_password": "YourPassword_123",
		"account_type":     "Super",
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
					"db_cluster_id":    "${alicloud_polardb_cluster.cluster.id}",
					"account_name":     "tftestsuper",
					"account_password": "YourPassword_123",
					"account_type":     "Super",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "from terraform super",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "from terraform super",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_12345",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "tf test super",
					"account_password":    "YourPassword_1234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "tf test super",
						"account_password":    "YourPassword_1234",
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

	resource "alicloud_polardb_cluster" "cluster" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id = "${alicloud_vswitch.default.id}"
		description = "${var.name}"
	}`, PolarDBCommonTestCase, name)
}
