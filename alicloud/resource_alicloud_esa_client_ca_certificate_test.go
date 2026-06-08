package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
const testEsaClientCaCertificate = `-----BEGIN CERTIFICATE-----\nMIIDLDCCAhSgAwIBAgIUbW7alHwt3HICHGVsLkqpW9Cd4oIwDQYJKoZIhvcNAQEL\nBQAwLjEQMA4GA1UEAwwHdGVzdC1jYTENMAsGA1UECgwEdGVzdDELMAkGA1UEBhMC\nQ04wHhcNMjYwNjAxMDMyMzEzWhcNMzYwNTI5MDMyMzEzWjAuMRAwDgYDVQQDDAd0\nZXN0LWNhMQ0wCwYDVQQKDAR0ZXN0MQswCQYDVQQGEwJDTjCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAJvUEYwZodySRvwKaHVFgJsGe85tXYDzcHMlBS/7\n904+hVg/zPn3W83vlGHg0S9TUjDDYz+mx9XtGqbVTLNsnAFw6jnPvg6LoK3BE6fN\nT/muS4tkVN7zfKekC6RIPqpDoZEAVfXNzIdg31WdRhfDsUZ7ZgseTsh/o+0ophq/\nz8MwlLd0PYJJNQx+lvmBsBaqXKZ15+TXN4WwZ4Y0euXrX5uCd/eU+KrtsultZb/5\nF1yGR7ON2QCeYe0F1Csi+mdfuniipBtt4evtGRFRkTewzLPSoo8R/IokO50vtfl7\nNTTlwVw+N777PJZF5j0XdjGWtlxL98XpxkugpXSzsWUAewMCAwEAAaNCMEAwDwYD\nVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0OBBYEFIPFTEjW2zDy\nKmmsULjCisAJyu9QMA0GCSqGSIb3DQEBCwUAA4IBAQCSbwv8pOa1Tmv3eUCju4lW\n09JIYAE+H6n3tYRCwao0UdJ69c8yyNVVzbYf/3IbYBPs73OwW8sk3kKlZAHaPcUF\nFsCixLskgWKrJyCcK+sCt0c8JC+GeBHmJcncaq0E7XzxNB5R0KY2w6S7KCjeK1G9\niNsxvAVqNpISuHByd1LNKOT5S1niUx9Yd9eQTCcXEzZnE8Y2tf+gkUPjG9+5kxKL\nXHT3Choi7HjUYmwduKpV7wOAJdlBuyHegBii2xzRhdSriqffbqfczy9D5XSC6Q+w\nyF2ugGWW1IXF10DkMH3ixzbdGTohIUHpRZea3ELKhshMEo6ZqIhCLG3dv1kyeRG+\n-----END CERTIFICATE-----`
