// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package sslcertificateservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sslcertificateservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Pre-existing TypeList reorder-coverage debt on webhook_list, unchanged by this pure directory
// migration (the resource file was relocated only; no schema attribute was added or modified).
// Marked per scripts/collection-order/README.md "Bypassing stock attributes".
// collection-order: bypass webhook_list -- pre-existing attribute outside this change

// Test SslCertificatesService Contact. >>> Resource test cases, automatically generated.
// Case Contact资源用例 12895
func TestAccAliCloudSslCertificatesServiceContact_basic12895(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_contact.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSslCertificatesServiceContactMap12895)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sslcertificateservice.NewSslCertificatesServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSslCertificatesServiceContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSslCertificatesServiceContactBasicDependence12895)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSslCertificatesServiceContactMapIdcard)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sslcertificateservice.NewSslCertificatesServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSslCertificatesServiceContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSslCertificatesServiceContactBasicDependenceIdcard)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSslCertificatesServiceContactMapWebhooks)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sslcertificateservice.NewSslCertificatesServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSslCertificatesServiceContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSslCertificatesServiceContactBasicDependenceWebhooks)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
						"mobile": acctest.CHECKSET,
						"email":  acctest.CHECKSET,
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
