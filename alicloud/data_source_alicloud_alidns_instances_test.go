package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func SkipTestAccAlicloudAlidnsInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_alidns_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAlidnsInstancesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_alidns_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_alidns_instance.default.id}-fake"},
		}),
	}

	var existAlidnsInstancesMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"ids.0":                      CHECKSET,
			"instances.#":                "1",
			"instances.0.dns_security":   "no",
			"instances.0.domain_numbers": "4",
			"instances.0.instance_id":    CHECKSET,
			"instances.0.version_code":   "version_personal",
			"instances.0.version_name":   CHECKSET,
		}
	}

	var fakeAlidnsInstancesMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var alidnsInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alidns_instances.default",
		existMapFunc: existAlidnsInstancesMapCheck,
		fakeMapFunc:  fakeAlidnsInstancesMapCheck,
	}

	alidnsInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceAlidnsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_alidns_instance" "default" {
   dns_security   = "no"
   domain_numbers = "4"
   period         = 1
   renew_period   = 1
   renewal_status = "ManualRenewal"
   version_code   = "version_personal"
}`)
}
