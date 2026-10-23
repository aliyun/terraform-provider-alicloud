package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Cms AddonRelease. >>> Resource test cases, automatically generated.
// Case AddonRelease常规测试 8556
func TestAccAliCloudCmsAddonRelease_basic8556(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_addon_release.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAddonReleaseMap8556)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAddonRelease")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAddonReleaseBasicDependence8556)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"integration_policy_id": "${alicloud_cms_integration_policy.cloud.id}",
					"addon_name":            "cloud-acs-ecs",
					"addon_version":         "2.0.7",
					"workspace":             "${alicloud_cms_integration_policy.cloud.workspace}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"integration_policy_id": CHECKSET,
						"addon_name":            "cloud-acs-ecs",
						"addon_version":         "2.0.7",
						"workspace":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addon_version": "2.0.6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addon_version": "2.0.6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": `{\"discoverType\": \"StaticTarget\",\"discover\": {},\"targets\": \"127.0.0.1:8999\",\"metricPath\": \"/metrics\",\"scrapeInterval\": 15,\"honorTimestamps\": false,\"honorLabels\": false,\"queryParams\": {}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entity_rules": []map[string]interface{}{
						{
							"resource_group_id": "resourceGroupId",
							"instance_ids":      []string{"a", "b", "c", "d", "e"},
							"entity_types":      []string{"a", "b", "c", "d", "e"},
							"region_ids":        []string{"a", "b", "c", "d", "e"},
							"ip_match_rule": []map[string]interface{}{
								{
									"ip_field_key": "ipFieldKey",
									"ip_cidr":      "ipCIDR",
								},
							},
							"tags": []map[string]interface{}{
								{
									"op":         "in",
									"tag_key":    "TF",
									"tag_values": []string{"a", "b", "c", "d", "e"},
								},
								{
									"op":         "notIn",
									"tag_key":    "TFTag",
									"tag_values": []string{"f", "g", "h", "i", "j"},
								},
								{
									"op":         "equals",
									"tag_key":    "TFTagKey",
									"tag_values": []string{"k", "l", "m", "n", "o"},
								},
							},
							"field_rules": []map[string]interface{}{
								{
									"op":           "startWith",
									"field_key":    "TF",
									"field_values": []string{"a", "b", "c", "d", "e"},
								},
								{
									"op":           "equals",
									"field_key":    "TFField",
									"field_values": []string{"f", "g", "h", "i", "j"},
								},
								{
									"op":           "inAll",
									"field_key":    "TFFieldKey",
									"field_values": []string{"k", "l", "m", "n", "o"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entity_rules.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAliCloudCmsAddonRelease_basic8556_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_addon_release.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAddonReleaseMap8556)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAddonRelease")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAddonReleaseBasicDependence8556)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"integration_policy_id": "${alicloud_cms_integration_policy.cloud.id}",
					"addon_name":            "cloud-acs-ecs",
					"addon_version":         "2.0.7",
					"workspace":             "${alicloud_cms_integration_policy.cloud.workspace}",
					"addon_release_name":    name,
					"aliyun_lang":           "en",
					"env_type":              "Cloud",
					"config":                `{\"discoverType\": \"StaticTarget\",\"discover\": {},\"targets\": \"127.0.0.1:1111\",\"metricPath\": \"/metrics\",\"scrapeInterval\": 15,\"honorTimestamps\": false,\"honorLabels\": false,\"queryParams\": {},\"customRelabelConfigs\": \"\"}`,
					"entity_rules": []map[string]interface{}{
						{
							"resource_group_id": "resourceGroupId",
							"instance_ids":      []string{"a", "b", "c", "d", "e"},
							"entity_types":      []string{"a", "b", "c", "d", "e"},
							"region_ids":        []string{"a", "b", "c", "d", "e"},
							"ip_match_rule": []map[string]interface{}{
								{
									"ip_field_key": "ipFieldKey",
									"ip_cidr":      "ipCIDR",
								},
							},
							"tags": []map[string]interface{}{
								{
									"op":         "in",
									"tag_key":    "TF",
									"tag_values": []string{"a", "b", "c", "d", "e"},
								},
								{
									"op":         "notIn",
									"tag_key":    "TFTag",
									"tag_values": []string{"f", "g", "h", "i", "j"},
								},
								{
									"op":         "equals",
									"tag_key":    "TFTagKey",
									"tag_values": []string{"k", "l", "m", "n", "o"},
								},
							},
							"field_rules": []map[string]interface{}{
								{
									"op":           "startWith",
									"field_key":    "TF",
									"field_values": []string{"a", "b", "c", "d", "e"},
								},
								{
									"op":           "equals",
									"field_key":    "TFField",
									"field_values": []string{"f", "g", "h", "i", "j"},
								},
								{
									"op":           "inAll",
									"field_key":    "TFFieldKey",
									"field_values": []string{"k", "l", "m", "n", "o"},
								},
							},
						},
					},
					"dry_run": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"integration_policy_id": CHECKSET,
						"addon_name":            "cloud-acs-ecs",
						"addon_version":         "2.0.7",
						"workspace":             CHECKSET,
						"addon_release_name":    name,
						"aliyun_lang":           "en",
						"env_type":              "Cloud",
						"config":                CHECKSET,
						"entity_rules.#":        "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AliCloudCmsAddonReleaseMap8556 = map[string]string{
	"addon_release_name": CHECKSET,
	"aliyun_lang":        CHECKSET,
	"create_time":        CHECKSET,
	"region_id":          CHECKSET,
}

func AliCloudCmsAddonReleaseBasicDependence8556(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_integration_policy" "cloud" {
  policy_type             = "Cloud"
  integration_policy_name = var.name
  workspace               = alicloud_cms_workspace.default.id
}
`, name)
}

// Test Cms AddonRelease. <<< Resource test cases, automatically generated.
