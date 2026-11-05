package cen_test

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudCenTransitRouterVpcAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.CenTransitRouterVpcAttachmentSupportRegions)
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCenServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterVpcAttachment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCenTransitRouterVpcAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":       "${alicloud_cen_transit_router.default.cen_id}",
					"vpc_id":       "${alicloud_vpc.default.id}",
					"force_delete": "true",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default_master.id}",
							"zone_id":    "${alicloud_vswitch.default_master.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.default.id}",
							"zone_id":    "${alicloud_vswitch.default.zone_id}",
						},
					},
					"route_table_association_enabled": "false",
					"route_table_propagation_enabled": "false",
					"order_type":                      "PayByCenOwner",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":          acctest.CHECKSET,
						"vpc_id":          acctest.CHECKSET,
						"zone_mappings.#": "2",
						"order_type":      "PayByCenOwner",
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
							"vswitch_id": "${alicloud_vswitch.default_master.id}",
							"zone_id":    "${alicloud_vswitch.default_master.zone_id}",
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
							"vswitch_id": "${alicloud_vswitch.default.id}",
							"zone_id":    "${alicloud_vswitch.default.zone_id}",
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
				ImportStateVerifyIgnore: []string{"route_table_association_enabled", "route_table_propagation_enabled", "force_delete"},
			},
		},
	})
}

func TestAccAliCloudCenTransitRouterVpcAttachment_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.CenTransitRouterVpcAttachmentSupportRegions)
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCenServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterVpcAttachment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCenTransitRouterVpcAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                                "${alicloud_cen_transit_router.default.cen_id}",
					"vpc_id":                                "${alicloud_vpc.default.id}",
					"transit_router_id":                     "${alicloud_cen_transit_router.default.transit_router_id}",
					"resource_type":                         "VPC",
					"payment_type":                          "PayAsYouGo",
					"vpc_owner_id":                          "${data.alicloud_account.default.id}",
					"auto_publish_route_enabled":            "false",
					"transit_router_attachment_name":        name,
					"transit_router_attachment_description": name,
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default_master.id}",
							"zone_id":    "${alicloud_vswitch.default_master.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.default.id}",
							"zone_id":    "${alicloud_vswitch.default.zone_id}",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "TransitRouterVpcAttachment",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                                acctest.CHECKSET,
						"vpc_id":                                acctest.CHECKSET,
						"transit_router_id":                     acctest.CHECKSET,
						"resource_type":                         "VPC",
						"payment_type":                          "PayAsYouGo",
						"vpc_owner_id":                          acctest.CHECKSET,
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

// TestAccAliCloudCenTransitRouterVpcAttachment_crossAccountOrderType mirrors
// the cross-account payer flow for VPC: account A (TerraformUT) owns the CEN
// and creates the transit_router_vpc_attachment; account B (TerraformTest)
// owns the VPC and the grant_attachment that authorizes A. order_type on the
// attachment is expected to follow the grant, and updating it requires the
// grant to be updated first (enforced here via depends_on).
func TestAccAliCloudCenTransitRouterVpcAttachment_crossAccountOrderType(t *testing.T) {
	testAccPreCheckCENCrossAccount(t)
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.CenTransitRouterVpcAttachmentSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap0)
	providerFactories, factoryProviders := cenCrossAccountProviderFactories()
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		client := cenCrossAccountClientByAK(*factoryProviders, sharedCENCrossAccountCreds.utAK)
		if client == nil {
			return cen.NewCenServiceV2(nil)
		}
		return cen.NewCenServiceV2(client)
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTrVpcAttX%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		ProviderFactories: providerFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCENTransitRouterVpcAttachmentCrossAccountConfig(name, "PayByCenOwner"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type": "PayByCenOwner",
					}),
				),
			},
			{
				Config: testAccCENTransitRouterVpcAttachmentCrossAccountConfig(name, "PayByResourceOwner"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type": "PayByResourceOwner",
					}),
				),
			},
			// ImportState intentionally omitted (see VPN cross-account test).
		},
	})
}

func testAccCENTransitRouterVpcAttachmentCrossAccountConfig(name, orderType string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
%s

data "alicloud_account" "b" {}

data "alicloud_account" "a" {
  provider = alicloud.a
}

# --- Account B: VPC + VSwitch + Grant ---
data "alicloud_zones" "b" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.b.ids.0
}

