package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks ProjectMember. >>> Resource test cases, automatically generated.
// Case projectMember对接terraform_成都 8920
func TestAccAliCloudDataWorksProjectMember_basic8920(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_project_member.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksProjectMemberMap8920)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksProjectMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksProjectMemberBasicDependence8920)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_id":    "${alicloud_ram_user.defaultKCTrB2.id}",
					"project_id": "${alicloud_data_works_project.defaultCoMnk8.id}",
					"roles": []map[string]interface{}{
						{
							"code": "${var.管理员角色码}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":    CHECKSET,
						"project_id": CHECKSET,
						"roles.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"roles": []map[string]interface{}{
						{
							"code": "${var.访客角色码}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"roles.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"roles": []map[string]interface{}{
						{
							"code": "${var.管理员角色码}",
						},
						{
							"code": "${var.部署角色码}",
						},
						{
							"code": "${var.访客角色码}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"roles.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"roles": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"roles.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"roles": []map[string]interface{}{
						{
							"code": "${var.管理员角色码}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"roles.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"roles": []map[string]interface{}{
						{
							"code": "${var.访客角色码}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"roles.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"roles": []map[string]interface{}{
						{
							"code": "${var.访客角色码}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"roles.#": "1",
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

var AlicloudDataWorksProjectMemberMap8920 = map[string]string{
	"status": CHECKSET,
}

func AlicloudDataWorksProjectMemberBasicDependence8920(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "部署角色码" {
  default = "role_project_deploy"
}

variable "管理员角色码" {
  default = "role_project_admin"
}

variable "工作空间名称" {
  default = "test_member_pop"
}

variable "访客角色码" {
  default = "role_project_guest"
}

variable "deploy名称" {
  default = "部署"
}

variable "guest名称" {
  default = "访客"
}

variable "admin名称" {
  default = "空间管理员"
}

resource "alicloud_ram_user" "defaultKCTrB2" {
  display_name = var.name
  name = var.name
}

resource "alicloud_data_works_project" "defaultCoMnk8" {
  status       = "Available"
  project_name = format("%%s1", var.name)
  display_name = var.工作空间名称
  description  = var.工作空间名称
  pai_task_enabled = true
}


`, name)
}

// Test DataWorks ProjectMember. <<< Resource test cases, automatically generated.
