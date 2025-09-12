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

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCenTransitRouterVpcAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenTransitRouterVpcAttachmentSupportRegions)
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterVpcAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterVpcAttachmentBasicDependence0)
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
					"cen_id": "${alicloud_cen_transit_router.default.cen_id}",
					"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${data.alicloud_vswitches.default_master.vswitches.0.id}",
							"zone_id":    "${data.alicloud_vswitches.default_master.vswitches.0.zone_id}",
						},
						{
							"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
							"zone_id":    "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
						},
					},
					"route_table_association_enabled": "false",
					"route_table_propagation_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":          CHECKSET,
						"vpc_id":          CHECKSET,
						"zone_mappings.#": "2",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"auto_publish_route_enabled": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"auto_publish_route_enabled": "true",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${data.alicloud_vswitches.default_master.vswitches.0.id}",
							"zone_id":    "${data.alicloud_vswitches.default_master.vswitches.0.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
							"zone_id":    "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "TransitRouterVpcAttachment",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "TransitRouterVpcAttachment",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"route_table_association_enabled", "route_table_propagation_enabled"},
			},
		},
	})
}

func TestAccAliCloudCenTransitRouterVpcAttachment_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenTransitRouterVpcAttachmentSupportRegions)
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterVpcAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterVpcAttachmentBasicDependence0)
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
					"cen_id":                                "${alicloud_cen_transit_router.default.cen_id}",
					"vpc_id":                                "${data.alicloud_vpcs.default.ids.0}",
					"transit_router_id":                     "${alicloud_cen_transit_router.default.transit_router_id}",
					"resource_type":                         "VPC",
					"payment_type":                          "PayAsYouGo",
					"vpc_owner_id":                          "${data.alicloud_account.default.id}",
					"auto_publish_route_enabled":            "false",
					"transit_router_attachment_name":        name,
					"transit_router_attachment_description": name,
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${data.alicloud_vswitches.default_master.vswitches.0.id}",
							"zone_id":    "${data.alicloud_vswitches.default_master.vswitches.0.zone_id}",
						},
						{
							"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
							"zone_id":    "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "TransitRouterVpcAttachment",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                                CHECKSET,
						"vpc_id":                                CHECKSET,
						"transit_router_id":                     CHECKSET,
						"resource_type":                         "VPC",
						"payment_type":                          "PayAsYouGo",
						"vpc_owner_id":                          CHECKSET,
						"auto_publish_route_enabled":            "false",
						"transit_router_attachment_name":        name,
						"transit_router_attachment_description": name,
						"zone_mappings.#":                       "2",
						"tags.%":                                "2",
						"tags.Created":                          "TF",
						"tags.For":                              "TransitRouterVpcAttachment",
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

var AliCloudCenTransitRouterVpcAttachmentMap0 = map[string]string{
	"transit_router_id":            CHECKSET,
	"resource_type":                CHECKSET,
	"payment_type":                 CHECKSET,
	"vpc_owner_id":                 CHECKSET,
	"transit_router_attachment_id": CHECKSET,
	"status":                       CHECKSET,
	"region_id":                    CHECKSET,
}

func AliCloudCenTransitRouterVpcAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.ids.0
	}

	data "alicloud_vswitches" "default_master" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.ids.1
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
  		protection_level  = "REDUCED"
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id = alicloud_cen_instance.default.id
	}
