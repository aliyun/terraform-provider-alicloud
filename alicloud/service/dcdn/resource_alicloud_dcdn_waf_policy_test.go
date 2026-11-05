package dcdn_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dcdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dcdn_waf_policy",
		&resource.Sweeper{
			Name: "alicloud_dcdn_waf_policy",
			F:    testSweepDcdnWafPolicy,
		})
}

func testSweepDcdnWafPolicy(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DCDNSupportRegions) {
		log.Printf("[INFO] Skipping Dcdn Waf Policy unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDcdnWafPolicies"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId

	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, false)
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

		resp, err := jsonpath.Get("$.Policies", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Policies", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["PolicyName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Dcdn Waf Policy: %s", item["PolicyName"].(string))
				continue
			}
			action := "DeleteDcdnWafPolicy"
			request := map[string]interface{}{
				"PolicyId": item["PolicyId"],
			}
			_, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Dcdn Waf Policy (%s): %s", item["PolicyName"].(string), err)
			}
			log.Printf("[INFO] Delete Dcdn Waf Policy success: %s ", item["PolicyName"].(string))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudDcdnWafPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_waf_policy.default"
	checkoutSupportedRegions(t, true, connectivity.DCDNSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDCDNWafPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnWafPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccdcdnwafpolicy%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDCDNWafPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_type":   "custom",
					"policy_name":   "${var.name}",
					"defense_scene": "waf_group",
					"status":        "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_type":   "custom",
						"policy_name":   name,
						"defense_scene": "waf_group",
						"status":        "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "off",
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

var AlicloudDCDNWafPolicyMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudDCDNWafPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
