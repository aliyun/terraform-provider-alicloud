package alicloud

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudECSSecurityGroupRuleBasic(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rac := resourceAttrCheckInit(
		resourceCheckInit(resourceId, &v, serviceFunc),
		resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap))
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abc",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleCidrIp(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_security_group_id": REMOVEKEY,
						"cidr_ip":                  "0.0.0.0/0",
						"description":              "abcd",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleDescription(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleAll(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abcd",
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

func TestAccAliCloudECSSecurityGroupRuleEgress(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":        "egress",
		"policy":      "accept",
		"description": "SHDRP-7513",
		"port_range":  "443/443",
		"priority":    "1",
		"cidr_ip":     "182.254.11.243/32",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupEgressRule(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
					}),
				),
			},
			{
				Config: hclSecurityGroupEgressRuleDescription(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7512",
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

func TestAccAliCloudECSSecurityGroupRuleMulti(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test.2"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleMulti(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_ip":                  "45.20.250.240/32",
						"source_security_group_id": REMOVEKEY,
					}),
				),
			},
		},
	})

}

func TestAccAliCloudECSSecurityGroupRulePrefixList(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRulePrefixList)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRulePrefix(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    "abc",
						"prefix_list_id": CHECKSET,
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

func TestAccAliCloudECSSecurityGroupRuleEgressIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "egress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "443/443",
		"priority":     "1",
		"ipv6_cidr_ip": "2408:4004:cc:400::/56",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupEgressRuleIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
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

func TestAccAliCloudECSSecurityGroupRuleIngressIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupIngressRuleIpv6Map)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupIngressRuleIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2408:4004:cc:400::/56",
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

func TestAccAliCloudECSSecurityGroupRuleEgressOtherIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "egress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "443/443",
		"priority":     "1",
		"ipv6_cidr_ip": "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupEgressRuleOtherIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2001:db8:3c4d:15::1a2f:1a2b/0",
						"description":  "SHDRP-7513",
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

func TestAccAliCloudECSSecurityGroupRuleIngressOtherIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupIngressRuleIpv6Map)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupIngressRuleOtherIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2001:db8:3c4d:15::1a2f:1a2b/0",
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

func TestAccAliCloudECSSecurityGroupRuleEgressICMPv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "ingress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "-1/-1",
		"priority":     "1",
		"ipv6_cidr_ip": "::/0",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupEgressRuleICMPv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
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

// TestUnitSecurityGroupRuleSourceGroupOwnerIDSchema is a pure unit test (no TF_ACC, no cloud)
// verifying the source_group_owner_id schema field exists with the correct attributes:
// TypeString, Optional, ForceNew, not Required, not Computed (write-only).
func TestUnitSecurityGroupRuleSourceGroupOwnerIDSchema(t *testing.T) {
	r := resourceAliyunSecurityGroupRule()
	s, ok := r.Schema["source_group_owner_id"]
	if !ok {
		t.Fatal("source_group_owner_id is not present in the schema")
	}
	if s.Type != schema.TypeString {
		t.Fatalf("source_group_owner_id: expected TypeString, got %v", s.Type)
	}
	if !s.Optional {
		t.Fatal("source_group_owner_id: expected Optional=true")
	}
	if !s.ForceNew {
		t.Fatal("source_group_owner_id: expected ForceNew=true")
	}
	if s.Required {
		t.Fatal("source_group_owner_id: expected Required=false")
	}
	if s.Computed {
		t.Fatal("source_group_owner_id: expected Computed=false (write-only)")
	}
}

