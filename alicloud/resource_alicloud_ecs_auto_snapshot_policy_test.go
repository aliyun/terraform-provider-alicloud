package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_auto_snapshot_policy", &resource.Sweeper{
		Name: "alicloud_ecs_auto_snapshot_policy",
		F:    testSweepEcsAutoSnapshotPolicy,
	})
}

func testSweepEcsAutoSnapshotPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId
	for {
		action := "DescribeAutoSnapshotPolicyEx"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_auto_snapshot_policy", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["AutoSnapshotPolicyName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecs SnapShot Policy: %s (%s)", item["AutoSnapshotPolicyName"], item["AutoSnapshotPolicyId"])
				continue
			}
			action = "DeleteAutoSnapshotPolicy"
			request := map[string]interface{}{
				"autoSnapshotPolicyId": item["AutoSnapshotPolicyId"],
				"regionId":             client.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs SnapShot Policy (%s (%s)): %s", item["AutoSnapshotPolicyName"], item["InstanceId"], err)
				continue
			}
			log.Printf("[INFO] Delete Ecs SnapShot Policy success: %s ", item["AutoSnapshotPolicyId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudEcsAutoSnapshotPolicyBasic(t *testing.T) {

	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	randInt := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccSnapshotPolicyBasic%d", randInt)
	basicMap := map[string]string{
		"name":              name,
		"repeat_weekdays.#": "1",
		"retention_days":    "-1",
		"time_points.#":     "1",
	}
	var v map[string]interface{}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_weekdays": []string{"1"},
					"retention_days":  "-1",
					"time_points":     []string{"1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_days": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_points": []string{"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_cross_region_copy":        "true",
					"target_copy_regions":             []string{"cn-beijing"},
					"copied_snapshots_retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_cross_region_copy":        "true",
						"target_copy_regions.#":           "1",
						"copied_snapshots_retention_days": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                            name,
					"repeat_weekdays":                 []string{"1"},
					"retention_days":                  "-1",
					"time_points":                     []string{"1"},
					"enable_cross_region_copy":        "false",
					"target_copy_regions":             []string{"cn-shanghai"},
					"copied_snapshots_retention_days": "-1",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                            name,
						"repeat_weekdays.#":               "1",
						"retention_days":                  "-1",
						"time_points.#":                   "1",
						"enable_cross_region_copy":        "false",
						"target_copy_regions.#":           "1",
						"copied_snapshots_retention_days": "-1",
						"tags.%":                          "2",
						"tags.Created":                    "TF-update",
						"tags.For":                        "Test-update",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSnapshotPolicyMulti(t *testing.T) {

	resourceId := "alicloud_snapshot_policy.default.4"
	randInt := acctest.RandIntRange(10000, 99999)
	var v map[string]interface{}
	name := fmt.Sprintf("tf-testAccSnapshotPolicyMulti%d", randInt)
	basicMap := map[string]string{
		"name":              name,
		"repeat_weekdays.#": "1",
		"retention_days":    "-1",
		"time_points.#":     "1",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":           "5",
					"name":            name,
					"repeat_weekdays": []string{"1"},
					"retention_days":  "-1",
					"time_points":     []string{"1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudEcsAutoSnapshotPolicyBasic1(t *testing.T) {

	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	randInt := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccSnapshotPolicyBasic%d", randInt)
	basicMap := map[string]string{
		"name":              name,
		"repeat_weekdays.#": "1",
		"retention_days":    "-1",
		"time_points.#":     "1",
	}
	var v map[string]interface{}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                            name,
					"repeat_weekdays":                 []string{"1"},
					"retention_days":                  "-1",
					"time_points":                     []string{"1"},
					"copied_snapshots_retention_days": "2",
					"enable_cross_region_copy":        "true",
					"target_copy_regions":             []string{"cn-beijing"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                            name,
						"repeat_weekdays.#":               "1",
						"retention_days":                  "-1",
						"time_points.#":                   "1",
						"copied_snapshots_retention_days": "2",
						"enable_cross_region_copy":        "true",
						"target_copy_regions.#":           "1",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
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
