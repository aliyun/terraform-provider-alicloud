package ebs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ebs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
// lintignore: AT001
func SkipTestAccAlicloudEbsDedicatedBlockStorageCluster_basic(t *testing.T) {
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, connectivity.EbsDedicatedBlockStorageClusterRegions)
	resourceId := "alicloud_ebs_dedicated_block_storage_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEbsDedicatedBlockStorageClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ebs.NewEbsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEbsDedicatedBlockStorageCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEbsDedicatedBlockStorageClusterBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
						"zone_id":                              acctest.CHECKSET,
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

`, name, acctest.DefaultRegionToTest)
}
