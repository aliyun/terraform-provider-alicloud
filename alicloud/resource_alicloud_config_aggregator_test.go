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
	resource.AddTestSweepers("alicloud_config_aggregator", &resource.Sweeper{
		Name: "alicloud_config_aggregator",
		F:    testSweepConfigAggregator,
		Dependencies: []string{
			"alicloud_config_aggregate_compliance_pack",
			"alicloud_config_aggregate_config_rule",
		},
	})
}

func testSweepConfigAggregator(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}

	// Get all AggregatorId
	aggregatorIds := make([]string, 0)
	action := "ListAggregators"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-09-07"), StringPointer("AK"), request, nil, &runtime)
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
			log.Println("List Config Aggregator Failed!", err)
		}
		resp, err := jsonpath.Get("$.AggregatorsResult.Aggregators", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AggregatorsResult.Aggregators", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["AggregatorName"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Aggregate: %v (%v)", item["AggregatorName"], item["AggregatorId"])
				continue
			}
			aggregatorIds = append(aggregatorIds, fmt.Sprint(item["AggregatorId"]))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	log.Printf("[INFO] Deleting Aggregate:  (%s)", strings.Join(aggregatorIds, ","))
	action = "DeleteAggregators"
	deleteRequest := map[string]interface{}{
		"AggregatorIds": strings.Join(aggregatorIds, ","),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, deleteRequest, &util.RuntimeOptions{})
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
		log.Printf("[INFO] Delete Aggregate Failed:  (%s)", strings.Join(aggregatorIds, ","))
	}
	return nil
}

func TestAccAlicloudConfigAggregator_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregatorMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregatorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.0.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.0.display_name}",
							"account_type": "ResourceDirectory",
						},
					},
					"aggregator_name": "${var.name}",
					"description":     "tf-create-aggregator",
					"aggregator_type": "CUSTOM",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_accounts.#": "1",
						"aggregator_name":       name,
						"description":           "tf-create-aggregator",
						"aggregator_type":       "CUSTOM",
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
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.0.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.0.display_name}",
							"account_type": "ResourceDirectory",
						},
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.1.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.1.display_name}",
							"account_type": "ResourceDirectory",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_accounts.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-modify-aggregator",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-modify-aggregator",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.0.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.0.display_name}",
							"account_type": "ResourceDirectory",
						},
					},
					"aggregator_name": "${var.name}",
					"description":     "tf-create-aggregator",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_accounts.#": "1",
						"aggregator_name":       name,
						"description":           "tf-create-aggregator",
					}),
				),
			},
		},
	})
}

var AlicloudConfigAggregatorMap0 = map[string]string{
	"aggregator_type": "CUSTOM",
}

func AlicloudConfigAggregatorBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}

`, name)
}
