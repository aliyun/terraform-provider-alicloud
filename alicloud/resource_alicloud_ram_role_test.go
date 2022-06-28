package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
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

func TestAccAlicloudRAMRole_basic(t *testing.T) {
	var v *ram.GetRoleResponse
	resourceId := "alicloud_ram_role.default"
	ra := resourceAttrInit(resourceId, ramRoleMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamRoleCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": fmt.Sprintf("tf-testAcc%sRamRoleConfig-%d", defaultRegionToTest, rand)}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
			{
				Config: testAccRamRoleNameConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": fmt.Sprintf("tf-testAcc%sRamRoleConfig-%d-N", defaultRegionToTest, rand)}),
				),
			},
			{
				Config: testAccRamRoleForceConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"force": "false"}),
				),
			},
			{
				Config: testAccRamRoleDocumentConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					// There is a bug when d.Set a set parameter. The new values can not overwrite the state
					// when a parameter is a TypeSet and Computed. https://github.com/hashicorp/terraform-plugin-sdk/issues/20504
					//testAccCheck(map[string]string{"services.#": "1"}),
					testAccCheck(nil),
				),
			},
			{
				Config: testAccRamRoleCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":       fmt.Sprintf("tf-testAcc%sRamRoleConfig-%d", defaultRegionToTest, rand),
						"services.#": "2",
						"force":      "true",
					}),
				),
			},
			{
				Config: testAccRamRoleMaxSessionDurationConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_session_duration": "7200",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRAMRole_multi(t *testing.T) {
	var v *ram.GetRoleResponse
	resourceId := "alicloud_ram_role.default.9"
	ra := resourceAttrInit(resourceId, ramRoleMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamRoleMultiConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccRamRoleCreateConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "default" {
	  name = "tf-testAcc%sRamRoleConfig-%d"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	}`, defaultRegionToTest, rand)
}

func TestAccAlicloudRAMRole_Description(t *testing.T) {
	var v *ram.GetRoleResponse
	resourceId := "alicloud_ram_role.default"
	ra := resourceAttrInit(resourceId, ramRoleMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamRoleCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "this is a test"}),
				),
			},
			{
				Config: testAccRamRoleUpdateDescription(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "this is a test_update"}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func testAccRamRoleUpdateDescription(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "default" {
	  name = "tf-testAcc%sRamRoleConfig-%d"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test_update"
	  force = true
	}`, defaultRegionToTest, rand)
}

func testAccRamRoleMaxSessionDurationConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "default" {
	  name = "tf-testAcc%sRamRoleConfig-%d"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	  max_session_duration = 7200
	}`, defaultRegionToTest, rand)
}

func testAccRamRoleNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "default" {
	  name = "tf-testAcc%sRamRoleConfig-%d-N"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	}`, defaultRegionToTest, rand)
}

func testAccRamRoleForceConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "default" {
	  name = "tf-testAcc%sRamRoleConfig-%d-N"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = false
	}`, defaultRegionToTest, rand)
}
func testAccRamRoleDocumentConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "default" {
	  name = "tf-testAcc%sRamRoleConfig-%d-N"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = false
	}`, defaultRegionToTest, rand)
}

func testAccRamRoleMultiConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "default" {
	  name = "tf-testAccRamRoleConfig-%d-${count.index}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	  count = 10
	}`, rand)
}

var ramRoleMap = map[string]string{
	"name":        CHECKSET,
	"services.#":  "2",
	"document":    CHECKSET,
	"description": "this is a test",
	"version":     "1",
	"force":       "true",
	"arn":         CHECKSET,
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

		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return WrapError(err)
		}
	}
	return nil
}
