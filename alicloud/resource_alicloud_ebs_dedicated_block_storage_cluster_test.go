package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func SkipTestAccAlicloudEbsDedicatedBlockStorageCluster_basic(t *testing.T) {
	var v map[string]interface{}
	testAccPreCheckWithRegions(t, true, connectivity.EbsDedicatedBlockStorageClusterRegions)
	resourceId := "alicloud_ebs_dedicated_block_storage_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsDedicatedBlockStorageClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsDedicatedBlockStorageCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsDedicatedBlockStorageClusterBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":                                 "Premium",
					"zone_id":                              "${data.alicloud_ebs_regions.default.regions[0].zones[0].zone_id}",
					"dedicated_block_storage_cluster_name": name,
					"total_capacity":                       "61440",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                                 "Premium",
						"zone_id":                              CHECKSET,
						"dedicated_block_storage_cluster_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_block_storage_cluster_name": name + "_update",
					"description":                          name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_block_storage_cluster_name": name + "_update",
						"description":                          name + "_update",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEbsDedicatedBlockStorageClusterMap = map[string]string{}

func AlicloudEbsDedicatedBlockStorageClusterBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region" {
  default = "%s"
}

data "alicloud_ebs_regions" "default"{
  region_id = var.region
}

`, name, defaultRegionToTest)
}
