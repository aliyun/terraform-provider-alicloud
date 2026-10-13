package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccEsaAigwRecordPreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	if os.Getenv("ALICLOUD_ESA_AIGW_INSTANCE_ID") != "" {
		return
	}
	rawClient, clientErr := sharedClientForRegion(os.Getenv("ALICLOUD_REGION"))
	if clientErr != nil {
		t.Fatalf("Failed to get AliCloud client for ESA AIGW record test: %v", clientErr)
	}
	client := rawClient.(*connectivity.AliyunClient)
	name := fmt.Sprintf("tf-testacc-aigw-%d", acctest.RandIntRange(10000, 99999))
	createAction := "CreateAIGWInstance"
	createReq := map[string]interface{}{
		"Name": name,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var resp map[string]interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = client.RpcPost("ESA", "2024-09-10", createAction, map[string]interface{}{}, createReq, true)
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
		t.Fatalf("Failed to create AIGW instance for test: %v", err)
	}
	addDebug(createAction, resp, createReq)

	instanceId := findEsaAigwInstanceIdByName(client, name)
	if instanceId == "" {
		t.Fatalf("Failed to find created AIGW instance named %s", name)
	}
	os.Setenv("ALICLOUD_ESA_AIGW_INSTANCE_ID", instanceId)
}

func findEsaAigwInstanceIdByName(client *connectivity.AliyunClient, name string) string {
	action := "ListAIGWInstances"
	request := map[string]interface{}{
		"PageSize":       50,
		"PageNumber":     1,
		"FuzzySearchKey": name,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var resp map[string]interface{}
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, err = client.RpcGet("ESA", "2024-09-10", action, request, nil)
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
		return ""
	}
	addDebug(action, resp, request)
	instances, _ := jsonpath.Get("$.Instances", resp)
	instanceList, _ := instances.([]interface{})
	for _, v := range instanceList {
		item := v.(map[string]interface{})
		instanceId := fmt.Sprint(firstNonEmpty(item["InstanceId"], item["AIGWInstanceId"]))
		itemName := fmt.Sprint(firstNonEmpty(item["Name"], item["AIGWInstanceName"]))
		if itemName == name {
			return instanceId
		}
	}
	if len(instanceList) > 0 {
		item := instanceList[0].(map[string]interface{})
		return fmt.Sprint(firstNonEmpty(item["InstanceId"], item["AIGWInstanceId"]))
	}
	return ""
}

func firstNonEmpty(values ...interface{}) interface{} {
	for _, v := range values {
		if v != nil && fmt.Sprint(v) != "" {
			return v
		}
	}
	return ""
}

func testAccEsaAigwRecordDestroyInstance() {
	instanceId := os.Getenv("ALICLOUD_ESA_AIGW_INSTANCE_ID")
	if instanceId == "" {
		return
	}
	rawClient, clientErr := sharedClientForRegion(os.Getenv("ALICLOUD_REGION"))
	if clientErr != nil {
		log.Printf("[ERROR] Failed to get AliCloud client for destroying AIGW instance: %s", clientErr)
		return
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "DeleteAIGWInstance"
	query := map[string]interface{}{
		"InstanceId": instanceId,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var resp map[string]interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = client.RpcPost("ESA", "2024-09-10", action, query, map[string]interface{}{}, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"20101", "NotFound"}) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, resp, query)
	os.Unsetenv("ALICLOUD_ESA_AIGW_INSTANCE_ID")
}

func testAccEsaAigwRecordBasicDependence(name string) string {
	instanceId := os.Getenv("ALICLOUD_ESA_AIGW_INSTANCE_ID")
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "aigw_instance_id" {
  default = "%s"
}

resource "alicloud_esa_rate_plan_instance" "default" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "default" {
  site_name   = "tf-testacc-aigw-%[1]s.com"
  instance_id = alicloud_esa_rate_plan_instance.default.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name, instanceId)
}

func TestAccAliCloudEsaAigwRecord_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_aigw_record.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaAigwRecordBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaAigwRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEsaAigwRecordBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccEsaAigwRecordPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${alicloud_esa_site.default.id}",
					"instance_id": "${var.aigw_instance_id}",
					"record_name": fmt.Sprintf("aigw-record-%d.example.com", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
	testAccEsaAigwRecordDestroyInstance()
}

var AliCloudEsaAigwRecordBasicMap = map[string]string{
	"site_id":     "${alicloud_esa_site.default.id}",
	"instance_id": "${var.aigw_instance_id}",
	"record_name": "CHECKSET",
}
