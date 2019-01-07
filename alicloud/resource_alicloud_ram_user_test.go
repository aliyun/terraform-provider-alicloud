package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"regexp"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_user", &resource.Sweeper{
		Name: "alicloud_ram_user",
		F:    testSweepRamUsers,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_ram_role",
		},
	})
}

func testSweepRamUsers(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var users []ram.User
	args := ram.ListUserRequest{}
	for {
		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListUsers(args)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Ram users: %s", err)
		}
		resp, _ := raw.(ram.ListUserResponse)
		if len(resp.Users.User) < 1 {
			break
		}
		users = append(users, resp.Users.User...)

		if !resp.IsTruncated {
			break
		}
		args.Marker = resp.Marker
	}
	sweeped := false

	for _, v := range users {
		name := v.UserName
		id := v.UserId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ram User: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Ram User: %s (%s)", name, id)
		req := ram.UserQueryRequest{
			UserName: name,
		}
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DeleteUser(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ram User (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudRamUser_basic(t *testing.T) {
	var v ram.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_user.user",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamUserConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists(
						"alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr(
						"alicloud_ram_user.user",
						"name",
						regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ram_user.user",
						"display_name",
						"displayname"),
					resource.TestCheckResourceAttr(
						"alicloud_ram_user.user",
						"comments",
						"yoyoyo"),
				),
			},
		},
	})

}

func testAccCheckRamUserExists(n string, user *ram.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.GetUser(request)
		})
		log.Printf("[WARN] User id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.(ram.UserResponse)
			*user = response.User
			return nil
		}
		return fmt.Errorf("Error finding user %#v", rs.Primary.ID)
	}
}

func testAccCheckRamUserDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user" {
			continue
		}

		// Try to find the user
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.GetUser(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

func testAccRamUserConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}`, rand)
}
