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

// Test Ebs ReplicaGroupDrill. >>> Resource test cases, automatically generated.
// Case 5254
func TestAccAliCloudEbsReplicaGroupDrill_basic5254(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_replica_group_drill.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEbsReplicaGroupDrillMap5254)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ebs.NewEbsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEbsReplicaGroupDrill")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsreplicagroupdrill%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEbsReplicaGroupDrillBasicDependence5254)
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
					"group_id": "pg-m1H9aaOUIGsDUwgZ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": acctest.CHECKSET,
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
	"status":                 acctest.CHECKSET,
	"replica_group_drill_id": acctest.CHECKSET,
}

func AlicloudEbsReplicaGroupDrillBasicDependence5254(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region" {
  default = "%s"
}

`, name, acctest.DefaultRegionToTest)
}

// Test Ebs ReplicaGroupDrill. <<< Resource test cases, automatically generated.
