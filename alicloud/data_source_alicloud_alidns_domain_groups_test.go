package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsDomainGroupDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_alidns_domain_groups.default"
	name := fmt.Sprintf("tf-testAccadnsDG%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlidnsDomainGroupConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_alidns_domain_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_alidns_domain_group.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_alidns_domain_group.default.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_alidns_domain_group.default.group_name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_alidns_domain_group.default.id}"},
			"name_regex": "${alicloud_alidns_domain_group.default.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_alidns_domain_group.default.id}"},
			"name_regex": "${alicloud_alidns_domain_group.default.group_name}-fake",
		}),
	}

	var existAlidnsDomainGroupMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"ids.0":                 CHECKSET,
			"groups.#":              "1",
			"groups.0.group_id":     CHECKSET,
			"groups.0.group_name":   name,
			"groups.0.domain_count": CHECKSET,
		}
	}

	var fakeAlidnsDomainGroupMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}

	var alidnsDomainGroupCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alidns_domain_groups.default",
		existMapFunc: existAlidnsDomainGroupMapCheck,
		fakeMapFunc:  fakeAlidnsDomainGroupMapCheck,
	}

	alidnsDomainGroupCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceAlidnsDomainGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_alidns_domain_group" "default" {
  group_name = "%s"
}`, name)
}
