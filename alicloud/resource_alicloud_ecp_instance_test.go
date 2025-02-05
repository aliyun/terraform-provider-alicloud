package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEcpInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	resourceId := "alicloud_ecp_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcpInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcpInstanceBasicDependence0)
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
					"instance_type":     "${data.alicloud_ecp_instance_types.default.instance_types.0.instance_type}",
					"image_id":          "android-image-release5501072_a11_20240530.raw",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     CHECKSET,
						"image_id":          "android-image-release5501072_a11_20240530.raw",
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resolution": "1280*720",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resolution": "1280*720",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": "${data.alicloud_ecp_key_pairs.default.pairs.0.key_pair_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vnc_password": "Ps1688",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
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
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vnc_password", "auto_pay", "period", "period_unit", "auto_renew", "force"},
			},
		},
	})
}

func TestAccAliCloudEcpInstance_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	resourceId := "alicloud_ecp_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcpInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcpInstanceBasicDependence0)
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
					"instance_type":     "${data.alicloud_ecp_instance_types.default.instance_types.0.instance_type}",
					"image_id":          "android-image-release5501072_a11_20240530.raw",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"eip_bandwidth":     "20",
					"resolution":        "1280*720",
					"key_pair_name":     "${data.alicloud_ecp_key_pairs.default.pairs.0.key_pair_name}",
					"vnc_password":      "Ps1688",
					"payment_type":      "PayAsYouGo",
					"instance_name":     name,
					"description":       name,
					"status":            "Stopped",
					"force":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     CHECKSET,
						"image_id":          "android-image-release5501072_a11_20240530.raw",
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
						"eip_bandwidth":     "20",
						"resolution":        "1280*720",
						"key_pair_name":     CHECKSET,
						"payment_type":      "PayAsYouGo",
						"instance_name":     name,
						"description":       name,
						"status":            "Stopped",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vnc_password", "auto_pay", "period", "period_unit", "auto_renew", "force"},
			},
		},
	})
}

func TestAccAliCloudEcpInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	resourceId := "alicloud_ecp_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcpInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcpInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":     "${data.alicloud_ecp_instance_types.default.instance_types.0.instance_type}",
					"image_id":          "android-image-release5501072_a11_20240530.raw",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"payment_type":      "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     CHECKSET,
						"image_id":          "android-image-release5501072_a11_20240530.raw",
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
						"payment_type":      "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resolution": "1280*720",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resolution": "1280*720",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": "${data.alicloud_ecp_key_pairs.default.pairs.0.key_pair_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vnc_password": "Ps1688",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
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
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vnc_password", "auto_pay", "period", "period_unit", "auto_renew", "force"},
			},
		},
	})
}

func TestAccAliCloudEcpInstance_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	resourceId := "alicloud_ecp_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcpInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcpInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":     "${data.alicloud_ecp_instance_types.default.instance_types.0.instance_type}",
					"image_id":          "android-image-release5501072_a11_20240530.raw",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"eip_bandwidth":     "20",
					"resolution":        "1280*720",
					"key_pair_name":     "${data.alicloud_ecp_key_pairs.default.pairs.0.key_pair_name}",
					"vnc_password":      "Ps1688",
					"payment_type":      "Subscription",
					"auto_pay":          "true",
					"period":            "1",
					"period_unit":       "Month",
					"auto_renew":        "false",
					"instance_name":     name,
					"description":       name,
					"status":            "Stopped",
					"force":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     CHECKSET,
						"image_id":          "android-image-release5501072_a11_20240530.raw",
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
						"eip_bandwidth":     "20",
						"resolution":        "1280*720",
						"key_pair_name":     CHECKSET,
						"payment_type":      "Subscription",
						"instance_name":     name,
						"description":       name,
						"status":            "Stopped",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vnc_password", "auto_pay", "period", "period_unit", "auto_renew", "force"},
			},
		},
	})
}

var AliCloudEcpInstanceMap0 = map[string]string{
	"resolution":    CHECKSET,
	"payment_type":  CHECKSET,
	"instance_name": CHECKSET,
	"status":        CHECKSET,
}

func AliCloudEcpInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ecp_zones" "default" {
	}

	data "alicloud_ecp_instance_types" "default" {
	}

	data "alicloud_ecp_key_pairs" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_ecp_zones.default.zones.0.zone_id
	}

	//data "alicloud_security_groups" "default" {
  	//	vpc_id = data.alicloud_vpcs.default.ids.0
	//}
    resource "alicloud_security_group" "default" {
		security_group_name = var.name
		vpc_id = data.alicloud_vpcs.default.ids.0
	}
`, name)
}

func TestUnitAliCloudEcpInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_name":     "RunInstancesValue",
		"description":       "RunInstancesValue",
		"force":             true,
		"payment_type":      "PayAsYouGo",
		"key_pair_name":     "RunInstancesValue",
		"security_group_id": "RunInstancesValue",
		"vswitch_id":        "RunInstancesValue",
		"image_id":          "RunInstancesValue",
		"instance_type":     "RunInstancesValue",
		"status":            "Running",
		"auto_pay":          true,
		"auto_renew":        true,
		"eip_bandwidth":     100,
		"period":            "1",
		"period_unit":       "Month",
		"resolution":        "RunInstancesValue",
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
		// ListInstances
		"Instances": map[string]interface{}{
			"Instance": []interface{}{
				map[string]interface{}{
					"InstanceId":      "RunInstancesValue",
					"Status":          "Running",
					"Description":     "RunInstancesValue",
					"ImageId":         "RunInstancesValue",
					"InstanceName":    "RunInstancesValue",
					"InstanceType":    "RunInstancesValue",
					"ChargeType":      "PostPaid",
					"KeyPairName":     "RunInstancesValue",
					"Resolution":      "RunInstancesValue",
					"SecurityGroupId": "RunInstancesValue",
					"VpcAttributes": map[string]interface{}{
						"VSwitchId": "RunInstancesValue",
					},
				},
			},
		},
		"InstanceIds": map[string]interface{}{
			"InstanceId": []interface{}{
				"RunInstancesValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// RunInstances
		"InstanceIds": map[string]interface{}{
			"InstanceId": []interface{}{
				"RunInstancesValue",
			},
		},
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecp_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudphoneClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudEcpInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListInstances Response
		"InstanceIds": map[string]interface{}{
			"InstanceId": []interface{}{
				"RunInstancesValue",
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RunInstances" {
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
		err := resourceAliCloudEcpInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudphoneClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudEcpInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateInstanceAttribute
	attributesDiff := map[string]interface{}{
		"instance_name": "UpdateInstanceAttributeValue",
		"description":   "UpdateInstanceAttributeValue",
		"key_pair_name": "UpdateInstanceAttributeValue",
		"vnc_password":  "UpdateInstanceAttributeValue",
	}
	diff, err := newInstanceDiff("alicloud_ecp_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListInstances Response
		"Instances": map[string]interface{}{
			"Instance": []interface{}{
				map[string]interface{}{
					"Description":  "UpdateInstanceAttributeValue",
					"InstanceName": "UpdateInstanceAttributeValue",
					"KeyPairName":  "UpdateInstanceAttributeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateInstanceAttribute" {
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
		err := resourceAliCloudEcpInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// StopInstances
	attributesDiff = map[string]interface{}{
		"status": "Stopped",
	}
	diff, err = newInstanceDiff("alicloud_ecp_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListInstances Response
		"Instances": map[string]interface{}{
			"Instance": []interface{}{
				map[string]interface{}{
					"Status": "Stopped",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "StopInstances" {
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
		err := resourceAliCloudEcpInstanceUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListInstances" {
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
		err := resourceAliCloudEcpInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudphoneClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudEcpInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "CloudPhoneInstances.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteInstances" {
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
		err := resourceAliCloudEcpInstanceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil", "CloudPhoneInstances.NotFound":
			assert.Nil(t, err)
		}
	}

}
