package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudFCServicesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_services.default"
	name := fmt.Sprintf("tf-testacc-fc-service-ds-basic-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceFCServicesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": alicloud_fc_service.default.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_service.default.name}_fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{alicloud_fc_service.default.service_id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_fc_service.default.service_id}_fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": alicloud_fc_service.default.name,
			"ids":        []string{alicloud_fc_service.default.service_id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_service.default.name}_fake",
			"ids":        []string{alicloud_fc_service.default.service_id},
		}),
	}

	var existFCServicesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"services.#":                        "1",
			"ids.#":                             "1",
			"names.#":                           "1",
			"services.0.id":                     CHECKSET,
			"services.0.name":                   name,
			"services.0.description":            name + "-description",
			"services.0.role":                   CHECKSET,
			"services.0.internet_access":        "true",
			"services.0.creation_time":          CHECKSET,
			"services.0.last_modification_time": CHECKSET,
			"services.0.log_config.#":           "1",
			"services.0.log_config.0.project":   name + "-project",
			"services.0.log_config.0.logstore":  name + "-store",
		}
	}

	var fakeFCServicesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"services.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var fcServicesRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existFCServicesMapFunc,
		fakeMapFunc:  fakeFCServicesMapFunc,
	}

	fcServicesRecordsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceFCServicesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_log_project" "default" {
    name = "${var.name}-project"
}

resource "alicloud_log_store" "default" {
    project = alicloud_log_project.default.name
    name = "${var.name}-store"
}

resource "alicloud_ram_role" "default" {
    name = var.name
    document = <<DEFINITION
    {
        "Statement": [
            {
                "Action": "sts:AssumeRole",
                "Effect": "Allow",
                "Principal": {
                    "Service": [
                        "fc.aliyuncs.com"
                    ]
                }
            }
        ],
        "Version": "1"
    }
    DEFINITION
    description = "this is a test"
    force = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
    role_name = alicloud_ram_role.default.name
    policy_name = "AliyunLogFullAccess"
    policy_type = "System"
}

resource "alicloud_fc_service" "default" {
    name = var.name
    description = "${var.name}-description"
    log_config {
	    project = alicloud_log_store.default.project
	    logstore = alicloud_log_store.default.name
    }
    role = alicloud_ram_role.default.arn
    depends_on = ["alicloud_ram_role_policy_attachment.default"]
    internet_access = true
}
`, name)
}
