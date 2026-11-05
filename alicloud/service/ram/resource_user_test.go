package ram_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	rawClient, err := ramTestClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	var users []ramsdk.User
	request := ramsdk.CreateListUsersRequest()
	for {
		raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.ListUsers(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram User: %s", err)
		}
		response, _ := raw.(*ramsdk.ListUsersResponse)
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
		raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			listPoliciesForUserRequest := ramsdk.CreateListPoliciesForUserRequest()
			listPoliciesForUserRequest.UserName = name
			return ramClient.ListPoliciesForUser(listPoliciesForUserRequest)
		})
		if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
			log.Printf("[ERROR] ListPoliciesForUser: %s (%s)", name, id)
		}
		listPoliciesForUserResponse, _ := raw.(*ramsdk.ListPoliciesForUserResponse)
		if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
			detachPolicyFromUserRequest := ramsdk.CreateDetachPolicyFromUserRequest()
			detachPolicyFromUserRequest.UserName = name

			for _, poloicy := range listPoliciesForUserResponse.Policies.Policy {
				detachPolicyFromUserRequest.PolicyName = poloicy.PolicyName
				detachPolicyFromUserRequest.PolicyType = poloicy.PolicyType
				_, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromUser(detachPolicyFromUserRequest)
				})
				if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
					log.Printf("[ERROR] DetachPolicyFromUser: %s (%s)", name, id)
				}
			}
		}

		raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			listAccessKeysRequest := ramsdk.CreateListAccessKeysRequest()
			listAccessKeysRequest.UserName = name
			return ramClient.ListAccessKeys(listAccessKeysRequest)
		})
		if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
			log.Printf("[ERROR] ListAccessKeys: %s (%s)", name, id)
		}
		listAccessKeysResponse, _ := raw.(*ramsdk.ListAccessKeysResponse)
		if len(listAccessKeysResponse.AccessKeys.AccessKey) > 0 {
			deleteAccessKeyRequest := ramsdk.CreateDeleteAccessKeyRequest()
			deleteAccessKeyRequest.UserName = name

			for _, accesskey := range listAccessKeysResponse.AccessKeys.AccessKey {
				deleteAccessKeyRequest.UserName = name
				deleteAccessKeyRequest.UserAccessKeyId = accesskey.AccessKeyId
				_, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
					return ramClient.DeleteAccessKey(deleteAccessKeyRequest)
				})
				if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
					log.Printf("[ERROR] ListAccessKeysResponse: %s (%s)", name, id)
				}
			}
		}

		listGroupsForUserRequest := ramsdk.CreateListGroupsForUserRequest()
		listGroupsForUserRequest.UserName = name
		raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.ListGroupsForUser(listGroupsForUserRequest)
		})
		if err != nil {
			log.Printf("[ERROR] ListGroupsForUser: %s (%s)", name, id)
		}
		listGroupsForUserResponse, _ := raw.(*ramsdk.ListGroupsForUserResponse)
		if len(listGroupsForUserResponse.Groups.Group) > 0 {
			for _, v := range listGroupsForUserResponse.Groups.Group {
				removeUserFromGroupRequest := ramsdk.CreateRemoveUserFromGroupRequest()
				removeUserFromGroupRequest.UserName = name
				removeUserFromGroupRequest.GroupName = v.GroupName
				_, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(removeUserFromGroupRequest)
				})
				if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
					log.Printf("[ERROR] RemoveUserFromGroup: %s (%s)", name, id)
				}
			}
		}
		log.Printf("[INFO] Deleting Ram User: %s (%s)", name, id)
		deleteUserRequest := ramsdk.CreateDeleteUserRequest()
		deleteUserRequest.UserName = name

		_, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
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

