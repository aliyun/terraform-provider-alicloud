package quotas_test

// collection-order: bypass aliyun_uids -- pre-existing TypeList unchanged by pure package migration; original schema, request/read ordering and one-account test preserved; no new reorder acceptance claimed

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/quotas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Quotas TemplateApplications. >>> Resource test cases, automatically generated.
// Case 5294
// lintignore: AT001
func TestAccAliCloudQuotasTemplateApplications_basic5294(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_applications.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudQuotasTemplateApplicationsMap5294)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return quotas.NewQuotasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeQuotasTemplateApplications")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotasapp%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudQuotasTemplateApplicationsBasicDependence5294)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
					"product_code":      "vpc",
					"quota_category":    "FlowControl",
					"aliyun_uids": []string{
						"${alicloud_resource_manager_account.account.id}"},
					"desire_value": "6",
					"notice_type":  "0",
					"env_language": "zh",
					"reason":       "测试",
					"dimensions": []map[string]interface{}{
						{
							"key":   "apiName",
							"value": "GetProductQuotaDimension",
						},
						{
							"key":   "apiVersion",
							"value": "2020-05-10",
						},
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
						"product_code":      "vpc",
						"quota_category":    "FlowControl",
						"aliyun_uids.#":     "1",
						"desire_value":      "6",
						"notice_type":       "0",
						"env_language":      "zh",
						"reason":            "测试",
						"dimensions.#":      "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_language", "notice_type"},
			},
		},
	})
}

var AlicloudQuotasTemplateApplicationsMap5294 = map[string]string{
	"quota_application_details.#": acctest.CHECKSET,
}

func AlicloudQuotasTemplateApplicationsBasicDependence5294(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_resource_manager_account" "account" {
  display_name = var.name
  abandon_able_check_id = ["SP_fc_fc"]
}


`, name)
}

// Case 5284
// lintignore: AT001
func TestAccAliCloudQuotasTemplateApplications_basic5284(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_applications.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudQuotasTemplateApplicationsMap5284)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return quotas.NewQuotasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeQuotasTemplateApplications")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotas%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudQuotasTemplateApplicationsBasicDependence5284)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "q_user_poc_instance_amount",
					"product_code":      "computenest",
					"quota_category":    "CommonQuota",
					"aliyun_uids": []string{
						"${alicloud_resource_manager_account.account.id}"},
					"desire_value": "7",
					"notice_type":  "0",
					"env_language": "zh",
					"reason":       "测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "q_user_poc_instance_amount",
						"product_code":      "computenest",
						"quota_category":    "CommonQuota",
						"aliyun_uids.#":     "1",
						"desire_value":      "7",
						"notice_type":       "0",
						"env_language":      "zh",
						"reason":            "测试",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_language", "notice_type"},
			},
		},
	})
}

var AlicloudQuotasTemplateApplicationsMap5284 = map[string]string{
	"quota_application_details.#": acctest.CHECKSET,
}

func AlicloudQuotasTemplateApplicationsBasicDependence5284(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_resource_manager_account" "account" {
  abandon_able_check_id = ["SP_fc_fc"]
  display_name = var.name
}


`, name)
}

// Case 5278
// lintignore: AT001
func TestAccAliCloudQuotasTemplateApplications_basic5278(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_applications.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudQuotasTemplateApplicationsMap5278)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return quotas.NewQuotasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeQuotasTemplateApplications")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastp%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudQuotasTemplateApplicationsBasicDependence5278)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "quotas.label_multi/A",
					"product_code":      "quotas",
					"effective_time":    "2023-12-03T16:00:00Z",
					"quota_category":    "WhiteListLabel",
					"aliyun_uids": []string{
						"${alicloud_resource_manager_account.account.id}"},
					"expire_time":  "2024-12-26T16:00:00Z",
					"desire_value": "1",
					"notice_type":  "0",
					"env_language": "zh",
					"reason":       "测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "quotas.label_multi/A",
						"product_code":      "quotas",
						"effective_time":    "2023-12-03T16:00:00Z",
						"quota_category":    "WhiteListLabel",
						"aliyun_uids.#":     "1",
						"expire_time":       "2024-12-26T16:00:00Z",
						"desire_value":      "1",
						"notice_type":       "0",
						"env_language":      "zh",
						"reason":            "测试",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_language", "notice_type"},
			},
		},
	})
}

var AlicloudQuotasTemplateApplicationsMap5278 = map[string]string{
	"quota_application_details.#": acctest.CHECKSET,
}

func AlicloudQuotasTemplateApplicationsBasicDependence5278(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_resource_manager_account" "account" {
  display_name = var.name
  abandon_able_check_id = ["SP_fc_fc"]
}


`, name)
}

// Test Quotas TemplateApplications. <<< Resource test cases, automatically generated.
