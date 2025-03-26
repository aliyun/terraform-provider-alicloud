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
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCenTransitRouterPeerAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_peer_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterPeerAttachmentMap)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterPeerAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterPeerAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
		},

		IDRefreshName:     resourceId,
		CheckDestroy:      testAccCheckCenTransitRouterPeerAttachmentDestroyWithProviders(&providers),
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":                       "alicloud.cn",
					"cen_id":                         "${alicloud_cen_instance.default.id}",
					"transit_router_id":              "${alicloud_cen_transit_router.default_0.transit_router_id}",
					"peer_transit_router_id":         "${alicloud_cen_transit_router.default_1.transit_router_id}",
					"peer_transit_router_region_id":  "cn-beijing",
					"transit_router_attachment_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"cen_id":                         CHECKSET,
						"peer_transit_router_id":         CHECKSET,
						"transit_router_id":              CHECKSET,
						"peer_transit_router_region_id":  "cn-beijing",
						"transit_router_attachment_name": name,
					}),
				),
			},
			// This step can not work in the multi region.
			//{
			//	ResourceName:            resourceId,
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"dry_run"},
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "tf-testaccdescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"transit_router_attachment_description": "tf-testaccdescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"transit_router_attachment_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_publish_route_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"auto_publish_route_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_bandwidth_package_id": "${alicloud_cen_bandwidth_package.default.id}",
					"bandwidth":                `2`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"cen_bandwidth_package_id": CHECKSET,
						"bandwidth":                "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_publish_route_enabled":            `true`,
					"bandwidth":                             `5`,
					"transit_router_attachment_description": "desp",
					"transit_router_attachment_name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"auto_publish_route_enabled":            "true",
						"bandwidth":                             "5",
						"transit_router_attachment_description": "desp",
						"transit_router_attachment_name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_link_type": "Gold",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_link_type": "Gold",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
						"tags.Updated": "TF",
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
						"tags.Updated": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCenTransitRouterPeerAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_peer_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterPeerAttachmentMap)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterPeerAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterPeerAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
		},

		IDRefreshName:     resourceId,
		CheckDestroy:      testAccCheckCenTransitRouterPeerAttachmentDestroyWithProviders(&providers),
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":                              "alicloud.cn",
					"cen_id":                                "${alicloud_cen_instance.default.id}",
					"transit_router_id":                     "${alicloud_cen_transit_router.default_0.transit_router_id}",
					"peer_transit_router_id":                "${alicloud_cen_transit_router.default_1.transit_router_id}",
					"peer_transit_router_region_id":         "cn-beijing",
					"transit_router_peer_attachment_name":   name,
					"auto_publish_route_enabled":            "false",
					"bandwidth":                             `5`,
					"cen_bandwidth_package_id":              "${alicloud_cen_bandwidth_package.default.id}",
					"dry_run":                               "false",
					"transit_router_attachment_description": "desp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"cen_id":                                CHECKSET,
						"peer_transit_router_id":                CHECKSET,
						"transit_router_id":                     CHECKSET,
						"peer_transit_router_region_id":         "cn-beijing",
						"transit_router_peer_attachment_name":   name,
						"auto_publish_route_enabled":            "false",
						"bandwidth":                             `5`,
						"cen_bandwidth_package_id":              CHECKSET,
						"dry_run":                               "false",
						"transit_router_attachment_description": "desp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_peer_attachment_name": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"transit_router_peer_attachment_name": name + "_Update",
					}),
				),
			},
		},
	})
}

var AliCloudCenTransitRouterPeerAttachmentMap = map[string]string{
	"auto_publish_route_enabled":            CHECKSET,
	"bandwidth":                             CHECKSET,
	"cen_bandwidth_package_id":              "",
	"cen_id":                                CHECKSET,
	"dry_run":                               NOSET,
	"peer_transit_router_id":                CHECKSET,
	"peer_transit_router_region_id":         "cn-beijing",
	"resource_type":                         "TR",
	"status":                                "Attached",
	"transit_router_attachment_description": "",
	"transit_router_attachment_name":        CHECKSET,
	"transit_router_id":                     CHECKSET,
}

func testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_cen_transit_router_peer_attachment ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			cbnService := CbnService{client}

			resp, err := cbnService.DescribeCenTransitRouterPeerAttachment(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}
			res = resp
			return nil
		}
		return fmt.Errorf("alicloud_cen_transit_router_peer_attachment not found")
	}
}

func testAccCheckCenTransitRouterPeerAttachmentDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenTransitRouterPeerAttachmentDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenTransitRouterPeerAttachmentDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {

	client := provider.Meta().(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_transit_router_peer_attachment" {
			continue
		}
		resp, err := cbnService.DescribeCenTransitRouterPeerAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Transit Router Attachment still exist,  ID %s ", fmt.Sprint(resp["TransitRouterAttachmentId"]))
		}
	}

	return nil
}

func AliCloudCenTransitRouterPeerAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {	
	default = "%s"
}

provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "cn"
  region = "cn-hangzhou"
}

resource "alicloud_cen_instance" "default" {
  provider = alicloud.cn
  name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_bandwidth_package" "default" {
  provider = alicloud.cn
  bandwidth                  = 5
  cen_bandwidth_package_name = var.name
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  provider = alicloud.cn
  instance_id          = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
}

resource "alicloud_cen_transit_router" "default_0" {
  provider = alicloud.cn
  cen_id = alicloud_cen_bandwidth_package_attachment.default.instance_id
  transit_router_name = "${var.name}-00"
}

resource "alicloud_cen_transit_router" "default_1" {
  provider = alicloud.bj
  cen_id = alicloud_cen_transit_router.default_0.cen_id
  transit_router_name = "${var.name}-01"
}

`, name)
}

func AliCloudCenTransitRouterPeerAttachmentBasicDependenceMutile(name string) string {
	return fmt.Sprintf(`

variable "name" {	
	default = "%s"
}

provider "alicloud" {
  alias = "eu"
  region = "eu-central-1"
}

provider "alicloud" {
  alias = "cn"
  region = "cn-hangzhou"
}

resource "alicloud_cen_instance" "default" {
  provider = alicloud.eu
  name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_bandwidth_package" "default" {
  provider = alicloud.eu
  bandwidth                  = 5
  cen_bandwidth_package_name = var.name
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  provider = alicloud.eu
  instance_id          = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
}

resource "alicloud_cen_transit_router" "default_0" {
  provider = alicloud.eu
  cen_id = alicloud_cen_bandwidth_package_attachment.default.instance_id
  transit_router_name = "${var.name}-00"
}

resource "alicloud_cen_transit_router" "default_1" {
  provider = alicloud.cn
  cen_id = alicloud_cen_transit_router.default_0.cen_id
  transit_router_name = "${var.name}-01"
}

