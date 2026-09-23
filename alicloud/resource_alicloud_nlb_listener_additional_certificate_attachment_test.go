package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Nlb ListenerAdditionalCertificateAttachment. >>> Resource test cases, automatically generated.
// Case 5384
func TestAccAliCloudNlbListenerAdditionalCertificateAttachment_basic5384(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_listener_additional_certificate_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbListenerAdditionalCertificateAttachmentMap5384)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbListenerAdditionalCertificateAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlblisteneradditionalcert%d", "cn-hangzhou", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbListenerAdditionalCertificateAttachmentBasicDependence5384)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.NLBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_id": "${alicloud_ssl_certificates_service_certificate.ssl.id}-cn-hangzhou",
					"listener_id":    "${alicloud_nlb_listener.create_listener.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_id": CHECKSET,
						"listener_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudNlbListenerAdditionalCertificateAttachmentMap5384 = map[string]string{
	"status": CHECKSET,
}

func AlicloudNlbListenerAdditionalCertificateAttachmentBasicDependence5384(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

resource "alicloud_vpc" "create_vpc" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "create_vsw_j" {
  vpc_id     = alicloud_vpc.create_vpc.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "create_vsw_k" {
  vpc_id     = alicloud_vpc.create_vpc.id
  zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_nlb_load_balancer" "lb" {
  address_ip_version = "Ipv4"
  zone_mappings {
    vswitch_id = alicloud_vswitch.create_vsw_j.id
    zone_id    = alicloud_vswitch.create_vsw_j.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.create_vsw_k.id
    zone_id    = alicloud_vswitch.create_vsw_k.zone_id
  }
  load_balancer_type = "Network"
  load_balancer_name = var.name

  vpc_id       = alicloud_vpc.create_vpc.id
  address_type = "Internet"
}

resource "alicloud_nlb_server_group" "create_sg" {
  address_ip_version = "Ipv4"
  scheduler          = "Wrr"
  health_check {
  }
  server_group_type = "Instance"
  vpc_id            = alicloud_vpc.create_vpc.id
  protocol          = "TCPSSL"
  server_group_name = var.name

}

resource "alicloud_ssl_certificates_service_certificate" "ssl0" {
  cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIIDhDCCAmwCCQCwJW4JChLBqTANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC
Q04xEDAOBgNVBAgMB0JlaWppbmcxEDAOBgNVBAcMB0JlaWppbmcxDDAKBgNVBAoM
A0FsaTEPMA0GA1UECwwGQWxpeXVuMRIwEAYDVQQDDAlUZXJyYWZvcm0xHTAbBgkq
hkiG9w0BCQEWDjEyM0BhbGl5dW0uY29tMB4XDTI0MTIyNTA3MjQ0OFoXDTI3MTIy
NTA3MjQ0OFowgYMxCzAJBgNVBAYTAkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYD
VQQHDAdCZWlqaW5nMQwwCgYDVQQKDANBbGkxDzANBgNVBAsMBkFsaXl1bjESMBAG
A1UEAwwJVGVycmFmb3JtMR0wGwYJKoZIhvcNAQkBFg4xMjNAYWxpeXVtLmNvbTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK4UufXydtJZeW6lX9VahVIk
ifblYCVkFcFoderF2FtD5AeMZJ+v+chHc7RiV+U7P3o0Fzk+cg7OL9dSEYBrwHK4
9yCwU/Mv+I/KsS8GjrRMOPjbrYvI0GjheEPJcILbt29tygrxX2PwV6FqWNknbGpk
Ej8L9zTL977IHBmgw8A2AeKlqV64s8ydAgGbWO3NTK64OlEJJNR+J+75uYskNT3s
8DqaQV/IWlGAiUmGVeorWkrAWCfx2zSwI9jU8pNHtSF7PyxlbRy1ir2Lv1WnQKHf
Bnhr/wXwKOL5IJRVZ144Z9TdcoPo4GbFmYMSTwYFIbjYZ3yxoygeXMk0UXPZxVsC
AwEAATANBgkqhkiG9w0BAQsFAAOCAQEAVPA+Q0/5T6VzVw+MFXjCxXH1mWgd767w
YWX4tvdGsTDkK6/ESm8m9GDp5F3p7Degk0isr9XkyzkWo/nPEPWQOeYR0kNTvpwY
mKz9/aJwxalHS6O/8K2Ed6pZcXW0SUfjdH0/9YHw+vu4i2cQGTICzrKuEvyck40y
fQocvFyw6O7W+tewLA4ntTuC6HhEQbC0p7zxGc3LSuayBgJrJiOAnGvFu+/OFQi+
zEXi1xt8uQR6q5DQDsfqNCxpRKsCmU+POzNg2Y31GDMv4ZPerou5jXa1gh8/TVBT
IX3OTy5aL4Ue8nBip3bVw+V/9L9xhmXbex6IMwwvrWI4OfMt6ECifQ==
-----END CERTIFICATE-----
EOF
  certificate_name = join("-", [var.name, 0])

  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEArhS59fJ20ll5bqVf1VqFUiSJ9uVgJWQVwWh16sXYW0PkB4xk
n6/5yEdztGJX5Ts/ejQXOT5yDs4v11IRgGvAcrj3ILBT8y/4j8qxLwaOtEw4+Nut
i8jQaOF4Q8lwgtu3b23KCvFfY/BXoWpY2SdsamQSPwv3NMv3vsgcGaDDwDYB4qWp
XrizzJ0CAZtY7c1Mrrg6UQkk1H4n7vm5iyQ1PezwOppBX8haUYCJSYZV6itaSsBY
J/HbNLAj2NTyk0e1IXs/LGVtHLWKvYu/VadAod8GeGv/BfAo4vkglFVnXjhn1N1y
g+jgZsWZgxJPBgUhuNhnfLGjKB5cyTRRc9nFWwIDAQABAoIBAC0D2Q6bc1RzpK4S
/5QZQ055el+o8tLYbbPEwnFCVe9LwASfrkmI5OuAZpAnuhjh2ElOfQ7lcfMYKFDi
vPnbYzmHUQhX8G17YygzvtutM2u2JilcDSWPeS0V2NaWmYyNKoMa/dsUjZk3RkHM
UUteIW/ljr5U5sj1UYw5DOMnqlbicy2cPPY4g1QKGW5t3p5Lxw5ojgqynzi8EKMq
j0apEoTXxmciOrwwiP2ynRTEN77+FUZkTvmxmSPoIfNTycDPRr4aUwVHV2d5FHPn
d1MdjSoUPbHdOLfynwXqTz9OlvMSUDrBvs6k5ripGY9qvh9PrOdj7zLXVRQXUuOR
YwoVHKECgYEA1NzNGifjW5cdcbkzc86QA9TM7yAyBmgnopzlm+dFIhxtJmydxN4V
820x1Lkfe2vLCyYQ6fcEKAtjC9qdw+E2uzHAbtvnR4JseF3z1D82xw3MgGT9l3zc
mMrgKmdCGGLWi6hboylX+2GBMVl2R0aRZrZje67jZcDXd06mlvW257ECgYEA0Vv1
X3Ubn8XA3AA1ybem8fWNwEXfcYP1lJq0cX1gUXlpQvxWN61//aFZUCJZw5cEPArQ
rEqhT81VCqXGO/by6D3fJD+4P8D6v6szJK2AGvXkZMfnJwAbHcOyGlxMin1CTJss
ZID0XI9xmbedm7Wi40+qXz8q4rs25kft9YjfzMsCgYBSPfE8vtaYJ52nt7+Kae+4
mzqG1XCeixVtPaN1BfjvAf6mDucyDgB7KeBL6S6ht/ceGpoEW30On7+n79JuwRAt
aT6JVotYVKrmIp63jajzZYByxxI3unVcz12m5HhkBaQRF344XxvwMy8ASyloxnod
LjDns52GTeix3wB8aPk/MQKBgGOQRwXpjISUKB64HtxacZN6ArqgwB2c8uqEFDIw
vOCiS7Pmix4ZbdfxpqbcXzIMHKBtSEXXjBWGgd35bmfQDj7yRa9Yekgff2Ati7ny
pQytSbu/8abzfvHNwmKU6HWoEiKaXSdCyHNIaG8BCnwlilxt44k+YifHftlO9dSi
DkS3AoGAYmF++8uEvQot5Yma4GraY+7ZyfWNLwClsOsrN2g19Vycg16fJk0olwDx
2kRWKqNn99HJJwiLie1nvsDRJLbmzmI91Pttpu/EYFDJ8OYQOr1OhhPwwTygf+7S
1o2RTXu3gKNG6fxOtHFatws3IzvovOASYyJR5XP2sIJURLOrSN0=
-----END RSA PRIVATE KEY-----
EOF
}

resource "alicloud_nlb_listener" "create_listener" {
  listener_port      = "443"
  server_group_id    = alicloud_nlb_server_group.create_sg.id
  load_balancer_id   = alicloud_nlb_load_balancer.lb.id
  listener_protocol  = "TCPSSL"
  certificate_ids    = ["${alicloud_ssl_certificates_service_certificate.ssl0.id}-cn-hangzhou"]
  ca_certificate_ids = []
}

resource "alicloud_ssl_certificates_service_certificate" "ssl" {
  cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIIDhDCCAmwCCQCwJW4JChLBqTANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC
Q04xEDAOBgNVBAgMB0JlaWppbmcxEDAOBgNVBAcMB0JlaWppbmcxDDAKBgNVBAoM
A0FsaTEPMA0GA1UECwwGQWxpeXVuMRIwEAYDVQQDDAlUZXJyYWZvcm0xHTAbBgkq
hkiG9w0BCQEWDjEyM0BhbGl5dW0uY29tMB4XDTI0MTIyNTA3MjQ0OFoXDTI3MTIy
NTA3MjQ0OFowgYMxCzAJBgNVBAYTAkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYD
VQQHDAdCZWlqaW5nMQwwCgYDVQQKDANBbGkxDzANBgNVBAsMBkFsaXl1bjESMBAG
A1UEAwwJVGVycmFmb3JtMR0wGwYJKoZIhvcNAQkBFg4xMjNAYWxpeXVtLmNvbTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK4UufXydtJZeW6lX9VahVIk
ifblYCVkFcFoderF2FtD5AeMZJ+v+chHc7RiV+U7P3o0Fzk+cg7OL9dSEYBrwHK4
9yCwU/Mv+I/KsS8GjrRMOPjbrYvI0GjheEPJcILbt29tygrxX2PwV6FqWNknbGpk
Ej8L9zTL977IHBmgw8A2AeKlqV64s8ydAgGbWO3NTK64OlEJJNR+J+75uYskNT3s
8DqaQV/IWlGAiUmGVeorWkrAWCfx2zSwI9jU8pNHtSF7PyxlbRy1ir2Lv1WnQKHf
Bnhr/wXwKOL5IJRVZ144Z9TdcoPo4GbFmYMSTwYFIbjYZ3yxoygeXMk0UXPZxVsC
AwEAATANBgkqhkiG9w0BAQsFAAOCAQEAVPA+Q0/5T6VzVw+MFXjCxXH1mWgd767w
YWX4tvdGsTDkK6/ESm8m9GDp5F3p7Degk0isr9XkyzkWo/nPEPWQOeYR0kNTvpwY
mKz9/aJwxalHS6O/8K2Ed6pZcXW0SUfjdH0/9YHw+vu4i2cQGTICzrKuEvyck40y
fQocvFyw6O7W+tewLA4ntTuC6HhEQbC0p7zxGc3LSuayBgJrJiOAnGvFu+/OFQi+
zEXi1xt8uQR6q5DQDsfqNCxpRKsCmU+POzNg2Y31GDMv4ZPerou5jXa1gh8/TVBT
IX3OTy5aL4Ue8nBip3bVw+V/9L9xhmXbex6IMwwvrWI4OfMt6ECifQ==
-----END CERTIFICATE-----
EOF
  certificate_name = join("-", [var.name, 1])

  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEArhS59fJ20ll5bqVf1VqFUiSJ9uVgJWQVwWh16sXYW0PkB4xk
n6/5yEdztGJX5Ts/ejQXOT5yDs4v11IRgGvAcrj3ILBT8y/4j8qxLwaOtEw4+Nut
i8jQaOF4Q8lwgtu3b23KCvFfY/BXoWpY2SdsamQSPwv3NMv3vsgcGaDDwDYB4qWp
XrizzJ0CAZtY7c1Mrrg6UQkk1H4n7vm5iyQ1PezwOppBX8haUYCJSYZV6itaSsBY
J/HbNLAj2NTyk0e1IXs/LGVtHLWKvYu/VadAod8GeGv/BfAo4vkglFVnXjhn1N1y
g+jgZsWZgxJPBgUhuNhnfLGjKB5cyTRRc9nFWwIDAQABAoIBAC0D2Q6bc1RzpK4S
/5QZQ055el+o8tLYbbPEwnFCVe9LwASfrkmI5OuAZpAnuhjh2ElOfQ7lcfMYKFDi
vPnbYzmHUQhX8G17YygzvtutM2u2JilcDSWPeS0V2NaWmYyNKoMa/dsUjZk3RkHM
UUteIW/ljr5U5sj1UYw5DOMnqlbicy2cPPY4g1QKGW5t3p5Lxw5ojgqynzi8EKMq
j0apEoTXxmciOrwwiP2ynRTEN77+FUZkTvmxmSPoIfNTycDPRr4aUwVHV2d5FHPn
d1MdjSoUPbHdOLfynwXqTz9OlvMSUDrBvs6k5ripGY9qvh9PrOdj7zLXVRQXUuOR
YwoVHKECgYEA1NzNGifjW5cdcbkzc86QA9TM7yAyBmgnopzlm+dFIhxtJmydxN4V
820x1Lkfe2vLCyYQ6fcEKAtjC9qdw+E2uzHAbtvnR4JseF3z1D82xw3MgGT9l3zc
mMrgKmdCGGLWi6hboylX+2GBMVl2R0aRZrZje67jZcDXd06mlvW257ECgYEA0Vv1
X3Ubn8XA3AA1ybem8fWNwEXfcYP1lJq0cX1gUXlpQvxWN61//aFZUCJZw5cEPArQ
rEqhT81VCqXGO/by6D3fJD+4P8D6v6szJK2AGvXkZMfnJwAbHcOyGlxMin1CTJss
ZID0XI9xmbedm7Wi40+qXz8q4rs25kft9YjfzMsCgYBSPfE8vtaYJ52nt7+Kae+4
mzqG1XCeixVtPaN1BfjvAf6mDucyDgB7KeBL6S6ht/ceGpoEW30On7+n79JuwRAt
aT6JVotYVKrmIp63jajzZYByxxI3unVcz12m5HhkBaQRF344XxvwMy8ASyloxnod
LjDns52GTeix3wB8aPk/MQKBgGOQRwXpjISUKB64HtxacZN6ArqgwB2c8uqEFDIw
vOCiS7Pmix4ZbdfxpqbcXzIMHKBtSEXXjBWGgd35bmfQDj7yRa9Yekgff2Ati7ny
pQytSbu/8abzfvHNwmKU6HWoEiKaXSdCyHNIaG8BCnwlilxt44k+YifHftlO9dSi
DkS3AoGAYmF++8uEvQot5Yma4GraY+7ZyfWNLwClsOsrN2g19Vycg16fJk0olwDx
2kRWKqNn99HJJwiLie1nvsDRJLbmzmI91Pttpu/EYFDJ8OYQOr1OhhPwwTygf+7S
1o2RTXu3gKNG6fxOtHFatws3IzvovOASYyJR5XP2sIJURLOrSN0=
-----END RSA PRIVATE KEY-----
EOF
}


`, name)
}

// Test Nlb ListenerAdditionalCertificateAttachment. <<< Resource test cases, automatically generated.
