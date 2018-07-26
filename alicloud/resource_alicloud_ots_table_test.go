package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudOtsTable_Basic(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsTable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist(
						"alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist(
						"alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr(
						"alicloud_ots_table.basic",
						"table_name", "testAccOtsTable"),
				),
			},
		},
	})

}

func testAccCheckOtsTableExist(n string, table *tablestore.DescribeTableResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		response, err := client.DescribeOtsTable(split[0], split[1])

		if err != nil {
			return fmt.Errorf("Error finding OTS table %s: %#v", rs.Primary.ID, err)
		}

		table = response
		return nil
	}
}

func testAccCheckOtsTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ots_table" {
			continue
		}

		client := testAccProvider.Meta().(*AliyunClient)
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if _, err := client.DescribeOtsTable(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("error! Ots table still exists")
	}

	return nil
}

const testAccOtsTable = `
variable "name" {
  default = "testAccOtsTable"
}
resource "alicloud_ots_instance" "foo" {
  name = "${var.name}"
  description = "${var.name}"
  accessed_by = "Any"
  tags {
    Created = "TF"
    For = "acceptance test"
  }
}

resource "alicloud_ots_table" "basic" {
  instance_name = "${alicloud_ots_instance.foo.name}"
  table_name = "${var.name}"
  primary_key = {
    name = "pk1"
    type = "Integer"
  }
  time_to_live = -1
  max_version = 1
}
`
