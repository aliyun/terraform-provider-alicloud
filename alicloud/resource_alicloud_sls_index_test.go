// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls Index. >>> Resource test cases, automatically generated.
// Case index_terraform 10982
func TestAccAliCloudSlsIndex_basic10982(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsIndex")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-nanjing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            "true",
							"case_sensitive": "true",
							"token": []string{
								"a"},
							"exclude_keys": []string{
								"t"},
						},
					},
					"logstore_name": "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":  "${var.project_name}",
					"keys":          "{\\\"test\\\":{\\\"caseSensitive\\\":false,\\\"token\\\":[\\\"\\\\n\\\",\\\"\\\\t\\\",\\\",\\\",\\\" \\\",\\\";\\\",\\\"\\\\\\\"\\\",\\\"'\\\",\\\"(\\\",\\\")\\\",\\\"{\\\",\\\"}\\\",\\\"[\\\",\\\"]\\\",\\\"<\\\",\\\">\\\",\\\"?\\\",\\\"/\\\",\\\"#\\\",\\\":\\\"],\\\"type\\\":\\\"text\\\",\\\"doc_value\\\":false,\\\"alias\\\":\\\"\\\",\\\"chn\\\":false}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": CHECKSET,
						"project_name":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            "false",
							"case_sensitive": "false",
							"token": []string{
								"tt"},
							"include_keys": []string{
								"tt"},
						},
					},
					"max_text_len": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_text_len": "500",
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

var AlicloudSlsIndexMap10982 = map[string]string{}

func AlicloudSlsIndexBasicDependence10982(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "logstore_name" {
  default = "logstore-test"
}

variable "project_name" {
  default = "project-for-index-terraform-test"
}

resource "alicloud_log_project" "defaultdCM1bA" {
  description = "terrafrom test"
  name        = var.project_name
}

resource "alicloud_log_store" "default7MW26R" {
  hot_ttl          = "7"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaultdCM1bA.project_name
  name             = var.logstore_name
}


`, name)
}

// Test Sls Index. <<< Resource test cases, automatically generated.
