package sls_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Sls MachineGroup. >>> Resource test cases, automatically generated.
// Case machine_group_terrafrom 10981
func TestAccAliCloudSlsMachineGroup_basic10981(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_machine_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsMachineGroupMap10981)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsMachineGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsMachineGroupBasicDependence10981)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":            "group1",
					"group_type":            "test",
					"project_name":          "${alicloud_log_project.defaultyJqrue.project_name}",
					"machine_identify_type": "ip",
					"group_attribute": []map[string]interface{}{
						{
							"group_topic":   "test",
							"external_name": "test",
						},
					},
					"machine_list": []string{
						"192.168.1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":            "group1",
						"group_type":            "test",
						"project_name":          acctest.CHECKSET,
						"machine_identify_type": "ip",
						"machine_list.#":        "1",
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

var AlicloudSlsMachineGroupMap10981 = map[string]string{}

func AlicloudSlsMachineGroupBasicDependence10981(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "project_name" {
  default = "project-for-machine-group-terraform"
}

resource "alicloud_log_project" "defaultyJqrue" {
  description = "for terraform test"
  name        = var.name
}


`, name)
}

// Case machine_group_unordered_list 10982
func TestAccAliCloudSlsMachineGroup_basic10982(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_machine_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsMachineGroupMap10982)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsMachineGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsMachineGroupBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":            "group_unordered",
					"group_type":            "test",
					"project_name":          "${alicloud_log_project.defaultyJqrue.project_name}",
					"machine_identify_type": "ip",
					"group_attribute": []map[string]interface{}{
						{
							"group_topic":   "test",
							"external_name": "test",
						},
					},
					"machine_list": []string{
						"192.168.1.1",
						"192.168.1.2",
						"192.168.1.3",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":            "group_unordered",
						"group_type":            "test",
						"project_name":          acctest.CHECKSET,
						"machine_identify_type": "ip",
						"machine_list.#":        "3",
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

var AlicloudSlsMachineGroupMap10982 = map[string]string{}

func AlicloudSlsMachineGroupBasicDependence10982(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "project_name" {
  default = "project-for-machine-group-terraform-unordered"
}

resource "alicloud_log_project" "defaultyJqrue" {
  description = "for terraform test"
  name        = var.name
}


`, name)
}

// Test Sls MachineGroup. <<< Resource test cases, automatically generated.
