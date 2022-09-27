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

func TestAccAlicloudCloudFirewallControlPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallControlPolicyBasicDependence0)
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
					"application_name": "ANY",
					"acl_action":       "accept",
					"description":      "放行所有流量",
					"destination_type": "net",
					"destination":      "100.1.1.0/24",
					"direction":        "out",
					"proto":            "ANY",
					"source":           "1.2.3.0/24",
					"source_type":      "net",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": "ANY",
						"acl_action":       "accept",
						"description":      "放行所有流量",
						"destination_type": "net",
						"destination":      "100.1.1.0/24",
						"direction":        "out",
						"proto":            "ANY",
						"source":           "1.2.3.0/24",
						"source_type":      "net",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"release": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"release": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action": "log",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action": "log",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "Any",
					"destination_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "Any",
						"destination_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "www.aliyun.com",
					"destination_type": "domain",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "www.aliyun.com",
						"destination_type": "domain",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":      "Any",
					"source_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":      "Any",
						"source_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port_group": "ANY",
					"dest_port_type":  "group",
					"dest_port":       "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port_group": "ANY",
						"dest_port_type":  "group",
						"dest_port":       "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       "100/200",
					"dest_port_type":  "port",
					"dest_port_group": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "100/200",
						"dest_port_type":  "port",
						"dest_port_group": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proto": "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "TCP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "放行TCP流量",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "放行TCP流量",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "100.1.1.0/24",
					"application_name": "ANY",
					"description":      "放行所有流量",
					"source_type":      "net",
					"acl_action":       "accept",
					"destination_type": "net",
					"direction":        "out",
					"source":           "1.2.3.0/24",
					"dest_port_type":   "port",
					"proto":            "ANY",
					"dest_port":        "0/0",
					"release":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "100.1.1.0/24",
						"application_name": "ANY",
						"description":      "放行所有流量",
						"source_type":      "net",
						"acl_action":       "accept",
						"destination_type": "net",
						"direction":        "out",
						"source":           "1.2.3.0/24",
						"dest_port_type":   "port",
						"proto":            "ANY",
						"dest_port":        "0/0",
						"release":          "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"ip_version", "source_ip"},
			},
		},
	})
}

func TestAccAlicloudCloudFirewallControlPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallControlPolicyBasicDependence0)
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
					"application_name": "ANY",
					"acl_action":       "accept",
					"description":      "放行所有流量",
					"destination_type": "net",
					"destination":      "100.1.1.0/24",
					"direction":        "in",
					"proto":            "ANY",
					"source":           "1.2.3.0/24",
					"source_type":      "net",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": "ANY",
						"acl_action":       "accept",
						"description":      "放行所有流量",
						"destination_type": "net",
						"destination":      "100.1.1.0/24",
						"direction":        "in",
						"proto":            "ANY",
						"source":           "1.2.3.0/24",
						"source_type":      "net",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"release": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"release": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action": "log",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action": "log",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "Any",
					"destination_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "Any",
						"destination_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":      "Any",
					"source_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":      "Any",
						"source_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":      "100/200",
					"dest_port_type": "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":      "100/200",
						"dest_port_type": "port",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proto": "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "TCP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "放行TCP流量",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "放行TCP流量",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "100.1.1.0/24",
					"application_name": "ANY",
					"description":      "放行所有流量",
					"source_type":      "net",
					"acl_action":       "accept",
					"destination_type": "net",
					"direction":        "out",
					"source":           "1.2.3.0/24",
					"dest_port_type":   "port",
					"proto":            "ANY",
					"dest_port":        "0/0",
					"release":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "100.1.1.0/24",
						"application_name": "ANY",
						"description":      "放行所有流量",
						"source_type":      "net",
						"acl_action":       "accept",
						"destination_type": "net",
						"direction":        "out",
						"source":           "1.2.3.0/24",
						"dest_port_type":   "port",
						"proto":            "ANY",
						"dest_port":        "0/0",
						"release":          "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"ip_version", "source_ip"},
			},
		},
	})
}

