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

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudGaListener_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudGaListenerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaListenerBasicDependence)
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
					"accelerator_id": "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}",
					"description":    "create_description",
					"name":           "${var.name}",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "70",
						},
					},
					"proxy_protocol": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id": CHECKSET,
						"description":    "create_description",
						"name":           name,
						"port_ranges.#":  "1",
						"proxy_protocol": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "SOURCE_IP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "SOURCE_IP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "100",
							"to_port":   "110",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_ranges.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "NONE",
					"protocol":        "UDP",
					"proxy_protocol":  "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "NONE",
						"protocol":        "UDP",
						"proxy_protocol":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "SOURCE_IP",
					"description":     "create_description",
					"protocol":        "TCP",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "70",
						},
					},
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "SOURCE_IP",
						"description":     "create_description",
						"protocol":        "TCP",
						"port_ranges.#":   "1",
						"name":            name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "20",
							"to_port":   "20",
						},
					},
					"certificates": []map[string]string{
						{
							"id": "${local.certificate_id}",
						},
					},
					"proxy_protocol":     "true",
					"protocol":           "HTTPS",
					"security_policy_id": "tls_cipher_policy_1_0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificates.#":     "1",
						"port_ranges.#":      "1",
						"proxy_protocol":     "true",
						"protocol":           "HTTPS",
						"security_policy_id": "tls_cipher_policy_1_0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_id": "tls_cipher_policy_1_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_id": "tls_cipher_policy_1_1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forwarded_for_config": []map[string]interface{}{
						{
							"forwarded_for_ga_id_enabled": "true",
							"forwarded_for_ga_ap_enabled": "true",
							"forwarded_for_proto_enabled": "true",
							"forwarded_for_port_enabled":  "true",
							"real_ip_enabled":             "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forwarded_for_config.#": "1",
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

func TestAccAliCloudGaListener_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudGaListenerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaListenerBasicDependence)
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
					"accelerator_id": "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}",
					"description":    "create_description",
					"name":           "${var.name}",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "10000",
							"to_port":   "16000",
						},
					},
					"proxy_protocol": "true",
					"listener_type":  "CustomRouting",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id": CHECKSET,
						"description":    "create_description",
						"name":           name,
						"port_ranges.#":  "1",
						"proxy_protocol": "true",
						"listener_type":  "CustomRouting",
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

func TestAccAliCloudGaListener_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudGaListenerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaListenerBasicDependence)
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
					"accelerator_id": "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}",
					"description":    "create_description",
					"name":           "${var.name}",
					"protocol":       "HTTP",
					"proxy_protocol": "true",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "60",
						},
					},
					"forwarded_for_config": []map[string]interface{}{
						{
							"forwarded_for_ga_id_enabled": "true",
							"forwarded_for_ga_ap_enabled": "true",
							"forwarded_for_proto_enabled": "true",
							"forwarded_for_port_enabled":  "true",
							"real_ip_enabled":             "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":         CHECKSET,
						"description":            "create_description",
						"name":                   name,
						"protocol":               "HTTP",
						"proxy_protocol":         "true",
						"port_ranges.#":          "1",
						"forwarded_for_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forwarded_for_config": []map[string]interface{}{
						{
							"forwarded_for_ga_id_enabled": "false",
							"forwarded_for_ga_ap_enabled": "false",
							"forwarded_for_proto_enabled": "false",
							"forwarded_for_port_enabled":  "false",
							"real_ip_enabled":             "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forwarded_for_config.#": "1",
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

var AliCloudGaListenerMap = map[string]string{
	"client_affinity": "NONE",
	"protocol":        "TCP",
	"status":          CHECKSET,
}

func AliCloudGaListenerBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	locals {
  		certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "%s"])
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
		bandwidth_billing_type = "BandwidthPackage"
	}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV
BAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP
MA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0
ZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow
djELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE
ChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG
9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ
AoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB
coG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook
KOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw
HQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy
+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC
QkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN
MAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ
AJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT
cQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1
Ofi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd
DUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=
-----END CERTIFICATE-----
EOF
  certificate_name = join("-", [var.name, 1])

  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV
kg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM
ywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB
AoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd
6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP
hwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4
MdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz
71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm
Ev9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE
qygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8
9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM
zWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe
DrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=
-----END RSA PRIVATE KEY-----
EOF
}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth              = 100
  		type                   = "Basic"
  		bandwidth_type         = "Basic"
  		payment_type           = "PayAsYouGo"
  		billing_type           = "PayBy95"
  		ratio                  = 30
  		bandwidth_package_name = var.name
  		auto_pay               = true
  		auto_use_coupon        = true
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}
`, name, defaultRegionToTest)
}

func TestUnitAliCloudGaListener(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id": "CreateListenerValue",
		"description":    "CreateListenerValue",
		"name":           "CreateListenerValue",
		"port_ranges": []map[string]interface{}{
			{
				"from_port": 60,
				"to_port":   70,
			},
		},
		"proxy_protocol": true,
		"certificates": []map[string]interface{}{
			{
				"id": "CreateListenerValue",
			},
		},
		"client_affinity": "CreateListenerValue",
		"protocol":        "CreateListenerValue",
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
		// DescribeListener
		"Description": "CreateListenerValue",
		"Certificates": []interface{}{
			map[string]interface{}{
				"Id": "CreateListenerValue",
			},
		},
		"ClientAffinity": "CreateListenerValue",
		"Name":           "CreateListenerValue",
		"PortRanges": []interface{}{
			map[string]interface{}{
				"FromPort": 60,
				"ToPort":   70,
			},
		},
		"State":      "active",
		"Protocol":   "CreateListenerValue",
		"ListenerId": "CreateListenerValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateListener
		"ListenerId": "CreateListenerValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_listener", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeListener Response
		"ListenerId": "CreateListenerValue",
	}
	errorCodes := []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateListener" {
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
		err := resourceAliCloudGaListenerCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateListener
	attributesDiff := map[string]interface{}{
		"accelerator_id": "UpdateListenerValue",
		"description":    "UpdateListenerValue",
		"name":           "UpdateListenerValue",
		"port_ranges": []map[string]interface{}{
			{
				"from_port": 70,
				"to_port":   80,
			},
		},
		"proxy_protocol": true,
		"certificates": []map[string]interface{}{
			{
				"id": "UpdateListenerValue",
			},
		},
		"client_affinity": "UpdateListenerValue",
		"protocol":        "UpdateListenerValue",
	}
	diff, err := newInstanceDiff("alicloud_ga_listener", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeListener Response
		"Description": "UpdateListenerValue",
		"Certificates": []interface{}{
			map[string]interface{}{
				"Id": "UpdateListenerValue",
			},
		},
		"ClientAffinity": "UpdateListenerValue",
		"Name":           "UpdateListenerValue",
		"PortRanges": []interface{}{
			map[string]interface{}{
				"FromPort": 70,
				"ToPort":   80,
			},
		},
		"Protocol":   "UpdateListenerValue",
		"ListenerId": "UpdateListenerValue",
	}
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateListener" {
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
		err := resourceAliCloudGaListenerUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeListener" {
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
		err := resourceAliCloudGaListenerRead(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteListener" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "DescribeListener" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaListenerDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
