package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudRealtimeComputeVariablesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_realtime_compute_variables.default"
	name := fmt.Sprintf("tfacc%d", rand)
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

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"workspace":  "${alicloud_realtime_compute_variable.default.workspace}",
			"namespace":  "${alicloud_realtime_compute_variable.default.namespace}",
			"name_regex": "${alicloud_realtime_compute_variable.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"workspace":  "${alicloud_realtime_compute_variable.default.workspace}",
			"namespace":  "${alicloud_realtime_compute_variable.default.namespace}",
			"name_regex": "${alicloud_realtime_compute_variable.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"workspace":  "${alicloud_realtime_compute_variable.default.workspace}",
			"namespace":  "${alicloud_realtime_compute_variable.default.namespace}",
			"ids":        []string{"${alicloud_realtime_compute_variable.default.id}"},
			"name_regex": "${alicloud_realtime_compute_variable.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"workspace":  "${alicloud_realtime_compute_variable.default.workspace}",
			"namespace":  "${alicloud_realtime_compute_variable.default.namespace}",
			"ids":        []string{"${alicloud_realtime_compute_variable.default.id}_fake"},
			"name_regex": "${alicloud_realtime_compute_variable.default.name}_fake",
		}),
	}

	var existAliCloudRealtimeComputeVariablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"variables.#":             "1",
			"variables.0.id":          CHECKSET,
			"variables.0.workspace":   CHECKSET,
			"variables.0.namespace":   CHECKSET,
			"variables.0.name":        CHECKSET,
			"variables.0.kind":        CHECKSET,
			"variables.0.value":       CHECKSET,
			"variables.0.description": CHECKSET,
		}
	}

	var fakeAliCloudRealtimeComputeVariablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
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

	aliCloudRealtimeComputeVariablesInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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

resource "alicloud_realtime_compute_variable" "default" {
  name        = var.name
  namespace   = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  workspace   = alicloud_realtime_compute_vvp_instance.default.resource_id
  kind        = "Clear"
  value       = "YourPassword123!"
  description = var.name
}
`, name)
}
