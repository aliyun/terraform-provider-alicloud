package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test RealtimeCompute Variable. >>> Resource test cases, automatically generated.
func TestAccAliCloudRealtimeComputeVariable_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_variable.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeVariableMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVariable")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-acc-var-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeVariableBasicDependence0)
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
					"workspace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"name":        name,
					"kind":        "Clear",
					"value":       "test-value",
					"description": "test variable description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":   CHECKSET,
						"namespace":   CHECKSET,
						"name":        name,
						"kind":        "Clear",
						"value":       "test-value",
						"description": "test variable description",
					}),
				),
			},
			// Step 2: modify ONLY value (description unchanged) — verify value update in isolation.
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"name":        name,
					"kind":        "Clear",
					"value":       "updated-value",
					"description": "test variable description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":   CHECKSET,
						"namespace":   CHECKSET,
						"name":        name,
						"kind":        "Clear",
						"value":       "updated-value",
						"description": "test variable description",
					}),
				),
			},
			// Step 3: modify ONLY description (value unchanged) — verify description update in isolation.
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"name":        name,
					"kind":        "Clear",
					"value":       "updated-value",
					"description": "updated variable description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":   CHECKSET,
						"namespace":   CHECKSET,
						"name":        name,
						"kind":        "Clear",
						"value":       "updated-value",
						"description": "updated variable description",
					}),
				),
			},
			// Step 4: modify ONLY description (clear it, value unchanged) — verify description removal in isolation.
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"namespace":   "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"name":        name,
					"kind":        "Clear",
					"value":       "updated-value",
					"description": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":   CHECKSET,
						"namespace":   CHECKSET,
						"name":        name,
						"kind":        "Clear",
						"value":       "updated-value",
						"description": "",
					}),
				),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kind", "value"},
				ResourceName:            resourceId,
			},
		},
	})
}

func AlicloudRealtimeComputeVariableBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_Vpc" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-variable"
}

resource "alicloud_vswitch" "create_Vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-variable"
}

resource "alicloud_oss_bucket" "create_bucket" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-test-tf-variable"
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

var AlicloudRealtimeComputeVariableMap0 = map[string]string{}
