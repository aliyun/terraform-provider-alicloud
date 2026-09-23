package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

// Case alicloud_ram_access_key_policy basic
func TestAccAliCloudRamAccessKeyPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_access_key_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudRamAccessKeyPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamAccessKeyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccramakp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRamAccessKeyPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_key_id":  "${alicloud_ram_access_key.default.id}",
					"user_principal_name": "${alicloud_ram_user.default.name}@${data.alicloud_account.default.id}.onaliyun.com",
					"access_key_policy":   "{\\\"Version\\\":1,\\\"Status\\\":\\\"Active\\\",\\\"Statements\\\":[{\\\"Type\\\":\\\"ClassicWhiteList\\\",\\\"IPList\\\":[\\\"10.0.0.1/32\\\"]}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_key_id":  CHECKSET,
						"user_principal_name": CHECKSET,
						"access_key_policy":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_key_policy": "{\\\"Version\\\":1,\\\"Status\\\":\\\"Active\\\",\\\"Statements\\\":[{\\\"Type\\\":\\\"ClassicWhiteList\\\",\\\"IPList\\\":[\\\"10.0.0.2/32\\\"]}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_key_policy": CHECKSET,
					}),
				),
			},
			{
				// Disable the policy while retaining the rule (Status Inactive, non-empty Statements).
				// Note: a disabled policy with no statements is the reset baseline and is treated as
				// "not exist", so a statement is retained here to keep the resource present.
				Config: testAccConfig(map[string]interface{}{
					"access_key_policy": "{\\\"Version\\\":1,\\\"Status\\\":\\\"Inactive\\\",\\\"Statements\\\":[{\\\"Type\\\":\\\"ClassicWhiteList\\\",\\\"IPList\\\":[\\\"10.0.0.2/32\\\"]}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_key_policy": CHECKSET,
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

var AlicloudRamAccessKeyPolicyMap0 = map[string]string{}

func AlicloudRamAccessKeyPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_ram_access_key" "default" {
  user_name = alicloud_ram_user.default.name
}
`, name)
}

func TestUnitRamAccessKeyPolicyStripVersion(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"removes version", `{"Version":1,"Status":"Active","Statements":[]}`, `{"Statements":[],"Status":"Active"}`},
		{"no version keeps document", `{"Status":"Active","Statements":[]}`, `{"Statements":[],"Status":"Active"}`},
		{"empty object", `{}`, `{}`},
		{"invalid json returned as-is", `not-json`, `not-json`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, stripAccessKeyPolicyVersion(c.in))
		})
	}
}

func TestUnitRamAccessKeyPolicyEquivalent(t *testing.T) {
	// Same policy, one with the server-managed Version field and one without.
	a := `{"Version":1,"Status":"Active","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.1/32"]}]}`
	b := `{"Status":"Active","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.1/32"]}]}`
	assert.True(t, accessKeyPolicyEquivalent(a, b), "documents differing only by Version should be equivalent")

	// Different IP list -> not equivalent.
	c := `{"Status":"Active","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.2/32"]}]}`
	assert.False(t, accessKeyPolicyEquivalent(a, c), "documents with different statements should not be equivalent")
}

func TestUnitRamAccessKeyPolicyIsEmpty(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"blank string", "   ", true},
		{"empty object", `{}`, true},
		{"disabled without statements", `{"Version":1,"Status":"Inactive","Statements":[]}`, true},
		{"disabled without statements field", `{"Status":"Inactive"}`, true},
		{"disabled with statements is not empty", `{"Status":"Inactive","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.1/32"]}]}`, false},
		{"active policy is not empty", `{"Status":"Active","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.1/32"]}]}`, false},
		{"invalid json is not empty", `not-json`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, isEmptyAccessKeyPolicy(c.in))
		})
	}
}

func TestUnitRamAccessKeyPolicyParseId(t *testing.T) {
	cases := []struct {
		name    string
		id      string
		wantUpn string
		wantAk  string
	}{
		{"composite id", "alice@1234567890.onaliyun.com:LTAI0123456789", "alice@1234567890.onaliyun.com", "LTAI0123456789"},
		{"access key only", "LTAI0123456789", "", "LTAI0123456789"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			upn, ak := parseRamAccessKeyPolicyId(c.id)
			assert.Equal(t, c.wantUpn, upn)
			assert.Equal(t, c.wantAk, ak)
		})
	}
}
