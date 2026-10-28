package brainindustrial_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudBrainIndustrialPidOrganizationsDataSource(t *testing.T) {
	resourceId := "data.alicloud_brain_industrial_pid_organizations.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceBrainIndustrialPidOrganizationsDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_brain_industrial_pid_organization.default.pid_organization_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_brain_industrial_pid_organization.default.pid_organization_name}-fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_brain_industrial_pid_organization.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_brain_industrial_pid_organization.default.id}-fake"},
		}),
	}
	var existBrainIndustrialPidOrganizationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "1",
			"ids.0":           acctest.CHECKSET,
			"names.#":         "1",
			"names.0":         name,
			"organizations.#": "1",
			"organizations.0.parent_pid_organization_id": "",
			"organizations.0.id":                         acctest.CHECKSET,
			"organizations.0.pid_organization_id":        acctest.CHECKSET,
			"organizations.0.pid_organization_level":     acctest.CHECKSET,
			"organizations.0.pid_organization_name":      name,
		}
	}

	var fakeBrainIndustrialPidOrganizationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"names.#":         "0",
			"organizations.#": "0",
		}
	}

	var BrainIndustrialPidOrganizationsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existBrainIndustrialPidOrganizationsMapFunc,
		FakeMapFunc:  fakeBrainIndustrialPidOrganizationsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
	}

	BrainIndustrialPidOrganizationsInfo.DataSourceTestCheckWithPreCheck(t, 0, preCheck, nameRegexConf, idsConf)
}

func dataSourceBrainIndustrialPidOrganizationsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_brain_industrial_pid_organization" "default" {
		pid_organization_name = "%s"
	}`, name)
}
