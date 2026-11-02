package alidns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAlidnsInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_dns_instances.default"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceAlidnsInstancesConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_instance.default.id}-fake"},
		}),
	}

	var existAlidnsInstancesMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"ids.0":                      acctest.CHECKSET,
			"instances.#":                "1",
			"instances.0.dns_security":   "no",
			"instances.0.domain_numbers": "4",
			"instances.0.instance_id":    acctest.CHECKSET,
			"instances.0.version_code":   "version_personal",
			"instances.0.version_name":   acctest.CHECKSET,
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

	var AlidnsInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_dns_instances.default",
		ExistMapFunc: existAlidnsInstancesMapCheck,
		FakeMapFunc:  fakeAlidnsInstancesMapCheck,
	}

	AlidnsInstancesCheckInfo.DataSourceTestCheck(t, rand, idsConf)
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithTime(t, []int{1})
	}
	AlidnsInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, 0, preCheck, idsConf)
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