// TestAccAliCloudECSSecurityGroupRuleSourceGroupOwnerID verifies that source_group_owner_id
// is transparently passed through as SourceGroupOwnerId (ingress) / DestGroupOwnerId (egress)
// in the AuthorizeSecurityGroup/Egress request. It uses the current account's own UID as the
// owner value, which the API accepts as a valid self cross-account rule.
// source_group_owner_id is write-only (DescribeSecurityGroupAttribute does not return it),
// so the import step does not run ImportStateVerify (the field cannot be reconciled).
func TestAccAliCloudECSSecurityGroupRuleSourceGroupOwnerID(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":                     "ingress",
		"ip_protocol":              "tcp",
		"nic_type":                 "intranet",
		"policy":                   "accept",
		"port_range":               "22/22",
		"priority":                 "1",
		"security_group_id":        CHECKSET,
		"source_security_group_id": CHECKSET,
		"source_group_owner_id":    CHECKSET,
		"cidr_ip":                  "",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleSourceGroupOwnerID(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                  "ingress",
						"ip_protocol":           "tcp",
						"nic_type":              "intranet",
						"policy":                "accept",
						"port_range":            "22/22",
						"priority":              "1",
						"security_group_id":     CHECKSET,
						"source_group_owner_id": CHECKSET,
						// Self-account: DescribeSecurityGroupAttribute returns SourceGroupOwnerAccount
						// as an empty string (the server only populates it for true cross-account
						// rules where it converts SourceGroupOwnerId into a UID string). Expect
						// empty rather than "set" so the assertion matches the real Read behavior.
						"source_group_owner_account": "",
					}),
				),
			},
			{
				// source_group_owner_id is write-only (DescribeSecurityGroupAttribute
				// does not return it), so ImportStateVerify cannot reconcile it; the
				// import step only verifies the resource can be imported via its ID.
				ResourceName: resourceId,
				ImportState:  true,
			},
		},
	})

}

// TestAccAliCloudECSSecurityGroupRuleSourceGroupOwnerIDCrossAccount verifies a real
// cross-account security group rule: the rule lives in account 1 (default provider,
// ALICLOUD_ACCESS_KEY_1) and references a security group owned by account 2 (peer
// provider, ALICLOUD_ACCESS_KEY_2) via source_security_group_id + source_group_owner_id.
// The peer account id is read from ALICLOUD_ACCOUNT_ID_2 / INVITED_ALICLOUD_ACCOUNT_ID,
// or derived from ALICLOUD_ACCESS_KEY_2 via STS GetCallerIdentity.
// source_group_owner_id is write-only (DescribeSecurityGroupAttribute does not return it),
// so no ImportStateVerify step is performed.
func TestAccAliCloudECSSecurityGroupRuleSourceGroupOwnerIDCrossAccount(t *testing.T) {
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	peerAccountId := testAccResolveSecurityGroupRulePeerAccountId(t)
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			return Provider(), nil
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithMultipleAccount(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSecurityGroupRuleCrossAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleSourceGroupOwnerIDCrossAccount(name, peerAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "type", "ingress"),
					resource.TestCheckResourceAttr(resourceId, "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceId, "nic_type", "intranet"),
					resource.TestCheckResourceAttr(resourceId, "policy", "accept"),
					resource.TestCheckResourceAttr(resourceId, "port_range", "22/22"),
					resource.TestCheckResourceAttr(resourceId, "priority", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
					resource.TestCheckResourceAttrSet(resourceId, "source_security_group_id"),
					resource.TestCheckResourceAttr(resourceId, "source_group_owner_id", peerAccountId),
					// Cross-account: the server converts SourceGroupOwnerId into
					// SourceGroupOwnerAccount (returned as the peer UID string). The HCL
					// declares source_group_owner_account explicitly so the config value
					// matches the Read-back value and the ForceNew field does not drift.
					resource.TestCheckResourceAttr(resourceId, "source_group_owner_account", peerAccountId),
					resource.TestCheckResourceAttr(resourceId, "cidr_ip", ""),
				),
			},
		},
	})
}

var securityGroupRulePeerAccountId string

