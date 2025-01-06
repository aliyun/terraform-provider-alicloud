package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudAlbAscript_basic2051(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_ascript.default"
	ra := resourceAttrInit(resourceId, AliCloudAlbAscriptMap2051)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbAscript")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testaccalbascript%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlbAscriptBasicDependence2051)
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
					"listener_id":    "${alicloud_alb_listener.default.id}",
					"position":       "RequestHead",
					"ascript_name":   "test",
					"script_content": "time()",
					"enabled":        "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_id":    CHECKSET,
						"position":       "RequestHead",
						"ascript_name":   "test",
						"script_content": "time()",
						"enabled":        "true",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"ascript_name": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ascript_name": "test_update",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"script_content": "ls",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"script_content": "ls",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"script_content":        "now()",
					"ascript_name":          "test",
					"enabled":               "true",
					"ext_attribute_enabled": "true",
					"ext_attributes": []map[string]interface{}{
						{
							"attribute_key":   "EsDebug",
							"attribute_value": "rdk",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"script_content":        "now()",
						"ascript_name":          "test",
						"enabled":               "true",
						"ext_attribute_enabled": "true",
						"ext_attributes.#":      "1",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"ext_attribute_enabled": "false",
					"ext_attributes":        REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ext_attribute_enabled": "false",
						"ext_attributes.#":      "0",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudAlbAscriptMap2051 = map[string]string{}

func AliCloudAlbAscriptBasicDependence2051(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "port" {
  default = "3366"
}

data "alicloud_alb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count        = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_alb_zones.default.zones.0.id
  vswitch_name = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count        = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id      = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_alb_load_balancer" "default_3" {
  vpc_id                 = data.alicloud_vpcs.default.ids.0
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Standard"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  tags = {
    Created = "TF"
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  modification_protection_config {
    status = "NonProtection"
  }
  access_log_config {
    log_project = alicloud_log_project.default.name
    log_store   = alicloud_log_store.default.name
  }
}

resource "alicloud_alb_server_group" "default_4" {
  protocol          = "HTTP"
  vpc_id            = data.alicloud_vpcs.default.vpcs.0.id
  server_group_name = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  health_check_config {
    health_check_enabled = "false"
  }
  sticky_session_config {
    sticky_session_enabled = "false"
  }
  tags = {
    Created = "TF"
  }
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = var.name
  cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID1zCCAr+gAwIBAgIRAOrWWz1qmkcSg90JDHjuzFwwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjQxMTI2MDczNjA4WhcNMjkxMTI1MDczNjA4WjAgMQswCQYDVQQGEwJDTjER
MA8GA1UEAxMIdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQDa7HDGbQ1Km0f4ZaFzYbjVN0q8KkvZ+oQUd4naGOZnlH5k0XFwmjg+TWf88YX3
5IF8c45/rXrTWucPLg7FeqR96Wq9HZEmzEhs6VG031V9Hqa32saRScCOAyhiW7Hj
OWf6BZveuxbZNbgQCR59QzX4CeAIC68xavIDAy3wcTAH9cIkD71BxEPJGGR7BIVH
9DcWXaMAnJqQfrkth0xHBjflZABHAI0wPYPfaw8fd9DRkMYOIkfjwrrcL5IvhI1u
D3wdHJQWA2vR8hjoU4dHiJLbUtQ+xV1UGVkF67CpQ6LDjSQdX7xlZ7WJMc/7dCJ9
a7tr0ZTwq4/3KSgcRvm62oGvAgMBAAGjgc0wgcowDgYDVR0PAQH/BAQDAgWgMB0G
A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQogSYF0TQa
P8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYBBQUHMAGGFWh0dHA6
Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDovL2NhLm15c3NsLmNv
bS9teXNzbHRlc3Ryc2EuY3J0MBMGA1UdEQQMMAqCCHRlc3QuY29tMA0GCSqGSIb3
DQEBCwUAA4IBAQAxPOlK5WBA9kITzxYyjqe/YvWzfMlsmj0yvpyHrPeZf7HZTTFz
ebYkzrHL8ZLyOHBhag0nL7Poj6ek98NoXTuCYCi8LspdadapOeYQzLce3beu/frk
sqU0A6WLHG9Ol9yUDMCX7xvLoAY/LDrcOM3Z87C/u/ykB4wKfFN2XfR3EZx3PQqw
sV77LOnyQixB4FMHpHlKuDoUkSN9uvxwEPOeGnLZXm96hPsjPwk1bDM8qerNPpVI
CwJ6kNuZ2eLz2Umqu2Gh3l4aADdIwxRY1OOjjZNut8STosABKWVGIwQbbAdRPQze
qHZ05oVTjFy9L1DAzhQ5Zn3oUjLl5KW4tYBA
-----END CERTIFICATE-----
EOF
  key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2uxwxm0NSptH+GWhc2G41TdKvCpL2fqEFHeJ2hjmZ5R+ZNFx
cJo4Pk1n/PGF9+SBfHOOf61601rnDy4OxXqkfelqvR2RJsxIbOlRtN9VfR6mt9rG
kUnAjgMoYlux4zln+gWb3rsW2TW4EAkefUM1+AngCAuvMWryAwMt8HEwB/XCJA+9
QcRDyRhkewSFR/Q3Fl2jAJyakH65LYdMRwY35WQARwCNMD2D32sPH3fQ0ZDGDiJH
48K63C+SL4SNbg98HRyUFgNr0fIY6FOHR4iS21LUPsVdVBlZBeuwqUOiw40kHV+8
ZWe1iTHP+3QifWu7a9GU8KuP9ykoHEb5utqBrwIDAQABAoIBAQCErEfIKOymKybZ
pZXLnAxswt563FMtngGPecZEM1TmrvpOVROffwbY0wZTJ3fd/FBwwIM6Y0MNdYiU
DYCMM0AewmeahqGh1qmJv3hx2eswMXQt9driz8RvDADcYt+SagbWYbHNsKovJrwO
k8gzd5jsYeewWIxqsXpLUxDzJ1VJbIqoHgkrirRRPo0onpixPWeA0RbElSwjwIUw
y43cC4WF8N7wot3cTST8yeKM8ujtqpN22ZtKnbkHTd03vnwQTMeUMJeDQmSmY5aJ
yFr7yw/Z66+7Amh6pkWhzZSDHsjI4y/S3CCdpwFlMA7ID590umJB6HFxWsmVacSe
MSs2vIJZAoGBAOiecPH1HVDQqH6PcrN/X9E3pDKSyAj+nHsVDGIZsie9f5g/qA0A
tcJtQLS0CzrpMTLsAnsfdh2T7Lg6pYFz5jnOUyMjOImAEbCtgvqBxqgFea//OhdP
8s/RmxKIAenBsk7Wbwx8/KPhbZLUNe8OnILVHDfS6kLSa49Iu+4UvrpNAoGBAPDt
mky5MMHKdHwbqxPo9jYrz1m3gqqIvv+VihO4t/DE6t2Zg43ctfFm1BVEDSwPjYs/
YV69KfVrVRUnzMZVdtHZ/dBK784YTY0OujemoaIzMKFIL8tbJFldVv2IgB+IelTX
e675hVdHjNUqZhHwccd8X6d/8icohZw62SNHb/HrAoGBAN1HSt1/c6Gau42Y212Q
fw9ARLuvEQYtXaFfxmXTV7uh8axccXndAQmwb+r1kfE6PojYJQwGQ4+jVX1ynFnm
bEz0zfUQ3gk+gJV2mK+/n7/ZZYZb3WCrtqimFUOtiVRZ40pHhV91zcX+/QK9R4je
d1elbbBUvG9QRu0IHW0+4qfJAoGAOmlQvIM1l/ZOsXw/yO71KoMKnXTJYDERJYQK
2ucw6VXEn39FjtJQ5jsI9jLugp0usvDl2YNBNfgUw7FHi1pTGWOhjqtsYmov+x/z
8+QZUerZQnDu7X2mXWgs3AEJFxwOlJ09pllmg5ecRF4oKvdBjpzP0BtMCURgyFTY
Kh56vIsCgYBMbneMvFY6PCESKIAXj16BF4lqYVXFqHVoxyfxIuVlAy3TMNwxvpbS
yDETk05Ux9yNES0WyTb1SWVG1o1wXc0dnDXCwJqLC1tzJUNUSD1AYvktoNIFErcN
gs3ercrzBTX5ezORPj9ErRAPrSq+V3z1Lge5Gl+EqgDvAfnknww75w==
-----END RSA PRIVATE KEY-----
EOF
}

resource "alicloud_alb_listener" "default" {
  load_balancer_id     = alicloud_alb_load_balancer.default_3.id
  listener_protocol    = "HTTPS"
  listener_port        = var.port
  listener_description = var.name
  default_actions {
    type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default_4.id
      }
    }
  }
  certificates {
    certificate_id = join("", [alicloud_ssl_certificates_service_certificate.default.id, "-cn-hangzhou"])
  }
}
`, name)
}

// Test Alb AScript. >>> Resource test cases, automatically generated.
// Case Ascript_test241227_RequestFoot 9755
func TestAccAliCloudAlbAScript_basic9755(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_ascript.default"
	ra := resourceAttrInit(resourceId, AlicloudAlbAScriptMap9755)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbAScript")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbAScriptBasicDependence9755)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"script_content":        "if and(match_re($uri, '^/1.txt$'), $arg_type) { rewrite(concat('/1.', $arg_type), 'break') }",
					"position":              "RequestFoot",
					"ascript_name":          name,
					"listener_id":           "${alicloud_alb_listener.defaultc8auOh.id}",
					"enabled":               "true",
					"ext_attribute_enabled": "true",
					"ext_attributes": []map[string]interface{}{
						{
							"attribute_key":   "EsDebug",
							"attribute_value": "aaa",
						},
					},
					"dry_run": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"script_content":        "if and(match_re($uri, '^/1.txt$'), $arg_type) { rewrite(concat('/1.', $arg_type), 'break') }",
						"position":              "RequestFoot",
						"ascript_name":          name,
						"listener_id":           CHECKSET,
						"enabled":               "true",
						"ext_attribute_enabled": "true",
						"ext_attributes.#":      "1",
						"dry_run":               "false",
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

var AlicloudAlbAScriptMap9755 = map[string]string{
	"status": CHECKSET,
}

func AlicloudAlbAScriptBasicDependence9755(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zoneid-1" {
  default = "cn-hangzhou-j"
}

variable "zoneid-2" {
  default = "cn-hangzhou-k"
}

variable "regionid" {
  default = "cn-hangzhou"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vsw-1" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = var.zoneid-1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "vsw-2" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = var.zoneid-2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%%s2", var.name)
}

resource "alicloud_alb_load_balancer" "lb" {
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw-1.id
    zone_id    = alicloud_vswitch.vsw-1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw-2.id
    zone_id    = alicloud_vswitch.vsw-2.zone_id
  }
}

resource "alicloud_alb_server_group" "defaultBJ6wNY" {
  vpc_id    = alicloud_vpc.vpc.id
  scheduler = "Wrr"
  sticky_session_config {
    cookie_timeout      = "1000"
    sticky_session_type = "Insert"
  }
  health_check_config {
    health_check_http_version = "HTTP1.1"
    health_check_interval     = "2"
    health_check_method       = "HEAD"
    health_check_timeout      = "5"
    healthy_threshold         = "3"
    unhealthy_threshold       = "3"
    health_check_enabled      = true
  }
  protocol          = "HTTP"
  server_group_name = format("%%s4", var.name)
  server_group_type = "Ip"
}

resource "alicloud_alb_listener" "defaultc8auOh" {
  request_timeout = "60"
  listener_port   = "7788"
  default_actions {
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.defaultBJ6wNY.id
      }
    }
    type = "ForwardGroup"
  }
  idle_timeout      = "15"
  load_balancer_id  = alicloud_alb_load_balancer.lb.id
  listener_protocol = "HTTP"
}


`, name)
}

