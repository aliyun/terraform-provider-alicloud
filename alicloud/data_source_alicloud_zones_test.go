package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudZonesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudZonesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.foo"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_disk_categories.#"),
				),
			},
		},
	})
}

func TestAccAlicloudZonesDataSource_filter(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudZonesDataSourceFilter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.foo"),
					testCheckZoneLength("data.alicloud_zones.foo"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_disk_categories.#"),
				),
			},

			{
				Config: testAccCheckAlicloudZonesDataSourceFilterIoOptimized,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.foo"),
					testCheckZoneLength("data.alicloud_zones.foo"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_disk_categories.#"),
				),
			},
		},
	})
}

func TestAccAlicloudZonesDataSource_unitRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudZonesDataSource_unitRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.foo"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_disk_categories.#"),
				),
			},
		},
	})
}

func TestAccAlicloudZonesDataSource_multiZone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudZonesDataSource_multiZone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.#"),
					resource.TestMatchResourceAttr("data.alicloud_zones.default", "zones.0.id", regexp.MustCompile(fmt.Sprintf(".%s.", MULTI_IZ_SYMBOL))),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.0.local_name", ""),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.0.available_instance_types.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.0.available_resource_creation.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.0.available_disk_categories.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudZonesDataSource_chargeType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudZonesDataSource_chargeType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.0.local_name", ""),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.0.available_instance_types.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.0.available_resource_creation.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.0.available_disk_categories.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudZonesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudZonesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.default"),
					resource.TestCheckResourceAttr("data.alicloud_zones.default", "zones.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_zones.default", "zones.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_zones.default", "zones.local_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_zones.default", "zones.available_instance_types"),
					resource.TestCheckNoResourceAttr("data.alicloud_zones.default", "zones.available_resource_creation"),
					resource.TestCheckNoResourceAttr("data.alicloud_zones.default", "zones.available_disk_categories"),
				),
			},
		},
	})
}

// the zone length changed occasionally
// check by range to avoid test case failure
func testCheckZoneLength(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s", name)
		}

		i, err := strconv.Atoi(is.Attributes["zones.#"])

		if err != nil {
			return fmt.Errorf("convert zone length err: %#v", err)
		}

		if i <= 0 {
			return fmt.Errorf("zone length expected greater than 0 got err: %d", i)
		}

		return nil
	}
}

const testAccCheckAlicloudZonesDataSourceBasicConfig = `
data "alicloud_zones" "foo" {
}
`

const testAccCheckAlicloudZonesDataSourceFilter = `
data "alicloud_zones" "foo" {
	available_resource_creation= "VSwitch"
	available_disk_category= "cloud_efficiency"
}
`

const testAccCheckAlicloudZonesDataSourceFilterIoOptimized = `
provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_zones" "foo" {
	available_resource_creation= "IoOptimized"
	available_disk_category= "cloud_ssd"
}
`

const testAccCheckAlicloudZonesDataSource_unitRegion = `
provider "alicloud" {
	alias = "northeast"
	region = "ap-southeast-1"
}

data "alicloud_zones" "foo" {
	provider = "alicloud.northeast"
	available_resource_creation= "VSwitch"
}
`

const testAccCheckAlicloudZonesDataSource_multiZone = `
provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_zones" "default" {
  available_resource_creation= "Rds"
  multi = true
}`

const testAccCheckAlicloudZonesDataSource_chargeType = `
provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_zones" "default" {
  instance_charge_type = "PrePaid"
  available_resource_creation= "Rds"
  multi = true
}`

const testAccCheckAlicloudZonesDataSourceEmpty = `
data "alicloud_zones" "default" {
  available_instance_type = "ecs.n1.fake"
}
`
