package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Esa ClientCaCertificate. >>> Resource test cases, automatically generated.
// Case resource_ClientCaCertificate_set_test 12003
func TestAccAliCloudEsaClientCaCertificate_basic12003(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_client_ca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaClientCaCertificateMap12003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaClientCaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaClientCaCertificateBasicDependence12003)
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
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"certificate": "-----BEGIN CERTIFICATE-----\\nMIIDRjCCAi6gAwIBAgIUMyxyjbmKN+WbxMEC+mMyojIIjLAwDQYJKoZIhvcNAQEF\\nBQAwPDEXMBUGA1UEAwwOU1NMZXllIFJvb3QgQ0ExFDASBgNVBAoMC1NTTGV5ZSBJ\\nbmMuMQswCQYDVQQGEwJDTjAeFw0yNjA0MDEwNTMwMDBaFw0yNjA2MzAwNTMwMDBa\\nMFwxCzAJBgNVBAYTAkNOMQ8wDQYDVQQIDAbljJfkuqwxDzANBgNVBAcMBuWMl+S6\\nrDENMAsGA1UECgwEMTEyMzEcMBoGA1UEAwwTZXhhbXBsZS10ZXN0MTIzLmNvbTCC\\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANGsJLaIkRXg37YMuIHDkZg2\\n/xJ2R/ckiRwzdvkF3PVt3ceG7ZyAUcREspw9ytuTnCaAZh7kCWduNxsmR1WhURI4\\nJ2xYRjzjWTkzWsuKcKpHg12kkZ38+s3d3cG/TfSqHWu213kFsoeAT3O6zauWu+9r\\nrl6N6aYUZZnPSguiRfqaw67+KR6g8KEC6dgmG2a3QzMUWxNz8yvxs+1j5xgMCzic\\nAl133UIryKrDwciXGocAdzXdtjRQZPaHipewG3XrPmyuh2f5K4imaRPrbO2yPAc8\\ns/XSEGjw3hEdGLjAEJDd1GFAnb/dDVA4DAmwPmZqT53GnrcgSXCH7mTjmcP34K0C\\nAwEAAaMgMB4wCwYDVR0RBAQwAoIAMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcN\\nAQEFBQADggEBAITEC5mQd7Wc7uY34FcttRNXQbOcQKQVoxfY8R6xThbsPOZZanvj\\nI6FZ30PAUm1A25WOZMHlPUXYpmyDD0RAQNDYgcLi9UcYY1YYqLqIrUoZNqnOB7ZF\\n2vdHPKCcHzUSroY5fx/9er0DJwY8x+PvSmNS8/9kq+jRa1QMWY7F383tAORs/vKu\\nW1rvvgBVco8e+RFamU54OVCrATNtpzOUNFJ2g0IDVYhUlpab7Lx5C6IGccCAIGfw\\nmc7EVDIDUcJlDTxx5N/ILWl+1M4be3335ZWTUhxp7gTOA+XIUKHCQE6c47OIUDLg\\niVUaDyfcbHu8BDlZi1bYL4pl+3RndU/LBMg=\\n-----END CERTIFICATE-----",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":     CHECKSET,
						"certificate": "-----BEGIN CERTIFICATE-----\nMIIDRjCCAi6gAwIBAgIUMyxyjbmKN+WbxMEC+mMyojIIjLAwDQYJKoZIhvcNAQEF\nBQAwPDEXMBUGA1UEAwwOU1NMZXllIFJvb3QgQ0ExFDASBgNVBAoMC1NTTGV5ZSBJ\nbmMuMQswCQYDVQQGEwJDTjAeFw0yNjA0MDEwNTMwMDBaFw0yNjA2MzAwNTMwMDBa\nMFwxCzAJBgNVBAYTAkNOMQ8wDQYDVQQIDAbljJfkuqwxDzANBgNVBAcMBuWMl+S6\nrDENMAsGA1UECgwEMTEyMzEcMBoGA1UEAwwTZXhhbXBsZS10ZXN0MTIzLmNvbTCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANGsJLaIkRXg37YMuIHDkZg2\n/xJ2R/ckiRwzdvkF3PVt3ceG7ZyAUcREspw9ytuTnCaAZh7kCWduNxsmR1WhURI4\nJ2xYRjzjWTkzWsuKcKpHg12kkZ38+s3d3cG/TfSqHWu213kFsoeAT3O6zauWu+9r\nrl6N6aYUZZnPSguiRfqaw67+KR6g8KEC6dgmG2a3QzMUWxNz8yvxs+1j5xgMCzic\nAl133UIryKrDwciXGocAdzXdtjRQZPaHipewG3XrPmyuh2f5K4imaRPrbO2yPAc8\ns/XSEGjw3hEdGLjAEJDd1GFAnb/dDVA4DAmwPmZqT53GnrcgSXCH7mTjmcP34K0C\nAwEAAaMgMB4wCwYDVR0RBAQwAoIAMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcN\nAQEFBQADggEBAITEC5mQd7Wc7uY34FcttRNXQbOcQKQVoxfY8R6xThbsPOZZanvj\nI6FZ30PAUm1A25WOZMHlPUXYpmyDD0RAQNDYgcLi9UcYY1YYqLqIrUoZNqnOB7ZF\n2vdHPKCcHzUSroY5fx/9er0DJwY8x+PvSmNS8/9kq+jRa1QMWY7F383tAORs/vKu\nW1rvvgBVco8e+RFamU54OVCrATNtpzOUNFJ2g0IDVYhUlpab7Lx5C6IGccCAIGfw\nmc7EVDIDUcJlDTxx5N/ILWl+1M4be3335ZWTUhxp7gTOA+XIUKHCQE6c47OIUDLg\niVUaDyfcbHu8BDlZi1bYL4pl+3RndU/LBMg=\n-----END CERTIFICATE-----",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_certificate_hostnames": []string{
						"test1.example-test123.com", "test2.example-test123.com", "test3.example-test123.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_certificate_hostnames.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_certificate_hostnames": []string{
						"test3.example-test123.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_certificate_hostnames.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_certificate_hostnames": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_certificate_hostnames.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_certificate_hostnames": []string{
						"test1.example-test123.com", "test2.example-test123.com", "test3.example-test123.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_certificate_hostnames.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudEsaClientCaCertificate_basic12003_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_client_ca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaClientCaCertificateMap12003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaClientCaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaClientCaCertificateBasicDependence12003)
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
					"client_ca_cert_name": name,
					"site_id":             "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"client_ca_certificate_hostnames": []string{
						"test1.example-test123.com", "test2.example-test123.com", "test3.example-test123.com"},
					"certificate": "-----BEGIN CERTIFICATE-----\\nMIIDRjCCAi6gAwIBAgIUMyxyjbmKN+WbxMEC+mMyojIIjLAwDQYJKoZIhvcNAQEF\\nBQAwPDEXMBUGA1UEAwwOU1NMZXllIFJvb3QgQ0ExFDASBgNVBAoMC1NTTGV5ZSBJ\\nbmMuMQswCQYDVQQGEwJDTjAeFw0yNjA0MDEwNTMwMDBaFw0yNjA2MzAwNTMwMDBa\\nMFwxCzAJBgNVBAYTAkNOMQ8wDQYDVQQIDAbljJfkuqwxDzANBgNVBAcMBuWMl+S6\\nrDENMAsGA1UECgwEMTEyMzEcMBoGA1UEAwwTZXhhbXBsZS10ZXN0MTIzLmNvbTCC\\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANGsJLaIkRXg37YMuIHDkZg2\\n/xJ2R/ckiRwzdvkF3PVt3ceG7ZyAUcREspw9ytuTnCaAZh7kCWduNxsmR1WhURI4\\nJ2xYRjzjWTkzWsuKcKpHg12kkZ38+s3d3cG/TfSqHWu213kFsoeAT3O6zauWu+9r\\nrl6N6aYUZZnPSguiRfqaw67+KR6g8KEC6dgmG2a3QzMUWxNz8yvxs+1j5xgMCzic\\nAl133UIryKrDwciXGocAdzXdtjRQZPaHipewG3XrPmyuh2f5K4imaRPrbO2yPAc8\\ns/XSEGjw3hEdGLjAEJDd1GFAnb/dDVA4DAmwPmZqT53GnrcgSXCH7mTjmcP34K0C\\nAwEAAaMgMB4wCwYDVR0RBAQwAoIAMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcN\\nAQEFBQADggEBAITEC5mQd7Wc7uY34FcttRNXQbOcQKQVoxfY8R6xThbsPOZZanvj\\nI6FZ30PAUm1A25WOZMHlPUXYpmyDD0RAQNDYgcLi9UcYY1YYqLqIrUoZNqnOB7ZF\\n2vdHPKCcHzUSroY5fx/9er0DJwY8x+PvSmNS8/9kq+jRa1QMWY7F383tAORs/vKu\\nW1rvvgBVco8e+RFamU54OVCrATNtpzOUNFJ2g0IDVYhUlpab7Lx5C6IGccCAIGfw\\nmc7EVDIDUcJlDTxx5N/ILWl+1M4be3335ZWTUhxp7gTOA+XIUKHCQE6c47OIUDLg\\niVUaDyfcbHu8BDlZi1bYL4pl+3RndU/LBMg=\\n-----END CERTIFICATE-----",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_cert_name":               name,
						"site_id":                           CHECKSET,
						"client_ca_certificate_hostnames.#": "3",
						"certificate":                       "-----BEGIN CERTIFICATE-----\nMIIDRjCCAi6gAwIBAgIUMyxyjbmKN+WbxMEC+mMyojIIjLAwDQYJKoZIhvcNAQEF\nBQAwPDEXMBUGA1UEAwwOU1NMZXllIFJvb3QgQ0ExFDASBgNVBAoMC1NTTGV5ZSBJ\nbmMuMQswCQYDVQQGEwJDTjAeFw0yNjA0MDEwNTMwMDBaFw0yNjA2MzAwNTMwMDBa\nMFwxCzAJBgNVBAYTAkNOMQ8wDQYDVQQIDAbljJfkuqwxDzANBgNVBAcMBuWMl+S6\nrDENMAsGA1UECgwEMTEyMzEcMBoGA1UEAwwTZXhhbXBsZS10ZXN0MTIzLmNvbTCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANGsJLaIkRXg37YMuIHDkZg2\n/xJ2R/ckiRwzdvkF3PVt3ceG7ZyAUcREspw9ytuTnCaAZh7kCWduNxsmR1WhURI4\nJ2xYRjzjWTkzWsuKcKpHg12kkZ38+s3d3cG/TfSqHWu213kFsoeAT3O6zauWu+9r\nrl6N6aYUZZnPSguiRfqaw67+KR6g8KEC6dgmG2a3QzMUWxNz8yvxs+1j5xgMCzic\nAl133UIryKrDwciXGocAdzXdtjRQZPaHipewG3XrPmyuh2f5K4imaRPrbO2yPAc8\ns/XSEGjw3hEdGLjAEJDd1GFAnb/dDVA4DAmwPmZqT53GnrcgSXCH7mTjmcP34K0C\nAwEAAaMgMB4wCwYDVR0RBAQwAoIAMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcN\nAQEFBQADggEBAITEC5mQd7Wc7uY34FcttRNXQbOcQKQVoxfY8R6xThbsPOZZanvj\nI6FZ30PAUm1A25WOZMHlPUXYpmyDD0RAQNDYgcLi9UcYY1YYqLqIrUoZNqnOB7ZF\n2vdHPKCcHzUSroY5fx/9er0DJwY8x+PvSmNS8/9kq+jRa1QMWY7F383tAORs/vKu\nW1rvvgBVco8e+RFamU54OVCrATNtpzOUNFJ2g0IDVYhUlpab7Lx5C6IGccCAIGfw\nmc7EVDIDUcJlDTxx5N/ILWl+1M4be3335ZWTUhxp7gTOA+XIUKHCQE6c47OIUDLg\niVUaDyfcbHu8BDlZi1bYL4pl+3RndU/LBMg=\n-----END CERTIFICATE-----",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudEsaClientCaCertificateMap12003 = map[string]string{
	"status":            CHECKSET,
	"create_time":       CHECKSET,
	"client_ca_cert_id": CHECKSET,
}

