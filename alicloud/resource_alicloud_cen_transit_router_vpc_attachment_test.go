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
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"cen_id":          CHECKSET,
						"vpc_id":          CHECKSET,
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	checkoutSupportedRegions(t, true, connectivity.CenTransitRouterVpcAttachmentSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterVpcAttachmentMap0)
	providerFactories, factoryProviders := cenCrossAccountProviderFactories()
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		client := cenCrossAccountClientByAK(*factoryProviders, sharedCENCrossAccountCreds.utAK)
		if client == nil {
			return &CenServiceV2{}
		}
		return &CenServiceV2{client}
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTrVpcAttX%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	p := Provider().ResourcesMap
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
		err := resourceAliCloudCenTransitRouterVpcAttachmentUpdate(resourceData1, rawClient)
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// Case VPC Attachment支持ApplianceMode+IPv6测试用例 12197
func TestAccAliCloudCenTransitRouterVpcAttachment_basic12197(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterVpcAttachmentMap12197)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterVpcAttachmentBasicDependence12197)
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
						"vpc_id":                                CHECKSET,
						"transit_router_vpc_attachment_name":    name,
						"cen_id":                                CHECKSET,
						"payment_type":                          "PayAsYouGo",
						"transit_router_attachment_description": "tf-test",
						"transit_router_id":                     CHECKSET,
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

var AlicloudCenTransitRouterVpcAttachmentMap12197 = map[string]string{
	"transit_router_attachment_id": CHECKSET,
	"status":                       CHECKSET,
	"create_time":                  CHECKSET,
	"region_id":                    CHECKSET,
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

// ---- Sub-resource sweepers (补齐) ----
// DetachCenChildInstance (the legacy path used by alicloud_cen_instance_attachment)
// returns Forbidden.TypeInvalid for enterprise-TR VPC children; the correct API
// to detach them is DeleteTransitRouterVpcAttachment so the parent
// alicloud_cen_instance sweeper can finally DeleteCen. Cross-region TR peers are
// detached with DeleteTransitRouterPeerAttachment. These sub-resource sweepers
// run BEFORE the parent (explicit ordered make sweep calls) and carry the same
// name-based NODELETING guard as the parent (protectedSweepVpcName) so shared
// infra CENs (e.g. "no-deleting-cen") and their children are never touched.

func init() {
	resource.AddTestSweepers("alicloud_cen_transit_router_vpc_attachment", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_vpc_attachment",
		F:    testSweepCenTransitRouterVpcAttachments,
	})
}

func testSweepCenTransitRouterVpcAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// Transit Router VPC attachments are regional to the VPC's region (e.g.
	// cn-shanghai), not the CEN's home region (cn-beijing). The LIST must query
	// each candidate region's endpoint; the DELETE works cross-region from the
	// home region with Force=true. Regional clients are cached so each is created
	// only once across all CENs.
	regionalClients := map[string]*connectivity.AliyunClient{}
	return forEachCenForSubSweep(client, func(cenId, cenName string) {
		sweepCenTransitRouterAttachments(client, cenId, cenName,
			"ListTransitRouterVpcAttachments", "DeleteTransitRouterVpcAttachment", "VPC",
			cenSubSweepRegions(client.RegionId), regionalClients)
	})
}

// cenSubSweepRegions returns the set of regions to probe for transit-router VPC
// attachments. VPC attachments live in the VPC's region, so the sweeper must
// query each region's endpoint. The CEN home region is always included; common
// test VPC regions are added so cross-region enterprise attachments (e.g.
// cn-shanghai VPCs under a cn-beijing CEN) are found.
func cenSubSweepRegions(homeRegion string) []string {
	return []string{
		homeRegion,
		"cn-shanghai",
		"cn-hangzhou",
		"cn-zhangjiakou",
		"cn-shenzhen",
		"cn-hongkong",
		"cn-huhehaote",
		"cn-chengdu",
		"cn-heyuan",
		"cn-guangzhou",
	}
}

