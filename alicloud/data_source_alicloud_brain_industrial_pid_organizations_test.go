package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudBrainIndustrialPidOrganizationsDataSource(t *testing.T) {
	resourceId := "data.alicloud_brain_industrial_pid_organizations.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceBrainIndustrialPidOrganizationsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_brain_industrial_pid_organization.default.pid_organization_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_brain_industrial_pid_organization.default.pid_organization_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_brain_industrial_pid_organization.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_brain_industrial_pid_organization.default.id}-fake"},
		}),
	}
	var existBrainIndustrialPidOrganizationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "1",
			"ids.0":           CHECKSET,
			"names.#":         "1",
			"names.0":         name,
			"organizations.#": "1",
			"organizations.0.parent_pid_organization_id": "",
			"organizations.0.id":                         CHECKSET,
			"organizations.0.pid_organization_id":        CHECKSET,
			"organizations.0.pid_organization_level":     CHECKSET,
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

	var BrainIndustrialPidOrganizationsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existBrainIndustrialPidOrganizationsMapFunc,
		fakeMapFunc:  fakeBrainIndustrialPidOrganizationsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
	}

	BrainIndustrialPidOrganizationsInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, nameRegexConf, idsConf)
}

func dataSourceBrainIndustrialPidOrganizationsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_brain_industrial_pid_organization" "default" {
		pid_organization_name = "%s"
	}`, name)
}
