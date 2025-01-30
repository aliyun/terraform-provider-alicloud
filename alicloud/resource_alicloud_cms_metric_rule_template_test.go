package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

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
	resource.AddTestSweepers(
		"alicloud_cms_metric_rule_template",
		&resource.Sweeper{
			Name: "alicloud_cms_metric_rule_template",
			F:    testSweepCmsMetricRuleTemplate,
		})
}

func testSweepCmsMetricRuleTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeMetricRuleTemplateList"
	request := map[string]interface{}{}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Templates.Template", response)

		if formatInt(response["Total"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Templates.Template", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Cms Metric Rule Template: %s", item["Name"].(string))
				continue
			}

			action := "DeleteMetricRuleTemplate"
			request := map[string]interface{}{
				"TemplateId": item["TemplateId"],
			}
			_, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Cms Metric Rule Template (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Cms Metric Rule Template success: %s ", item["Name"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudCmsMetricRuleTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_metric_rule_template.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsMetricRuleTemplateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMetricRuleTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicemetricruletemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsMetricRuleTemplateBasicDependence0)
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
					"metric_rule_template_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_rule_template_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_templates": []map[string]interface{}{
						{
							"rule_name":   name,
							"metric_name": "cpu_total",
							"namespace":   "acs_ecs_dashboard",
							"category":    "ecs",
							"escalations": []map[string]interface{}{
								{
									"critical": []map[string]interface{}{
										{
											"comparison_operator": "GreaterThanThreshold",
											"statistics":          "Average",
											"threshold":           "90",
											"times":               "1",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_templates.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_templates": []map[string]interface{}{
						{
							"rule_name":   name,
							"metric_name": "cpu_total",
							"namespace":   "acs_ecs_dashboard",
							"category":    "ecs",
							"webhook":     "https://www.aliyun.com",
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
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_templates.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_templates": []map[string]interface{}{
						{
							"rule_name":   name,
							"metric_name": "cpu_total",
							"namespace":   "acs_ecs_dashboard",
							"category":    "ecs",
							"webhook":     "https://www.aliyun.com",
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
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_templates.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_templates": []map[string]interface{}{
						{
							"rule_name":   name,
							"metric_name": "cpu_total",
							"namespace":   "acs_ecs_dashboard",
							"category":    "ecs",
							"webhook":     "https://www.aliyun.com",
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
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_templates.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_id", "apply_mode", "notify_level", "silence_time", "webhook", "enable_start_time", "enable_end_time"},
			},
		},
	})
}

func TestAccAliCloudCmsMetricRuleTemplate_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_metric_rule_template.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsMetricRuleTemplateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMetricRuleTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicemetricruletemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsMetricRuleTemplateBasicDependence0)
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
					"metric_rule_template_name": name,
					"description":               name,
					"group_id":                  "${alicloud_cms_monitor_group.default.id}",
					"apply_mode":                "GROUP_INSTANCE_FIRST",
					"notify_level":              "2",
					"silence_time":              "80000",
					"webhook":                   "https://www.aliyun.com",
					"enable_start_time":         "00",
					"enable_end_time":           "23",
					"alert_templates": []map[string]interface{}{
						{
							"rule_name":   name,
							"metric_name": "cpu_total",
							"namespace":   "acs_ecs_dashboard",
							"category":    "ecs",
							"webhook":     "https://www.aliyun.com",
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
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_rule_template_name": name,
						"description":               name,
						"alert_templates.#":         "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_id", "apply_mode", "notify_level", "silence_time", "webhook", "enable_start_time", "enable_end_time"},
			},
		},
	})
}

var AliCloudCmsMetricRuleTemplateMap0 = map[string]string{
	"rest_version": CHECKSET,
}

