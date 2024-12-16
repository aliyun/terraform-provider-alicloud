package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudKvStoreBackupPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.KvStoreSupportRegions)
	resourceId := "alicloud_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudKvStoreBackupPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreBackupPolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKvStoreBackupPolicyBasicDependence0)
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
					"instance_id":   "${alicloud_kvstore_instance.default.id}",
					"backup_period": []string{"Tuesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":     CHECKSET,
						"backup_period.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Tuesday", "Wednesday", "Sunday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudKvStoreBackupPolicy_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.KvStoreSupportRegions)
	resourceId := "alicloud_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudKvStoreBackupPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreBackupPolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKvStoreBackupPolicyBasicDependence0)
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
					"instance_id":   "${alicloud_kvstore_instance.default.id}",
					"backup_time":   "10:00Z-11:00Z",
					"backup_period": []string{"Tuesday", "Wednesday", "Sunday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":     CHECKSET,
						"backup_time":     "10:00Z-11:00Z",
						"backup_period.#": "3",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudKvStoreBackupPolicyMap0 = map[string]string{}

func AliCloudKvStoreBackupPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_kvstore_zones" "default" {
  		instance_charge_type = "PostPaid"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_kvstore_zones.default.zones.0.id
  		vpc_id  = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_kvstore_instance" "default" {
  		zone_id          = data.alicloud_kvstore_zones.default.zones.0.id
  		instance_class   = "redis.master.small.default"
  		db_instance_name = var.name
  		engine_version   = "5.0"
  		vswitch_id       = data.alicloud_vswitches.default.ids.0
	}
`, name)
}
