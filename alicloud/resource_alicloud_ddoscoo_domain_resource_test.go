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
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdoscooDomainResource")
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
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdoscooDomainResource")
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
					"domain":       name,
					"instance_ids": []string{"${data.alicloud_ddoscoo_instances.default.ids.0}"},
					"real_servers": []string{"177.167.32.11", "177.167.32.12", "177.167.32.13"},
					"rs_type":      `0`,
					"ocsp_enabled": "true",
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
			testAccPreCheckWithTime(t, []int{1})
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
MIID4TCCAsmgAwIBAgIRANZGvLwT8kuWpPlZ/Aj+uPgwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjQwODIzMDk0NzA0WhcNMjUwODIzMDk0NzA0WjAlMQswCQYDVQQGEwJDTjEW
MBQGA1UEAxMNdGVzdGxkLnFxLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBALmZY2geTFi+50gAVyDQH9Y5sTv8LLX6+MET1l3larzjX1M0Az9ZEIc0
TNrAp8mtJRlpQCzyDPZg88AwSdEwqSOSsnGzfS2DUcPJmdn2a2n5PLvWE28qPuSf
6fl3IhNiPzLYR51+7ccJKEQRhfOK2usmJo6oTG/0Lhh4BRH5owcclKv6n3YHaBVj
JNigiq1/tlqU46toZvotPOORjpy21kJPZioHqOVCDO4zreMy2xuIiYtpSSmXwkEO
zcQQ3K8sbRx9ED8SCdb229h7ioTug02YBXs0YOQZ024HFaIF8Nz1M+mdHy1jCbLd
yJoT/jzE4RdldZKZJFaSKV1c7EYlzhkCAwEAAaOB0jCBzzAOBgNVHQ8BAf8EBAMC
BaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFCiB
JgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEFBQcwAYYV
aHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8vY2EubXlz
c2wuY29tL215c3NsdGVzdHJzYS5jcnQwGAYDVR0RBBEwD4INdGVzdGxkLnFxLmNv
bTANBgkqhkiG9w0BAQsFAAOCAQEAnPJl1GrePDIulWfsETPbGnrZv3j3ZRXuou0o
K32X/nydS/i/j+AUzKSyezmnR1edkgY1hbGaza702SLQJuGh2IqJvAFyifwV/CZ5
cpJIi5G7kWTBjZo9NgVnDMhR8y5DCKE8BhiUBwcSvKKC8se2yWHm1fk9pRxG0Mc6
0fstl40jtR5XZYsW1GhX4fzwrWuBodPKticgXPn2e24ec+4rVrziu5R7D77AzJjG
Y/wzNYvAUWEzEya7Ve53nhu+WpIuIQn0ux8nPDioFdOjckn4jK3ePYdS2mWT6EBU
BC74GYiBNDz0QgHADq1VTExeLzC0tw9PPdWl0WfoTgCCKLz0yA==
-----END CERTIFICATE-----
EOF
  certificate_name = format("%%s3", var.name)
  key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAuZljaB5MWL7nSABXINAf1jmxO/wstfr4wRPWXeVqvONfUzQD
