package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAliCloudVPCForwardEntry(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_forward_entry"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_forward_entry"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"forward_table_id": "forward_table_id",
		"external_ip":      "external_ip",
		"external_port":    "80",
		"internal_ip":      "internal_ip",
		"internal_port":    "8080",
		"ip_protocol":      "tcp",
		"name":             "name",
		"port_break":       false,
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
		"ForwardTableEntries": map[string]interface{}{
			"ForwardTableEntry": []interface{}{
				map[string]interface{}{
					"ForwardEntryName": "forward_entry_name",
					"ExternalPort":     "external_port",
					"Name":             "forward_entry_name",
					"ExternalIp":       "external_ip",
					"ForwardTableId":   "MockForwardTableId",
					"InternalIp":       "internal_ip",
					"InternalPort":     "internal_port",
					"IpProtocol":       "ip_protocol",
					"Status":           "Available",
					"ForwardEntryId":   "MockForwardEntryId",
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
			result := ReadMockResponse
			return result, nil
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
			result["ForwardEntryId"] = "MockForwardEntryId"
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
		"DescribeForwardTableEntriesNotFound": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"ForwardTableEntries": map[string]interface{}{
					"ForwardTableEntry": []interface{}{},
				},
			}
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudNatGatewayForwardEntryCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("InvalidIp.NotInNatgw")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudNatGatewayForwardEntryCreate(d, rawClient)
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
		err := resourceAliCloudNatGatewayForwardEntryCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateNonRetryableError", func(t *testing.T) {
		retryFlag := false
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
		err := resourceAliCloudNatGatewayForwardEntryCreate(dCreate, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("MockForwardTableId", ":", "MockForwardEntryId"))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudNatGatewayForwardEntryUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateForwardEntryAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"external_ip", "external_port", "forward_entry_name", "name", "internal_ip", "internal_port", "ip_protocol", "port_break"} {
			switch p["alicloud_forward_entry"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_forward_entry"].Schema).Data(nil, diff)
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
		err := resourceAliCloudNatGatewayForwardEntryUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateForwardEntryNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"external_ip", "external_port", "forward_entry_name", "name", "internal_ip", "internal_port", "ip_protocol", "port_break"} {
			switch p["alicloud_forward_entry"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_forward_entry"].Schema).Data(nil, diff)
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
		err := resourceAliCloudNatGatewayForwardEntryUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_forward_entry"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
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
		err := resourceAliCloudNatGatewayForwardEntryUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudNatGatewayForwardEntryDelete(d, rawClient)
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
		err := resourceAliCloudNatGatewayForwardEntryDelete(d, rawClient)
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
		patcheDescribeForwardEntry := gomonkey.ApplyMethod(reflect.TypeOf(&VpcService{}), "DescribeForwardEntry", func(*VpcService, string) (map[string]interface{}, error) {
			return responseMock["RetryError"]("InvalidForwardEntryId.NotFound")
		})
		err := resourceAliCloudNatGatewayForwardEntryDelete(d, rawClient)
		patches.Reset()
		patcheDescribeForwardEntry.Reset()
		assert.NotNil(t, err)
	})

	t.Run("DeleteMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_forward_entry"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
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
		err := resourceAliCloudNatGatewayForwardEntryDelete(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("DescribeForwardTableEntriesNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["DescribeForwardTableEntriesNotFound"]("")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudNatGatewayForwardEntryRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeDescribeVpcNatIpAbnormal", func(t *testing.T) {
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
		err := resourceAliCloudNatGatewayForwardEntryRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})

	t.Run("ReadMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_forward_entry"].Schema).Data(nil, nil)
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
		err := resourceAliCloudNatGatewayForwardEntryRead(resourceData1, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
}

// Test NatGateway ForwardEntry. >>> Resource test cases, automatically generated.

// Case 全生命周期_ForwardEntry_1.0.1 9600
func TestAccAliCloudNatGatewayForwardEntry_basic9600(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_forward_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayForwardEntryMap9600)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGatewayForwardEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnatgateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayForwardEntryBasicDependence9600)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"external_port":      "80",
					"external_ip":        "${alicloud_vpc_nat_ip.NatIp1.nat_ip}",
					"ip_protocol":        "tcp",
					"internal_port":      "8080",
					"internal_ip":        "172.16.0.115",
					"forward_table_id":   "${alicloud_nat_gateway.NATGateway.forward_table_ids}",
					"forward_entry_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port":      "80",
						"external_ip":        CHECKSET,
						"ip_protocol":        "tcp",
						"internal_port":      "8080",
						"internal_ip":        "172.16.0.115",
						"forward_table_id":   CHECKSET,
						"forward_entry_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_ip": "${alicloud_vpc_nat_ip.NatIp2.nat_ip}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_ip": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_port": "90",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port": "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forward_entry_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forward_entry_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internal_ip": "172.16.0.55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internal_ip": "172.16.0.55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internal_port": "9090",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internal_port": "9090",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol": "udp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol": "udp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_break": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_break"},
			},
		},
	})
}

