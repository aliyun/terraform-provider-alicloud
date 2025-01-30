package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_group_metric_rule", &resource.Sweeper{
		Name: "alicloud_cms_group_metric_rule",
		F:    testSweepCmsGroupMetricRule,
	})
}

func testSweepCmsGroupMetricRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeMetricRuleList"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["Page"] = 1
	var response map[string]interface{}
	for {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_group_metric_rules", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Alarms.Alarm", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Alarms.Alarm", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["RuleName"].(string)
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Cms Metric Rule: %s ", name)
					continue
				}
			}
			log.Printf("[INFO] Delete Cms Metric Rule: %s ", name)

			delAction := "DeleteMetricRules"
			delRequest := map[string]interface{}{
				"Id": []string{item["RuleId"].(string)},
			}

			_, err = client.RpcPost("Cms", "2019-01-01", delAction, nil, delRequest, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Cms Metric Rule (%s): %s", name, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["Page"] = request["Page"].(int) + 1
	}
	return nil
}

func TestAccAliCloudCmsGroupMetricRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_group_metric_rule.default"
	ra := resourceAttrInit(resourceId, resourceAliCloudCmsGroupMetricRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsGroupMetricRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCmsGroupMetricRuletf-testacc-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAliCloudCmsGroupMetricRuleBasicDependence)
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
					"rule_id":                "4a9a8978-a9cc-55ca-aa7c-530ccd91ae57",
					"group_id":               "${alicloud_cms_monitor_group.default.id}",
					"group_metric_rule_name": name,
					"metric_name":            "disk_writebytes",
					"namespace":              "acs_ecs_dashboard",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanOrEqualToThreshold",
									"statistics":          "Average",
									"threshold":           "90",
									"times":               "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_id":                "4a9a8978-a9cc-55ca-aa7c-530ccd91ae57",
						"group_id":               CHECKSET,
						"group_metric_rule_name": name,
						"metric_name":            "disk_writebytes",
						"namespace":              "acs_ecs_dashboard",
						"escalations.#":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "${alicloud_cms_monitor_group.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_metric_rule_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_metric_rule_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"metric_name": "diskusage_used",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_name": "diskusage_used",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups": "${alicloud_cms_monitor_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dimensions": `{\"device\":\"C:\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dimensions": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email_subject": "tf-testacc-rule-name-warning-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email_subject": "tf-testacc-rule-name-warning-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-22:59",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-22:59",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_effective_interval": "00:00-06:30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_effective_interval": "00:00-06:30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "180",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "85800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": "85800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook": "https://www.aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook": "https://www.aliyun.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"id":          "1",
							"json_params": `{\"a\":\"b\"}`,
							"level":       "Warn",
							"arn":         "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"statistics":          "Maximum",
									"threshold":           "10",
									"times":               "2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations": []map[string]interface{}{
						{
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanOrEqualToThreshold",
									"statistics":          "Minimum",
									"threshold":           "20",
									"times":               "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations": []map[string]interface{}{
						{
							"warn": []map[string]interface{}{
								{
									"comparison_operator": "LessThanOrEqualToThreshold",
									"statistics":          "Average",
									"threshold":           "30",
									"times":               "5",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"category", "interval"},
			},
		},
	})
}

func TestAccAliCloudCmsGroupMetricRule_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_group_metric_rule.default"
	ra := resourceAttrInit(resourceId, resourceAliCloudCmsGroupMetricRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsGroupMetricRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCmsGroupMetricRuletf-testacc-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAliCloudCmsGroupMetricRuleBasicDependence)
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
					"rule_id":                "4a9a8978-a9cc-55ca-aa7c-530ccd91ae57",
					"group_id":               "${alicloud_cms_monitor_group.default.id}",
					"group_metric_rule_name": name,
					"metric_name":            "disk_writebytes",
					"namespace":              "acs_ecs_dashboard",
					"category":               "ecs",
					"contact_groups":         "${alicloud_cms_monitor_group.default.id}",
					"dimensions":             `{\"device\":\"C:\"}`,
					"email_subject":          "tf-testacc-rule-name-warning",
					"effective_interval":     "00:00-22:59",
					"no_effective_interval":  "00:00-06:30",
					"interval":               "60",
					"period":                 "180",
					"silence_time":           "85800",
					"webhook":                "https://www.aliyun.com",
					"targets": []map[string]interface{}{
						{
							"id":          "1",
							"json_params": `{\"a\":\"b\"}`,
							"level":       "Warn",
							"arn":         "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
						},
					},
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanOrEqualToThreshold",
									"statistics":          "Average",
									"threshold":           "90",
									"times":               "1",
								},
							},
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanOrEqualToThreshold",
									"statistics":          "Minimum",
									"threshold":           "20",
									"times":               "3",
								},
							},
							"warn": []map[string]interface{}{
								{
									"comparison_operator": "LessThanOrEqualToThreshold",
									"statistics":          "Average",
									"threshold":           "30",
									"times":               "5",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_id":                "4a9a8978-a9cc-55ca-aa7c-530ccd91ae57",
						"group_id":               CHECKSET,
						"group_metric_rule_name": name,
						"metric_name":            "disk_writebytes",
						"namespace":              "acs_ecs_dashboard",
						"contact_groups":         CHECKSET,
						"dimensions":             CHECKSET,
						"email_subject":          "tf-testacc-rule-name-warning",
						"effective_interval":     "00:00-22:59",
						"no_effective_interval":  "00:00-06:30",
						"period":                 "180",
						"silence_time":           "85800",
						"webhook":                "https://www.aliyun.com",
						"targets.#":              "1",
						"escalations.#":          "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"category", "interval"},
			},
		},
	})
}