// Case Ascript_test241227_RequestHead 9640
func TestAccAliCloudAlbAScript_basic9640(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_ascript.default"
	ra := resourceAttrInit(resourceId, AlicloudAlbAScriptMap9640)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbAScript")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbAScriptBasicDependence9640)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"script_content":        "if and(match_re($uri, '^/1.txt$'), $arg_type) { rewrite(concat('/1.', $arg_type), 'break') }",
					"position":              "RequestHead",
					"ascript_name":          name,
					"listener_id":           "${alicloud_alb_listener.defaultc8auOh.id}",
					"enabled":               "true",
					"ext_attribute_enabled": "true",
					"ext_attributes": []map[string]interface{}{
						{
							"attribute_key":   "EsDebug",
							"attribute_value": "aaa",
						},
					},
					"dry_run": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"script_content":        "if and(match_re($uri, '^/1.txt$'), $arg_type) { rewrite(concat('/1.', $arg_type), 'break') }",
						"position":              "RequestHead",
						"ascript_name":          name,
						"listener_id":           CHECKSET,
						"enabled":               "true",
						"ext_attribute_enabled": "true",
						"ext_attributes.#":      "1",
						"dry_run":               "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"script_content": "if and(match_re($uri, '^/2.txt$'), $arg_type) { rewrite(concat('/1.', $arg_type), 'break') }",
					"ascript_name":   name + "_update",
					"ext_attributes": []map[string]interface{}{
						{
							"attribute_key":   "EsDebug",
							"attribute_value": "bbb",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"script_content":   "if and(match_re($uri, '^/2.txt$'), $arg_type) { rewrite(concat('/1.', $arg_type), 'break') }",
						"ascript_name":     name + "_update",
						"ext_attributes.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled":               "false",
					"ext_attribute_enabled": "false",
					"ext_attributes":        REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled":               "false",
						"ext_attribute_enabled": "false",
						"ext_attributes.#":      "0",
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

var AlicloudAlbAScriptMap9640 = map[string]string{
	"status": CHECKSET,
}

func AlicloudAlbAScriptBasicDependence9640(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zoneid-1" {
  default = "cn-hangzhou-j"
}

variable "zoneid-2" {
  default = "cn-hangzhou-k"
}

variable "regionid" {
  default = "cn-hangzhou"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vsw-1" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = var.zoneid-1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "vsw-2" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = var.zoneid-2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%%s2", var.name)
}

resource "alicloud_alb_load_balancer" "lb" {
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw-1.id
    zone_id    = alicloud_vswitch.vsw-1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw-2.id
    zone_id    = alicloud_vswitch.vsw-2.zone_id
  }
}

resource "alicloud_alb_server_group" "defaultBJ6wNY" {
  vpc_id    = alicloud_vpc.vpc.id
  scheduler = "Wrr"
  sticky_session_config {
    cookie_timeout      = "1000"
    sticky_session_type = "Insert"
  }
  health_check_config {
    health_check_http_version = "HTTP1.1"
    health_check_interval     = "2"
    health_check_method       = "HEAD"
    health_check_timeout      = "5"
    healthy_threshold         = "3"
    unhealthy_threshold       = "3"
    health_check_enabled      = true
  }
  protocol          = "HTTP"
  server_group_name = format("%%s4", var.name)
  server_group_type = "Ip"
}

resource "alicloud_alb_listener" "defaultc8auOh" {
  request_timeout = "60"
  listener_port   = "7788"
  default_actions {
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.defaultBJ6wNY.id
      }
    }
    type = "ForwardGroup"
  }
  idle_timeout      = "15"
  load_balancer_id  = alicloud_alb_load_balancer.lb.id
  listener_protocol = "HTTP"
}


