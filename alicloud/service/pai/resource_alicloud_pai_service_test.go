package pai_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/pai"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPaiService_basic7678_modified(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPaiServiceMap7678_modified)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return pai.NewPaiServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePaiService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc7678%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPaiServiceBasicDependence7678_modified)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					//"labels": map[string]interface{}{
					//	"\"0\"": "{\\\"LabelKey\\\":\\\"testkey\\\",\\\"LabelValue\\\":\\\"testvalue\\\"}",
					//},
					"develop":        "false",
					"service_config": "{\\\"metadata\\\":{\\\"cpu\\\":1,\\\"gpu\\\":0,\\\"instance\\\":1,\\\"memory\\\":2000,\\\"name\\\":\\\"${var.name}\\\",\\\"rpc\\\":{\\\"keepalive\\\":70000},\\\"workspace_id\\\":\\\"${alicloud_pai_workspace_workspace.default.id}\\\"},\\\"workspace_id\\\":\\\"${alicloud_pai_workspace_workspace.default.id}\\\",\\\"model_path\\\":\\\"http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz\\\",\\\"processor_entry\\\":\\\"libecho.so\\\",\\\"processor_path\\\":\\\"http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz\\\",\\\"processor_type\\\":\\\"cpp\\\"}",
					"workspace_id":   "${alicloud_pai_workspace_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"develop":        "false",
						"service_config": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_config": "{\\\"metadata\\\":{\\\"cpu\\\":1,\\\"gpu\\\":0,\\\"instance\\\":1,\\\"memory\\\":2000,\\\"name\\\":\\\"${var.name}\\\",\\\"rpc\\\":{\\\"keepalive\\\":70000},\\\"workspace_id\\\":\\\"${alicloud_pai_workspace_workspace.update.id}\\\"},\\\"workspace_id\\\":\\\"${alicloud_pai_workspace_workspace.update.id}\\\",\\\"model_path\\\":\\\"http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz\\\",\\\"processor_entry\\\":\\\"libecho.so\\\",\\\"processor_path\\\":\\\"http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz\\\",\\\"processor_type\\\":\\\"cpp\\\"}",
					"workspace_id":   "${alicloud_pai_workspace_workspace.update.id}",
					//"labels": map[string]interface{}{
					//	"\"0\"": "{\\\"LabelKey\\\":\\\"testkeyupdate\\\",\\\"LabelValue\\\":\\\"testvalueupdate\\\"}",
					//	"\"1\"": "{\\\"LabelKey\\\":\\\"testkeyupdate1\\\",\\\"LabelValue\\\":\\\"testvalueupdate1\\\"}",
					//},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_config": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_config": "{\\\"metadata\\\":{\\\"cpu\\\":1,\\\"gpu\\\":0,\\\"instance\\\":1,\\\"memory\\\":2000,\\\"name\\\":\\\"${var.name}\\\",\\\"rpc\\\":{\\\"keepalive\\\":70000},\\\"workspace_id\\\":\\\"${alicloud_pai_workspace_workspace.update.id}\\\"},\\\"workspace_id\\\":\\\"${alicloud_pai_workspace_workspace.update.id}\\\",\\\"model_path\\\":\\\"http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz\\\",\\\"processor_entry\\\":\\\"libecho.so\\\",\\\"processor_path\\\":\\\"http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz\\\",\\\"processor_type\\\":\\\"cpp\\\"}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"develop"},
			},
		},
	})
}

var AlicloudPaiServiceMap7678_modified = map[string]string{
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
	"service_uid": acctest.CHECKSET,
}

func AlicloudPaiServiceBasicDependence7678_modified(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

variable "space_name" {
  default = "%[1]s"
}

variable "space_name_update" {
  default = "%[1]s_update"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_pai_workspace_workspace" "default" {
  description    = var.space_name
  workspace_name = var.space_name
  display_name   = var.space_name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_workspace" "update" {
  description    = var.space_name_update
  workspace_name = var.space_name_update
  display_name   = var.space_name_update
  env_types      = ["prod"]
}


`, name)
}

// Case Service 3213
func TestAccAliCloudPaiService_basic3213(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPaiServiceMap3213)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return pai.NewPaiServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePaiService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc%spaiservice%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPaiServiceBasicDependence3213)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"develop":        "false",
					"service_config": fmt.Sprintf("{\\\"metadata\\\":{\\\"cpu\\\":1,\\\"gpu\\\":0,\\\"instance\\\":1,\\\"memory\\\":2000,\\\"name\\\":\\\"tftestacc234%d\\\",\\\"rpc\\\":{\\\"keepalive\\\":70000}},\\\"model_path\\\":\\\"http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz\\\",\\\"processor_entry\\\":\\\"libecho.so\\\",\\\"processor_path\\\":\\\"http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz\\\",\\\"processor_type\\\":\\\"cpp\\\"}", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"develop":        "false",
						"service_config": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"develop"},
			},
		},
	})
}

var AlicloudPaiServiceMap3213 = map[string]string{
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
	"service_uid": acctest.CHECKSET,
}

func AlicloudPaiServiceBasicDependence3213(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
