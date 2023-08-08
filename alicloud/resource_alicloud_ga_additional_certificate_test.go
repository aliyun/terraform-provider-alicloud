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

func TestAccAliCloudGaAdditionalCertificate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_additional_certificate.default"
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudGaAdditionalCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAdditionalCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaadditionalcertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaAdditionalCertificateBasicDependence0)
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
					"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":    "${alicloud_ga_listener.default.id}",
					"domain":         "${local.domain}",
					"certificate_id": "${local.certificate_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id": CHECKSET,
						"listener_id":    CHECKSET,
						"domain":         CHECKSET,
						"certificate_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_id": "${local.certificate_id_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_id": CHECKSET,
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

var AliCloudGaAdditionalCertificateMap0 = map[string]string{}

func AliCloudGaAdditionalCertificateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
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

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		count            = 2
  		certificate_name = join("-", [var.name, count.index])
		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID7zCCAtegAwIBAgIRAKi2/Fx1cUTyhV839x42ockwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjMwODA5MDQ1NDU3WhcNMjYwODA4MDQ1NDU3WjAsMQswCQYDVQQGEwJDTjEd
MBsGA1UEAxMUYWxpY2xvdWQtcHJvdmlkZXIuY24wggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQDdkot9e0pMCTPAtA29Sz5sF+aPT/l9+3sOnQeJ1kKLNkqK
iQgwADexoAqlmTaZM03gh/GnkqPw9gxN/fJHWdVzxE03Fs8bKgMdS6cf0v/xArrQ
zm6N4vmsbuE8SX2eu303PAsyBMqPByTODZ5i+5LkZcrxMFQsbA3xnBouzS5e+T+a
7YTyyVv5WDy871/sdRAYTfnUttdnqkKGeMKgQgRlJ2pDk5/k2iwmQmSh/wbk465+
1U5w2npPYGPvGAkzl7RRc4/VckqlV8P0cmgguqIRyllJwFEnvcpqpOHTxBOBq9iZ
4b/h7ynrfB/GbAw574eSEl0gzLBW60bT9YedbTeXAgMBAAGjgdkwgdYwDgYDVR0P
AQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSME
GDAWgBQogSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYB
BQUHMAGGFWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDov
L2NhLm15c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MB8GA1UdEQQYMBaCFGFsaWNs
b3VkLXByb3ZpZGVyLmNuMA0GCSqGSIb3DQEBCwUAA4IBAQCwUBeznv6cAjcTLCDb
SSvgkM9HFcbWnuGS8Nf5P4YfmSs52VuHZyjzwphjAU6B/danI/nMdZe52PXyvjVV
02Y8ld/tMpqPV5SpaOadLtdg6TGBNJieOAt9doM8WNEgq/JycAL9ivIOjChUetZf
ZEV7HDIgiHSpqAPWMZYL71MS/p5zYkyOnPqmGyLNdi1neotwVCQopQXRNC2iLlVV
yQONfXH5iijqr1iTWkB0ESK/xBt1PB655PlTjzFQUOovE1SyoQS8K3u7TP6+BqtD
G9TYNTNZvxl5I/iU/KdWVip+qJbxRA8Skc8gHkkzeIEStw3l5cjnrp9h7EhnhkOh
ltGN
-----END CERTIFICATE-----
EOF
		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA3ZKLfXtKTAkzwLQNvUs+bBfmj0/5fft7Dp0HidZCizZKiokI
