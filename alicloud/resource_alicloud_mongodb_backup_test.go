package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Mongodb Backup. >>> Resource test cases.
func TestAccAliCloudMongodbBackup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_backup.default"
	ra := resourceAttrInit(resourceId, AlicloudMongodbBackupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongodbBackupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_mongodb_instance.defaultHrZmxC.id}",
					"backup_method":           "Snapshot",
					"backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"backup_method":           "Snapshot",
						"backup_retention_period": "7",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_retention_period"},
			},
		},
	})
}

var AlicloudMongodbBackupMap0 = map[string]string{
	"backup_id":         CHECKSET,
	"backup_job_id":     CHECKSET,
	"status":            CHECKSET,
	"backup_type":       CHECKSET,
	"backup_mode":       CHECKSET,
	"backup_size":       CHECKSET,
	"backup_start_time": CHECKSET,
	"backup_end_time":   CHECKSET,
}

func AlicloudMongodbBackupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "region_id" {
  default = "cn-shanghai"
}

resource "alicloud_vpc" "defaultie35CW" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultg0DCAR" {
  vpc_id     = alicloud_vpc.defaultie35CW.id
  zone_id    = var.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_instance" "defaultHrZmxC" {
  engine_version      = "4.4"
  storage_type        = "cloud_essd1"
  vswitch_id          = alicloud_vswitch.defaultg0DCAR.id
  db_instance_storage = "20"
  vpc_id              = alicloud_vpc.defaultie35CW.id
  db_instance_class   = "mdb.shard.4x.large.d"
  storage_engine      = "WiredTiger"
  network_type        = "VPC"
  zone_id             = var.zone_id
}


`, name)
}

// Test Mongodb Backup. <<< Resource test cases.
