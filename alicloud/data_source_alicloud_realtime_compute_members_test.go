package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudRealtimeComputeMembersDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_realtime_compute_members.default"
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRealtimeComputeMembersConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace":   "${alicloud_realtime_compute_member.default.namespace}",
			"resource_id": "${alicloud_realtime_compute_member.default.resource_id}",
			"ids":         []string{"${alicloud_realtime_compute_member.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace":   "${alicloud_realtime_compute_member.default.namespace}",
			"resource_id": "${alicloud_realtime_compute_member.default.resource_id}",
			"ids":         []string{"${alicloud_realtime_compute_member.default.id}_fake"},
		}),
	}

	var existAliCloudRealtimeComputeMembersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"members.#":             "1",
			"members.0.id":          CHECKSET,
			"members.0.resource_id": CHECKSET,
			"members.0.namespace":   CHECKSET,
			"members.0.member":      CHECKSET,
			"members.0.role":        CHECKSET,
		}
	}

	var fakeAliCloudRealtimeComputeMembersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"members.#": "0",
		}
	}

	var aliCloudRealtimeComputeMembersInfo = dataSourceAttr{
		resourceId:   "data.alicloud_realtime_compute_members.default",
		existMapFunc: existAliCloudRealtimeComputeMembersMapFunc,
		fakeMapFunc:  fakeAliCloudRealtimeComputeMembersMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudRealtimeComputeMembersInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func dataSourceRealtimeComputeMembersConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
   default = "%s"
}

data "alicloud_oss_buckets" "default" {
}

resource "alicloud_vpc" "default" {
 is_default = false
 cidr_block = "172.16.0.0/16"
 vpc_name   = "test-tf-vpc"
}

resource "alicloud_vswitch" "default" {
 is_default   = false
 vpc_id       = alicloud_vpc.default.id
 zone_id      = "cn-hangzhou-i"
 cidr_block   = "172.16.0.0/24"
 vswitch_name = "test-tf-vSwitch"
}

resource "alicloud_ram_user" "default" {
 name         = var.name
 display_name = "displayname"
 mobile       = "86-18888888888"
 email        = "hello.uuu@aaa.com"
 comments     = "yoyoyo"
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
 vvp_instance_name = var.name
 storage {
   oss {
     bucket = data.alicloud_oss_buckets.default.buckets.0.name
   }
 }
 vpc_id      = alicloud_vpc.default.id
 vswitch_ids = [alicloud_vswitch.default.id]
 resource_spec {
   cpu       = "8"
   memory_gb = "32"
 }
 payment_type = "PayAsYouGo"
 zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_member" "default" {
 member      = alicloud_ram_user.default.id
 namespace   = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
 resource_id = alicloud_realtime_compute_vvp_instance.default.resource_id
 role        = "viewer"
}
`, name)
}
