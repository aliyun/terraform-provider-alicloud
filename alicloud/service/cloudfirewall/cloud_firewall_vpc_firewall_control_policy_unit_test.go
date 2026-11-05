package cloudfirewall

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	client "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	tea "github.com/alibabacloud-go/tea/tea"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	assert "github.com/stretchr/testify/assert"
)

// collection-order: bypass application_name_list -- pre-existing list unchanged by this PR; fixes cover scalar release fixtures and order updates only
// collection-order: bypass repeat_days -- pre-existing list unchanged by this PR; fixes cover scalar release fixtures and order updates only

// lintignore: R001
func TestUnitAliCloudCloudFirewallVpcFirewallControlPolicy(t *testing.T) {
	p := unitTestProviderResources()
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
		"release":          "true",
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
	rawClient, err := unitTestSharedClientForRegion(region)
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
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_cloud_firewall_vpc_firewall_control_policy", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfirewallClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = resourceVpcFirewallControlPolicyCreate(dInit, rawClient)
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
		err := resourceVpcFirewallControlPolicyCreate(dInit, rawClient)
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
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = resourceVpcFirewallControlPolicyUpdate(dExisted, rawClient)
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
		err := resourceVpcFirewallControlPolicyUpdate(dExisted, rawClient)
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
		err := resourceVpcFirewallControlPolicyRead(dExisted, rawClient)
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
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = resourceVpcFirewallControlPolicyDelete(dExisted, rawClient)
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
		err := resourceVpcFirewallControlPolicyDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
