package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
		"tf-testacc",
		"tf_testacc",
	}

	var users []ram.UserInListUsers
	request := ram.CreateListUsersRequest()
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsers(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram User: %s", err)
		}
		response, _ := raw.(*ram.ListUsersResponse)
		if len(response.Users.User) < 1 {
			break
		}
		users = append(users, response.Users.User...)

		if !response.IsTruncated {
			break
		}
		request.Marker = response.Marker
	}
	sweeped := false

	for _, v := range users {
		name := v.UserName
		id := v.UserId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), prefix) {
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
			listPoliciesForUserRequest := ram.CreateListPoliciesForUserRequest()
			listPoliciesForUserRequest.UserName = name
			return ramClient.ListPoliciesForUser(listPoliciesForUserRequest)
		})
		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
			log.Printf("[ERROR] ListPoliciesForUser: %s (%s)", name, id)
		}
		listPoliciesForUserResponse, _ := raw.(*ram.ListPoliciesForUserResponse)
		if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
			detachPolicyFromUserRequest := ram.CreateDetachPolicyFromUserRequest()
			detachPolicyFromUserRequest.UserName = name

			for _, poloicy := range listPoliciesForUserResponse.Policies.Policy {
				detachPolicyFromUserRequest.PolicyName = poloicy.PolicyName
				detachPolicyFromUserRequest.PolicyType = poloicy.PolicyType
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromUser(detachPolicyFromUserRequest)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					log.Printf("[ERROR] DetachPolicyFromUser: %s (%s)", name, id)
				}
			}
		}

		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			listAccessKeysRequest := ram.CreateListAccessKeysRequest()
			listAccessKeysRequest.UserName = name
			return ramClient.ListAccessKeys(listAccessKeysRequest)
		})
		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
			log.Printf("[ERROR] ListAccessKeys: %s (%s)", name, id)
		}
		listAccessKeysResponse, _ := raw.(*ram.ListAccessKeysResponse)
		if len(listAccessKeysResponse.AccessKeys.AccessKey) > 0 {
			deleteAccessKeyRequest := ram.CreateDeleteAccessKeyRequest()
			deleteAccessKeyRequest.UserName = name

			for _, accesskey := range listAccessKeysResponse.AccessKeys.AccessKey {
				deleteAccessKeyRequest.UserName = name
				deleteAccessKeyRequest.UserAccessKeyId = accesskey.AccessKeyId
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DeleteAccessKey(deleteAccessKeyRequest)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					log.Printf("[ERROR] ListAccessKeysResponse: %s (%s)", name, id)
				}
			}
		}

		listGroupsForUserRequest := ram.CreateListGroupsForUserRequest()
		listGroupsForUserRequest.UserName = name
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroupsForUser(listGroupsForUserRequest)
		})
		if err != nil {
			log.Printf("[ERROR] ListGroupsForUser: %s (%s)", name, id)
		}
		listGroupsForUserResponse, _ := raw.(*ram.ListGroupsForUserResponse)
		if len(listGroupsForUserResponse.Groups.Group) > 0 {
			for _, v := range listGroupsForUserResponse.Groups.Group {
				removeUserFromGroupRequest := ram.CreateRemoveUserFromGroupRequest()
				removeUserFromGroupRequest.UserName = name
				removeUserFromGroupRequest.GroupName = v.GroupName
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(removeUserFromGroupRequest)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					log.Printf("[ERROR] RemoveUserFromGroup: %s (%s)", name, id)
				}
			}
		}
		log.Printf("[INFO] Deleting Ram User: %s (%s)", name, id)
		deleteUserRequest := ram.CreateDeleteUserRequest()
		deleteUserRequest.UserName = name

		_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteUser(deleteUserRequest)
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

func TestAccAlicloudRAMUser(t *testing.T) {
	var v *ram.UserInGetUser
	randInt := acctest.RandIntRange(1000000, 99999999)

	resourceId := "alicloud_ram_user.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamUserNameConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAcc%sRamUserConfig-%d", defaultRegionToTest, randInt),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
			{
				Config: testAccRamUserDisplayNameConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "displayname",
					}),
				),
			},
			{
				Config: testAccRamUserMobileConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile": "86-18888888888",
					}),
				),
			},
			{
				Config: testAccRamUserEmailConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "hello.uuu@aaa.com",
					}),
				),
			},
			{
				Config: testAccRamUserCommentsConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comments": "yoyoyo",
					}),
				),
			},
			{
				Config: testAccRamUserConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testAcc%sRamUserConfig-%d_all", defaultRegionToTest, randInt),
						"display_name": "displayname_all",
						"mobile":       "86-18888888889",
						"email":        "hello.uuu@aaa_all.com",
						"comments":     "yoyoyo_all",
					}),
				),
			},
		},
	})
}

func testAccCheckRamUserExists(n string, user *ram.UserInGetUser) resource.TestCheckFunc {
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

		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return WrapError(err)
		}
	}
	return nil
}

func testAccRamUserNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamUserConfig-%d"
	}`, defaultRegionToTest, rand)
}

func testAccRamUserDisplayNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamUserConfig-%d"
	  display_name = "displayname"
	}`, defaultRegionToTest, rand)
}

func testAccRamUserMobileConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamUserConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	}`, defaultRegionToTest, rand)
}

func testAccRamUserEmailConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamUserConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	}`, defaultRegionToTest, rand)
}

func testAccRamUserCommentsConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamUserConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}`, defaultRegionToTest, rand)
}

func testAccRamUserConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamUserConfig-%d_all"
	  display_name = "displayname_all"
	  mobile = "86-18888888889"
	  email = "hello.uuu@aaa_all.com"
	  comments = "yoyoyo_all"
	}`, defaultRegionToTest, rand)
}
