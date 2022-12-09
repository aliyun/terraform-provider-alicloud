package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudFirewallVpcFirewallControlPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallvpcfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcFirewallControlPolicyBasicDependence0)
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
					"order":            "1",
					"destination":      "127.0.0.2/32",
					"application_name": "ANY",
					"description":      "${var.name}",
					"source_type":      "net",
					"dest_port":        "80/88",
					"acl_action":       "accept",
					"lang":             "zh",
					"destination_type": "net",
					"source":           "127.0.0.1/32",
					"dest_port_type":   "port",
					"proto":            "TCP",
					"release":          "true",
					"member_uid":       "${data.alicloud_account.default.id}",
					"vpc_firewall_id":  "${alicloud_cen_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order":            "1",
						"destination":      "127.0.0.2/32",
						"application_name": "ANY",
						"description":      name,
						"source_type":      "net",
						"dest_port":        "80/88",
						"acl_action":       "accept",
						"destination_type": "net",
						"source":           "127.0.0.1/32",
						"dest_port_type":   "port",
						"proto":            "TCP",
						"release":          "true",
						"member_uid":       CHECKSET,
						"vpc_firewall_id":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AlicloudCloudFirewallVpcFirewallControlPolicyMap0 = map[string]string{}

func AlicloudCloudFirewallVpcFirewallControlPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {}

resource "alicloud_cen_instance" "default" {
	cen_instance_name = var.name
	description = "tf-testAccCenConfigDescription"
	tags 		= {
		Created = "TF"
		For 	= "acceptance test"
	}
}
`, name)
}

func TestAccAlicloudCloudFirewallVpcFirewallControlPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcFirewallControlPolicyMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallvpcfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcFirewallControlPolicyBasicDependence1)
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
					"order":            "1",
					"destination":      "127.0.0.2/32",
					"application_name": "ANY",
					"description":      "${var.name}",
					"source_type":      "net",
					"dest_port":        "80/88",
					"acl_action":       "accept",
					"lang":             "zh",
					"destination_type": "net",
					"source":           "127.0.0.1/32",
					"dest_port_type":   "port",
					"proto":            "TCP",
					"vpc_firewall_id":  "${alicloud_cen_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order":            "1",
						"application_name": "ANY",
						"description":      name,
						"source_type":      "net",
						"source":           "127.0.0.1/32",
						"dest_port_type":   "port",
						"dest_port":        "80/88",
						"proto":            "TCP",
						"destination_type": "net",
						"destination":      "127.0.0.2/32",
						"acl_action":       "accept",
						"vpc_firewall_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proto": "UDP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "UDP",
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
					"application_name": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action": "drop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action": "drop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "127.0.0.3/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": "127.0.0.3/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_type": "group",
					"destination":      "${alicloud_cloud_firewall_address_book.default.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_type": "group",
						"destination":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "127.0.0.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "127.0.0.2/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type": "group",
					"source":      "${alicloud_cloud_firewall_address_book.default.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type": "group",
						"source":      CHECKSET,
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
					"dest_port": "20/22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port": "20/22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "127.0.0.2/32",
					"application_name": "ANY",
					"description":      "${var.name}",
					"source_type":      "net",
					"dest_port":        "80/88",
					"acl_action":       "accept",
					"destination_type": "net",
					"source":           "127.0.0.1/32",
					"dest_port_type":   "port",
					"release":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "127.0.0.2/32",
						"application_name": "ANY",
						"description":      name,
						"source_type":      "net",
						"dest_port":        "80/88",
						"acl_action":       "accept",
						"destination_type": "net",
						"source":           "127.0.0.1/32",
						"dest_port_type":   "port",
						"release":          "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AlicloudCloudFirewallVpcFirewallControlPolicyMap1 = map[string]string{}

func AlicloudCloudFirewallVpcFirewallControlPolicyBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cloud_firewall_address_book" "default" {
  description      = "tf-testAccAddressBook"
  group_name       = "${var.name}"
  group_type       = "ip"
  address_list     = ["10.21.0.0/16", "10.168.0.0/16"]
  auto_add_tag_ecs = 0
}

resource "alicloud_cen_instance" "default" {
	cen_instance_name = var.name
	description = "tf-testAccCenConfigDescription"
	tags 		= {
		Created = "TF"
		For 	= "acceptance test"
	}
}
`, name)
}