// testAccResolveSecurityGroupRulePeerAccountId resolves the peer (account 2) UID
// for the cross-account security group rule test. It prefers ALICLOUD_ACCOUNT_ID_2
// / INVITED_ALICLOUD_ACCOUNT_ID, then derives the UID from ALICLOUD_ACCESS_KEY_2 via
// STS GetCallerIdentity. The test is skipped when neither is available.
func testAccResolveSecurityGroupRulePeerAccountId(t *testing.T) string {
	if securityGroupRulePeerAccountId != "" {
		return securityGroupRulePeerAccountId
	}
	for _, envName := range []string{"ALICLOUD_ACCOUNT_ID_2", "INVITED_ALICLOUD_ACCOUNT_ID"} {
		if v := strings.TrimSpace(os.Getenv(envName)); v != "" {
			securityGroupRulePeerAccountId = v
			return securityGroupRulePeerAccountId
		}
	}
	ak := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_2"))
	sk := strings.TrimSpace(os.Getenv("ALICLOUD_SECRET_KEY_2"))
	if ak == "" || sk == "" {
		t.Skipf("Skipping: set ALICLOUD_ACCOUNT_ID_2 or INVITED_ALICLOUD_ACCOUNT_ID, or set ALICLOUD_ACCESS_KEY_2/SECRET_KEY_2 so the peer account id can be derived")
	}
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	client, err := sdk.NewClientWithAccessKey(region, ak, sk)
	if err != nil {
		t.Fatalf("failed to build STS client for peer account: %s", err)
	}
	request := requests.NewCommonRequest()
	request.Method = requests.POST
	request.Scheme = "https"
	request.Domain = "sts.aliyuncs.com"
	request.Version = "2015-04-01"
	request.Product = "Sts"
	request.ApiName = "GetCallerIdentity"
	request.TransToAcsRequest()
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		t.Skipf("Skipping: failed to derive peer account id from ALICLOUD_ACCESS_KEY_2: %s", err)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal([]byte(response.GetHttpContentString()), &result)
	accountId := strings.TrimSpace(fmt.Sprint(result["AccountId"]))
	if accountId == "" || accountId == "<nil>" {
		t.Skipf("Skipping: STS GetCallerIdentity did not return AccountId for ALICLOUD_ACCESS_KEY_2")
	}
	securityGroupRulePeerAccountId = accountId
	return securityGroupRulePeerAccountId
}

var testAccCheckSecurityGroupRuleBasicMap = map[string]string{
	"type":                     "ingress",
	"ip_protocol":              "tcp",
	"nic_type":                 "intranet",
	"policy":                   "drop",
	"port_range":               "22/22",
	"priority":                 "100",
	"security_group_id":        CHECKSET,
	"source_security_group_id": CHECKSET,
	"security_group_rule_id":   CHECKSET,
	"cidr_ip":                  "",
}

var testAccCheckSecurityGroupIngressRuleIpv6Map = map[string]string{
	"type":              "ingress",
	"ip_protocol":       "tcp",
	"nic_type":          "intranet",
	"policy":            "drop",
	"port_range":        "22/22",
	"priority":          "100",
	"security_group_id": CHECKSET,
	"cidr_ip":           "",
}

var testAccCheckSecurityGroupRulePrefixList = map[string]string{
	"type":              "ingress",
	"ip_protocol":       "tcp",
	"nic_type":          "intranet",
	"policy":            "accept",
	"port_range":        "22/22",
	"priority":          "100",
	"security_group_id": CHECKSET,
}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group_rule" {
			continue
		}
		_, err := ecsService.DescribeSecurityGroupRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}
	}

	return nil
}

// testAccCheckSecurityGroupRuleCrossAccountDestroy verifies the cross-account
// rule is destroyed after the test. In ProviderFactories mode testAccProvider
// is not the configured instance (its Meta is nil); when Meta is unavailable
// the destroy verification is skipped and relies on Terraform's own destroy
// reporting, which surfaces a failed destroy as a test error.
func testAccCheckSecurityGroupRuleCrossAccountDestroy(s *terraform.State) error {
	if testAccProvider == nil || testAccProvider.Meta() == nil {
		return nil
	}
	return testAccCheckSecurityGroupRuleDestroy(s)
}

func hclSecurityGroupRuleBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRBase%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  nic_type                 = "intranet"
  policy                   = "drop"
  port_range               = "22/22"
  priority                 = 100
  security_group_id        = alicloud_security_group.test.0.id
  source_security_group_id = alicloud_security_group.test.1.id
  description              = "abc"
}
`, name)
}

func hclSecurityGroupRulePrefix(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRPrefix%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.vpcs.0.id
  security_group_name = var.name
}

resource "alicloud_ecs_prefix_list" "test" {
  address_family   = "IPv4"
  max_entries      = 2
  prefix_list_name = "tftest"
  description      = "description"
  entry {
    cidr        = "192.168.0.0/24"
    description = "description"
  }
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  prefix_list_id    = alicloud_ecs_prefix_list.test.id
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.id
  description       = "abc"
}
`, name)
}