func TestAccAliCloudNatGatewayForwardEntry_basic9600_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_forward_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayForwardEntryMap9600)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGatewayForwardEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnatgateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayForwardEntryBasicDependence9600)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"external_port":      "80",
					"external_ip":        "${alicloud_vpc_nat_ip.NatIp1.nat_ip}",
					"ip_protocol":        "tcp",
					"internal_port":      "8080",
					"internal_ip":        "172.16.0.115",
					"forward_table_id":   "${alicloud_nat_gateway.NATGateway.forward_table_ids}",
					"forward_entry_name": name,
					"port_break":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port":      "80",
						"external_ip":        CHECKSET,
						"ip_protocol":        "tcp",
						"internal_port":      "8080",
						"internal_ip":        "172.16.0.115",
						"forward_table_id":   CHECKSET,
						"forward_entry_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_break"},
			},
		},
	})
}

// Case 适配字段name 9602
func TestAccAliCloudNatGatewayForwardEntry_basic9602(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_forward_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayForwardEntryMap9600)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGatewayForwardEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnatgateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayForwardEntryBasicDependence9600)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"external_port":    "80",
					"external_ip":      "${alicloud_vpc_nat_ip.NatIp1.nat_ip}",
					"ip_protocol":      "tcp",
					"internal_port":    "8080",
					"internal_ip":      "172.16.0.115",
					"forward_table_id": "${alicloud_nat_gateway.NATGateway.forward_table_ids}",
					"name":             name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port":    "80",
						"external_ip":      CHECKSET,
						"ip_protocol":      "tcp",
						"internal_port":    "8080",
						"internal_ip":      "172.16.0.115",
						"forward_table_id": CHECKSET,
						"name":             name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_ip": "${alicloud_vpc_nat_ip.NatIp2.nat_ip}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_ip": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_port": "90",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port": "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internal_ip": "172.16.0.55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internal_ip": "172.16.0.55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internal_port": "9090",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internal_port": "9090",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol": "udp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol": "udp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_break": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_break"},
			},
		},
	})
}

func TestAccAliCloudNatGatewayForwardEntry_basic9602_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_forward_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayForwardEntryMap9600)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGatewayForwardEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnatgateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayForwardEntryBasicDependence9600)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"external_port":    "80",
					"external_ip":      "${alicloud_vpc_nat_ip.NatIp1.nat_ip}",
					"ip_protocol":      "tcp",
					"internal_port":    "8080",
					"internal_ip":      "172.16.0.115",
					"forward_table_id": "${alicloud_nat_gateway.NATGateway.forward_table_ids}",
					"name":             name,
					"port_break":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port":    "80",
						"external_ip":      CHECKSET,
						"ip_protocol":      "tcp",
						"internal_port":    "8080",
						"internal_ip":      "172.16.0.115",
						"forward_table_id": CHECKSET,
						"name":             name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_break"},
			},
		},
	})
}

var AliCloudNatGatewayForwardEntryMap9600 = map[string]string{
	"status":           CHECKSET,
	"forward_entry_id": CHECKSET,
}

func AliCloudNatGatewayForwardEntryBasicDependence9600(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "tf-test-nat-forward-entry"
}