`, name)
}

func TestUnitAliCloudCenTransitRouterVpcAttachment(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_cen_transit_router_vpc_attachment"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_cen_transit_router_vpc_attachment"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"cen_id":                                "cen_id",
		"transit_router_id":                     "transit_router_id",
		"transit_router_attachment_name":        "transit_router_attachment_name",
		"transit_router_attachment_description": "transit_router_attachment_description",
		"vpc_id":                                "vpc_id",
		"vpc_owner_id":                          "vpc_owner_id",
		"dry_run":                               false,
		"resource_type":                         "VPC",
		"payment_type":                          "PayAsYouGo",
		"route_table_association_enabled":       false,
		"route_table_propagation_enabled":       false,
		"zone_mappings": []map[string]interface{}{
			{
				"vswitch_id": "vswitch_id",
				"zone_id":    "zone_id",
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
		"TransitRouterAttachments": []interface{}{
			map[string]interface{}{
				"Status":                             "Attached",
				"AutoPublishRouteEnabled":            false,
				"ResourceType":                       "VPC",
				"TransitRouterAttachmentDescription": "transit_router_attachment_description",
				"TransitRouterAttachmentName":        "transit_router_attachment_name",
				"TransitRouterId":                    "transit_router_id",
				"VpcId":                              "vpc_id",
				"VpcOwnerId":                         "vpc_owner_id",
				"CenId":                              "cen_id",
				"TransitRouterAttachmentId":          "MockTransitRouterAttachmentId",
				"ChargeType":                         "POSTPAY",
				"ZoneMappings": []interface{}{
					map[string]interface{}{
						"VSwitchId": "vswitch_id",
						"ZoneId":    "zone_id",
					},
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cen_transit_router_vpc_attachment", "MockTransitRouterAttachmentId"))
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
			result["TransitRouterAttachmentId"] = "MockTransitRouterAttachmentId"
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCenTransitRouterVpcAttachmentCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Operation.Blocking")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudCenTransitRouterVpcAttachmentCreate(d, rawClient)
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
		err := resourceAliCloudCenTransitRouterVpcAttachmentCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("CenId", ":", "MockTransitRouterAttachmentId"))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudCenTransitRouterVpcAttachmentUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateTransitRouterVpcAttachmentAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_type", "dry_run", "transit_router_attachment_description", "transit_router_attachment_name"} {
			switch p["alicloud_cen_transit_router_vpc_attachment"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_cen_transit_router_vpc_attachment"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Operation.Blocking")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudCenTransitRouterVpcAttachmentUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateTransitRouterVpcAttachmentAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_type", "dry_run", "transit_router_attachment_description", "transit_router_attachment_name"} {
			switch p["alicloud_cen_transit_router_vpc_attachment"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_cen_transit_router_vpc_attachment"].Schema).Data(nil, diff)
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
		err := resourceAliCloudCenTransitRouterVpcAttachmentUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCenTransitRouterVpcAttachmentDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Operation.Blocking")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudCenTransitRouterVpcAttachmentDelete(d, rawClient)
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
		patcheDescribeCenTransitRouterVpcAttachment := gomonkey.ApplyMethod(reflect.TypeOf(&CbnService{}), "DescribeCenTransitRouterVpcAttachment", func(*CbnService, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("ResourceNotfound")
		})
		err := resourceAliCloudCenTransitRouterVpcAttachmentDelete(d, rawClient)
		patches.Reset()
		patcheDescribeCenTransitRouterVpcAttachment.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeCenTransitRouterVpcAttachmentNotFound", func(t *testing.T) {
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
		err := resourceAliCloudCenTransitRouterVpcAttachmentRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeCenTransitRouterVpcAttachmentAbnormal", func(t *testing.T) {
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
		err := resourceAliCloudCenTransitRouterVpcAttachmentRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Cen TransitRouterVpcAttachment. >>> Resource test cases, automatically generated.
// Case TR支持IPv6_善问_线上_副本1726040691220 7866
func TestAccAliCloudCenTransitRouterVpcAttachment_basic7866(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap7866)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitroutervpcattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterVpcAttachmentBasicDependence7866)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id": "${alicloud_vpc.defaultJLRlxW.id}",
					"cen_id": "${alicloud_cen_instance.defaultJ6HrUE.id}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultoxj0Cs.id}",
							"zone_id":    "${alicloud_vswitch.defaultoxj0Cs.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Dd.id}",
							"zone_id":    "${alicloud_vswitch.defaulteKv3Dd.zone_id}",
						},
					},
					"transit_router_id":                     "${alicloud_cen_transit_router.defaults5WvfD.transit_router_id}",
					"transit_router_vpc_attachment_name":    name,
					"transit_router_attachment_description": "test",
					"auto_publish_route_enabled":            "true",
					"transit_router_vpc_attachment_options": map[string]interface{}{
						"\"ipv6Support\"": "enable",
					},
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":                                CHECKSET,
						"cen_id":                                CHECKSET,
						"zone_mappings.#":                       "2",
						"transit_router_id":                     CHECKSET,
						"transit_router_vpc_attachment_name":    name,
						"transit_router_attachment_description": "test",
						"auto_publish_route_enabled":            "true",
						"payment_type":                          "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultBY6Ody.id}",
							"zone_id":    "${alicloud_vswitch.defaultBY6Ody.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultoxj0Cs.id}",
							"zone_id":    "${alicloud_vswitch.defaultoxj0Cs.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Dd.id}",
							"zone_id":    "${alicloud_vswitch.defaulteKv3Dd.zone_id}",
						},
					},
					"transit_router_vpc_attachment_name":    name + "_update",
					"transit_router_attachment_description": "testupdate",
					"auto_publish_route_enabled":            "false",
					"transit_router_vpc_attachment_options": map[string]interface{}{
						"\"ipv6Support\"": "disable",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#":                       "3",
						"transit_router_vpc_attachment_name":    name + "_update",
						"transit_router_attachment_description": "testupdate",
						"auto_publish_route_enabled":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultoxj0Cs.id}",
							"zone_id":    "${alicloud_vswitch.defaultoxj0Cs.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultBY6Ody.id}",
							"zone_id":    "${alicloud_vswitch.defaultBY6Ody.zone_id}",
						},
					},
					"transit_router_vpc_attachment_name": name + "_update",
					"transit_router_vpc_attachment_options": map[string]interface{}{
						"\"ipv6Support\"": "enable",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#":                    "2",
						"transit_router_vpc_attachment_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cen_id", "dry_run"},
			},
		},
	})
}

var AliCloudCenTransitRouterVpcAttachmentMap7866 = map[string]string{
	"status":       CHECKSET,
	"create_time":  CHECKSET,
	"vpc_owner_id": CHECKSET,
	"payment_type": CHECKSET,
	"region_id":    CHECKSET,
}

func AliCloudCenTransitRouterVpcAttachmentBasicDependence7866(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultJLRlxW" {
  cidr_block  = "192.168.0.0/16"
  ipv6_isp    = "BGP"
  description = "ttt"
  enable_ipv6 = true
  vpc_name    = "ttt"
}

resource "alicloud_vswitch" "defaulteKv3Dd" {
  vpc_id               = alicloud_vpc.defaultJLRlxW.id
  cidr_block           = "192.168.3.0/24"
  zone_id              = "cn-hangzhou-h"
  vswitch_name         = "v1"
  ipv6_cidr_block_mask = "3"
}

resource "alicloud_vswitch" "defaultoxj0Cs" {
  vpc_id               = alicloud_vpc.defaultJLRlxW.id
  zone_id              = "cn-hangzhou-i"
  cidr_block           = "192.168.4.0/24"
  vswitch_name         = "v2"
  ipv6_cidr_block_mask = "4"
}

resource "alicloud_vswitch" "defaultBY6Ody" {
  vpc_id               = alicloud_vpc.defaultJLRlxW.id
  zone_id              = "cn-hangzhou-j"
  cidr_block           = "192.168.6.0/24"
  vswitch_name         = "v3"
  ipv6_cidr_block_mask = "6"
}

resource "alicloud_cen_instance" "defaultJ6HrUE" {
  cen_instance_name = "rdktest01"
}

resource "alicloud_cen_transit_router" "defaults5WvfD" {
  cen_id = alicloud_cen_instance.defaultJ6HrUE.id
}


`, name)
}

