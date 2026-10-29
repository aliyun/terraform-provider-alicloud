// collection-order: bypass real_servers -- pre-existing TypeList attribute outside this directory-only migration; reorder+convergent coverage debt predates this change

package ddoscoo_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	ddoscoo "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudDdosCooPort_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_port.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDdoscooPortMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdoscooPort")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 65530)
	name := fmt.Sprintf("tf-testacc%sddoscooport%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDdoscooPortBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"frontend_port":     fmt.Sprint(rand),
					"backend_port":      fmt.Sprint(rand + 1),
					"instance_id":       "${data.alicloud_ddoscoo_instances.default.ids.0}",
					"frontend_protocol": "tcp",
					"real_servers":      []string{"1.1.1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"frontend_port":     fmt.Sprint(rand),
						"backend_port":      fmt.Sprint(rand + 1),
						"frontend_protocol": "tcp",
						"real_servers.#":    "1",
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
					"real_servers": []string{"1.1.1.1", "2.2.2.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{"1.1.1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "1",
					}),
				),
			},
		},
	})
}

var AlicloudDdoscooPortMap0 = map[string]string{
	"instance_id":       acctest.CHECKSET,
	"frontend_protocol": "tcp",
	"real_servers.#":    "1",
}

func AlicloudDdoscooPortBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

data "alicloud_ddoscoo_instances" "default" {}
`, name)
}

func TestAccAliCloudDdosCooPort_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_port.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDdoscooPortMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdoscooPort")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 65530)
	name := fmt.Sprintf("tf-testacc%sddoscooport1%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDdoscooPortBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"frontend_port":     fmt.Sprint(rand),
					"backend_port":      fmt.Sprint(rand + 1),
					"instance_id":       "${data.alicloud_ddoscoo_instances.default.ids.0}",
					"frontend_protocol": "udp",
					"real_servers":      []string{"1.1.1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"frontend_port":     fmt.Sprint(rand),
						"backend_port":      fmt.Sprint(rand + 1),
						"frontend_protocol": "udp",
						"real_servers.#":    "1",
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
					"real_servers": []string{"1.1.1.1", "2.2.2.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{"1.1.1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "1",
					}),
				),
			},
		},
	})
}

var AlicloudDdoscooPortMap1 = map[string]string{
	"instance_id":       acctest.CHECKSET,
	"frontend_protocol": "udp",
	"real_servers.#":    "1",
}

func AlicloudDdoscooPortBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

data "alicloud_ddoscoo_instances" "default" {}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDdoscooPort(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	acctest.CheckoutSupportedRegions(t, true, connectivity.DdoscooSupportedRegions)
	dInit, _ := schema.InternalMap(p["alicloud_ddoscoo_port"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ddoscoo_port"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"frontend_port":     "CreatePortValue",
		"backend_port":      "443",
		"instance_id":       "CreatePortValue",
		"frontend_protocol": "CreatePortValue",
		"real_servers":      []string{"CreatePortValue"},
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribePort
		"NetworkRules": []interface{}{
			map[string]interface{}{
				"FrontendPort":     "CreatePortValue",
				"FrontendProtocol": "CreatePortValue",
				"InstanceIds":      "CreatePortValue",
				"BackendPort":      "443",
				"RealServers":      "CreatePortValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ddoscoo_port", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdoscooClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_ddoscoo_port"].Create(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreatePort" {
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
		err := p["alicloud_ddoscoo_port"].Create(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ddoscoo_port"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdoscooClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_ddoscoo_port"].Update(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	//ModifyPort
	attributesDiff := map[string]interface{}{
		"real_servers": []string{"ModifyPortValue"},
		"backend_port": "80",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_ddoscoo_port", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ddoscoo_port"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribePort Response
		"NetworkRules": []interface{}{
			map[string]interface{}{
				"BackendPort": "80",
				"RealServers": "ModifyPortValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyPort" {
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
		err := p["alicloud_ddoscoo_port"].Update(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ddoscoo_port"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	//Read
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_ddoscoo_port", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ddoscoo_port"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribePort" {
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
		err := p["alicloud_ddoscoo_port"].Read(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdoscooClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_ddoscoo_port"].Delete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_ddoscoo_port", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ddoscoo_port"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeletePort" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Success": true,
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := p["alicloud_ddoscoo_port"].Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test DdosCoo Port. >>> Resource test cases, automatically generated.
// Case 端口资源测试 7456
func TestAccAliCloudDdosCooPort_basic7456(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_port.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDdosCooPortMap7456)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosCooPort")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sddoscooport%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDdosCooPortBasicDependence7456)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${data.alicloud_ddoscoo_instances.default.ids.0}",
					"backend_port": "8000",
					"real_servers": []string{
						"1.1.1.1"},
					"config": []map[string]interface{}{
						{
							"persistence_timeout": "900",
						},
					},
					"frontend_port":     "8082",
					"frontend_protocol": "tcp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       acctest.CHECKSET,
						"backend_port":      "8000",
						"real_servers.#":    "1",
						"frontend_port":     "8082",
						"frontend_protocol": "tcp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{
						"1.1.1.1", "2.2.2.2", "3.3.3.3"},
					"config": []map[string]interface{}{
						{
							"persistence_timeout": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{
						"5.5.5.5", "6.6.6.6", "7.7.7.7"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{
						"5.5.5.5", "6.6.6.6"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDdosCooPortMap7456 = map[string]string{}

func AlicloudDdosCooPortBasicDependence7456(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_ddoscoo_instances" "default" {}

`, name)
}

// Test DdosCoo Port. <<< Resource test cases, automatically generated.
