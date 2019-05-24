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

func TestAccAlicloudRamAccessKey_basic(t *testing.T) {
	var v ram.AccessKey
	var u ram.User
	resourceAKId := "alicloud_ram_access_key.default"
	resourceUserId := "alicloud_ram_user.default"
	ra := resourceAttrInit("alicloud_ram_access_key.default", accessKeyBasicMap)
	rand := acctest.RandIntRange(1000000, 9999999)

	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceAKId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamAccessKeyCreate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists(resourceAKId, &v),
					testAccCheckRamUserExists(resourceUserId, &u),
					testAccCheck(map[string]string{"user_name": fmt.Sprintf("tf-testAcc%sRamAccessKeyConfig%d", defaultRegionToTest, rand)}),
				),
			},
			{
				Config: testAccRamAccessKeyStatus(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists(resourceAKId, &v),
					testAccCheckRamUserExists(resourceUserId, &u),
					testAccCheck(accessKeyBasicMap),
				),
			},
		},
	})
}

func TestAccAlicloudRamAccessKey_multi(t *testing.T) {
	var v ram.AccessKey
	var u ram.User
	resourceAKId := "alicloud_ram_access_key.default.1"
	resourceUserId := "alicloud_ram_user.default"
	ra := resourceAttrInit(resourceAKId, accessKeyMultiMap)
	rand := acctest.RandIntRange(1000000, 9999999)

	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceAKId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeyMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyExists(resourceAKId, &v),
					testAccCheckRamUserExists(resourceUserId, &u),
					testAccCheck(map[string]string{"user_name": fmt.Sprintf("tf-testAcc%sRamAccessKeyConfig%d", defaultRegionToTest, rand)}),
				),
			},
		},
	})
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

func testAccRamAccessKeyCreate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
		name = "tf-testAcc%sRamAccessKeyConfig%d"
		display_name = "displayname"
		mobile = "86-18888888888"
		email = "hello.uuu@aaa.com"
		comments = "yoyoyo"
	}

	resource "alicloud_ram_access_key" "default" {
		user_name = "${alicloud_ram_user.default.name}"
		status = "Active"
		secret_file = "/hello.txt"
	}`, defaultRegionToTest, rand)
}

func testAccRamAccessKeyStatus(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamAccessKeyConfig%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_access_key" "default" {
	  user_name = "${alicloud_ram_user.default.name}"
	  status = "Inactive"
	  secret_file = "/hello.txt"
	}`, defaultRegionToTest, rand)
}

func testAccAccessKeyMulti(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamAccessKeyConfig%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"

	}

	resource "alicloud_ram_access_key" "default" {
	  user_name = "${alicloud_ram_user.default.name}"
	  status = "Active"
	  secret_file = "/hello.txt"
	  count = 2
	}`, defaultRegionToTest, rand)
}

var accessKeyBasicMap = map[string]string{
	"user_name":   CHECKSET,
	"status":      CHECKSET,
	"secret_file": "/hello.txt",
}
var accessKeyMultiMap = map[string]string{
	"user_name":   CHECKSET,
	"status":      "Active",
	"secret_file": "/hello.txt",
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
