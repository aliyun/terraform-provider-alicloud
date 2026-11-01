package quotas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/quotas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Quotas TemplateService. >>> Resource test cases, automatically generated.
// Case 启用模版服务测试case_线上 6241
// lintignore: AT001
func TestAccAliCloudQuotasTemplateService_basic6241(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudQuotasTemplateServiceMap6241)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return quotas.NewQuotasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeQuotasTemplateService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastemplateservice%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudQuotasTemplateServiceBasicDependence6241)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
// lintignore: AT001
func TestAccAliCloudQuotasTemplateService_basic5843(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudQuotasTemplateServiceMap5843)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return quotas.NewQuotasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeQuotasTemplateService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastemplateservice%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudQuotasTemplateServiceBasicDependence5843)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
