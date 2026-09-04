package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudRealtimeComputeMember_basic11888(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_member.default"
	ra := resourceAttrInit(resourceId, AliCloudRealtimeComputeMemberMap11888)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRealtimeComputeMemberBasicDependence11888)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"member":      "${alicloud_ram_user.default.id}",
					"resource_id": "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"role":        "viewer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"member":      CHECKSET,
						"resource_id": CHECKSET,
						"namespace":   CHECKSET,
						"role":        "viewer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role": "owner",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role": "owner",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role": "editor",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role": "editor",
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

func TestAccAliCloudRealtimeComputeMember_basic11888_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_member.default"
	ra := resourceAttrInit(resourceId, AliCloudRealtimeComputeMemberMap11888)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRealtimeComputeMemberBasicDependence11888)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"member":      "${alicloud_ram_user.default.id}",
					"resource_id": "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"role":        "viewer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"member":      CHECKSET,
						"resource_id": CHECKSET,
						"namespace":   CHECKSET,
						"role":        "viewer",
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

var AliCloudRealtimeComputeMemberMap11888 = map[string]string{}

func AliCloudRealtimeComputeMemberBasicDependence11888(name string) string {
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
  vvp_instance_name = "code-test-tf"
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
`, name)
}
