package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/actiontrail"
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
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	request := actiontrail.CreateDescribeTrailsRequest()
	raw, err := client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.DescribeTrails(request)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Trails: %s", err)
	}
	response := raw.(*actiontrail.DescribeTrailsResponse)

	swept := false

	for _, v := range response.TrailList {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Trail: %s", name)
			continue
		}
		swept = true

		log.Printf("[INFO] Deleting Trail: %s", name)

		request := actiontrail.CreateDeleteTrailRequest()
		request.Name = name

		_, err := client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
			return actiontrailClient.DeleteTrail(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete trail (%s): %s", name, err)
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudActiontrailTrail_basic(t *testing.T) {
	var v actiontrail.TrailListItem
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
