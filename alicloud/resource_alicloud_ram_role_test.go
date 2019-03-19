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
	resource.AddTestSweepers("alicloud_ram_role", &resource.Sweeper{
		Name: "alicloud_ram_role",
		F:    testSweepRamRoles,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_fc_service",
		},
	})
}

func testSweepRamRoles(region string) error {
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

	request := ram.CreateListRolesRequest()
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListRoles(request)
	})
	if err != nil {
		return WrapError(err)
	}
	resp, _ := raw.(*ram.ListRolesResponse)
	if len(resp.Roles.Role) < 1 {
		return nil
	}

	sweeped := false

	for _, v := range resp.Roles.Role {
		name := v.RoleName
		id := v.RoleId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ram Role: %s (%s)", name, id)
			continue
		}
		sweeped = true

		log.Printf("[INFO] Detaching Ram Role: %s (%s) policies.", name, id)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateListPoliciesForRoleRequest()
			request.RoleName = name
			return ramClient.ListPoliciesForRole(request)
		})
		resp, _ := raw.(*ram.ListPoliciesForRoleResponse)
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram Role (%s (%s)) policies: %s", name, id, err)
		} else if len(resp.Policies.Policy) > 0 {
			for _, v := range resp.Policies.Policy {
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					request := ram.CreateDetachPolicyFromRoleRequest()
					request.PolicyName = v.PolicyName
					request.RoleName = name
					request.PolicyType = v.PolicyType
					return ramClient.DetachPolicyFromRole(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					log.Printf("[ERROR] Failed detach Policy %s: %#v", v.PolicyName, err)
				}
			}
		}

		log.Printf("[INFO] Deleting Ram Role: %s (%s)", name, id)
		request := ram.CreateDeleteRoleRequest()
		request.RoleName = name

		_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteRole(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ram Role (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

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
			{
				Config: testAccRamRoleConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists(
						"alicloud_ram_role.role", &v),
					resource.TestMatchResourceAttr(
						"alicloud_ram_role.role",
						"name",
						regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ram_role.role",
						"description",
						"this is a test"),
				),
			},
		},
	})

}

func TestAccAlicloudRamRole_RamUsers(t *testing.T) {
	var v ram.Role
	var u ram.User
	randInt := acctest.RandIntRange(1000000, 99999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_role.role",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamRoleConfig_RamUsers(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists("alicloud_ram_role.role", &v),
					testAccCheckRamUserExists("alicloud_ram_user.user1", &u),
					resource.TestMatchResourceAttr("alicloud_ram_role.role","name",regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","ram_users.0.user_name",fmt.Sprintf("tf-testAccRamGroupConfig-%v.a", randInt)),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","description","this is a test"),
				),
			},
			{
				Config: testAccRamRoleConfig_newRamUsers(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists("alicloud_ram_role.role", &v),
					testAccCheckRamUserExists("alicloud_ram_user.user2", &u),
					resource.TestMatchResourceAttr("alicloud_ram_role.role","name",regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","ram_users.0.user_name",fmt.Sprintf("tf-testAccRamGroupConfig-%v.b", randInt)),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","description","this is a test"),
				),
			},
		},
	})

}

func TestAccAlicloudRamRole_reDocument(t *testing.T) {
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
			{
				Config: testAccRamRoleConfig_document1(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists("alicloud_ram_role.role", &v),
					resource.TestMatchResourceAttr("alicloud_ram_role.role","name",regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","document","false"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","description","this is a test"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","force", "false"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role", "arn", string(Available)),
				),
			},
			{
				Config: testAccRamRoleConfig_document2(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists("alicloud_ram_role.role", &v),
					resource.TestMatchResourceAttr("alicloud_ram_role.role","name",regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","document","true"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","description","this is a test"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","force", "false"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role", "arn", string(Available)),
				),
			},
		},
	})

}

func TestAccAlicloudRamRole_reServices(t *testing.T) {
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
			{
				Config: testAccRamRoleConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists("alicloud_ram_role.role", &v),
					resource.TestMatchResourceAttr("alicloud_ram_role.role","name",regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","services.0","apigateway.aliyuncs.com"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","services.1", "ecs.aliyuncs.com"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","description","this is a test"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","force", "true"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role", "arn", string(Available)),
				),
			},
			{
				Config: testAccRamRoleConfig_services(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists("alicloud_ram_role.role", &v),
					resource.TestMatchResourceAttr("alicloud_ram_role.role","name",regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","services.0","rds.aliyuncs.com"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","services.1", "oss.aliyuncs.com"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","description","this is a test"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","force", "true"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role", "arn", string(Available)),
				),
			},
		},
	})

}

func TestAccAlicloudRamRole_version(t *testing.T) {
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
			{
				Config: testAccRamRoleConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists("alicloud_ram_role.role", &v),
					resource.TestMatchResourceAttr("alicloud_ram_role.role","name",regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","description","this is a test"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","version","1"),
				),
			},
			{
				Config: testAccRamRoleConfig_version(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists("alicloud_ram_role.role", &v),
					resource.TestMatchResourceAttr("alicloud_ram_role.role","name",regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","description","this is a test"),
					resource.TestCheckResourceAttr("alicloud_ram_role.role","version","2"),
				),
			},
		},
	})

}

func testAccCheckRamRoleExists(n string, role *ram.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No role ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetRoleRequest()
		request.RoleName = rs.Primary.Attributes["name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		log.Printf("[WARN] Role id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.(*ram.GetRoleResponse)
			*role = response.Role
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckRamRoleDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_role" {
			continue
		}

		// Try to find the role
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetRoleRequest()
		request.RoleName = rs.Primary.Attributes["name"]

		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
	}
	return nil
}

func testAccRamRoleConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "role" {
	  name = "tf-testAccRamRoleConfig-%d"
	  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
	  description = "this is a test"
	  force = true
	}`, rand)
}

func testAccRamRoleConfig_services(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "role" {
	  name = "tf-testAccRamRoleConfig-%d"
	  services = ["rds.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = true
	}`, rand)
}

func testAccRamRoleConfig_document1(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "role" {
	  name = "tf-testAccRamRoleConfig-%d"
	  document = false
	  description = "this is a test"
	}`, rand)
}

func testAccRamRoleConfig_document2(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "role" {
	  name = "tf-testAccRamRoleConfig-%d"
	  document = true
	  description = "this is a test"
	}`, rand)
}

func testAccRamRoleConfig_version(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "role" {
	  name = "tf-testAccRamRoleConfig-%d"
	  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
	  version = "2"
	  description = "this is a test"
	  force = true
	}`, rand)
}

func testAccRamRoleConfig_RamUsers(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamRoleConfig-%v"
	}
	resource "alicloud_ram_user" "user1" {
	  name = "${var.name}.a"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}
	resource "alicloud_ram_role" "role" {
	  name = "${var.name}"
	  ram_users = ["${alicloud_ram_user.user1}"]
	  description = "this is a test"
	}`, rand)
}

func testAccRamRoleConfig_newRamUsers(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamRoleConfig-%v"
	}
	resource "alicloud_ram_user" "user2" {
	  name = "${var.name}.b"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}
	resource "alicloud_ram_role" "role" {
	  name = "${var.name}"
	  ram_users = ["${alicloud_ram_user.user2}"]
	  description = "this is a test"
	}`, rand)
}