func AliCloudEsaClientCaCertificateBasicDependence12003(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}
`, name)
}

// Test Esa ClientCaCertificate. <<< Resource test cases, automatically generated.

// generated by https://www.ssleye.com/ssltool/self_sign.html
const testEsaClientCaCertificate = `-----BEGIN CERTIFICATE-----\nMIIDRTCCAi2gAwIBAgIUHRPTIPKP2zN9on/NCzBe0BV68UUwDQYJKoZIhvcNAQEF\nBQAwMzEPMA0GA1UEAwwGU1NMZXllMRMwEQYDVQQKDApTU0xleWUgSW5jMQswCQYD\nVQQGEwJDTjAeFw0yNTA3MzAwODQzMDBaFw0yNTEyMzEwODQwMDBaMGQxCzAJBgNV\nBAYTAkNOMQ8wDQYDVQQIDAbljJfkuqwxEDAOBgNVBAcMB0JlaWppbmcxGzAZBgNV\nBAoMEuenkeaKgOaciemZkOWFrOWPuDEVMBMGA1UEAwwMZ29zaXRlY2RuLmNuMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtu2oW3t2bj9LsFnXj1C2EmaR\nJYJwNgHsTBKl3DxeL2+Ext0qN2Z+UgTqYM1c1HOdwN9x13pnAVe4PmiLAkxpp/4u\n5gKsH1+6p3aXFUk0NvEoLXfESoQpyvoB0o/8oryxNs3+iUfvAk+a7IKAr99a1P9F\nTkpyE6t+dgSLYhHc49ZRdYImmZcYQLmpygYOwWBdv6hlQUFi/tvX16fRZ0GgyUOK\n7xsTWG6qUhPJyLRtj9zn+0khgh5DJhfJQ4KTWZMX63UPiIx7sPu9sR+TPWqJsEuq\nVipxouMys+NNMjDtn55+PE6/sDbkvULHeFUglGMZ9qHcl3ej31zmkhu6bmvNcQID\nAQABoyAwHjALBgNVHREEBDACggAwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0B\nAQUFAAOCAQEAF6J9TdaDYQ96EaWvb2ttQ6jNrDe4k3t1cdfhPEWMJzxZFxoDBYZ2\nAl9vB2JICEsGDkCwpqYz2UXJsGnq2rHjUxouYo1568K/loownWjwdCgdLGbQpnXY\nQeqPSTRLT71ikH+RqCpoYxcN63i3j9oYWm9KoD5F4arcqlLrEUZ1TqW5csGSY1h6\n2HmGPsINl9KCxwUS+76dxsdHIqLFx0qdnD6S5vmd0sin33jdYhj9ltp0KvhEgMvS\nXMuzECVRvI4MZxebf7gkV3EByqV6XvazBSxuMhplygpAaLra11yV1M/m9wzVwlnS\nS2GNvRkNym9WnH0IQ0kn9hS8hj52Nh12JQ==\n-----END CERTIFICATE-----`
