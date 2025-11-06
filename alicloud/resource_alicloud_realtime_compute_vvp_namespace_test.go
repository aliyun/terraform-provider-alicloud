// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRealtimeComputeVvpNamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_vvp_namespace.default"
	ra := resourceAttrInit(resourceId, AliCloudRealtimeComputeVvpNamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVvpNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnamespace%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRealtimeComputeVvpNamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_realtime_compute_vvp_instance.default.id}",
					"namespace":   name,
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "1",
							"memory_gb": "8",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":     CHECKSET,
						"namespace":       name,
						"resource_spec.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "2",
							"memory_gb": "16",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_spec.#": "1",
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

func TestAccAliCloudRealtimeComputeVvpNamespace_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_vvp_namespace.default"
	ra := resourceAttrInit(resourceId, AliCloudRealtimeComputeVvpNamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVvpNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnamespace%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRealtimeComputeVvpNamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_realtime_compute_vvp_instance.default.id}",
					"namespace":   name,
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "1",
							"memory_gb": "8",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":     CHECKSET,
						"namespace":       name,
						"resource_spec.#": "1",
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

var AliCloudRealtimeComputeVvpNamespaceMap0 = map[string]string{
	"status": CHECKSET,
}

func AliCloudRealtimeComputeVvpNamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-i"
	}

	resource "alicloud_oss_bucket" "defaultOSS" {
  		bucket = var.name
	}

	resource "alicloud_realtime_compute_vvp_instance" "default" {
  		vvp_instance_name = var.name
  		vpc_id            = data.alicloud_vpcs.default.ids.0
  		payment_type      = "PayAsYouGo"
  		zone_id           = "cn-hangzhou-i"
  		vswitch_ids       = [data.alicloud_vswitches.default.ids.0]
  		storage {
    		oss {
      			bucket = alicloud_oss_bucket.defaultOSS.bucket
    		}
  		}
	}
`, name)
}
