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
