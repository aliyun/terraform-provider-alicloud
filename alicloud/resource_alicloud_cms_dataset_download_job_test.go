package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms DatasetDownloadJob. >>> Resource test cases, hand-written.
func TestAccAliCloudCmsDatasetDownloadJob_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dataset_download_job.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsDatasetDownloadJobMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDatasetDownloadJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccddlbasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsDatasetDownloadJobBasicDependence)
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
					"workspace":        "${alicloud_cms_workspace.default.workspace_name}",
					"dataset_name":     "${alicloud_cms_dataset.default.dataset_name}",
					"job_name":         name,
					"query":            "* | select count(*) as total",
					"compression_type": "gzip",
					"content_type":     "csv",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":        CHECKSET,
						"dataset_name":     CHECKSET,
						"job_name":         name,
						"query":            "* | select count(*) as total",
						"compression_type": "gzip",
						"content_type":     "csv",
						"create_time":      CHECKSET,
						"region_id":        CHECKSET,
						"update_time":      CHECKSET,
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

func TestAccAliCloudCmsDatasetDownloadJob_basic_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dataset_download_job.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsDatasetDownloadJobMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDatasetDownloadJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccddltwin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsDatasetDownloadJobBasicDependence)
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
					"workspace":        "${alicloud_cms_workspace.default.workspace_name}",
					"dataset_name":     "${alicloud_cms_dataset.default.dataset_name}",
					"job_name":         name,
					"query":            "* | select count(*) as total",
					"compression_type": "gzip",
					"content_type":     "csv",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":        CHECKSET,
						"dataset_name":     CHECKSET,
						"job_name":         name,
						"query":            "* | select count(*) as total",
						"compression_type": "gzip",
						"content_type":     "csv",
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

func TestAccAliCloudCmsDatasetDownloadJob_disappear(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dataset_download_job.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsDatasetDownloadJobMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDatasetDownloadJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccddldisp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsDatasetDownloadJobBasicDependence)
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
					"workspace":        "${alicloud_cms_workspace.default.workspace_name}",
					"dataset_name":     "${alicloud_cms_dataset.default.dataset_name}",
					"job_name":         name,
					"query":            "* | select count(*) as total",
					"compression_type": "gzip",
					"content_type":     "csv",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":        CHECKSET,
						"dataset_name":     CHECKSET,
						"job_name":         name,
						"query":            "* | select count(*) as total",
						"compression_type": "gzip",
						"content_type":     "csv",
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

var AliCloudCmsDatasetDownloadJobMap = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"update_time": CHECKSET,
}

func AliCloudCmsDatasetDownloadJobBasicDependence(name string) string {
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

resource "alicloud_cms_dataset" "default" {
  workspace    = alicloud_cms_workspace.default.workspace_name
  dataset_name = var.name
  description  = "terraform-example"
  schema       = "{\"metric\":{\"type\":\"text\"}}"
}
`, name)
}
