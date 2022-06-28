package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSDedicatedHostsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_ecs_dedicated_hosts.default"
	name := fmt.Sprintf("tf_testAccEcsDedicatedHostsDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsDedicatedHostsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_dedicated_host.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_dedicated_host.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"name_regex": name + "fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"status": "UnderAssessment",
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"dedicated_host_type": "ddh.g6",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"dedicated_host_type": "ddh.g5",
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"tags": map[string]string{
				"Create": "TF",
				"For":    "ddh-test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"tags": map[string]string{
				"Create": "ddh-test",
				"For":    "TF",
			},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"dedicated_host_type": "ddh.g6",
			"status":              "Available",
			"name_regex":          name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_ecs_dedicated_host.default.id}"},
			"dedicated_host_type": "ddh.g6",
			"name_regex":          name + "fake",
			"status":              "UnderAssessment",
		}),
	}

	var existKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"ids.0":                                 CHECKSET,
			"names.#":                               "1",
			"names.0":                               CHECKSET,
			"hosts.0.action_on_maintenance":         "Migrate",
			"hosts.0.auto_placement":                "on",
			"hosts.0.auto_release_time":             "",
			"hosts.0.id":                            CHECKSET,
			"hosts.0.dedicated_host_id":             CHECKSET,
			"hosts.0.dedicated_host_name":           CHECKSET,
			"hosts.0.dedicated_host_type":           CHECKSET,
			"hosts.0.description":                   "From_Terraform",
			"hosts.0.expired_time":                  CHECKSET,
			"hosts.0.gpu_spec":                      "",
			"hosts.0.machine_id":                    CHECKSET,
			"hosts.0.payment_type":                  "PostPaid",
			"hosts.0.physical_gpus":                 CHECKSET,
			"hosts.0.resource_group_id":             "",
			"hosts.0.sale_cycle":                    "",
			"hosts.0.sockets":                       CHECKSET,
			"hosts.0.status":                        "Available",
			"hosts.0.supported_instance_types_list": NOSET,
			"hosts.0.zone_id":                       CHECKSET,
		}
	}

	var fakeKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"hosts.#": "0",
		}
	}

	var ecsDedicatedHostsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsSecretVersionsMapFunc,
		fakeMapFunc:  fakeKmsSecretVersionsMapFunc,
	}

	ecsDedicatedHostsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, typeConf, statusConf, tagsConf, allConf)
}

func dataSourceEcsDedicatedHostsConfigDependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_ecs_dedicated_host" "default" {
		  dedicated_host_type = "ddh.g6"
		  description = "From_Terraform"
		  dedicated_host_name = "%s"
          action_on_maintenance = "Migrate"
          tags = {
			Create = "TF"
    		For = "ddh-test",
  			}
		}
	`, name)
}
