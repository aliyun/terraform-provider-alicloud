package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRamAccessKey_status(t *testing.T) {
	var v ram.AccessKey

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_access_key.ak",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamAccessKeyConfig_status,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists("alicloud_ram_access_key.ak", &v),
					resource.TestCheckNoResourceAttr("alicloud_ram_access_key.ak", "user_name"),
					resource.TestCheckNoResourceAttr("alicloud_ram_access_key.ak", "secret_file"),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "status", "Active"),
				),
			},
			{
				Config: testAccRamAccessKeyConfig_statuschange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists("alicloud_ram_access_key.ak", &v),
					resource.TestCheckNoResourceAttr("alicloud_ram_access_key.ak", "user_name"),
					resource.TestCheckNoResourceAttr("alicloud_ram_access_key.ak", "secret_file"),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "status", "Inactive"),
				),
			},
		},
	})

}

func TestAccAlicloudRamAccessKey_username(t *testing.T) {
	var v ram.AccessKey
	randInt := acctest.RandIntRange(1000000, 99999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_access_key.ak",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamAccessKeyConfig_username(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists("alicloud_ram_access_key.ak", &v),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "user_name", fmt.Sprintf("tf-testAccRamAccessKeyConfig%d", randInt)),
					resource.TestCheckNoResourceAttr("alicloud_ram_access_key.ak", "secret_file"),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "status", "Active"),
				),
			},
		},
	})

}

func TestAccAlicloudRamAccessKey_scretfile(t *testing.T) {
	var v ram.AccessKey

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_access_key.ak",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamAccessKeyConfig_secretfile,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists("alicloud_ram_access_key.ak", &v),
					resource.TestCheckNoResourceAttr("alicloud_ram_access_key.ak", "user_name"),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "secret_file", "/world.txt"),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "status", "Active"),
				),
			},
		},
	})

}

func TestAccAlicloudRamAccessKey_multi(t *testing.T) {
	var v ram.AccessKey
	randInt := acctest.RandIntRange(1000000, 99999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_access_key.ak",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamAccessKeyConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists("alicloud_ram_access_key.ak", &v),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "user_name", fmt.Sprintf("tf-testAccRamAccessKeyConfig%d", randInt)),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "status", "Active"),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "secret_file", "/hello.txt"),
				),
			},
			{
				Config: testAccRamAccessKeyConfig_multi(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists("alicloud_ram_access_key.ak", &v),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "user_name", fmt.Sprintf("tf-testAccRamAccessKeyConfig%d", randInt)),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "status", "Inactive"),
					resource.TestCheckResourceAttr("alicloud_ram_access_key.ak", "secret_file", "/hello.txt"),
				),
			},
		},
	})

}

func testAccCheckRamAccessKeyExists(n string, ak *ram.AccessKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Access key ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListAccessKeysRequest()
		request.UserName = rs.Primary.Attributes["user_name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListAccessKeys(request)
		})

		if err == nil {
			response, _ := raw.(*ram.ListAccessKeysResponse)
			if len(response.AccessKeys.AccessKey) > 0 {
				for _, v := range response.AccessKeys.AccessKey {
					if v.AccessKeyId == rs.Primary.ID {
						*ak = v
						return nil
					}
				}
			}
			return WrapError(fmt.Errorf("Error finding access key %s", rs.Primary.ID))
		}
		return WrapError(err)
	}
}

func testAccCheckRamAccessKeyDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_access_key" {
			continue
		}

		// Try to find the ak
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListAccessKeysRequest()
		request.UserName = rs.Primary.Attributes["user_name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListAccessKeys(request)
		})

		response, _ := raw.(*ram.ListAccessKeysResponse)
		if len(response.AccessKeys.AccessKey) > 0 {
			for _, v := range response.AccessKeys.AccessKey {
				if v.AccessKeyId == rs.Primary.ID {
					return WrapError(Error("Error Access Key still exist"))
				}
			}
		}
		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
	}
	return nil
}

func testAccRamAccessKeyConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamAccessKeyConfig%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_access_key" "ak" {
	  user_name = "${alicloud_ram_user.user.name}"
	  status = "Active"
	  secret_file = "/hello.txt"
	}`, rand)
}

const testAccRamAccessKeyConfig_status = `
resource "alicloud_ram_access_key" "ak" {
	  status = "Active"
}
`

const testAccRamAccessKeyConfig_statuschange = `
resource "alicloud_ram_access_key" "ak" {
	  status = "Inactive"
}
`

func testAccRamAccessKeyConfig_username(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamAccessKeyConfig%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_access_key" "ak" {
	  user_name = "${alicloud_ram_user.user.name}"
	}`, rand)
}

const testAccRamAccessKeyConfig_secretfile = `
resource "alicloud_ram_access_key" "ak" {
	  status = "Active"
	  secret_file = "/world.txt"
}
`

func testAccRamAccessKeyConfig_multi(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamAccessKeyConfig%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_access_key" "ak" {
	  user_name = "${alicloud_ram_user.user.name}"
	  status = "Inactive"
	  secret_file = "/hello.txt"
	}`, rand)
}
