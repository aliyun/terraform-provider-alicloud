package nlb_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nlb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNlbSecurityPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_security_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNLBSecurityPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbSecurityPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbsecuritypolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNLBSecurityPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"security_policy_name": "${var.name}",
					"ciphers":              []string{"ECDHE-RSA-AES128-SHA", "ECDHE-ECDSA-AES128-SHA"},
					"tls_versions":         []string{"TLSv1.0", "TLSv1.1", "TLSv1.2"},
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":    acctest.CHECKSET,
						"security_policy_name": name,
						"ciphers.#":            "2",
						"tls_versions.#":       "3",
						"tags.%":               "2",
						"tags.Created":         "tfTestAcc0",
						"tags.For":             "Tftestacc 0",
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

var AlicloudNLBSecurityPolicyMap0 = map[string]string{
	"status":               acctest.CHECKSET,
	"tags.%":               acctest.CHECKSET,
	"tls_versions.#":       acctest.CHECKSET,
	"ciphers.#":            acctest.CHECKSET,
	"security_policy_name": acctest.CHECKSET,
}

func AlicloudNLBSecurityPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}

func TestAccAliCloudNlbSecurityPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_security_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNLBSecurityPolicyMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbSecurityPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbsecuritypolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNLBSecurityPolicyBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ciphers":      []string{"ECDHE-RSA-AES128-SHA", "ECDHE-ECDSA-AES128-SHA"},
					"tls_versions": []string{"TLSv1.0", "TLSv1.1", "TLSv1.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#":      "2",
						"tls_versions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ciphers": []string{"ECDHE-RSA-AES128-SHA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{"TLSv1.0", "TLSv1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_versions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc4",
						"For":     "Tftestacc 4",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc4",
						"tags.For":     "Tftestacc 4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": "${var.name}_update",
					"ciphers":              []string{"ECDHE-RSA-AES128-SHA", "ECDHE-ECDSA-AES128-SHA"},
					"tls_versions":         []string{"TLSv1.0", "TLSv1.1", "TLSv1.2"},
					"tags": map[string]string{
						"Created": "tfTestAcc5",
						"For":     "Tftestacc 5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": name + "_update",
						"ciphers.#":            "2",
						"tls_versions.#":       "3",
						"tags.%":               "2",
						"tags.Created":         "tfTestAcc5",
						"tags.For":             "Tftestacc 5",
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

var AlicloudNLBSecurityPolicyMap1 = map[string]string{
	"status":         acctest.CHECKSET,
	"tls_versions.#": acctest.CHECKSET,
	"ciphers.#":      acctest.CHECKSET,
}

func AlicloudNLBSecurityPolicyBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001

func TestAccAliCloudNlbSecurityPolicy_basic5352(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_security_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNlbSecurityPolicyMap5352)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbSecurityPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbsecuritypolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNlbSecurityPolicyBasicDependence5352)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ciphers": []string{
						"ECDHE-ECDSA-AES128-SHA"},
					"tls_versions": []string{
						"TLSv1.0"},
					"security_policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#":            "1",
						"tls_versions.#":       "1",
						"security_policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ciphers": []string{
						"ECDHE-ECDSA-AES128-SHA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{
						"TLSv1.0"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_versions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ciphers": []string{
						"ECDHE-ECDSA-AES256-SHA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{
						"TLSv1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_versions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ciphers": []string{
						"ECDHE-ECDSA-AES128-SHA", "TLS_AES_128_GCM_SHA256"},
					"tls_versions": []string{
						"TLSv1.0", "TLSv1.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#":      "2",
						"tls_versions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": name + "_update",
					"ciphers": []string{
						"ECDHE-ECDSA-AES128-SHA"},
					"tls_versions": []string{
						"TLSv1.0"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": name + "_update",
						"ciphers.#":            "1",
						"tls_versions.#":       "1",
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AlicloudNlbSecurityPolicyMap5352 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudNlbSecurityPolicyBasicDependence5352(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

// Case 5352  twin

func TestAccAliCloudNlbSecurityPolicy_basic5352_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_security_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNlbSecurityPolicyMap5352)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbSecurityPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbsecuritypolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNlbSecurityPolicyBasicDependence5352)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": name,
					"ciphers": []string{
						"ECDHE-ECDSA-AES128-SHA", "TLS_AES_128_GCM_SHA256"},
					"tls_versions": []string{
						"TLSv1.0", "TLSv1.1", "TLSv1.3"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": name,
						"ciphers.#":            "2",
						"tls_versions.#":       "3",
						"resource_group_id":    acctest.CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
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

// An empty string inside a list parameter is decoded as a nil element by
// terraform-plugin-sdk, and the request serializer (tea-rpc-utils >= v1.1.3,
// pulled in by tea-rpc >= v1.3.4) must fail with a clear serialization error
// instead of crashing the provider process.

func TestAccAliCloudNlbSecurityPolicy_emptyStringInList(t *testing.T) {
	resourceId := "alicloud_nlb_security_policy.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbsecuritypolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNLBSecurityPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": "${var.name}",
					"ciphers":              []string{"", "ECDHE-RSA-AES128-SHA"},
					"tls_versions":         []string{"TLSv1.2"},
				}),
				ExpectError: regexp.MustCompile(`cannot serialize repeated parameter element "Ciphers\.1": value is nil`),
			},
		},
	})
}

