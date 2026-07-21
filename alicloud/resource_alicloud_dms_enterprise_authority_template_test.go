package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DMSEnterprise AuthorityTemplate. >>> Resource test cases, automatically generated.
// Case 4696
func TestAccAliCloudDMSEnterpriseAuthorityTemplate_basic4696(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_authority_template.default"
	ra := resourceAttrInit(resourceId, AlicloudDMSEnterpriseAuthorityTemplateMap4696)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DMSEnterpriseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDMSEnterpriseAuthorityTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdmsenterpriseauthoritytemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDMSEnterpriseAuthorityTemplateBasicDependence4696)
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
					"authority_template_name": name,
					"tid":                     "${data.alicloud_dms_user_tenants.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authority_template_name": name,
						"tid":                     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "资源用例测试权限模板",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "资源用例测试权限模板",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"authority_template_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authority_template_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "资源用例测试权限模板-更新",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "资源用例测试权限模板-更新",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"authority_template_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authority_template_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "资源用例测试权限模板",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "资源用例测试权限模板",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"authority_template_name": name + "_update",
					"description":             "资源用例测试权限模板",
					"tid":                     "${data.alicloud_dms_user_tenants.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authority_template_name": name + "_update",
						"description":             "资源用例测试权限模板",
						"tid":                     CHECKSET,
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

var AlicloudDMSEnterpriseAuthorityTemplateMap4696 = map[string]string{
	"create_time":           CHECKSET,
	"authority_template_id": CHECKSET,
	"tid":                   CHECKSET,
}

func AlicloudDMSEnterpriseAuthorityTemplateBasicDependence4696(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "v_tid" {
  default = "1"
}

data "alicloud_dms_user_tenants" "default" {
	status = "ACTIVE"
}

`, name)
}

// Case 4696  twin
func TestAccAliCloudDMSEnterpriseAuthorityTemplate_basic4696_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_authority_template.default"
	ra := resourceAttrInit(resourceId, AlicloudDMSEnterpriseAuthorityTemplateMap4696)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DMSEnterpriseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDMSEnterpriseAuthorityTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdmsenterpriseauthoritytemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDMSEnterpriseAuthorityTemplateBasicDependence4696)
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
					"authority_template_name": name,
					"description":             "资源用例测试权限模板",
					"tid":                     "${data.alicloud_dms_user_tenants.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authority_template_name": name,
						"description":             "资源用例测试权限模板",
						"tid":                     CHECKSET,
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

// Test DMSEnterprise AuthorityTemplate. <<< Resource test cases, automatically generated.