MAA3saAKpZk2mTNN4Ifxp5Kj8PYMTf3yR1nVc8RNNxbPGyoDHUunH9L/8QK60M5u
jeL5rG7hPEl9nrt9NzwLMgTKjwckzg2eYvuS5GXK8TBULGwN8ZwaLs0uXvk/mu2E
8slb+Vg8vO9f7HUQGE351LbXZ6pChnjCoEIEZSdqQ5Of5NosJkJkof8G5OOuftVO
cNp6T2Bj7xgJM5e0UXOP1XJKpVfD9HJoILqiEcpZScBRJ73KaqTh08QTgavYmeG/
4e8p63wfxmwMOe+HkhJdIMywVutG0/WHnW03lwIDAQABAoIBAQCe5rHS09B8pzzO
PlJ8JrIlox5eOOScTPX7jPITD+25GL5si8mrYvyODlCUYkSdqgV3uQa9PpUEAfDh
HfXa5boGxAj8MQdmW8LQB6lbUV7r4SFJDkKKzvRvjTVKnwnQBHXQXudIf9ckq+Lh
QzMLmY/G7JmWTyqOkQ+O7nx4g/11bcU7uQrQdvWPfc0+IiT1TYQdyLQ/Chlj3RF/
iwF8ZL2sfKF+Z5O49+Q6cXvUcQOvqtkIXbQijayyVNBMJwDB7aOZRA7JBNj9/ib6
N0iTo81dJVz/nnpbWRaFTVinIsDF1heDfQ1qDx06T/Mpi6pjoWjRUcyIHEbZJTel
0nXDJD1BAoGBAPZB/PN8MP+o9gkf2jnoU9LzctDJrQwD1J2XElq4RomimPIMqDQP
5TRAJThf0O0X4Mv2n9EzV457OpJL+fz9htRWEYogWl9bkbzZ1AoX4K/acuGeawTT
YEhPjJ2ZETsBsCeDkDDuHHzYwRQv+EfoXH36z9PBDxG1ZDb7kWwAILXdAoGBAOZW
jXG7m4I7cxUtXGtjwydh4K7nwH/5QoH2m928HM2AT48eQCl3CMQ089+qeJGgfHQv
GyVOO/FGhcFsFi10FMQ7IlwWgZODg64qnrNhi4zbV1M2wKem1T2dlEpkd82EFdnS
GYRIEkFORMxEDyzx3Th2TajpWC8YKKG3Tnm0bQ4DAoGBAIZTEEswHvoVi78GZN7Z
X3/d028X0xCOtlcPpK9ffPpuesbtKILdeMS7iJHrkecB81jOOfa+7q+FgDl0v/PD
xtvj5sVVSHZjWGeO2h53T9QccDWpV+7V7dsDqUv9xmxNS20CUpCeEWP4R7lfQSrY
EDuXp+11jWa3buae6n/iwfTxAoGABEYW2cVhXUk9GWd+D4AKXvCx+ozSRY2abk7l
FXgoEKgQ0db92ccboohY/g1rr0gLBxzYpBiPhCqK0MvwnWdJ+1odiRfhz5rhFpoz
16A3tqVbOXAKoxG1Yy9JURgMIQQSY7hCQPIVZKDPJfsdTPgv4pxPVJL/z9/i4R1F
l3yBiYECgYEA0+vpzL24nHZYdwgBF4qbmYhv8baRi07/BNgV1+d6vESuO/MwwoE/
2UZ9Drf5yoX2Bvi5/vVMbyc7cSluO7icPBkl0D8F7E3x0v5mzwPxtpR8BTRoJKOL
/rMdLscMz2VQsL5DJd/9OZg60fHRaRtWtV0afXzL5zUxnfDLot24IG4=
-----END RSA PRIVATE KEY-----
EOF
	}

	resource "alicloud_ssl_certificates_service_certificate" "update" {
  		certificate_name = "${var.name}-update"
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID7jCCAtagAwIBAgIQUNnSVa/sQNeb9pBN9NhkwTANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yMzA4MDkwMzM4MThaFw0yODA4MDcwMzM4MThaMCwxCzAJBgNVBAYTAkNOMR0w
GwYDVQQDExRhbGljbG91ZC1wcm92aWRlci5jbjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAOgskr8dEfZYdjr0xaIqlCkmE802vABoj3SQNn3rLWnUj+1v
Wqbpsj6Bu61Scb8mtl/OZOOM7sgq0Q1hpdO8xvMGxTMuZ2bjX0EqCMqh4AvFofHL
a/iVD07hfoM1Jo8CEidh1uvcOuXP1TlaqU020x1TX3a3niJu4JVkmCkCOwAbWYuj
O8IsgBCsFaF9d4+C1JRYOtRbIHCNhd0sxG8AGovUDLvlkePeH5NF7DNvFXgGJ4iv
EQcY9pP08RBFUkaznOw/r64Up7zhLb+Ie4SyAvs1FulhMAmIXOcbsND39hJ+/WIP
8beWvIN1eCS8zcvgAvDgMkV8oqqVbQu1dqx5WuMCAwEAAaOB2TCB1jAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQY
MBaAFCiBJgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEF
BQcwAYYVaHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8v
Y2EubXlzc2wuY29tL215c3NsdGVzdHJzYS5jcnQwHwYDVR0RBBgwFoIUYWxpY2xv
dWQtcHJvdmlkZXIuY24wDQYJKoZIhvcNAQELBQADggEBALd0hFZAd2XHJgETbHQs
h4YUBNKxrIy6JiWfxffhIL1ZK5pI443DC4VRGfxVi3zWqs01WbNtJ2b1KdfSoovH
Zwi3hdMF1IwoAB/Y2sS4zjqS0H1od7MN9KKHes6bl3yCgpmaYs5cHbyg0IJHmeq3
rCgbKsvHfUwtzBNNPHlpANakAYd/5O1pztmUskWMUVaExfpMoQLo/AX9Lqm8pVjw
xs921I703l/E5zEnd3PVSYagy/KQJrwVt+wQZS11HsAryfO9kct/9f+c85VDo6Ht
iRirW/EnNPQRSno4z0V2x1Rn5+ZaoJo8cWzPvKrdfCG9TUozt4AR/LIudNLb6NNW
n7g=
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA6CySvx0R9lh2OvTFoiqUKSYTzTa8AGiPdJA2festadSP7W9a
pumyPoG7rVJxvya2X85k44zuyCrRDWGl07zG8wbFMy5nZuNfQSoIyqHgC8Wh8ctr
+JUPTuF+gzUmjwISJ2HW69w65c/VOVqpTTbTHVNfdreeIm7glWSYKQI7ABtZi6M7
wiyAEKwVoX13j4LUlFg61FsgcI2F3SzEbwAai9QMu+WR494fk0XsM28VeAYniK8R
Bxj2k/TxEEVSRrOc7D+vrhSnvOEtv4h7hLIC+zUW6WEwCYhc5xuw0Pf2En79Yg/x
t5a8g3V4JLzNy+AC8OAyRXyiqpVtC7V2rHla4wIDAQABAoIBABKGQ+sluaIrKrvH
feFTfmDOHfRYsqVhslh9jSt80THJePZb1SLOMJ+WIFBS7Kpwv0pjoF8bho3IBMgJ
i36aaFFJsABGao+mApqjbPIl+kdWLHarYWEDG6aSjVKQshPk+WfVAZ3uA3EEpSGf
XzS+9Bc56LsDKYXbzOV+kjlraSO35AMec3CpISdx4K1caEAhKX6it9bvPq4pSYXi
PQspba0Jv46VV7MaabVjLzsinz5/md4vxyYHNIJAukHUfwJIsVC9ZNxukwSw+CzE
MMO64ylq2DGokNerGsLetuViV8UWi7qmUmms2fAmchodW16olgNkYTz27+V/A42S
eex63pkCgYEA+CqKhqp3qPe2E9KVrycrwjoycxmhOn3Iz1xiN7uAEv+DzfKtfZVf
mcOIiqw4Z82RkgjHb9vJuTigKdDkB1zE2gSDnep44sDWJM/5nPjGlMgnkiJWJhci
CnD0P4d6cT5wyDt7Q0/tS6ql2UrCpW4ktw1AP0Rm/z/VBD8jGkVenjcCgYEA74DM
Z2Qmh3bPt1TykpOlw+H+sEuvlkYxqMlbtn3Rv3WgEPIBekOFrgP7n/uLW1Aizn8w
EhNBBAE8w5jvklqZWYbpFMJQc09eqUkI8aTbLooZbzYj1f3CrzBRKn1GoTPmN9V0
j9r+TbH3/5CEoqlsJdmeQPofuv5Qid2oEutZcrUCgYBuZ16hco0xmqJiRzlYZvDM
w99V3X0g7Hy947e+W6gqy4nzwZb1W9LgMWE5cEzXwViVw1oWpY0k3dBDSi9oJxlc
dM2pH3sQRgH+9pdyAis2XaVdGfGBmKEITCAdc0RBxSmfqva3h4NmOlD2TpAx0MJ8
vWRrwR6hR+CYtw4CzgG+GQKBgQDGmi5lugW9JUe/xeBUrbyyv0+cT1auLUz2ouq7
XIA23Mo74wJYqW9Lyp+4nTWFJeGHDK8G/hJWyNPjeomG+jvZombbQPrHc9SSWi7h
eowKfpfywZlb1M7AyTc1HacY+9l3CTlcJQPl16NHuEZUQFue02NIjGENhd+xQy4h
ainFVQKBgAoPs9ebtWbBOaqGpOnquzb7WibvW/ifOfzv5/aOhkbpp+KCjcSON6CB
QF3BEXMcNMGWlpPrd8PaxCAzR4MyU///ekJri2icS9lrQhGSz2TtYhdED4pv1Aag
7eTPl5L7xAwphCSwy8nfCKmvlqcX/MSJ7A+LHB/2hdbuuEOyhpbu
-----END RSA PRIVATE KEY-----
EOF
	}

	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		name           = var.name
  		protocol       = "HTTPS"
  		port_ranges {
    		from_port = 8080
    		to_port   = 8080
  		}
  		certificates {
    		id = join("-", [alicloud_ssl_certificates_service_certificate.default.0.id, "%s"])
  		}
	}

