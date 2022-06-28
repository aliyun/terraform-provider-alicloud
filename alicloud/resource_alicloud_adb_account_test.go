package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudADBAccount_update_forSuper(t *testing.T) {
	var v *adb.DBAccount
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccadbaccount-%d", rand)
	var basicMap = map[string]string{
		"db_cluster_id":    CHECKSET,
		"account_name":     "tftestsuper",
		"account_password": "YourPassword_123",
	}
	resourceId := "alicloud_adb_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAdbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbAccountConfigDependence)
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
					"db_cluster_id":    "${alicloud_adb_db_cluster.cluster.id}",
					"account_name":     "tftestsuper",
					"account_password": "YourPassword_123",
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

func resourceAdbAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	resource "alicloud_adb_db_cluster" "cluster" {
		db_cluster_category = "MixedStorage"
		mode = "flexible"
		compute_resource = "8Core32GB"
		vswitch_id              = local.vswitch_id
		description             = "${var.name}"
	}`, AdbCommonTestCase, name)
}