func TestAccAliCloudRAMUser(t *testing.T) {
	var v *ramsdk.User
	randInt := tfacctest.RandIntRange(1000000, 99999999)

	resourceId := "alicloud_ram_user.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return ram.NewRamService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamUserNameConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAcc%sRamUserConfig-%d", acctest.DefaultRegionToTest, randInt),
					}),
				),
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
				Config: testAccRamUserEmailConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comments": acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccRamUserMobileConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccRamUserDisplayNameConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile": acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccRamUserNameConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccRamUserConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testAcc%sRamUserConfig-%d_all", acctest.DefaultRegionToTest, randInt),
						"display_name": "displayname_all",
						"mobile":       "86-18888888889",
						"email":        "hello.uuu@aaa_all.com",
						"comments":     "yoyoyo_all",
					}),
				),
			},
			{
				Config: testAccRamUserNameConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testAcc%sRamUserConfig-%d", acctest.DefaultRegionToTest, randInt),
						"display_name": acctest.REMOVEKEY,
						"mobile":       acctest.REMOVEKEY,
						"email":        acctest.REMOVEKEY,
						"comments":     acctest.REMOVEKEY,
					}),
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

func testAccCheckRamUserExists(n string, user *ramsdk.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return errs.WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return errs.WrapError(fmt.Errorf("No user ID is set"))
		}

		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		request := ramsdk.CreateGetUserRequest()
		request.UserName = rs.Primary.Attributes["name"]

		raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.GetUser(request)
		})
		log.Printf("[WARN] User id %#v", rs.Primary.ID)
		if err == nil {
			response, _ := raw.(*ramsdk.GetUserResponse)
			*user = response.User
			return nil
		}
		return errs.WrapError(err)
	}
}

func testAccCheckRamUserDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user" {
			continue
		}

		// Try to find the user
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		request := ramsdk.CreateGetUserRequest()
		request.UserName = rs.Primary.Attributes["name"]

		_, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.GetUser(request)
		})

		if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return errs.WrapError(err)
		}
	}
	return nil
}

func testAccRamUserNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
  		name = "tf-testAcc%sRamUserConfig-%d"
	}
`, acctest.DefaultRegionToTest, rand)
}

func testAccRamUserDisplayNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
  		name         = "tf-testAcc%sRamUserConfig-%d"
  		display_name = "displayname"
	}
`, acctest.DefaultRegionToTest, rand)
}

func testAccRamUserMobileConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
  		name         = "tf-testAcc%sRamUserConfig-%d"
  		display_name = "displayname"
  		mobile       = "86-18888888888"
	}
`, acctest.DefaultRegionToTest, rand)
}

func testAccRamUserEmailConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
  		name         = "tf-testAcc%sRamUserConfig-%d"
  		display_name = "displayname"
  		mobile       = "86-18888888888"
  		email        = "hello.uuu@aaa.com"
	}
`, acctest.DefaultRegionToTest, rand)
}

func testAccRamUserCommentsConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
  		name         = "tf-testAcc%sRamUserConfig-%d"
  		display_name = "displayname"
  		mobile       = "86-18888888888"
  		email        = "hello.uuu@aaa.com"
  		comments     = "yoyoyo"
	}
`, acctest.DefaultRegionToTest, rand)
}

func testAccRamUserConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
  		name         = "tf-testAcc%sRamUserConfig-%d_all"
  		display_name = "displayname_all"
  		mobile       = "86-18888888889"
  		email        = "hello.uuu@aaa_all.com"
  		comments     = "yoyoyo_all"
	}
`, acctest.DefaultRegionToTest, rand)
}

// TestAccAliCloudRAMUser_tags covers the tags attribute, whose CRUD is
// managed via the IMS (2019-08-15) tagging API with ResourceType=user,
// because the legacy RAM (2015-05-01) TagResources API does not enumerate
// "user" as a taggable resource type.
func TestAccAliCloudRAMUser_tags(t *testing.T) {
	var v *ramsdk.User
	randInt := tfacctest.RandIntRange(1000000, 99999999)

	resourceId := "alicloud_ram_user.default"
	name := fmt.Sprintf("tf-testAcc%sRamUserTags-%d", acctest.DefaultRegionToTest, randInt)
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return ram.NewRamService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, func(string) string { return "" })

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: testAccCheckRamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":  name,
					"force": true,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"force":        "true",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
						"New":     "Tag",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
						"tags.New":     "Tag",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.Created": "TF",
						"tags.For":     acctest.REMOVEKEY,
						"tags.New":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
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
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
		},
	})
}
