package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		"tf-testAcc",
		"tf_testAcc",
	}

	var groups []ram.GroupInListGroups
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

func TestAccAlicloudRAMGroup_basic(t *testing.T) {
	var v *ram.GetGroupResponse
	resourceId := "alicloud_ram_group.default"
	ra := resourceAttrInit(resourceId, ramGroupBasicMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sRamGroupConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRamGroupConfigDependence)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"name": name + "_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comments": "group comments",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comments": "group comments",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comments": "group comments update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comments": "group comments update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force": "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name":     fmt.Sprintf("tf-testAcc%sRamGroupConfig-%d", defaultRegionToTest, rand),
					"comments": "group comments",
					"force":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     fmt.Sprintf("tf-testAcc%sRamGroupConfig-%d", defaultRegionToTest, rand),
						"comments": "group comments",
						"force":    "false",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRAMGroup_multi(t *testing.T) {
	var v *ram.GetGroupResponse
	resourceId := "alicloud_ram_group.default.9"
	ra := resourceAttrInit(resourceId, ramGroupBasicMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sRamGroupConfig-%d-9", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRamGroupConfigDependence)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":  fmt.Sprintf("tf-testAcc%sRamGroupConfig-%d-${count.index}", defaultRegionToTest, rand),
					"count": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var ramGroupBasicMap = map[string]string{
	"comments": "",
	"force":    CHECKSET,
}

func resourceRamGroupConfigDependence(name string) string {
	return ""
}
