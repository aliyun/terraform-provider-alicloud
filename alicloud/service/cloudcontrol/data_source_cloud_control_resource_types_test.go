package cloudcontrol_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// CloudControl resource-type data sources are Hangzhou-only; the flat suite's
// testAccPreCheckWithRegions converged to cn-hangzhou, so pin it here.
func TestAccAlicloudCloudControlResourceTypeDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// ids filter, matching
			{
				Config: testAccCheckAlicloudCloudControlResourceTypeSourceConfig(rand, map[string]string{
					"product": `"VPC"`,
					"ids":     `["VSwitch"]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloud_control_resource_types.default", "types.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.primary_identifier"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.sensitive_info_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.public_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.update_type_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.list_response_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.list_only_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.product"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.get_only_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.create_only_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.properties"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.info.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.handlers.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.read_only_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.required.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.filter_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.get_response_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.update_only_properties.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_resource_types.default", "types.0.delete_only_properties.#"),
				),
			},
			// product filter, no match
			{
				Config: testAccCheckAlicloudCloudControlResourceTypeSourceConfig(rand, map[string]string{
					"product": `"VPC_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloud_control_resource_types.default", "types.#", "0"),
				),
			},
		},
	})
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
