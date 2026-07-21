package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCrEEReposDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cr_ee_repo.default.repo_id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cr_ee_repo.default.repo_id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cr_ee_repo.default.name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cr_ee_repo.default.name}_fake"`,
		}),
	}

	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"namespace": `"${alicloud_cr_ee_repo.default.namespace}"`,
		}),
		fakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cr_ee_repo.default.name}_fake"`,
			"namespace":  `"${alicloud_cr_ee_repo.default.namespace}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"name_regex": `"${alicloud_cr_ee_repo.default.name}"`,
			"namespace":  `"${alicloud_cr_ee_repo.default.namespace}"`,
		}),
		fakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cr_ee_repo.default.repo_id}_fake"]`,
			"name_regex": `"${alicloud_cr_ee_repo.default.name}_fake"`,
			"namespace":  `"${alicloud_cr_ee_repo.default.namespace}"`,
		}),
	}

	var existAliCloudCrEEReposDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"repos.#":             "1",
			"repos.0.id":          CHECKSET,
			"repos.0.instance_id": CHECKSET,
			"repos.0.namespace":   CHECKSET,
			"repos.0.name":        CHECKSET,
			"repos.0.summary":     CHECKSET,
			"repos.0.repo_type":   CHECKSET,
			"repos.0.tags.#":      "0",
		}
	}

	var fakeAliCloudCrEEReposDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"repos.#": "0",
		}
	}

	var alicloudCrEEReposCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cr_ee_repos.default",
		existMapFunc: existAliCloudCrEEReposDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudCrEEReposDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudCrEEReposCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, namespaceConf, allConf)
}

func TestAccAliCloudCrEEReposDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"enable_details": "false",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_cr_ee_repo.default.name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_cr_ee_repo.default.name}"`,
			"enable_details": "false",
		}),
	}

	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"namespace":      `"${alicloud_cr_ee_repo.default.namespace}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"namespace":      `"${alicloud_cr_ee_repo.default.namespace}"`,
			"enable_details": "false",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"name_regex":     `"${alicloud_cr_ee_repo.default.name}"`,
			"namespace":      `"${alicloud_cr_ee_repo.default.namespace}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"name_regex":     `"${alicloud_cr_ee_repo.default.name}"`,
			"namespace":      `"${alicloud_cr_ee_repo.default.namespace}"`,
			"enable_details": "false",
		}),
	}

	var existAliCloudCrEEReposDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"repos.#":             "1",
			"repos.0.id":          CHECKSET,
			"repos.0.instance_id": CHECKSET,
			"repos.0.namespace":   CHECKSET,
			"repos.0.name":        CHECKSET,
			"repos.0.summary":     CHECKSET,
			"repos.0.repo_type":   CHECKSET,
			// Currently, Terraform doesn't have an environment testing ListRepoTag
			//"repos.0.tags.#":              "1",
			//"repos.0.tags.0.tag":          CHECKSET,
			//"repos.0.tags.0.image_id":     CHECKSET,
			//"repos.0.tags.0.image_size":   CHECKSET,
			//"repos.0.tags.0.digest":       CHECKSET,
			//"repos.0.tags.0.status":       CHECKSET,
			//"repos.0.tags.0.image_create": CHECKSET,
			//"repos.0.tags.0.image_update": CHECKSET,
		}
	}

	var fakeAliCloudCrEEReposDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"repos.#":             "1",
			"repos.0.id":          CHECKSET,
			"repos.0.instance_id": CHECKSET,
			"repos.0.namespace":   CHECKSET,
			"repos.0.name":        CHECKSET,
			"repos.0.summary":     CHECKSET,
			"repos.0.repo_type":   CHECKSET,
			"repos.0.tags.#":      "0",
		}
	}

	var alicloudCrEEReposCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cr_ee_repos.default",
		existMapFunc: existAliCloudCrEEReposDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudCrEEReposDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudCrEEReposCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, namespaceConf, allConf)
}

func testAccCheckAliCloudCrEEReposDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testacc-repo-%d"
	}

	data "alicloud_cr_ee_instances" "default"{
		name_regex = "default-nodeleting"
	}

	resource "alicloud_cr_ee_namespace" "default" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.0
  		name               = var.name
  		auto_create        = true
  		default_visibility = "PRIVATE"
	}

	resource "alicloud_cr_ee_repo" "default" {
  		instance_id = alicloud_cr_ee_namespace.default.instance_id
  		namespace   = alicloud_cr_ee_namespace.default.name
  		name        = var.name
  		repo_type   = "PRIVATE"
  		summary     = var.name
	}

	data "alicloud_cr_ee_repos" "default" {
  		instance_id = alicloud_cr_ee_repo.default.instance_id
 		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
