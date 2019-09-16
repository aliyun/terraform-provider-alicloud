package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDdosbgpInstanceDataSource_basic(t *testing.T) {
	resource.Test(t, getTestCase(t))
}

func testAccCheckAlicloudDdosbgpInstanceDataSource_from_name(randInt int, region string, base_bandwidth int, bandwidth int, ip_count int) string {
	return fmt.Sprintf(`
    data "alicloud_ddosbgp_instances" "instance" {
        name_regex = "${alicloud_ddosbgp_instance.foo.name}"
        region = "%s"
    }

    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddosbgp_instance" "foo" {
      name                    = "tf_testAcc_%v"
      region                  = "%s"
      base_bandwidth          = "%d"	
      bandwidth               = "%d"
      ip_count                = "%d"
      ip_type                 = "v4"
      type                    = "1"
	}

	`, region, randInt, region, base_bandwidth, bandwidth, ip_count)
}

func testAccCheckAlicloudDdosbgpInstanceDataSource_from_ids(randInt int, region string, base_bandwidth int, bandwidth int, ip_count int) string {
	return fmt.Sprintf(`
    data "alicloud_ddosbgp_instances" "instance" {
        ids = [ "${alicloud_ddosbgp_instance.foo.id}" ]
        region = "%s"
    }

    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddosbgp_instance" "foo" {
      name                    = "tf_testAcc_%v"
      region                  = "%s"
      base_bandwidth          = "%d"	
      bandwidth               = "%d"
      ip_count                = "%d"
      ip_type                 = "v4"
      type                    = "1"
	}

	`, region, randInt, region, base_bandwidth, bandwidth, ip_count)
}

func testAccCheckAlicloudDdosbgpInstanceDataSource_from_both(randInt int, region string, base_bandwidth int, bandwidth int, ip_count int) string {
	return fmt.Sprintf(`
    data "alicloud_ddosbgp_instances" "instance" {
        name_regex = "${alicloud_ddosbgp_instance.foo.name}"
        ids = [ "${alicloud_ddosbgp_instance.foo.id}" ]
        region = "%s"
    }

    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddosbgp_instance" "foo" {
      name                    = "tf_testAcc_%v"
      region                  = "%s"
      base_bandwidth          = "%d"	
      bandwidth               = "%d"
      ip_count                = "%d"
      ip_type                 = "v4"
      type                    = "1"
	}`, region, randInt, region, base_bandwidth, bandwidth, ip_count)
}

func testAccCheckAlicloudDdosbgpInstanceDataSourceNameRegexConfig_mismatch_name(randInt int, region string, base_bandwidth int, bandwidth int, ip_count int) string {
	return fmt.Sprintf(`
    data "alicloud_ddosbgp_instances" "instance" {
        name_regex = "${alicloud_ddosbgp_instance.foo.name}-fake"
        region = "%s"
    }

    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddosbgp_instance" "foo" {
      name                    = "tf_testAcc_%v"
      region                  = "%s"
      base_bandwidth          = "%d"	
      bandwidth               = "%d"
      ip_count                = "%d"
      ip_type                 = "v4"
      type                    = "1"
	}`, region, randInt, region, base_bandwidth, bandwidth, ip_count)
}

func testAccCheckAlicloudDdosbgpInstanceDataSourceNameRegexConfig_mismatch_ids(randInt int, region string, base_bandwidth int, bandwidth int, ip_count int) string {
	return fmt.Sprintf(`
    data "alicloud_ddosbgp_instances" "instance" {
        ids = [ "${alicloud_ddosbgp_instance.foo.id}-fake" ]
        region = "%s"
    }

    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddosbgp_instance" "foo" {
      name                    = "tf_testAcc_%v"
      region                  = "%s"
      base_bandwidth          = "%d"	
      bandwidth               = "%d"
      ip_count                = "%d"
      ip_type                 = "v4"
      type                    = "1"
	}`, region, randInt, region, base_bandwidth, bandwidth, ip_count)
}

func testAccCheckAlicloudDdosbgpInstanceDataSourceNameRegexConfig_mismatch_all(randInt int, region string, base_bandwidth int, bandwidth int, ip_count int) string {
	return fmt.Sprintf(`
    data "alicloud_ddosbgp_instances" "instance" {
        name_regex = "${alicloud_ddosbgp_instance.foo.name}-fake"
        ids = [ "${alicloud_ddosbgp_instance.foo.id}-fake" ]
        region = "%s"
    }

    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddosbgp_instance" "foo" {
      name                    = "tf_testAcc_%v"
      region                  = "%s"
      base_bandwidth          = "%d"	
      bandwidth               = "%d"
      ip_count                = "%d"
      ip_type                 = "v4"
      type                    = "1"
	}`, region, randInt, region, base_bandwidth, bandwidth, ip_count)
}

