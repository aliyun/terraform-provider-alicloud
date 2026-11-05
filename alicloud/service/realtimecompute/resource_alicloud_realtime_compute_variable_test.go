package realtimecompute_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/realtimecompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRealtimeComputeVariable_basic11188(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_variable.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRealtimeComputeVariableMap11188)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVariable")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRealtimeComputeVariableBasicDependence11188)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"workspace": acctest.CHECKSET,
						"namespace": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRealtimeComputeVariableMap11188)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVariable")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRealtimeComputeVariableBasicDependence11188)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"workspace":   acctest.CHECKSET,
						"namespace":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRealtimeComputeVariableMap11188)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVariable")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRealtimeComputeVariableBasicDependence11188)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"workspace": acctest.CHECKSET,
						"namespace": acctest.CHECKSET,
						"name":      name,
						"kind":      "Encrypted",
						"value":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "YourPassword123!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRealtimeComputeVariableMap11188)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVariable")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRealtimeComputeVariableBasicDependence11188)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"workspace":   acctest.CHECKSET,
						"namespace":   acctest.CHECKSET,
						"name":        name,
						"kind":        "Encrypted",
						"value":       acctest.CHECKSET,
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
