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
		"alicloud_actiontrail_trail",
		&resource.Sweeper{
			Name: "alicloud_actiontrail_trail",
			F:    testSweepActiontrailTrail,
		})
}

func testSweepActiontrailTrail(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	action := "DescribeTrails"
	conn, err := client.NewActiontrailClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-04"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_actiontrail_trails", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.TrailList", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TrailList", response)
	}
	sweeped := false
	for _, v := range resp.([]interface{}) {
		item := v.(map[string]interface{})
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ActionTrail Trails: %s", item["Name"].(string))
			continue
		}
		sweeped = true
		action = "DeleteTrail"
		request := map[string]interface{}{
			"Name": item["Name"],
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-04"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete ActionTrail Trail (%s): %s", item["Name"].(string), err)
		}
		if sweeped {
			// Waiting 5 seconds to ensure these ActionTrail Trails have been deleted.
			time.Sleep(5 * time.Second)
		}
		log.Printf("[INFO] Delete ActionTrail Trail success: %s ", item["Name"].(string))
	}
	return nil
}

func TestAccAlicloudActiontrailTrail_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_actiontrail_trail.default"
	ra := resourceAttrInit(resourceId, ActiontrailTrailMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailTrail")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccactiontrail%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ActiontrailTrailBasicdependence)
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
					"trail_name":      name,
					"role_name":       "aliyunactiontraildefaultrole",
					"oss_bucket_name": "${alicloud_oss_bucket.default.id}",
					"status":          "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"trail_name":      name,
						"role_name":       "aliyunactiontraildefaultrole",
						"oss_bucket_name": name,
						"status":          "Disable",
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
					"role_name": "aliyunserviceroleforactiontrail",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name": "aliyunserviceroleforactiontrail",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"event_rw": "All",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_rw": "All",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_bucket_name": "${alicloud_oss_bucket.default2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_bucket_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"trail_region": "cn-beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"trail_region": "cn-beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_bucket_name": "${alicloud_oss_bucket.default.id}",
					"role_name":       "aliyunactiontraildefaultrole",
					"trail_region":    "All",
					"event_rw":        "Write",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_bucket_name": name,
						"role_name":       "aliyunactiontraildefaultrole",
						"trail_region":    "All",
						"event_rw":        "Write",
					}),
				),
			},
		},
	})
}

var ActiontrailTrailMap = map[string]string{}

func ActiontrailTrailBasicdependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_oss_bucket" "default" {
		bucket  = "${var.name}"
	}

	resource "alicloud_oss_bucket" "default2" {
		bucket  = "${var.name}-update"
	}

`, name)
}