func TestUnitAlicloudCloudFirewallVpcFirewallControlPolicy(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cloud_firewall_vpc_firewall_control_policy"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cloud_firewall_vpc_firewall_control_policy"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"acl_action":       "CreateVpcFirewallControlPolicyValue",
		"application_name": "CreateVpcFirewallControlPolicyValue",
		"description":      "CreateVpcFirewallControlPolicyValue",
		"dest_port":        "CreateVpcFirewallControlPolicyValue",
		"dest_port_group":  "CreateVpcFirewallControlPolicyValue",
		"dest_port_type":   "CreateVpcFirewallControlPolicyValue",
		"destination":      "CreateVpcFirewallControlPolicyValue",
		"destination_type": "CreateVpcFirewallControlPolicyValue",
		"lang":             "CreateVpcFirewallControlPolicyValue",
		"member_uid":       "CreateVpcFirewallControlPolicyValue",
		"order":            1,
		"proto":            "CreateVpcFirewallControlPolicyValue",
		"release":          true,
		"source":           "CreateVpcFirewallControlPolicyValue",
		"source_type":      "CreateVpcFirewallControlPolicyValue",
		"vpc_firewall_id":  "CreateVpcFirewallControlPolicyValue",
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
		// DescribeVpcFirewallControlPolicy
		"Policys": []interface{}{
			map[string]interface{}{
				"AclAction":             "CreateVpcFirewallControlPolicyValue",
				"AclUuid":               "CreateVpcFirewallControlPolicyValue",
				"ApplicationId":         "DefaultValue",
				"ApplicationName":       "CreateVpcFirewallControlPolicyValue",
				"Description":           "CreateVpcFirewallControlPolicyValue",
				"DestPort":              "CreateVpcFirewallControlPolicyValue",
				"DestPortGroup":         "CreateVpcFirewallControlPolicyValue",
				"DestPortGroupPorts":    []interface{}{},
				"DestPortType":          "CreateVpcFirewallControlPolicyValue",
				"Destination":           "CreateVpcFirewallControlPolicyValue",
				"DestinationGroupCidrs": []interface{}{},
				"DestinationGroupType":  "DefaultValue",
				"DestinationType":       "CreateVpcFirewallControlPolicyValue",
				"HitTimes":              0,
				"MemberUid":             "CreateVpcFirewallControlPolicyValue",
				"Order":                 1,
				"Proto":                 "CreateVpcFirewallControlPolicyValue",
				"Release":               "true",
				"Source":                "CreateVpcFirewallControlPolicyValue",
				"SourceGroupCidrs":      []interface{}{},
				"SourceGroupType":       "DefaultValue",
				"SourceType":            "CreateVpcFirewallControlPolicyValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateVpcFirewallControlPolicy
		"AclUuid": "CreateVpcFirewallControlPolicyValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cloud_firewall_vpc_firewall_control_policy", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfirewallClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudCloudFirewallVpcFirewallControlPolicyCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeVpcFirewallControlPolicy Response
		"Policys": []interface{}{
			map[string]interface{}{
				"AclUuid": "CreateVpcFirewallControlPolicyValue",
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVpcFirewallControlPolicy" {
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
		err := resourceAlicloudCloudFirewallVpcFirewallControlPolicyCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cloud_firewall_vpc_firewall_control_policy"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfirewallClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudCloudFirewallVpcFirewallControlPolicyUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyVpcFirewallControlPolicy
	attributesDiff := map[string]interface{}{
		"acl_action":       "ModifyVpcFirewallControlPolicyValue",
		"application_name": "ModifyVpcFirewallControlPolicyValue",
		"description":      "ModifyVpcFirewallControlPolicyValue",
		"dest_port":        "ModifyVpcFirewallControlPolicyValue",
		"dest_port_group":  "ModifyVpcFirewallControlPolicyValue",
		"dest_port_type":   "ModifyVpcFirewallControlPolicyValue",
		"destination":      "ModifyVpcFirewallControlPolicyValue",
		"destination_type": "ModifyVpcFirewallControlPolicyValue",
		"proto":            "ModifyVpcFirewallControlPolicyValue",
		"release":          false,
		"source":           "ModifyVpcFirewallControlPolicyValue",
		"source_type":      "ModifyVpcFirewallControlPolicyValue",
	}
	diff, err := newInstanceDiff("alicloud_cloud_firewall_vpc_firewall_control_policy", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cloud_firewall_vpc_firewall_control_policy"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeVpcFirewallControlPolicy Response
		"Policys": []interface{}{
			map[string]interface{}{
				"AclAction":       "ModifyVpcFirewallControlPolicyValue",
				"ApplicationName": "ModifyVpcFirewallControlPolicyValue",
				"Description":     "ModifyVpcFirewallControlPolicyValue",
				"DestPort":        "ModifyVpcFirewallControlPolicyValue",
				"DestPortGroup":   "ModifyVpcFirewallControlPolicyValue",
				"DestPortType":    "ModifyVpcFirewallControlPolicyValue",
				"Destination":     "ModifyVpcFirewallControlPolicyValue",
				"DestinationType": "ModifyVpcFirewallControlPolicyValue",
				"Proto":           "ModifyVpcFirewallControlPolicyValue",
				"Release":         "ModifyVpcFirewallControlPolicyValue",
				"Source":          "ModifyVpcFirewallControlPolicyValue",
				"SourceType":      "ModifyVpcFirewallControlPolicyValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyVpcFirewallControlPolicy" {
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
		err := resourceAlicloudCloudFirewallVpcFirewallControlPolicyUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cloud_firewall_vpc_firewall_control_policy"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeVpcFirewallControlPolicy" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
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
		err := resourceAlicloudCloudFirewallVpcFirewallControlPolicyRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfirewallClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudCloudFirewallVpcFirewallControlPolicyDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_cloud_firewall_vpc_firewall_control_policy", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cloud_firewall_vpc_firewall_control_policy"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVpcFirewallControlPolicy" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudCloudFirewallVpcFirewallControlPolicyDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