locals {
  domain                = "alicloud-provider.cn"
  certificate_id        = join("-", [alicloud_ssl_certificates_service_certificate.default.1.id, "%s"])
  certificate_id_update = join("-", [alicloud_ssl_certificates_service_certificate.update.id, "%s"])
}
`, name, defaultRegionToTest, defaultRegionToTest, defaultRegionToTest)
}

func TestUnitAliCloudGaAdditionalCertificate(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ga_additional_certificate"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_additional_certificate"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id": "AssociateAdditionalCertificatesWithListenerValue",
		"domain":         "AssociateAdditionalCertificatesWithListenerValue",
		"certificate_id": "AssociateAdditionalCertificatesWithListenerValue",
		"listener_id":    "AssociateAdditionalCertificatesWithListenerValue",
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
		// ListListenerCertificates
		"Certificates": []interface{}{
			map[string]interface{}{
				"CertificateId": "AssociateAdditionalCertificatesWithListenerValue",
				"Domain":        "AssociateAdditionalCertificatesWithListenerValue",
			},
		},
		"State": "active",
	}
	CreateMockResponse := map[string]interface{}{
		// AssociateAdditionalCertificatesWithListener

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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_additional_certificate", errorCode))
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGaAdditionalCertificateCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListListenerCertificates Response
		"Certificates": []interface{}{
			map[string]interface{}{
				"CertificateId": "AssociateAdditionalCertificatesWithListenerValue",
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AssociateAdditionalCertificatesWithListener" {
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
		err := resourceAliCloudGaAdditionalCertificateCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_additional_certificate"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_ga_additional_certificate", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_additional_certificate"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListListenerCertificates" {
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
		err := resourceAliCloudGaAdditionalCertificateRead(dExisted, rawClient)
		patches.Reset()
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGaAdditionalCertificateDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ga_additional_certificate", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_additional_certificate"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DissociateAdditionalCertificatesFromListener" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"State": "active",
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaAdditionalCertificateDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
