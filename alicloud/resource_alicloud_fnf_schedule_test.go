package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_fnf_schedule", &resource.Sweeper{
		Name: "alicloud_fnf_schedule",
		F:    testSweepFnfSchedule,
	})
}

func testSweepFnfSchedule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "ListSchedules"
	request := make(map[string]interface{})
	var response map[string]interface{}
	conn, err := client.NewFnfClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fnf_schedules", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.Schedules", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Schedules", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["ScheduleName"].(string)
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Fnf Schedule: %s ", name)
			continue
		}
		log.Printf("[Info] Delete Fnf Schedule: %s", name)

		action := "DeleteSchedule"
		conn, err := client.NewFnfClient()
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"FlowName":     item["FlowName"],
			"ScheduleName": name,
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Fnf Schedule (%s): %s", name, err)
		}
	}
	return nil
}

func TestAccAlicloudFnfSchedule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fnf_schedule.default"
	ra := resourceAttrInit(resourceId, AlicloudFnfScheduleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &FnfService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFnfSchedule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudFnfSchedule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFnfScheduleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cron_expression": "30 9 * * * *",
					"flow_name":       "${alicloud_fnf_flow.default.name}",
					"schedule_name":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cron_expression": "30 9 * * * *",
						"flow_name":       CHECKSET,
						"schedule_name":   name,
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
					"cron_expression": "30 18 * * * *",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cron_expression": "30 18 * * * *",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testaccFnFSchedule813242",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccFnFSchedule813242",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": `false`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payload": `{\"tf-testchange\": \"test success\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payload": `{"tf-testchange": "test success"}`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cron_expression": "30 9 * * * *",
					"description":     "tf-testaccFnFSchedule983041",
					"enable":          `true`,
					"payload":         `{\"tf-test\": \"test success\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cron_expression": "30 9 * * * *",
						"description":     "tf-testaccFnFSchedule983041",
						"enable":          "true",
						"payload":         `{"tf-test": "test success"}`,
					}),
				),
			},
		},
	})
}

var AlicloudFnfScheduleMap0 = map[string]string{
	"enable":             "true",
	"last_modified_time": CHECKSET,
	"schedule_id":        CHECKSET,
}

func AlicloudFnfScheduleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_fnf_flow" "default" {
definition= "version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld"
description= "tf-testaccFnFFlow983041"
name = var.name
type= "FDL"
}
`, name)
}
