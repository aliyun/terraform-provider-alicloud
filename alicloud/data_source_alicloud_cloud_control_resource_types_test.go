package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudControlResourceTypeDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudControlResourceTypeSourceConfig(rand, map[string]string{
			"product": `"VPC"`,
			"ids":     `["VSwitch"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudControlResourceTypeSourceConfig(rand, map[string]string{
			"product": `"VPC_fake"`,
		}),
	}
	CloudControlResourceTypeCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existCloudControlResourceTypeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"types.#":                             "1",
		"types.0.primary_identifier":          CHECKSET,
		"types.0.sensitive_info_properties.#": CHECKSET,
		"types.0.public_properties.#":         CHECKSET,
		"types.0.update_type_properties.#":    CHECKSET,
		"types.0.list_response_properties.#":  CHECKSET,
		"types.0.list_only_properties.#":      CHECKSET,
		"types.0.product":                     CHECKSET,
		"types.0.get_only_properties.#":       CHECKSET,
		"types.0.resource_type":               CHECKSET,
		"types.0.create_only_properties.#":    CHECKSET,
		"types.0.properties":                  CHECKSET,
		"types.0.info.#":                      CHECKSET,
		"types.0.handlers.#":                  CHECKSET,
		"types.0.read_only_properties.#":      CHECKSET,
		"types.0.required.#":                  CHECKSET,
		"types.0.filter_properties.#":         CHECKSET,
		"types.0.get_response_properties.#":   CHECKSET,
		"types.0.update_only_properties.#":    CHECKSET,
		"types.0.delete_only_properties.#":    CHECKSET,
	}
}

var fakeCloudControlResourceTypeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"types.#": "0",
	}
}

var CloudControlResourceTypeCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_control_resource_types.default",
	existMapFunc: existCloudControlResourceTypeMapFunc,
	fakeMapFunc:  fakeCloudControlResourceTypeMapFunc,
}

func testAccCheckAlicloudCloudControlResourceTypeSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudControlResourceType%d"
}

data "alicloud_cloud_control_resource_types" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
