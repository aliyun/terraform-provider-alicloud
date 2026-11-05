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

// Test Ebs ReplicaPairDrill. >>> Resource test cases, automatically generated.
// Case 5513
func TestAccAliCloudEbsReplicaPairDrill_basic5513(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_replica_pair_drill.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEbsReplicaPairDrillMap5513)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ebs.NewEbsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEbsReplicaPairDrill")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsreplicapairdrill%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEbsReplicaPairDrillBasicDependence5513)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"pair_id": "pair-cn-wwo3kjfq5001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pair_id": acctest.CHECKSET,
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
	"status":                acctest.CHECKSET,
	"replica_pair_drill_id": acctest.CHECKSET,
}

func AlicloudEbsReplicaPairDrillBasicDependence5513(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region" {
  default = "%s"
}
`, name, acctest.DefaultRegionToTest)
}

// Test Ebs ReplicaPairDrill. <<< Resource test cases, automatically generated.
