package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_digital_employee", &resource.Sweeper{
		Name: "alicloud_cms_digital_employee",
		F:    testSweepCmsDigitalEmployee,
	})
}

func testSweepCmsDigitalEmployee(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	prefixes := []string{
		"tf-acc-cms",
		"tfacccms",
	}

	query := make(map[string]*string)
	query["maxResults"] = StringPointer(fmt.Sprintf("%d", PageSizeLarge))
	objects, _, err := cmsServiceV2.ListCmsDigitalEmployees(query)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve Cms DigitalEmployee list: %s", err)
		return nil
	}

	for _, item := range objects {
		name := fmt.Sprintf("%v", item["name"])
		skip := true
		for _, prefix := range prefixes {
			if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
				skip = false
				break
			}
		}
		if skip {
			continue
		}
		log.Printf("[INFO] delete cms digital employee: %s", name)
		action := fmt.Sprintf("/digital-employee/%s", name)
		_, err := client.RoaDelete("Cms", "2024-03-30", action, make(map[string]*string), nil, nil, true)
		if err != nil {
			log.Printf("[ERROR] Failed to delete cms digital employee (%s): %s", name, err)
		}
	}

	return nil
}

func testAccCmsDigitalEmployeeDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ram_role" "default" {
  name     = "tf-acc-cms-de-role-%[1]s"
  document = <<EOF
{
  "Version": "1",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "RAM": ["acs:ram::1632449562862807:root"]
      }
    }
  ]
}
EOF
}

