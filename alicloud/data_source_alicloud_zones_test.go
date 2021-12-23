package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
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
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
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
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
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
				Config: testAccCheckAlicloudZonesDataSourceUnitRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.foo"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.available_disk_categories.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudZonesDataSource_slb(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudZonesDataSource_slb,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.0.local_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.0.available_instance_types.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.0.available_resource_creation.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.0.available_disk_categories.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "zones.0.slb_slave_zone_ids.#"),
				),
			},
		},
	})
}

func TestAccAlicloudZonesDataSource_enable_details(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudZonesDataSourceEnableDetails,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_zones.foo"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "zones.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.local_name", ""),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.available_instance_types.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.available_resource_creation.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.available_disk_categories.#", "0"),
					resource.TestCheckResourceAttrSet("data.alicloud_zones.foo", "ids.#"),
					resource.TestCheckResourceAttr("data.alicloud_zones.foo", "zones.0.slb_slave_zone_ids.#", "0"),
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
					resource.TestCheckResourceAttrSet("data.alicloud_zones.default", "ids.#"),
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
	enable_details = true
}
`

const testAccCheckAlicloudZonesDataSourceFilter = `
data "alicloud_zones" "foo" {
	available_resource_creation= "VSwitch"
	available_disk_category= "cloud_efficiency"
	enable_details = true
}
`

const testAccCheckAlicloudZonesDataSourceFilterIoOptimized = `
data "alicloud_zones" "foo" {
	available_resource_creation= "IoOptimized"
	available_disk_category= "cloud_ssd"
	enable_details = true
}
`

const testAccCheckAlicloudZonesDataSourceUnitRegion = `
data "alicloud_zones" "foo" {
	available_resource_creation= "VSwitch"
	enable_details = true
}
`
const testAccCheckAlicloudZonesDataSource_slb = `
data "alicloud_zones" "default" {
  available_resource_creation= "Slb"
  enable_details = true
  available_slb_address_ip_version= "ipv4"
  available_slb_address_type="Vpc"
}`

const testAccCheckAlicloudZonesDataSourceEnableDetails = `
data "alicloud_zones" "foo" {}
`
const testAccCheckAlicloudZonesDataSourceEmpty = `
data "alicloud_zones" "default" {
  available_instance_type = "ecs.n1.fake"
}
`