`, name)
}

// Case Ascript_test241227_ResponseHead 9756
func TestAccAliCloudAlbAScript_basic9756(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_ascript.default"
	ra := resourceAttrInit(resourceId, AlicloudAlbAScriptMap9756)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbAScript")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbAScriptBasicDependence9756)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"script_content":        "id1 = rand(1, 99999)",
					"position":              "ResponseHead",
					"ascript_name":          name,
					"listener_id":           "${alicloud_alb_listener.defaultc8auOh.id}",
					"enabled":               "true",
					"ext_attribute_enabled": "true",
					"ext_attributes": []map[string]interface{}{
						{
							"attribute_key":   "EsDebug",
							"attribute_value": "aaa",
						},
					},
					"dry_run": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"script_content":        "id1 = rand(1, 99999)",
						"position":              "ResponseHead",
						"ascript_name":          name,
						"listener_id":           CHECKSET,
						"enabled":               "true",
						"ext_attribute_enabled": "true",
						"ext_attributes.#":      "1",
						"dry_run":               "false",
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

var AlicloudAlbAScriptMap9756 = map[string]string{
	"status": CHECKSET,
}

func AlicloudAlbAScriptBasicDependence9756(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zoneid-1" {
  default = "cn-hangzhou-j"
}

variable "zoneid-2" {
  default = "cn-hangzhou-k"
}

variable "regionid" {
  default = "cn-hangzhou"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vsw-1" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = var.zoneid-1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "vsw-2" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = var.zoneid-2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%%s2", var.name)
}

resource "alicloud_alb_load_balancer" "lb" {
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw-1.id
    zone_id    = alicloud_vswitch.vsw-1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw-2.id
    zone_id    = alicloud_vswitch.vsw-2.zone_id
  }
}

resource "alicloud_alb_server_group" "defaultBJ6wNY" {
  vpc_id    = alicloud_vpc.vpc.id
  scheduler = "Wrr"
  sticky_session_config {
    cookie_timeout      = "1000"
    sticky_session_type = "Insert"
  }
  health_check_config {
    health_check_http_version = "HTTP1.1"
    health_check_interval     = "2"
    health_check_method       = "HEAD"
    health_check_timeout      = "5"
    healthy_threshold         = "3"
    unhealthy_threshold       = "3"
    health_check_enabled      = true
  }
  protocol          = "HTTP"
  server_group_name = format("%%s4", var.name)
  server_group_type = "Ip"
}

resource "alicloud_alb_listener" "defaultc8auOh" {
  request_timeout = "60"
  listener_port   = "7788"
  default_actions {
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.defaultBJ6wNY.id
      }
    }
    type = "ForwardGroup"
  }
  idle_timeout      = "15"
  load_balancer_id  = alicloud_alb_load_balancer.lb.id
  listener_protocol = "HTTP"
}


`, name)
}

// Test Alb AScript. <<< Resource test cases, automatically generated.