resource "alicloud_ram_role" "second" {
  name     = "tf-acc-cms-de-role-2-%[1]s"
  document = <<EOF
{
  "Version": "1",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "RAM": ["acs:ram::1632449562862807:root"]
      }
    }
  ]
}
EOF
}
`, name)
}

// hclEscapeString escapes double quotes in a string so it can be safely embedded
// into an HCL double-quoted string produced by testAccConfig's valueConvert,
// which renders string values via fmt.Sprintf("\"%s\"", v) without escaping.
// Without this, JSON strings containing double quotes (e.g. bailian attributes)
// would prematurely close the HCL string and break config parsing.
func hclEscapeString(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

func TestAccAliCloudCmsDigitalEmployee_basic0(t *testing.T) {
	resourceId := "alicloud_cms_digital_employee.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCmsDigitalEmployeeDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCmsDigitalEmployeeDestroy(resourceId),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"digital_employee_name": name,
					"role_arn":              "${alicloud_ram_role.default.arn}",
					"description":           "description-basic-1",
					"display_name":          "display-basic-1",
					"default_rule":          "rule-default",
					"resource_group_id":     "",
					"attributes": map[string]interface{}{
						"k1": "v1",
						"k2": "v2",
					},
					"knowledges": []map[string]interface{}{
						{
							"bailian": []map[string]interface{}{
								{
									"workspace_id": "ws-tfacc-" + name,
									"index_id":     "idx-tfacc-" + name,
									"region":       "cn-beijing",
									"attributes":   hclEscapeString(`{"source":"tfacc"}`),
								},
							},
						},
					},
					"sandbox_network_policy": []map[string]interface{}{
						{
							"allow_fqdns": []string{"api.example.com", "cdn.example.com"},
							"allow_cidrs": []string{"10.0.0.0/8", "192.168.1.1"},
							"enable_acl":  false,
						},
					},
					"tool_policy": []map[string]interface{}{
						{
							"aliyun": []map[string]interface{}{
								{
									"enable":           true,
									"deny_policy":      []string{"ecs:Delete*"},
									"auto_pass_policy": []string{"log:Get*", "log:List*"},
									"statements": []map[string]interface{}{
										{
											"decision":    "user_ack",
											"product":     "Sls",
											"api_version": "2020-12-30",
											"actions":     []string{"log:GetProject", "log:CreateDashboard"},
										},
									},
								},
							},
						},
					},
					"tags": []map[string]interface{}{
						{
							"key":   "env",
							"value": "tf-acc",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "digital_employee_name", name),
					resource.TestCheckResourceAttrSet(resourceId, "role_arn"),
					resource.TestCheckResourceAttr(resourceId, "description", "description-basic-1"),
					resource.TestCheckResourceAttr(resourceId, "display_name", "display-basic-1"),
					resource.TestCheckResourceAttr(resourceId, "default_rule", "rule-default"),
					resource.TestCheckResourceAttr(resourceId, "attributes.k1", "v1"),
					resource.TestCheckResourceAttr(resourceId, "attributes.k2", "v2"),
					resource.TestCheckResourceAttr(resourceId, "sandbox_network_policy.0.enable_acl", "false"),
					resource.TestCheckResourceAttr(resourceId, "sandbox_network_policy.0.allow_fqdns.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "knowledges.0.bailian.0.workspace_id", "ws-tfacc-"+name),
					resource.TestCheckResourceAttr(resourceId, "tool_policy.0.aliyun.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "tags.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttrSet(resourceId, "update_time"),
					resource.TestCheckResourceAttrSet(resourceId, "region_id"),
					resource.TestCheckResourceAttr(resourceId, "resource_type", "ALIYUN::CMS::DIGITALEMPLOYEE"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_arn":          "${alicloud_ram_role.second.arn}",
					"description":       "description-basic-2-updated",
					"display_name":      "display-basic-2-updated",
					"default_rule":      "rule-updated",
					"resource_group_id": "rg-acctest-modified",
					"attributes": map[string]interface{}{
						"k1": "v1-updated",
						"k3": "v3",
					},
					"knowledges": []map[string]interface{}{
						{
							"bailian": []map[string]interface{}{
								{
									"workspace_id": "ws-tfacc-upd-" + name,
									"index_id":     "idx-tfacc-upd-" + name,
									"region":       "cn-shanghai",
									"attributes":   hclEscapeString(`{"source":"tfacc-updated"}`),
								},
							},
						},
					},
					"sandbox_network_policy": []map[string]interface{}{
						{
							"allow_fqdns": []string{"api2.example.com"},
							"allow_cidrs": []string{"172.16.0.0/12"},
							"enable_acl":  true,
						},
					},
					"tool_policy": []map[string]interface{}{
						{
							"aliyun": []map[string]interface{}{
								{
									"enable":           false,
									"deny_policy":      []string{"ecs:Delete*", "vpc:DeleteVpc"},
									"auto_pass_policy": []string{"log:Get*", "log:List*", "oss:Get*"},
									"statements": []map[string]interface{}{
										{
											"decision":    "auto_pass",
											"product":     "Ecs",
											"api_version": "2014-05-26",
											"actions":     []string{"ecs:DescribeInstances"},
										},
									},
								},
							},
						},
					},
					"tags": []map[string]interface{}{
						{
							"key":   "env-updated",
							"value": "tf-acc-updated",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "description-basic-2-updated"),
					resource.TestCheckResourceAttr(resourceId, "display_name", "display-basic-2-updated"),
					resource.TestCheckResourceAttr(resourceId, "default_rule", "rule-updated"),
					resource.TestCheckResourceAttr(resourceId, "digital_employee_name", name),
					resource.TestCheckResourceAttr(resourceId, "sandbox_network_policy.0.enable_acl", "true"),
					resource.TestCheckResourceAttr(resourceId, "tool_policy.0.aliyun.0.enable", "false"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tool_policy", "attributes"},
			},
		},
	})
}

func TestAccAliCloudCmsDigitalEmployee_basic0_nonExist(t *testing.T) {
	resourceId := "alicloud_cms_digital_employee.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCmsDigitalEmployeeDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCmsDigitalEmployeeDestroy(resourceId),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"digital_employee_name": name,
					"role_arn":              "${alicloud_ram_role.default.arn}",
					"description":           "description-nonexist",
					"display_name":          "display-nonexist",
					"default_rule":          "rule-default",
					"resource_group_id":     "",
					"attributes": map[string]interface{}{
						"k1": "v1",
					},
					"sandbox_network_policy": []map[string]interface{}{
						{
							"allow_fqdns": []string{"api.example.com"},
							"allow_cidrs": []string{"10.0.0.0/8"},
							"enable_acl":  false,
						},
					},
					"tags": []map[string]interface{}{
						{
							"key":   "env",
							"value": "tf-acc",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "digital_employee_name", name),
				),
			},
			{
				Config:  testAccConfig(map[string]interface{}{}),
				Destroy: true,
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckCmsDigitalEmployeeDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cmsServiceV2 := CmsServiceV2{client}
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return nil
		}
		if rs.Primary.ID == "" {
			return nil
		}
		_, err := cmsServiceV2.DescribeCmsDigitalEmployee(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
		return WrapError(Error("cms digital employee %s still exists", rs.Primary.ID))
	}
}
