package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudInstanceTypesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.4c8g"),

					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.memory_size", "8"),
				),
			},

			resource.TestStep{
				Config: testAccCheckAlicloudInstanceTypesDataSourceBasicConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.4c8g"),

					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.#", "1"),

					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.memory_size", "8"),
				),
			},
		},
	})
}

func TestAccAlicloudInstanceTypesDataSource_gpu(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceGpu,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.gpu"),

					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.id", "ecs.gn5-c4g1.xlarge"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.gpu.amount", "1"),
				),
			},
		},
	})
}

const testAccCheckAlicloudInstanceTypesDataSourceBasicConfig = `
data "alicloud_instance_types" "4c8g" {
	cpu_core_count = 4
	memory_size = 8
}
`

const testAccCheckAlicloudInstanceTypesDataSourceBasicConfigUpdate = `
data "alicloud_instance_types" "4c8g" {
	instance_type_family= "ecs.c4"
	cpu_core_count = 4
	memory_size = 8
}
`

const testAccCheckAlicloudInstanceTypesDataSourceGpu = `
provider "alicloud" {
	region = "cn-hangzhou"
}
data "alicloud_instance_types" "gpu" {
	instance_type_family = "ecs.gn5"
	instance_charge_type = "PrePaid"
	cpu_core_count = 4
	memory_size = 30
}
`
