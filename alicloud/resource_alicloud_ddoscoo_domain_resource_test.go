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

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/alibabacloud-go/tea-rpc/client"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDdosCooDomainResource_https_ext(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_domain_resource.default"
	ra := resourceAttrInit(resourceId, AliCloudDdoscooDomainResourceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosCooServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosCooDomainResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandString(10)
	name := fmt.Sprintf("tf-testacc%s.alibaba.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdoscooDomainResourceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain":       name,
					"instance_ids": []string{"${data.alicloud_ddoscoo_instances.default.ids.0}"},
					"real_servers": []string{"177.167.32.11", "177.167.32.12", "177.167.32.13"},
					"rs_type":      `0`,
					"https_ext":    `{\"Http2\":1,\"Http2https\":0,\"Https2http\":0}`,
					"proxy_types": []map[string]interface{}{
						{
							"proxy_ports": []string{"80", "8080"},
							"proxy_type":  "http",
						},
						{
							"proxy_ports": []string{"443", "8443"},
							"proxy_type":  "https",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":         name,
						"instance_ids.#": "1",
						"real_servers.#": "3",
						"rs_type":        "0",
						"https_ext":      CHECKSET,
						"proxy_types.#":  "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_ext": `{\"Http2\":1,\"Http2https\":0,\"Https2http\":0}`,
					"proxy_types": []map[string]interface{}{
						{
							"proxy_ports": []string{"443"},
							"proxy_type":  "https",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_ext":     CHECKSET,
						"proxy_types.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_ids": []string{"${data.alicloud_ddoscoo_instances.default.ids.0}", "${data.alicloud_ddoscoo_instances.default.ids.1}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{"aliyun.com", "taobao.com", "alibaba.com"},
					"rs_type":      `1`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "3",
						"rs_type":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_ids": []string{"${data.alicloud_ddoscoo_instances.default.ids.0}"},
					"real_servers": []string{"177.167.32.11", "177.167.32.12", "177.167.32.13", "177.167.32.14", "177.167.32.15"},
					"rs_type":      `0`,
					"https_ext":    `{\"Http2\":0,\"Http2https\":0,\"Https2http\":0}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "1",
						"real_servers.#": "5",
						"rs_type":        "0",
						"https_ext":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ocsp_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ocsp_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ocsp_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ocsp_enabled": "false",
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

func TestAccAliCloudDdosCooDomainResource_none_https_ext(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_domain_resource.default"
	ra := resourceAttrInit(resourceId, AliCloudDdoscooDomainResourceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosCooServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosCooDomainResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%d.alibaba.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdoscooDomainResourceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain":         name,
					"instance_ids":   []string{"${data.alicloud_ddoscoo_instances.default.ids.0}"},
					"real_servers":   []string{"177.167.32.11", "177.167.32.12", "177.167.32.13"},
					"rs_type":        `0`,
					"ocsp_enabled":   "true",
					"custom_headers": "{\\\"22\\\":\\\"$ReqClientIP\\\",\\\"77\\\":\\\"88\\\",\\\"99\\\":\\\"$ReqClientPort\\\"}",
					"proxy_types": []map[string]interface{}{
						{
							"proxy_ports": []string{"80", "8080"},
							"proxy_type":  "http",
						},
						{
							"proxy_ports": []string{"443", "8443"},
							"proxy_type":  "https",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":         name,
						"instance_ids.#": "1",
						"real_servers.#": "3",
						"rs_type":        "0",
						"ocsp_enabled":   "true",
						"custom_headers": CHECKSET,
						"proxy_types.#":  "2",
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

var AliCloudDdoscooDomainResourceMap0 = map[string]string{
	"https_ext":      CHECKSET,
	"cname":          CHECKSET,
	"instance_ids.#": "1",
	"proxy_types.#":  "1",
	"real_servers.#": "1",
	"rs_type":        "0",
}

func AliCloudDdoscooDomainResourceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_ddoscoo_instances" "default" {
	}
`)
}

func TestUnitAliCloudDdoscooDomainResource(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	checkoutSupportedRegions(t, true, connectivity.DdoscooSupportedRegions)
	dInit, _ := schema.InternalMap(p["alicloud_ddoscoo_domain_resource"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ddoscoo_domain_resource"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"domain":    "CreateDomainResourceValue",
		"https_ext": "CreateDomainResourceValue",
		"proxy_types": []map[string]interface{}{
			{
				"proxy_ports": []int{443},
				"proxy_type":  "https",
			},
		},
		"instance_ids": []string{"CreateDomainResourceValue"},
		"real_servers": []string{"CreateDomainResourceValue"},
		"rs_type":      1,
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
		// DescribeDomainResource
		"WebRules": []interface{}{
			map[string]interface{}{
				"Domain":   "CreateDomainResourceValue",
				"HttpsExt": "CreateDomainResourceValue",
				"InstanceIds": []interface{}{
					"CreateDomainResourceValue",
				},
				"ProxyTypes": []interface{}{
					map[string]interface{}{
						"ProxyPorts": []int{443},
						"ProxyType":  "https",
					},
				},
				"RealServers": "CreateDomainResourceValue",
				"RsType":      1,
			},
		},
	}
	CreateMockResponse := map[string]interface{}{}
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ddoscoo_domain_resource", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdoscooClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudDdosCooDomainResourceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDomainResource" {
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
		err := resourceAliCloudDdosCooDomainResourceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ddoscoo_domain_resource"].Schema).Data(dInit.State(), nil)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudDdosCooDomainResourceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	//ModifyDomainResource
	attributesDiff := map[string]interface{}{
		"proxy_types": []map[string]interface{}{
			{
				"proxy_ports": []int{80},
				"proxy_type":  "http",
			},
		},
		"real_servers": []string{"ModifyDomainResourceValue"},
		"rs_type":      2,
		"https_ext":    "ModifyDomainResourceValue",
		"instance_ids": []string{"ModifyDomainResourceValue"},
	}
	diff, err := newInstanceDiff("alicloud_ddoscoo_domain_resource", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ddoscoo_domain_resource"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDomainResource Response
		"WebRules": []interface{}{
			map[string]interface{}{
				"HttpsExt": "ModifyDomainResourceValue",
				"ProxyTypes": []interface{}{
					map[string]interface{}{
						"ProxyPorts": []int{80},
						"ProxyType":  "http",
					},
				},
				"InstanceIds": []interface{}{
					"ModifyDomainResourceValue",
				},
				"RealServers": "ModifyDomainResourceValue",
				"RsType":      2,
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDomainResource" {
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
		err := resourceAliCloudDdosCooDomainResourceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ddoscoo_domain_resource"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDomainResource" {
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
		err := resourceAliCloudDdosCooDomainResourceRead(dExisted, rawClient)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudDdosCooDomainResourceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDomainResource" {
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
		err := resourceAliCloudDdosCooDomainResourceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test DdosCoo DomainResource. >>> Resource test cases, automatically generated.
// Case 测试选择证书-TF接入 7932
func TestAccAliCloudDdosCooDomainResource_basic7932(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_domain_resource.default"
	ra := resourceAttrInit(resourceId, AliCloudDdosCooDomainResourceMap7932)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosCooServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosCooDomainResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sddoscoodomainresource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdosCooDomainResourceBasicDependence7932)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rs_type":      "0",
					"ocsp_enabled": "false",
					"proxy_types": []map[string]interface{}{
						{
							"proxy_ports": []string{
								"80"},
							"proxy_type": "http",
						},
						{
							"proxy_ports": []string{
								"443"},
							"proxy_type": "https",
						},
						{
							"proxy_ports": []string{
								"80"},
							"proxy_type": "websocket",
						},
					},
					"real_servers": []string{
						"1.1.1.1", "2.2.2.2", "3.3.3.3"},
					"domain": "testld.qq.com",
					"instance_ids": []string{
						"${alicloud_ddoscoo_instance.defaultSJe7n8.id}", "${alicloud_ddoscoo_instance.default6lyurZ.id}", "${alicloud_ddoscoo_instance.defaultTTvY0D.id}"},
					"https_ext": "{\\\"Https2http\\\":1,\\\"Http2\\\":1,\\\"Http2https\\\":0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rs_type":        "0",
						"ocsp_enabled":   "false",
						"proxy_types.#":  "3",
						"real_servers.#": "3",
						"domain":         "testld.qq.com",
						"instance_ids.#": "3",
						"https_ext":      "{\"Https2http\":1,\"Http2\":1,\"Http2https\":0}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rs_type":      "1",
					"ocsp_enabled": "true",
					"proxy_types": []map[string]interface{}{
						{
							"proxy_ports": []string{
								"80", "8080"},
							"proxy_type": "http",
						},
					},
					"real_servers": []string{
						"1.qq.com"},
					"instance_ids": []string{
						"${alicloud_ddoscoo_instance.defaultTTvY0D.id}"},
					"https_ext":       "{\\\"Https2http\\\":0,\\\"Http2\\\":0,\\\"Http2https\\\":0}",
					"cert_identifier": "${local.certificate_id}",
					"cert_region":     "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rs_type":        "1",
						"ocsp_enabled":   "true",
						"proxy_types.#":  "1",
						"real_servers.#": "1",
						"instance_ids.#": "1",
						"https_ext":      "{\"Https2http\":0,\"Http2\":0,\"Http2https\":0}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cert_identifier": "${local.certificate_id_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert_identifier": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_headers": "{\\\"22\\\":\\\"$ReqClientIP\\\",\\\"77\\\":\\\"88\\\",\\\"99\\\":\\\"$ReqClientPort\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_headers": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert", "cert_identifier", "cert_region", "key"},
			},
		},
	})
}

var AliCloudDdosCooDomainResourceMap7932 = map[string]string{
	"cname": CHECKSET,
}

func AliCloudDdosCooDomainResourceBasicDependence7932(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ddoscoo_instance" "defaultTTvY0D" {
  normal_qps        = "3000"
  bandwidth_mode    = "2"
  product_type      = "ddoscoo"
  period            = "1"
  port_count        = "50"
  name              = "test"
  service_bandwidth = "200"
  base_bandwidth    = "30"
  bandwidth         = "50"
  function_version  = "0"
  address_type      = "Ipv4"
  edition_sale      = "coop"
  domain_count      = "50"
  product_plan      = "9"
}

resource "alicloud_ddoscoo_instance" "defaultSJe7n8" {
  normal_qps        = "3000"
  bandwidth_mode    = "2"
  product_type      = "ddoscoo"
  period            = "1"
  port_count        = "50"
  name              = "test2"
  service_bandwidth = "200"
  base_bandwidth    = "30"
  bandwidth         = "50"
  function_version  = "1"
  address_type      = "Ipv4"
  edition_sale      = "coop"
  domain_count      = "50"
  product_plan      = "9"
}

resource "alicloud_ddoscoo_instance" "default6lyurZ" {
  normal_qps        = "3000"
  bandwidth_mode    = "2"
  product_type      = "ddoscoo"
  period            = "1"
  port_count        = "50"
  name              = "test2"
  service_bandwidth = "200"
  base_bandwidth    = "30"
  bandwidth         = "50"
  function_version  = "1"
  address_type      = "Ipv4"
  edition_sale      = "coop"
  domain_count      = "50"
  product_plan      = "9"
}

resource "alicloud_ssl_certificates_service_certificate" "default3MYZEt" {
  cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID4TCCAsmgAwIBAgIRANbwhGnf1Ev/iAkbZUMiw8kwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjUxMDE2MDgzMTQ3WhcNMzAxMDE1MDgzMTQ3WjAlMQswCQYDVQQGEwJDTjEW
MBQGA1UEAxMNdGVzdGxkLnFxLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBANrqclqIRWI4U8nEspfEER88636gd8TkcIrUFjO2K5fhfhTxN8jScEfS
byycP/WfiL2KSsbVv/HzzqOXlzrW872DxIxCicMcv7Iswj2Cy38NDZyZAmR21vp3
XtL82+jCAQtz+IX85p4l3kkaZU+JWJqEY0q7E5M9zXeY99j/Zuni4f9elX312Ixf
iERl+Mcj1CjUu4T6ub+SoTZjblk6cZilIRqZCFCTYReuD9zI8F4b/87OUXyWOGCg
/fyzJXpK6iW2+2qJJuUpz9eUOy1OJDsQczZPiNm1Z4DBBlWVi2DrE01RBAGx5fWk
duXYIflPsA4mGHIEEduTZHIw/URJH4ECAwEAAaOB0jCBzzAOBgNVHQ8BAf8EBAMC
BaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFCiB
JgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEFBQcwAYYV
aHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8vY2EubXlz
c2wuY29tL215c3NsdGVzdHJzYS5jcnQwGAYDVR0RBBEwD4INdGVzdGxkLnFxLmNv
bTANBgkqhkiG9w0BAQsFAAOCAQEAXUcojuc/mE7SxaFrnaQ+KrLZLG1pg8Tlfko3
KrCZS6ebaviZeVsXTl0j2IRX7bMH9Y8zw5aEVOiZp6IWsBDLkmKcERFHIxpMMkVs
XVj3c1OaLYxYx0d+A385bHaBKJIx7qcyNTzI4BAb0RJs7JOhDg6f4nH42bFnvrYS
6WKShoYvZAqaMN+jfsWp+WIZd3wb7WNcdaRuJ5zWwTCIc3N51s+zVqfn3XrXwRFz
hhDRZIMLJM2dw3p0jSM+GdaIUd8gXW/jGJ55N3ZcHReDsgVGhvr+l+7O8oEtdEQ8
BEw9VgiKxQBBc5HmxcGhQbXLUAZ7TWxBnroMewH8oz0bL+s/IQ==
-----END CERTIFICATE-----
EOF
  certificate_name = format("%%s3", var.name)
  key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA2upyWohFYjhTycSyl8QRHzzrfqB3xORwitQWM7Yrl+F+FPE3
yNJwR9JvLJw/9Z+IvYpKxtW/8fPOo5eXOtbzvYPEjEKJwxy/sizCPYLLfw0NnJkC
ZHbW+nde0vzb6MIBC3P4hfzmniXeSRplT4lYmoRjSrsTkz3Nd5j32P9m6eLh/16V
ffXYjF+IRGX4xyPUKNS7hPq5v5KhNmNuWTpxmKUhGpkIUJNhF64P3MjwXhv/zs5R
fJY4YKD9/LMlekrqJbb7aokm5SnP15Q7LU4kOxBzNk+I2bVngMEGVZWLYOsTTVEE
AbHl9aR25dgh+U+wDiYYcgQR25NkcjD9REkfgQIDAQABAoIBABSoiCcL8gQ9TYbe
VX4l5fm9LXnzGapOZmJrdjbmC4IXKOdABiQg27CjZpNeoViD+Aru6HSQCj+CYu8k
KITIcRLiwuL7inWLmnltaN0WIS75o92xwLyLTGkxZ5TggL4bxK54gKzgO0EUUMA0
Sfgx/VcDhD0ynzvHWsLdABKNs4ABmOFWckCQGTNa0i85HG3SUQgKyDF0jVyvs4sU
uRR8jKOV9SNd2tr/YDkIyH9AFPvIpWvXBOE/k4QES3oQTQvWPaLyeeTHWBZLA2p6
cdePJF3furvaZMg4aE+oOuTuZh02o33jLpSMNNK6K0lfH0KCaA+qC9Hw8mD8y/2P
7t4G2kMCgYEA5s36XnrLhEyE55jynrumlaP7zAC3CXX65i0tMf1TUqmGl2+7it9e
G0iK4C/MLcSSov+S+808WymOqzArpIblygbNJq9oln3kJA7VvEhMnCNLv5cNMN89
eIKrmt67ux8f8TrZEyBd7Jl+ln+bI5IyCawIXn4bUSppHsOsLzaB6dcCgYEA8tA6
B+fknGvlaLljWMOzxBT3uQ2j6oTIgPLjuzOfDIyzlmyMDyePFo/TS32xD9idiJWC
tZAJ1sZIAu95kfU5pHLwK7ruW+PChtx5WtHXYAXh2++umoOjxK90VWpWeBVjXvtn
rKykfYh7JbDgBJClyo77iSZmCdu7O/XCu6KfBmcCgYAfejQNMp4S+wSdOWTNdTYw
7l5m4ioUZjzDq6GgUbZNbcVnXdusAu6oteoKzToBe++rv0NiiAkVPcOxYS5yj9tD
BE5yWjXfYGf+6u8HcKzSFpY8GPO5mJifmOKiioH78TDAC5CTZTSqEf0LtXeJEGU2
oHm7uWMsXKZdhb4z6jEpnQKBgC/S61sbRV+5sJmLyhF1mjaImrIMCbjrJkKflFMO
u8jQ/Z4nCv8BH6Gl+kvoGbOxSnXYXMI9+HIg45YQbLVew1ese7lhPAlFNs8xJYXJ
xs3W2sFi19T/EIZwuE0KgLVuIQBYK/dKmatP8lFeIQFFLCJVPx2oPni6mooYwZ4L
TZ8JAoGAJcWGPMmd3Ctr/58xjudUzRVWLGHJ2yG7THPXfnhjj2mw740Ful18zTd0
4QmYo6jSZ/1tQPaqW2i44kEOcDjCdKP6nyPIQDGgI3itw+iULEoxaN6lixrhP06R
KrNQWmTBvbU2fF0hd4Gy5CyIFxVxXoP7g0dV64BoK/FlBR3Yljk=
-----END RSA PRIVATE KEY-----
EOF
}

resource "alicloud_ssl_certificates_service_certificate" "update" {
  		certificate_name = "${var.name}-update"
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID4TCCAsmgAwIBAgIRALqrfbGoV0N+k/Y1LE+G8aowDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjUxMDE2MDgzMjIzWhcNMzAxMDE1MDgzMjIzWjAlMQswCQYDVQQGEwJDTjEW
MBQGA1UEAxMNdGVzdGxkLnFxLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAKrt0fPFAwUpwFv+o4MhuOR2RMhA46Imua/haTbzoMubjqEL+8qlgGIZ
KPDjJvNk1ALmJjpeXU1DN9bkFgcghh+Q1U5f4/aHOo5rUQM3KS45teMNji9R8epY
HIVCRma/teF8Cs2BsE2bY7Hh0bfl+Mmwwt6gzMgBYNtax6Srr+fiKfpy8E1C4rSk
A+Mqhv6H2qvLZ2ELcJoTvtbmnOIMF/xyBpuduTQ6TK8LQC35HhI2VSSL5k7erEFU
C+1JUTEGgZ18i8sjqeXPF0dln+LsVQwz3UlltevsRvkw1OkEtJeVLyu5cDHG9COZ
7xG2fR8BCraVCaeOjrseMo4LALTXzwMCAwEAAaOB0jCBzzAOBgNVHQ8BAf8EBAMC
BaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFCiB
JgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEFBQcwAYYV
aHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8vY2EubXlz
c2wuY29tL215c3NsdGVzdHJzYS5jcnQwGAYDVR0RBBEwD4INdGVzdGxkLnFxLmNv
bTANBgkqhkiG9w0BAQsFAAOCAQEAC7IU6qY8afXWn3X5zmoOhbPUSmeUqifWCrgd
D5E5SGqVP9eWzMz0BzB76G38hOEuhaAXpL1SjfOAFgHzPGVZt1VTDMvg2ajpPaPX
4IGNXdZxt+yIFxoBesjYv1YvEj6PW0+af5FV4gql4LhxNJilptysc8gRcxCeYgT9
6PgowhXNUxnDxIqIxAg93MMRHLqWFTF4Cvx8U3eEm4eUFuPlhfxYFvggTMCTwd7m
2Z6lMu8SqLmF57JMK820naibIQIKaa12BJlZoz4jx6c0Nav92dOTCSOjPgdUM5Zb
UcYZ8K63OFiZ98M0f3SCltqBGCTSl8wo4RMJUBrSRTizkd8Swg==
-----END CERTIFICATE-----
EOF
  key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAqu3R88UDBSnAW/6jgyG45HZEyEDjoia5r+FpNvOgy5uOoQv7
yqWAYhko8OMm82TUAuYmOl5dTUM31uQWByCGH5DVTl/j9oc6jmtRAzcpLjm14w2O
L1Hx6lgchUJGZr+14XwKzYGwTZtjseHRt+X4ybDC3qDMyAFg21rHpKuv5+Ip+nLw
TULitKQD4yqG/ofaq8tnYQtwmhO+1uac4gwX/HIGm525NDpMrwtALfkeEjZVJIvm
Tt6sQVQL7UlRMQaBnXyLyyOp5c8XR2Wf4uxVDDPdSWW16+xG+TDU6QS0l5UvK7lw
Mcb0I5nvEbZ9HwEKtpUJp46Oux4yjgsAtNfPAwIDAQABAoIBAEIrZ6QZR/aHN67F
UFJSyysyN6VYLWcX27lhJyR9QumfUiM5KuPDlwQi3k0Geo0tor9ujiz5W+Atnd/E
E1z188YjgNfi2jKVHg+FLuryPzBkaeu4Uysxa1e/fWb/BZcALy5XoSz2QCSC+6Cg
nVm2Hs4hbgbWNABXPEIejfvK9QFsU4aIrVEl5S4/UeVPBYtOo0yfDmQK657frGb4
RqB40dCeOybsDDi/sDKPLI4OO7h2n1ywTkMtXwD260VIBnQ4slje8T0Y8u2YrajO
12roKLLh/VDva7n0cNx5HQmX8JPnEFsim6yNFij9SaV1ylmo+0EViYMUz2FxDs1A
4nl8jaECgYEAxiuaETJabJaHBrRZYrs2iaygVAPl+gaSU2RoyJu7PKbxtLD+0jnr
eL1hR923T9adRsh+d5Az0gkZGGudQiNe+GI4x8PD3BrcexlBonS3CNqtrrb9SKzd
E1UK7L59TgyLJcDzHtEsPcTUe5LzFCfO3/O49HvTr3jXSLuwg4/3QSECgYEA3M8k
jSwex7KS/EIgOiQLpAS4L4AGopulRLGBMDdoYTmUmqg1VZyS0WYRm6Fce2Ndr4Lb
0sbYKCE299fpiuqRggVslJ5OMNM7G4HACBasNSQxWCkWG23Y3xqtCH5uajPXUIHl
t111Yrl9duXwXCGMHT/XnJVmKHUkxRRFxp/Dd6MCgYBrVq6q5eVIr/gPT5yi99jA
lbp6B2qIFQspFFgVYRT38000nDJKWIkM6zdIH/Xszsh90Jd/16HaAIeRTKjvbA1C
6KDsw0LRc9M88h81CZciuqAc5I0o0kkk8YlrVnq0zeKI3oxRguc9xeF51czIfA94
CqGB+5hbkU663L7tZAt/QQKBgBFxR4jjWFccEyJcMuGE4WqGeOo/qcaElwyTHQpr
BhLQEp4Y9YWaxbpG3tM1bvHMSqVHqAfBb2fUH9x6MNepae8kcIxY6QJQXVXx7PJ2
oAnenwtAy59FESGmoM6P9jbre3G/oR7YAiLXVkLjLRaKC+Bvn5+d6aD+h/YNgOmM
y0sTAoGAWPl6xPMjPwZS3EFi3tv/Nuuyd15qakQydUudYSsKcNhwMJbJ6WXuWRu7
RZJSfpv+8AnRK53E9FSjbyitpRK7G/XFRnNB7VPtXl/XU2nMfzv4PSI8HunE9eTJ
fDW8gVoHSvmTLXZVB0RVEHcDM7CsXCyVIpG/y6g6VheyXELwGOI=
-----END RSA PRIVATE KEY-----
EOF
	}

locals {
  certificate_id        = join("", [alicloud_ssl_certificates_service_certificate.default3MYZEt.id, "-cn-hangzhou"])
  certificate_id_update = join("", [alicloud_ssl_certificates_service_certificate.update.id, "-cn-hangzhou"])
}
`, name)
}

// Case 手动上传证书_TF接入 7935
func TestAccAliCloudDdosCooDomainResource_basic7935(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_domain_resource.default"
	ra := resourceAttrInit(resourceId, AliCloudDdosCooDomainResourceMap7935)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosCooServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosCooDomainResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sddoscoodomainresource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdosCooDomainResourceBasicDependence7935)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-1"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rs_type": "0",
					"proxy_types": []map[string]interface{}{
						{
							"proxy_ports": []string{
								"80"},
							"proxy_type": "http",
						},
						{
							"proxy_ports": []string{
								"443"},
							"proxy_type": "https",
						},
					},
					"real_servers": []string{
						"1.1.1.1"},
					"domain": "testcert.qq.com",
					"instance_ids": []string{
						"${alicloud_ddoscoo_instance.defaultaRfzZ9.id}"},
					"https_ext":   "{\\\"Https2http\\\":1,\\\"Http2\\\":1,\\\"Http2https\\\":0}",
					"cert_name":   name,
					"cert":        "${var.cert}",
					"key":         "${var.key}",
					"cert_region": "ap-southeast-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rs_type":        "0",
						"proxy_types.#":  "2",
						"real_servers.#": "1",
						"domain":         "testcert.qq.com",
						"instance_ids.#": "1",
						"https_ext":      "{\"Https2http\":1,\"Http2\":1,\"Http2https\":0}",
						"cert_name":      CHECKSET,
						"cert":           CHECKSET,
						"key":            CHECKSET,
						"cert_region":    "ap-southeast-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_ext":   "{\\\"Https2http\\\":0,\\\"Http2\\\":0,\\\"Http2https\\\":0}",
					"cert_name":   name + "1",
					"cert":        "${var.cert_update}",
					"key":         "${var.key_update}",
					"cert_region": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_ext":   "{\"Https2http\":0,\"Http2\":0,\"Http2https\":0}",
						"cert_name":   CHECKSET,
						"cert":        CHECKSET,
						"key":         CHECKSET,
						"cert_region": "cn-hangzhou",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert", "cert_identifier", "cert_region", "key"},
			},
		},
	})
}

var AliCloudDdosCooDomainResourceMap7935 = map[string]string{
	"cname": CHECKSET,
}

func AliCloudDdosCooDomainResourceBasicDependence7935(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ddoscoo_instance" "defaultaRfzZ9" {
  normal_bandwidth = "100"
  normal_qps       = "500"
  bandwidth_mode   = "2"
  product_plan     = "3"
  product_type     = "ddosDip"
  period           = "1"
  port_count       = "5"
  name             = "测试手动上传证书"
  function_version = "0"
  domain_count     = "10"
}

	variable "cert" {
  		default = <<EOF
-----BEGIN CERTIFICATE-----
MIID5DCCAsygAwIBAgIQL+BHedSOQ9OaO6KUb8CxSjANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yNTEwMTYwOTUxNDRaFw0zMDEwMTUwOTUxNDRaMCcxCzAJBgNVBAYTAkNOMRgw
FgYDVQQDEw90ZXN0Y2VydC5xcS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQCe+TPJ8P4huN4akgJow9yiw+hMYa7cE4ocB3ibFLIqK0FGTkNWKK7j
gxUU7PxmBCG2auIwMEYw8CbLBZdoptQXkJMlaJp5eY2bNZYbBK4i8/c4OUbrw1vD
f/iHwJD2kJ7CXSWABM29vw1KvKb07GVgQtESdUNynNd+x+BvSZg48MwPvVni+jV3
GBRX1WBm1700/Yvx+rgKPOz0ptZWnXaDLZ2brswNrSic2EgcHvpkd3YQzMmF3Hj7
lF9GBFSBqsRLymitMDM2BWdn0afspIxRxSlNmUkHku170NUTxJTyQQIVW7CGR4po
PSMwc/j7+fX41Q0bEVcs2YNrtOhvc+R/AgMBAAGjgdQwgdEwDgYDVR0PAQH/BAQD
AgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQo
gSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYBBQUHMAGG
FWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDovL2NhLm15
c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MBoGA1UdEQQTMBGCD3Rlc3RjZXJ0LnFx
LmNvbTANBgkqhkiG9w0BAQsFAAOCAQEALW7I3aulCKCZYhpDYObKjWxf0vA34A8c
CcNWlBqsBsbEFrYc3AvBt9esTj7ifWMYPn/o/z1GmhxVJNZhYXhm/bMsfzYVpDZF
mP5J8kCgahO6kM3qY1l0mIklRxp7QKQheDYnUezD+EsHxLWReHVCdtWEa7MGP+BD
XsnKLkjPPmcZ71oQM5vw/Zt+3RNuUWUyQEoc99Ioy8OeANqL1akQsZZWErs6bahf
0ZyY+bTsSgFBMGsuMgTfS3AY6VH3n1TLuxUbz+Mzc03i2hURxLCh91MYBp57H/cl
sXiuIjHumMV4zIuAgeCPVtC99FS+A7kX+aljdedSDDKkD1Qhf2RfmQ==
-----END CERTIFICATE-----
EOF
	}

	variable "key" {
  		default = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEpAIBAAKCAQEAnvkzyfD+IbjeGpICaMPcosPoTGGu3BOKHAd4mxSyKitBRk5D
Viiu44MVFOz8ZgQhtmriMDBGMPAmywWXaKbUF5CTJWiaeXmNmzWWGwSuIvP3ODlG
68Nbw3/4h8CQ9pCewl0lgATNvb8NSrym9OxlYELREnVDcpzXfsfgb0mYOPDMD71Z
4vo1dxgUV9VgZte9NP2L8fq4Cjzs9KbWVp12gy2dm67MDa0onNhIHB76ZHd2EMzJ
hdx4+5RfRgRUgarES8porTAzNgVnZ9Gn7KSMUcUpTZlJB5Lte9DVE8SU8kECFVuw
hkeKaD0jMHP4+/n1+NUNGxFXLNmDa7Tob3PkfwIDAQABAoIBABBFBSGsiiMDC6lg
z7dPculY2bwrc8euhkAXxvuUHXxxqyoEq3aKDlqsGSsS6pfkN1u3H56lP3cVtNqq
n/7ECB5rRkFErdPdy/oaUgTs2zIqPAm4hENvsxipzv0nCORQaQxk7QgSNcjZYIUJ
aUrVVgteezguNnC95X4M+wauBlKur4nhDV1XwsJiaurK5afK8p8em8gWqdfK4zCB
YGk2M6KlmwOBcoF3o/i+jfWju1lG2UPWKm4yDcCRptev0l/HxZ5oHdG/5eZGZ7Ay
Ke0J552LpSV/Lq/bIvg4fH7qbvyUTiXUwkICIWUaj2/461Hn6WIQM4hrPGFSD8Qs
YqkIIVUCgYEAzUaa4ozJitqhV9wOLi1m6IqdPXpRJo2B7mDFd7UakwY/5GFHdbwP
8RyTCili2Rm1mANEF8GJCNZgFp1z8acAhSoPfxTx+61lj/jGD6F2u7wMLEeMXkqj
YmhWuVRxv7++gQQotbXbM8LQXHliG6ItmxSQY5Lq2AnZVynkujOpgpUCgYEAxkGW
+/Fwx1Vhs1JNhKdn0eiTOgB1hBXljxe3hyFZ46CbSxibdjFEdYqYDXyAdgMglHur
tbdbETxP9Ycvokxk5IHnowvISnIc1lKyySXC9q21K6wiBFQLMKRj7KhxyNYyBJvC
ZkKM16H/cst9JOQQ8eLAQNnlcLuCkIqOTWnpecMCgYEAszQJMPARXkvhAG+WXY+7
QBUKkkn/IDX3ESCgIxISgfm5u2mFVe34yNfWMc/RgI/mLS/kuQx20iU8O2H3fyX4
2UfPwXSKj9lfSaG3XpvpqJjQ07Mego6MNfO6ig6DQw9kgwMbew6or3ZKKgC5ukAJ
qlH4f0UaCcIHYAWtrTQ+rkUCgYEAm3BbZ2dSTAbmVgkmW+ZA4PPfUq9/c7MTS9CF
hT4h0vVeLE+7u7w+94VVV+WQdnZXOfOImi2LCgVmj5ORRkdtJzeunEglnjC/6U3n
fQvNQ0jIbdhEx235ZAbPjYI3zAYcKz7P+QsekAYkWSWwFZd2rZ9hqrbsTCnH4Xmw
voNWma8CgYBfl1MZWLhgSwqz2yjtfy9JcXXOFKDEnbVQasDC6F+CKnUTtAUyLpWW
+xGQ8qMzA4CH1dR8ZBlwlfVj9EzZrw9zV3FYZW1w+XbSKW49FM4gmQ7V3eOza5oy
bMBBiK28dmU0fQRR7mXeiR2vuridlK1R01c6p9WFtZvXtVPy633BkQ==
-----END PRIVATE KEY-----
EOF
	}

	variable "cert_update" {
  		default = <<EOF
-----BEGIN CERTIFICATE-----
MIID5TCCAs2gAwIBAgIRAII7yjVVTE4dsnN/hCbI4Z0wDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjUxMDE2MDk1MzM5WhcNMzAxMDE1MDk1MzM5WjAnMQswCQYDVQQGEwJDTjEY
MBYGA1UEAxMPdGVzdGNlcnQucXEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAu+zZRv6YrSNBnQrcb67bghbvwgERbkJwCryIGCmC0oKNn+B657hk
tNopvnIlHHPBsyybvsLE7mxoycAo7fltvsZdY+OehMQ8EzfFRh7oM4dkIvh7fPwo
Mgiu8OrlN+JmhdhZO4iOs/f+zqENiZ6gEp8FNofvsz5sq8KaycY64OiHK5HmZ3xx
azxjom8ckCUBSogj5AVaLDj/SMnQiAVNwEz7OCsHYwIe540CMy5mtGTTgbcvly65
6c7t50+lyjxeXY2PSXzwBirt3O5uaqnvGHevKbzj1426sQhPk1gOguxaVQSwGcZR
xB21FWKTsF+xSu1W/0c2ZqtKFjJpj8P+BwIDAQABo4HUMIHRMA4GA1UdDwEB/wQE
AwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHwYDVR0jBBgwFoAU
KIEmBdE0Gj/Bcw+7k88VHD8Dv38wYwYIKwYBBQUHAQEEVzBVMCEGCCsGAQUFBzAB
hhVodHRwOi8vb2NzcC5teXNzbC5jb20wMAYIKwYBBQUHMAKGJGh0dHA6Ly9jYS5t
eXNzbC5jb20vbXlzc2x0ZXN0cnNhLmNydDAaBgNVHREEEzARgg90ZXN0Y2VydC5x
cS5jb20wDQYJKoZIhvcNAQELBQADggEBAJZdmNX3xuFEiYC4ad5d72nBDaBvbjtX
fPolfEiDF47lHFH+eQpdo4C5YR2dUQY2o3hP5LWrjEPmLHcZFTkobwyycTCwLKTs
8wkGhOYVcLcCBsg9hzU6t6GDTRSueglZlmRN8XOD3qG2AJefcJdB+idONw/1qheT
T3kxaHvrqaS7vPXNNe3U4u9cDfeEyEguYbhWdXyALBbQ12YF9/DJr4LHQNMKPqI1
Q68KJ/fT59zRPLmXpj57dgZylNPGo9FOis0XNISJ7j9u2baXb80afKZNzo56+Flf
aqyCp5f8p3Kmw6YQsgLJH/Ir+tBGUsJfaAlxGGoTLCf+Zqos2q/BvTI=
-----END CERTIFICATE-----
EOF
	}

	variable "key_update" {
  		default = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEpAIBAAKCAQEAu+zZRv6YrSNBnQrcb67bghbvwgERbkJwCryIGCmC0oKNn+B6
57hktNopvnIlHHPBsyybvsLE7mxoycAo7fltvsZdY+OehMQ8EzfFRh7oM4dkIvh7
fPwoMgiu8OrlN+JmhdhZO4iOs/f+zqENiZ6gEp8FNofvsz5sq8KaycY64OiHK5Hm
Z3xxazxjom8ckCUBSogj5AVaLDj/SMnQiAVNwEz7OCsHYwIe540CMy5mtGTTgbcv
ly656c7t50+lyjxeXY2PSXzwBirt3O5uaqnvGHevKbzj1426sQhPk1gOguxaVQSw
GcZRxB21FWKTsF+xSu1W/0c2ZqtKFjJpj8P+BwIDAQABAoIBAAZGhADzaozJkxzX
6oGOQMVI18vONlNMw6oQHqlT5Yr7EhinKeOIDFDfwioabLPVB8Bgenj1zxa5Jwyp
rpQ30prey+qUhMwhM3Je1+ceDBoAaO8kBhen4f29vX3NEkd593t7vIsY5c8LtoYW
6blRQz4r8kQeaPo+2Okpb/rR9FBjVMgqnYl993wwqohyJesgVr7qNNQVvilSB3TZ
eFkJlTMyhK7ymbvIUm0IHi3FZxgk6xCg/+LHegtKKpBh7UmdoNf/plOYPNeUyv5B
mQcsXcvFK8IKxaEtOdzWNMI4UsP5x5p4IFbB/9hndwPCh/6IcQu6AIu/hm6GKMtS
SgTBz4UCgYEAz2Dlnw7kBs7uwSeRCFhjJf0XyAUq+xr7Dy3v0AeTDClmmy6z9UNF
Dq+L7i8zH07yzVDzeWnTSRnpQLP1lnapdYfTeRWcKFz5ZDXwAcLpfCSyYHLgPDgV
8UKnS6wXx0cxPLZhmeaCrWXE8xxmiTWb8FqzQhQpwGckHBEMwJz5HysCgYEA5/xX
RcEct8u4DlFg7nYTK1efD9xyWgevqe8ciad3Mysy5AwfUXrMEf2R7V5G3EzSGEHz
yWMQmOdEgIKpd0/1HWZ7/k+m8XE0gwmX3M19EzEvwGOhcf7vvX4Lq7Z2f3afYiNR
hN9LhpwSwfJAyuVlHJGWoriwxSBFOiEYKwTxjpUCgYEAxlw552XX6UdAitNs392j
oO+xMqr2zM+m+4MGEydbmVN0iNUoX15kDMMPhtnw/W6Hwqo+6ZC3AAJf4XsBW1XP
i9NLDVQFVXpxNlB9bUHiEdQMJ0Nah19iZa5K1ZAcAopvZ1JQk2Qw9OkWdTBiR7ZW
nZY0Ru2AbkB6ArqwRwEfLZcCgYEA2uyJE5vdVRncVS65EfC9wF5NDnPUOmAch3rO
bJ1sYQ54VTuXZpZC9Qtd5irdJlMcxaWfwcJKTHGbdMdZ0+3R/G/VvbY/boSNsMeh
187YJP96981N8z1J04Ka0u47P6ibWsrHyGPNa3foP701JgR7eg1uoZs3vp/olKXc
n+RnbU0CgYB3GrAbyrwxkODiv5pPqMbRp4cdn7kbJatS5nAFQeyXkk2ilPbDvVlM
Oca4t9iJVF6kDepSHfabIJPxRnckpdLFD+3oXMsKp05SdFnulDadvV0xKDzznuWn
ESgznzdSAoMa+N48/IAoPVMkf8lV3g+Lby7d4m2XFR5vdu+vjqtq7g==
-----END PRIVATE KEY-----
EOF
	}
`, name)
}

// Test DdosCoo DomainResource. <<< Resource test cases, automatically generated.
