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

// Test Governance Baseline. >>> Resource test cases, automatically generated.
// Case 账号工厂基线资源测试 7320
func TestAccAliCloudGovernanceBaseline_basic7320(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGovernanceBaselineMap7320)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return governance.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGovernanceBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernancebaseline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceBaselineBasicDependence7320)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						{
							"version": "1.0",
							"name":    "${var.item_password_policy}",
							"config":  "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf auto test baseline",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf auto test baseline",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						{
							"config":  "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}",
							"name":    "${var.item_password_policy}",
							"version": "1.0",
						},
						{
							"version": "1.0",
							"name":    "${var.item_ram_security}",
							"config":  "{\\\"LoginNetworkMasks\\\":\\\"\\\",\\\"EnableSaveMfaTicket\\\":false,\\\"AllowUserToChangePassword\\\":true,\\\"AllowUserToManageAccessKeys\\\":false,\\\"AllowUserToManageMfaDevices\\\":true,\\\"EnforceMfaForLogin\\\":false,\\\"LoginSessionDuration\\\":6}",
						},
						{
							"config":  "{\\\"EnabledServices\\\":[\\\"CONFIG\\\"],\\\"EnabledServiceConfigurations\\\":[]}",
							"name":    "${var.item_services}",
							"version": "1.0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf auto test baseline update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf auto test baseline update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						{
							"config":  "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}",
							"name":    "${var.item_password_policy}",
							"version": "1.0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "abc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						{
							"version": "1.0",
							"name":    "${var.item_password_policy}",
							"config":  "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}",
						},
					},
					"description":   "tf auto test baseline",
					"baseline_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "1",
						"description":      "tf auto test baseline",
						"baseline_name":    name + "_update",
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

// TestAccAliCloudGovernanceBaseline_baselineItemsOrder exercises a convergent
// reorder of the baseline_items TypeList: apply the members in one order,
// then apply a permutation of the same members and require the post-apply
// plan to be empty. Satisfies the scripts/collection-order TypeList order
// coverage requirement; whether the Cloud Governance Center API echoes the
// submitted order is established when the acceptance test runs.
func TestAccAliCloudGovernanceBaseline_baselineItemsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGovernanceBaselineMap7320)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return governance.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGovernanceBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernancebaseline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceBaselineBasicDependence7320)
	passwordPolicyConfig := "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}"
	ramSecurityConfig := "{\\\"LoginNetworkMasks\\\":\\\"\\\",\\\"EnableSaveMfaTicket\\\":false,\\\"AllowUserToChangePassword\\\":true,\\\"AllowUserToManageAccessKeys\\\":false,\\\"AllowUserToManageMfaDevices\\\":true,\\\"EnforceMfaForLogin\\\":false,\\\"LoginSessionDuration\\\":6}"
	passwordPolicyItem := map[string]interface{}{
		"version": "1.0",
		"name":    "${var.item_password_policy}",
		"config":  passwordPolicyConfig,
	}
	ramSecurityItem := map[string]interface{}{
		"version": "1.0",
		"name":    "${var.item_ram_security}",
		"config":  ramSecurityConfig,
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_name": name,
					"baseline_items": []map[string]interface{}{
						passwordPolicyItem,
						ramSecurityItem,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_name":    name,
						"baseline_items.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						ramSecurityItem,
						passwordPolicyItem,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "2",
					}),
				),
			},
		},
	})
}

var AlicloudGovernanceBaselineMap7320 = map[string]string{}

func AlicloudGovernanceBaselineBasicDependence7320(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "item_password_policy" {
  default = "ACS-BP_ACCOUNT_FACTORY_RAM_USER_PASSWORD_POLICY"
}

variable "baseline_name_update" {
  default = "tf-auto-test-baseline-update"
}

variable "item_services" {
  default = "ACS-BP_ACCOUNT_FACTORY_SUBSCRIBE_SERVICES"
}

variable "baseline_name" {
  default = "tf-auto-test-baseline"
}

variable "item_ram_security" {
  default = "ACS-BP_ACCOUNT_FACTORY_RAM_SECURITY_PREFERENCE"
}


`, name)
}

// Case 账号工厂基线资源测试 7320  twin
func TestAccAliCloudGovernanceBaseline_basic7320_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGovernanceBaselineMap7320)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return governance.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGovernanceBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernancebaseline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceBaselineBasicDependence7320)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						{
							"version": "1.0",
							"name":    "${var.item_password_policy}",
							"config":  "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}",
						},
					},
					"description":   "tf auto test baseline",
					"baseline_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "1",
						"description":      "tf auto test baseline",
						"baseline_name":    name,
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

// Case 账号工厂基线资源测试 7320  raw
func TestAccAliCloudGovernanceBaseline_basic7320_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_governance_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGovernanceBaselineMap7320)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return governance.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGovernanceBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgovernancebaseline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGovernanceBaselineBasicDependence7320)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						{
							"version": "1.0",
							"name":    "${var.item_password_policy}",
							"config":  "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}",
						},
					},
					"description":   "tf auto test baseline",
					"baseline_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "1",
						"description":      "tf auto test baseline",
						"baseline_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						{
							"config":  "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}",
							"name":    "${var.item_password_policy}",
							"version": "1.0",
						},
						{
							"version": "1.0",
							"name":    "${var.item_ram_security}",
							"config":  "{\\\"LoginNetworkMasks\\\":\\\"\\\",\\\"EnableSaveMfaTicket\\\":false,\\\"AllowUserToChangePassword\\\":true,\\\"AllowUserToManageAccessKeys\\\":false,\\\"AllowUserToManageMfaDevices\\\":true,\\\"EnforceMfaForLogin\\\":false,\\\"LoginSessionDuration\\\":6}",
						},
						{
							"config":  "{\\\"EnabledServices\\\":[\\\"CONFIG\\\"],\\\"EnabledServiceConfigurations\\\":[]}",
							"name":    "${var.item_services}",
							"version": "1.0",
						},
					},
					"description":   "tf auto test baseline update",
					"baseline_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "3",
						"description":      "tf auto test baseline update",
						"baseline_name":    name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_items": []map[string]interface{}{
						{
							"config":  "{\\\"MinimumPasswordLength\\\":8,\\\"RequireLowercaseCharacters\\\":true,\\\"RequireUppercaseCharacters\\\":true,\\\"RequireNumbers\\\":true,\\\"RequireSymbols\\\":true,\\\"MaxPasswordAge\\\":0,\\\"HardExpiry\\\":false,\\\"PasswordReusePrevention\\\":0,\\\"MaxLoginAttempts\\\":0}",
							"name":    "${var.item_password_policy}",
							"version": "1.0",
						},
					},
					"description":   "abc",
					"baseline_name": name + "_update2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_items.#": "1",
						"description":      "abc",
						"baseline_name":    name + "_update2",
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

// Test Governance Baseline. <<< Resource test cases, automatically generated.
