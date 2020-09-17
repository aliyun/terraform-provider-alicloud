package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGpdbConnectionUpdate(t *testing.T) {
	var v *gpdb.DBInstanceNetInfo

	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": REGEXMATCH + fmt.Sprintf("^tf-testacc%s.gpdb.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com", rand),
		"port":              "3306",
	}

	resourceId := "alicloud_gpdb_connection.default"
	serverFunc := func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeGpdbConnection")
	ra := resourceAttrInit(resourceId, basicMap)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", testGpdbConnectionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_gpdb_instance.default.id}",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_gpdb_instance.default.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%s", rand),
					"port":              "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
		},
	})
}

func testGpdbConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`
        data "alicloud_zones" "default" {
            available_resource_creation = "Gpdb"
        }
        variable "name" {
            default              = "tf-testAccGpdbInstance"
        }
        data "alicloud_vswitches" "default" {
		  zone_id    = data.alicloud_zones.default.ids[0]
		  name_regex = "default-tf--testAcc-00"
		}
        resource "alicloud_gpdb_instance" "default" {
            vswitch_id           = "${data.alicloud_vswitches.default.ids.0}"
            engine               = "gpdb"
            engine_version       = "4.3"
            instance_class       = "gpdb.group.segsdx2"
            instance_group_count = "2"
            description          = "${var.name}"
	}`)
}
