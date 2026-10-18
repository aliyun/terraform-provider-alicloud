// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
	resource.AddTestSweepers("alicloud_cr_scan_rule", &resource.Sweeper{
		Name: "alicloud_cr_scan_rule",
		F:    testSweepCrScanRule,
	})
}

// testSweepCrScanRule cleans up Container Registry image scan rules left on CR
// EE instances by acceptance tests. Rules on instances (or with rule names)
// marked NODELETING are preserved; in non-SWEEP_ALL mode only rules on tf-*
// named instances are removed.
func testSweepCrScanRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}

	prefixes := []string{"tf-testAcc", "tf_testAcc", "tftestAcc"}

	pageNo, pageSize := 1, PageSizeLarge
	for {
		resp, lerr := crService.ListCrEEInstances(pageNo, pageSize)
		if lerr != nil {
			log.Printf("[ERROR] Failed to ListCrEEInstances: %s", lerr)
			return nil
		}
		for _, inst := range resp.Instances {
			if protectedSweepVpcName(inst.InstanceName) {
				log.Printf("[INFO] Skipping CR scan rule on instance %s (%s): protected (NODELETING)", inst.InstanceId, inst.InstanceName)
				continue
			}
			if !sweepAll() {
				match := false
				for _, p := range prefixes {
					if strings.HasPrefix(strings.ToLower(inst.InstanceName), strings.ToLower(p)) {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			listAction := "ListScanRule"
			listQuery := map[string]interface{}{
				"RegionId":   client.RegionId,
				"InstanceId": inst.InstanceId,
				"PageSize":   pageSize,
				"PageNo":     1,
			}
			for {
				wait := incrementalWait(3*time.Second, 3*time.Second)
				var response map[string]interface{}
				rerr := retry.Retry(5*time.Minute, func() *retry.RetryError {
					response, err = client.RpcGet("cr", "2018-12-01", listAction, listQuery, nil)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if rerr != nil {
					log.Printf("[ERROR] Failed to ListScanRule for instance %s: %s", inst.InstanceId, rerr)
					break
				}
				raw, jerr := jsonpath.Get("$.ScanRules", response)
				if jerr != nil {
					log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", listAction, "$.ScanRules", response)
					break
				}
				result, _ := raw.([]interface{})
				for _, v := range result {
					item := v.(map[string]interface{})
					ruleId := fmt.Sprint(item["ScanRuleId"])
					ruleName := fmt.Sprint(item["RuleName"])
					if protectedSweepVpcName(ruleName) {
						log.Printf("[INFO] Skipping CR scan rule %s:%s: protected (NODELETING)", inst.InstanceId, ruleName)
						continue
					}
					log.Printf("[INFO] Deleting CR scan rule: %s:%s (%s)", inst.InstanceId, ruleId, ruleName)
					delReq := map[string]interface{}{
						"RegionId":   client.RegionId,
						"InstanceId": inst.InstanceId,
						"ScanRuleId": ruleId,
					}
					dwait := incrementalWait(3*time.Second, 3*time.Second)
					if derr := retry.Retry(5*time.Minute, func() *retry.RetryError {
						_, e := client.RpcPost("cr", "2018-12-01", "DeleteScanRule", map[string]interface{}{}, delReq, true)
						if e != nil {
							if IsExpectedErrors(e, []string{"SCAN_RULE_NOT_EXIST"}) || NotFoundError(e) {
								return nil
							}
							if NeedRetry(e) {
								dwait()
								return retry.RetryableError(e)
							}
							return retry.NonRetryableError(e)
						}
						return nil
					}); derr != nil {
						log.Printf("[ERROR] Failed to delete CR scan rule %s:%s: %s", inst.InstanceId, ruleId, derr)
						continue
					}
				}
				if len(result) < pageSize {
					break
				}
				listQuery["PageNo"] = listQuery["PageNo"].(int) + 1
			}
		}
		if len(resp.Instances) < pageSize {
			break
		}
		pageNo++
	}
	return nil
}

// Test Cr ScanRule. >>> Resource test cases, automatically generated.
// Case ScanRule-1_pl 11745
func TestAccAliCloudCrScanRule_basic11745(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_scan_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCrScanRuleMap11745)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrScanRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrScanRuleBasicDependence11745)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": ".*",
					"scan_scope":              "REPO",
					"trigger_type":            "MANUAL",
					"scan_type":               "VUL",
					"rule_name":               "302",
					"namespaces": []string{
						"aa"},
					"repo_names": []string{
						"bb", "cc", "dd", "ee"},
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": ".*",
						"scan_scope":              "REPO",
						"trigger_type":            "MANUAL",
						"scan_type":               "VUL",
						"rule_name":               CHECKSET,
						"namespaces.#":            "1",
						"repo_names.#":            "4",
						"instance_id":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "aab",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "aab",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{
						"cc"},
					"repo_names": []string{
						"bb", "cc", "dd"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespaces.#": "1",
						"repo_names.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{
						"aa"},
					"repo_names": []string{
						"bb", "cc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespaces.#": "1",
						"repo_names.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{
						"bb"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespaces.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_names": []string{
						"bb", "aa"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_names.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"trigger_type": "MANUAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"trigger_type": "MANUAL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "cc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "cc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "a",
					"scan_scope":              "NAMESPACE",
					"trigger_type":            "AUTO",
					"rule_name":               "aa",
					"namespaces": []string{
						"aa", "bb", "cc", "dd"},
					"repo_names": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "a",
						"scan_scope":              "NAMESPACE",
						"trigger_type":            "AUTO",
						"rule_name":               "aa",
						"namespaces.#":            "4",
						"repo_names.#":            "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "bb",
					"scan_scope":              "REPO",
					"trigger_type":            "MANUAL",
					"rule_name":               "bb",
					"namespaces": []string{
						"aa"},
					"repo_names": []string{
						"aa"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "bb",
						"scan_scope":              "REPO",
						"trigger_type":            "MANUAL",
						"rule_name":               "bb",
						"namespaces.#":            "1",
						"repo_names.#":            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "cc",
					"scan_scope":              "INSTANCE",
					"trigger_type":            "AUTO",
					"rule_name":               "dd",
					"namespaces":              []string{},
					"repo_names":              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "cc",
						"scan_scope":              "INSTANCE",
						"trigger_type":            "AUTO",
						"rule_name":               "dd",
						"namespaces.#":            "0",
						"repo_names.#":            "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "dd",
					"scan_scope":              "NAMESPACE",
					"trigger_type":            "MANUAL",
					"rule_name":               "ff",
					"namespaces": []string{
						"aa", "bb"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "dd",
						"scan_scope":              "NAMESPACE",
						"trigger_type":            "MANUAL",
						"rule_name":               "ff",
						"namespaces.#":            "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "aa",
					"scan_scope":              "INSTANCE",
					"trigger_type":            "AUTO",
					"rule_name":               "gg",
					"namespaces":              []string{},
					"repo_names":              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "aa",
						"scan_scope":              "INSTANCE",
						"trigger_type":            "AUTO",
						"rule_name":               "gg",
						"namespaces.#":            "0",
						"repo_names.#":            "0",
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

var AlicloudCrScanRuleMap11745 = map[string]string{
	"scan_rule_id": CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudCrScanRuleBasicDependence11745(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 1
  renewal_status = "AutoRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}
`, name)
}

// Test Cr ScanRule. <<< Resource test cases, automatically generated.
