package dmsenterprise_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	dmsenterprise "github.com/aliyun/terraform-provider-alicloud/alicloud/service/dmsenterprise"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test DMSEnterprise AuthorityTemplate. >>> Resource test cases, automatically generated.
// Case 4696
func TestAccAliCloudDMSEnterpriseAuthorityTemplate_basic4696(t *testing.T) {
	tfacctest.CheckoutSupportedRegions(t, true, connectivity.DMSEnterpriseSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_authority_template.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AlicloudDMSEnterpriseAuthorityTemplateMap4696)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dmsenterprise.NewDMSEnterpriseServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDMSEnterpriseAuthorityTemplate")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdmsenterpriseauthoritytemplate%d", tfacctest.DefaultRegionToTest, rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AlicloudDMSEnterpriseAuthorityTemplateBasicDependence4696)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"authority_template_name": name,
					"tid":                     "${data.alicloud_dms_user_tenants.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authority_template_name": name,
						"tid":                     tfacctest.CHECKSET,
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
						"tid":                     tfacctest.CHECKSET,
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
	"create_time":           tfacctest.CHECKSET,
	"authority_template_id": tfacctest.CHECKSET,
	"tid":                   tfacctest.CHECKSET,
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
	tfacctest.CheckoutSupportedRegions(t, true, connectivity.DMSEnterpriseSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_authority_template.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AlicloudDMSEnterpriseAuthorityTemplateMap4696)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dmsenterprise.NewDMSEnterpriseServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDMSEnterpriseAuthorityTemplate")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdmsenterpriseauthoritytemplate%d", tfacctest.DefaultRegionToTest, rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AlicloudDMSEnterpriseAuthorityTemplateBasicDependence4696)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"tid":                     tfacctest.CHECKSET,
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