var resourceAliCloudCmsGroupMetricRuleMap = map[string]string{
	"contact_groups": CHECKSET,
	"dimensions":     CHECKSET,
	"email_subject":  CHECKSET,
	"period":         "60",
	"silence_time":   "86400",
	"status":         CHECKSET,
}

func resourceAliCloudCmsGroupMetricRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cms_alarm_contact_group" "default" {
  		alarm_contact_group_name = var.name
  		describe                 = "tf-testacc"
  		contacts                 = ["test1", "test2", "test3"]
	}

	resource "alicloud_cms_monitor_group" "default" {
  		monitor_group_name = var.name
  		contact_groups     = [alicloud_cms_alarm_contact_group.default.id]
	}

	resource "alicloud_cms_alarm_contact_group" "update" {
  		alarm_contact_group_name = "${var.name}-update"
  		describe                 = "tf-testacc"
  		contacts                 = ["test5", "test6", "test7"]
	}

	resource "alicloud_cms_monitor_group" "update" {
  		monitor_group_name = "${var.name}-update"
  		contact_groups     = [alicloud_cms_alarm_contact_group.update.id]
	}
`, name)
}

func TestUnitAliCloudCmsGroupMetricRule(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cms_group_metric_rule"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cms_group_metric_rule"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"category":           "CreateGroupMetricRuleValue",
		"contact_groups":     "CreateGroupMetricRuleValue",
		"dimensions":         "CreateGroupMetricRuleValue",
		"effective_interval": "CreateGroupMetricRuleValue",
		"email_subject":      "CreateGroupMetricRuleValue",
		"escalations": []map[string]interface{}{
			{
				"critical": []map[string]interface{}{
					{
						"comparison_operator": "CreateGroupMetricRuleValue",
						"statistics":          "CreateGroupMetricRuleValue",
						"threshold":           "CreateGroupMetricRuleValue",
						"times":               1,
					},
				},
				"info": []map[string]interface{}{
					{
						"comparison_operator": "CreateGroupMetricRuleValue",
						"statistics":          "CreateGroupMetricRuleValue",
						"threshold":           "CreateGroupMetricRuleValue",
						"times":               1,
					},
				},
				"warn": []map[string]interface{}{
					{
						"comparison_operator": "CreateGroupMetricRuleValue",
						"statistics":          "CreateGroupMetricRuleValue",
						"threshold":           "CreateGroupMetricRuleValue",
						"times":               1,
					},
				},
			},
		},
		"group_id":               "CreateGroupMetricRuleValue",
		"group_metric_rule_name": "CreateGroupMetricRuleValue",
		"interval":               "CreateGroupMetricRuleValue",
		"metric_name":            "CreateGroupMetricRuleValue",
		"namespace":              "CreateGroupMetricRuleValue",
		"no_effective_interval":  "CreateGroupMetricRuleValue",
		"period":                 1,
		"rule_id":                "CreateGroupMetricRuleValue",
		"silence_time":           10,
		"webhook":                "CreateGroupMetricRuleValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeMetricRuleList
		"Alarms": map[string]interface{}{
			"Alarm": []interface{}{
				map[string]interface{}{
					"ContactGroups":     "CreateGroupMetricRuleValue",
					"Dimensions":        "CreateGroupMetricRuleValue",
					"EffectiveInterval": "CreateGroupMetricRuleValue",
					"MailSubject":       "CreateGroupMetricRuleValue",
					"Escalations": map[string]interface{}{
						"Critical": map[string]interface{}{
							"ComparisonOperator": "CreateGroupMetricRuleValue",
							"Statistics":         "CreateGroupMetricRuleValue",
							"Threshold":          "CreateGroupMetricRuleValue",
							"Times":              1,
						},
						"Info": map[string]interface{}{
							"ComparisonOperator": "CreateGroupMetricRuleValue",
							"Statistics":         "CreateGroupMetricRuleValue",
							"Threshold":          "CreateGroupMetricRuleValue",
							"Times":              1,
						},
						"Warn": map[string]interface{}{
							"ComparisonOperator": "CreateGroupMetricRuleValue",
							"Statistics":         "CreateGroupMetricRuleValue",
							"Threshold":          "CreateGroupMetricRuleValue",
							"Times":              1,
						},
					},
					"GroupId":             "CreateGroupMetricRuleValue",
					"RuleName":            "CreateGroupMetricRuleValue",
					"MetricName":          "CreateGroupMetricRuleValue",
					"Namespace":           "CreateGroupMetricRuleValue",
					"NoEffectiveInterval": "CreateGroupMetricRuleValue",
					"Period":              1,
					"SilenceTime":         10,
					"AlertState":          "OK",
					"Webhook":             "CreateGroupMetricRuleValue",
					"RuleId":              "CreateGroupMetricRuleValue",
				},
			},
		},
		"Targets": map[string]interface{}{
			"Target": []interface{}{},
		},
		"Code":    200,
		"Message": "Message",
	}
	CreateMockResponse := map[string]interface{}{
		// PutGroupMetricRule
		"Alarms": map[string]interface{}{
			"Alarm": []interface{}{
				map[string]interface{}{
					"RuleId": "CreateGroupMetricRuleValue",
				},
			},
		},
		"Code":    200,
		"Message": "Message",
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cms_group_metric_rule", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAliCloudCmsGroupMetricRuleCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		errorCodes := []string{"NonRetryableError", "Throttling", "ExceedingQuota", "Throttling.User", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "PutGroupMetricRule" {
					switch errorCode {
					case "NonRetryableError":
						return failedResponseMock(errorCode)
					default:
						retryIndex++
						if retryIndex >= len(errorCodes)-1 {
							successResponseMock(ReadMockResponseDiff)
							return CreateMockResponse, nil
						}
						return failedResponseMock(errorCodes[retryIndex])
					}
				}
				return ReadMockResponse, nil
			})
			err := resourceAliCloudCmsGroupMetricRuleCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cms_group_metric_rule"].Schema).Data(dInit.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAliCloudCmsGroupMetricRuleUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		// PutGroupMetricRule
		attributesDiff := map[string]interface{}{
			"rule_id":                "CreateGroupMetricRuleValue",
			"group_id":               "PutGroupMetricRuleValue",
			"group_metric_rule_name": "PutGroupMetricRuleValue",
			"metric_name":            "PutGroupMetricRuleValue",
			"namespace":              "PutGroupMetricRuleValue",
			"contact_groups":         "PutGroupMetricRuleValue",
			"dimensions":             "PutGroupMetricRuleValue",
			"effective_interval":     "PutGroupMetricRuleValue",
			"email_subject":          "PutGroupMetricRuleValue",
			"escalations": []map[string]interface{}{
				{
					"critical": []map[string]interface{}{
						{
							"comparison_operator": "PutGroupMetricRuleValue",
							"statistics":          "PutGroupMetricRuleValue",
							"threshold":           "PutGroupMetricRuleValue",
							"times":               2,
						},
					},
					"info": []map[string]interface{}{
						{
							"comparison_operator": "PutGroupMetricRuleValue",
							"statistics":          "PutGroupMetricRuleValue",
							"threshold":           "PutGroupMetricRuleValue",
							"times":               2,
						},
					},
					"warn": []map[string]interface{}{
						{
							"comparison_operator": "PutGroupMetricRuleValue",
							"statistics":          "PutGroupMetricRuleValue",
							"threshold":           "PutGroupMetricRuleValue",
							"times":               2,
						},
					},
				},
			},
			"no_effective_interval": "PutGroupMetricRuleValue",
			"period":                2,
			"silence_time":          20,
			"webhook":               "PutGroupMetricRuleValue",
			"category":              "PutGroupMetricRuleValue",
			"interval":              "PutGroupMetricRuleValue",
		}
		diff, err := newInstanceDiff("alicloud_cms_group_metric_rule", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cms_group_metric_rule"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeMetricRuleList Response
			"Alarms": map[string]interface{}{
				"Alarm": []interface{}{
					map[string]interface{}{
						"RuleId":            "CreateGroupMetricRuleValue",
						"GroupId":           "PutGroupMetricRuleValue",
						"RuleName":          "PutGroupMetricRuleValue",
						"MetricName":        "PutGroupMetricRuleValue",
						"Namespace":         "PutGroupMetricRuleValue",
						"ContactGroups":     "PutGroupMetricRuleValue",
						"Dimensions":        "PutGroupMetricRuleValue",
						"EffectiveInterval": "PutGroupMetricRuleValue",
						"MailSubject":       "PutGroupMetricRuleValue",
						"Escalations": map[string]interface{}{
							"Critical": map[string]interface{}{
								"ComparisonOperator": "PutGroupMetricRuleValue",
								"Statistics":         "PutGroupMetricRuleValue",
								"Threshold":          "PutGroupMetricRuleValue",
								"Times":              2,
							},
							"Info": map[string]interface{}{
								"ComparisonOperator": "PutGroupMetricRuleValue",
								"Statistics":         "PutGroupMetricRuleValue",
								"Threshold":          "PutGroupMetricRuleValue",
								"Times":              2,
							},
							"Warn": map[string]interface{}{
								"ComparisonOperator": "PutGroupMetricRuleValue",
								"Statistics":         "PutGroupMetricRuleValue",
								"Threshold":          "PutGroupMetricRuleValue",
								"Times":              2,
							},
						},
						"NoEffectiveInterval": "PutGroupMetricRuleValue",
						"Period":              2,
						"SilenceTime":         20,
						"AlertState":          "OK",
						"Webhook":             "PutGroupMetricRuleValue",
						"Category":            "PutGroupMetricRuleValue",
						"Interval":            "PutGroupMetricRuleValue",
					},
				},
			},
			"Targets": map[string]interface{}{
				"Target": []interface{}{},
			},
			"Code":    200,
			"Message": "Message",
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "Throttling.User", "ExceedingQuota", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "PutGroupMetricRule" {
					switch errorCode {
					case "NonRetryableError":
						return failedResponseMock(errorCode)
					default:
						retryIndex++
						if retryIndex >= len(errorCodes)-1 {
							return successResponseMock(ReadMockResponseDiff)
						}
						return failedResponseMock(errorCodes[retryIndex])
					}
				}
				return ReadMockResponse, nil
			})
			err := resourceAliCloudCmsGroupMetricRuleUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cms_group_metric_rule"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Read
	t.Run("Read", func(t *testing.T) {
		errorCodes := []string{"nil", "{}"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DescribeMetricRuleList" {
					switch errorCode {
					case "{}":
						return notFoundResponseMock(errorCode)
					default:
						retryIndex++
						if errorCodes[retryIndex] == "nil" {
							return ReadMockResponse, nil
						}
						return failedResponseMock(errorCodes[retryIndex])
					}
				}
				return ReadMockResponse, nil
			})
			err := resourceAliCloudCmsGroupMetricRuleRead(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "{}":
				assert.Nil(t, err)
			}
		}
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAliCloudCmsGroupMetricRuleDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		attributesDiff := map[string]interface{}{}
		diff, err := newInstanceDiff("alicloud_cms_group_metric_rule", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cms_group_metric_rule"].Schema).Data(dInit.State(), diff)
		errorCodes := []string{"NonRetryableError", "Throttling", "ExceedingQuota", "Throttling.User", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteMetricRules" {
					switch errorCode {
					case "NonRetryableError":
						return failedResponseMock(errorCode)
					default:
						retryIndex++
						if errorCodes[retryIndex] == "nil" {
							ReadMockResponse = map[string]interface{}{
								"Code":    200,
								"Message": "Message",
							}
							return ReadMockResponse, nil
						}
						return failedResponseMock(errorCodes[retryIndex])
					}
				}
				return ReadMockResponse, nil
			})
			err := resourceAliCloudCmsGroupMetricRuleDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "nil":
				assert.Nil(t, err)
			}
		}
	})
}