// forEachCenForSubSweep enumerates CEN instances in the region
// forEachCenForSubSweep enumerates CEN instances in the region and invokes cb
// for each CEN eligible for sub-resource sweeping: CENs carrying a do-not-delete
// marker (protectedSweepVpcName) are always skipped (even in SWEEP_ALL mode),
// and when SWEEP_ALL is not active only CENs whose name matches a test prefix
// are visited. Mirrors the enumeration + guard in testSweepCenInstances.
func forEachCenForSubSweep(client *connectivity.AliyunClient, cb func(cenId, cenName string)) error {
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_example",
		"tf-example",
		"terraform-example",
	}
	describeCensRequest := cbn.CreateDescribeCensRequest()
	describeCensRequest.RegionId = client.RegionId
	describeCensRequest.PageSize = requests.NewInteger(PageSizeLarge)
	describeCensRequest.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCbnClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DescribeCens(describeCensRequest)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving CEN Instances: %s", err)
		}
		describeCensResponse, _ := raw.(*cbn.DescribeCensResponse)
		if len(describeCensResponse.Cens.Cen) < 1 {
			break
		}
		for _, cenInstance := range describeCensResponse.Cens.Cen {
			name := cenInstance.Name
			cenId := cenInstance.CenId
			if protectedSweepVpcName(name) {
				log.Printf("[INFO] Skipping CEN sub-resource sweep (do-not-delete marker): %s (%s)", name, cenId)
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
					log.Printf("[INFO] Skipping CEN sub-resource sweep (name does not match test prefixes): %s (%s)", name, cenId)
					continue
				}
			}
			cb(cenId, name)
		}
		if len(describeCensResponse.Cens.Cen) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(describeCensRequest.PageNumber)
		if err != nil {
			return err
		}
		describeCensRequest.PageNumber = page
	}
	return nil
}

// cenIdsProtectedFromSweep enumerates CEN instances in the region and returns
// the set of CEN IDs whose name matches a do-not-delete marker
// (protectedSweepVpcName). Sub-resource sweepers (flowlog, bandwidth package)
// use this to skip resources attached to a protected parent CEN even in
// SWEEP_ALL mode, so a shared infra CEN's children are never orphaned or
// deleted ahead of the protected parent.
func cenIdsProtectedFromSweep(client *connectivity.AliyunClient) (map[string]bool, error) {
	protected := map[string]bool{}
	describeCensRequest := cbn.CreateDescribeCensRequest()
	describeCensRequest.RegionId = client.RegionId
	describeCensRequest.PageSize = requests.NewInteger(PageSizeLarge)
	describeCensRequest.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCbnClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DescribeCens(describeCensRequest)
		})
		if err != nil {
			return protected, fmt.Errorf("error retrieving CEN Instances: %s", err)
		}
		describeCensResponse, _ := raw.(*cbn.DescribeCensResponse)
		if len(describeCensResponse.Cens.Cen) < 1 {
			break
		}
		for _, cenInstance := range describeCensResponse.Cens.Cen {
			if protectedSweepVpcName(cenInstance.Name) {
				protected[cenInstance.CenId] = true
			}
		}
		if len(describeCensResponse.Cens.Cen) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(describeCensRequest.PageNumber)
		if err != nil {
			return protected, err
		}
		describeCensRequest.PageNumber = page
	}
	return protected, nil
}

