package dataworks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudDataWorksResourceFile_basic covers the Framework resource
// alicloud_dataworks_resource_file (dataworks-public/2020-05-18): create an SQL
// resource file inside a DataWorks workspace, update its content/description,
// import it, and verify it is gone after destroy. The resource file is created
// under a freshly created workspace in cn-chengdu, where the new DataWorks
// workspace API is available (see TestAccAliCloudDataWorksProject_basic8903).
func TestAccAliCloudDataWorksResourceFile_basic(t *testing.T) {
	resourceId := "alicloud_dataworks_resource_file.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDataWorksResourceFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudDataWorksResourceFileBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "file_id"),
					resource.TestCheckResourceAttr(resourceId, "file_name", name),
					resource.TestCheckResourceAttr(resourceId, "file_type", "10"),
					resource.TestCheckResourceAttr(resourceId, "content", "SHOW TABLES;"),
					resource.TestCheckResourceAttr(resourceId, "file_description", "tf resource file"),
					resource.TestCheckResourceAttr(resourceId, "file_folder_path", "Business Flow/tfTestAcc/folderDi"),
					resource.TestCheckResourceAttrSet(resourceId, "region_id"),
				),
			},
			{
				Config: testAccAliCloudDataWorksResourceFileUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "file_name", fmt.Sprintf("%s-updated", name)),
					resource.TestCheckResourceAttr(resourceId, "content", "SELECT 1;"),
					resource.TestCheckResourceAttr(resourceId, "file_description", "tf resource file updated"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"file_folder_path"},
			},
		},
	})
}

func testAccCheckDataWorksResourceFileDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dataworks_resource_file" {
			continue
		}
		_, err := dataworks.NewService(client).DescribeDataWorksResourceFile(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing DataWorks resource file %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("DataWorks resource file %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudDataWorksResourceFileBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "default" {
  project_name      = var.name
  display_name      = var.name
  pai_task_enabled  = false
  description       = "tf workspace for resource file"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  status            = "Available"
}

resource "alicloud_dataworks_resource_file" "default" {
  project_id      = alicloud_data_works_project.default.id
  file_name       = var.name
  file_type       = 10
  content         = "SHOW TABLES;"
  file_description = "tf resource file"
  file_folder_path = "Business Flow/tfTestAcc/folderDi"
}
`, name)
}

func testAccAliCloudDataWorksResourceFileUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "default" {
  project_name      = var.name
  display_name      = var.name
  pai_task_enabled  = false
  description       = "tf workspace for resource file"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  status            = "Available"
}

resource "alicloud_dataworks_resource_file" "default" {
  project_id      = alicloud_data_works_project.default.id
  file_name       = "${var.name}-updated"
  file_type       = 10
  content         = "SELECT 1;"
  file_description = "tf resource file updated"
  file_folder_path = "Business Flow/tfTestAcc/folderDi"
}
`, name)
}
