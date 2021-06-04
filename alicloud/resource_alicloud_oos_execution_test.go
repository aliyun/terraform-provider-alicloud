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
	resource.AddTestSweepers("alicloud_oos_execution", &resource.Sweeper{
		Name: "alicloud_oos_execution",
		F:    testSweepOosExecution,
	})
}

func testSweepOosExecution(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	action := "ListExecutions"
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "testSweepOosExecution", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Executions", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Executions", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			name := item["TemplateName"]
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name.(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping OOS Execution: %s", item["ExecutionId"].(string))
				continue
			}
			sweeped = true
			action = "DeleteExecutions"
			request := map[string]interface{}{
				"ExecutionIds": convertListToJsonString(convertListStringToListInterface([]string{item["ExecutionId"].(string)})),
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete OOS Execution (%s): %s", item["ExecutionId"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure OOS Execution have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete OOS Execution success: %s ", item["ExecutionId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudOosExecution_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_execution.default"
	ra := resourceAttrInit(resourceId, OosExecutionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosExecution")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosExecution%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, OosExecutionBasicdependence)
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
					"template_name": "${alicloud_oos_template.default.template_name}",
					"description":   "From TF Test",
					"parameters":    `{\"Status\":\"Running\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name": CHECKSET,
						"description":   "From TF Test",
						"parameters":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "loop_mode", "safety_check", "parameters"},
			},
		},
	})
}

var OosExecutionMap = map[string]string{
	"create_date": CHECKSET,
	"executed_by": CHECKSET,
	"is_parent":   CHECKSET,
	"start_date":  CHECKSET,
	"status":      CHECKSET,
	"template_id": CHECKSET,
	"update_date": CHECKSET,
}

func OosExecutionBasicdependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_oos_template" "default" {
		  content= <<EOF
		  {
			"FormatVersion": "OOS-2019-06-01",
			"Description": "Describe instances of given status",
			"Parameters":{
			  "Status":{
				"Type": "String",
				"Description": "(Required) The status of the Ecs instance."
			  }
			},
			"Tasks": [
			  {
				"Properties" :{
				  "Parameters":{
					"Status": "{{ Status }}"
				  },
				  "API": "DescribeInstances",
				  "Service": "Ecs"
				},
				"Name": "foo",
				"Action": "ACS::ExecuteApi"
			  }]
		  }
		  EOF
		  template_name = "%s"
		  version_name = "test"
		  tags = {
			"Created" = "TF",
			"For" = "template Test"
		  }
		}
	`, name)
}