resource "alicloud_cen_transit_router_grant_attachment" "default" {
  instance_type = "VPC"
  instance_id   = alicloud_vpc.default.id
  cen_owner_id  = data.alicloud_account.a.id
  cen_id        = alicloud_cen_instance.default.id
  order_type    = %q
}

# --- Account A: CEN + Transit Router + VPC attachment consuming B's grant ---
resource "alicloud_cen_instance" "default" {
  provider          = alicloud.a
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  provider = alicloud.a
  cen_id   = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  provider          = alicloud.a
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  vpc_owner_id      = data.alicloud_account.b.id
  force_delete      = true
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = data.alicloud_zones.b.ids.0
  }
  order_type = %q
  depends_on = [alicloud_cen_transit_router_grant_attachment.default]
}
`, name, cenCrossAccountProviderBlocks(), orderType, orderType)
}

var AliCloudCenTransitRouterVpcAttachmentMap0 = map[string]string{
	"transit_router_id":            acctest.CHECKSET,
	"resource_type":                acctest.CHECKSET,
	"payment_type":                 acctest.CHECKSET,
	"vpc_owner_id":                 acctest.CHECKSET,
	"transit_router_attachment_id": acctest.CHECKSET,
	"status":                       acctest.CHECKSET,
	"region_id":                    acctest.CHECKSET,
}

func AliCloudCenTransitRouterVpcAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_account" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = "cn-hangzhou-h"
  vswitch_name = var.name
}

resource "alicloud_vswitch" "default_master" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = "cn-hangzhou-i"
  vswitch_name = "${var.name}-master"
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

// lintignore: R001
func TestUnitAliCloudCenTransitRouterVpcAttachment(t *testing.T) {
	p := acctest.Provider.ResourcesMap
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_cen_transit_router_vpc_attachment", "MockTransitRouterAttachmentId"))
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
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentCreate(d, rawClient)
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
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentCreate(d, rawClient)
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
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("CenId", ":", "MockTransitRouterAttachmentId"))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateTransitRouterVpcAttachmentAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_type", "dry_run", "transit_router_attachment_description", "transit_router_attachment_name"} {
			switch p["alicloud_cen_transit_router_vpc_attachment"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
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
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateTransitRouterVpcAttachmentAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_type", "dry_run", "transit_router_attachment_description", "transit_router_attachment_name"} {
			switch p["alicloud_cen_transit_router_vpc_attachment"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
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
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentDelete(d, rawClient)
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
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentDelete(d, rawClient)
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
		patcheDescribeCenTransitRouterVpcAttachment := gomonkey.ApplyMethod(reflect.TypeOf(cen.NewCbnService(nil)), "DescribeCenTransitRouterVpcAttachment", func(*cen.CbnService, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("ResourceNotfound")
		})
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentDelete(d, rawClient)
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
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentRead(d, rawClient)
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
		err := cen.ResourceAliCloudCenTransitRouterVpcAttachmentRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Cen TransitRouterVpcAttachment. >>> Resource test cases, automatically generated.
// Case TR支持IPv6_善问_线上_副本1726040691220 7866
func TestAccAliCloudCenTransitRouterVpcAttachment_basic7866(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap7866)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCenServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitroutervpcattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCenTransitRouterVpcAttachmentBasicDependence7866)
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
						"vpc_id":                                acctest.CHECKSET,
						"cen_id":                                acctest.CHECKSET,
						"zone_mappings.#":                       "2",
						"transit_router_id":                     acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":       acctest.CHECKSET,
	"create_time":  acctest.CHECKSET,
	"vpc_owner_id": acctest.CHECKSET,
	"payment_type": acctest.CHECKSET,
	"region_id":    acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap7248)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCenServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCenTransitRouterVpcAttachmentBasicDependence7248)
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
						"vpc_id":                                acctest.CHECKSET,
						"cen_id":                                acctest.CHECKSET,
						"zone_mappings.#":                       "3",
						"transit_router_id":                     acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"transit_router_attachment_id": acctest.CHECKSET,
	"status":                       acctest.CHECKSET,
	"create_time":                  acctest.CHECKSET,
	"region_id":                    acctest.CHECKSET,
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

// Case VPC Attachment支持ApplianceMode+IPv6测试用例 12197
func TestAccAliCloudCenTransitRouterVpcAttachment_basic12197(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCenTransitRouterVpcAttachmentMap12197)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCenServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCenTransitRouterVpcAttachmentBasicDependence12197)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"options": []map[string]interface{}{
						{
							"ipv6_support":           "enable",
							"appliance_mode_support": "enable",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultWZIJKW.id}",
							"zone_id":    "${alicloud_vswitch.defaultWZIJKW.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.default5XiiOt.id}",
							"zone_id":    "${alicloud_vswitch.default5XiiOt.zone_id}",
						},
					},
					"auto_publish_route_enabled":            "true",
					"vpc_id":                                "${alicloud_vpc.defaultpzLsWV.id}",
					"transit_router_vpc_attachment_name":    name,
					"cen_id":                                "${alicloud_cen_instance.defaultv7lEPH.id}",
					"payment_type":                          "PayAsYouGo",
					"transit_router_attachment_description": "tf-test",
					"transit_router_id":                     "${alicloud_cen_transit_router.defaultrjuavE.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#":                       "2",
						"auto_publish_route_enabled":            "true",
						"vpc_id":                                acctest.CHECKSET,
						"transit_router_vpc_attachment_name":    name,
						"cen_id":                                acctest.CHECKSET,
						"payment_type":                          "PayAsYouGo",
						"transit_router_attachment_description": "tf-test",
						"transit_router_id":                     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": []map[string]interface{}{
						{
							"ipv6_support":           "disable",
							"appliance_mode_support": "enable",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultjQs0qh.id}",
							"zone_id":    "${alicloud_vswitch.defaultjQs0qh.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultWZIJKW.id}",
							"zone_id":    "${alicloud_vswitch.defaultWZIJKW.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.default5XiiOt.id}",
							"zone_id":    "${alicloud_vswitch.default5XiiOt.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": []map[string]interface{}{
						{
							"ipv6_support":           "disable",
							"appliance_mode_support": "disable",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_vpc_attachment_name":    name + "_update",
					"transit_router_attachment_description": "ss",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_vpc_attachment_name":    name + "_update",
						"transit_router_attachment_description": "ss",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultWZIJKW.id}",
							"zone_id":    "${alicloud_vswitch.defaultWZIJKW.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.default5XiiOt.id}",
							"zone_id":    "${alicloud_vswitch.default5XiiOt.zone_id}",
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default5XiiOt.id}",
							"zone_id":    "${alicloud_vswitch.default5XiiOt.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultjQs0qh.id}",
							"zone_id":    "${alicloud_vswitch.defaultjQs0qh.zone_id}",
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
					"auto_publish_route_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_publish_route_enabled": "false",
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AlicloudCenTransitRouterVpcAttachmentMap12197 = map[string]string{
	"transit_router_attachment_id": acctest.CHECKSET,
	"status":                       acctest.CHECKSET,
	"create_time":                  acctest.CHECKSET,
	"region_id":                    acctest.CHECKSET,
}

func AlicloudCenTransitRouterVpcAttachmentBasicDependence12197(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "defaultpzLsWV" {
  ipv6_isp    = "BGP"
  description = "jiatu-terraform-test"
  cidr_block  = "192.168.0.0/16"
  enable_ipv6 = true
  vpc_name    = "jiatu-terraform-test"
}

resource "alicloud_vswitch" "defaultWZIJKW" {
  ipv6_cidr_block_mask = "3"
  vpc_id               = alicloud_vpc.defaultpzLsWV.id
  zone_id              = "cn-hangzhou-h"
  cidr_block           = "192.168.3.0/24"
  vswitch_name         = "tf"
}

resource "alicloud_vswitch" "default5XiiOt" {
  vpc_id               = alicloud_vpc.defaultpzLsWV.id
  zone_id              = "cn-hangzhou-i"
  cidr_block           = "192.168.4.0/24"
  vswitch_name         = "tf2"
  ipv6_cidr_block_mask = "4"
}

resource "alicloud_vswitch" "defaultjQs0qh" {
  description          = "tf"
  ipv6_cidr_block_mask = "6"
  vpc_id               = alicloud_vpc.defaultpzLsWV.id
  zone_id              = "cn-hangzhou-j"
  cidr_block           = "192.168.6.0/24"
  vswitch_name         = "tf"
}

resource "alicloud_cen_instance" "defaultv7lEPH" {
  cen_instance_name = "tf-test"
}

resource "alicloud_cen_transit_router" "defaultrjuavE" {
  cen_id = alicloud_cen_instance.defaultv7lEPH.id
}


`, name)
}

// Test Cen TransitRouterVpcAttachment. <<< Resource test cases, automatically generated.
