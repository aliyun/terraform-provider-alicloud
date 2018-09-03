package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCen_basic(t *testing.T) {
	var cen cbn.Cen

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenExists("alicloud_cen.foo", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen.foo", "name", "testAccCenConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_cen.foo", "description", "testAccCenConfigDescription"),
				),
			},
		},
	})

}

func TestAccAlicloudCen_update(t *testing.T) {
	var cen cbn.Cen

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenExists("alicloud_cen.foo", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen.foo", "name", "testAccCenConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_cen.foo", "description", "testAccCenConfigDescription"),
				),
			},
			resource.TestStep{
				Config: testAccCenConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenExists("alicloud_cen.foo", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen.foo", "name", "testAccCenConfigUpdate"),
					resource.TestCheckResourceAttr(
						"alicloud_cen.foo", "description", "testAccCenConfigDescriptionUpdate"),
				),
			},
		},
	})
}

func TestAccAlicloudCen_multi(t *testing.T) {
	var cen cbn.Cen

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenExists("alicloud_cen.bar_1", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen.bar_1", "name", "testAccCenConfig-1"),
					resource.TestCheckResourceAttr(
						"alicloud_cen.bar_1", "description", "testAccCenConfigDescription-1"),
					testAccCheckCenExists("alicloud_cen.bar_2", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen.bar_2", "name", "testAccCenConfig-2"),
					resource.TestCheckResourceAttr(
						"alicloud_cen.bar_2", "description", "testAccCenConfigDescription-2"),
					testAccCheckCenExists("alicloud_cen.bar_3", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen.bar_3", "name", "testAccCenConfig-3"),
					resource.TestCheckResourceAttr(
						"alicloud_cen.bar_3", "description", "testAccCenConfigDescription-3"),
				),
			},
		},
	})
}

func testAccCheckCenExists(n string, cen *cbn.Cen) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CEN ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeCen(rs.Primary.ID)

		if err != nil {
			return err
		}

		*cen = instance
		return nil
	}
}

func testAccCheckCenDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen" {
			continue
		}

		// Try to find the CEN
		instance, err := client.DescribeCen(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.CenId != "" {
			return fmt.Errorf("CEN %s still exist", instance.CenId)
		}
	}

	return nil
}

const testAccCenConfig = `
resource "alicloud_cen" "foo" {
	name = "testAccCenConfig"
	description = "testAccCenConfigDescription"
}
`

const testAccCenConfigUpdate = `
resource "alicloud_cen" "foo" {
	name = "testAccCenConfigUpdate"
	description = "testAccCenConfigDescriptionUpdate"
}
`

const testAccCenConfigMulti = `
resource "alicloud_cen" "bar_1" {
	name = "testAccCenConfig-1"
	description = "testAccCenConfigDescription-1"
}
resource "alicloud_cen" "bar_2" {
	name = "testAccCenConfig-2"
	description = "testAccCenConfigDescription-2"
}
resource "alicloud_cen" "bar_3" {
	name = "testAccCenConfig-3"
	description = "testAccCenConfigDescription-3"
}
`
