package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamAccessKey_basic(t *testing.T) {
	var v ram.AccessKey
	var u ram.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_access_key.ak",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamAccessKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamAccessKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists(
						"alicloud_ram_user.user", &u),
					testAccCheckRamAccessKeyExists(
						"alicloud_ram_access_key.ak", &v),
					resource.TestCheckResourceAttr(
						"alicloud_ram_access_key.ak",
						"status",
						"Active"),
				),
			},
		},
	})

}

func testAccCheckRamAccessKeyExists(n string, ak *ram.AccessKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access key ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		response, err := conn.ListAccessKeys(request)

		if err == nil {
			if len(response.AccessKeys.AccessKey) > 0 {
				for _, v := range response.AccessKeys.AccessKey {
					if v.AccessKeyId == rs.Primary.ID {
						*ak = v
						return nil
					}
				}
			}
			return fmt.Errorf("Error finding access key %s", rs.Primary.ID)
		}
		return fmt.Errorf("Error finding access key %s: %#v", rs.Primary.ID, err)
	}
}

func testAccCheckRamAccessKeyDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_access_key" {
			continue
		}

		// Try to find the ak
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		response, err := conn.ListAccessKeys(request)

		if len(response.AccessKeys.AccessKey) > 0 {
			for _, v := range response.AccessKeys.AccessKey {
				if v.AccessKeyId == rs.Primary.ID {
					return fmt.Errorf("Error Access Key still exist")
				}
			}
		}
		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

const testAccRamAccessKeyConfig = `
resource "alicloud_ram_user" "user" {
  name = "username"
  display_name = "displayname"
  mobile = "86-18888888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
}

resource "alicloud_ram_access_key" "ak" {
  user_name = "${alicloud_ram_user.user.name}"
  status = "Active"
  secret_file = "/hello.txt"
}`
