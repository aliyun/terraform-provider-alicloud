package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms Dataset. >>> Resource test cases, hand-written.
func TestAccAliCloudCmsDataset_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dataset.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsDatasetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdsbasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsDatasetBasicDependence)
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
					"workspace":    "${alicloud_cms_workspace.default.workspace_name}",
					"dataset_name": name,
					"description":  name,
					"schema":       `{\"metric\":{\"type\":\"text\"}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":    CHECKSET,
						"dataset_name": name,
						"description":  name,
						"schema":       CHECKSET,
						"create_time":  CHECKSET,
						"region_id":    CHECKSET,
						"update_time":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("%s-updated", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("%s-updated", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "",
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

func TestAccAliCloudCmsDataset_basic_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dataset.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsDatasetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdstwin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsDatasetBasicDependence)
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
					"workspace":    "${alicloud_cms_workspace.default.workspace_name}",
					"dataset_name": name,
					"description":  name,
					"schema":       `{\"metric\":{\"type\":\"text\"}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":    CHECKSET,
						"dataset_name": name,
						"description":  name,
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

func TestAccAliCloudCmsDataset_disappear(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dataset.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsDatasetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdsdisp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsDatasetBasicDependence)
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
					"workspace":    "${alicloud_cms_workspace.default.workspace_name}",
					"dataset_name": name,
					"description":  name,
					"schema":       `{\"metric\":{\"type\":\"text\"}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":    CHECKSET,
						"dataset_name": name,
						"description":  name,
					}),
				),
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				ResourceName:       resourceId,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

var AliCloudCmsDatasetMap = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"update_time": CHECKSET,
}

func AliCloudCmsDatasetBasicDependence(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}
`, name)
}
