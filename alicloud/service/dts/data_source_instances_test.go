package dts_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const dtsInstancesDataSourceId = "data.alicloud_dts_instances.default"

func TestAccAlicloudDtsInstanceDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// idsConf
			{
				Config: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
					"ids": `["${alicloud_dts_instance.default.id}"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsInstanceExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
					"ids": `["${alicloud_dts_instance.default.id}_fake"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsInstanceFakeChecks()...),
			},
			// nameRegexConf
			{
				Config: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
					"name_regex": `"${alicloud_dts_instance.default.instance_name}"`,
					"ids":        `["${alicloud_dts_instance.default.id}"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsInstanceExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
					"ids":        `["${alicloud_dts_instance.default.id}"]`,
					"name_regex": `"${alicloud_dts_instance.default.instance_name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsInstanceFakeChecks()...),
			},
			// ResourceGroupIdConf
			{
				Config: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
					"ids":               `["${alicloud_dts_instance.default.id}"]`,
					"resource_group_id": `"${alicloud_dts_instance.default.resource_group_id}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsInstanceExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
					"ids":               `["${alicloud_dts_instance.default.id}_fake"]`,
					"resource_group_id": `"${alicloud_dts_instance.default.resource_group_id}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsInstanceFakeChecks()...),
			},
			// allConf
			{
				Config: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
					"ids":               `["${alicloud_dts_instance.default.id}"]`,
					"resource_group_id": `"${alicloud_dts_instance.default.resource_group_id}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsInstanceExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
					"ids":               `["${alicloud_dts_instance.default.id}_fake"]`,
					"resource_group_id": `"${alicloud_dts_instance.default.resource_group_id}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsInstanceFakeChecks()...),
			},
		},
	})
}

func dtsInstanceExistChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "ids.#", "1"),
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "names.#", "1"),
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "instances.#", "1"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.id"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.create_time"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.destination_endpoint_engine_name"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.dts_instance_id"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.instance_name"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.instance_class"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.payment_type"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.resource_group_id"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.source_endpoint_engine_name"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.source_region"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.destination_region"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.status"),
		resource.TestCheckResourceAttrSet(dtsInstancesDataSourceId, "instances.0.type"),
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "instances.0.tags.%", "2"),
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "instances.0.tags.Create", "TF"),
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "instances.0.tags.For", "acceptance test"),
	}
}

func dtsInstanceFakeChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "ids.#", "0"),
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "names.#", "0"),
		resource.TestCheckResourceAttr(dtsInstancesDataSourceId, "instances.#", "0"),
	}
}

func testAccCheckAlicloudDtsInstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDtsInstance%d"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_dts_instance" "default" {
  type                             = "sync"
  resource_group_id                = data.alicloud_resource_manager_resource_groups.default.ids.0
  payment_type                     = "PayAsYouGo"
  instance_class                   = "large"
  source_endpoint_engine_name      = "MySQL"
  source_region                    = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_region               = data.alicloud_regions.default.regions.0.id
  tags = {
	  Create = "TF"
	  For = "acceptance test",
	}
}

data "alicloud_dts_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
