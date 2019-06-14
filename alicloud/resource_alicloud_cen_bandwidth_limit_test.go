package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCenBandwidthLimit_basic(t *testing.T) {
	var v cbn.CenInterRegionBandwidthLimit
	resourceId := "alicloud_cen_bandwidth_limit.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, cenBandwidthLimitMap)
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenBandwidthLimitDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccCenBandwidthLimitCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(nil),
				),
			},
			{
				Config: testAccCenBandwidthLimitNumberConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"bandwidth_limit": "3",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthLimit_multi(t *testing.T) {
	var v cbn.CenInterRegionBandwidthLimit
	resourceId := "alicloud_cen_bandwidth_limit.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, cenBandwidthLimitMultiMap)
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenBandwidthLimitDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccCenBandwidthLimitMultiConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

var cenBandwidthLimitMap = map[string]string{
	"instance_id":     CHECKSET,
	"region_ids.#":    "2",
	"bandwidth_limit": "5",
}

var cenBandwidthLimitMultiMap = map[string]string{
	"instance_id":     CHECKSET,
	"region_ids.#":    "2",
	"bandwidth_limit": "2",
}

func testAccCenBandwidthLimitCreateConfig(rand int) string {
	return fmt.Sprintf(`
variable "name"{
    default = "tf-testAcc%sCenBandwidthLimitConfig-%d"
}

provider "alicloud" {
    alias = "fra"
    region = "eu-central-1"
}

provider "alicloud" {
    alias = "sh"
    region = "cn-shanghai"
}

resource "alicloud_vpc" "default" {
  provider = "alicloud.fra"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "default1" {
  provider = "alicloud.sh"
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "default" {
     name = "${var.name}"
     description = "tf-testAccCenBandwidthLimitConfigDescription"
}

resource "alicloud_cen_bandwidth_package" "default" {
    bandwidth = 5
    geographic_region_ids = [
		"Europe",
		"China"]
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.default.id}"
}

resource "alicloud_cen_instance_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default.id}"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "default1" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default1.id}"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidth_limit" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    region_ids = [
        "eu-central-1",
        "cn-shanghai"]
     bandwidth_limit = 5
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.default",
        "alicloud_cen_instance_attachment.default",
        "alicloud_cen_instance_attachment.default1"]
}
`, defaultRegionToTest, rand)
}

func testAccCenBandwidthLimitNumberConfig(rand int) string {
	return fmt.Sprintf(`
variable "name"{
    default = "tf-testAcc%sCenBandwidthLimitConfig-%d"
}

provider "alicloud" {
    alias = "fra"
    region = "eu-central-1"
}

provider "alicloud" {
    alias = "sh"
    region = "cn-shanghai"
}

resource "alicloud_vpc" "default" {
  provider = "alicloud.fra"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "default1" {
  provider = "alicloud.sh"
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "default" {
     name = "${var.name}"
     description = "tf-testAccCenBandwidthLimitConfigDescription"
}

resource "alicloud_cen_bandwidth_package" "default" {
    bandwidth = 5
    geographic_region_ids = [
		"Europe",
		"China"]
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.default.id}"
}

resource "alicloud_cen_instance_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default.id}"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "default1" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default1.id}"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidth_limit" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    region_ids = [
        "eu-central-1",
        "cn-shanghai"]
     bandwidth_limit = 3
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.default",
        "alicloud_cen_instance_attachment.default",
        "alicloud_cen_instance_attachment.default1"]
}
`, defaultRegionToTest, rand)
}

func testAccCenBandwidthLimitMultiConfig(rand int) string {
	return fmt.Sprintf(`
variable "name"{
    default = "tf-testAcc%sCenBandwidthLimitMulti-%d"
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

resource "alicloud_vpc" "default" {
  provider = "alicloud.fra"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "default1" {
  provider = "alicloud.sh"
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "default2" {
  provider = "alicloud.hz"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_instance" "default" {
    name = "${var.name}"
}

resource "alicloud_cen_bandwidth_package" "default" {
    bandwidth = 5
    geographic_region_ids = [
		"Europe",
		"China"]
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.default.id}"
}

resource "alicloud_cen_instance_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default.id}"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "default1" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default1.id}"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_instance_attachment" "default2" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default2.id}"
    child_instance_region_id = "cn-hangzhou"
}

resource "alicloud_cen_bandwidth_limit" "default" {
	provider = "alicloud.fra"
    instance_id = "${alicloud_cen_instance.default.id}"
    region_ids = [
        "eu-central-1",
        "cn-shanghai"]
     bandwidth_limit = 2
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.default",
        "alicloud_cen_instance_attachment.default",
        "alicloud_cen_instance_attachment.default1"]
}

resource "alicloud_cen_bandwidth_limit" "default1" {
	provider = "alicloud.fra"
    instance_id = "${alicloud_cen_instance.default.id}"
    region_ids = [
        "eu-central-1",
        "cn-hangzhou"]
     bandwidth_limit = 3
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.default",
        "alicloud_cen_instance_attachment.default",
        "alicloud_cen_instance_attachment.default2"]
}
`, defaultRegionToTest, rand)
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

			instance, err := cenService.DescribeCenBandwidthLimit(rs.Primary.ID)
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

		instance, err := cenService.DescribeCenBandwidthLimit(rs.Primary.ID)
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
