package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
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
		Providers:     testAccProviders,
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
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${alicloud_esa_site.resource_Site_set_test.id}",
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
					"site_id":      "${alicloud_esa_site.resource_Site_set_test.id}",
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
				Config: testAccConfig(map[string]interface{}{}),
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

resource "alicloud_esa_site" "resource_Site_set_test" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_certificate_set_test.id
  coverage    = "overseas"
  access_type = "NS"
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
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${alicloud_esa_site.resource_Site_cas_test.id}",
					"created_type": "cas",
					"cert_name":    "hyhtest2",
					"cas_id":       "16884280",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cas_id": "16884280",
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


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_cas_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_cas_test" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_cas_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA Certificate. <<< Resource test cases, automatically generated.

const testEsaCertificate = `-----BEGIN CERTIFICATE-----\nMIIDMjCCAhqgAwIBAgIUCSTQpf73Fdt1aWBI4R50edUIma0wDQYJKoZIhvcNAQEF\nBQAwLTEOMAwGA1UEAwwFdGVzdDExDjAMBgNVBAoMBXRlc3QyMQswCQYDVQQGEwJD\nTjAeFw0yNTAxMDIwOTM1MDBaFw0yNTA2MjAwOTM1MDBaMFoxCzAJBgNVBAYTAkNO\nMQ8wDQYDVQQIDAbljJfkuqwxDzANBgNVBAcMBuWMl+S6rDENMAsGA1UECgwEdGVz\ndDEaMBgGA1UEAwwRaHloNC5nb3NpdGVjZG4uY24wggEiMA0GCSqGSIb3DQEBAQUA\nA4IBDwAwggEKAoIBAQDf0y5yMMpCjKv55g/1bwWh+N/En7LTbkMa3DPN9kNwhlib\nkQEHmrWsOCh2P7qaE+UBwXMBxJLXBUVkZr9R7QOydOiQg/RcnYfzMI3GWFsQoQJi\nwsCzhD8qJ6zyi+sRqzp15AqJTkhljrI2YOO12C7b80cza3rEM6gV71sIoUXMGEfX\nZMGenDdgSMTFg67f95JRnazNxtcZS8L6SKjFqDS7py6aO+Tkwk5uM+s2s48nQVKh\nt9AymU0Az5eyG5evNryNBbt4sGSMRaijeU0r7FjfN6/WgG2n0a+UqTmPapm/MMPd\nMQqXoerXFPcWYOA2PW+rWFMR5IrlrZYNsjGJiAcnAgMBAAGjHTAbMAsGA1UdEQQE\nMAKCADAMBgNVHRMBAf8EAjAAMA0GCSqGSIb3DQEBBQUAA4IBAQCoQvwaqc9dQ5vy\nqsqxJP08hLPym41xSnEOXSH6miaiDNJnFxugq9ZSXXyeqYJA+IHmqaWCT4sMlN+J\nWX0SGNUIK/BbCscgORUnHcVajuug7Mrnko67habVslw4o8SXphWKeEiUd4hljtQK\nYFHFkw+1mm4EI97cKRtTCzPf5W4xHxQu1+GYEwbxFMVYKWrPgNTXjZmQfLpbIHCj\n8OqWf9/I4sMxAa71nA4HYUbVsJVt3ecwqlwtIMQS0AIQ2yaROBRbksj+R+80DTw7\nH6uTqH6k0J4s0RE+KWjTTFdG8J7wxziVOko4hYMEbfBrxjlQKEBE/LipWHFtNOtL\n3S4JqWvM\n-----END CERTIFICATE-----`
const testEsaPrivateKey = `-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDf0y5yMMpCjKv5\n5g/1bwWh+N/En7LTbkMa3DPN9kNwhlibkQEHmrWsOCh2P7qaE+UBwXMBxJLXBUVk\nZr9R7QOydOiQg/RcnYfzMI3GWFsQoQJiwsCzhD8qJ6zyi+sRqzp15AqJTkhljrI2\nYOO12C7b80cza3rEM6gV71sIoUXMGEfXZMGenDdgSMTFg67f95JRnazNxtcZS8L6\nSKjFqDS7py6aO+Tkwk5uM+s2s48nQVKht9AymU0Az5eyG5evNryNBbt4sGSMRaij\neU0r7FjfN6/WgG2n0a+UqTmPapm/MMPdMQqXoerXFPcWYOA2PW+rWFMR5IrlrZYN\nsjGJiAcnAgMBAAECggEAHoPlNi5OtQoGvFhQXq7XPsD2iREYyVikD3psGa10flfi\nprr7P/UoxaUWZyDDKRSDSVz9eAU729LdJhYYCWxd76uetW04GJRln5NEEQGk0LyS\n3bIdoZvrHK0yGBNZZhxJKR+BDD0/A5GT3HQTN4yUvuoJEAqcPzmnte3fJGsQYAXC\nFt6mh4whgMEylv43O9Oar3BN9VBdMEt7aeIQICLPSvxGEm5p3rughNc9xe5t0My7\noDmVUxARrEt0p3mcD3kEzdBbYpODP2rDxQ5+1kC5Ehud6brtEWrPRkJ85f/tS4CD\nQa4MDAmLfESvwPx24luyyGZ87zai2akLZnIeNPXZiQKBgQD6k4FGAlND6n3yE77H\n34Dpe/nAN4o5/wFw4PLzofD1Px/Yo/elVj22rkZsAp0nv/09bu874O3GZ3k+MUBw\nQpLZPRB83s+yGxmPGpJ2MEKmNo5yiq+b1ms4LgGoeBJ0Au2Q1D2JgDKNYhAbhXc0\ni2wBOm0Dsmo6vKIYn2tHH5QFfQKBgQDkq3EuXFhHLsYr8B6/Hi38JsgDdiqp3MGw\nz1fihICC916XhlcmSDF3yj/EGuWs2QZbBeF9hc3i5zyX78lf83/61wQUivQKsLh8\nccXAIGC+Anaj+QSR27ffsKOjEf3UcG2GO6fJWQ7wL6YqDTI4NRKqooHHB5ZyazuJ\nXUWbkCTQcwKBgQCTn2cfqasIbhO5FGznMAOwint/BLmIpXVh3QUFB7j/oyrN5Pu3\nCnHdOBsA8yFHE9LL6JlEu6UZqEhDnZyLBo2hMlqOVm4iTdjm+A6lVpVNewK89Hu0\n4cPVGzWa3PJyKm9vTbrbUQ13QfifXif2atU8fAFRlkEDTbJpszHueuonuQKBgEwT\nCU1eJXRRCFbXxLLabHwFvub/6gOm0L1szZUrdcGcYFjSta2juOlcXMh2FlDhxq6D\nj1f+KtfBDDST7o/AZ6Rg4heffr8LueSUyTA41VaBZUz57F5GFa8Sie4XbzW8zYMp\n15Ex0rlvcL1utHVkmrgv5jL7wD6ClEUHxb+SPd0nAoGBANP4TorCH0l0HaE7iVDG\n4fdg6bxwuIvl2gZZFJykz16Lc2v2wow8E9/8pRQp1Emkigi98yco9mUs+ckSqzQU\ncwF9ZO15wpjWbFWoBz3bam59vAL7L4gJ1a2xotBhKLSCOEGakPPF/wO/o6QiIH3u\nL3L439BFPjMU73CD5H+F9Lco\n-----END PRIVATE KEY-----`
