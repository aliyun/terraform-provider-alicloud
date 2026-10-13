package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test ESA Certificate. >>> Resource test cases, automatically generated.
// Case resource_Certificate_apply_test
func TestAccAliCloudESACertificateresource_Certificate_apply_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudESACertificateresource_Certificate_apply_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESACertificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESACertificateresource_Certificate_apply_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domains":      "101.gositecdn.cn",
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"type":         "lets_encrypt",
					"created_type": "free",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"created_type"},
			},
		},
	})
}

var AliCloudESACertificateresource_Certificate_apply_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESACertificateresource_Certificate_apply_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name = "gositecdn.cn"
}

`, name)
}

// Case resource_Certificate_set_test
func TestAccAliCloudESACertificateresource_Certificate_set_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudESACertificateresource_Certificate_set_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESACertificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESACertificateresource_Certificate_set_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"certificate":  testEsaCertificate,
					"private_key":  testEsaPrivateKey,
					"created_type": "upload",
					"region":       "cn-hangzhou",
					"cert_name":    "hyhtestname44",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"created_type": "upload",
					"region":       "cn-beijing",
					"cert_name":    "hyhtestname44",
					"certificate":  testEsaCertificate,
					"private_key":  testEsaPrivateKey,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"created_type", "private_key"},
			},
		},
	})
}

var AliCloudESACertificateresource_Certificate_set_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESACertificateresource_Certificate_set_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_certificate_set_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name = "gositecdn.cn"
}

`, name)
}

// Case resource_Certificate_cas_test
func TestAccAliCloudESACertificateresource_Certificate_cas_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudESACertificateresource_Certificate_cas_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESACertificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESACertificateresource_Certificate_cas_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"created_type": "cas",
					"cert_name":    "hyhtest2",
					"cas_id":       "${local.cert_id}",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"created_type"},
			},
		},
	})
}

var AliCloudESACertificateresource_Certificate_cas_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESACertificateresource_Certificate_cas_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = var.name
  cert = <<EOF
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

locals {
  # cert_id = join("",[alicloud_ssl_certificates_service_certificate.default.id,"cn-hangzhou"])
  cert_id = alicloud_ssl_certificates_service_certificate.default.id
}

resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_cas_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name = "gositecdn.cn"
}

