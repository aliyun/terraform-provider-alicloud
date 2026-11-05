package expressconnect_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnect"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAliCloudExpressConnectBgpNetwork(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_vpc_bgp_network"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_vpc_bgp_network"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"dst_cidr_block": "192.168.0.0/24",
		"router_id":      "router_id",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"BgpNetworks": map[string]interface{}{
			"BgpNetwork": []interface{}{
				map[string]interface{}{
					"DstCidrBlock": "192.168.0.0/24",
					"Status":       "Available",
					"RouterId":     "router_id",
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_vpc_bgp_network", "router_id:192.168.0.0/24"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkCreate(d, rawClient)
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
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkCreate(d, rawClient)
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
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("router_id:192.168.0.0/24")

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
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
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		patchDescribeExpressConnectBgpNetwork := gomonkey.ApplyMethod(reflect.TypeOf(&expressconnect.ExpressConnectServiceV2{}), "DescribeExpressConnectBgpNetwork", func(*expressconnect.ExpressConnectServiceV2, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkDelete(d, rawClient)
		patches.Reset()
		patchDescribeExpressConnectBgpNetwork.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeExpressConnectBgpNetworkNotFound", func(t *testing.T) {
		patchRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkRead(d, rawClient)
		patchRequest.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeExpressConnectBgpNetworkAbnormal", func(t *testing.T) {
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
		err := expressconnect.ResourceAliCloudExpressConnectBgpNetworkRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test ExpressConnect BgpNetwork. >>> Resource test cases, automatically generated.
// Case BGP Network V1 11500
func TestAccAliCloudExpressConnectBgpNetwork_basic11500(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_bgp_network.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectBgpNetworkMap11500)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectBgpNetwork")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectBgpNetworkBasicDependence11500)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dst_cidr_block": "192.168.0.0/16",
					"router_id":      "${data.alicloud_express_connect_router_interfaces.default.interfaces.0.router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_cidr_block": "192.168.0.0/16",
						"router_id":      acctest.CHECKSET,
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

func TestAccAliCloudExpressConnectBgpNetwork_basic11500_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_bgp_network.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectBgpNetworkMap11500)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectBgpNetwork")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectBgpNetworkBasicDependence11500)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dst_cidr_block": "192.168.0.0/16",
					"router_id":      "${data.alicloud_express_connect_router_interfaces.default.interfaces.0.router_id}",
					"vpc_id":         "${data.alicloud_express_connect_router_interfaces.default.interfaces.0.opposite_vpc_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_cidr_block": "192.168.0.0/16",
						"router_id":      acctest.CHECKSET,
						"vpc_id":         acctest.CHECKSET,
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

var AliCloudExpressConnectBgpNetworkMap11500 = map[string]string{
	"status": acctest.CHECKSET,
}

func AliCloudExpressConnectBgpNetworkBasicDependence11500(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_express_connect_router_interfaces" "default" {
  		name_regex = "^default-NODELETING-VBR$"
	}
`, name)
}

// Test ExpressConnect BgpNetwork. <<< Resource test cases, automatically generated.
