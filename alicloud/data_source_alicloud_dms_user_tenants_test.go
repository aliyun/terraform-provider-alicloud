package alicloud

import (
	"testing"
)

func TestAccAlicloudDMSUserTenantsDataSource(t *testing.T) {
	resourceId := "data.alicloud_dms_user_tenants.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceDmsUserTenantsConfigDependence)

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "ACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "IN_ACTIVE",
		}),
	}
	var existDmsUserTenantsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"ids.0":                 CHECKSET,
			"tenants.#":             "1",
			"tenants.0.tenant_name": CHECKSET,
			"tenants.0.status":      "ACTIVE",
			"tenants.0.id":          CHECKSET,
			"tenants.0.tid":         CHECKSET,
		}
	}

	var fakeDmsUserTenantsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"tenants.#": "0",
		}
	}

	var userTenantsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDmsUserTenantsMapFunc,
		fakeMapFunc:  fakeDmsUserTenantsMapFunc,
	}

	userTenantsCheckInfo.dataSourceTestCheck(t, 0, statusConf)
}

func dataSourceDmsUserTenantsConfigDependence(name string) string {
	return ""
}
