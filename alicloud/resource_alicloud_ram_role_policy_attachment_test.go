package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Test Ram RolePolicyAttachment. >>> Resource test cases, automatically generated.
// Case  RolePolicyAttachment测试 9050
func TestAccAliCloudRamRolePolicyAttachment_basic9050(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role_policy_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudRamRolePolicyAttachmentMap9050)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRolePolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRolePolicyAttachmentBasicDependence9050)
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
					"role_name":   "${alicloud_ram_role.default.id}",
					"policy_name": "${alicloud_ram_policy.default.id}",
					"policy_type": "Custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":   CHECKSET,
						"policy_name": CHECKSET,
						"policy_type": "Custom",
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

func TestAccAliCloudRamRolePolicyAttachment_basic9051(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role_policy_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudRamRolePolicyAttachmentMap9050)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRolePolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRolePolicyAttachmentBasicDependence9050)
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
					"role_name":   "${alicloud_ram_role.default.id}",
					"policy_name": "AliyunECSReadOnlyAccess",
					"policy_type": "System",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":   CHECKSET,
						"policy_name": "AliyunECSReadOnlyAccess",
						"policy_type": "System",
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

var AliCloudRamRolePolicyAttachmentMap9050 = map[string]string{}

func AliCloudRamRolePolicyAttachmentBasicDependence9050(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_ram_policy" "default" {
	  name = var.name
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role" "default" {
	  name = var.name
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	}
`, name)
}

// Test Ram RolePolicyAttachment. <<< Resource test cases, automatically generated.

func TestAccAliCloudRamRolePolicyAttachment_resourceGroupScopes(t *testing.T) {
	resourceId := "alicloud_ram_role_policy_attachment.default"
	accountResourceId := "alicloud_ram_role_policy_attachment.account"
	scopeAResourceId := "alicloud_ram_role_policy_attachment.scope_a"
	name := fmt.Sprintf("tf-testacc-ram-scope-%d", acctest.RandIntRange(1000000, 9999999))
	accountID, accountAlias := "", ""
	if os.Getenv("TF_ACC") != "" {
		testAccPreCheck(t)
		rawClient, err := sharedClientForRegion(os.Getenv("ALICLOUD_REGION"))
		if err != nil {
			t.Fatal(err)
		}
		client := rawClient.(*connectivity.AliyunClient)
		accountID, err = client.AccountId()
		if err != nil {
			t.Fatal(err)
		}
		response, err := client.RpcPost("ims", "2019-08-15", "GetDefaultDomain", nil, map[string]interface{}{}, true)
		if err != nil {
			t.Fatal(err)
		}
		domain, ok := response["DefaultDomainName"].(string)
		if !ok || !strings.HasSuffix(domain, ".onaliyun.com") {
			t.Fatal("invalid RAM default domain")
		}
		accountAlias = strings.TrimSuffix(domain, ".onaliyun.com")
		if accountAlias == "" {
			t.Fatal("empty RAM account alias")
		}
	}
	allScopes := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return RamRolePolicyAttachmentScopeDependence(name, true, true)
	})
	withoutA := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return RamRolePolicyAttachmentScopeDependence(name, false, true)
	})
	singleAttachment := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return RamRolePolicyAttachmentScopeDependence(name, false, false)
	})
	var accountIDInitial, scopeAID, scopeBID, movedID, accountIDFinal string
	ids := map[string]bool{}
	check := func(resourceName string, saved *string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok || rs.Primary.ID == "" {
				return fmt.Errorf("attachment %s missing from state", resourceName)
			}
			if *saved != "" && *saved != rs.Primary.ID {
				return fmt.Errorf("attachment %s unexpectedly changed ID", resourceName)
			}
			*saved = rs.Primary.ID
			ids[*saved] = true
			return testAccRamRolePolicyAttachmentRawScope(*saved, accountID, accountAlias, true)
		}
	}
	absent := func(id *string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			if *id == "" {
				return fmt.Errorf("missing saved attachment ID")
			}
			return testAccRamRolePolicyAttachmentRawScope(*id, accountID, accountAlias, false)
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders, IDRefreshName: resourceId,
		CheckDestroy: func(s *terraform.State) error {
			for _, rs := range s.RootModule().Resources {
				if rs.Type == "alicloud_ram_role_policy_attachment" {
					ids[rs.Primary.ID] = true
				}
			}
			for id := range ids {
				if id != "" {
					if err := testAccRamRolePolicyAttachmentRawScope(id, accountID, accountAlias, false); err != nil {
						return err
					}
				}
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: allScopes(map[string]interface{}{
					"policy_name": "${alicloud_ram_policy.scope.policy_name}", "policy_type": "Custom",
					"role_name":         "${alicloud_ram_role.scope.role_name}",
					"resource_group_id": "${alicloud_resource_manager_resource_group.scope_b.id}",
				}),
				Check: resource.ComposeTestCheckFunc(check(accountResourceId, &accountIDInitial), check(scopeAResourceId, &scopeAID), check(resourceId, &scopeBID),
					resource.TestCheckResourceAttr(accountResourceId, "resource_group_id", ""),
					resource.TestCheckResourceAttrPair(scopeAResourceId, "resource_group_id", "alicloud_resource_manager_resource_group.scope_a", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "alicloud_resource_manager_resource_group.scope_b", "id"),
					func(s *terraform.State) error {
						if len(strings.Split(accountIDInitial, ":")) != 4 || len(strings.Split(scopeAID, ":")) != 5 || len(strings.Split(scopeBID, ":")) != 5 || scopeAID == scopeBID {
							return fmt.Errorf("account/resource-group IDs are not distinct and importable")
						}
						return nil
					}),
			},
			{ResourceName: accountResourceId, ImportState: true, ImportStateVerify: true},
			{ResourceName: scopeAResourceId, ImportState: true, ImportStateVerify: true},
			{ResourceName: resourceId, ImportState: true, ImportStateVerify: true},
			{
				Config: withoutA(map[string]interface{}{
					"policy_name": "${alicloud_ram_policy.scope.policy_name}", "policy_type": "Custom",
					"role_name":         "${alicloud_ram_role.scope.role_name}",
					"resource_group_id": "${alicloud_resource_manager_resource_group.scope_b.id}",
				}),
				Check: resource.ComposeTestCheckFunc(absent(&scopeAID), check(accountResourceId, &accountIDInitial), check(resourceId, &scopeBID)),
			},
			{
				Config: withoutA(map[string]interface{}{
					"policy_name": "${alicloud_ram_policy.scope.policy_name}", "policy_type": "Custom",
					"role_name":         "${alicloud_ram_role.scope.role_name}",
					"resource_group_id": "${alicloud_resource_manager_resource_group.scope_a.id}",
				}),
				Check: resource.ComposeTestCheckFunc(absent(&scopeBID), check(accountResourceId, &accountIDInitial), check(resourceId, &movedID),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "alicloud_resource_manager_resource_group.scope_a", "id")),
			},
			{
				Config: singleAttachment(map[string]interface{}{
					"policy_name": "${alicloud_ram_policy.scope.policy_name}", "policy_type": "Custom",
					"role_name":         "${alicloud_ram_role.scope.role_name}",
					"resource_group_id": "${alicloud_resource_manager_resource_group.scope_a.id}",
				}),
				Check: resource.ComposeTestCheckFunc(absent(&accountIDInitial), check(resourceId, &movedID)),
			},
			{
				Config: singleAttachment(map[string]interface{}{
					"policy_name": "${alicloud_ram_policy.scope.policy_name}", "policy_type": "Custom",
					"role_name": "${alicloud_ram_role.scope.role_name}", "resource_group_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(absent(&movedID), check(resourceId, &accountIDFinal),
					resource.TestCheckResourceAttr(resourceId, "resource_group_id", ""),
					func(s *terraform.State) error {
						if accountIDFinal != accountIDInitial {
							return fmt.Errorf("clearing scope did not restore the legacy account ID")
						}
						return nil
					}),
			},
			{
				Config: singleAttachment(map[string]interface{}{
					"policy_name": "${alicloud_ram_policy.scope.policy_name}", "policy_type": "Custom",
					"role_name": "${alicloud_ram_role.scope.role_name}",
				}),
				PlanOnly: true,
			},
		},
	})
}

func testAccRamRolePolicyAttachmentRawScope(id, accountID, accountAlias string, present bool) error {
	parts := strings.Split(id, ":")
	if len(parts) != 4 && len(parts) != 5 {
		return fmt.Errorf("invalid saved attachment ID")
	}
	scope := accountID
	if len(parts) == 5 {
		scope = parts[4]
	}
	request := map[string]interface{}{
		"PolicyName": parts[1], "PolicyType": parts[2], "PrincipalType": "ServiceRole",
		"PrincipalName": parts[3] + "@role." + accountAlias + ".onaliyunservice.com", "ResourceGroupId": scope,
	}
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPost("ResourceManager", "2020-03-31", "ListPolicyAttachments", nil, request, true)
		if err != nil {
			if !present && IsExpectedErrors(err, []string{"EntityNotExists.ResourceGroup", "EntityNotExist.Policy"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		container, ok := response["PolicyAttachments"].(map[string]interface{})
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("missing policy attachments response"))
		}
		items, ok := container["PolicyAttachment"].([]interface{})
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("invalid policy attachments list"))
		}
		if !present {
			if len(items) == 0 {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("detached scope still has an attachment"))
		}
		if len(items) != 1 {
			return resource.RetryableError(fmt.Errorf("expected one attachment, got %d", len(items)))
		}
		item, ok := items[0].(map[string]interface{})
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("invalid policy attachment"))
		}
		for field, expected := range request {
			if actual, ok := item[field].(string); !ok || actual != expected {
				return resource.NonRetryableError(fmt.Errorf("attachment returned a different %s", field))
			}
		}
		return nil
	})
}

func RamRolePolicyAttachmentScopeDependence(name string, includeA, includeAccount bool) string {
	config := fmt.Sprintf(`