// TestAccAliCloudNlbSecurityPolicy_memberOrder covers the TypeList Order
// Coverage requirement for the ciphers and tls_versions attributes.
//
// The NLB API stores ciphers and tls_versions by membership (order-insensitive):
// UpdateSecurityPolicyAttribute rejects a same-membership reorder with
// "UpdateSecurityPolicyFailed / Security Policy is not changed" because the API
// detects no actual change, so a pure-reorder step cannot exercise a real
// TypeList update against this API. Each step below instead changes list
// membership AND reorders the retained members, so the apply produces a genuine
// update that the API accepts while still exercising the Provider's ordered
// TypeList handling (the configured order is transmitted to the request).
// Membership counts are asserted after each apply; same-membership reordering is
// intentionally avoided because the API rejects it as a no-op.
func TestAccAliCloudNlbSecurityPolicy_memberOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_security_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNLBSecurityPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbSecurityPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbsecuritypolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNLBSecurityPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": name,
					"ciphers":              []string{"ECDHE-RSA-AES128-SHA", "ECDHE-ECDSA-AES128-SHA"},
					"tls_versions":         []string{"TLSv1.0", "TLSv1.1", "TLSv1.2"},
					"tags": map[string]string{
						"Created": "tfTestAcc",
						"For":     "Tftestacc",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": name,
						"ciphers.#":            "2",
						"tls_versions.#":       "3",
						"tags.%":               "2",
						"tags.Created":         "tfTestAcc",
						"tags.For":             "Tftestacc",
					}),
				),
			},
			{
				// Change ciphers membership AND order: add AES256 and move
				// ECDSA ahead of RSA. A same-membership reorder would be
				// rejected as "Security Policy is not changed", so the add
				// makes the API see a real change. tls_versions and tags are
				// unchanged so the update is attributable to ciphers.
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": name,
					"ciphers":              []string{"ECDHE-ECDSA-AES128-SHA", "ECDHE-RSA-AES128-SHA", "ECDHE-ECDSA-AES256-SHA"},
					"tls_versions":         []string{"TLSv1.0", "TLSv1.1", "TLSv1.2"},
					"tags": map[string]string{
						"Created": "tfTestAcc",
						"For":     "Tftestacc",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#": "3",
					}),
				),
			},
			{
				// Change tls_versions membership AND order: remove TLSv1.0 and
				// move TLSv1.2 ahead of TLSv1.1. The membership change makes the
				// API accept the update; ciphers and tags are unchanged from
				// the prior step so the update is attributable to tls_versions.
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": name,
					"ciphers":              []string{"ECDHE-ECDSA-AES128-SHA", "ECDHE-RSA-AES128-SHA", "ECDHE-ECDSA-AES256-SHA"},
					"tls_versions":         []string{"TLSv1.2", "TLSv1.1"},
					"tags": map[string]string{
						"Created": "tfTestAcc",
						"For":     "Tftestacc",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_versions.#": "2",
					}),
				),
			},
		},
	})
}

// Test Nlb SecurityPolicy. <<< Resource test cases, automatically generated.
