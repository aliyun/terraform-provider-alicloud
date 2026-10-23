package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudVodDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_domain.default"
	ra := resourceAttrInit(resourceId, AliCloudVodDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VodService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVodDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVodDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VodSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": name,
					"scope":       "overseas",
					"sources": []map[string]interface{}{
						{
							"source_type":    "ipaddr",
							"source_content": "1.1.1.1",
							"source_port":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": name,
						"scope":       "overseas",
						"sources.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "on",
					"cert_id":      "${alicloud_ssl_certificates_service_certificate.default.0.id}",
					"cert_region":  "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "on",
						"cert_id":      CHECKSET,
						"cert_region":  "cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cert_id":   "${alicloud_ssl_certificates_service_certificate.default.1.id}",
					"cert_name": "${alicloud_ssl_certificates_service_certificate.default.1.certificate_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert_id":   CHECKSET,
						"cert_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "off",
					"cert_id":      REMOVEKEY,
					"cert_name":    REMOVEKEY,
					"cert_region":  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "off",
						"cert_id":      REMOVEKEY,
						"cert_name":    REMOVEKEY,
						"cert_region":  REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "on",
					"cert_name":    name,
					"cert_type":    "upload",
					"ssl_pub":      "${var.cert}",
					"ssl_pri":      "${var.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "on",
						"cert_name":    name,
						"cert_type":    "upload",
						"ssl_pub":      CHECKSET,
						"ssl_pri":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cert_name": name + "update",
					"ssl_pub":   "${var.cert_update}",
					"ssl_pri":   "${var.key_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert_name": name + "update",
						"ssl_pub":   CHECKSET,
						"ssl_pri":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"source_type":     "domain",
							"source_content":  "www.aliyuntest.com",
							"source_port":     "443",
							"source_priority": "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "1",
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
				ImportStateVerifyIgnore: []string{"ssl_pri", "check_url", "top_level_domain", "tags"},
			},
		},
	})
}

func TestAccAliCloudVodDomain_basic0_twin0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_domain.default"
	ra := resourceAttrInit(resourceId, AliCloudVodDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VodService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVodDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVodDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VodSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":  name,
					"scope":        "overseas",
					"ssl_protocol": "on",
					"cert_id":      "${alicloud_ssl_certificates_service_certificate.default.0.id}",
					"cert_name":    "${alicloud_ssl_certificates_service_certificate.default.0.certificate_name}",
					"cert_type":    "cas",
					"cert_region":  "cn-hangzhou",
					"check_url":    "http://" + name + "/test.html",
					"sources": []map[string]interface{}{
						{
							"source_type":     "ipaddr",
							"source_content":  "1.1.1.1",
							"source_port":     "80",
							"source_priority": "30",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":  name,
						"scope":        "overseas",
						"ssl_protocol": "on",
						"cert_id":      CHECKSET,
						"cert_name":    CHECKSET,
						"cert_type":    "cas",
						"cert_region":  "cn-hangzhou",
						"sources.#":    "1",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_pri", "check_url", "top_level_domain", "tags"},
			},
		},
	})
}

func TestAccAliCloudVodDomain_basic0_twin1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_domain.default"
	ra := resourceAttrInit(resourceId, AliCloudVodDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VodService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVodDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVodDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VodSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":  name,
					"scope":        "overseas",
					"ssl_protocol": "on",
					"cert_name":    name,
					"cert_type":    "upload",
					"ssl_pub":      "${var.cert}",
					"ssl_pri":      "${var.key}",
					"check_url":    "http://" + name + "/test.html",
					"sources": []map[string]interface{}{
						{
							"source_type":     "ipaddr",
							"source_content":  "1.1.1.1",
							"source_port":     "80",
							"source_priority": "30",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":  name,
						"scope":        "overseas",
						"ssl_protocol": "on",
						"cert_name":    name,
						"cert_type":    "upload",
						"ssl_pub":      CHECKSET,
						"ssl_pri":      CHECKSET,
						"sources.#":    "1",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_pri", "check_url", "top_level_domain", "tags"},
			},
		},
	})
}

var AliCloudVodDomainMap0 = map[string]string{
	"ssl_protocol": CHECKSET,
	"cname":        CHECKSET,
	"description":  "",
	"status":       CHECKSET,
	"gmt_created":  CHECKSET,
	"gmt_modified": CHECKSET,
}

func AliCloudVodDomainBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "cert_ssl" {
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

variable "key_ssl" {
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

variable "cert" {
  default = <<EOF
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
}

variable "key" {
  default = <<EOF
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

variable "cert_update" {
  default = <<EOF
-----BEGIN CERTIFICATE-----
MIIC1TCCAb2gAwIBAgIJALlQrZ+xRM5/MA0GCSqGSIb3DQEBCwUAMBoxGDAWBgNV
BAMMD3dvMi5leGFtcGxlLmNvbTAeFw0yNjA5MTcwNzEzMDJaFw0zNjA5MTQwNzEz
MDJaMBoxGDAWBgNVBAMMD3dvMi5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAMgnNMtWdGsE0PGSIZPAKdNETllbh61BA9nBXBS3xrhQ
jP/gESmyNJ0acROKB7xTH/uvrm0zQtO4Hh86zowhY4VArmFkMsU0/Gev30klWm0Q
QrCu8bkvM2YO46mb/Jv8VkAWQwDfiZ61W7YSUsod4pr7Who3E/2Fi0yiG3MIZHOL
hXYHpspuxM1dIINIFJ3L0XTjD5CBkaMczZpHqOlWXOYBfuzDJeVQOnvynmfGFB95
VbwSUko0wdVCq70HcR15hM+6UiT4i2K9LG2gvG4Hjh77Zt9r2Zeuk0nopXCznZVf
TOsnxrFfbXh1Wk/kdGy/zU6QqVZfUUamiutqstIoOjkCAwEAAaMeMBwwGgYDVR0R
BBMwEYIPd28yLmV4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQB4YxefozgB
lNtiDGR4BGviXVGLPWXesRgJ8dG5oUnQ2T0vCiy847YLKUY1aKSgC/lftGtNX62M
Xg5Hw8TuYYILJyfgBy0hr70RwZK8gpM0CnlkcVfsuIA8no0D4r5oFmAMQsG0QcWl
cUyuULlkOQ1XZ9uoQkq/c9XljyUjO9kBa+OFS5MDLl94bwrmZV6vuUQJaUyQn5OT
QN+qcvE642vecqtDkJ0jFC2GPQTZjA0pwWcg5ijzQydBIr5bcIkQnPGQ2C8pflqK
KgapaB2rHwkFuUZs+rUmq+31C6z17Fep0hbGanjYpSdNz4rNKzMnHmqayvhnlwyR
+ZSbN9AmkVBd
-----END CERTIFICATE-----
EOF
}

variable "key_update" {
  default = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDIJzTLVnRrBNDx
kiGTwCnTRE5ZW4etQQPZwVwUt8a4UIz/4BEpsjSdGnETige8Ux/7r65tM0LTuB4f
Os6MIWOFQK5hZDLFNPxnr99JJVptEEKwrvG5LzNmDuOpm/yb/FZAFkMA34metVu2
ElLKHeKa+1oaNxP9hYtMohtzCGRzi4V2B6bKbsTNXSCDSBSdy9F04w+QgZGjHM2a
R6jpVlzmAX7swyXlUDp78p5nxhQfeVW8ElJKNMHVQqu9B3EdeYTPulIk+ItivSxt
oLxuB44e+2bfa9mXrpNJ6KVws52VX0zrJ8axX214dVpP5HRsv81OkKlWX1FGporr
arLSKDo5AgMBAAECggEAZptawEbw5ZARgVXVqsXfazPgKoT2uAXVcINkriZMVF7O
+bmncsWtnUBZ6ZdgQdYnVkKI2+hkgN8paf8etc3EFCLFgp34PQrdvlX11Ky2f6TP
W35D4xhkdBR7945E4L5qhO4kqUTvaxI0EtGv6XeKnMLwPPIxaE1WzRFpmv4ZaPWc
bYVLQVzt3ZJf6/ZARvGJP8l9EUfLGTgGzBkHAqj09yG/cRS50vKVzpjAK55NqP1t
ElOqy6OygBci8zu90cfqJfFAAcxqybwtOYAVozW8BMhwoVO1pl9vE6p/QTPksAg4
RtZSN1RwFWNXFgX+rBhIK9kOIDEjE+OI7tY8lfabpQKBgQDnyUOZBZGZjgxcchQs
Dynv8WTgBdUkGSQyFhskjpyFnU5ogzVj4mYhsD9IuAdbzheEInyUkqhm2si71xeu
zWN6Ms3jPfSCJPqGVz19Bc5xF11SAexdqzop3JD3iGkKRAzoq73FXnOtyxz+KPAT
u76m2Qjj5APAJguOziZ8vKgVxwKBgQDdD/gdDJZH5rv8qAa+lhIUdtyE4HUe+r39
IbaOIBHMXzH3hfOBOx8/D4MI9+CaURl6AWP14TGmKfUl1d5wYaypQb4smk/2U/MC
8aSb0SLkfMhWXb43rFGf1MqKKmBG+9ogh9+GGa8jZdiUxYc0XeV/B6MM05Fzmcju
4jZXE20v/wKBgBydmW3b5STbB6WmI91opdnMqJKLxcumGvTgC8vQm3xXdmWqxT2m
pstD8hwtNCNeZ0LV6yGAla/VFPEz1eWrz7s3aHnEnVov7IoHXT0BKYKGyuuj9yV9
TTn71B3rRmSREDLteQlT5t4u9E5giXs2M4xvwQw5i0j+qXCg/NlrhSjjAoGAaTr5
0/Iy84UthYcCtuS/NIikToh+mEGqAP50ZVF0x4GlSnu0P5dCN4WNzc4ZV7w1UXOp
/ti1BrFAZkNZgMuuImAORDMsz0iFDKi0wAuc/YE86If1TZE0xc3IoF6tt/yZiBTF
3ZMKYOJRZ/hjdRUhcwfSoP9xZMMW+Egsr+SqfEcCgYAcr0u6Jx3lYVWlu18Xp67l
PkqHk8u5dG5tRse63yuofm/Ab4MyH1Qpvsjb+nS6VAmQDZo/DZrmtyFjvwyMMhza
mtWyYjxLCN4kqMYPfvLWacI10BH7oPI8GhYu4ODMG/jNfn7Znx/iOWu1ZZeKa5mK
5yBbAEwZ8BcG7vL6Ji4qrQ==
-----END PRIVATE KEY-----
EOF
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  count            = 2
  certificate_name = join("", [var.name, count.index])
  cert             = var.cert_ssl
  key              = var.key_ssl
}
`, name)
}
