package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCrEEReposDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	repoName := fmt.Sprintf("tf-testacc-cr-ee-repo-%d", rand)
	namespaceName := fmt.Sprintf("tf-testacc-cr-ee-ns-%d", rand)
	resourceId := "data.alicloud_cr_ee_repos.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprint(rand),
		dataSourceCrEEReposConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${data.alicloud_cr_ee_instances.default.ids.0}",
			"namespace":   "${alicloud_cr_ee_namespace.default.name}",
			"name_regex":  "${alicloud_cr_ee_repo.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${data.alicloud_cr_ee_instances.default.ids.0}",
			"namespace":   "${alicloud_cr_ee_namespace.default.name}",
			"name_regex":  "${alicloud_cr_ee_repo.default.name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${data.alicloud_cr_ee_instances.default.ids.0}",
			"namespace":   "${alicloud_cr_ee_namespace.default.name}",
			"ids":         []string{"${alicloud_cr_ee_repo.default.repo_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${data.alicloud_cr_ee_instances.default.ids.0}",
			"namespace":   "${alicloud_cr_ee_namespace.default.name}",
			"ids":         []string{"${alicloud_cr_ee_repo.default.repo_id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${data.alicloud_cr_ee_instances.default.ids.0}",
			"namespace":   "${alicloud_cr_ee_namespace.default.name}",
			"name_regex":  "${alicloud_cr_ee_repo.default.name}",
			"ids":         []string{"${alicloud_cr_ee_repo.default.repo_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${data.alicloud_cr_ee_instances.default.ids.0}",
			"namespace":   "${alicloud_cr_ee_namespace.default.name}",
			"ids":         []string{"${alicloud_cr_ee_repo.default.repo_id}-fake"},
			"name_regex":  "${alicloud_cr_ee_repo.default.name}-fake",
		}),
	}

	var existCrEEReposMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"ids.0":               CHECKSET,
			"names.#":             "1",
			"names.0":             repoName,
			"repos.#":             "1",
			"repos.0.instance_id": CHECKSET,
			"repos.0.namespace":   namespaceName,
			"repos.0.id":          CHECKSET,
			"repos.0.name":        repoName,
			"repos.0.summary":     "test summary",
			"repos.0.repo_type":   "PRIVATE",
		}
	}

	var fakeCrEEReposMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"repos.#": "0",
		}
	}

	var crEEReposCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCrEEReposMapFunc,
		fakeMapFunc:  fakeCrEEReposMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithCrEE(t)
	}
	crEEReposCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceCrEEReposConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testacc-cr-ee-repo-%s"
	}

	data "alicloud_cr_ee_instances" "default" {
	}
	
	resource "alicloud_cr_ee_namespace" "default" {
		instance_id = "${data.alicloud_cr_ee_instances.default.ids.0}"
		name = "tf-testacc-cr-ee-ns-%s"
		auto_create	= true
		default_visibility = "PRIVATE"
	}
	
	resource "alicloud_cr_ee_repo" "default" {
		instance_id = "${data.alicloud_cr_ee_instances.default.ids.0}"
		namespace = "${alicloud_cr_ee_namespace.default.name}"
		name = "${var.name}"
		summary = "test summary"
		repo_type = "PRIVATE"
		detail = "test detail"
	}

	`, name, name)
}