resource "alicloud_vswitch" "VSwitch" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = "eu-central-1b"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "tf-test-nat-forward-entry-vsw"
}

resource "alicloud_nat_gateway" "NATGateway" {
  vpc_id           = alicloud_vpc.vpc.id
  nat_gateway_name = var.name
  nat_type         = "Enhanced"
  vswitch_id       = alicloud_vswitch.VSwitch.id
  network_type     = "intranet"
}

resource "alicloud_vpc_nat_ip" "NatIp1" {
  nat_ip         = "172.16.0.66"
  nat_ip_name    = "tf-test-dnat-natip1"
  nat_gateway_id = alicloud_nat_gateway.NATGateway.id
  nat_ip_cidr    = alicloud_vswitch.VSwitch.cidr_block
}

resource "alicloud_vpc_nat_ip" "NatIp2" {
  nat_ip         = "172.16.0.88"
  nat_ip_name    = "tf-test-dnat-natip2"
  nat_gateway_id = alicloud_vpc_nat_ip.NatIp1.nat_gateway_id
  nat_ip_cidr    = alicloud_vswitch.VSwitch.cidr_block
}


`, name)
}

// Test NatGateway ForwardEntry. <<< Resource test cases, automatically generated.

func init() {
	resource.AddTestSweepers("alicloud_forward_entry", &resource.Sweeper{
		Name: "alicloud_forward_entry",
		F:    testSweepForwardEntries,
	})
}

func testSweepForwardEntries(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-test-",
		"tf_test_",
		"tf-example",
		"tf_example",
	}

	natGateways, err := sweepListNatGateways(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}

	for _, nat := range natGateways {
		natName := fmt.Sprint(nat["Name"])
		natID := fmt.Sprint(nat["NatGatewayId"])
		// If the parent NAT gateway is explicitly do-not-delete, protect its entries too.
		if protectedSweepVpcName(natName) {
			log.Printf("[INFO] Skipping ForwardEntries of NAT Gateway %s (%s): do-not-delete parent", natName, natID)
			continue
		}
		for _, forwardTableId := range sweepNatTableIds(nat, "ForwardTableIds", "ForwardTableId") {
			entries, err := sweepListForwardEntries(client, forwardTableId)
			if err != nil {
				log.Printf("[ERROR] Failed to list ForwardEntries for table %s: %s", forwardTableId, err)
				continue
			}
			for _, entry := range entries {
				name := fmt.Sprint(entry["ForwardEntryName"])
				forwardEntryId := fmt.Sprint(entry["ForwardEntryId"])
				if protectedSweepVpcName(name) {
					log.Printf("[INFO] Skipping ForwardEntry %s (%s): do-not-delete name", name, forwardEntryId)
					continue
				}
				if !sweepAll() {
					skip := true
					for _, prefix := range prefixes {
						if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
							skip = false
							break
						}
					}
					if skip {
						log.Printf("[INFO] Skipping ForwardEntry: %s (%s)", name, forwardEntryId)
						continue
					}
				}
				log.Printf("[INFO] Deleting ForwardEntry: %s (%s) in table %s", name, forwardEntryId, forwardTableId)
				action := "DeleteForwardEntry"
				request := map[string]interface{}{
					"ForwardTableId": forwardTableId,
					"ForwardEntryId": forwardEntryId,
					"RegionId":       client.RegionId,
					"ClientToken":    buildClientToken(action),
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(5*time.Minute, func() *retry.RetryError {
					_, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectStatus.NATGW", "OperationConflict", "UnknownError"}) || NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					if NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete ForwardEntry %s (%s): %s", name, forwardEntryId, err)
				}
			}
		}
	}
	return nil
}

func sweepListForwardEntries(client *connectivity.AliyunClient, forwardTableId string) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	action := "DescribeForwardTableEntries"
	request := map[string]interface{}{
		"ForwardTableId": forwardTableId,
		"RegionId":       client.RegionId,
		"PageSize":       PageSizeLarge,
		"PageNumber":     1,
	}
	for {
		response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			return nil, err
		}
		resp, err := jsonpath.Get("$.ForwardTableEntries.ForwardTableEntry", response)
		if err != nil {
			return nil, err
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return items, nil
}