var AlicloudCloudFirewallControlPolicyMap0 = map[string]string{

	"release":          "true",
	"source_ip":        NOSET,
	"application_name": "ANY",
	"description":      "放行所有流量",
	"destination":      "100.1.1.0/24",
	"proto":            "ANY",
	"source":           "1.2.3.0/24",
	"acl_action":       "accept",
	"destination_type": "net",
	"source_type":      "net",
}

func AlicloudCloudFirewallControlPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAlicloudCloudFirewallControlPolicy(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"application_name": "ANY",
		"acl_action":       "accept",
		"dest_port":        "0/0",
		"description":      "description",
		"destination_type": "net",
		"destination":      "100.1.1.0/24",
		"direction":        "direction",
		"proto":            "ANY",
		"source":           "1.2.3.0/24",
		"source_type":      "net",
		"dest_port_group":  "group",
		"dest_port_type":   "port",
		"ip_version":       "ip_version",
		"lang":             "lang",
		"source_ip":        "source_ip",
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
		"Policys": []interface{}{
			map[string]interface{}{
				"AclUuid":         "MockAclUuid",
				"Direction":       "direction",
				"AclAction":       "accept",
				"ApplicationName": "ANY",
				"Description":     "description",
				"DestPort":        "dest_port",
				"DestPortGroup":   "ANY",
				"DestPortType":    "group",
				"Destination":     "100.1.1.0/24",
				"DestinationType": "net",
				"Proto":           "proto",
				"Release":         "release",
				"Source":          "source",
				"SourceType":      "source_type",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cloud_firewall_control_policy", "MockAclUuid"))
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
			result["AclUuid"] = "MockAclUuid"
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
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudCloudFirewallControlPolicyCreate(d, rawClient)
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
		err := resourceAlicloudCloudFirewallControlPolicyCreate(d, rawClient)
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
		err := resourceAlicloudCloudFirewallControlPolicyCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("MockAclUuid", ":", "direction"))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudCloudFirewallControlPolicyUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyControlPolicyAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"acl_action", "application_name", "application_name", "description", "destination", "destination_type", "proto", "source", "source_type", "dest_port", "dest_port_group", "dest_port_type", "lang", "release", "source_ip"} {
			switch p["alicloud_cloud_firewall_control_policy"].Schema[key].Type {
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
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, diff)
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
		err := resourceAlicloudCloudFirewallControlPolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyControlPolicyNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"acl_action", "application_name", "application_name", "description", "destination", "destination_type", "proto", "source", "source_type", "dest_port", "dest_port_group", "dest_port_type", "lang", "release", "source_ip"} {
			switch p["alicloud_cloud_firewall_control_policy"].Schema[key].Type {
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
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, diff)
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
		err := resourceAlicloudCloudFirewallControlPolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, nil)
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
		err := resourceAlicloudCloudFirewallControlPolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudCloudFirewallControlPolicyDelete(d, rawClient)
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
		err := resourceAlicloudCloudFirewallControlPolicyDelete(d, rawClient)
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
		patcheDescribeCenInstance := gomonkey.ApplyMethod(reflect.TypeOf(&CbnService{}), "DescribeCenInstance", func(*CbnService, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("ResourceNotfound")
		})
		err := resourceAlicloudCloudFirewallControlPolicyDelete(d, rawClient)
		patches.Reset()
		patcheDescribeCenInstance.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, nil)
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
		err := resourceAlicloudCloudFirewallControlPolicyDelete(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeCloudFirewallControlPolicyNotFound", func(t *testing.T) {
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
		err := resourceAlicloudCloudFirewallControlPolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeCloudFirewallControlPolicyAbnormal", func(t *testing.T) {
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
		err := resourceAlicloudCloudFirewallControlPolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
