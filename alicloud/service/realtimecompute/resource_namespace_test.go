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

func TestAccAliCloudRealtimeComputeNamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_namespace.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRealtimeComputeNamespaceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeNamespace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRealtimeComputeNamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_realtime_compute_vvp_instance.default.id}",
					"namespace":   name,
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "2",
							"memory_gb": "8",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":     acctest.CHECKSET,
						"namespace":       name,
						"resource_spec.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "3",
							"memory_gb": "12",
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
				Config: testAccConfig(map[string]interface{}{
					"ha": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha": "false",
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

func TestAccAliCloudRealtimeComputeNamespace_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_namespace.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRealtimeComputeNamespaceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeNamespace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRealtimeComputeNamespaceBasicDependence0)
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
					"instance_id": "${alicloud_realtime_compute_vvp_instance.default.id}",
					"namespace":   name,
					"ha":          false,
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "2",
							"memory_gb": "8",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":     acctest.CHECKSET,
						"namespace":       name,
						"ha":              "false",
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

var AliCloudRealtimeComputeNamespaceMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

func AliCloudRealtimeComputeNamespaceBasicDependence0(name string) string {
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
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = data.alicloud_oss_buckets.default.buckets.0.name
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = [alicloud_vswitch.default.id]
  resource_spec {
    cpu       = "10"
    memory_gb = "40"
  }
  default_namespace_resource_spec {
    cpu       = "2"
    memory_gb = "8"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}
`, name)
}
