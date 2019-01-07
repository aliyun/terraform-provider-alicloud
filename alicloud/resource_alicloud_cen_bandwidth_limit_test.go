package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCenBandwidthLimit_basic(t *testing.T) {
	var cenBwpLimit cbn.CenInterRegionBandwidthLimit

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenBandwidthLimitDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthLimitConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExistsWithProviders("alicloud_cen_bandwidth_limit.foo", &cenBwpLimit, &providers),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.foo", "bandwidth_limit", "4"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.foo", "region_ids.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "eu-central-1", "cn-shanghai"),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthLimit_update(t *testing.T) {
	var cenBwpLimit cbn.CenInterRegionBandwidthLimit

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenBandwidthLimitDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthLimitConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExistsWithProviders("alicloud_cen_bandwidth_limit.foo", &cenBwpLimit, &providers),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.foo", "bandwidth_limit", "4"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.foo", "region_ids.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "eu-central-1", "cn-shanghai"),
				),
			},
			resource.TestStep{
				Config: testAccCenBandwidthLimitUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExistsWithProviders("alicloud_cen_bandwidth_limit.foo", &cenBwpLimit, &providers),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.foo", "bandwidth_limit", "5"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.foo", "region_ids.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "eu-central-1", "cn-shanghai"),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthLimit_multi(t *testing.T) {
	var cenBwpLimit cbn.CenInterRegionBandwidthLimit

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenBandwidthLimitDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthLimitMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExistsWithProviders("alicloud_cen_bandwidth_limit.bar1", &cenBwpLimit, &providers),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.bar1", "bandwidth_limit", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.bar1", "region_ids.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "eu-central-1", "cn-shanghai"),

					testAccCheckCenBandwidthLimitExistsWithProviders("alicloud_cen_bandwidth_limit.bar2", &cenBwpLimit, &providers),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.bar2", "bandwidth_limit", "3"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_limit.bar2", "region_ids.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "eu-central-1", "cn-hangzhou"),
				),
			},
		},
	})
}

func testAccCheckCenBandwidthLimitExistsWithProviders(n string, cenBwpLimit *cbn.CenInterRegionBandwidthLimit, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CEN bandwidth limit ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			cenService := CenService{client}

			params, err := cenService.GetCenAndRegionIds(rs.Primary.ID)
			if err != nil {
				return err
			}
			cenId := params[0]
			localRegionId := params[1]
			oppositeRegionId := params[2]
			instance, err := cenService.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)
			if err != nil {
				return err
			}

			*cenBwpLimit = instance
			return nil
		}
		return fmt.Errorf("Cen bandwidth not found")
	}
}

func testAccCheckCenBandwidthLimitDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenBandwidthLimitDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenBandwidthLimitDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	cenService := CenService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_bandwidth_limit" {
			continue
		}

		params, err := cenService.GetCenAndRegionIds(rs.Primary.ID)
		if err != nil {
			return err
		}
		cenId := params[0]
		localRegionId := params[1]
		oppositeRegionId := params[2]

		instance, err := cenService.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("CEN Bandwidth Limit still exist, CEN ID %s localRegionId %s oppositeRegionId %s",
				instance.CenId, instance.LocalRegionId, instance.OppositeRegionId)
		}
	}

	return nil
}

func testAccCheckCenBandwidthLimitRegionId(cenBwpLimit *cbn.CenInterRegionBandwidthLimit, regionAId string, regionBId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if (cenBwpLimit.LocalRegionId == regionAId && cenBwpLimit.OppositeRegionId == regionBId) ||
			(cenBwpLimit.LocalRegionId == regionBId && cenBwpLimit.OppositeRegionId == regionAId) {
			return nil
		} else {
			return fmt.Errorf("CEN %s BandwidthLimit Region ID error", cenBwpLimit.CenId)
		}
	}
}

const testAccCenBandwidthLimitConfig = `
variable "name"{
    default = "tf-testAccCenBandwidthLimitConfig"
}

provider "alicloud" {
    alias = "fra"
    region = "eu-central-1"
}

provider "alicloud" {
    alias = "sh"
    region = "cn-shanghai"
}

resource "alicloud_vpc" "vpc1" {
  provider = "alicloud.fra"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  provider = "alicloud.sh"
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "cen" {
     name = "${var.name}"
     description = "tf-testAccCenBandwidthLimitConfigDescription"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
    bandwidth = 5
    geographic_region_ids = [
		"Europe",
		"China"]
}

resource "alicloud_cen_bandwidth_package_attachment" "bwp_attach" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidth_limit" "foo" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    region_ids = [
        "eu-central-1",
        "cn-shanghai"]
     bandwidth_limit = 4
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}
`

const testAccCenBandwidthLimitUpdate = `
variable "name"{
    default = "tf-testAccCenBandwidthLimitConfig"
}

provider "alicloud" {
    alias = "fra"
    region = "eu-central-1"
}

provider "alicloud" {
    alias = "sh"
    region = "cn-shanghai"
}

resource "alicloud_vpc" "vpc1" {
  provider = "alicloud.fra"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  provider = "alicloud.sh"
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "cen" {
     name = "${var.name}"
     description = "tf-testAccCenBandwidthLimitConfigDescription"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
    bandwidth = 5
    geographic_region_ids = [
		"Europe",
		"China"]
}

resource "alicloud_cen_bandwidth_package_attachment" "bwp_attach" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidth_limit" "foo" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    region_ids = [
        "eu-central-1",
        "cn-shanghai"]
     bandwidth_limit = 5
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}
`

const testAccCenBandwidthLimitMulti = `
variable "name"{
    default = "tf-testAccCenBandwidthLimitMulti"
}

provider "alicloud" {
    alias = "fra"
    region = "eu-central-1"
}

provider "alicloud" {
  alias = "sh"
  region = "cn-shanghai"
}

provider "alicloud" {
  alias = "hz"
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "vpc1" {
  provider = "alicloud.fra"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  provider = "alicloud.sh"
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "vpc3" {
  provider = "alicloud.hz"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_instance" "cen" {
    name = "${var.name}"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
    bandwidth = 5
    geographic_region_ids = [
		"Europe",
		"China"]
}

resource "alicloud_cen_bandwidth_package_attachment" "bwp_attach" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_3" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc3.id}"
    child_instance_region_id = "cn-hangzhou"
}

resource "alicloud_cen_bandwidth_limit" "bar1" {
	provider = "alicloud.fra"
    instance_id = "${alicloud_cen_instance.cen.id}"
    region_ids = [
        "eu-central-1",
        "cn-shanghai"]
     bandwidth_limit = 2
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}

resource "alicloud_cen_bandwidth_limit" "bar2" {
	provider = "alicloud.fra"
    instance_id = "${alicloud_cen_instance.cen.id}"
    region_ids = [
        "eu-central-1",
        "cn-hangzhou"]
     bandwidth_limit = 3
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_3"]
}
`
