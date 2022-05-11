package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCREESyncRulesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ee-sr-%d", rand)
	resourceId := "data.alicloud_cr_ee_sync_rules.default"
	region := os.Getenv("ALICLOUD_REGION")

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name,
		dataSourceCrEESyncRulesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"name_regex":  "${alicloud_cr_ee_sync_rule.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"name_regex":  "${alicloud_cr_ee_sync_rule.default.name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"ids":         []string{"${alicloud_cr_ee_sync_rule.default.rule_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"ids":         []string{"${alicloud_cr_ee_sync_rule.default.rule_id}-fake"},
		}),
	}

	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"namespace_name": "${alicloud_cr_ee_namespace.source_ns.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"namespace_name": "${alicloud_cr_ee_namespace.source_ns.name}-fake",
		}),
	}

	repoConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"repo_name":   "${alicloud_cr_ee_repo.source_repo.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"repo_name":   "${alicloud_cr_ee_repo.source_repo.name}-fake",
		}),
	}

	targetInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":        "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"target_instance_id": "${alicloud_cr_ee_namespace.target_ns.instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":        "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"target_instance_id": "${alicloud_cr_ee_namespace.target_ns.instance_id}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":        "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"namespace_name":     "${alicloud_cr_ee_namespace.source_ns.name}",
			"name_regex":         "${alicloud_cr_ee_sync_rule.default.name}",
			"ids":                []string{"${alicloud_cr_ee_sync_rule.default.rule_id}"},
			"repo_name":          "${alicloud_cr_ee_repo.source_repo.name}",
			"target_instance_id": "${alicloud_cr_ee_namespace.target_ns.instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":        "${alicloud_cr_ee_namespace.source_ns.instance_id}",
			"namespace_name":     "${alicloud_cr_ee_namespace.source_ns.name}-fake",
			"name_regex":         "${alicloud_cr_ee_sync_rule.default.name}-fake",
			"ids":                []string{"${alicloud_cr_ee_sync_rule.default.rule_id}-fake"},
			"repo_name":          "${alicloud_cr_ee_repo.source_repo.name}-fake",
			"target_instance_id": "${alicloud_cr_ee_namespace.target_ns.instance_id}-fake",
		}),
	}

	var existCrEESyncRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         CHECKSET,
			"names.#":                       "1",
			"names.0":                       name,
			"rules.#":                       "1",
			"rules.0.instance_id":           CHECKSET,
			"rules.0.namespace_name":        name,
			"rules.0.id":                    CHECKSET,
			"rules.0.name":                  name,
			"rules.0.region_id":             region,
			"rules.0.repo_name":             name,
			"rules.0.sync_direction":        "FROM",
			"rules.0.sync_scope":            "REPO",
			"rules.0.sync_trigger":          "PASSIVE",
			"rules.0.tag_filter":            ".*",
			"rules.0.target_instance_id":    CHECKSET,
			"rules.0.target_namespace_name": name,
			"rules.0.target_region_id":      region,
			"rules.0.target_repo_name":      name,
		}
	}

	var fakeCrEESyncRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}

	var crEESyncRulesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCrEESyncRulesMapFunc,
		fakeMapFunc:  fakeCrEESyncRulesMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	crEESyncRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, namespaceConf,
		repoConf, targetInstanceIdConf, allConf)
}

func dataSourceCrEESyncRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "region" {
	default = "%s"
}

variable "name" {
	default = "%s"
}
data "alicloud_cr_ee_instances" "default" {}

resource "alicloud_cr_ee_namespace" "source_ns" {
	instance_id = data.alicloud_cr_ee_instances.default.ids.0
	name = "${var.name}"
	auto_create	= true
	default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_namespace" "target_ns" {
	instance_id = data.alicloud_cr_ee_instances.default.ids.1
	name = "${var.name}"
	auto_create	= true
	default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "source_repo" {
	instance_id = "${alicloud_cr_ee_namespace.source_ns.instance_id}"
	namespace = "${alicloud_cr_ee_namespace.source_ns.name}"
	name = "${var.name}"
	summary = "test"
	repo_type = "PRIVATE"
	detail = "test"
}

resource "alicloud_cr_ee_repo" "target_repo" {
	instance_id = "${alicloud_cr_ee_namespace.target_ns.instance_id}"
	namespace = "${alicloud_cr_ee_namespace.target_ns.name}"
	name = "${var.name}"
	summary = "test"
	repo_type = "PRIVATE"
	detail = "test"
}

resource "alicloud_cr_ee_sync_rule" "default" {
	instance_id = "${alicloud_cr_ee_namespace.source_ns.instance_id}"
	namespace_name = "${alicloud_cr_ee_namespace.source_ns.name}"
	name = "${var.name}"
	target_region_id = "${var.region}"
	target_instance_id = "${alicloud_cr_ee_namespace.target_ns.instance_id}"
	target_namespace_name = "${alicloud_cr_ee_namespace.target_ns.name}"
	tag_filter = ".*"
	repo_name = "${alicloud_cr_ee_repo.source_repo.name}"
	target_repo_name = "${alicloud_cr_ee_repo.target_repo.name}"
}
`, defaultRegionToTest, name)
}
