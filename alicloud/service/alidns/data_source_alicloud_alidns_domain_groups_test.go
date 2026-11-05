package alidns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAlidnsDomainGroupDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_alidns_domain_groups.default"
	name := fmt.Sprintf("tf-testAccadnsDG%v", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceAlidnsDomainGroupConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_alidns_domain_group.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_alidns_domain_group.default.id}-fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_alidns_domain_group.default.group_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_alidns_domain_group.default.group_name}-fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_alidns_domain_group.default.id}"},
			"name_regex": "${alicloud_alidns_domain_group.default.group_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_alidns_domain_group.default.id}"},
			"name_regex": "${alicloud_alidns_domain_group.default.group_name}-fake",
		}),
	}

	var existAlidnsDomainGroupMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"ids.0":                 acctest.CHECKSET,
			"groups.#":              "1",
			"groups.0.group_id":     acctest.CHECKSET,
			"groups.0.group_name":   name,
			"groups.0.domain_count": acctest.CHECKSET,
		}
	}

	var fakeAlidnsDomainGroupMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}

	var alidnsDomainGroupCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_alidns_domain_groups.default",
		ExistMapFunc: existAlidnsDomainGroupMapCheck,
		FakeMapFunc:  fakeAlidnsDomainGroupMapCheck,
	}

	alidnsDomainGroupCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceAlidnsDomainGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_alidns_domain_group" "default" {
  group_name = "%s"
}`, name)
}
