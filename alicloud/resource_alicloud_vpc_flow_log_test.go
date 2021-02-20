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
	resource.AddTestSweepers("alicloud_vpc_flow_log", &resource.Sweeper{
		Name: "alicloud_vpc_flow_log",
		F:    testSweepVpcFlowLog,
	})
}

func testSweepVpcFlowLog(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}
	client := rawClient.(*connectivity.AliyunClient)
	request := make(map[string]interface{}, 0)
	action := "DescribeFlowLogs"
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	var instances []string
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve VpcFlowLog service list: %s", err)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.FlowLogs.FlowLog", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.FlowLogs.FlowLog", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if v, ok := item["FlowLogId"]; !ok || v.(string) == "" {
				continue
			}
			instances = append(instances, fmt.Sprint(item["FlowLogId"], ":", item["FlowLogName"]))
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, instance := range instances {
		instanceId := strings.Split(instance, ":")[0]
		instanceName := strings.Split(instance, ":")[1]
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping vpc_flow_log: %s ", instanceId)
			continue
		}

		action := "DeleteFlowLog"
		var response map[string]interface{}
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"FlowLogId": instanceId,
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*10, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve vpc_flow_log %s %v", instanceId, err)
			continue
		}
		log.Printf("[INFO] Delete vpc_flow_log instance: %s ", instanceId)
	}
	return nil
}

func TestAccAlicloudVpcFlowLog_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_flow_log.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcFlowLogMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFlowLogs")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccVpcFlowLog-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcFlowLogBasicDependence)
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
					"resource_id":    "${alicloud_vpc.default.id}",
					"resource_type":  "VPC",
					"traffic_type":   "All",
					"log_store_name": "${var.log_store_name}",
					"project_name":   "${var.project_name}",
					"flow_log_name":  name,
					"description":    "test",
					"status":         "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         "Active",
						"resource_type":  "VPC",
						"traffic_type":   "All",
						"resource_id":    CHECKSET,
						"project_name":   "vpc-flow-log-for-vpc",
						"log_store_name": "vpc-flow-log-for-vpc",
						"flow_log_name":  name,
						"description":    "test",
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
					"flow_log_name": name + "change",
					"description":   name + "change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flow_log_name": name + "change",
						"description":   name + "change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "Active",
					"flow_log_name": name + "change_again",
					"description":   name + "change_again",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":        "Active",
						"flow_log_name": name + "change_again",
						"description":   name + "change_again",
					}),
				),
			},
		},
	})
}

var AlicloudVpcFlowLogMap = map[string]string{
	"status": "Active",
}

func AlicloudVpcFlowLogBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

variable "log_store_name" {
  default = "vpc-flow-log-for-vpc"
}

variable "project_name" {
  default = "vpc-flow-log-for-vpc"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  name       = var.name
}

`, name)
}
