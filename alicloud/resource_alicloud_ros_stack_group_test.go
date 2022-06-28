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
	resource.AddTestSweepers(
		"alicloud_ros_stack_group",
		&resource.Sweeper{
			Name: "alicloud_ros_stack_group",
			F:    testSweepRosStackGroup,
		})
}

func testSweepRosStackGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   region,
		"Status":     "ACTIVE",
	}
	var response map[string]interface{}
	action := "ListStackGroups"
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ros_stack_group", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.StackGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.StackGroups", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["StackGroupName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ros StackGroup: %s", item["StackGroupName"].(string))
				continue
			}
			sweeped = true
			action := "DeleteStackGroup"
			request := map[string]interface{}{
				"StackGroupName": item["StackGroupName"],
				"RegionId":       region,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ros StackGroup (%s): %s", item["StackGroupName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Ros StackGroup have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Ros StackGroup success: %s ", item["StackGroupName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudROSStackGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRosStackGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosStackGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudRosStackGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosStackGroupBasicDependence)
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
					"stack_group_name": name,
					"template_body":    `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "VpcName",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "InstanceType",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_group_name": name,
						"template_body":    CHECKSET,
						"parameters.#":     "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_ids", "operation_description", "template_body", "operation_preferences", "region_ids", "template_url", "template_version"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_body": `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Description\" : \"模板描述信息，可用于说明模板的适用场景、架构说明等。\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "tf-testacc",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "ECS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_body": CHECKSET,
						"parameters.#":  "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_body": `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "VpcName",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "InstanceType",
						},
					},
					"description": "test for tf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_body": CHECKSET,
						"parameters.#":  "2",
						"description":   "test for tf",
					}),
				),
			},
		},
	})
}

var AlicloudRosStackGroupMap = map[string]string{
	"stack_group_id": CHECKSET,
	"status":         CHECKSET,
}

func AlicloudRosStackGroupBasicDependence(name string) string {
	return ""
}
