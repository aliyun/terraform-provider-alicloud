package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       alicloud_gpdb_instance.default.id,
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
					"instance_id":       alicloud_gpdb_instance.default.id,
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
        resource "alicloud_vpc" "default" {
            name                 = var.name
            cidr_block           = "172.16.0.0/16"
        }
        resource "alicloud_vswitch" "default" {
            availability_zone    = data.alicloud_zones.default.zones.0.id
            vpc_id               = alicloud_vpc.default.id
            cidr_block           = "172.16.0.0/24"
            name                 = var.name
        }
        resource "alicloud_gpdb_instance" "default" {
            vswitch_id           = alicloud_vswitch.default.id
            engine               = "gpdb"
            engine_version       = "4.3"
            instance_class       = "gpdb.group.segsdx2"
            instance_group_count = "2"
            description          = var.name
	}`)
}
