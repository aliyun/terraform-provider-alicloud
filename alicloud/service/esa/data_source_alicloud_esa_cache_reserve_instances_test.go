package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEsaCacheReserveInstancesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_cache_reserve_instances.default"
	name := fmt.Sprintf("tf-testAcc-EsaCacheReserveInstance%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaCacheReserveInstancesConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}_fake"},
		}),
	}

	cacheReserveInstanceId := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cache_reserve_instance_id": "${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cache_reserve_instance_id": "${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}_fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}"},
			"status": "offline",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"status": "overdue",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":                       []string{"${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}"},
			"cache_reserve_instance_id": "${data.alicloud_esa_cache_reserve_instances.test.instances.0.id}",
			"status":                    "offline",
			"sort_by":                   "CreateTime",
			"sort_order":                "asc",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"instances.0.id":                        acctest.CHECKSET,
			"instances.0.cache_reserve_instance_id": acctest.CHECKSET,
			"instances.0.quota_gb":                  acctest.CHECKSET,
			"instances.0.cr_region":                 acctest.CHECKSET,
			"instances.0.payment_type":              acctest.CHECKSET,
			"instances.0.period":                    acctest.CHECKSET,
			"instances.0.status":                    acctest.CHECKSET,
			"instances.0.create_time":               acctest.CHECKSET,
			"instances.0.expire_time":               acctest.CHECKSET,
		}
	}

	var fakeAliCloudEsaCacheReserveInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var aliCloudEsaCacheReserveInstancesInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_esa_cache_reserve_instances.default",
		ExistMapFunc: existAliCloudEsaCacheReserveInstancesMapFunc,
		FakeMapFunc:  fakeAliCloudEsaCacheReserveInstancesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaCacheReserveInstancesInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, cacheReserveInstanceId, statusConf, allConf)
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
