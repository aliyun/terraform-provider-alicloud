package cloudmonitorservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Cms Dataset. >>> Resource test cases, hand-written.
func TestAccAliCloudCmsDataset_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dataset.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsDatasetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsDataset")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdsbasic%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsDatasetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
						"workspace":    acctest.CHECKSET,
						"dataset_name": name,
						"description":  name,
						"schema":       acctest.CHECKSET,
						"create_time":  acctest.CHECKSET,
						"region_id":    acctest.CHECKSET,
						"update_time":  acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsDatasetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsDataset")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdstwin%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsDatasetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
						"workspace":    acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsDatasetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsDataset")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdsdisp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsDatasetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
						"workspace":    acctest.CHECKSET,
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
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
	"update_time": acctest.CHECKSET,
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
