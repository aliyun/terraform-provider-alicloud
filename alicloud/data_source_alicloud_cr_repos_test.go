package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudCRReposDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cr_repos.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-cr-repo-%d", rand),
		dataSourceCRReposConfigDependence)

	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace": alicloud_cr_repo.default.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace": "${alicloud_cr_repo.default.name}_fake",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": alicloud_cr_repo.default.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_cr_repo.default.name}_fake",
		}),
	}

	enableDetailsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     alicloud_cr_repo.default.name,
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace":      alicloud_cr_repo.default.name,
			"name_regex":     alicloud_cr_repo.default.name,
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace":      "${alicloud_cr_repo.default.name}_fake",
			"name_regex":     alicloud_cr_repo.default.name,
			"enable_details": "true",
		}),
	}

	var existCRReposMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":           "1",
			"names.0":           fmt.Sprintf("tf-testacc-cr-repo-%d", rand),
			"repos.#":           "1",
			"repos.0.name":      fmt.Sprintf("tf-testacc-cr-repo-%d", rand),
			"repos.0.namespace": fmt.Sprintf("tf-testacc-cr-repo-%d", rand),
			"repos.0.summary":   "OLD",
			"repos.0.repo_type": "PUBLIC",
			"repos.0.tags.#":    "0",
		}
	}

	var fakeCRReposMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#": "0",
			"repos.#": "0",
		}
	}

	var crReposCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCRReposMapFunc,
		fakeMapFunc:  fakeCRReposMapFunc,
	}

	crReposCheckInfo.dataSourceTestCheck(t, rand, namespaceConf, nameRegexConf, enableDetailsConf, allConf)
}

func dataSourceCRReposConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_namespace" "default" {
    name = var.name
    auto_create	= false
    default_visibility = "PUBLIC"
}

resource "alicloud_cr_repo" "default" {
    namespace = alicloud_cr_namespace.default.name
    name = var.name
    summary = "OLD"
    repo_type = "PUBLIC"
    detail  = "OLD"
}
`, name)
}
