package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
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
		return WrapError(err)
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
	request := ram.CreateListUsersRequest()
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsers(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram User: %s", err)
		}
		resp, _ := raw.(*ram.ListUsersResponse)
		if len(resp.Users.User) < 1 {
			break
		}
		users = append(users, resp.Users.User...)

		if !resp.IsTruncated {
			break
		}
		request.Marker = resp.Marker
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
		log.Printf("[INFO] Detaching Ram User policy: %s (%s)", name, id)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateListPoliciesForUserRequest()
			request.UserName = name
			return ramClient.ListPoliciesForUser(request)
		})
		if err != nil && !RamEntityNotExist(err) {
			log.Printf("[ERROR] ListPoliciesForUser: %s (%s)", name, id)
		}
		response, _ := raw.(*ram.ListPoliciesForUserResponse)
		if len(response.Policies.Policy) > 1 {
			request := ram.CreateDetachPolicyFromUserRequest()
			request.UserName = name

			for _, poloicy := range response.Policies.Policy {
				request.PolicyName = poloicy.PolicyName
				request.PolicyType = poloicy.PolicyType
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromUser(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					log.Printf("[ERROR] DetachPolicyFromUser: %s (%s)", name, id)
				}
			}
		}
		request1 := ram.CreateListGroupsForUserRequest()
		request1.UserName = name
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroupsForUser(request1)
		})
		if err != nil {
			log.Printf("[ERROR] ListGroupsForUser: %s (%s)", name, id)
		}
		groupResp, _ := raw.(*ram.ListGroupsForUserResponse)
		if len(groupResp.Groups.Group) > 0 {
			for _, v := range groupResp.Groups.Group {
				request := ram.CreateRemoveUserFromGroupRequest()
				request.UserName = name
				request.GroupName = v.GroupName
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					log.Printf("[ERROR] RemoveUserFromGroup: %s (%s)", name, id)
				}
			}
		}
		log.Printf("[INFO] Deleting Ram User: %s (%s)", name, id)
		request := ram.CreateDeleteUserRequest()
		request.UserName = name

		_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteUser(request)
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

func TestAccAlicloudRamUser_withrename(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig_default(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
			{
				Config: testAccRamUserConfig_withrename(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-new-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_withredispname(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig_withdispname(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
			{
				Config: testAccRamUserConfig_withdispnameupdate(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "new_displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_withmobile(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig_withmobile(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
			{
				Config: testAccRamUserConfig_withmobileupdate(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-16666666666"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_withemail(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig_withemail(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
			{
				Config: testAccRamUserConfig_withemailupdate(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.world@163.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_withcomments(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig_withcomments(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
			{
				Config: testAccRamUserConfig_withcommentsupdate(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", ""),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "RamUser_pls"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "force", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_multirename(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
			{
				Config: testAccRamUserConfig_multirename(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-new-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_multiredispname(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
			{
				Config: testAccRamUserConfig_multiredispname(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "new_displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_multiremobile(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
			{
				Config: testAccRamUserConfig_multiremobile(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-16666666666"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_multireEmail(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
			{
				Config: testAccRamUserConfig_multireEmail(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.world@163.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
		},
	})
}

func TestAccAlicloudRamUser_multirecomments(t *testing.T) {
	var v ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

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
				Config: testAccRamUserConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "yoyoyo"),
				),
			},
			{
				Config: testAccRamUserConfig_multirecomments(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists("alicloud_ram_user.user", &v),
					resource.TestMatchResourceAttr("alicloud_ram_user.user", "name", regexp.MustCompile("^tf-testAccRamUserConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "display_name", "displayname"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "mobile", "86-18888888888"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "email", "hello.uuu@aaa.com"),
					resource.TestCheckResourceAttr("alicloud_ram_user.user", "comments", "RamUser_pls"),
				),
			},
		},
	})
}

func testAccCheckRamUserExists(n string, user *ram.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No user ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetUserRequest()
		request.UserName = rs.Primary.Attributes["name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetUser(request)
		})
		log.Printf("[WARN] User id %#v", rs.Primary.ID)
		if err == nil {
			response, _ := raw.(*ram.GetUserResponse)
			*user = response.User
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckRamUserDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user" {
			continue
		}

		// Try to find the user
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetUserRequest()
		request.UserName = rs.Primary.Attributes["name"]

		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetUser(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
	}
	return nil
}

func testAccRamUserConfig_default(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	}`, rand)
}

func testAccRamUserConfig_withrename(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-new-%d"
	}`, rand)
}

func testAccRamUserConfig_withdispname(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  display_name = "displayname"
	}`, rand)
}

func testAccRamUserConfig_withdispnameupdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  display_name = "new_displayname"
	}`, rand)
}

func testAccRamUserConfig_withmobile(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  mobile = "86-18888888888"
	}`, rand)
}

func testAccRamUserConfig_withmobileupdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  mobile = "86-16666666666"
	}`, rand)
}

func testAccRamUserConfig_withemail(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  email = "hello.uuu@aaa.com"
	}`, rand)
}

func testAccRamUserConfig_withemailupdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  email = "hello.world@163.com"
	}`, rand)
}

func testAccRamUserConfig_withcomments(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  comments = "yoyoyo"
	}`, rand)
}

func testAccRamUserConfig_withcommentsupdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  comments = "RamUser_pls"
	}`, rand)
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

func testAccRamUserConfig_multirename(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-new-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}`, rand)
}

func testAccRamUserConfig_multiredispname(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  display_name = "new_displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}`, rand)
}

func testAccRamUserConfig_multiremobile(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  display_name = "displayname"
	  mobile = "86-16666666666"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}`, rand)
}

func testAccRamUserConfig_multireEmail(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.world@163.com"
	  comments = "yoyoyo"
	}`, rand)
}

func testAccRamUserConfig_multirecomments(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamUserConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "RamUser_pls"
	}`, rand)
}
