package ecs_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudECSInstanceTypesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.c4g8"),
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
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.0.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.category", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.nvme_support"),
				),
			},
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceBasicConfigSortedByPrice,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.c4g8"),
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
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.0.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.category", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.nvme_support"),
				),
			},
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceBasicConfigSortedByPriceWithZone,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.zone"),
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
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.gpu.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.gpu.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.gpu.0.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.burstable_instance.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.local_storage.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.zone", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.zone", "instance_types.0.local_storage.0.category", ""),
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
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceGpu,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.gpu"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.cpu_core_count"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.memory_size"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.family", "ecs.gn5"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.gpu.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.gpu.0.amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.gpu.0.category"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.0.category"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.nvme_support"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "ids.#"),
				),
			},
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceGpuK8SMaster,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.gpu"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.gpu.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "ids.#"),
				),
			},
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceGpuK8SWorker,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.gpu"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.cpu_core_count"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.memory_size"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.family", "ecs.gn5"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.gpu.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.gpu.0.amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.gpu.0.category"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.gpu", "instance_types.0.local_storage.0.category"),
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
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.empty"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.empty", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.gpu.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.burstable_instance.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.empty", "instance_types.0.local_storage.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.empty", "ids.#"),
				),
			},
		},
	})
}

func TestAccAliCloudECSInstanceTypesDataSource_k8sSpec(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceK8Sc1g2,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.c1g2"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.gpu.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.burstable_instance.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.c1g2", "instance_types.0.local_storage.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c1g2", "ids.#"),
				),
			},
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceK8Sc2g4,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.c2g4"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.cpu_core_count", "2"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.memory_size", "4"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.gpu.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.gpu.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.gpu.0.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.burstable_instance.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.local_storage.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c2g4", "instance_types.0.local_storage.0.category", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c2g4", "ids.#"),
				),
			},
		},
	})
}

func TestAccAliCloudECSInstanceTypesDataSource_k8sFamily(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceK8ST5,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.t5"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.t5", "instance_types.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.cpu_core_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.family"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.eni_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.availability_zones.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.gpu.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.burstable_instance.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instance_types.t5", "instance_types.0.local_storage.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.t5", "ids.#"),
				),
			},
		},
	})
}

func TestAccAliCloudECSInstanceTypesDataSource_imageId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceImageId,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.c4g8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.family"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.eni_amount"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.availability_zones.#"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.gpu.0.category", ""),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.0.initial_credit"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.burstable_instance.0.baseline_credit"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.amount"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.c4g8", "instance_types.0.local_storage.0.category", ""),
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
  sorted_by      = "Price"
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
  sorted_by = "Price"
}
`

func TestAccAliCloudECSInstanceTypesDataSource_sortedBy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceSortedByConfig("CPU"),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.sorted"),
					testAccCheckInstanceTypesAscending("data.alicloud_instance_types.sorted", "cpu_core_count"),
				),
			},
			{
				Config: testAccCheckAliCloudInstanceTypesDataSourceSortedByConfig("Memory"),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_instance_types.sorted"),
					testAccCheckInstanceTypesAscending("data.alicloud_instance_types.sorted", "memory_size"),
				),
			},
		},
	})
}

func testAccCheckAliCloudInstanceTypesDataSourceSortedByConfig(sortedBy string) string {
	return fmt.Sprintf(`
data "alicloud_instance_types" "sorted" {
  instance_type_family = "ecs.g6"
  sorted_by            = "%s"
}
`, sortedBy)
}

func testAccCheckInstanceTypesAscending(name string, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("data source %s not found in state", name)
		}
		count, err := strconv.Atoi(rs.Primary.Attributes["instance_types.#"])
		if err != nil {
			return fmt.Errorf("reading instance_types.# got an error: %v", err)
		}
		if count < 2 {
			return fmt.Errorf("expected multiple instance types sorted by %s, got %d", attr, count)
		}
		parse := func(idx int) (float64, error) {
			return strconv.ParseFloat(rs.Primary.Attributes[fmt.Sprintf("instance_types.%d.%s", idx, attr)], 64)
		}
		prev, err := parse(0)
		if err != nil {
			return fmt.Errorf("reading instance_types.0.%s got an error: %v", attr, err)
		}
		for i := 1; i < count; i++ {
			cur, err := parse(i)
			if err != nil {
				return fmt.Errorf("reading instance_types.%d.%s got an error: %v", i, attr, err)
			}
			if cur < prev {
				return fmt.Errorf("instance_types not ascending by %s: index %d (%v) < index %d (%v)", attr, i, cur, i-1, prev)
			}
			prev = cur
		}
		return nil
	}
}
