package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRealtimeComputeMember_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_member.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeMemberMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeFlinkMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-flinkmember%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeMemberBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_id": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"member":      "${alicloud_ram_user.user.id}",
					"role":        "VIEWER",
				}),
				Check: testAccCheck(map[string]string{}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_id": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"member":      "${alicloud_ram_user.user.id}",
					"role":        "EDITOR",
				}),
				Check: testAccCheck(map[string]string{}),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudRealtimeComputeMemberMap0 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudRealtimeComputeMemberBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ram_user" "user" {
  name = var.name
}

resource "alicloud_vpc" "create_Vpc" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-flink-member"
}

resource "alicloud_vswitch" "create_Vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-flink-member"
}

resource "alicloud_oss_bucket" "create_bucket" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket.id
    }
  }
  vpc_id      = alicloud_vpc.create_Vpc.id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch.zone_id
}

`, name)
}
