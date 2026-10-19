package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("alicloud_ocean_base_instance", &resource.Sweeper{
		Name: "alicloud_ocean_base_instance",
		F:    testSweepOceanBaseInstances,
	})
}

// testSweepOceanBaseInstances cleans up OceanBase instances created by acceptance tests.
// Subscription (PrePay) instances cannot be destroyed via the OpenAPI and are skipped
// (recorded as leftover). Name-based NODELETING protection is always applied, even under
// SWEEP_ALL, by reusing protectedSweepVpcName from the VPC sweeper helpers.
func testSweepOceanBaseInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
		"tf_testAcc",
	}

	action := "DescribeInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	objects := make([]interface{}, 0)
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, err := client.RpcPost("OceanBasePro", "2019-09-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		resp, jsonErr := jsonpath.Get("$.Instances", response)
		if jsonErr != nil {
			log.Printf("[ERROR] DescribeInstances parse $.Instances failed: %v", jsonErr)
			break
		}
		result, _ := resp.([]interface{})
		objects = append(objects, result...)
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, v := range objects {
		item := v.(map[string]interface{})
		id := fmt.Sprint(item["InstanceId"])
		name := fmt.Sprint(item["InstanceName"])
		payType := fmt.Sprint(item["PayType"])

		if protectedSweepVpcName(name) {
			log.Printf("[INFO] Skipping OceanBase Instance (NODELETING marker): %s (%s)", name, id)
			continue
		}

		if !sweepAll() {
			matched := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					matched = true
					break
				}
			}
			if !matched {
				log.Printf("[INFO] Skipping OceanBase Instance (prefix mismatch): %s (%s)", name, id)
				continue
			}
		}

		// Subscription (PrePay) instances cannot be destroyed via DeleteInstances; the
		// resource Delete path refuses them as well. Skip and record as leftover.
		if payType == "PrePay" || payType == "PREPAY" {
			log.Printf("[INFO] Skipping OceanBase Instance (Subscription cannot be destroyed): %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting OceanBase Instance: %s (%s)", name, id)
		deleteAction := "DeleteInstances"
		deleteRequest := map[string]interface{}{
			"InstanceIds":       fmt.Sprintf("[\"%s\"]", id),
			"BackupRetainMode":  "delete_all",
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		deleteErr := retry.Retry(10*time.Minute, func() *retry.RetryError {
			resp, err := client.RpcPost("OceanBasePro", "2019-09-01", deleteAction, nil, deleteRequest, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"UnknownError", "IllegalOperation.Resource"}) {
					return nil
				}
				return retry.NonRetryableError(err)
			}
			addDebug(deleteAction, resp, deleteRequest)
			return nil
		})
		if deleteErr != nil {
			log.Printf("[ERROR] Failed to delete OceanBase Instance (%s (%s)): %s", name, id, deleteErr)
		}
	}

	return nil
}

func TestAccAliCloudOceanBaseInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	checkoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OceanBaseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOceanBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOceanBaseInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":  "${var.name}",
					"series":         "normal",
					"auto_renew":     "false",
					"disk_size":      "100",
					"payment_type":   "PayAsYouGo",
					"instance_class": "8C32GB",
					//"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"zones":              []string{"${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"},
					"backup_retain_mode": "delete_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name,
						"series":         "normal",
						"auto_renew":     "false",
						"disk_size":      "100",
						"instance_class": CHECKSET,
						//"resource_group_id": CHECKSET,
						"zones.#":      "3",
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_num": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_num": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size":      "200",
					"instance_class": "14C70GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size":      "200",
						"instance_class": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size":      "100",
					"instance_class": "8C32G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size":      "100",
						"instance_class": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "auto_renew_period", "backup_retain_mode", "period"},
			},
		},
	})
}

func TestAccAliCloudOceanBaseInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	checkoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OceanBaseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOceanBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOceanBaseInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":  "${var.name}",
					"series":         "normal",
					"auto_renew":     "false",
					"disk_size":      "100",
					"payment_type":   "PayAsYouGo",
					"instance_class": "8C32GB",
					//"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"zones":              []string{"${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"},
					"backup_retain_mode": "delete_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name,
						"series":         "normal",
						"auto_renew":     "false",
						"disk_size":      "100",
						"instance_class": CHECKSET,
						//"resource_group_id": CHECKSET,
						"zones.#":      "3",
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "auto_renew_period", "backup_retain_mode", "period"},
			},
		},
	})
}

func TestAccAliCloudOceanBaseInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	checkoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OceanBaseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOceanBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOceanBaseInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"period":            "1",
					"period_unit":       "Month",
					"auto_renew_period": "1",
					"instance_name":     "${var.name}",
					"series":            "normal",
					"auto_renew":        "false",
					"disk_size":         "300",
					"disk_type":         "cloud_essd_pl1",
					"ob_version":        "4.2.1",
					"payment_type":      "Subscription",
					"instance_class":    "4C16G",
					//"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"zones":              []string{"${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"},
					"backup_retain_mode": "delete_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name,
						"series":         "normal",
						"auto_renew":     "false",
						"disk_size":      "300",
						"disk_type":      "cloud_essd_pl1",
						"ob_version":     "4.2.1",
						"instance_class": CHECKSET,
						//"resource_group_id": CHECKSET,
						"zones.#":      "3",
						"payment_type": "Subscription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "auto_renew_period", "backup_retain_mode", "period"},
			},
		},
	})
}

var AlicloudOceanBaseInstanceMap0 = map[string]string{
	"status":  CHECKSET,
	"zones.#": CHECKSET,
}

func AlicloudOceanBaseInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {}


data "alicloud_resource_manager_resource_groups" "default"{}

`, name)
}

// Test OceanBase Instance. >>> Resource test cases, automatically generated.
// Case tf_0822_qianyu_zhubei 7602
func TestAccAliCloudOceanBaseInstance_basic7602(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap7602)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OceanBaseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOceanBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soceanbaseinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOceanBaseInstanceBasicDependence7602)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "4C16G",
					"zones": []string{
						"cn-hangzhou-j", "cn-hangzhou-i"},
					"instance_name":    name,
					"disk_type":        "cloud_essd_pl1",
					"series":           "normal",
					"disk_size":        "100",
					"payment_type":     "PayAsYouGo",
					"ob_version":       "4.2.1",
					"node_num":         "3",
					"cpu_arch":         "X86",
					"primary_region":   "cn-hangzhou",
					"primary_instance": "${alicloud_ocean_base_instance.createZhu.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":   CHECKSET,
						"zones.#":          "2",
						"instance_name":    name,
						"disk_type":        "cloud_essd_pl1",
						"series":           "normal",
						"disk_size":        "100",
						"payment_type":     "PayAsYouGo",
						"ob_version":       "4.2.1",
						"node_num":         "3",
						"cpu_arch":         "X86",
						"primary_region":   CHECKSET,
						"primary_instance": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":      "8C32GB",
					"instance_name":       name + "_update",
					"disk_size":           "200",
					"upgrade_spec_native": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":      CHECKSET,
						"instance_name":       name + "_update",
						"disk_size":           "200",
						"upgrade_spec_native": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"node_num":      "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
						"node_num":      "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":  name + "_update",
					"disk_size":      "100",
					"instance_class": "4C16G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name + "_update",
						"disk_size":      "100",
						"instance_class": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "period", "period_unit", "upgrade_spec_native"},
			},
		},
	})
}

var AlicloudOceanBaseInstanceMap7602 = map[string]string{
	"node_num":       CHECKSET,
	"cpu":            CHECKSET,
	"disk_type":      CHECKSET,
	"commodity_code": CHECKSET,
	"ob_version":     CHECKSET,
	"status":         CHECKSET,
	"create_time":    CHECKSET,
	"instance_name":  CHECKSET,
}

func AlicloudOceanBaseInstanceBasicDependence7602(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ocean_base_instance" "createZhu" {
  instance_class = "4C16G"
  zones          = ["cn-hangzhou-j", "cn-hangzhou-i"]
  instance_name  = var.name
  disk_type      = "cloud_essd_pl1"
  series         = "normal"
  disk_size      = "100"
  payment_type   = "PayAsYouGo"
  ob_version     = "4.2.1"
  period_unit    = "Hour"
  node_num       = "3"
}


`, name)
}

// Test OceanBase Instance. <<< Resource test cases, automatically generated.
