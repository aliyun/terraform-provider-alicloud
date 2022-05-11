package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCTrafficMirrorFiltersDataSource(t *testing.T) {
	resourceId := "data.alicloud_vpc_traffic_mirror_filters.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	name := fmt.Sprintf("tf-testacc-vpctrafficmirrorfilter-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcTrafficMirrorFiltersDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_vpc_traffic_mirror_filter.default.traffic_mirror_filter_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_vpc_traffic_mirror_filter.default.traffic_mirror_filter_name}-fake",
		}),
	}
	filterNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"traffic_mirror_filter_name": "${alicloud_vpc_traffic_mirror_filter.default.traffic_mirror_filter_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"traffic_mirror_filter_name": "${alicloud_vpc_traffic_mirror_filter.default.traffic_mirror_filter_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpc_traffic_mirror_filter.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpc_traffic_mirror_filter.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_traffic_mirror_filter.default.id}"},
			"status": "Created",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_traffic_mirror_filter.default.id}"},
			"status": "Deleting",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_vpc_traffic_mirror_filter.default.traffic_mirror_filter_name}",
			"ids":        []string{"${alicloud_vpc_traffic_mirror_filter.default.id}"},
			"status":     "Created",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_vpc_traffic_mirror_filter.default.traffic_mirror_filter_name}-fake",
			"ids":        []string{"${alicloud_vpc_traffic_mirror_filter.default.id}"},
			"status":     "Deleting",
		}),
	}
	var existActiontrailTrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"ids.0":                                CHECKSET,
			"filters.#":                            "1",
			"filters.0.traffic_mirror_filter_name": fmt.Sprintf("tf-testacc-vpctrafficmirrorfilter-%d", rand),
			"filters.0.status":                     "Created",
		}
	}

	var fakeActiontrailTrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"filters.#": "0",
		}
	}

	var vpcTrafficMirrorFilterCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existActiontrailTrailMapFunc,
		fakeMapFunc:  fakeActiontrailTrailMapFunc,
	}

	vpcTrafficMirrorFilterCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, filterNameConf, idsConf, statusConf, allConf)
}

func dataSourceVpcTrafficMirrorFiltersDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
	  default = "%s"
	}

	resource "alicloud_vpc_traffic_mirror_filter" "default" {
	  traffic_mirror_filter_name = var.name
	  traffic_mirror_filter_description = var.name
	}`, name)
}
