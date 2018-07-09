package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamRole_basic(t *testing.T) {
	var v ram.Role

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_role.role",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamRoleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists(
						"alicloud_ram_role.role", &v),
					resource.TestCheckResourceAttr(
						"alicloud_ram_role.role",
						"name",
						"rolename"),
					resource.TestCheckResourceAttr(
						"alicloud_ram_role.role",
						"description",
						"this is a test"),
				),
			},
		},
	})

}

func testAccCheckRamRoleExists(n string, role *ram.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.RoleQueryRequest{
			RoleName: rs.Primary.Attributes["name"],
		}

		response, err := conn.GetRole(request)
		log.Printf("[WARN] Role id %#v", rs.Primary.ID)

		if err == nil {
			*role = response.Role
			return nil
		}
		return fmt.Errorf("Error finding role %#v", rs.Primary.ID)
	}
}

func testAccCheckRamRoleDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_role" {
			continue
		}

		// Try to find the role
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.RoleQueryRequest{
			RoleName: rs.Primary.Attributes["name"],
		}

		_, err := conn.GetRole(request)

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

const testAccRamRoleConfig = `
resource "alicloud_ram_role" "role" {
  name = "rolename"
  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
  description = "this is a test"
  force = true
}`
