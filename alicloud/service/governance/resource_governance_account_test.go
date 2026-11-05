package governance_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/governance"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudGovernanceAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGovernanceAccountMap7372)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return governance.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGovernanceAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernanceaccount%d", acctest.DefaultRegionToTest, rand)
	defaultDomainName := fmt.Sprintf("%s.onaliyun.com", name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceAccountBasicDependence7372)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_id":          "${data.alicloud_resource_manager_accounts.default.ids.0}",
					"baseline_id":         "${data.alicloud_governance_baselines.default.ids.0}",
					"payer_account_id":    "${data.alicloud_account.default.id}",
					"display_name":        name,
					"account_name_prefix": name,
					"folder_id":           "${data.alicloud_resource_manager_folders.default.ids.0}",
					"default_domain_name": defaultDomainName,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_id":         acctest.CHECKSET,
						"payer_account_id":    acctest.CHECKSET,
						"display_name":        acctest.CHECKSET,
						"account_name_prefix": acctest.CHECKSET,
						"folder_id":           acctest.CHECKSET,
						"default_domain_name": acctest.CHECKSET,
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

// lintignore: AT001
func TestAccAliCloudGovernanceAccount_basic7372(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGovernanceAccountMap7372)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return governance.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGovernanceAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernanceaccount%d", acctest.DefaultRegionToTest, rand)
	defaultDomainName := fmt.Sprintf("%s.onaliyun.com", name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceAccountBasicDependence7372)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
						"baseline_id":         acctest.CHECKSET,
						"payer_account_id":    acctest.CHECKSET,
						"display_name":        acctest.CHECKSET,
						"account_name_prefix": acctest.CHECKSET,
						"folder_id":           acctest.CHECKSET,
						"default_domain_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_id": "${data.alicloud_governance_baselines.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_tags": []map[string]interface{}{
						{
							"tag_key":   "tag-key2",
							"tag_value": "tag-value2",
						},
						{
							"tag_key":   "tag-key3",
							"tag_value": "tag-value3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_tags.#": "2",
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

// TestAccAliCloudGovernanceAccount_accountTagsOrder exercises a convergent
// reorder of the account_tags TypeList: apply the members in one order, then
// apply a permutation of the same members and require the post-apply plan to
// be empty. Satisfies the scripts/collection-order TypeList order coverage
// requirement; whether the Cloud Governance Center API echoes the submitted
// order is established when the acceptance test runs.
// lintignore: AT001
func TestAccAliCloudGovernanceAccount_accountTagsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGovernanceAccountMap7372)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return governance.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGovernanceAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernanceaccount%d", acctest.DefaultRegionToTest, rand)
	defaultDomainName := fmt.Sprintf("%s.onaliyun.com", name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceAccountBasicDependence7372)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_id":          "${data.alicloud_resource_manager_accounts.default.ids.0}",
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
						"baseline_id":         acctest.CHECKSET,
						"payer_account_id":    acctest.CHECKSET,
						"display_name":        name,
						"account_name_prefix": name,
						"folder_id":           acctest.CHECKSET,
						"default_domain_name": defaultDomainName,
						"account_tags.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_tags": []map[string]interface{}{
						{
							"tag_key":   "tag-key2",
							"tag_value": "tag-value2",
						},
						{
							"tag_key":   "tag-key1",
							"tag_value": "tag-value1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_tags.#": "2",
					}),
				),
			},
		},
	})
}

var AlicloudGovernanceAccountMap7372 = map[string]string{
	"status":     acctest.CHECKSET,
	"account_id": acctest.CHECKSET,
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

data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}


`, name)
}
