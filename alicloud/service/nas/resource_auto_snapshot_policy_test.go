package nas_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_nas_auto_snapshot_policy",
		&resource.Sweeper{
			Name: "alicloud_nas_auto_snapshot_policy",
			F:    testSweepNasAutoSnapshotPolicy,
		})
}

func testSweepNasAutoSnapshotPolicy(region string) error {
	if nasTestSweepPreCheckWithRegions(region, true, connectivity.NASSupportRegions) {
		log.Printf("[INFO] Skipping Nas Auto Snapshot Policy unsupported region: %s", region)
		return nil
	}

	rawClient, err := nasTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeAutoSnapshotPolicies"
	request := map[string]interface{}{}
	request["FileSystemType"] = "extreme"
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.AutoSnapshotPolicies.AutoSnapshotPolicy", action, err)
			return nil
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
				log.Printf("[INFO] Skipping Nas Auto Snapshot Policy: %s", item["AutoSnapshotPolicyName"].(string))
				continue
			}
			action := "DeleteAutoSnapshotPolicy"
			request := map[string]interface{}{
				"AutoSnapshotPolicyId": item["AutoSnapshotPolicyId"],
			}
			_, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Nas Auto Snapshot Policy (%s): %s", item["AutoSnapshotPolicyName"].(string), err)
			}
			log.Printf("[INFO] Delete Nas Auto Snapshot Policy success: %s ", item["AutoSnapshotPolicyName"].(string))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudNASAutoSnapshotPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_auto_snapshot_policy.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNASAutoSnapshotPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasAutoSnapshotPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-nasautosnapshotpolicy%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNASAutoSnapshotPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"2", "3", "4"},
					"time_points":     []string{"0", "1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "3",
						"time_points.#":     "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_snapshot_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"3", "4", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_points": []string{"1", "2", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_days": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "30",
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

func TestAccAliCloudNASAutoSnapshotPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_auto_snapshot_policy.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNASAutoSnapshotPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasAutoSnapshotPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-nasautosnapshotpolicy%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNASAutoSnapshotPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_name": "${var.name}",
					"repeat_weekdays":           []string{"3", "4"},
					"time_points":               []string{"1", "2"},
					"retention_days":            "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_snapshot_policy_name": name,
						"repeat_weekdays.#":         "2",
						"time_points.#":             "2",
						"retention_days":            "30",
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

var AlicloudNASAutoSnapshotPolicyMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudNASAutoSnapshotPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
var AlicloudNasAutoSnapshotPolicyMap6532 = map[string]string{
	"status":         acctest.CHECKSET,
	"create_time":    acctest.CHECKSET,
	"retention_days": acctest.CHECKSET,
}

func AlicloudNasAutoSnapshotPolicyBasicDependence6532(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Case 自动快照策略 6532  raw
func TestAccAliCloudNasAutoSnapshotPolicy_basic6532_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_auto_snapshot_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNasAutoSnapshotPolicyMap6532)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasAutoSnapshotPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasautosnapshotpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNasAutoSnapshotPolicyBasicDependence6532)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"time_points":               []string{"0", "1", "2"},
					"retention_days":            "1",
					"repeat_weekdays":           []string{"2", "3", "4"},
					"auto_snapshot_policy_name": name,
					"file_system_type":          "extreme",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#":             "3",
						"retention_days":            "1",
						"repeat_weekdays.#":         "3",
						"auto_snapshot_policy_name": name,
						"file_system_type":          "extreme",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_points":               []string{"0", "5"},
					"retention_days":            "2",
					"repeat_weekdays":           []string{"4", "5", "6", "7"},
					"auto_snapshot_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#":             "2",
						"retention_days":            "2",
						"repeat_weekdays.#":         "4",
						"auto_snapshot_policy_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Nas AutoSnapshotPolicy. <<< Resource test cases, automatically generated.
