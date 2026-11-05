package nas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNASDataFlow_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_data_flow.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASCPFSSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNASDataFlowMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasDataFlow")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasdataflow%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNASDataFlowBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fset_id":        "${alicloud_nas_fileset.default.fileset_id}",
					"throughput":     "600",
					"source_storage": "oss://${alicloud_oss_bucket.default.bucket}",
					"file_system_id": "${alicloud_nas_file_system.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fset_id":        acctest.CHECKSET,
						"throughput":     "600",
						"source_storage": acctest.CHECKSET,
						"file_system_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"throughput": "1200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"throughput": "1200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}
func TestAccAliCloudNASDataFlow_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_data_flow.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASCPFSSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNASDataFlowMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasDataFlow")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasdataflow%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNASDataFlowBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fset_id":              "${alicloud_nas_fileset.default.fileset_id}",
					"description":          "${var.name}",
					"throughput":           "600",
					"source_storage":       "oss://${alicloud_oss_bucket.default.bucket}",
					"source_security_type": "SSL",
					"file_system_id":       "${alicloud_nas_file_system.default.id}",
					"dry_run":              "false",
					"status":               "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fset_id":              acctest.CHECKSET,
						"description":          name,
						"throughput":           "600",
						"source_storage":       acctest.CHECKSET,
						"source_security_type": "SSL",
						"file_system_id":       acctest.CHECKSET,
						"status":               "Running",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudNASDataFlowMap0 = map[string]string{
	"status":         acctest.CHECKSET,
	"file_system_id": acctest.CHECKSET,
	"data_flow_id":   acctest.CHECKSET,
}

func AlicloudNASDataFlowBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "cpfs"
  storage_type     = "advance_100"
  file_system_type = "cpfs"
  capacity         = 3600
  description      = "tf-testacc"
  zone_id          = local.zone_id
  vpc_id           = data.alicloud_vpcs.default.ids.0
  vswitch_id       = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_nas_mount_target" "default" {
	file_system_id = "${alicloud_nas_file_system.default.id}"
	vswitch_id = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "private"
  tags   = {
    cpfs-dataflow = "true"
  }
}

resource "alicloud_nas_fileset" "default" {
  depends_on       = ["alicloud_nas_mount_target.default"]
  file_system_id   = alicloud_nas_file_system.default.id
  description      = var.name
  file_system_path = "/tf-testAcc-Path/"
}
`, name)
}

// lintignore: R001
