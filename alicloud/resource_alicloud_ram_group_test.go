package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_group", &resource.Sweeper{
		Name: "alicloud_ram_group",
		F:    testSweepRamGroups,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_ram_user",
		},
	})
}

func testSweepRamGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		fmt.Sprintf("tf-testAcc%s", defaultRegionToTest),
		fmt.Sprintf("tf_testAcc%s", defaultRegionToTest),
	}

	var groups []ram.Group
	request := ram.CreateListGroupsRequest()
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroups(request)
		})
		if err != nil {
			return WrapError(err)
		}
		resp, _ := raw.(*ram.ListGroupsResponse)
		if len(resp.Groups.Group) < 1 {
			break
		}
		groups = append(groups, resp.Groups.Group...)

		if !resp.IsTruncated {
			break
		}
		request.Marker = resp.Marker
	}
	sweeped := false

	for _, v := range groups {
		name := v.GroupName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ram Group: %s", name)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Ram Group: %s", name)
		request := ram.CreateListPoliciesForGroupRequest()
		request.GroupName = name

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram Group (%s): %s", name, err)
		}
		response, _ := raw.(*ram.ListPoliciesForGroupResponse)
		for _, p := range response.Policies.Policy {
			request := ram.CreateDetachPolicyFromGroupRequest()
			request.PolicyType = p.PolicyType
			request.GroupName = name
			request.PolicyName = p.PolicyName
			log.Printf("[INFO] Detaching Ram policy %s from group: %s", p.PolicyName, name)
			_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.DetachPolicyFromGroup(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to detach policy from Group (%s): %s", name, err)
			}
		}
		_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateDeleteGroupRequest()
			request.GroupName = name
			return ramClient.DeleteGroup(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ram Group (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudRamGroup_basic(t *testing.T) {
	var v *ram.GetGroupResponse
	resourceId := "alicloud_ram_group.default"
	ra := resourceAttrInit(resourceId, ramGroupBasicMap)
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
		CheckDestroy:  testAccCheckRamGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamGroupNameConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": fmt.Sprintf("tf-testAcc%sRamGroupConfig-%d", defaultRegionToTest, rand)}),
				),
			},
			{
				Config: testAccRamGroupConmmentsConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"comments": "group comments"}),
				),
			},
			{
				Config: testAccRamGroupForceConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"force": "true"}),
				),
			},

			{
				Config: testAccRamGroupAllConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     fmt.Sprintf("tf-testAcc%sRamGroupConfig-all-%d", defaultRegionToTest, rand),
						"comments": "group comments all",
						"force":    "false",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRamGroup_multi(t *testing.T) {
	var v *ram.GetGroupResponse
	resourceId := "alicloud_ram_group.default.9"
	ra := resourceAttrInit(resourceId, nil)
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
		CheckDestroy:  testAccCheckRamGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamGroupMultiConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     fmt.Sprintf("tf-testAcc%sRamGroupConfig-%d-9", defaultRegionToTest, rand),
						"comments": "group comments",
						"force":    "false",
					}),
				),
			},
		},
	})
}

var ramGroupBasicMap = map[string]string{
	"comments": "",
	"force":    CHECKSET,
}

func testAccRamGroupNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "default" {
	  name = "tf-testAcc%sRamGroupConfig-%d"
	}`, defaultRegionToTest, rand)
}

func testAccRamGroupConmmentsConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "default" {
	  name = "tf-testAcc%sRamGroupConfig-%d"
	  comments = "group comments"
	}`, defaultRegionToTest, rand)
}

func testAccRamGroupForceConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "default" {
	  name = "tf-testAcc%sRamGroupConfig-%d"
	  comments = "group comments"
	  force=true
	}`, defaultRegionToTest, rand)
}

func testAccRamGroupAllConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "default" {
	  name = "tf-testAcc%sRamGroupConfig-all-%d"
	  comments = "group comments all"
	  force=false
	}`, defaultRegionToTest, rand)
}

func testAccRamGroupMultiConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "default" {
	  name = "tf-testAcc%sRamGroupConfig-%d-${count.index}"
	  comments = "group comments"
	  count=10
	}`, defaultRegionToTest, rand)
}

func testAccCheckRamGroupExists(n string, group *ram.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Group ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetGroupRequest()
		request.GroupName = rs.Primary.ID
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetGroup(request)
		})

		if err == nil {
			response, _ := raw.(*ram.GetGroupResponse)
			*group = response.Group
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckRamGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group" {
			continue
		}

		// Try to find the group
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetGroupRequest()
		request.GroupName = rs.Primary.ID

		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetGroup(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
	}
	return nil
}