// Case 全生命周期_副本_xiaohu_线上 7248
func TestAccAliCloudCenTransitRouterVpcAttachment_basic7248(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap7248)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterVpcAttachmentBasicDependence7248)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id": "${alicloud_vpc.defaultJLRlxW.id}",
					"cen_id": "${alicloud_cen_instance.defaultJ6HrUE.id}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Da.id}",
							"zone_id":    "cn-hangzhou-j",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultoxj0Cs.id}",
							"zone_id":    "cn-hangzhou-i",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Dd.id}",
							"zone_id":    "cn-hangzhou-h",
						},
					},
					"transit_router_id":                     "${alicloud_cen_transit_router.defaults5WvfD.transit_router_id}",
					"transit_router_vpc_attachment_name":    name,
					"transit_router_attachment_description": "test",
					"auto_publish_route_enabled":            "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":                                CHECKSET,
						"cen_id":                                CHECKSET,
						"zone_mappings.#":                       "3",
						"transit_router_id":                     CHECKSET,
						"transit_router_vpc_attachment_name":    name,
						"transit_router_attachment_description": "test",
						"auto_publish_route_enabled":            "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultoxj0Cs.id}",
							"zone_id":    "cn-hangzhou-i",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Dd.id}",
							"zone_id":    "cn-hangzhou-h",
						},
					},
					"transit_router_vpc_attachment_name":    name + "_update",
					"transit_router_attachment_description": "testupdate",
					"auto_publish_route_enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#":                       "2",
						"transit_router_vpc_attachment_name":    name + "_update",
						"transit_router_attachment_description": "testupdate",
						"auto_publish_route_enabled":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Dd.id}",
							"zone_id":    "cn-hangzhou-h",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Da.id}",
							"zone_id":    "cn-hangzhou-j",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Dd.id}",
							"zone_id":    "cn-hangzhou-h",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaulteKv3Da.id}",
							"zone_id":    "cn-hangzhou-j",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cen_id", "dry_run"},
			},
		},
	})
}

var AliCloudCenTransitRouterVpcAttachmentMap7248 = map[string]string{
	"transit_router_attachment_id": CHECKSET,
	"status":                       CHECKSET,
	"create_time":                  CHECKSET,
	"region_id":                    CHECKSET,
}

func AliCloudCenTransitRouterVpcAttachmentBasicDependence7248(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultJLRlxW" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaulteKv3Dd" {
  vpc_id     = alicloud_vpc.defaultJLRlxW.id
  cidr_block = "192.168.3.0/24"
  zone_id    = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "defaultoxj0Cs" {
  vpc_id     = alicloud_vpc.defaultJLRlxW.id
  zone_id    = "cn-hangzhou-i"
  cidr_block = "192.168.4.0/24"
}

resource "alicloud_vswitch" "defaulteKv3Da" {
  vpc_id     = alicloud_vpc.defaultJLRlxW.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "192.168.5.0/24"
}

resource "alicloud_cen_instance" "defaultJ6HrUE" {
  cen_instance_name = "rdktest01"
}

resource "alicloud_cen_transit_router" "defaults5WvfD" {
  cen_id = alicloud_cen_instance.defaultJ6HrUE.id
}


`, name)
}

// Test Cen TransitRouterVpcAttachment. <<< Resource test cases, automatically generated.