`, name)
}

func TestAccAliCloudCenTransitRouterPeerAttachment_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_peer_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterPeerAttachmentMap)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterPeerAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterPeerAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
		},

		IDRefreshName:     resourceId,
		CheckDestroy:      testAccCheckCenTransitRouterPeerAttachmentDestroyWithProviders(&providers),
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":                       "alicloud.cn",
					"resource_type":                  "TR",
					"cen_id":                         "${alicloud_cen_instance.default.id}",
					"transit_router_id":              "${alicloud_cen_transit_router.default_0.transit_router_id}",
					"peer_transit_router_id":         "${alicloud_cen_transit_router.default_1.transit_router_id}",
					"peer_transit_router_region_id":  "cn-beijing",
					"transit_router_attachment_name": name,
					"bandwidth":                      "5",
					"bandwidth_type":                 "DataTransfer",
					"default_link_type":              "Platinum",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"cen_id":                         CHECKSET,
						"peer_transit_router_id":         CHECKSET,
						"transit_router_id":              CHECKSET,
						"resource_type":                  "TR",
						"peer_transit_router_region_id":  "cn-beijing",
						"transit_router_attachment_name": name,
						"bandwidth_type":                 "DataTransfer",
						"default_link_type":              "Platinum",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_link_type": "Gold",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"default_link_type": "Gold",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_type":           "BandwidthPackage",
					"cen_bandwidth_package_id": "${alicloud_cen_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"bandwidth_type":           "BandwidthPackage",
						"cen_bandwidth_package_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestUnitAliCloudCenTransitRouterPeerAttachment(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_cen_transit_router_peer_attachment"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_cen_transit_router_peer_attachment"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"cen_id":                                "cen_id",
		"transit_router_id":                     "transit_router_id",
		"peer_transit_router_id":                "peer_transit_router_id",
		"peer_transit_router_region_id":         "cn-beijing",
		"transit_router_attachment_name":        "transit_router_attachment_name",
		"auto_publish_route_enabled":            false,
		"bandwidth":                             2,
		"bandwidth_type":                        "BandwidthPackage",
		"cen_bandwidth_package_id":              "cen_bandwidth_package_id",
		"dry_run":                               false,
		"resource_type":                         "TR",
		"transit_router_attachment_description": "transit_router_attachment_description",
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
				"TransitRouterAttachmentId":          "MockTransitRouterAttachmentId",
				"CenId":                              "cen_id",
				"AutoPublishRouteEnabled":            false,
				"Bandwidth":                          2,
				"BandwidthType":                      "BandwidthPackage",
				"CenBandwidthPackageId":              "cen_bandwidth_package_id",
				"PeerTransitRouterId":                "peer_transit_router_region_id",
				"PeerTransitRouterRegionId":          "peer_transit_router_region_id",
				"ResourceType":                       "TR",
				"TransitRouterAttachmentDescription": "transit_router_attachment_description",
				"TransitRouterAttachmentName":        "transit_router_attachment_name",
				"TransitRouterId":                    "transit_router_id",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cen_transit_router_peer_attachment", "MockTransitRouterAttachmentId"))
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentCreate(d, rawClient)
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentCreate(d, rawClient)
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("MockCenId", ":", "MockTransitRouterAttachmentId"))
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

		err := resourceAliCloudCenTransitRouterPeerAttachmentUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyCenAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"auto_publish_route_enabled", "bandwidth", "dry_run", "cen_bandwidth_package_id", "transit_router_attachment_description", "transit_router_attachment_name"} {
			switch p["alicloud_cen_transit_router_peer_attachment"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_cen_transit_router_peer_attachment"].Schema).Data(nil, diff)
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyCenAttributeAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"auto_publish_route_enabled", "bandwidth", "dry_run", "cen_bandwidth_package_id", "transit_router_attachment_description", "transit_router_attachment_name"} {
			switch p["alicloud_cen_transit_router_peer_attachment"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_cen_transit_router_peer_attachment"].Schema).Data(nil, diff)
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentUpdate(resourceData1, rawClient)
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentDelete(d, rawClient)
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentDelete(d, rawClient)
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
		patcheDescribeCenTransitRouterPeerAttachment := gomonkey.ApplyMethod(reflect.TypeOf(&CbnService{}), "DescribeCenTransitRouterPeerAttachment", func(*CbnService, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("ResourceNotfound")
		})
		err := resourceAliCloudCenTransitRouterPeerAttachmentDelete(d, rawClient)
		patches.Reset()
		patcheDescribeCenTransitRouterPeerAttachment.Reset()
		assert.Nil(t, err)
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("DeleteMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_cen_transit_router_peer_attachment"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudCenTransitRouterPeerAttachmentDelete(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeCenInstanceNotFound", func(t *testing.T) {
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeCenInstanceAbnormal", func(t *testing.T) {
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
		err := resourceAliCloudCenTransitRouterPeerAttachmentRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

func SkipTestAccAliCloudCenTransitRouterPeerAttachment_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_peer_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterPeerAttachmentMap)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterPeerAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterPeerAttachmentBasicDependenceMutile)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
		},

		IDRefreshName:     resourceId,
		CheckDestroy:      testAccCheckCenTransitRouterPeerAttachmentDestroyWithProviders(&providers),
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_type":                  "TR",
					"provider":                       "alicloud.eu",
					"cen_id":                         "${alicloud_cen_instance.default.id}",
					"transit_router_id":              "${alicloud_cen_transit_router.default_1.transit_router_id}",
					"peer_transit_router_id":         "${alicloud_cen_transit_router.default_0.transit_router_id}",
					"peer_transit_router_region_id":  "cn-hangzhou",
					"transit_router_attachment_name": name,
					"auto_publish_route_enabled":     "false",
					//"bandwidth":                      `5`,
					//"cen_bandwidth_package_id":              "${alicloud_cen_bandwidth_package.default.id}",
					"dry_run":                               "false",
					"transit_router_attachment_description": "desp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"resource_type":                  "TR",
						"cen_id":                         CHECKSET,
						"peer_transit_router_id":         CHECKSET,
						"transit_router_id":              CHECKSET,
						"peer_transit_router_region_id":  "cn-hangzhou",
						"transit_router_attachment_name": name,
						"auto_publish_route_enabled":     "false",
						//"bandwidth":                      `5`,
						//"cen_bandwidth_package_id":              CHECKSET,
						"dry_run":                               "false",
						"transit_router_attachment_description": "desp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_type": "DataTransfer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"bandwidth_type": "DataTransfer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_type": "BandwidthPackage",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"bandwidth_type": "BandwidthPackage",
					}),
				),
			},
		},
	})
}

// Test Cen TransitRouterPeerAttachment. >>> Resource test cases, automatically generated.
var AliCloudCenTransitRouterPeerAttachmentMap6600 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudCenTransitRouterPeerAttachmentBasicDependence6600(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cen_instance" "default3Y1UVp" {
  description       = "铂金链路支持Terraform"
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultcmzeKb" {
  cen_id                     = alicloud_cen_instance.default3Y1UVp.id
  transit_router_name        = var.name
  transit_router_description = "铂金链路支持Terraform"
}

resource "alicloud_cen_transit_router" "defaultlAxHgL" {
  cen_id                     = alicloud_cen_instance.default3Y1UVp.id
  transit_router_name        = var.name
  transit_router_description = "铂金链路支持Terraform"
}


`, name)
}

