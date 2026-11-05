package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Ens Snapshot. >>> Resource test cases, automatically generated.
// Case 5162
func TestAccAliCloudEnsSnapshot_basic5162(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_snapshot.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEnsSnapshotMap5162)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsSnapshot")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senssnapshot%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEnsSnapshotBasicDependence5162)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
					"disk_id":       "${alicloud_ens_disk.disk.id}",
					"snapshot_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
						"disk_id":       acctest.CHECKSET,
						"snapshot_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "SnapShotDescription_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SnapShotDescription_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "SnapShotDescription_UPDATE_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SnapShotDescription_UPDATE_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   "SnapShotDescription_autotest",
					"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
					"snapshot_name": name + "_update",
					"disk_id":       "${alicloud_ens_disk.disk.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "SnapShotDescription_autotest",
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
						"snapshot_name": name + "_update",
						"disk_id":       acctest.CHECKSET,
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

var AlicloudEnsSnapshotMap5162 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudEnsSnapshotBasicDependence5162(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_disk" "disk" {
  category      = "cloud_efficiency"
  size          = "20"
  payment_type  = "PayAsYouGo"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}


`, name)
}

// Case 5162  twin
func TestAccAliCloudEnsSnapshot_basic5162_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_snapshot.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEnsSnapshotMap5162)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsSnapshot")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senssnapshot%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEnsSnapshotBasicDependence5162)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   "SnapShotDescription_UPDATE_autotest",
					"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
					"snapshot_name": name,
					"disk_id":       "${alicloud_ens_disk.disk.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "SnapShotDescription_UPDATE_autotest",
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
						"snapshot_name": name,
						"disk_id":       acctest.CHECKSET,
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

// Test Ens Snapshot. <<< Resource test cases, automatically generated.
