package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_dns_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAlidnsInstancesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_instance.default.id}-fake"},
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
			"instances.0.payment_type":   "Subscription",
			"instances.0.domain":         "",
		}
	}

	var fakeAlidnsInstancesMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var AlidnsInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dns_instances.default",
		existMapFunc: existAlidnsInstancesMapCheck,
		fakeMapFunc:  fakeAlidnsInstancesMapCheck,
	}

	AlidnsInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
	}
	AlidnsInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, idsConf)
}

func dataSourceAlidnsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns_instance" "default" {
   dns_security   = "no"
   domain_numbers = "4"
   period         = 1
   renew_period   = 1
   renewal_status = "ManualRenewal"
   version_code   = "version_personal"
}`)
}
