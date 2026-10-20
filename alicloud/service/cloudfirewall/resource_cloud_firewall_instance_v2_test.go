package cloudfirewall_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The account-site and monthly-day gates live in the root package's test files
// and are not importable from a service package, so they are kept here with the
// same semantics. The region is pinned in PreCheck, the same way the Cloudauth
// split package does it.
type AccountSite string

const (
	DomesticSite = AccountSite("Domestic")
	IntlSite     = AccountSite("International")
)

func testAccPreCheckWithAccountSiteType(t *testing.T, account AccountSite) {
	defaultAccount := string(DomesticSite)
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_ACCOUNT_SITE")); v != "" {
		defaultAccount = v
	} else if v = strings.TrimSpace(os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")); v != "" {
		defaultAccount = v
	} else if v = strings.TrimSpace(os.Getenv("ALICLOUD_ACCOUNT_TYPE")); v != "" {
		defaultAccount = v
	}
	if defaultAccount != string(account) {
		t.Skipf("Skipping unsupported account type %s-Site. It only supports %s-Site.", defaultAccount, account)
	}
}

func testAccPreCheckWithTime(t *testing.T, days []int) {
	skipped := true
	for _, d := range days {
		if time.Now().Day() == d {
			skipped = false
			break
		}
	}
	if skipped {
		t.Skipf("Skipping the test case when not in specified days %#v of every month", days)
	}
}

// Test CloudFirewall Instance. >>> Resource test cases, automatically generated.
// Case 国内版按量付费2.0 11709
func TestAccAliCloudCloudFirewallInstanceV2_basic11709(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
						"spec":         "payg_version",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sdl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sdl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":     "false",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11709_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
					"sdl":          "true",
					"cfw_log":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
						"spec":         "payg_version",
						"sdl":          "true",
						"cfw_log":      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

var AliCloudCloudFirewallInstanceV2Map11709 = map[string]string{
	"cfw_log":     acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"end_time":    acctest.CHECKSET,
	"user_status": acctest.CHECKSET,
	"status":      acctest.CHECKSET,
}

func AliCloudCloudFirewallInstanceV2BasicDependence11709(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 国际版按量付费2.0 11710
func TestAccAliCloudCloudFirewallInstanceV2_basic11710(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_intl",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_intl",
						"spec":         "payg_version",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sdl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sdl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":     "false",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11710_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_intl",
					"spec":         "payg_version",
					"sdl":          "true",
					"cfw_log":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_intl",
						"spec":         "payg_version",
						"sdl":          "true",
						"cfw_log":      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

// Case 国内版预付费2.0 11711
func TestAccAliCloudCloudFirewallInstanceV2_basic11711(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11711)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"product_code": "cfw",
					"product_type": "cfw_sub_public_cn",
					"spec":         "premium_version",
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"product_code": "cfw",
						"product_type": "cfw_sub_public_cn",
						"spec":         "premium_version",
						"period":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11711_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11711)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":          "Subscription",
					"product_code":          "cfw",
					"product_type":          "cfw_sub_public_cn",
					"spec":                  "enterprise_version",
					"sdl":                   "true",
					"cfw_log":               "false",
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
					"period":                "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "Subscription",
						"product_code":          "cfw",
						"product_type":          "cfw_sub_public_cn",
						"spec":                  "enterprise_version",
						"sdl":                   "true",
						"cfw_log":               "false",
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

var AliCloudCloudFirewallInstanceV2Map11711 = map[string]string{
	"cfw_log":               acctest.CHECKSET,
	"renewal_duration_unit": acctest.CHECKSET,
	"renewal_status":        acctest.CHECKSET,
	"create_time":           acctest.CHECKSET,
	"end_time":              acctest.CHECKSET,
	"user_status":           acctest.CHECKSET,
	"status":                acctest.CHECKSET,
}

// Case 国际版预付费2.0 11712
func TestAccAliCloudCloudFirewallInstanceV2_basic11712(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11711)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"product_code": "cfw",
					"product_type": "cfw_sub_public_intl",
					"spec":         "premium_version",
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"product_code": "cfw",
						"product_type": "cfw_sub_public_intl",
						"spec":         "premium_version",
						"period":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11712_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11711)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":          "Subscription",
					"product_code":          "cfw",
					"product_type":          "cfw_sub_public_intl",
					"spec":                  "ultimate_version",
					"sdl":                   "true",
					"cfw_log":               "false",
					"renewal_duration":      "1",
					"renewal_duration_unit": "Y",
					"renewal_status":        "AutoRenewal",
					"period":                "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "Subscription",
						"product_code":          "cfw",
						"product_type":          "cfw_sub_public_intl",
						"spec":                  "ultimate_version",
						"sdl":                   "true",
						"cfw_log":               "false",
						"renewal_duration":      "1",
						"renewal_duration_unit": "Y",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

// Test CloudFirewall Instance. <<< Resource test cases, automatically generated.

// Case 按量付费版修改 billing 方式: the API does not support converting an
// existing instance, so the plan must stay an in-place update (no replacement),
// the apply must fail with the console smooth-conversion guidance, and the
// instance ID must remain unchanged afterwards.
func TestAccAliCloudCloudFirewallInstanceV2_paymentTypeChange(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)

	var instanceId string
	captureInstanceId := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary.ID == "" {
			return fmt.Errorf("cannot find resource %s or its ID", resourceId)
		}
		instanceId = rs.Primary.ID
		return nil
	}
	checkInstanceIdUnchanged := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary.ID == "" {
			return fmt.Errorf("cannot find resource %s or its ID", resourceId)
		}
		if rs.Primary.ID != instanceId {
			return fmt.Errorf("instance was replaced: ID changed from %s to %s", instanceId, rs.Primary.ID)
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
						"spec":         "payg_version",
					}),
					captureInstanceId,
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionUpdate),
					},
				},
				ExpectError: regexp.MustCompile(`(?s)cannot change payment_type of Cloud Firewall instance .* is not supported by the API.*Convert Pay-As-You-Go to Subscription.*The existing instance is kept unchanged`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
					checkInstanceIdUnchanged,
				),
			},
		},
	})
}

