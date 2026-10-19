// Package alicloud
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var AliCloudCmsContextStoreMap = map[string]string{
	"context_store_name": CHECKSET,
	"context_type":       CHECKSET,
	"workspace":          CHECKSET,
	"description":        CHECKSET,
	"status":             CHECKSET,
	"region_id":          CHECKSET,
	"create_time":        CHECKSET,
	"update_time":        CHECKSET,
}

func AliCloudCmsContextStoreBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

	provider "alicloud" {
		alias = "hz"
		region = "cn-hangzhou"
	}
	resource "alicloud_log_project" "default" {
		provider = alicloud.hz
		project_name = var.name
	}
	resource "alicloud_log_store" "default" {
		provider = alicloud.hz
		project_name = alicloud_log_project.default.project_name
		logstore_name = "${var.name}-default"
	}
	resource "alicloud_log_project" "update" {
		provider = alicloud.hz
		project_name = "${var.name}-update"
	}
	resource "alicloud_log_store" "update" {
		provider = alicloud.hz
		project_name = alicloud_log_project.update.project_name
		logstore_name = "${var.name}-update"
	}
	resource "alicloud_cms_workspace" "default" {
		provider = alicloud.hz
		sls_project = alicloud_log_project.default.project_name
		workspace_name = var.name
	}
`, name)
}

func TestAccAliCloudCmsContextStore_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_context_store.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsContextStoreMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsContextStore")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsContextStoreBasicDependence)
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
					"provider":           "alicloud.hz",
					"workspace":          "${alicloud_cms_workspace.default.workspace_name}",
					"context_store_name": "${var.name}",
					"context_type":       "private",
					"description":        "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"context_store_name": CHECKSET,
						"context_type":       "private",
						"workspace":          CHECKSET,
						"description":        CHECKSET,
						"status":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":           "alicloud.hz",
					"workspace":          "${alicloud_cms_workspace.default.workspace_name}",
					"context_store_name": "${var.name}",
					"context_type":       "public",
					"description":        "${var.name}-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"context_type": "public",
						"description":  name + "-update",
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

func TestAccAliCloudCmsContextStore_config(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_context_store.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsContextStoreMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsContextStore")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsContextStoreBasicDependence)
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
					"provider":           "alicloud.hz",
					"workspace":          "${alicloud_cms_workspace.default.workspace_name}",
					"context_store_name": "${var.name}",
					"context_type":       "private",
					"description":        "${var.name}",
					"config": []interface{}{map[string]interface{}{
						"source": []interface{}{map[string]interface{}{
							"project":    "${alicloud_log_project.default.project_name}",
							"logstore":   "${alicloud_log_store.default.logstore_name}",
							"start_time": "2024-01-01T00:00:00Z",
						}},
						"metadata_field": map[string]interface{}{
							"key1": "value1",
							"key2": "value2",
						},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"context_store_name": CHECKSET,
						"workspace":          CHECKSET,
						"config.#":           "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":           "alicloud.hz",
					"workspace":          "${alicloud_cms_workspace.default.workspace_name}",
					"context_store_name": "${var.name}",
					"context_type":       "private",
					"description":        "${var.name}-update",
					"config": []interface{}{map[string]interface{}{
						"source": []interface{}{map[string]interface{}{
							"project":    "${alicloud_log_project.update.project_name}",
							"logstore":   "${alicloud_log_store.update.logstore_name}",
							"start_time": "2024-02-01T00:00:00Z",
						}},
						"metadata_field": map[string]interface{}{
							"key1": "value1-updated",
							"key2": "value2",
						},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"context_store_name": CHECKSET,
						"config.#":           "1",
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
