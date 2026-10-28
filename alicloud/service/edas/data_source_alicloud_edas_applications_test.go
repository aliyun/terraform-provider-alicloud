package edas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEdasApplicationsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_edas_applications.default"
	name := fmt.Sprintf("tf-testacc-edas-applications%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasApplicationConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_edas_application.default.application_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_edas_application.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_edas_application.default.id}_fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_edas_application.default.id}"},
			"name_regex": "${alicloud_edas_application.default.application_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_edas_application.default.id}_fake"},
			"name_regex": "${alicloud_edas_application.default.application_name}",
		}),
	}

	var existEdasApplicationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"applications.#":                  "1",
			"applications.0.app_name":         name,
			"applications.0.app_id":           acctest.CHECKSET,
			"applications.0.application_type": "War",
			"applications.0.build_package_id": acctest.CHECKSET,
		}
	}

	var fakeEdasApplicationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"applications.#": "0",
		}
	}

	var edasApplicationCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existEdasApplicationsMapFunc,
		FakeMapFunc:  fakeEdasApplicationsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
	}

	edasApplicationCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceEdasApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		data "alicloud_vpcs" "default" {
			name_regex = "^default-NODELETING$"
		}

		resource "alicloud_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = data.alicloud_vpcs.default.ids.0
		}

		resource "alicloud_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = alicloud_edas_cluster.default.id
		  package_type = "WAR"
		}
		`, name)
}