// Case 铂金链路支持Terraform 6600  twin
func SkipTestAccAliCloudCenTransitRouterPeerAttachment_basic6600_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_peer_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterPeerAttachmentMap6600)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterPeerAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitrouterpeerattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterPeerAttachmentBasicDependence6600)
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
					"transit_router_attachment_name":        name,
					"default_link_type":                     "Platinum",
					"bandwidth_type":                        "DataTransfer",
					"cen_id":                                "${alicloud_cen_instance.default3Y1UVp.id}",
					"transit_router_attachment_description": "铂金链路支持Terraform - Create",
					"peer_transit_router_region_id":         defaultRegionToTest,
					"transit_router_id":                     "${alicloud_cen_transit_router.defaultcmzeKb.id}",
					"peer_transit_router_id":                "${alicloud_cen_transit_router.defaultlAxHgL.id}",
					"bandwidth":                             "10",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_name":        name,
						"default_link_type":                     "Platinum",
						"bandwidth_type":                        "DataTransfer",
						"cen_id":                                CHECKSET,
						"transit_router_attachment_description": "铂金链路支持Terraform - Create",
						"peer_transit_router_region_id":         CHECKSET,
						"transit_router_id":                     CHECKSET,
						"peer_transit_router_id":                CHECKSET,
						"bandwidth":                             "10",
						"tags.%":                                "2",
						"tags.Created":                          "TF",
						"tags.For":                              "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_link_type": "Gold",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_link_type": "Gold",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
						"tags.Updated": "TF",
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
				ImportStateVerifyIgnore: []string{"resource_type"},
			},
		},
	})
}

// Test Cen TransitRouterPeerAttachment. <<< Resource test cases, automatically generated.
