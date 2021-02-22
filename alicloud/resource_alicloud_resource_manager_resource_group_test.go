package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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
	prefixes := []string{
		"tf-rd",
		"tf-",
	}

	action := "ListResourceGroups"
	request := make(map[string]interface{})
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}

	var groupIds []string
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resoure manager group in service list: %s", err)
		}
		resp, err := jsonpath.Get("$.ResourceGroups.ResourceGroup", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ResourceGroups.ResourceGroup", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			// Skip Invalid resource group.
			if item["Status"].(string) != "OK" {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping resource manager group: %s ", item["Name"])
			} else {
				groupIds = append(groupIds, item["Id"].(string))
			}
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1

	}

	for _, groupId := range groupIds {
		log.Printf("[INFO] Delete resource manager group: %s ", groupId)

		request := map[string]interface{}{
			"ResourceGroupId": groupId,
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete resource manager group (%s): %s", groupId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerResourceGroup_basic(t *testing.T) {
	var v map[string]interface{}
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
