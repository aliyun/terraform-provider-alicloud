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
	resource.AddTestSweepers(
		"alicloud_threat_detection_anti_brute_force_rule",
		&resource.Sweeper{
			Name: "alicloud_threat_detection_anti_brute_force_rule",
			F:    testSweepThreatDetectionAntiBruteForceRule,
		})
}

func testSweepThreatDetectionAntiBruteForceRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{"tf-testAcc", "tf_testAcc"}
	action := "DescribeAntiBruteForceRules"
	query := make(map[string]interface{})
	request := make(map[string]interface{})
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, e := client.RpcPost("Sas", "2018-12-03", action, query, request, true)
		response = resp
		if e != nil {
			if NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch Threat Detection AntiBruteForceRules: %s", WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_anti_brute_force_rule", action, AlibabaCloudSdkGoERROR))
		return nil
	}
	v, err := jsonpath.Get("$.Rules", response)
	if err != nil {
		log.Printf("[WARN] DescribeAntiBruteForceRules returned no Rules (empty account?): %v", err)
		return nil
	}
	result, ok := v.([]interface{})
	if !ok || len(result) == 0 {
		log.Printf("[INFO] DescribeAntiBruteForceRules: 0 rules in region %s", client.RegionId)
		return nil
	}
	for _, item := range result {
		rule := item.(map[string]interface{})
		ruleId := fmt.Sprint(rule["Id"])
		name := fmt.Sprint(rule["Name"])
		// System-default rules (DefaultRule == true) are cloud-managed and must
		// not be removed by the sweeper in either NORMAL or SWEEP_ALL mode.
		isDefault := false
		switch dv := rule["DefaultRule"].(type) {
		case bool:
			isDefault = dv
		case string:
			isDefault = strings.EqualFold(dv, "true")
		case float64:
			isDefault = dv != 0
		}
		if isDefault {
			log.Printf("[INFO] Skipping Threat Detection AntiBruteForceRule (system default): %s (%s)", name, ruleId)
			continue
		}
		if protectedSweepVpcName(name) {
			log.Printf("[INFO] Skipping Threat Detection AntiBruteForceRule (NODELETING): %s (%s)", name, ruleId)
			continue
		}
		if !sweepAll() {
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Threat Detection AntiBruteForceRule (no prefix): %s (%s)", name, ruleId)
				continue
			}
		}
		log.Printf("[INFO] Deleting Threat Detection AntiBruteForceRule: %s (%s)", name, ruleId)
		deleteAction := "DeleteAntiBruteForceRule"
		deleteQuery := make(map[string]interface{})
		deleteRequest := make(map[string]interface{})
		deleteRequest["Ids.1"] = ruleId
		if _, e := client.RpcPost("Sas", "2018-12-03", deleteAction, deleteQuery, deleteRequest, true); e != nil {
			log.Printf("[ERROR] Failed to delete Threat Detection AntiBruteForceRule (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete Threat Detection AntiBruteForceRule success: %s (%s)", name, ruleId)
	}
	return nil
}

// Test ThreatDetection AntiBruteForceRule. >>> Resource test cases, automatically generated.
// Case 对接Terraform 7780
func TestAccAliCloudThreatDetectionAntiBruteForceRule_basic7780(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_anti_brute_force_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAntiBruteForceRuleMap7780)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAntiBruteForceRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAntiBruteForceRuleBasicDependence7780)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"anti_brute_force_rule_name": name,
					"forbidden_time":             "360",
					"uuid_list": []string{
						"032b618f-b220-4a0d-bd37-fbdc6ef58b6a"},
					"fail_count": "80",
					"span":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": name,
						"forbidden_time":             "360",
						"uuid_list.#":                "1",
						"fail_count":                 "80",
						"span":                       "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"anti_brute_force_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forbidden_time": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forbidden_time": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fail_count": "70",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fail_count": "70",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"span": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"span": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"uuid_list": []string{"inet-eae024ff-c25b-4168-b839-f6ecc67e9db4", "1d343293-52db-4684-a5e1-5feb1e0670fa"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uuid_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_rule": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type": []map[string]interface{}{
						{
							"rdp":        "off",
							"ssh":        "off",
							"sql_server": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type.#": "1",
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

func TestAccAliCloudThreatDetectionAntiBruteForceRule_basic7780_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_anti_brute_force_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAntiBruteForceRuleMap7780)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAntiBruteForceRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAntiBruteForceRuleBasicDependence7780)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"anti_brute_force_rule_name": name,
					"forbidden_time":             "360",
					"uuid_list": []string{
						"032b618f-b220-4a0d-bd37-fbdc6ef58b6a"},
					"fail_count":   "80",
					"span":         "10",
					"default_rule": "true",
					"protocol_type": []map[string]interface{}{
						{
							"rdp":        "off",
							"ssh":        "off",
							"sql_server": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": name,
						"forbidden_time":             "360",
						"uuid_list.#":                "1",
						"fail_count":                 "80",
						"span":                       "10",
						"default_rule":               "true",
						"protocol_type.#":            "1",
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

var AliCloudThreatDetectionAntiBruteForceRuleMap7780 = map[string]string{}

func AliCloudThreatDetectionAntiBruteForceRuleBasicDependence7780(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ThreatDetection AntiBruteForceRule. <<< Resource test cases, automatically generated.
