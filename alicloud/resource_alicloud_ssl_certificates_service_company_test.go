// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test SslCertificatesService Company. >>> Resource test cases, automatically generated.
// Case Company资源用例 12894
func TestAccAliCloudSslCertificatesServiceCompany_basic12894(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_company.default"
	ra := resourceAttrInit(resourceId, AlicloudSslCertificatesServiceCompanyMap12894)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceCompany")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSslCertificatesServiceCompanyBasicDependence12894)
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
					"company_address": "西安市",
					"company_name":    name,
					"department":      "测试部门1",
					"city":            "西安",
					"company_type":    "1",
					"country_code":    "111122",
					"post_code":       "11112233",
					"company_code":    "12312311",
					"company_phone":   "15101081174",
					"province":        "陕西",
					"lang":            "zh",
					"company_email":   "test@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"company_address": "西安市",
						"company_name":    name,
						"department":      "测试部门1",
						"city":            "西安",
						"company_type":    "1",
						"country_code":    CHECKSET,
						"post_code":       CHECKSET,
						"company_code":    CHECKSET,
						"company_phone":   CHECKSET,
						"province":        "陕西",
						"lang":            CHECKSET,
						"company_email":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"company_address": "西安市12333",
					"company_name":    name + "_update",
					"department":      "测试部门2",
					"city":            "北京",
					"company_type":    "2",
					"country_code":    "222211",
					"post_code":       "22221111",
					"company_code":    "111122",
					"company_phone":   "15201081174",
					"province":        "西安市111",
					"lang":            "en",
					"company_email":   "update@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"company_address": "西安市12333",
						"company_name":    name + "_update",
						"department":      "测试部门2",
						"city":            "北京",
						"company_type":    "2",
						"country_code":    CHECKSET,
						"post_code":       CHECKSET,
						"company_code":    CHECKSET,
						"company_phone":   CHECKSET,
						"province":        "西安市111",
						"lang":            CHECKSET,
						"company_email":   CHECKSET,
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

var AlicloudSslCertificatesServiceCompanyMap12894 = map[string]string{}

func AlicloudSslCertificatesServiceCompanyBasicDependence12894(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test SslCertificatesService Company. <<< Resource test cases, automatically generated.
