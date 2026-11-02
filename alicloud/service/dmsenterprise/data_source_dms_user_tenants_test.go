package dmsenterprise_test

import (
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAliCloudDMSUserTenantsDataSource(t *testing.T) {
	tfacctest.CheckoutSupportedRegions(t, true, connectivity.DMSEnterpriseSupportRegions)
	resourceId := "data.alicloud_dms_user_tenants.default"
	testAccConfig := tfacctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceDmsUserTenantsConfigDependence)

	existConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"status": "ACTIVE",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"999999999"},
		}),
	}
	var existDmsUserTenantsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 tfacctest.CHECKSET,
			"ids.0":                 tfacctest.CHECKSET,
			"tenants.#":             tfacctest.CHECKSET,
			"tenants.0.tenant_name": tfacctest.CHECKSET,
			"tenants.0.status":      "ACTIVE",
			"tenants.0.id":          tfacctest.CHECKSET,
			"tenants.0.tid":         tfacctest.CHECKSET,
		}
	}

	var fakeDmsUserTenantsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"tenants.#": "0",
		}
	}

	var userTenantsCheckInfo = tfacctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDmsUserTenantsMapFunc,
		FakeMapFunc:  fakeDmsUserTenantsMapFunc,
	}

	userTenantsCheckInfo.DataSourceTestCheck(t, 0, existConf)
}

func dataSourceDmsUserTenantsConfigDependence(name string) string {
	return ""
}
