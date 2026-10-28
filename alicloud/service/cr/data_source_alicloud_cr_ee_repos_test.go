package cr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCrEEReposDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cr_ee_repo.default.repo_id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cr_ee_repo.default.repo_id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cr_ee_repo.default.name}"`,
		}),
		FakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cr_ee_repo.default.name}_fake"`,
		}),
	}

	namespaceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"namespace": `"${alicloud_cr_ee_repo.default.namespace}"`,
		}),
		FakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cr_ee_repo.default.name}_fake"`,
			"namespace":  `"${alicloud_cr_ee_repo.default.namespace}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"name_regex": `"${alicloud_cr_ee_repo.default.name}"`,
			"namespace":  `"${alicloud_cr_ee_repo.default.namespace}"`,
		}),
		FakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
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
			"repos.0.id":          acctest.CHECKSET,
			"repos.0.instance_id": acctest.CHECKSET,
			"repos.0.namespace":   acctest.CHECKSET,
			"repos.0.name":        acctest.CHECKSET,
			"repos.0.summary":     acctest.CHECKSET,
			"repos.0.repo_type":   acctest.CHECKSET,
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

	var alicloudCrEEReposCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cr_ee_repos.default",
		ExistMapFunc: existAliCloudCrEEReposDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudCrEEReposDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	alicloudCrEEReposCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, namespaceConf, allConf)
}

func TestAccAliCloudCrEEReposDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"enable_details": "false",
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_cr_ee_repo.default.name}"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_cr_ee_repo.default.name}"`,
			"enable_details": "false",
		}),
	}

	namespaceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"namespace":      `"${alicloud_cr_ee_repo.default.namespace}"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"namespace":      `"${alicloud_cr_ee_repo.default.namespace}"`,
			"enable_details": "false",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cr_ee_repo.default.repo_id}"]`,
			"name_regex":     `"${alicloud_cr_ee_repo.default.name}"`,
			"namespace":      `"${alicloud_cr_ee_repo.default.namespace}"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAliCloudCrEEReposDataSourceName(rand, map[string]string{
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
			"repos.0.id":          acctest.CHECKSET,
			"repos.0.instance_id": acctest.CHECKSET,
			"repos.0.namespace":   acctest.CHECKSET,
			"repos.0.name":        acctest.CHECKSET,
			"repos.0.summary":     acctest.CHECKSET,
			"repos.0.repo_type":   acctest.CHECKSET,
			// Currently, Terraform doesn't have an environment testing ListRepoTag
			//"repos.0.tags.#":              "1",
			//"repos.0.tags.0.tag":          acctest.CHECKSET,
			//"repos.0.tags.0.image_id":     acctest.CHECKSET,
			//"repos.0.tags.0.image_size":   acctest.CHECKSET,
			//"repos.0.tags.0.digest":       acctest.CHECKSET,
			//"repos.0.tags.0.status":       acctest.CHECKSET,
			//"repos.0.tags.0.image_create": acctest.CHECKSET,
			//"repos.0.tags.0.image_update": acctest.CHECKSET,
		}
	}

	var fakeAliCloudCrEEReposDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"repos.#":             "1",
			"repos.0.id":          acctest.CHECKSET,
			"repos.0.instance_id": acctest.CHECKSET,
			"repos.0.namespace":   acctest.CHECKSET,
			"repos.0.name":        acctest.CHECKSET,
			"repos.0.summary":     acctest.CHECKSET,
			"repos.0.repo_type":   acctest.CHECKSET,
			"repos.0.tags.#":      "0",
		}
	}

	var alicloudCrEEReposCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cr_ee_repos.default",
		ExistMapFunc: existAliCloudCrEEReposDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudCrEEReposDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	alicloudCrEEReposCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, namespaceConf, allConf)
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