// sweepCenTransitRouterAttachments lists transit-router attachments of the
// given kind for a CEN and force-deletes each. listAction/deleteAction are the
// cbn RPC action names. regions is the set of regions to probe for the LIST:
// VPC attachments live in the VPC's region (e.g. cn-shanghai) while peer
// attachments live in the CEN's home region, so the caller passes a candidate
// list for VPC and a single-element list for peer. regionalClients caches
// per-region AliyunClients so each regional endpoint is created only once.
//
// Force=true is always set on the delete so route-table associations and other
// dependencies are removed automatically (the API otherwise rejects with
// OperationFailed.RouteTableAssociationExist). The delete is issued from the
// home-region client; cross-region force-delete works by attachment ID.
func sweepCenTransitRouterAttachments(client *connectivity.AliyunClient, cenId, cenName, listAction, deleteAction, kind string, regions []string, regionalClients map[string]*connectivity.AliyunClient) {
	for _, region := range regions {
		// Determine the client to use for LIST. VPC attachments are regional
		// to the VPC's region, so a separate regional client (with its own
		// endpoint) is required for each region. Peer attachments live in the
		// CEN's home region, so the home client suffices.
		listClient := client
		if region != client.RegionId {
			if rc, ok := regionalClients[region]; ok {
				listClient = rc
			} else {
				raw, err := sharedClientForRegion(region)
				if err != nil {
					log.Printf("[ERROR] sharedClientForRegion(%s) failed for CEN %s (%s): %v", region, cenName, cenId, err)
					continue
				}
				listClient = raw.(*connectivity.AliyunClient)
				regionalClients[region] = listClient
			}
		}
		request := map[string]interface{}{
			"CenId":      cenId,
			"RegionId":   region,
			"MaxResults": 100,
		}
		for {
			var response map[string]interface{}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				resp, e := listClient.RpcPost("Cbn", "2017-09-12", listAction, nil, request, true)
				response = resp
				if e != nil {
					// IllegalParam.CenIdOrRegionId means this region has no transit
					// router for this CEN — non-fatal, just skip this region.
					if IsExpectedErrors(e, []string{"IllegalParam.CenIdOrRegionId", "IllegalParam.Region"}) {
						return retry.NonRetryableError(e)
					}
					if NeedRetry(e) {
						wait()
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				addDebug(listAction, response, request)
				return nil
			})
			if err != nil {
				// Non-fatal: many regions will have no TR for this CEN.
				log.Printf("[DEBUG] %s in %s for CEN %s (%s): %v", listAction, region, cenName, cenId, err)
				break
			}
			att, jerr := jsonpath.Get("$.TransitRouterAttachments", response)
			if jerr != nil {
				log.Printf("[ERROR] jsonpath.Get $.TransitRouterAttachments failed for CEN %s (%s) in %s: %v", cenName, cenId, region, jerr)
				break
			}
			result, _ := att.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				attName := fmt.Sprint(item["TransitRouterAttachmentName"])
				id := fmt.Sprint(item["TransitRouterAttachmentId"])
				log.Printf("[INFO] Deleting CEN TR %s attachment: %s (%s) parent CEN %s (%s) region %s", kind, attName, id, cenName, cenId, region)
				delReq := map[string]interface{}{
					"TransitRouterAttachmentId": id,
					"Force":                     true,
				}
				var delResp map[string]interface{}
				delWait := incrementalWait(3*time.Second, 5*time.Second)
				delErr := retry.Retry(1*time.Minute, func() *retry.RetryError {
					resp, e := client.RpcPost("Cbn", "2017-09-12", deleteAction, nil, delReq, false)
					delResp = resp
					if e != nil {
						// NotFound means already deleted — treat as success.
						if IsExpectedErrors(e, []string{"InvalidTransitRouterAttachmentId.NotFound"}) {
							return nil
						}
						// NotSupportForceDelete is non-retryable.
						if IsExpectedErrors(e, []string{"OperationFailed.NotSupportForceDelete"}) {
							return retry.NonRetryableError(e)
						}
						if IsExpectedErrors(e, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status", "IncorrectStatus.Attachment"}) || NeedRetry(e) {
							delWait()
							return retry.RetryableError(e)
						}
						return retry.NonRetryableError(e)
					}
					return nil
				})
				addDebug(deleteAction, delResp, delReq)
				if delErr != nil {
					log.Printf("[ERROR] %s failed for attachment %s (%s): %v", deleteAction, attName, id, delErr)
				}
			}
			// NextToken pagination (MaxResults/NextToken, not PageSize/PageNumber).
			nextToken := ""
			if v, ok := response["NextToken"]; ok && v != nil {
				if s := fmt.Sprint(v); s != "" {
					nextToken = s
				}
			}
			if nextToken == "" {
				break
			}
			request["NextToken"] = nextToken
		}
	}
}
