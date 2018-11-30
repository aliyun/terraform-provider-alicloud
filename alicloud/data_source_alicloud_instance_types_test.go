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
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.gpu.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.4c8g", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.local_storage.category", ""),
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
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.cpu_core_count"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.memory_size"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.family", "ecs.gn5"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.gpu.category"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.category"),
				),
			},
		},
	})
}

func TestAccAlicloudInstanceTypesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.empty"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.empty", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.gpu.%"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.local_storage.%"),
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

const testAccCheckAlicloudInstanceTypesDataSourceGpu = `
provider "alicloud" {
	region = "cn-hangzhou"
}
data "alicloud_instance_types" "gpu" {
	instance_type_family = "ecs.gn5"
	instance_charge_type = "PrePaid"
}
`

const testAccCheckAlicloudInstanceTypesDataSourceEmpty = `
data "alicloud_instance_types" "empty" {
	instance_type_family = "ecs.fake"
}
`