P1kQhzRM2sCnya0lGWlALPIM9mDzwDBJ0TCpI5KycbN9LYNRw8mZ2fZrafk8u9YT
byo+5J/p+XciE2I/MthHnX7txwkoRBGF84ra6yYmjqhMb/QuGHgFEfmjBxyUq/qf
dgdoFWMk2KCKrX+2WpTjq2hm+i0845GOnLbWQk9mKgeo5UIM7jOt4zLbG4iJi2lJ
KZfCQQ7NxBDcryxtHH0QPxIJ1vbb2HuKhO6DTZgFezRg5BnTbgcVogXw3PUz6Z0f
LWMJst3ImhP+PMThF2V1kpkkVpIpXVzsRiXOGQIDAQABAoIBAQCBPiw4A+k8X2vk
+r+xjNyurCwcTmXAL81rfmnnputmL5tg8DZWtanJzQS7zC7LRPQxttZGtiOKqkbz
DW1J6+3MZMo4XToNKIYWpduqKWvxNusxDkkoPy3evPEMlAY5o0/JE00DgrEHyfut
MtqplocN+tocu1vHFi3HQkSdmM4LE46ZfFu5w1FRbNI1Gqjj/cwlF/T93V3qMap1
WfsJjhMIX9LjBq3y9GAfAtAw7JYwkztr2AhYzCsK25wAj72zFY6FJTZ8LklfS41q
DrVtdjMx42IonDQtkzrqzfYlXdzzhZzuQHxn+qJODseoU8oDG9j3eKhp1dKgqLfx
tv1o3km1AoGBAOqXGEw2w94uVchCjuTum3XFYieEla0IUbHJCaWKU/hoSbht7j23
K7tA9//epBuRLGtYE0sPBK6i31mQT216muspO1g3pwGJhPy8VSpsJ1GQhB2G2UNz
kZlRK+2/gx35TdTi9x0C6UWk5XkhgWO3R35BlEnuV7EOyJunobiUcoObAoGBAMqJ
reuSbJajNGfeBzPel7F62ZDufC85hWaGKzeXIk3DkXcsEpeR5ogMGHCZ3dGf4Yz/
pcfjnCMIWjc+MkA4ppFd4432FJkxNQQP0z7njpXW3e5tionsMN+UwPzi4wZPKufK
osjw43JpBzpHGxG3ynLgZg7bSfrQZhPDTQw84nJbAoGBAKpy8mKeAB71R7rkMXNB
s48Uxca03RQGUWV+DxZKtcxt6fKpXUtWRd4ezJMLL+4fw0iTjCEjXmGNUf9/jVac
mOd44/erKBtD0m7YYIEcaE0pVfUmP8J0vDvL8MEkP56Nv/GIn8hijx/dOiaTI7JS
Pw4LlDVLikfJ2BTQ7f5xTes1AoGBAJ7h4HiDFgIZp1uvtfC/tjn5CEGEhBC7y+VA
bRify747I5rcDP2v66tf6bAzU+pExLhKN++Vov9sZvEdLmhoyGoSwBa2KzR9gHxe
ObYICjeLJfALKHnHuhM6ayY2iieB5UOOF6MQLSysLYpPC3IbvonddNJEvkUuRFVO
iNuHy5AvAoGAPqbzSNA05gf85zRO/JZAmZuWXG3o3pVougbkc211p+ynMpS+/bMb
c/nR36kOE551lFjIoAjoeIs16Wbq+00u9GlcQmyAdFpfaFCNHa3dayJwKJMW9Nia
fKbiiiOAQE8s2v8Paa+b00GspeWLow4u0G5lBVau4JjEVnl6ivLXlzY=
-----END RSA PRIVATE KEY-----
EOF
}