func generateRegionTestCase(rand int, region string, bandwidth int) []resource.TestStep {
	return []resource.TestStep{
		{
			Config: testAccCheckAlicloudDdosbgpInstanceDataSource_from_name(rand, region, 20, bandwidth, 100),
			Check: resource.ComposeTestCheckFunc(
				testAccCheckAlicloudDataSourceID("data.alicloud_ddosbgp_instances.instance"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.#", "1"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "names.#", "1"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "ids.#", "1"),
				resource.TestCheckResourceAttrSet("data.alicloud_ddosbgp_instances.instance", "instances.0.id"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.bandwidth", strconv.Itoa(bandwidth)),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.base_bandwidth", "20"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.ip_count", "100"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.ip_type", "IPv4"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.region", region),
			),
		},
		{
			Config: testAccCheckAlicloudDdosbgpInstanceDataSource_from_ids(rand, region, 20, bandwidth, 100),
			Check: resource.ComposeTestCheckFunc(
				testAccCheckAlicloudDataSourceID("data.alicloud_ddosbgp_instances.instance"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.#", "1"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "names.#", "1"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "ids.#", "1"),
				resource.TestCheckResourceAttrSet("data.alicloud_ddosbgp_instances.instance", "instances.0.id"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.bandwidth", strconv.Itoa(bandwidth)),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.base_bandwidth", "20"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.ip_count", "100"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.ip_type", "IPv4"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.region", region),
			),
		},
		{
			Config: testAccCheckAlicloudDdosbgpInstanceDataSource_from_both(rand, region, 20, bandwidth, 100),
			Check: resource.ComposeTestCheckFunc(
				testAccCheckAlicloudDataSourceID("data.alicloud_ddosbgp_instances.instance"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.#", "1"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "names.#", "1"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "ids.#", "1"),
				resource.TestCheckResourceAttrSet("data.alicloud_ddosbgp_instances.instance", "instances.0.id"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.bandwidth", strconv.Itoa(bandwidth)),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.base_bandwidth", "20"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.ip_count", "100"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.ip_type", "IPv4"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.0.region", region),
			),
		},
		{
			Config: testAccCheckAlicloudDdosbgpInstanceDataSourceNameRegexConfig_mismatch_name(rand, region, 20, bandwidth, 100),
			Check: resource.ComposeTestCheckFunc(
				testAccCheckAlicloudDataSourceID("data.alicloud_ddosbgp_instances.instance"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.#", "0"),
			),
		},
		{
			Config: testAccCheckAlicloudDdosbgpInstanceDataSourceNameRegexConfig_mismatch_ids(rand, region, 20, bandwidth, 100),
			Check: resource.ComposeTestCheckFunc(
				testAccCheckAlicloudDataSourceID("data.alicloud_ddosbgp_instances.instance"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.#", "0"),
			),
		},
		{
			Config: testAccCheckAlicloudDdosbgpInstanceDataSourceNameRegexConfig_mismatch_all(rand, region, 20, bandwidth, 100),
			Check: resource.ComposeTestCheckFunc(
				testAccCheckAlicloudDataSourceID("data.alicloud_ddosbgp_instances.instance"),
				resource.TestCheckResourceAttr("data.alicloud_ddosbgp_instances.instance", "instances.#", "0"),
			),
		},
	}
}

func getTestCase(t *testing.T) resource.TestCase {
	rand := acctest.RandInt()
	var testCase = resource.TestCase{}
	testCase.PreCheck = func() {
		testAccPreCheckWithRegions(t, true, connectivity.DdosbgpSupportedRegions)
		testAccPreCheck(t)
	}
	testCase.Providers = testAccProviders

	var steps []resource.TestStep
	steps = append(steps, generateRegionTestCase(rand, "cn-hangzhou", 101)...)
	steps = append(steps, generateRegionTestCase(rand, "cn-shanghai", 101)...)
	steps = append(steps, generateRegionTestCase(rand, "cn-qingdao", 51)...)
	steps = append(steps, generateRegionTestCase(rand, "cn-beijing", 101)...)
	steps = append(steps, generateRegionTestCase(rand, "cn-zhangjiakou", 101)...)
	steps = append(steps, generateRegionTestCase(rand, "cn-huhehaote", 101)...)
	steps = append(steps, generateRegionTestCase(rand, "cn-shenzhen", 51)...)

	testCase.Steps = steps
	return testCase
}