variable "name" { default = %q }
data "alicloud_account" "scope" {}
resource "alicloud_ram_role" "scope" {
  role_name = var.name
  assume_role_policy_document = jsonencode({
    Version = "1"
    Statement = [{ Action = "sts:AssumeRole", Effect = "Allow", Principal = { RAM = [format("acs:ram::%%s:root", data.alicloud_account.scope.id)] } }]
  })
}
resource "alicloud_ram_policy" "scope" {
  policy_name = var.name
  policy_document = jsonencode({
    Version = "1"
    Statement = [{ Action = ["ecs:DescribeInstances"], Effect = "Allow", Resource = ["*"] }]
  })
}
resource "alicloud_resource_manager_resource_group" "scope_a" {
  resource_group_name = "${var.name}-a"
  display_name = "${var.name}-a"
}
resource "alicloud_resource_manager_resource_group" "scope_b" {
  resource_group_name = "${var.name}-b"
  display_name = "${var.name}-b"
}
`, name)
	if includeA {
		config += `
resource "alicloud_ram_role_policy_attachment" "scope_a" {
  policy_name = alicloud_ram_policy.scope.policy_name
  policy_type = "Custom"
  role_name = alicloud_ram_role.scope.role_name
  resource_group_id = alicloud_resource_manager_resource_group.scope_a.id
}
`
	}
	if includeAccount {
		config += `
resource "alicloud_ram_role_policy_attachment" "account" {
  policy_name = alicloud_ram_policy.scope.policy_name
  policy_type = "Custom"
  role_name = alicloud_ram_role.scope.role_name
}
`
	}
	return config
}
