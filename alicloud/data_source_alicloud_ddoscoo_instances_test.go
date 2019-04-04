package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDdoscooInstanceDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDdoscooInstanceDataSource_from_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ddoscoo_instances.instance"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ddoscoo_instances.instance", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.bandwidth", "30"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.base_bandwidth", "30"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.service_bandwidth", "100"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.port_count", "50"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.domain_count", "50"),
				),
			},
			{
				Config: testAccCheckAlicloudDdoscooInstanceDataSource_from_ids(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ddoscoo_instances.instance"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ddoscoo_instances.instance", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.bandwidth", "30"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.base_bandwidth", "30"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.service_bandwidth", "100"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.port_count", "50"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.domain_count", "50"),
				),
			},
			{
				Config: testAccCheckAlicloudDdoscooInstanceDataSource_from_both(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ddoscoo_instances.instance"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ddoscoo_instances.instance", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.bandwidth", "30"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.base_bandwidth", "30"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.service_bandwidth", "100"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.port_count", "50"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.0.domain_count", "50"),
				),
			},
			{
				Config: testAccCheckAlicloudDdoscooInstanceDataSourceNameRegexConfig_mismatch_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ddoscoo_instances.instance"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.#", "0"),
				),
			},
			{
				Config: testAccCheckAlicloudDdoscooInstanceDataSourceNameRegexConfig_mismatch_ids(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ddoscoo_instances.instance"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.#", "0"),
				),
			},
			{
				Config: testAccCheckAlicloudDdoscooInstanceDataSourceNameRegexConfig_mismatch_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ddoscoo_instances.instance"),
					resource.TestCheckResourceAttr("data.alicloud_ddoscoo_instances.instance", "instances.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDdoscooInstanceDataSource_from_name(randInt int) string {
	return fmt.Sprintf(`
    data "alicloud_ddoscoo_instances" "instance" {
        name_regex = "${alicloud_ddoscoo_instance.foo.name}"
    }

    provider "alicloud" {
        endpoints = {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddoscoo_instance" "foo" {
      name                    = "tf_testAcc%v"
      bandwidth               = "30"
      base_bandwidth          = "30"
      service_bandwidth       = "100"
      port_count              = "50"
      domain_count            = "50"
	}`, randInt)
}

func testAccCheckAlicloudDdoscooInstanceDataSource_from_ids(randInt int) string {
	return fmt.Sprintf(`
    data "alicloud_ddoscoo_instances" "instance" {
        ids = [ "${alicloud_ddoscoo_instance.foo.id}" ]
    }

    provider "alicloud" {
        endpoints = {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddoscoo_instance" "foo" {
      name                    = "tf_testAcc%v"
      bandwidth               = "30"
      base_bandwidth          = "30"
      service_bandwidth       = "100"
      port_count              = "50"
      domain_count            = "50"
	}`, randInt)
}

func testAccCheckAlicloudDdoscooInstanceDataSource_from_both(randInt int) string {
	return fmt.Sprintf(`
    data "alicloud_ddoscoo_instances" "instance" {
        name_regex = "${alicloud_ddoscoo_instance.foo.name}"
        ids = [ "${alicloud_ddoscoo_instance.foo.id}" ]
    }

    provider "alicloud" {
        endpoints = {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddoscoo_instance" "foo" {
      name                    = "tf_testAcc%v"
      bandwidth               = "30"
      base_bandwidth          = "30"
      service_bandwidth       = "100"
      port_count              = "50"
      domain_count            = "50"
	}`, randInt)
}

func testAccCheckAlicloudDdoscooInstanceDataSourceNameRegexConfig_mismatch_name(randInt int) string {
	return fmt.Sprintf(`
    data "alicloud_ddoscoo_instances" "instance" {
        name_regex = "${alicloud_ddoscoo_instance.foo.name}-fake"
    }

    provider "alicloud" {
        endpoints = {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddoscoo_instance" "foo" {
      name                    = "tf_testAcc%v"
      bandwidth               = "30"
      base_bandwidth          = "30"
      service_bandwidth       = "100"
      port_count              = "50"
      domain_count            = "50"
	}`, randInt)
}

func testAccCheckAlicloudDdoscooInstanceDataSourceNameRegexConfig_mismatch_ids(randInt int) string {
	return fmt.Sprintf(`
    data "alicloud_ddoscoo_instances" "instance" {
        ids = [ "${alicloud_ddoscoo_instance.foo.id}-fake" ]
    }

    provider "alicloud" {
        endpoints = {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddoscoo_instance" "foo" {
      name                    = "tf_testAcc%v"
      bandwidth               = "30"
      base_bandwidth          = "30"
      service_bandwidth       = "100"
      port_count              = "50"
      domain_count            = "50"
	}`, randInt)
}

func testAccCheckAlicloudDdoscooInstanceDataSourceNameRegexConfig_mismatch_all(randInt int) string {
	return fmt.Sprintf(`
    data "alicloud_ddoscoo_instances" "instance" {
        name_regex = "${alicloud_ddoscoo_instance.foo.name}-fake"
        ids = [ "${alicloud_ddoscoo_instance.foo.id}-fake" ]
    }

    provider "alicloud" {
        endpoints = {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddoscoo_instance" "foo" {
      name                    = "tf_testAcc%v"
      bandwidth               = "30"
      base_bandwidth          = "30"
      service_bandwidth       = "100"
      port_count              = "50"
      domain_count            = "50"
	}`, randInt)
}