`, name)
}

// Test ESA Certificate. <<< Resource test cases, automatically generated.

const testEsaCertificate = `-----BEGIN CERTIFICATE-----\nMIIDHzCCAgegAwIBAgIUX6GfQdAQ0ou+qmcohIBfYrlj5kgwDQYJKoZIhvcNAQEL\nBQAwODEaMBgGA1UEAwwRaHloNC5nb3NpdGVjZG4uY24xDTALBgNVBAoMBHRlc3Qx\nCzAJBgNVBAYTAkNOMB4XDTI2MDYwMTA2MjU1OVoXDTM2MDUyOTA2MjU1OVowODEa\nMBgGA1UEAwwRaHloNC5nb3NpdGVjZG4uY24xDTALBgNVBAoMBHRlc3QxCzAJBgNV\nBAYTAkNOMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsu3fSD1f6c7x\nWK/iuVCOyUkdcW5dQ6oQ52/51Sz9FiW4z4QLs9rhTM1KK1BB+lNw2xJhu5MJT3+x\n1AxQa6eyVsaiiWZjn++QCpQaIsP6GbdAGvaQcbOMv6eB20ez/q67ev8i0M7ie0+F\nt1scDxCJxHlaULI5JEQaT4gbVm5pX0KARdW7xvIMnw18rGvocD2XbDO9YGMn/Q4l\nqg1u1W8FVbhKkm/p4kl5djMQmoPk8CW/t9uOf5Lm9bcR1z/dXhJEjLz0/ltRMd/w\n5x4qqk73ODKJ+mSSUVUw+7tWO1bSkKqsYW3l/FuZwoOORaP9u/LPDOIiCPPgOUt5\nyLWMC6zrUwIDAQABoyEwHzAdBgNVHQ4EFgQUVFJdgGnh9ukIbV15P8DjNanAutww\nDQYJKoZIhvcNAQELBQADggEBACxBJKrC8EsTT6MQZ4OrP4ms8QSPmYJHbjP17VB6\nBTFYHX7Lww2EB68T4yDwI9ydarDtaXlT0hVyLeWcrguK+tmLtd7h0tRTdO1eE/cj\njs5R9zhbL0bMYPcDr9mjEPfHS8YaIumhitQaeaXv36SSeF5m7IhMHPr5gPGNLj0s\nO4xif4yPP2eR4qiXPHuR5nysYzUW432gocjWdP6Zx++AZaFp33X1COyicFf5uAfZ\nG2uf1Vb6KDQYR6itvMvMzFXrCcOTDxHUuxeJY9PpQumt5vIb9CFG/9ovuV4jrPF8\nBzLfQo6Y5+H+FzXTMwjzJ2rE+Lm70ziOktAA1dvPuny24vI=\n-----END CERTIFICATE-----`
const testEsaPrivateKey = `-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCy7d9IPV/pzvFY\nr+K5UI7JSR1xbl1DqhDnb/nVLP0WJbjPhAuz2uFMzUorUEH6U3DbEmG7kwlPf7HU\nDFBrp7JWxqKJZmOf75AKlBoiw/oZt0Aa9pBxs4y/p4HbR7P+rrt6/yLQzuJ7T4W3\nWxwPEInEeVpQsjkkRBpPiBtWbmlfQoBF1bvG8gyfDXysa+hwPZdsM71gYyf9DiWq\nDW7VbwVVuEqSb+niSXl2MxCag+TwJb+3245/kub1txHXP91eEkSMvPT+W1Ex3/Dn\nHiqqTvc4Mon6ZJJRVTD7u1Y7VtKQqqxhbeX8W5nCg45Fo/278s8M4iII8+A5S3nI\ntYwLrOtTAgMBAAECggEAQyWFT0TJA6MHazLGMKkMjHkFtZWnJkdiBJg+90LkzzTk\nv+tbwOj496tqlAqQV/KMPYoOZyfsrIrNHzhnzZ6nDG2KfWmRJWnvcijWDgnhh1j/\nk57H4gNxZFLmJnYoFAFalfO9CwM4dvIGyiJEy1p3eOXZgMMBMpzkCsiXCb5xiK94\nZRVCK/LvK6eQhNOjYcmUe+mQXhx1hv54T4EIMCcsavBv1MLeKrGdVLv4IHd03QZx\nwWl+VZZh1N9c7Z3sJEFMZVdl3JUWVxxJBP+v8uBoRdgPuXi820zARfNT5CPtcYfF\nPUjyExXo8S3FifPemRBZV+fYjbgpmtcnZl4/SBLMYQKBgQDwGh6dlHCTfD7mK2nv\nqr3e+QXz4I0UbwOWg313r6cuFXxzNSdG5toVZcMvTKh679Jugry5FVg6dIVWXfzA\n75pO0P6h3lU9FVRHBGc2IKtcIx4NUREEDFPM8df7VX4ybi7zTlYs0alzKba0/Zon\nZwtzCUx0PbywhTHD6qKNJkAS+QKBgQC+xtWfzFGvHHMnL4DVwu2jHd2nxwaGGgWn\nQhyNm74XjorhKEC3U3TNpmBMmJPZNNtQSwIQ0Yq5tgn4yybSYTaa6xCrJOf3TBZq\nMJEjxJzumhVS/Q1KkZlorHtEK4a179dQg+IU7r1vEvlfBQf6i5/oU0ORTGJkEkYH\nlVw8w6L3qwKBgCWPMHLeIa4wpXZEHFJNl14l/nRkEC2+IAWPlDUA2VowKkOrcPV6\nb2shfCMODt0MXxLCiNs7J44dZC5ajYtw7+accvjHWvYvO/vQCIVDHwtOwwi6Qbss\nYn+Q5YR/nzosWlPdUUW5lpRZVieB9HdtezEHp1oXvkiuzVYkgkEqVqOhAoGAHLAL\nqabw1ZNCoa7cAcj5MSEplrQv//Rjyz3+yzCTSjmOGsORz7+F/fK54mrDONNg81cE\nLYFFCh4cq8Pox5QEwRD+Ba5cD2zqpfc9rBJBwwN6l2skF4WDeyEMvDiLXkp9p0bd\ntWYdKFnDFA3OoFdkqWvz6iKBXSj+TN+h6iVFGVcCgYAxTS5hKcwZ2m8CvUSDcBup\nn4UB/DAwRF6wND3gsIfYl5wDsU2bCicEUl0XKwNfwfj9WhuqEGmkvP7pDXv/AanP\nsdJU2xeE7HgMg2g25DM3JUAgk8dKnG3TBUL4gV/3PGNBepiV2VEGYgkLZWAIEEb6\nXnkDZruiM1eTYkXgDHZkuA==\n-----END PRIVATE KEY-----`
