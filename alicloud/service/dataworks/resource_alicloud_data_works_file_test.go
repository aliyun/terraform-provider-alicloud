package dataworks_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const dataWorksFileResourceId = "alicloud_data_works_file.default"

func TestAccAliCloudDataWorksFile_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testacc%sdataworksfile%d", region, rand)
	projectId := os.Getenv("DATA_WORKS_PROJECT_ID")
	if projectId == "" {
		projectId = "34051"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDataWorksFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataWorksFileBasicConfig(projectId, name, "SHOW TABLES;"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataWorksFileResourceId, "file_name", name),
					resource.TestCheckResourceAttr(dataWorksFileResourceId, "file_type", "10"),
					resource.TestCheckResourceAttr(dataWorksFileResourceId, "project_id", projectId),
					resource.TestCheckResourceAttr(dataWorksFileResourceId, "content", "SHOW TABLES;"),
					resource.TestCheckResourceAttrSet(dataWorksFileResourceId, "id"),
					resource.TestCheckResourceAttrSet(dataWorksFileResourceId, "file_id"),
				),
			},
			{
				// content is updatable via UpdateFile; file_type is ForceNew.
				Config: testAccDataWorksFileBasicConfig(projectId, name, "SELECT 1;"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataWorksFileResourceId, "content", "SELECT 1;"),
					resource.TestCheckResourceAttr(dataWorksFileResourceId, "file_name", name),
				),
			},
			{
				ResourceName:            dataWorksFileResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deploy", "submit", "project_identifier"},
			},
		},
	})
}

func testAccDataWorksFileBasicConfig(projectId, name, content string) string {
	return fmt.Sprintf(`
resource "alicloud_data_works_file" "default" {
  project_id = "%[1]s"
  file_name  = "%[2]s"
  file_type  = "10"
  content    = %[3]q
}
`, projectId, name, content)
}

func testAccCheckDataWorksFileDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := dataworks.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_data_works_file" {
			continue
		}
		_, err := service.DescribeDataWorksFile(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("checking alicloud_data_works_file destroy (%s): %w", rs.Primary.ID, err)
		}
		return fmt.Errorf("alicloud_data_works_file %s still exists", rs.Primary.ID)
	}
	return nil
}
