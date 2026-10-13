package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAliCloudDMSUserTenantsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.DMSEnterpriseSupportRegions)
	resourceId := "data.alicloud_dms_user_tenants.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceDmsUserTenantsConfigDependence)

	existConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "ACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"999999999"},
		}),
	}
	var existDmsUserTenantsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 CHECKSET,
			"ids.0":                 CHECKSET,
			"tenants.#":             CHECKSET,
			"tenants.0.tenant_name": CHECKSET,
			"tenants.0.status":      "ACTIVE",
			"tenants.0.id":          CHECKSET,
			"tenants.0.tid":         CHECKSET,
		}
	}

	var fakeDmsUserTenantsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"tenants.#": "0",
		}
	}

	var userTenantsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDmsUserTenantsMapFunc,
		fakeMapFunc:  fakeDmsUserTenantsMapFunc,
	}

	userTenantsCheckInfo.dataSourceTestCheck(t, 0, existConf)
}

func dataSourceDmsUserTenantsConfigDependence(name string) string {
	return ""
}
