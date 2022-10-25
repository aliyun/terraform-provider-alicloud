package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAlicloudSAEApplicationScalingRule_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application_scaling_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEApplicationScalingRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplicationScalingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEApplicationScalingRuleBasicDependence0)
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
					"app_id":                   "${alicloud_sae_application.default.id}",
					"scaling_rule_name":        "${var.name}",
					"scaling_rule_type":        "timing",
					"min_ready_instances":      "3",
					"min_ready_instance_ratio": "-1",
					"scaling_rule_enable":      "true",
					"scaling_rule_timer": []map[string]interface{}{
						{
							"begin_date": "2022-02-25",
							"end_date":   "2022-03-25",
							"period":     "* * *",
							"schedules": []map[string]interface{}{
								{
									"at_time":         "08:00",
									"target_replicas": "10",
								},
								{
									"at_time":         "20:00",
									"target_replicas": "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":               CHECKSET,
						"scaling_rule_name":    name,
						"scaling_rule_type":    "timing",
						"scaling_rule_enable":  "true",
						"scaling_rule_timer.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_rule_enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_rule_enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_enable": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_ready_instances": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_ready_instance_ratio": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_rule_timer": []map[string]interface{}{
						{
							"begin_date": "2022-04-25",
							"end_date":   "2022-05-25",
							"period":     "* * *",
							"schedules": []map[string]interface{}{
								{
									"at_time":         "07:00",
									"target_replicas": "10",
								},
								{
									"at_time":         "22:00",
									"target_replicas": "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_timer.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_ready_instances":      "3",
					"min_ready_instance_ratio": "-1",
					"scaling_rule_timer": []map[string]interface{}{
						{
							"begin_date": "2022-02-25",
							"end_date":   "2022-03-25",
							"period":     "* * *",
							"schedules": []map[string]interface{}{
								{
									"at_time":         "08:00",
									"target_replicas": "10",
								},
								{
									"at_time":         "20:00",
									"target_replicas": "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_timer.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_ready_instances", "min_ready_instance_ratio"},
			},
		},
	})
}

func TestAccAlicloudSAEApplicationScalingRule_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application_scaling_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEApplicationScalingRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplicationScalingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEApplicationScalingRuleBasicDependence0)
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
					"app_id":                   "${alicloud_sae_application.default.id}",
					"scaling_rule_name":        "${var.name}",
					"scaling_rule_type":        "metric",
					"min_ready_instances":      "3",
					"min_ready_instance_ratio": "-1",
					"scaling_rule_enable":      "true",
					"scaling_rule_metric": []map[string]interface{}{
						{
							"max_replicas": "50",
							"min_replicas": "3",
							"metrics": []map[string]interface{}{
								{
									"metric_type":                       "CPU",
									"metric_target_average_utilization": "20",
								},
								{
									"metric_type":                       "MEMORY",
									"metric_target_average_utilization": "30",
								},
								{
									"metric_type":                       "tcpActiveConn",
									"metric_target_average_utilization": "20",
								},
							},
							"scale_up_rules": []map[string]interface{}{
								{
									"step":                         "100",
									"disabled":                     "false",
									"stabilization_window_seconds": "0",
								},
							},
							"scale_down_rules": []map[string]interface{}{
								{
									"step":                         "100",
									"disabled":                     "false",
									"stabilization_window_seconds": "300",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":                CHECKSET,
						"scaling_rule_name":     name,
						"scaling_rule_type":     "metric",
						"scaling_rule_enable":   "true",
						"scaling_rule_metric.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_rule_metric": []map[string]interface{}{
						{
							"max_replicas": "40",
							"min_replicas": "5",
							"metrics": []map[string]interface{}{
								{
									"metric_type":                       "CPU",
									"metric_target_average_utilization": "30",
								},
								{
									"metric_type":                       "MEMORY",
									"metric_target_average_utilization": "40",
								},
								{
									"metric_type":                       "tcpActiveConn",
									"metric_target_average_utilization": "30",
								},
							},
							"scale_up_rules": []map[string]interface{}{
								{
									"step":                         "90",
									"disabled":                     "false",
									"stabilization_window_seconds": "10",
								},
							},
							"scale_down_rules": []map[string]interface{}{
								{
									"step":                         "90",
									"disabled":                     "false",
									"stabilization_window_seconds": "200",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_metric.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_ready_instances", "min_ready_instance_ratio"},
			},
		},
	})
}

func TestAccAlicloudSAEApplicationScalingRule_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application_scaling_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEApplicationScalingRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplicationScalingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEApplicationScalingRuleBasicDependence0)
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
					"app_id":                   "${alicloud_sae_application.default.id}",
					"scaling_rule_name":        "${var.name}",
					"scaling_rule_type":        "mix",
					"min_ready_instances":      "3",
					"min_ready_instance_ratio": "-1",
					"scaling_rule_enable":      "true",
					"scaling_rule_timer": []map[string]interface{}{
						{
							"period": "* * *",
							"schedules": []map[string]interface{}{
								{
									"at_time":      "08:00",
									"max_replicas": "20",
									"min_replicas": "3",
								},
								{
									"at_time":      "20:00",
									"max_replicas": "10",
									"min_replicas": "3",
								},
							},
						},
					},
					"scaling_rule_metric": []map[string]interface{}{
						{
							"max_replicas": "50",
							"min_replicas": "3",
							"metrics": []map[string]interface{}{
								{
									"metric_type":                       "CPU",
									"metric_target_average_utilization": "20",
								},
								{
									"metric_type":                       "MEMORY",
									"metric_target_average_utilization": "30",
								},
								{
									"metric_type":                       "tcpActiveConn",
									"metric_target_average_utilization": "20",
								},
							},
							"scale_up_rules": []map[string]interface{}{
								{
									"step":                         "100",
									"disabled":                     "false",
									"stabilization_window_seconds": "0",
								},
							},
							"scale_down_rules": []map[string]interface{}{
								{
									"step":                         "100",
									"disabled":                     "false",
									"stabilization_window_seconds": "300",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":                CHECKSET,
						"scaling_rule_name":     name,
						"scaling_rule_type":     "mix",
						"scaling_rule_enable":   "true",
						"scaling_rule_timer.#":  "1",
						"scaling_rule_metric.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_ready_instances", "min_ready_instance_ratio"},
			},
		},
	})
}

func AlicloudSAEApplicationScalingRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = join(":",["%s",var.name])
  namespace_name        = var.name
}
resource "alicloud_sae_application" "default" {
  app_description = var.name
  app_name        = var.name
  namespace_id    = alicloud_sae_namespace.default.namespace_id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  jdk             = "Open JDK 8"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone        = "Asia/Shanghai"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}
`, name, defaultRegionToTest)
}

var AlicloudSAEApplicationScalingRuleMap0 = map[string]string{
	"app_id": CHECKSET,
}

func TestUnitAlicloudSAEApplicationScalingRule(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_sae_application_scaling_rule"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_sae_application_scaling_rule"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"app_id":                   "app_id",
		"scaling_rule_name":        "scaling_rule_name",
		"scaling_rule_type":        "mix",
		"min_ready_instances":      3,
		"min_ready_instance_ratio": -1,
		"scaling_rule_enable":      true,
		"scaling_rule_timer": []map[string]interface{}{
			{
				"begin_date": "2021-03-25",
				"end_date":   "2021-04-25",
				"period":     "* * *",
				"schedules": []map[string]interface{}{
					{
						"at_time":         "08:00",
						"target_replicas": 3,
						"max_replicas":    50,
						"min_replicas":    1,
					},
				},
			},
		},
		"scaling_rule_metric": []map[string]interface{}{
			{
				"max_replicas": 3,
				"min_replicas": 1,
				"metrics": []map[string]interface{}{
					{
						"metric_type":                       "CPU",
						"metric_target_average_utilization": 20,
					},
				},
				"scale_up_rules": []map[string]interface{}{
					{
						"step":                         100,
						"disabled":                     false,
						"stabilization_window_seconds": 0,
					},
				},
				"scale_down_rules": []map[string]interface{}{
					{
						"step":                         100,
						"disabled":                     false,
						"stabilization_window_seconds": 300,
					},
				},
			},
		},
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"body": map[string]interface{}{
			"Data": map[string]interface{}{
				"CurrentPage": 1,
				"TotalSize":   3,
				"PageSize":    10,
				"ApplicationScalingRules": []interface{}{
					map[string]interface{}{
						"Timer": map[string]interface{}{
							"EndDate":   "2021-04-25",
							"BeginDate": "2021-03-25",
							"Schedules": []interface{}{
								map[string]interface{}{
									"AtTime":         "08:00",
									"TargetReplicas": 3,
									"MaxReplicas":    50,
									"MinReplicas":    1,
								}},
							"Period": "* * *",
						},
						"AppId":            "app_id",
						"ScaleRuleEnabled": true,
						"ScaleRuleType":    "mix",
						"Metric": map[string]interface{}{
							"Metrics": []interface{}{
								map[string]interface{}{
									"MetricTargetAverageUtilization": 20,
									"MetricType":                     "CPU",
								}},
							"MetricsStatus": map[string]interface{}{
								"DesiredReplicas":     2,
								"NextScaleTimePeriod": 3,
								"CurrentReplicas":     2,
								"LastScaleTime":       "2022-01-11T08:14:32Z",
								"CurrentMetrics": []interface{}{
									map[string]interface{}{
										"Type":         "Resource",
										"CurrentValue": 0,
										"Name":         "cpu",
									},
								},
								"NextScaleMetrics": []interface{}{
									map[string]interface{}{
										"NextScaleOutAverageUtilization": 21,
										"NextScaleInAverageUtilization":  10,
										"Name":                           "cpu",
									},
								},
								"MaxReplicas": 3,
								"MinReplicas": 1,
							},
							"MaxReplicas": 3,
							"MinReplicas": 1,
							"ScaleUpRules": map[string]interface{}{
								"Step":                       100,
								"StabilizationWindowSeconds": 0,
								"Disabled":                   false,
							},
							"ScaleDownRules": map[string]interface{}{
								"Step":                       100,
								"StabilizationWindowSeconds": 300,
								"Disabled":                   false,
							},
						},
						"ScaleRuleName": "scale_rule_name",
					},
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_sae_application_scaling_rule", "xxx_id"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["body"].(map[string]interface{})["Data"].(map[string]interface{})["ScaleRuleName"] = "scale_rule_name"
			return result, nil
		},
		"CreateNoBodyError": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{}
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateStatusNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}

	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServerlessClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudSaeApplicationScalingRuleCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("CreateRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNoBodyError", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["CreateNoBodyError"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServerlessClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudSaeApplicationScalingRuleUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateParseResourceIdAbnormal", func(t *testing.T) {
		d.SetId("app_id")
		err := resourceAlicloudSaeApplicationScalingRuleUpdate(d, rawClient)
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("app_id:scaling_rule_name")

	t.Run("UpdateModifyAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"min_ready_instances"} {
			switch p["alicloud_sae_application_scaling_rule"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_sae_application_scaling_rule"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyAttributeNoRetryErrorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"min_ready_instances"} {
			switch p["alicloud_sae_application_scaling_rule"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData, _ := schema.InternalMap(p["alicloud_sae_application_scaling_rule"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&SaeService{}), "DescribeSaeApplicationScalingRule", func(*SaeService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudSaeApplicationScalingRuleUpdate(resourceData, rawClient)
		patches.Reset()
		patcheDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"min_ready_instances", "min_ready_instance_ratio"} {
			switch p["alicloud_sae_application_scaling_rule"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		diff.SetAttribute("scaling_rule_metric.0.max_replicas", &terraform.ResourceAttrDiff{Old: "3", New: "4"})
		diff.SetAttribute("scaling_rule_metric.0.scale_up_rules.0.step", &terraform.ResourceAttrDiff{Old: "100", New: "110"})
		diff.SetAttribute("scaling_rule_metric.0.scale_down_rules.0.step", &terraform.ResourceAttrDiff{Old: "100", New: "110"})
		diff.SetAttribute("scaling_rule_metric.0.metrics.0.metric_target_average_utilization", &terraform.ResourceAttrDiff{Old: "20", New: "30"})
		diff.SetAttribute("scaling_rule_timer.0.begin_date", &terraform.ResourceAttrDiff{Old: "2021-03-25", New: "2021-03-26"})
		diff.SetAttribute("scaling_rule_timer.0.schedules.0.max_replicas", &terraform.ResourceAttrDiff{Old: "50", New: "60"})
		resourceData1, _ := schema.InternalMap(p["alicloud_sae_application_scaling_rule"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateModifyStatusRunningAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("scaling_rule_enable", &terraform.ResourceAttrDiff{Old: "false", New: "true"})
		resourceData1, _ := schema.InternalMap(p["alicloud_sae_application_scaling_rule"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&SaeService{}), "DescribeSaeApplicationScalingRule", func(*SaeService, string) (map[string]interface{}, error) {
			object := map[string]interface{}{
				"AppId":            "app_id",
				"ScaleRuleEnabled": false,
				"ScaleRuleType":    "mix",
				"ScaleRuleName":    "scale_rule_name",
			}
			return object, nil
		})
		err := resourceAlicloudSaeApplicationScalingRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyStatusRunningAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("scaling_rule_enable", &terraform.ResourceAttrDiff{Old: strconv.FormatBool(false), New: strconv.FormatBool(true)})
		resourceData1, _ := schema.InternalMap(p["alicloud_sae_application_scaling_rule"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateStatusNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&SaeService{}), "DescribeSaeApplicationScalingRule", func(*SaeService, string) (map[string]interface{}, error) {
			object := map[string]interface{}{
				"AppId":            "app_id",
				"ScaleRuleEnabled": false,
				"ScaleRuleType":    "mix",
				"ScaleRuleName":    "scale_rule_name",
			}
			return object, nil
		})
		err := resourceAlicloudSaeApplicationScalingRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServerlessClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudSaeApplicationScalingRuleDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteParseResourceIdAbnormal", func(t *testing.T) {
		d.SetId("app_id")
		err := resourceAlicloudSaeApplicationScalingRuleDelete(d, rawClient)
		assert.NotNil(t, err)
	})

	d.SetId("app_id:scaling_rule_name")
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				// retry until the timeout comes
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeNotFound", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleRead(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudSaeApplicationScalingRuleRead(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
}
