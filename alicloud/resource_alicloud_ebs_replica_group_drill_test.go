package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ebs ReplicaGroupDrill. >>> Resource test cases, automatically generated.
// Case 5254
func TestAccAliCloudEbsReplicaGroupDrill_basic5254(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_replica_group_drill.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsReplicaGroupDrillMap5254)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsReplicaGroupDrill")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsreplicagroupdrill%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsReplicaGroupDrillBasicDependence5254)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "pg-m1H9aaOUIGsDUwgZ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": CHECKSET,
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

var AlicloudEbsReplicaGroupDrillMap5254 = map[string]string{
	"status":                 CHECKSET,
	"replica_group_drill_id": CHECKSET,
}

func AlicloudEbsReplicaGroupDrillBasicDependence5254(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region" {
  default = "%s"
}

`, name, defaultRegionToTest)
}

// Test Ebs ReplicaGroupDrill. <<< Resource test cases, automatically generated.
