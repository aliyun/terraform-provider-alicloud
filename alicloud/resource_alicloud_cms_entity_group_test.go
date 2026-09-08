package alicloud

import (
	"fmt"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_entity_group", &resource.Sweeper{
		Name: "alicloud_cms_entity_group",
		F:    testSweepCmsEntityGroup,
	})
}

func testSweepCmsEntityGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{"tf-testAcc", "tf_testAcc"}
	action := "/entity-groups"
	query := map[string]*string{
		"MaxResults": StringPointer(fmt.Sprintf("%d", PageSizeLarge)),
	}
	for {
		response, err := client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_entity_group", action, AlibabaCloudSdkGoERROR)
		}
		resp, _ := jsonpath.Get("$.entityGroups[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["entityGroupName"])
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
			entityGroupId := fmt.Sprint(item["entityGroupId"])
			ws := fmt.Sprint(item["workspace"])
			delAction := fmt.Sprintf("/entity-groups/%s", entityGroupId)
			delQuery := map[string]*string{}
			if ws != "" {
				delQuery["workspace"] = StringPointer(ws)
			}
			if _, err := client.RoaDelete("Cms", "2024-03-30", delAction, delQuery, nil, nil, true); err != nil {
				if !IsExpectedErrors(err, []string{"404", "500"}) && !NotFoundError(err) {
					return WrapErrorf(err, DefaultErrorMsg, entityGroupId, delAction, AlibabaCloudSdkGoERROR)
				}
			}
		}
		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}
	return nil
}

func testAccCheckAlicloudCmsEntityGroupExists(n string, v map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No entity group ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cmsServiceV2 := CmsServiceV2{client}
		object, err := cmsServiceV2.DescribeCmsEntityGroup(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("DescribeCmsEntityGroup failed: %v", err)
		}
		egRaw, _ := jsonpath.Get("$.entityGroup", object)
		if egRaw != nil {
			if eg, ok := egRaw.(map[string]interface{}); ok {
				for k, val := range eg {
					v[k] = val
				}
			}
		}
		return nil
	}
}

func testAccCheckAlicloudCmsEntityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_entity_group" {
			continue
		}
		_, err := cmsServiceV2.DescribeCmsEntityGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"404", "500"}) {
				continue
			}
			return fmt.Errorf("DescribeCmsEntityGroup error checking destroy: %v", err)
		}
		return fmt.Errorf("Entity group %s still exists", rs.Primary.ID)
	}
	return nil
}

func AliCloudCmsEntityGroupBasicDependence(name string) string {
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
`, name)
}

func TestAccAliCloudCmsEntityGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_entity_group.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	ra := resourceAttrInit(resourceId, AliCloudCmsEntityGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEntityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEntityGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudCmsEntityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"entity_group_name": name,
					"description":       "tf test entity group",
					"resource_group_id": "rg-initial456",
					"workspace":         "${alicloud_cms_workspace.default.id}",
					"entity_rules": []map[string]interface{}{
						{
							"resource_group_id": "rg-rules-initial",
							"entity_types":      []string{"ECS"},
							"instance_ids":      []string{"i-bp1test0001"},
							"tags": []map[string]interface{}{
								{
									"op":         "add",
									"tag_key":    "Env",
									"tag_values": []string{"test"},
								},
							},
							"labels": []map[string]interface{}{
								{
									"op":         "add",
									"tag_key":    "Team",
									"tag_values": []string{"tf"},
								},
							},
							"annotations": []map[string]interface{}{
								{
									"op":         "add",
									"tag_key":    "Source",
									"tag_values": []string{"terraform"},
								},
							},
							"ip_match_rule": []map[string]interface{}{
								{
									"ip_field_key": "ip",
									"ip_cidr":      "192.168.0.0/16",
								},
							},
							"field_rules": []map[string]interface{}{
								{
									"field_key":    "region",
									"op":           "eq",
									"field_values": []string{"cn-hangzhou"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudCmsEntityGroupExists(resourceId, v),
					testAccCheck(map[string]string{
						"entity_group_name": name,
						"description":       "tf test entity group",
						"workspace":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entity_group_name": name + "_update",
					"description":       "tf test entity group updated",
					"resource_group_id": "rg-updated789",
					"workspace":         "${alicloud_cms_workspace.default.id}",
					"entity_rules": []map[string]interface{}{
						{
							"resource_group_id": "rg-rules-updated",
							"entity_types":      []string{"CS"},
							"instance_ids":      []string{"i-bp1test0002"},
							"tags": []map[string]interface{}{
								{
									"op":         "remove",
									"tag_key":    "Team",
									"tag_values": []string{"updated"},
								},
							},
							"labels": []map[string]interface{}{
								{
									"op":         "remove",
									"tag_key":    "Project",
									"tag_values": []string{"updated"},
								},
							},
							"annotations": []map[string]interface{}{
								{
									"op":         "remove",
									"tag_key":    "Owner",
									"tag_values": []string{"updated"},
								},
							},
							"ip_match_rule": []map[string]interface{}{
								{
									"ip_field_key": "clientIp",
									"ip_cidr":      "10.0.0.0/8",
								},
							},
							"field_rules": []map[string]interface{}{
								{
									"field_key":    "zone",
									"op":           "ne",
									"field_values": []string{"cn-beijing"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudCmsEntityGroupExists(resourceId, v),
					testAccCheck(map[string]string{
						"entity_group_name": name + "_update",
						"description":       "tf test entity group updated",
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

func TestAccAliCloudCmsEntityGroup_emptyRules(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_entity_group.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	ra := resourceAttrInit(resourceId, AliCloudCmsEntityGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEntityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEntityGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudCmsEntityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"entity_group_name": name,
					"description":       "tf test minimal entity group",
					"workspace":         "${alicloud_cms_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudCmsEntityGroupExists(resourceId, v),
					testAccCheck(map[string]string{
						"entity_group_name": name,
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

var AliCloudCmsEntityGroupMap = map[string]string{
	"entity_group_id": CHECKSET,
}
