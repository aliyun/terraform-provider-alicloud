package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudECSInstanceTypesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.c4g8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.primary_eni_queue_number"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.secondary_eni_queue_number"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_ipv6_address_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.maximum_queue_number_per_eni"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.total_eni_queue_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_private_ip_address_quantity"),
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
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceBasicConfigSortedByPrice,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.c4g8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.primary_eni_queue_number"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.secondary_eni_queue_number"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_ipv6_address_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.maximum_queue_number_per_eni"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.total_eni_queue_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_private_ip_address_quantity"),
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
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceBasicConfigSortedByPriceWithZone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instance_types.zone"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.cpu_core_count"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.memory_size"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.eni_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.primary_eni_queue_number"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.secondary_eni_queue_number"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.eni_ipv6_address_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.maximum_queue_number_per_eni"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.total_eni_queue_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.eni_private_ip_address_quantity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.gpu.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.gpu.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.gpu.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.burstable_instance.%", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.burstable_instance.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.burstable_instance.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.local_storage.%", "3"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.local_storage.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.local_storage.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.local_storage.category", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.nvme_support"),
				),
			},
		},
	})
}

func TestAccAliCloudECSInstanceTypesDataSource_gpu(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceGpu,
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
				Config: testAccCheckAliCloudInstanceTypesDataSourceGpuK8SMaster,
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
				Config: testAccCheckAliCloudInstanceTypesDataSourceGpuK8SWorker,
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

func TestAccAliCloudECSInstanceTypesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceEmpty,
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

func TestAccAliCloudECSInstanceTypesDataSource_k8sSpec(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceK8Sc1g2,
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
				Config: testAccCheckAliCloudInstanceTypesDataSourceK8Sc2g4,
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

func TestAccAliCloudECSInstanceTypesDataSource_k8sFamily(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceK8ST5,
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

func TestAccAliCloudECSInstanceTypesDataSource_imageId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceImageId,
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

const testAccCheckAliCloudInstanceTypesDataSourceBasicConfig = `
data "alicloud_instance_types" "c4g8" {
  cpu_core_count = 4
  memory_size    = 8
}
`

const testAccCheckAlicloudInstanceTypesDataSourceBasicConfigSortedByPrice = `
data "alicloud_instance_types" "c4g8" {
  cpu_core_count = 4
  memory_size    = 8
  sorted_by    = "Price"
}
`

const testAccCheckAliCloudInstanceTypesDataSourceGpu = `
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_instance_types" "gpu" {
  instance_type_family = "ecs.gn5"
  instance_charge_type = "PrePaid"
}
`

const testAccCheckAliCloudInstanceTypesDataSourceGpuK8SMaster = `
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_instance_types" "gpu" {
  kubernetes_node_role = "Master"
  instance_type_family = "ecs.gn5"
}
`

const testAccCheckAliCloudInstanceTypesDataSourceGpuK8SWorker = `
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_instance_types" "gpu" {
  kubernetes_node_role = "Worker"
  instance_type_family = "ecs.gn5"
}
`

const testAccCheckAliCloudInstanceTypesDataSourceEmpty = `
data "alicloud_instance_types" "empty" {
  instance_type_family = "ecs.fake"
}
`

const testAccCheckAliCloudInstanceTypesDataSourceK8Sc1g2 = `
data "alicloud_instance_types" "c1g2" {
  cpu_core_count       = 1
  memory_size          = 2
  kubernetes_node_role = "Master"
}
`

const testAccCheckAliCloudInstanceTypesDataSourceK8Sc2g4 = `
data "alicloud_instance_types" "c2g4" {
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}
`

const testAccCheckAliCloudInstanceTypesDataSourceK8ST5 = `
data "alicloud_instance_types" "t5" {
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Master"
  instance_type_family = "ecs.t5"
}
`

const testAccCheckAliCloudInstanceTypesDataSourceImageId = `
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "c4g8" {
  image_id       = data.alicloud_images.default.ids.0
  cpu_core_count = 4
  memory_size    = 8
}
`

const testAccCheckAliCloudInstanceTypesDataSourceBasicConfigSortedByPriceWithZone = `
data "alicloud_instance_types" "zone" {
  sorted_by                         = "Price"
}
`
