package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"log"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaListener_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AlicloudGaListenerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaListenerBasicDependence)
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
					"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
					"description":    "create_description",
					"name":           "${var.name}",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "70",
						},
					},
					"proxy_protocol": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id": CHECKSET,
						"description":    "create_description",
						"name":           name,
						"port_ranges.#":  "1",
						"proxy_protocol": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accelerator_id", "proxy_protocol"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "SOURCE_IP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "SOURCE_IP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "100",
							"to_port":   "110",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_ranges.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "NONE",
					"protocol":        "UDP",
					"proxy_protocol":  "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "NONE",
						"protocol":        "UDP",
						"proxy_protocol":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "SOURCE_IP",
					"description":     "create_description",
					"protocol":        "TCP",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "70",
						},
					},
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "SOURCE_IP",
						"description":     "create_description",
						"protocol":        "TCP",
						"port_ranges.#":   "1",
						"name":            name,
					}),
				),
			},
		},
	})
}

var AlicloudGaListenerMap = map[string]string{
	"client_affinity": "NONE",
	"protocol":        "TCP",
	"status":          CHECKSET,
}

func AlicloudGaListenerBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_ga_accelerators" "default"{
  
}
`, name)
}

func TestAccAlicloudGaListener_unit(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id": "CreateListenerValue",
		"certificates": []interface{}{
			map[string]interface{}{
				"id": "CreateListenerValue",
			},
		},
		"client_affinity": "CreateListenerValue",
		"description":     "CreateListenerValue",
		"name":            "CreateListenerValue",
		"port_ranges": []interface{}{
			map[string]interface{}{
				"from_port": 10,
				"to_port":   10,
			},
		},
		"protocol":       "CreateListenerValue",
		"proxy_protocol": true,
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
		// DescribeListener
		"AclType": "DefaultValue",
		"Certificates": []interface{}{
			map[string]interface{}{
				"Id": "CreateListenerValue",
			},
		},
		"ClientAffinity": "CreateListenerValue",
		"CreateTime":     "DefaultValue",
		"Description":    "CreateListenerValue",
		"ListenerId":     "MockListenerId",
		"Name":           "CreateListenerValue",
		"Protocol":       "CreateListenerValue",
		"State":          "active",
		"PortRanges": []interface{}{
			map[string]interface{}{
				"FromPort": 10,
				"ToPort":   10,
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateListener
		"ListenerId": "MockListenerId",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_listener", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudGaListenerCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeListener Response
			"ListenerId": "MockListenerId",
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "StateError.Accelerator", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CreateListener" {
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
			err := resourceAlicloudGaListenerCreate(dInit, rawClient)
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dInit.State(), nil)
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudGaListenerUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		// UpdateListener
		attributesDiff := map[string]interface{}{
			"certificates": []interface{}{
				map[string]interface{}{
					"id": "UpdateListenerValue",
				},
			},
			"client_affinity": "UpdateListenerValue",
			"description":     "UpdateListenerValue",
			"name":            "UpdateListenerValue",
			"port_ranges": []interface{}{
				map[string]interface{}{
					"from_port": 15,
					"to_port":   15,
				},
			},
			"protocol": "UpdateListenerValue",
		}
		diff, err := newInstanceDiff("alicloud_ga_listener", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeListener Response
			"Certificates": []interface{}{
				map[string]interface{}{
					"Id": "UpdateListenerValue",
				},
			},
			"ClientAffinity": "UpdateListenerValue",
			"Description":    "UpdateListenerValue",
			"Name":           "UpdateListenerValue",
			"Protocol":       "UpdateListenerValue",
			"PortRanges": []interface{}{
				map[string]interface{}{
					"FromPort": 15,
					"ToPort":   15,
				},
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "StateError.Accelerator", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "UpdateListener" {
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
			err := resourceAlicloudGaListenerUpdate(dExisted, rawClient)
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dExisted.State(), nil)
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
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "NotExist.Listener", "{}"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DescribeListener" {
					switch errorCode {
					case "{}", "NotExist.Listener":
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
			err := resourceAlicloudGaListenerRead(dExisted, rawClient)
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "{}", "NotExist.Listener":
				assert.Nil(t, err)
			}
		}
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudGaListenerDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		errorCodes := []string{"NonRetryableError", "Throttling", "StateError.Accelerator", "nil", "Forbidden", "NotActive.Listener", "NotExist.Accelerator", "NotExist.Listener", "Resource.QuotaFull", "UnknownError"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteListener" {
					switch errorCode {
					case "NonRetryableError", "Forbidden", "NotActive.Listener", "NotExist.Accelerator", "NotExist.Listener", "Resource.QuotaFull", "UnknownError":
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
			err := resourceAlicloudGaListenerDelete(dExisted, rawClient)
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "Forbidden", "NotActive.Listener", "NotExist.Accelerator", "NotExist.Listener", "Resource.QuotaFull", "UnknownError":
				assert.Nil(t, err)
			}
		}
	})
}
