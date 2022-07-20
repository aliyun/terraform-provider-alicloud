package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudARMSPrometheusAlertRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_alert_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudARMSPrometheusAlertRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsPrometheusAlertRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusalertrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudARMSPrometheusAlertRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPrePaidResources(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"cluster_id":                 "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
					"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
					"message":                    "node available memory is less than 10%",
					"duration":                   "1",
					"notify_type":                "ALERT_MANAGER",
					"type":                       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"cluster_id":                 CHECKSET,
						"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
						"message":                    "node available memory is less than 10%",
						"duration":                   "1",
						"notify_type":                "ALERT_MANAGER",
						"type":                       name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"duration": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duration": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"name":  "TF",
							"value": "test1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"annotations": []map[string]interface{}{
						{
							"name":  "TF",
							"value": "test1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"annotations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"message": "node available memory is less than 20%",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"message": "node available memory is less than 20%",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expression": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expression": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"type":                       name,
					"duration":                   "1",
					"labels": []map[string]interface{}{
						{
							"name":  "TF2",
							"value": "test2",
						},
					},
					"annotations": []map[string]interface{}{
						{
							"name":  "TF2",
							"value": "test2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"duration":                   "1",
						"type":                       name,
						"labels.#":                   "1",
						"annotations.#":              "1",
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

func TestAccAlicloudARMSPrometheusAlertRule_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_alert_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudARMSPrometheusAlertRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsPrometheusAlertRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusalertrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudARMSPrometheusAlertRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPrePaidResources(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"cluster_id":                 "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
					"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
					"message":                    "node available memory is less than 10%",
					"duration":                   "1",
					"notify_type":                "DISPATCH_RULE",
					"dispatch_rule_id":           "${alicloud_arms_dispatch_rule.default.id}",
					"type":                       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"cluster_id":                 CHECKSET,
						"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
						"message":                    "node available memory is less than 10%",
						"duration":                   "1",
						"notify_type":                "DISPATCH_RULE",
						"dispatch_rule_id":           CHECKSET,
						"type":                       name,
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

var AlicloudARMSPrometheusAlertRuleMap0 = map[string]string{}

func AlicloudARMSPrometheusAlertRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "Default"
}
`, name)
}

func AlicloudARMSPrometheusAlertRuleBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}
data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "Default"
}
resource "alicloud_arms_alert_contact" "default" {
  alert_contact_name = var.name
  email              = "${var.name}@aaa.com"
}
resource "alicloud_arms_alert_contact_group" "default" {
  alert_contact_group_name = var.name
  contact_ids              = [alicloud_arms_alert_contact.default.id]
}

resource "alicloud_arms_dispatch_rule" "default" {
  dispatch_rule_name = var.name
  dispatch_type      = "CREATE_ALERT"
  group_rules {
    group_wait_time = 5
    group_interval  = 15
    repeat_interval = 100
    grouping_fields = [
      "alertname"]
  }
  label_match_expression_grid {
   label_match_expression_groups {
     label_match_expressions {
       key      = "_aliyun_arms_involvedObject_kind"
       value    = "app"
       operator = "eq"
     }
   }
  }

  notify_rules {
    notify_objects {
      notify_object_id = alicloud_arms_alert_contact.default.id
      notify_type      = "ARMS_CONTACT"
      name             = var.name
    }
    notify_objects {
      notify_object_id = alicloud_arms_alert_contact_group.default.id
      notify_type      = "ARMS_CONTACT_GROUP"
      name             = var.name
    }
    notify_channels = ["dingTalk", "wechat"]
  }
}
`, name)
}

func TestUnitAlicloudARMSPrometheusAlertRule(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	rand := acctest.RandIntRange(10000, 99999)
	d, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"prometheus_alert_rule_name": "prometheus_alert_rule_name",
		"cluster_id":                 "ClusterId",
		"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
		"message":                    "node available memory is less than 10%",
		"duration":                   "1",
		"notify_type":                "DISPATCH_RULE",
		"type":                       "type",
		"dispatch_rule_id":           "1",
		"annotations": []map[string]interface{}{
			{
				"name":  "TF",
				"value": "test1",
			},
		},
		"labels": []map[string]interface{}{
			{
				"name":  "TF2",
				"value": "test2",
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
		"PrometheusAlertRules": []interface{}{
			map[string]interface{}{
				"AlertId":               "MockAlertId",
				"ClusterId":             "ClusterId",
				"PrometheusAlertRuleId": "MockAlertId",
				"Duration":              "duration",
				"Expression":            "expression",
				"Annotations": []interface{}{
					map[string]interface{}{
						"Name":  "TF",
						"Value": "test1",
					},
				},
				"Labels": []interface{}{
					map[string]interface{}{
						"Name":  "TF",
						"Value": "test1",
					},
				},
				"DispatchRuleId": "1",
				"Message":        "message",
				"NotifyType":     "notify_type",
				"AlertName":      "prometheus_alert_rule_name",
				"Status":         "1",
				"Type":           "type",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_arms_prometheus_alert_rule", "MockAlertId"))
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
			result["PrometheusAlertRule"] = map[string]interface{}{
				"AlertId": "MockAlertId",
			}
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
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
		"AnnotationsName": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"PrometheusAlertRules": []interface{}{
					map[string]interface{}{
						"AlertId":               "MockAlertId",
						"ClusterId":             "ClusterId",
						"PrometheusAlertRuleId": "MockAlertId",
						"Annotations": []interface{}{
							map[string]interface{}{
								"Name":  "message",
								"Value": "test1",
							},
						},
					},
				},
			}
			result["PrometheusAlertRule"] = map[string]interface{}{
				"AlertId": "MockAlertId",
			}
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudArmsPrometheusAlertRuleCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateAnnotationsNameMock", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["AnnotationsName"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("ClusterId", ":", "MockAlertId"))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudArmsPrometheusAlertRuleUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdatePrometheusAlertRuleAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"duration", "expression", "message", "prometheus_alert_rule_name", "annotations", "dispatch_rule_id", "labels", "notify_type"} {
			switch p["alicloud_arms_prometheus_alert_rule"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(2)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			case schema.TypeSet:
				diff.SetAttribute(fmt.Sprintf("%s.#", key), &terraform.ResourceAttrDiff{Old: "1", New: "1"})
				for _, ipConfig := range d.Get(key).(*schema.Set).List() {
					ipConfigArg := ipConfig.(map[string]interface{})
					for field, _ := range p["alicloud_arms_prometheus_alert_rule"].Schema[key].Elem.(*schema.Resource).Schema {
						diff.SetAttribute(fmt.Sprintf("%s.%d.%s", key, rand, field), &terraform.ResourceAttrDiff{Old: ipConfigArg[field].(string), New: ipConfigArg[field].(string) + "_update"})
					}
				}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdatePrometheusAlertRuleNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"duration", "expression", "message", "prometheus_alert_rule_name", "annotations", "dispatch_rule_id", "labels", "notify_type"} {
			switch p["alicloud_arms_prometheus_alert_rule"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			case schema.TypeSet:
				diff.SetAttribute(fmt.Sprintf("%s.#", key), &terraform.ResourceAttrDiff{Old: "1", New: "1"})
				for _, ipConfig := range d.Get(key).(*schema.Set).List() {
					ipConfigArg := ipConfig.(map[string]interface{})
					for field, _ := range p["alicloud_arms_prometheus_alert_rule"].Schema[key].Elem.(*schema.Resource).Schema {
						diff.SetAttribute(fmt.Sprintf("%s.%d.%s", key, rand, field), &terraform.ResourceAttrDiff{Old: ipConfigArg[field].(string), New: ipConfigArg[field].(string) + "_update"})
					}
				}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudArmsPrometheusAlertRuleDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleDelete(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeArmsPrometheusAlertRuleNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeArmsPrometheusAlertRuleAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})

	t.Run("ReadMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := false
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		patcheDescribeArmsPrometheusAlertRule := gomonkey.ApplyMethod(reflect.TypeOf(&ArmsService{}), "DescribeArmsPrometheusAlertRule", func(*ArmsService, string) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudArmsPrometheusAlertRuleRead(resourceData1, rawClient)
		patcheDorequest.Reset()
		patcheDescribeArmsPrometheusAlertRule.Reset()
		assert.NotNil(t, err)
	})
}
