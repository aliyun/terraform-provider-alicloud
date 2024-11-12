package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudGovernanceAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_account.default"
	ra := resourceAttrInit(resourceId, AlicloudGovernanceAccountMap7372)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GovernanceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGovernanceAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernanceaccount%d", defaultRegionToTest, rand)
	defaultDomainName := fmt.Sprintf("%s.onaliyun.com", name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceAccountBasicDependence7372)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_id":          "1493822914031335",
					"baseline_id":         "${data.alicloud_governance_baselines.default.ids.0}",
					"payer_account_id":    "${data.alicloud_account.default.id}",
					"display_name":        name,
					"account_name_prefix": name,
					"folder_id":           "${data.alicloud_resource_manager_folders.default.ids.0}",
					"default_domain_name": defaultDomainName,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_id":         CHECKSET,
						"payer_account_id":    CHECKSET,
						"display_name":        CHECKSET,
						"account_name_prefix": CHECKSET,
						"folder_id":           CHECKSET,
						"default_domain_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_tags": []map[string]interface{}{
						{
							"tag_key":   "tag-key1",
							"tag_value": "tag-value1",
						},
						{
							"tag_key":   "tag-key2",
							"tag_value": "tag-value2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_name_prefix", "display_name", "folder_id", "payer_account_id", "baseline_id", "default_domain_name"},
			},
		},
	})
}

func TestAccAliCloudGovernanceAccount_basic7372(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_account.default"
	ra := resourceAttrInit(resourceId, AlicloudGovernanceAccountMap7372)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GovernanceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGovernanceAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernanceaccount%d", defaultRegionToTest, rand)
	defaultDomainName := fmt.Sprintf("%s.onaliyun.com", name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceAccountBasicDependence7372)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_id":         "${data.alicloud_governance_baselines.default.ids.0}",
					"payer_account_id":    "${data.alicloud_account.default.id}",
					"display_name":        name,
					"account_name_prefix": name,
					"folder_id":           "${data.alicloud_resource_manager_folders.default.ids.0}",
					"default_domain_name": defaultDomainName,
					"account_tags": []map[string]interface{}{
						{
							"tag_key":   "tag-key1",
							"tag_value": "tag-value1",
						},
						{
							"tag_key":   "tag-key2",
							"tag_value": "tag-value2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_id":         CHECKSET,
						"payer_account_id":    CHECKSET,
						"display_name":        CHECKSET,
						"account_name_prefix": CHECKSET,
						"folder_id":           CHECKSET,
						"default_domain_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_id": "${data.alicloud_governance_baselines.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_name_prefix", "display_name", "folder_id", "payer_account_id", "baseline_id", "default_domain_name"},
			},
		},
	})
}

var AlicloudGovernanceAccountMap7372 = map[string]string{
	"status":     CHECKSET,
	"account_id": CHECKSET,
}

func AlicloudGovernanceAccountBasicDependence7372(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
}

data "alicloud_governance_baselines" "default" {
}

data "alicloud_resource_manager_folders" "default" {
}


`, name)
}