locals {
  certificate_id = join("", [alicloud_ssl_certificates_service_certificate.default3MYZEt.id, "-cn-hangzhou"])
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
					"cert":        "-----BEGIN CERTIFICATE----- MIID5DCCAsygAwIBAgIQWen3GebvT0GcE/a1MJVFgjANBgkqhkiG9w0BAQsFADBe MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe Fw0yNDA4MjMwOTE1MjNaFw0yNTA4MjMwOTE1MjNaMCcxCzAJBgNVBAYTAkNOMRgw FgYDVQQDEw90ZXN0Y2VydC5xcS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw ggEKAoIBAQDF/H2/Oas8trhkHs7B4B8mN00eup/5Tqar2QO4Hm499GJECKc3eMiC v+aAlW74Iymb2Varnv+WMdFRVMQgpXesi3akvVp0QxecvcDliilkh4ddTK731Rd7 PaSK1JdQX1jdGGhVnhQz+cPNFBGZ3tMYGhUkgNfqa3UFucJcBuRub/Ircr+5Ob4D FxSglfTHi+/EFcp7vMAOztLD4zXmEz3NysDNP6NzN7SD72DwPp0nxyRjrBlHSOVg szB/bFasQdAhZGeo64MvSb+SivdWEMhHkwKA5MhhYOkDeNPPSmlxbw0Z3nOyeMmI YkaxzhpO5DZN382duTQmiQ+Yg60OfL3NAgMBAAGjgdQwgdEwDgYDVR0PAQH/BAQD AgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQo gSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYBBQUHMAGG FWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDovL2NhLm15 c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MBoGA1UdEQQTMBGCD3Rlc3RjZXJ0LnFx LmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAKUCvcbUCJRFbVowd1YorILivqRmS6ztR 9vLdj4YZBWxmmQrgkDlkl78r/rXlJbqHunSh2Wbag7y+GaQQwg8xcL4Z3KrKj4zg aHEP1DyFiaxMTEuC/L2RgSX3xlXcf6fQ46D3y3Ja3iHFQnjx6npNaSZ2bSEULvJg IjPiJ/nbid1TYR5vtg8vtwDQfY7/+q8/3DWKcQq+SGcd9dDS6u9vulNdlW8e14bS CKDEzS/axjoICl9JagASLcElWIit/eD5zGIKzTPC9mEXiX/J/gUr70y9GiE7Ue++ 4nFSGQOwjMh/wO2HlmRfToeZ3g6rRCijibBHKHBmVym7NCai2voZJQ== -----END CERTIFICATE-----",
					"key":         "-----BEGIN RSA PRIVATE KEY----- MIIEowIBAAKCAQEAxfx9vzmrPLa4ZB7OweAfJjdNHrqf+U6mq9kDuB5uPfRiRAin N3jIgr/mgJVu+CMpm9lWq57/ljHRUVTEIKV3rIt2pL1adEMXnL3A5YopZIeHXUyu 99UXez2kitSXUF9Y3RhoVZ4UM/nDzRQRmd7TGBoVJIDX6mt1BbnCXAbkbm/yK3K/ uTm+AxcUoJX0x4vvxBXKe7zADs7Sw+M15hM9zcrAzT+jcze0g+9g8D6dJ8ckY6wZ R0jlYLMwf2xWrEHQIWRnqOuDL0m/kor3VhDIR5MCgOTIYWDpA3jTz0ppcW8NGd5z snjJiGJGsc4aTuQ2Td/Nnbk0JokPmIOtDny9zQIDAQABAoIBAQDD6fU4y8UhwCG4 mS+5c6D/PQvoU35Hwkd1l7pxcFNgpTqz3egyISgxEdny9WwoyQq8eJWmICEEK+nY VEv7jiFdMWhG3kTq9RUhejeuLEiHfQE7Fs2w2kFxJ29yHapZ0u/pYOSljFarlATo I2rDW1aB7BVt2L1P7+ONteKZFAzpJckft5ceRUzs5Jm1Cqt8OWO3Km+FBbCROv8M TevW44aoMwBGXuqs06FV1Z4dafglskjt2O38V4acZpH8Nc8j+nCONKL3OxwKY6HQ WfnbXnTLCF3IuMiy8ntrY8HYU6EABiCdr+Pl5HmhI2nmtSFTFbD4Gq70vgPL0P1m iULJGJ7hAoGBANQPrOGe9qHcBydvcBHE7qA9v1+IaTj03qzDTopTi/jxcd9pEkei skLyHNQ5yJT0QjTxB9iYRLfZccOGFyqz/Sdz6CwwTWBZeXOQ2AX7FPEcCnNr1TpF yMrgOY3H93KJISEVS6kYskByjK7XzXCp0KQNS2EeIhAXcqXxNmSwylvpAoGBAO8C PdZHd6aLLEZyVO1aZVHDxqmbhmGVoY9wZ+uwR4K2Hu/fjk0qlR9cYpw8+N675Wr9 E9Ff5/wjK2+/+uocQV9Zoap2vgrwX7GASuO5KYdCOBn6oUOSa+Ru+LgBNyUkXYES mM8eFC1QqfcSrETLAQqd2lmLcuaMq6jJtbBpvzhFAoGAYd8mNC9wtr1dE+dLuvfA BnbZJ1dG8QGa7/NoAVGT7X5JxwmwZR2C1oD1q0FMAOtGzzZbH60PMicKaWousQfH E/lbs2FLpOdGtX6pJQF/5dPCQwkGrVFd3bxk87nRy6vcfW9drxp10mbL5To2WAQY Bk8Ydic5I2IfCNVt/ETX8FkCgYA+OkAtVQgi9WM+qC/SaFGu2yETManoKFQbC3IT HB9SOeaOH49mKesPcjc+ZGWLYDJYC7IoNicpL2L0wnAqmdavY5/CyQ2rvW+8wCE/ bwsP6z6+DNIFzM6IeBgLmE1qPzCVFWlxq2wnbDQEXvk5I/2ObRDXdYYh3ogm9vV2 C+I8XQKBgHIquvifRVvWf1q9WFZLQXZMv1flPhNaLmR+2k6gNpJ8SeiOmBtE7gT6 Je+YOXEKvfr6jaaJwYHPi6IhWHs4fQbgdK4jei30sRL7c8QdKEuwRdHPimmGNAPb UapzHY7xq0Wk9enAnM/SXkjTAJEkrpiQiDuPZVi4sIYCOqb+Ovu5 -----END RSA PRIVATE KEY-----",
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
						"cert":           "-----BEGIN CERTIFICATE----- MIID5DCCAsygAwIBAgIQWen3GebvT0GcE/a1MJVFgjANBgkqhkiG9w0BAQsFADBe MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe Fw0yNDA4MjMwOTE1MjNaFw0yNTA4MjMwOTE1MjNaMCcxCzAJBgNVBAYTAkNOMRgw FgYDVQQDEw90ZXN0Y2VydC5xcS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw ggEKAoIBAQDF/H2/Oas8trhkHs7B4B8mN00eup/5Tqar2QO4Hm499GJECKc3eMiC v+aAlW74Iymb2Varnv+WMdFRVMQgpXesi3akvVp0QxecvcDliilkh4ddTK731Rd7 PaSK1JdQX1jdGGhVnhQz+cPNFBGZ3tMYGhUkgNfqa3UFucJcBuRub/Ircr+5Ob4D FxSglfTHi+/EFcp7vMAOztLD4zXmEz3NysDNP6NzN7SD72DwPp0nxyRjrBlHSOVg szB/bFasQdAhZGeo64MvSb+SivdWEMhHkwKA5MhhYOkDeNPPSmlxbw0Z3nOyeMmI YkaxzhpO5DZN382duTQmiQ+Yg60OfL3NAgMBAAGjgdQwgdEwDgYDVR0PAQH/BAQD AgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQo gSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYBBQUHMAGG FWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDovL2NhLm15 c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MBoGA1UdEQQTMBGCD3Rlc3RjZXJ0LnFx LmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAKUCvcbUCJRFbVowd1YorILivqRmS6ztR 9vLdj4YZBWxmmQrgkDlkl78r/rXlJbqHunSh2Wbag7y+GaQQwg8xcL4Z3KrKj4zg aHEP1DyFiaxMTEuC/L2RgSX3xlXcf6fQ46D3y3Ja3iHFQnjx6npNaSZ2bSEULvJg IjPiJ/nbid1TYR5vtg8vtwDQfY7/+q8/3DWKcQq+SGcd9dDS6u9vulNdlW8e14bS CKDEzS/axjoICl9JagASLcElWIit/eD5zGIKzTPC9mEXiX/J/gUr70y9GiE7Ue++ 4nFSGQOwjMh/wO2HlmRfToeZ3g6rRCijibBHKHBmVym7NCai2voZJQ== -----END CERTIFICATE-----",
						"key":            "-----BEGIN RSA PRIVATE KEY----- MIIEowIBAAKCAQEAxfx9vzmrPLa4ZB7OweAfJjdNHrqf+U6mq9kDuB5uPfRiRAin N3jIgr/mgJVu+CMpm9lWq57/ljHRUVTEIKV3rIt2pL1adEMXnL3A5YopZIeHXUyu 99UXez2kitSXUF9Y3RhoVZ4UM/nDzRQRmd7TGBoVJIDX6mt1BbnCXAbkbm/yK3K/ uTm+AxcUoJX0x4vvxBXKe7zADs7Sw+M15hM9zcrAzT+jcze0g+9g8D6dJ8ckY6wZ R0jlYLMwf2xWrEHQIWRnqOuDL0m/kor3VhDIR5MCgOTIYWDpA3jTz0ppcW8NGd5z snjJiGJGsc4aTuQ2Td/Nnbk0JokPmIOtDny9zQIDAQABAoIBAQDD6fU4y8UhwCG4 mS+5c6D/PQvoU35Hwkd1l7pxcFNgpTqz3egyISgxEdny9WwoyQq8eJWmICEEK+nY VEv7jiFdMWhG3kTq9RUhejeuLEiHfQE7Fs2w2kFxJ29yHapZ0u/pYOSljFarlATo I2rDW1aB7BVt2L1P7+ONteKZFAzpJckft5ceRUzs5Jm1Cqt8OWO3Km+FBbCROv8M TevW44aoMwBGXuqs06FV1Z4dafglskjt2O38V4acZpH8Nc8j+nCONKL3OxwKY6HQ WfnbXnTLCF3IuMiy8ntrY8HYU6EABiCdr+Pl5HmhI2nmtSFTFbD4Gq70vgPL0P1m iULJGJ7hAoGBANQPrOGe9qHcBydvcBHE7qA9v1+IaTj03qzDTopTi/jxcd9pEkei skLyHNQ5yJT0QjTxB9iYRLfZccOGFyqz/Sdz6CwwTWBZeXOQ2AX7FPEcCnNr1TpF yMrgOY3H93KJISEVS6kYskByjK7XzXCp0KQNS2EeIhAXcqXxNmSwylvpAoGBAO8C PdZHd6aLLEZyVO1aZVHDxqmbhmGVoY9wZ+uwR4K2Hu/fjk0qlR9cYpw8+N675Wr9 E9Ff5/wjK2+/+uocQV9Zoap2vgrwX7GASuO5KYdCOBn6oUOSa+Ru+LgBNyUkXYES mM8eFC1QqfcSrETLAQqd2lmLcuaMq6jJtbBpvzhFAoGAYd8mNC9wtr1dE+dLuvfA BnbZJ1dG8QGa7/NoAVGT7X5JxwmwZR2C1oD1q0FMAOtGzzZbH60PMicKaWousQfH E/lbs2FLpOdGtX6pJQF/5dPCQwkGrVFd3bxk87nRy6vcfW9drxp10mbL5To2WAQY Bk8Ydic5I2IfCNVt/ETX8FkCgYA+OkAtVQgi9WM+qC/SaFGu2yETManoKFQbC3IT HB9SOeaOH49mKesPcjc+ZGWLYDJYC7IoNicpL2L0wnAqmdavY5/CyQ2rvW+8wCE/ bwsP6z6+DNIFzM6IeBgLmE1qPzCVFWlxq2wnbDQEXvk5I/2ObRDXdYYh3ogm9vV2 C+I8XQKBgHIquvifRVvWf1q9WFZLQXZMv1flPhNaLmR+2k6gNpJ8SeiOmBtE7gT6 Je+YOXEKvfr6jaaJwYHPi6IhWHs4fQbgdK4jei30sRL7c8QdKEuwRdHPimmGNAPb UapzHY7xq0Wk9enAnM/SXkjTAJEkrpiQiDuPZVi4sIYCOqb+Ovu5 -----END RSA PRIVATE KEY-----",
						"cert_region":    "ap-southeast-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_ext":   "{\\\"Https2http\\\":0,\\\"Http2\\\":0,\\\"Http2https\\\":0}",
					"cert_name":   name + "1",
					"cert":        "-----BEGIN CERTIFICATE----- MIID5DCCAsygAwIBAgIQFO5STIOlR/KkRB3gDHsi5zANBgkqhkiG9w0BAQsFADBe MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe Fw0yNDA4MTMwNTEyMDdaFw0yNTA4MTMwNTEyMDdaMCcxCzAJBgNVBAYTAkNOMRgw FgYDVQQDEw90ZXN0Y2VydC5xcS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw ggEKAoIBAQDwg/HpVcY9yaLrxApP4kvrGwwPbpqOXw9NDWqSD2ms7Qvhi84dTnZT 9xyLXjW0MYI+W13kes4Prl37L8VHOZFvxdzZdzuHvZMkJDaRcekIiUYp+5SFqzzT EtlQodDZBLy0uSlwHnDHymUNJg2nXma+cOOwBBhvb9h+j5v4uuUkQoHRzIUHCkd6 8LR4zbOe+zrhxUU5AYm8C77ZJphDXi9GYm/moajpk+0biCDZ/vZSjTngZEujQKam dgRfVCiWgoUrueiijh1cLGf+W15A3kNo5UXfQrhbHVRQY0vy9dqU1ZBms0pWdc9G dJGd2kXkGNWtEbk4NCN6c64AW2L3G2JNAgMBAAGjgdQwgdEwDgYDVR0PAQH/BAQD AgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQo gSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYBBQUHMAGG FWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDovL2NhLm15 c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MBoGA1UdEQQTMBGCD3Rlc3RjZXJ0LnFx LmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAFl4tmH+m/lGr7aKxxz5uYdqepcjCHTzW KpxxEHr+sL9fOkap6I/NA9MYKsuVfnrTMnEJXeAj24g8q6t38bGTEMcUvh3YGS+P 7DLIF2T2/f9oiO6LVD72zC/HkqaenzWla536G4FLAW4+MzAdn+R7hd5QN4CuFaKI nWWr5KMjZfgyxPkbe+WiD8QryR/EOd3za5vdycZRnBVIcbZQzEHwxizm3JyPNYy8 KeIjt5kGjFt9ccb1oPXjVibB2SUO322kbhMeCXKq2tswnZBICy+1EcbmaxNlkxwc hcdb2wz/3pzH0jbexgTFrNWeoJyYykM8s/3cMY/t0dKppve4FXLG9A== -----END CERTIFICATE-----",
					"key":         "-----BEGIN RSA PRIVATE KEY----- MIIEowIBAAKCAQEA8IPx6VXGPcmi68QKT+JL6xsMD26ajl8PTQ1qkg9prO0L4YvO HU52U/cci141tDGCPltd5HrOD65d+y/FRzmRb8Xc2Xc7h72TJCQ2kXHpCIlGKfuU has80xLZUKHQ2QS8tLkpcB5wx8plDSYNp15mvnDjsAQYb2/Yfo+b+LrlJEKB0cyF BwpHevC0eM2znvs64cVFOQGJvAu+2SaYQ14vRmJv5qGo6ZPtG4gg2f72Uo054GRL o0CmpnYEX1QoloKFK7nooo4dXCxn/lteQN5DaOVF30K4Wx1UUGNL8vXalNWQZrNK VnXPRnSRndpF5BjVrRG5ODQjenOuAFti9xtiTQIDAQABAoIBAQCkuln3bA3ox69U NuKxL9a7Ybzy3NfyZtz98xBolTHVhE083xn+LH0SqQ7dzVqO3dHMj5tRH2L+jnhD z8YYMC+SFDxcnTMilw6uFDdjilcGx65Mlsh0fIGeNyyr8wgtevcb+C2PYunvjImF Zei4FwnbqUnohgWOXVYz6Hv08Vx7ZdW+QiH62I/LS73G7d2EPb26Zo3zMKg/H5AD xuNk82MDW0lCgrw869Yqhcd3GkkmgWi+S71AE0ftY03QeBrsSZbzz4Zsgk2GsEBt fGclOu2c5sNRhLy4o7GiZghPS32zkiec9H76Ip5n/nwXgcYCoxfvOQL+b8U5vpap gbgjWGLBAoGBAPQsJGrYAzt1naGbzwEgYbbyItKD3Bf30YVcZxnZW1fZn2H9NpWl oIBAiO8ls1WP32Mf1us91Bp0V4ESmmklFb5ZllRZXnGt6U2pvxl2FZbK18Nvcvw/ kc66t643mKdfsTAF3sRSmeeWxNSB5C+/yZfHARbBoZ64hYqSut2sivjlAoGBAPwq dDBrz4P/PGn191ZfpWyTC4NRMpC2zYkHoWtMJNs96bsZub4phWyo1KLL4sbAZ+uV GfpRpE2u8mKQWEUU2Gz/KI0cY1e40Icg/ZBlzNglb7ssBXKMrI3R2pigIzEAArhU KwsjhreGF7ix5DEyNT4i7PQL/tOu3uh9SUXl9zVJAoGANta/KxvuxejpiUVUHZ2n NI53Ua55vQxUi04wfba6dCWVTU2wd7WmMYfM+WEPQPU6J6ob++N8AqEEkiGaemjw 1DqMr88OjhuQHXg1SkOiH6bZBLTAL3Ubi0GWRVOJPnYYdn+rA47FsCTFejDeDfdW EHeKgBDm+p3YqEHCJE0/PR0CgYA+mw+zweCIhgLqz81znVWFylAubyddtHT9E27p I8N2xz1TXYS3CLn+i0AXlwUbkUN7ws3rTv+65bd57xprNEyzavoXZrfnXJQxKGir xAqCk3DVCI3lrbVdlH9wKznxfW4vc34oSs60m88h5NChwjRj0+n+gUfoKF9hW1Go z/p7OQKBgCqfExvjjmBibsr0ZtTprMGV2qI55LoVnUXcM1xSBbgU0rLwvw2YBrWg MRl3ixYN851wGD1LhgpwFjr7SnwEhpKdDloKSE1ANM5LE7zvJHsPY2uvmY/Rbn7i RqIzeCbUVYt1Ow0WEqy0DUy/fGLQEz9viwLnvTcDNWeuSsljeAO7 -----END RSA PRIVATE KEY-----",
					"cert_region": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_ext":   "{\"Https2http\":0,\"Http2\":0,\"Http2https\":0}",
						"cert_name":   CHECKSET,
						"cert":        "-----BEGIN CERTIFICATE----- MIID5DCCAsygAwIBAgIQFO5STIOlR/KkRB3gDHsi5zANBgkqhkiG9w0BAQsFADBe MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe Fw0yNDA4MTMwNTEyMDdaFw0yNTA4MTMwNTEyMDdaMCcxCzAJBgNVBAYTAkNOMRgw FgYDVQQDEw90ZXN0Y2VydC5xcS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw ggEKAoIBAQDwg/HpVcY9yaLrxApP4kvrGwwPbpqOXw9NDWqSD2ms7Qvhi84dTnZT 9xyLXjW0MYI+W13kes4Prl37L8VHOZFvxdzZdzuHvZMkJDaRcekIiUYp+5SFqzzT EtlQodDZBLy0uSlwHnDHymUNJg2nXma+cOOwBBhvb9h+j5v4uuUkQoHRzIUHCkd6 8LR4zbOe+zrhxUU5AYm8C77ZJphDXi9GYm/moajpk+0biCDZ/vZSjTngZEujQKam dgRfVCiWgoUrueiijh1cLGf+W15A3kNo5UXfQrhbHVRQY0vy9dqU1ZBms0pWdc9G dJGd2kXkGNWtEbk4NCN6c64AW2L3G2JNAgMBAAGjgdQwgdEwDgYDVR0PAQH/BAQD AgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQo gSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYBBQUHMAGG FWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDovL2NhLm15 c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MBoGA1UdEQQTMBGCD3Rlc3RjZXJ0LnFx LmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAFl4tmH+m/lGr7aKxxz5uYdqepcjCHTzW KpxxEHr+sL9fOkap6I/NA9MYKsuVfnrTMnEJXeAj24g8q6t38bGTEMcUvh3YGS+P 7DLIF2T2/f9oiO6LVD72zC/HkqaenzWla536G4FLAW4+MzAdn+R7hd5QN4CuFaKI nWWr5KMjZfgyxPkbe+WiD8QryR/EOd3za5vdycZRnBVIcbZQzEHwxizm3JyPNYy8 KeIjt5kGjFt9ccb1oPXjVibB2SUO322kbhMeCXKq2tswnZBICy+1EcbmaxNlkxwc hcdb2wz/3pzH0jbexgTFrNWeoJyYykM8s/3cMY/t0dKppve4FXLG9A== -----END CERTIFICATE-----",
						"key":         "-----BEGIN RSA PRIVATE KEY----- MIIEowIBAAKCAQEA8IPx6VXGPcmi68QKT+JL6xsMD26ajl8PTQ1qkg9prO0L4YvO HU52U/cci141tDGCPltd5HrOD65d+y/FRzmRb8Xc2Xc7h72TJCQ2kXHpCIlGKfuU has80xLZUKHQ2QS8tLkpcB5wx8plDSYNp15mvnDjsAQYb2/Yfo+b+LrlJEKB0cyF BwpHevC0eM2znvs64cVFOQGJvAu+2SaYQ14vRmJv5qGo6ZPtG4gg2f72Uo054GRL o0CmpnYEX1QoloKFK7nooo4dXCxn/lteQN5DaOVF30K4Wx1UUGNL8vXalNWQZrNK VnXPRnSRndpF5BjVrRG5ODQjenOuAFti9xtiTQIDAQABAoIBAQCkuln3bA3ox69U NuKxL9a7Ybzy3NfyZtz98xBolTHVhE083xn+LH0SqQ7dzVqO3dHMj5tRH2L+jnhD z8YYMC+SFDxcnTMilw6uFDdjilcGx65Mlsh0fIGeNyyr8wgtevcb+C2PYunvjImF Zei4FwnbqUnohgWOXVYz6Hv08Vx7ZdW+QiH62I/LS73G7d2EPb26Zo3zMKg/H5AD xuNk82MDW0lCgrw869Yqhcd3GkkmgWi+S71AE0ftY03QeBrsSZbzz4Zsgk2GsEBt fGclOu2c5sNRhLy4o7GiZghPS32zkiec9H76Ip5n/nwXgcYCoxfvOQL+b8U5vpap gbgjWGLBAoGBAPQsJGrYAzt1naGbzwEgYbbyItKD3Bf30YVcZxnZW1fZn2H9NpWl oIBAiO8ls1WP32Mf1us91Bp0V4ESmmklFb5ZllRZXnGt6U2pvxl2FZbK18Nvcvw/ kc66t643mKdfsTAF3sRSmeeWxNSB5C+/yZfHARbBoZ64hYqSut2sivjlAoGBAPwq dDBrz4P/PGn191ZfpWyTC4NRMpC2zYkHoWtMJNs96bsZub4phWyo1KLL4sbAZ+uV GfpRpE2u8mKQWEUU2Gz/KI0cY1e40Icg/ZBlzNglb7ssBXKMrI3R2pigIzEAArhU KwsjhreGF7ix5DEyNT4i7PQL/tOu3uh9SUXl9zVJAoGANta/KxvuxejpiUVUHZ2n NI53Ua55vQxUi04wfba6dCWVTU2wd7WmMYfM+WEPQPU6J6ob++N8AqEEkiGaemjw 1DqMr88OjhuQHXg1SkOiH6bZBLTAL3Ubi0GWRVOJPnYYdn+rA47FsCTFejDeDfdW EHeKgBDm+p3YqEHCJE0/PR0CgYA+mw+zweCIhgLqz81znVWFylAubyddtHT9E27p I8N2xz1TXYS3CLn+i0AXlwUbkUN7ws3rTv+65bd57xprNEyzavoXZrfnXJQxKGir xAqCk3DVCI3lrbVdlH9wKznxfW4vc34oSs60m88h5NChwjRj0+n+gUfoKF9hW1Go z/p7OQKBgCqfExvjjmBibsr0ZtTprMGV2qI55LoVnUXcM1xSBbgU0rLwvw2YBrWg MRl3ixYN851wGD1LhgpwFjr7SnwEhpKdDloKSE1ANM5LE7zvJHsPY2uvmY/Rbn7i RqIzeCbUVYt1Ow0WEqy0DUy/fGLQEz9viwLnvTcDNWeuSsljeAO7 -----END RSA PRIVATE KEY-----",
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


`, name)
}

// Test DdosCoo DomainResource. <<< Resource test cases, automatically generated.
