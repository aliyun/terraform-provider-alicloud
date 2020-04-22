package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEdasApplicationsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_edas_applications.default"
	name := fmt.Sprintf("tf-testacc-edas-applications%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasApplicationConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_edas_application.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_edas_application.default.id}_fake"},
		}),
	}

	var existEdasApplicationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"applications.#":                  "1",
			"applications.0.app_name":         fmt.Sprintf("tf-testacc-edas-applications%v", rand),
			"applications.0.app_id":           CHECKSET,
			"applications.0.application_type": "War",
			"applications.0.build_package_id": CHECKSET,
		}
	}

	var fakeEdasApplicationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"applications.#": "0",
		}
	}

	var edasApplicationCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEdasApplicationsMapFunc,
		fakeMapFunc:  fakeEdasApplicationsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
	}

	edasApplicationCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func dataSourceEdasApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		resource "alicloud_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}

		resource "alicloud_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = "${alicloud_vpc.default.id}"
		}

		resource "alicloud_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = alicloud_edas_cluster.default.id
		  package_type = "WAR"
		}
		`, name)
}
