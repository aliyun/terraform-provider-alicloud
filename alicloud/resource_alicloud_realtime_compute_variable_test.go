package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRealtimeComputeVariable_basic11188(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_variable.default"
	ra := resourceAttrInit(resourceId, AliCloudRealtimeComputeVariableMap11188)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVariable")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRealtimeComputeVariableBasicDependence11188)
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
					"workspace": "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace": "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":      name,
					"kind":      "Clear",
					"value":     "YourPassword123!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace": CHECKSET,
						"namespace": CHECKSET,
						"name":      name,
						"kind":      "Clear",
						"value":     "YourPassword123!",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "YourPassword123!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "YourPassword123!update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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

func TestAccAliCloudRealtimeComputeVariable_basic11188_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_variable.default"
	ra := resourceAttrInit(resourceId, AliCloudRealtimeComputeVariableMap11188)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVariable")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRealtimeComputeVariableBasicDependence11188)
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
					"workspace":   "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":        name,
					"kind":        "Clear",
					"value":       "YourPassword123!",
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":   CHECKSET,
						"namespace":   CHECKSET,
						"name":        name,
						"kind":        "Clear",
						"value":       "YourPassword123!",
						"description": name,
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

func TestAccAliCloudRealtimeComputeVariable_basic11190(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_variable.default"
	ra := resourceAttrInit(resourceId, AliCloudRealtimeComputeVariableMap11188)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVariable")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRealtimeComputeVariableBasicDependence11188)
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
					"workspace": "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace": "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":      name,
					"kind":      "Encrypted",
					"value":     "YourPassword123!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace": CHECKSET,
						"namespace": CHECKSET,
						"name":      name,
						"kind":      "Encrypted",
						"value":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "YourPassword123!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func TestAccAliCloudRealtimeComputeVariable_basic11190_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_variable.default"
	ra := resourceAttrInit(resourceId, AliCloudRealtimeComputeVariableMap11188)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVariable")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRealtimeComputeVariableBasicDependence11188)
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
					"workspace":   "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":        name,
					"kind":        "Encrypted",
					"value":       "YourPassword123!",
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":   CHECKSET,
						"namespace":   CHECKSET,
						"name":        name,
						"kind":        "Encrypted",
						"value":       CHECKSET,
						"description": name,
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

var AliCloudRealtimeComputeVariableMap11188 = map[string]string{}

func AliCloudRealtimeComputeVariableBasicDependence11188(name string) string {
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
