// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test SslCertificatesService Contact. >>> Resource test cases, automatically generated.
// Case Contact资源用例 12895
func TestAccAliCloudSslCertificatesServiceContact_basic12895(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_contact.default"
	ra := resourceAttrInit(resourceId, AlicloudSslCertificatesServiceContactMap12895)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSslCertificatesServiceContactBasicDependence12895)
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
					"name":   name,
					"mobile": "13312345678",
					"email":  "test1@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name,
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   name + "update",
					"mobile": "13312345678",
					"email":  "test1@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   name + "update",
					"mobile": "13300001111",
					"email":  "test1@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   name + "update",
					"mobile": "13300001111",
					"email":  "test2@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"email", "idcard", "mobile", "webhook_list"},
			},
		},
	})
}

var AlicloudSslCertificatesServiceContactMap12895 = map[string]string{}

func AlicloudSslCertificatesServiceContactBasicDependence12895(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// TestAccAliCloudSslCertificatesServiceContact_basicIdcard covers the idcard
// attribute. CreateContact accepts idcard on any account; CFCA only consumes the
// idcard when the contact is used to apply for a CFCA-brand certificate (other
// brands ignore it), so no account gating is needed.
func TestAccAliCloudSslCertificatesServiceContact_basicIdcard(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_contact.default"
	ra := resourceAttrInit(resourceId, AlicloudSslCertificatesServiceContactMapIdcard)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSslCertificatesServiceContactBasicDependenceIdcard)
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
					"name":   name,
					"mobile": "13312345678",
					"email":  "test1@example.com",
					"idcard": "110101199003078515",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name,
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   name + "update",
					"mobile": "13312345678",
					"email":  "test1@example.com",
					"idcard": "110101199003078515",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   name + "update",
					"mobile": "13300001111",
					"email":  "test1@example.com",
					"idcard": "110101199003078515",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   name + "update",
					"mobile": "13300001111",
					"email":  "test2@example.com",
					"idcard": "110101199003078515",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"email", "idcard", "mobile", "webhook_list"},
			},
		},
	})
}

var AlicloudSslCertificatesServiceContactMapIdcard = map[string]string{}

func AlicloudSslCertificatesServiceContactBasicDependenceIdcard(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// TestAccAliCloudSslCertificatesServiceContact_basicWebhooks covers the
// webhook_list attribute. CreateContact accepts webhooks whose elements are
// DingTalk robot URLs (https://oapi.dingtalk.com/robot/send?access_token=...);
// other URLs are rejected with HTTP 400 InvalidParameter, so the test uses a
// DingTalk-format URL.
func TestAccAliCloudSslCertificatesServiceContact_basicWebhooks(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_contact.default"
	ra := resourceAttrInit(resourceId, AlicloudSslCertificatesServiceContactMapWebhooks)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSslCertificatesServiceContactBasicDependenceWebhooks)
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
					"name":         name,
					"mobile":       "13312345678",
					"email":        "test1@example.com",
					"webhook_list": []interface{}{"https://oapi.dingtalk.com/robot/send?access_token=testacc0123456789"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name,
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         name + "update",
					"mobile":       "13312345678",
					"email":        "test1@example.com",
					"webhook_list": []interface{}{"https://oapi.dingtalk.com/robot/send?access_token=testacc0123456789"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         name + "update",
					"mobile":       "13300001111",
					"email":        "test1@example.com",
					"webhook_list": []interface{}{"https://oapi.dingtalk.com/robot/send?access_token=testacc0123456789"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         name + "update",
					"mobile":       "13300001111",
					"email":        "test2@example.com",
					"webhook_list": []interface{}{"https://oapi.dingtalk.com/robot/send?access_token=testacc0123456789"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         name + "update",
					"mobile":       "13300001111",
					"email":        "test2@example.com",
					"webhook_list": []interface{}{"https://oapi.dingtalk.com/robot/send?access_token=testaccupdate9876"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name + "update",
						"mobile": CHECKSET,
						"email":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"email", "idcard", "mobile", "webhook_list"},
			},
		},
	})
}

var AlicloudSslCertificatesServiceContactMapWebhooks = map[string]string{}

func AlicloudSslCertificatesServiceContactBasicDependenceWebhooks(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test SslCertificatesService Contact. <<< Resource test cases, automatically generated.
