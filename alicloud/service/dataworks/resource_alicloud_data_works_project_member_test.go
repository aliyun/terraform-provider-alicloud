package dataworks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The following TypeList attribute is pre-existing stock surfaced by the
// dataworks subpackage migration (schema unchanged by this PR). A convergent
// reorder ACC requires a DataWorks-enabled account, which is not opened this
// round, so no fabricated reorder configuration is added; it is marked as
// known stock debt per scripts/collection-order/README.md.
// collection-order: bypass roles -- pre-existing stock TypeList; schema unchanged by migration; reorder ACC needs a DataWorks account (not opened this round)

// Test DataWorks ProjectMember. >>> Resource test cases, automatically generated.
// Case projectMember对接terraform_成都 8920
func TestAccAliCloudDataWorksProjectMember_basic8920(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_project_member.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDataWorksProjectMemberMap8920)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dataworks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDataWorksProjectMember")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDataWorksProjectMemberBasicDependence8920)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"user_id":    acctest.CHECKSET,
						"project_id": acctest.CHECKSET,
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
					"roles": acctest.REMOVEKEY,
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
	"status": acctest.CHECKSET,
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
