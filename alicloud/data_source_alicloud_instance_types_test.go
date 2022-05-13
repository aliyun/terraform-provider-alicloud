package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECSInstanceTypesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.c4g8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.category", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.nvme_support"),
				),
			},
		},
	})
}

func TestAccAlicloudECSInstanceTypesDataSource_gpu(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.nvme_support"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "ids.#"),
				),
			},
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceGpuK8SMaster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.gpu"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.gpu.%"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "ids.#"),
				),
			},
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceGpuK8SWorker,
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
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.nvme_support"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "ids.#"),
				),
			},
		},
	})
}

func TestAccAlicloudECSInstanceTypesDataSource_empty(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.empty", "ids.#"),
				),
			},
		},
	})
}

func TestAccAlicloudECSInstanceTypesDataSource_k8sSpec(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceK8Sc1g2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.c1g2"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.gpu.%"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c1g2", "ids.#"),
				),
			},
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceK8Sc2g4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.c2g4"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.cpu_core_count", "2"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.memory_size", "4"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.gpu.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.local_storage.category", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "ids.#"),
				),
			},
		},
	})
}
func TestAccAlicloudECSInstanceTypesDataSource_k8sFamily(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceK8ST5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.t5"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.t5", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.gpu.%"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.burstable_instance.%"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.local_storage.%"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.t5", "ids.#"),
				),
			},
		},
	})
}
func TestAccAlicloudECSInstanceTypesDataSource_imageId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceImageId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.c4g8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.category", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlicloudInstanceTypesDataSourceBasicConfig = `
data "alicloud_instance_types" "c4g8" {
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
const testAccCheckAlicloudInstanceTypesDataSourceGpuK8SMaster = `
provider "alicloud" {
	region = "cn-hangzhou"
}
data "alicloud_instance_types" "gpu" {
	kubernetes_node_role = "Master"
	instance_type_family = "ecs.gn5"
}
`
const testAccCheckAlicloudInstanceTypesDataSourceGpuK8SWorker = `
provider "alicloud" {
	region = "cn-hangzhou"
}
data "alicloud_instance_types" "gpu" {
	kubernetes_node_role = "Worker"
	instance_type_family = "ecs.gn5"
}
`

const testAccCheckAlicloudInstanceTypesDataSourceEmpty = `
data "alicloud_instance_types" "empty" {
	instance_type_family = "ecs.fake"
}
`

const testAccCheckAlicloudInstanceTypesDataSourceK8Sc1g2 = `
data "alicloud_instance_types" "c1g2" {
	cpu_core_count = 1
	memory_size = 2
	kubernetes_node_role = "Master"
}
`
const testAccCheckAlicloudInstanceTypesDataSourceK8Sc2g4 = `
data "alicloud_instance_types" "c2g4" {
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}
`
const testAccCheckAlicloudInstanceTypesDataSourceK8ST5 = `
data "alicloud_instance_types" "t5" {
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Master"
	instance_type_family = "ecs.t5"
}
`
const testAccCheckAlicloudInstanceTypesDataSourceImageId = `
data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "c4g8" {
	image_id = data.alicloud_images.default.ids.0
	cpu_core_count = 4
	memory_size = 8
}
`
