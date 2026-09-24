package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test IaCService Module. >>> Resource test cases, automatically generated.
// Case Module lifecycle test
func TestAccAliCloudIacServiceModule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_iac_service_module.default"
	ra := resourceAttrInit(resourceId, AlicloudIacServiceModuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &IacServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeIacServiceModule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	projectId, groupId, groupIdUpdated, emptyProjectId := prepareIacServiceModuleGroupFixtures(t)
	name := fmt.Sprintf("tf-testacc%siacmodule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudIacServiceModuleBasicDependence0)
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
					"module_name":      name,
					"source":           "Registry",
					"source_path":      "alibaba/security-group:2.4.1",
					"version_strategy": "Manual",
					"description":      "tf-testacc module",
					"group_info": []map[string]interface{}{
						{
							"group_id":   groupId,
							"project_id": projectId,
						},
					},
					"tags": []map[string]interface{}{
						{
							"tag_key":   "Created",
							"tag_value": "TF",
						},
						{
							"tag_key":   "For",
							"tag_value": "Test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"module_name":               name,
						"source":                    "Registry",
						"source_path":               "alibaba/security-group:2.4.1",
						"version_strategy":          "Manual",
						"description":               "tf-testacc module",
						"group_info.#":              "1",
						"group_info.0.group_id":     groupId,
						"group_info.0.project_id":   projectId,
						"group_info.0.group_name":   CHECKSET,
						"group_info.0.project_name": CHECKSET,
						"tags.#":                    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"module_name":      name,
					"source":           "Registry",
					"source_path":      "alibaba/security-group:2.4.1",
					"version_strategy": "Manual",
					"description":      "tf-testacc module",
					"group_info": []map[string]interface{}{
						{
							"group_id":   groupIdUpdated,
							"project_id": emptyProjectId,
						},
					},
					"tags": []map[string]interface{}{
						{
							"tag_key":   "Created",
							"tag_value": "TF",
						},
						{
							"tag_key":   "For",
							"tag_value": "Test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_info.#":              "1",
						"group_info.0.group_id":     groupIdUpdated,
						"group_info.0.project_id":   emptyProjectId,
						"group_info.0.group_name":   CHECKSET,
						"group_info.0.project_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"module_name": name + "_update",
					"description": "tf-testacc module update",
					"source_path": "terraform-alicloud-modules/mongodb:3.0.0",
					"state_path":  "oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform.tfstate",
					"tags": []map[string]interface{}{
						{
							"tag_key":   "Created-update",
							"tag_value": "TF-update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"module_name": name + "_update",
						"description": "tf-testacc module update",
						"source_path": "terraform-alicloud-modules/mongodb:3.0.0",
						"state_path":  "oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform.tfstate",
						"tags.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_strategy": "SourcePathUpdated",
					"state_path":       "oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform-update.tfstate",
					"tags":             REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_strategy": "SourcePathUpdated",
						"state_path":       "oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform-update.tfstate",
						"tags.#":           "0",
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

var AlicloudIacServiceModuleMap0 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"output_path": CHECKSET,
}

func AlicloudIacServiceModuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// prepareIacServiceModuleGroupFixtures creates a project and two groups used to
// verify that group_info can be specified on create and changed on update, plus
// an empty project used to verify filters that match nothing. The fixtures are
// deleted when the test finishes. It runs before resource.Test, so it resolves
// the test region the same way testAccPreCheck does.
func prepareIacServiceModuleGroupFixtures(t *testing.T) (projectId string, groupId string, groupIdUpdated string, emptyProjectId string) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
		os.Setenv("ALICLOUD_REGION", region)
	}
	defaultRegionToTest = region
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Fatalf("Error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	rand := acctest.RandIntRange(10000, 99999)
	query := make(map[string]*string)

	newProject := func(suffix string) string {
		projectResp, err := client.RoaPost("IaCService", "2021-08-06", "/project", query, nil, map[string]interface{}{
			"name":        fmt.Sprintf("tf-acc-project-%d-%s", rand, suffix),
			"clientToken": buildClientToken("/project"),
		}, true)
		if err != nil {
			t.Fatalf("Creating IaCService project fixture failed: %s", err)
		}
		return fmt.Sprint(projectResp["projectId"])
	}
	projectId = newProject("a")
	emptyProjectId = newProject("empty")

	newGroup := func(suffix, projId string) string {
		groupResp, err := client.RoaPost("IaCService", "2021-08-06", "/group", query, nil, map[string]interface{}{
			"name":        fmt.Sprintf("tf-acc-group-%d-%s", rand, suffix),
			"projectId":   projId,
			"clientToken": buildClientToken("/group"),
		}, true)
		if err != nil {
			t.Fatalf("Creating IaCService group fixture failed: %s", err)
		}
		return fmt.Sprint(groupResp["groupId"])
	}
	groupId = newGroup("a", projectId)
	// groupIdUpdated is created in emptyProjectId so the resource update step
	// moves the module into a group of a different project, exercising both
	// group_id and project_id modification (the testing coverage rate check
	// requires project_id to change across steps). emptyProjectId still holds
	// no modules, so datasource filter tests that expect no match are unaffected.
	groupIdUpdated = newGroup("b", emptyProjectId)

	t.Cleanup(func() {
		for _, id := range []string{groupId, groupIdUpdated} {
			if _, err := client.RoaDelete("IaCService", "2021-08-06", fmt.Sprintf("/group/%s", id), query, nil, nil, true); err != nil {
				t.Logf("[WARN] deleting IaCService group fixture %s failed: %s", id, err)
			}
		}
		for _, id := range []string{projectId, emptyProjectId} {
			if _, err := client.RoaDelete("IaCService", "2021-08-06", fmt.Sprintf("/project/%s", id), query, nil, nil, true); err != nil {
				t.Logf("[WARN] deleting IaCService project fixture %s failed: %s", id, err)
			}
		}
	})
	return projectId, groupId, groupIdUpdated, emptyProjectId
}

// Test IaCService Module. <<< Resource test cases, automatically generated.
