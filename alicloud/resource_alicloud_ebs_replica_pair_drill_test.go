package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ebs ReplicaPairDrill. >>> Resource test cases, automatically generated.
// Case 5513
func TestAccAliCloudEbsReplicaPairDrill_basic5513(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_replica_pair_drill.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsReplicaPairDrillMap5513)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsReplicaPairDrill")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsreplicapairdrill%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsReplicaPairDrillBasicDependence5513)
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
					"pair_id": "pair-cn-wwo3kjfq5001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pair_id": CHECKSET,
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

var AlicloudEbsReplicaPairDrillMap5513 = map[string]string{
	"status":                CHECKSET,
	"replica_pair_drill_id": CHECKSET,
}

func AlicloudEbsReplicaPairDrillBasicDependence5513(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region" {
  default = "%s"
}
`, name, defaultRegionToTest)
}

// Test Ebs ReplicaPairDrill. <<< Resource test cases, automatically generated.
