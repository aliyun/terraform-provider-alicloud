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
		"alicloud_actiontrail_history_delivery_job",
		&resource.Sweeper{
			Name: "alicloud_actiontrail_history_delivery_job",
			F:    testSweepActionTrailHistoryDeliveryJob,
		})
}

func testSweepActionTrailHistoryDeliveryJob(region string) error {
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.ActiontrailSupportRegions)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListDeliveryHistoryJobs"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewActiontrailClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &runtime)
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
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.DeliveryHistoryJobs", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DeliveryHistoryJobs", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["TrailName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if fmt.Sprint(item["JobStatus"]) == "0" || fmt.Sprint(item["JobStatus"]) == "1" {
				skip = true
			}
			if skip {
				log.Printf("[INFO] Skipping Delivery History Job: %s", item["TrailName"].(string))
				continue
			}
			action := "DeleteDeliveryHistoryJob"
			request := map[string]interface{}{
				"JobId": item["JobId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Delivery History Job (%s): %s", item["TrailName"].(string), err)
			}
			log.Printf("[INFO] Delete Delivery History Job success: %s ", item["TrailName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudActionTrailHistoryDeliveryJob_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ActiontrailSupportRegions)
	resourceId := "alicloud_actiontrail_history_delivery_job.default"
	ra := resourceAttrInit(resourceId, AlicloudActionTrailHistoryDeliveryJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailHistoryDeliveryJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccactiontrailjob%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudActionTrailHistoryDeliveryJobBasicDependence0)
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
					"trail_name": "${alicloud_actiontrail_trail.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"trail_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudActionTrailHistoryDeliveryJobMap0 = map[string]string{}

func AlicloudActionTrailHistoryDeliveryJobBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_account" "default" {}

resource "alicloud_log_project" "default" {
  name = var.name
  description = "tf actiontrail test"
}

data "alicloud_ram_roles" "default" {
  name_regex = "AliyunActionTrailDefaultRole"
}

resource "alicloud_actiontrail_trail" "default" {
  trail_name = var.name
  sls_write_role_arn = data.alicloud_ram_roles.default.roles.0.arn
  sls_project_arn = "acs:log:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:project/${alicloud_log_project.default.name}"
}
`, name)
}
