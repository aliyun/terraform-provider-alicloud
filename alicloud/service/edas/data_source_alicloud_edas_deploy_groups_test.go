package edas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEdasDeployGroupDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_edas_deploy_groups.default"
	name := fmt.Sprintf("tf-testacc-edas-deploy-groups%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasDeployGroupConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_edas_deploy_group.default.group_name}",
			"app_id":     "${alicloud_edas_application.default.id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
			"app_id":     "${alicloud_edas_application.default.id}",
		}),
	}
	var existEdasDeployGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":             "1",
			"groups.0.group_name":  name,
			"groups.0.app_id":      acctest.CHECKSET,
			"groups.0.group_type":  acctest.CHECKSET,
			"groups.0.cluster_id":  acctest.CHECKSET,
			"groups.0.create_time": acctest.CHECKSET,
			"groups.0.update_time": acctest.CHECKSET,
			"groups.0.group_id":    acctest.CHECKSET,
		}
	}

	var fakeEdasDeployGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
		}
	}

	var edasApplicationCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existEdasDeployGroupsMapFunc,
		FakeMapFunc:  fakeEdasDeployGroupsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
	}

	edasApplicationCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceEdasDeployGroupConfigDependence(name string) string {
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

		resource "alicloud_edas_deploy_group" "default" {
		  app_id = alicloud_edas_application.default.id
		  group_name = "${var.name}"
		}
		`, name)
}
