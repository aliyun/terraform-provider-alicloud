package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Quotas TemplateService. >>> Resource test cases, automatically generated.
// Case 启用模版服务测试case_线上 6241
func TestAccAliCloudQuotasTemplateService_basic6241(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_service.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasTemplateServiceMap6241)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasTemplateService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastemplateservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasTemplateServiceBasicDependence6241)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_status": "-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_status": "-1",
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

var AlicloudQuotasTemplateServiceMap6241 = map[string]string{}

func AlicloudQuotasTemplateServiceBasicDependence6241(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 启用模版服务测试case 5843
func TestAccAliCloudQuotasTemplateService_basic5843(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_service.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasTemplateServiceMap5843)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasTemplateService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastemplateservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasTemplateServiceBasicDependence5843)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_status": "-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_status": "-1",
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

var AlicloudQuotasTemplateServiceMap5843 = map[string]string{}

func AlicloudQuotasTemplateServiceBasicDependence5843(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Quotas TemplateService. <<< Resource test cases, automatically generated.
