package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_resource_manager_resource_group", &resource.Sweeper{
		Name: "alicloud_resource_manager_resource_group",
		F:    testSweepResourceManagerResourceGroup,
	})
}

func testSweepResourceManagerResourceGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}

	prefixes := []string{
		"tf-rd",
		"tf-",
	}

	var groupIds []string
	request := resourcemanager.CreateListResourceGroupsRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.ListResourceGroups(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resoure manager group in service list: %s", err)
		}
		response, _ := raw.(*resourcemanager.ListResourceGroupsResponse)

		for _, v := range response.ResourceGroups.ResourceGroup {
			// Skip Invalid resource group.
			if v.Status != "OK" {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(v.Name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping resource manager group: %s ", v.Name)
			} else {
				groupIds = append(groupIds, v.Id)
			}
		}

		if len(response.ResourceGroups.ResourceGroup) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	for _, groupId := range groupIds {
		log.Printf("[INFO] Delete resource manager group: %s ", groupId)

		request := resourcemanager.CreateDeleteResourceGroupRequest()
		request.ResourceGroupId = groupId

		_, err = resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.DeleteResourceGroup(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete resource manager group (%s): %s", groupId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerResourceGroup_basic(t *testing.T) {
	var v resourcemanager.ResourceGroup
	resourceId := "alicloud_resource_manager_resource_group.default"
	ra := resourceAttrInit(resourceId, ResourceManagerResourceGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-rd%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerResourceGroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "terraform-test",
					"name":         name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "terraform-test",
						"name":         name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "terraform-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "terraform-test",
					}),
				),
			},
		},
	})
}

var ResourceManagerResourceGroupMap = map[string]string{}

func ResourceManagerResourceGroupBasicdependence(name string) string {
	return ""
}