// Case 按量付费版同时修改 billing 方式与 ForceNew 属性: product_type and spec
// are ForceNew, so the combined change would plan a replacement and release the
// existing instance before the Update guard could reject it. The plan must be
// rejected before any apply and the instance must be kept unchanged.
func TestAccAliCloudCloudFirewallInstanceV2_paymentTypeChangeWithReplacement(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCfwV2%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)

	var instanceId string
	captureInstanceId := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary.ID == "" {
			return fmt.Errorf("cannot find resource %s or its ID", resourceId)
		}
		instanceId = rs.Primary.ID
		return nil
	}
	checkInstanceIdUnchanged := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary.ID == "" {
			return fmt.Errorf("cannot find resource %s or its ID", resourceId)
		}
		if rs.Primary.ID != instanceId {
			return fmt.Errorf("instance was replaced: ID changed from %s to %s", instanceId, rs.Primary.ID)
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
						"spec":         "payg_version",
					}),
					captureInstanceId,
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"product_type": "cfw_sub_public_cn",
					"spec":         "premium_version",
				}),
				ExpectError: regexp.MustCompile(`(?s)cannot change payment_type of Cloud Firewall instance .* is not supported by the API.*Convert Pay-As-You-Go to Subscription.*The existing instance is kept unchanged`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
						"spec":         "payg_version",
					}),
					checkInstanceIdUnchanged,
				),
			},
		},
	})
}

// TestUnitAliCloudCloudFirewallInstanceV2PaymentTypeChangeWithReplacement
// verifies the plan-time guard through the SDK differ without any API call: a
// payment_type change alone stays an in-place update, while the same change
// combined with a ForceNew product_type or spec change must be rejected before
// Terraform can plan a replacement. A create with Subscription attributes must
// still be accepted, because the guard only applies to existing instances.
func TestUnitAliCloudCloudFirewallInstanceV2PaymentTypeChangeWithReplacement(t *testing.T) {
	resourceSchema := cloudfirewall.ResourceCloudFirewallInstanceV2()
	existingState := func() *sdkterraform.InstanceState {
		return &sdkterraform.InstanceState{
			ID: "cfw-unit-test-instance",
			Attributes: map[string]string{
				"payment_type": "PayAsYouGo",
				"product_code": "cfw",
				"product_type": "cfw_elasticity_public_cn",
				"spec":         "payg_version",
			},
		}
	}

	cases := []struct {
		name       string
		state      *sdkterraform.InstanceState
		config     map[string]interface{}
		wantError  bool
		wantCreate bool
	}{
		{
			name:  "payment type only keeps an in-place update",
			state: existingState(),
			config: map[string]interface{}{
				"payment_type": "Subscription",
				"product_code": "cfw",
				"product_type": "cfw_elasticity_public_cn",
				"spec":         "payg_version",
			},
		},
		{
			name:  "payment type with product type is rejected",
			state: existingState(),
			config: map[string]interface{}{
				"payment_type": "Subscription",
				"product_code": "cfw",
				"product_type": "cfw_sub_public_cn",
				"spec":         "payg_version",
			},
			wantError: true,
		},
		{
			name:  "payment type with spec is rejected",
			state: existingState(),
			config: map[string]interface{}{
				"payment_type": "Subscription",
				"product_code": "cfw",
				"product_type": "cfw_elasticity_public_cn",
				"spec":         "premium_version",
			},
			wantError: true,
		},
		{
			name:  "payment type with product type and spec is rejected",
			state: existingState(),
			config: map[string]interface{}{
				"payment_type": "Subscription",
				"product_code": "cfw",
				"product_type": "cfw_sub_public_cn",
				"spec":         "premium_version",
			},
			wantError: true,
		},
		{
			name:  "create with subscription attributes is allowed",
			state: &sdkterraform.InstanceState{Attributes: map[string]string{}},
			config: map[string]interface{}{
				"payment_type": "Subscription",
				"product_code": "cfw",
				"product_type": "cfw_sub_public_cn",
				"spec":         "premium_version",
			},
			wantCreate: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			diff, err := resourceSchema.Diff(context.Background(), tc.state, sdkterraform.NewResourceConfigRaw(tc.config), nil)
			if tc.wantError {
				if err == nil {
					t.Fatalf("expected the plan to be rejected, got diff %#v", diff)
				}
				if !strings.Contains(err.Error(), "cannot change payment_type of Cloud Firewall instance "+tc.state.ID) {
					t.Fatalf("unexpected plan error: %s", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected plan error: %s", err)
			}
			if tc.wantCreate {
				if diff == nil || diff.Destroy {
					t.Fatalf("expected a create plan, got diff %#v", diff)
				}
				paymentTypeDiff, ok := diff.Attributes["payment_type"]
				if !ok || paymentTypeDiff.Old != "" || paymentTypeDiff.New != "Subscription" {
					t.Fatalf("expected a create plan setting payment_type to Subscription, got %#v", paymentTypeDiff)
				}
				return
			}
			if diff != nil && diff.RequiresNew() {
				t.Fatalf("expected the plan to stay in place, but it requires a replacement")
			}
		})
	}
}
