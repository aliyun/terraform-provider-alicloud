// Package alicloud
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudModelstudioWorkspace_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_modelstudio_workspace.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-acc-ms-%d", rand)
	updatedName := fmt.Sprintf("tf-acc-ms-upd-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string { return "" })

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckModelstudioWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_name": name,
					"service_site":   "global",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelstudioWorkspaceExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "workspace_name", name),
					resource.TestCheckResourceAttr(resourceId, "service_site", "global"),
					resource.TestCheckResourceAttrSet(resourceId, "workspace_id"),
					resource.TestCheckResourceAttrSet(resourceId, "api_host"),
					resource.TestCheckResourceAttrSet(resourceId, "region_id"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_name": updatedName,
					"service_site":   "global",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModelstudioWorkspaceExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "workspace_name", updatedName),
					resource.TestCheckResourceAttr(resourceId, "service_site", "global"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckModelstudioWorkspaceExists(n string, v *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found in state", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for resource %s", n)
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		modelstudioService := ModelstudioService{client}
		obj, err := modelstudioService.DescribeModelstudioWorkspace(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("DescribeModelstudioWorkspace failed: %v", err)
		}
		*v = obj
		return nil
	}
}

func testAccCheckModelstudioWorkspaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	modelstudioService := ModelstudioService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_modelstudio_workspace" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}
		_, err := modelstudioService.DescribeModelstudioWorkspace(rs.Primary.ID)
		if err != nil && NotFoundError(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("unexpected error checking destroy for modelstudio workspace %s: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("modelstudio workspace %s still exists", rs.Primary.ID)
	}
	return nil
}
