package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudRealtimeComputeVariablesDataSource_basic11922(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_realtime_compute_variables.default"
	name := fmt.Sprintf("tfaccvariable%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRealtimeComputeVariablesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"workspace": "${alicloud_realtime_compute_variable.default.workspace}",
			"namespace": "${alicloud_realtime_compute_variable.default.namespace}",
			"ids":       []string{"${alicloud_realtime_compute_variable.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"workspace": "${alicloud_realtime_compute_variable.default.workspace}",
			"namespace": "${alicloud_realtime_compute_variable.default.namespace}",
			"ids":       []string{"${alicloud_realtime_compute_variable.default.id}_fake"},
		}),
	}

	var existAliCloudRealtimeComputeVariablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"variables.#":           "1",
			"variables.0.id":        CHECKSET,
			"variables.0.workspace": CHECKSET,
			"variables.0.namespace": CHECKSET,
			"variables.0.name":      CHECKSET,
			"variables.0.kind":      "Plain",
			"variables.0.value":     "test-value-1",
		}
	}

	var fakeAliCloudRealtimeComputeVariablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"variables.#": "0",
		}
	}

	var aliCloudRealtimeComputeVariablesInfo = dataSourceAttr{
		resourceId:   "data.alicloud_realtime_compute_variables.default",
		existMapFunc: existAliCloudRealtimeComputeVariablesMapFunc,
		fakeMapFunc:  fakeAliCloudRealtimeComputeVariablesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudRealtimeComputeVariablesInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func dataSourceRealtimeComputeVariablesConfig(name string) string {
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
    cpu       = "8"
    memory_gb = "32"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_variable" "default" {
  workspace   = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace   = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name        = var.name
  kind        = "Plain"
  value       = "test-value-1"
  description = "test-description-1"
}
`, name)
}
