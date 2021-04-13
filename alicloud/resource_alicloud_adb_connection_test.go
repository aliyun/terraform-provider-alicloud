package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAdbConnectionConfig(t *testing.T) {
	var v *adb.Address
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccAdbConnection%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id":     CHECKSET,
		"connection_string": CHECKSET,
		"ip_address":        CHECKSET,
		"port":              CHECKSET,
	}
	resourceId := "alicloud_adb_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAdbConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbConnectionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
			testAccPreCheckWithNoDefaultVswitch(t)
			testAccPreCheckWithRegions(t, false, connectivity.AdbDBClusterUnSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":     "${alicloud_adb_cluster.cluster.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%s", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceAdbConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	resource "alicloud_adb_cluster" "cluster" {
        db_cluster_version      = "3.0"
        db_cluster_category     = "Cluster"
        db_node_class           = data.alicloud_adb_db_cluster_classes.default.available_zone_list[0].classes[0]
        db_node_count           = 2
        db_node_storage         = 200
		pay_type                = "PostPaid"
		vswitch_id              = "${data.alicloud_vswitches.default.ids.0}"
		description             = "${var.name}"
	}`, AdbCommonTestCase, name)
}
