package natgateway_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	natgateway "github.com/aliyun/terraform-provider-alicloud/alicloud/service/natgateway"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAliCloudVPCForwardEntry(t *testing.T) {
	p := acctest.Provider.ResourcesMap
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
	rawClient := unitTestMockClient(t)
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
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := p["alicloud_forward_entry"].Create(d, rawClient)
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
		err := p["alicloud_forward_entry"].Create(d, rawClient)
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
		err := p["alicloud_forward_entry"].Create(dCreate, rawClient)
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
		err := p["alicloud_forward_entry"].Create(dCreate, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("MockForwardTableId", ":", "MockForwardEntryId"))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := p["alicloud_forward_entry"].Update(d, rawClient)
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
		err := p["alicloud_forward_entry"].Update(resourceData1, rawClient)
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
		err := p["alicloud_forward_entry"].Update(resourceData1, rawClient)
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
		err := p["alicloud_forward_entry"].Update(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := p["alicloud_forward_entry"].Delete(d, rawClient)
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
		err := p["alicloud_forward_entry"].Delete(d, rawClient)
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
		// DeleteMockNormal keeps assert.NotNil: the post-Delete Describe query is mocked to
		// return a non-retryable, non-NotFound error, so the state refresh fast-exits via the
		// real error return path (refresh -> WaitForState -> Delete WrapErrorf) without polling.
		patcheDescribeForwardEntry := gomonkey.ApplyMethod(reflect.TypeOf(&natgateway.NATGatewayServiceV2{}), "DescribeNatGatewayForwardEntry", func(*natgateway.NATGatewayServiceV2, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NonRetryableError")
		})
		err := p["alicloud_forward_entry"].Delete(d, rawClient)
		patcheDescribeForwardEntry.Reset()
		patches.Reset()
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
		patcheDescribeForwardEntry := gomonkey.ApplyMethod(reflect.TypeOf(&natgateway.NATGatewayServiceV2{}), "DescribeNatGatewayForwardEntry", func(*natgateway.NATGatewayServiceV2, string) (map[string]interface{}, error) {
			return responseMock["RetryError"]("InvalidForwardEntryId.NotFound")
		})
		err := p["alicloud_forward_entry"].Delete(resourceData1, rawClient)
		patcheDescribeForwardEntry.Reset()
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
		err := p["alicloud_forward_entry"].Read(d, rawClient)
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
		err := p["alicloud_forward_entry"].Read(d, rawClient)
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
		err := p["alicloud_forward_entry"].Read(resourceData1, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
}

// Test NatGateway ForwardEntry. >>> Resource test cases, automatically generated.

// Case 全生命周期_ForwardEntry_1.0.1 9600
func TestAccAliCloudNatGatewayForwardEntry_basic9600(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_forward_entry.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNatGatewayForwardEntryMap9600)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return natgateway.NewNATGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNatGatewayForwardEntry")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnatgateway%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNatGatewayForwardEntryBasicDependence9600)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"external_ip":        acctest.CHECKSET,
						"ip_protocol":        "tcp",
						"internal_port":      "8080",
						"internal_ip":        "172.16.0.115",
						"forward_table_id":   acctest.CHECKSET,
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
						"external_ip": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNatGatewayForwardEntryMap9600)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return natgateway.NewNATGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNatGatewayForwardEntry")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnatgateway%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNatGatewayForwardEntryBasicDependence9600)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"external_ip":        acctest.CHECKSET,
						"ip_protocol":        "tcp",
						"internal_port":      "8080",
						"internal_ip":        "172.16.0.115",
						"forward_table_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNatGatewayForwardEntryMap9600)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return natgateway.NewNATGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNatGatewayForwardEntry")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnatgateway%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNatGatewayForwardEntryBasicDependence9600)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"external_ip":      acctest.CHECKSET,
						"ip_protocol":      "tcp",
						"internal_port":    "8080",
						"internal_ip":      "172.16.0.115",
						"forward_table_id": acctest.CHECKSET,
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
						"external_ip": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNatGatewayForwardEntryMap9600)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return natgateway.NewNATGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNatGatewayForwardEntry")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnatgateway%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNatGatewayForwardEntryBasicDependence9600)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"external_ip":      acctest.CHECKSET,
						"ip_protocol":      "tcp",
						"internal_port":    "8080",
						"internal_ip":      "172.16.0.115",
						"forward_table_id": acctest.CHECKSET,
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
	"status":           acctest.CHECKSET,
	"forward_entry_id": acctest.CHECKSET,
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
