package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_bp_studio_application", &resource.Sweeper{
		Name: "alicloud_bp_studio_application",
		F:    testSweepBpStudioApplication,
	})
}

func testSweepBpStudioApplication(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	bpStudioService := BpStudioService{client}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListApplication"
	request := make(map[string]interface{})
	request["MaxResults"] = PageSizeLarge
	request["NextToken"] = 1
	var response map[string]interface{}
	BpStudioApplicationIds := make([]string, 0)
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("BPStudio", "2021-09-31", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_bp_studio_application", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["Name"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping BpStudioApplication Instance: %v", item["ApplicationId"])
					continue
				}
			}
			BpStudioApplicationIds = append(BpStudioApplicationIds, fmt.Sprint(item["ApplicationId"]))
		}
		if len(result) < request["MaxResults"].(int) {
			break
		}
		request["NextToken"] = request["NextToken"].(int) + 1
	}

	for _, id := range BpStudioApplicationIds {
		log.Printf("[INFO] Deleting BpStudioApplication Instance: %s", id)
		releaseAction := "ReleaseApplication"
		deleteAction := "DeleteApplication"
		if err != nil {
			return WrapError(err)
		}
		request = map[string]interface{}{
			"ApplicationId": id,
		}
		object, err := bpStudioService.DescribeBpStudioApplication(id)
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)

		if fmt.Sprint(object["Status"]) == "Deployed_Failure" || fmt.Sprint(object["Status"]) == "PartiallyDeployedSuccess" || fmt.Sprint(object["Status"]) == "Deployed_Success" || fmt.Sprint(object["Status"]) == "Destroyed_Failure" || fmt.Sprint(object["Status"]) == "PartiallyDestroyedSuccess" {
			err = resource.Retry(120*time.Minute, func() *resource.RetryError {
				_, err = client.RpcPost("BPStudio", "2021-09-31", releaseAction, nil, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			releaseApplicationStateConf := BuildStateConf([]string{}, []string{"Destroyed_Success"}, 120*time.Minute, 5*time.Second, bpStudioService.BpStudioApplicationStateRefreshFunc(id, []string{}))
			if _, err := releaseApplicationStateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, id)
			}
		}

		err = resource.Retry(120*time.Minute, func() *resource.RetryError {
			_, err = client.RpcPost("BPStudio", "2021-09-31", deleteAction, nil, request, true)
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
			log.Printf("[ERROR] Failed to delete BpStudioApplication Instance (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAlicloudBpStudioApplication_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bp_studio_application.default"
	checkoutSupportedRegions(t, true, connectivity.BpStudioApplicationSupportRegions)
	ra := resourceAttrInit(resourceId, resourceAlicloudBpStudioApplicationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BpStudioService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBpStudioApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBpStudioApplication-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudBpStudioApplicationBasicDependence)
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
					"application_name":  name,
					"template_id":       "YAUUQIYRSV1CMFGX",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"area_id":           defaultRegionToTest,
					"instances": []map[string]string{
						{
							"id":        "${data.alicloud_instances.default.instances.0.id}",
							"node_name": "${data.alicloud_instances.default.instances.0.name}",
							"node_type": "ecs",
						},
					},
					"configuration": map[string]string{
						"enableMonitor": "1",
					},
					"variables": map[string]string{
						"test": "1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name":  name,
						"template_id":       "YAUUQIYRSV1CMFGX",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"area_id", "instances", "configuration", "variables"},
			},
		},
	})
}

var resourceAlicloudBpStudioApplicationMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudBpStudioApplicationBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_instances" "default" {
		status = "Running"
	}
`, name)
}