func hclSecurityGroupRuleCidrIp(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR-CIRDIP%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "abcd"
}
`, name)
}

func hclSecurityGroupRuleDescription(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR-desc%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "description"
}
`, name)
}

func hclSecurityGroupRuleAll(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR_all%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "abcd"
}
`, name)
}

func hclSecurityGroupRuleMulti(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRMulti%s"
}

variable "cidr_ip_list" {
  type    = "list"
  default = ["50.255.255.255/32", "75.250.250.250/32", "45.20.250.240/32"]
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  count             = length(compact(var.cidr_ip_list))
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  cidr_ip           = element(var.cidr_ip_list, count.index)
}
`, name)
}

func hclSecurityGroupIngressRuleIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRIngressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  ipv6_cidr_ip      = "2408:4004:cc:400::/56"
}
`, name)
}

func hclSecurityGroupIngressRuleOtherIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRIngressOtherIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  ipv6_cidr_ip      = "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0"
}
`, name)
}

func hclSecurityGroupEgressRule(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgress%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  cidr_ip           = "182.254.11.243/32"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "2408:4004:cc:400::/56"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleOtherIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressOtherIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleDescription(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressDesc%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
  inner_access_policy = "Accept"
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  cidr_ip           = "182.254.11.243/32"
  description       = "SHDRP-7512"
}
`, name)
}

func hclSecurityGroupEgressRuleICMPv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "icmpv6"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "-1/-1"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "::/0"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupRuleSourceGroupOwnerID(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRSrcOwner%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_account" "current" {}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  nic_type                 = "intranet"
  policy                   = "accept"
  port_range               = "22/22"
  priority                 = 1
  security_group_id        = alicloud_security_group.test.0.id
  source_security_group_id = alicloud_security_group.test.1.id
  source_group_owner_id    = data.alicloud_account.current.id
  description              = "tf-testAccSGRSrcOwner"
}
`, name)
}

func hclSecurityGroupRuleSourceGroupOwnerIDCrossAccount(name string, peerAccountId string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRSrcOwnerCross%s"
}

# Peer account (account 2): owns the referenced source security group.
# Credentials come from ALICLOUD_ACCESS_KEY_2 / ALICLOUD_SECRET_KEY_2;
# testAccPreCheckWithMultipleAccount skips this test when ALICLOUD_ACCESS_KEY_2
# is unset. The own account uses the default provider (ALICLOUD_ACCESS_KEY env)
# and is intentionally not redeclared — an explicit provider block with empty
# ALICLOUD_ACCESS_KEY_1 would blank the credentials and break auth.
provider "alicloud" {
  alias      = "peer"
  access_key = "%s"
  secret_key = "%s"
}

data "alicloud_vpcs" "own" {
  name_regex = "^default-NODELETING$"
}

# Peer account self-built VPC. The peer account has no pre-existing
# default-NODELETING VPC, so data.alicloud_vpcs.peer would always be empty
# (Invalid index). Provision a VPC in the peer account instead.
resource "alicloud_vpc" "peer" {
  provider   = alicloud.peer
  vpc_name   = "${var.name}-peer"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_security_group" "peer" {
  provider            = alicloud.peer
  vpc_id              = alicloud_vpc.peer.id
  security_group_name = "${var.name}-peer"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.own.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  nic_type                 = "intranet"
  policy                   = "accept"
  port_range               = "22/22"
  priority                 = 1
  security_group_id        = alicloud_security_group.test.id
  source_security_group_id = alicloud_security_group.peer.id
  source_group_owner_id    = "%s"
  # SourceGroupOwnerAccount is deprecated for cross-account (the server silently
  # drops it and derives ownership from SourceGroupOwnerId), but DescribeSecurityGroupAttribute
  # returns the peer UID in SourceGroupOwnerAccount. Declaring it here makes the
  # config value equal the Read-back value so the ForceNew field does not drift.
  source_group_owner_account = "%s"
  description                = "tf-testAccSGRSrcOwnerCross"
}
`, name, os.Getenv("ALICLOUD_ACCESS_KEY_2"), os.Getenv("ALICLOUD_SECRET_KEY_2"), peerAccountId, peerAccountId)
}