func AliCloudCmsMetricRuleTemplateBasicDependence0(name string) string {
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
`, name)
}

func TestUnitAliCloudCmsMetricRuleTemplate(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cms_metric_rule_template"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cms_metric_rule_template"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"alert_templates": []map[string]interface{}{
			{
				"category":    "CreateMetricRuleTemplateValue",
				"metric_name": "CreateMetricRuleTemplateValue",
				"namespace":   "CreateMetricRuleTemplateValue",
				"rule_name":   "CreateMetricRuleTemplateValue",
				"webhook":     "CreateMetricRuleTemplateValue",
				"escalations": []map[string]interface{}{
					{
						"critical": []map[string]interface{}{
							{
								"comparison_operator": "CreateMetricRuleTemplateValue",
								"statistics":          "CreateMetricRuleTemplateValue",
								"threshold":           "CreateMetricRuleTemplateValue",
								"times":               "CreateMetricRuleTemplateValue",
							},
						},
						"info": []map[string]interface{}{
							{
								"comparison_operator": "CreateMetricRuleTemplateValue",
								"statistics":          "CreateMetricRuleTemplateValue",
								"threshold":           "CreateMetricRuleTemplateValue",
								"times":               "CreateMetricRuleTemplateValue",
							},
						},
						"warn": []map[string]interface{}{
							{
								"comparison_operator": "CreateMetricRuleTemplateValue",
								"statistics":          "CreateMetricRuleTemplateValue",
								"threshold":           "CreateMetricRuleTemplateValue",
								"times":               "CreateMetricRuleTemplateValue",
							},
						},
					},
				},
			},
		},
		"description":               "CreateMetricRuleTemplateValue",
		"metric_rule_template_name": "CreateMetricRuleTemplateValue",
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
		// DescribeMetricRuleTemplateAttribute
		"Resource": map[string]interface{}{
			"AlertTemplates": map[string]interface{}{
				"AlertTemplate": []interface{}{
					map[string]interface{}{
						"Category":   "CreateMetricRuleTemplateValue",
						"Webhook":    "CreateMetricRuleTemplateValue",
						"Namespace":  "CreateMetricRuleTemplateValue",
						"RuleName":   "CreateMetricRuleTemplateValue",
						"MetricName": "CreateMetricRuleTemplateValue",
						"Escalations": map[string]interface{}{
							"Critical": map[string]interface{}{
								"ComparisonOperator": "CreateMetricRuleTemplateValue",
								"Statistics":         "CreateMetricRuleTemplateValue",
								"Threshold":          "CreateMetricRuleTemplateValue",
								"Times":              "CreateMetricRuleTemplateValue",
							},
							"Info": map[string]interface{}{
								"ComparisonOperator": "CreateMetricRuleTemplateValue",
								"Statistics":         "CreateMetricRuleTemplateValue",
								"Threshold":          "CreateMetricRuleTemplateValue",
								"Times":              "CreateMetricRuleTemplateValue",
							},
							"Warn": map[string]interface{}{
								"ComparisonOperator": "CreateMetricRuleTemplateValue",
								"Statistics":         "CreateMetricRuleTemplateValue",
								"Threshold":          "CreateMetricRuleTemplateValue",
								"Times":              "CreateMetricRuleTemplateValue",
							},
						},
					},
				},
			},
			"Description": "CreateMetricRuleTemplateValue",
			"Name":        "CreateMetricRuleTemplateValue",
			"RestVersion": "CreateMetricRuleTemplateValue",
		},
		"Success": "true",
		"Id":      "CreateMetricRuleTemplateValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateMetricRuleTemplate
		"Success": "true",
		"Id":      "CreateMetricRuleTemplateValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cms_metric_rule_template", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCmsMetricRuleTemplateCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeMetricRuleTemplateAttribute Response
		"Alarms": map[string]interface{}{
			"Alarm": []interface{}{
				map[string]interface{}{
					"RuleId": "CreateMetricRuleTemplateValue",
				},
			},
		},
		"Code":    200,
		"Message": "Message",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateMetricRuleTemplate" {
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
		err := resourceAliCloudCmsMetricRuleTemplateCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_metric_rule_template"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCmsMetricRuleTemplateUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ApplyMetricRuleTemplate
	attributesDiff := map[string]interface{}{
		"group_id":          "ApplyMetricRuleTemplateValue",
		"apply_mode":        "ApplyMetricRuleTemplateValue",
		"enable_end_time":   "ApplyMetricRuleTemplateValue",
		"enable_start_time": "ApplyMetricRuleTemplateValue",
		"notify_level":      "ApplyMetricRuleTemplateValue",
		"silence_time":      60,
		"webhook":           "ApplyMetricRuleTemplateValue",
	}
	diff, err := newInstanceDiff("alicloud_cms_metric_rule_template", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_metric_rule_template"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeMetricRuleTemplateAttribute Response
		"Resource": map[string]interface{}{
			"GroupId":         "ApplyMetricRuleTemplateValue",
			"ApplyMode":       "ApplyMetricRuleTemplateValue",
			"EnableEndTime":   "ApplyMetricRuleTemplateValue",
			"EnableStartTime": "ApplyMetricRuleTemplateValue",
			"NotifyLevel":     "ApplyMetricRuleTemplateValue",
			"SilenceTime":     60,
			"Webhook":         "ApplyMetricRuleTemplateValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ApplyMetricRuleTemplate" {
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
		err := resourceAliCloudCmsMetricRuleTemplateUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_metric_rule_template"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyMetricRuleTemplate
	attributesDiff = map[string]interface{}{
		"rest_version":              "ModifyMetricRuleTemplateValue",
		"description":               "ModifyMetricRuleTemplateValue",
		"metric_rule_template_name": "ModifyMetricRuleTemplateValue",
		"alert_templates": []map[string]interface{}{
			{
				"category":    "ModifyMetricRuleTemplateValue",
				"metric_name": "ModifyMetricRuleTemplateValue",
				"namespace":   "ModifyMetricRuleTemplateValue",
				"rule_name":   "ModifyMetricRuleTemplateValue",
				"webhook":     "ModifyMetricRuleTemplateValue",
				"escalations": []map[string]interface{}{
					{
						"critical": []map[string]interface{}{
							{
								"comparison_operator": "ModifyMetricRuleTemplateValue",
								"statistics":          "ModifyMetricRuleTemplateValue",
								"threshold":           "ModifyMetricRuleTemplateValue",
								"times":               "ModifyMetricRuleTemplateValue",
							},
						},
						"info": []map[string]interface{}{
							{
								"comparison_operator": "ModifyMetricRuleTemplateValue",
								"statistics":          "ModifyMetricRuleTemplateValue",
								"threshold":           "ModifyMetricRuleTemplateValue",
								"times":               "ModifyMetricRuleTemplateValue",
							},
						},
						"warn": []map[string]interface{}{
							{
								"comparison_operator": "ModifyMetricRuleTemplateValue",
								"statistics":          "ModifyMetricRuleTemplateValue",
								"threshold":           "ModifyMetricRuleTemplateValue",
								"times":               "ModifyMetricRuleTemplateValue",
							},
						},
					},
				},
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_cms_metric_rule_template", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_metric_rule_template"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeMetricRuleTemplateAttribute Response
		"Resource": map[string]interface{}{
			"AlertTemplates": map[string]interface{}{
				"AlertTemplate": []interface{}{
					map[string]interface{}{
						"Category":   "ModifyMetricRuleTemplateValue",
						"Webhook":    "ModifyMetricRuleTemplateValue",
						"Namespace":  "ModifyMetricRuleTemplateValue",
						"RuleName":   "ModifyMetricRuleTemplateValue",
						"MetricName": "ModifyMetricRuleTemplateValue",
						"Escalations": map[string]interface{}{
							"Critical": map[string]interface{}{
								"ComparisonOperator": "ModifyMetricRuleTemplateValue",
								"Statistics":         "ModifyMetricRuleTemplateValue",
								"Threshold":          "ModifyMetricRuleTemplateValue",
								"Times":              "ModifyMetricRuleTemplateValue",
							},
							"Info": map[string]interface{}{
								"ComparisonOperator": "ModifyMetricRuleTemplateValue",
								"Statistics":         "ModifyMetricRuleTemplateValue",
								"Threshold":          "ModifyMetricRuleTemplateValue",
								"Times":              "ModifyMetricRuleTemplateValue",
							},
							"Warn": map[string]interface{}{
								"ComparisonOperator": "ModifyMetricRuleTemplateValue",
								"Statistics":         "ModifyMetricRuleTemplateValue",
								"Threshold":          "ModifyMetricRuleTemplateValue",
								"Times":              "ModifyMetricRuleTemplateValue",
							},
						},
					},
				},
			},
			"Description": "ModifyMetricRuleTemplateValue",
			"Name":        "ModifyMetricRuleTemplateValue",
			"RestVersion": "ModifyMetricRuleTemplateValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyMetricRuleTemplate" {
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
		err := resourceAliCloudCmsMetricRuleTemplateUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_metric_rule_template"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeMetricRuleTemplateAttribute" {
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
		err := resourceAliCloudCmsMetricRuleTemplateRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCmsMetricRuleTemplateDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteMetricRuleTemplate" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Code":    "",
							"Success": "true",
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudCmsMetricRuleTemplateDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
