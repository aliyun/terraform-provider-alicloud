package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEsaCacheReserveInstancesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_cache_reserve_instances.default"
	name := fmt.Sprintf("tf-testAcc-EsaCacheReserveInstance%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaCacheReserveInstancesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}_fake"},
		}),
	}

	cacheReserveInstanceId := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cache_reserve_instance_id": "${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cache_reserve_instance_id": "${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}_fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}"},
			"status": "offline",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "overdue",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                       []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}"},
			"cache_reserve_instance_id": "${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}",
			"status":                    "offline",
			"sort_by":                   "CreateTime",
			"sort_order":                "asc",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                       []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}_fake"},
			"cache_reserve_instance_id": "${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}_fake",
			"status":                    "overdue",
			"sort_by":                   "CreateTime",
			"sort_order":                "asc",
		}),
	}

	var existAliCloudEsaCacheReserveInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"instances.#":                           "1",
			"instances.0.id":                        CHECKSET,
			"instances.0.cache_reserve_instance_id": CHECKSET,
			"instances.0.quota_gb":                  CHECKSET,
			"instances.0.cr_region":                 CHECKSET,
			"instances.0.payment_type":              CHECKSET,
			"instances.0.period":                    CHECKSET,
			"instances.0.status":                    CHECKSET,
			"instances.0.create_time":               CHECKSET,
			"instances.0.expire_time":               CHECKSET,
		}
	}

	var fakeAliCloudEsaCacheReserveInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var aliCloudEsaCacheReserveInstancesInfo = dataSourceAttr{
		resourceId:   "data.alicloud_esa_cache_reserve_instances.default",
		existMapFunc: existAliCloudEsaCacheReserveInstancesMapFunc,
		fakeMapFunc:  fakeAliCloudEsaCacheReserveInstancesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaCacheReserveInstancesInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, cacheReserveInstanceId, statusConf, allConf)
}

func dataSourceEsaCacheReserveInstancesConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_cache_reserve_instances" "test" {
}
`, name)
}
