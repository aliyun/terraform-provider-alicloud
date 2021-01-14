package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBrainIndustrialPidLoop_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_brain_industrial_pid_loop.default"
	ra := resourceAttrInit(resourceId, AlicloudBrainIndustrialPidLoopMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Brain_industrialService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBrainIndustrialPidLoop")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBrainIndustrialPidLoopBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_dcs_type":      "standard",
					"pid_loop_desc":          "Test For Terraform",
					"pid_loop_is_crucial":    "false",
					"pid_loop_name":          name,
					"pid_loop_type":          "0",
					"pid_project_id":         "${alicloud_brain_industrial_pid_project.default.id}",
					"pid_loop_configuration": `{\"baseParam\":{\"forwardController\":false,\"integral\":false,\"kd\":{\"tagValue\":\"20\"},\"kp\":{},\"op\":\"PIDBenchmark.FOPDT_OP\",\"opParam\":{\"increment\":{\"max\":10},\"operate\":{\"max\":115,\"min\":-15},\"range\":{\"max\":115,\"min\":-15},\"trend\":0},\"openLoopTime\":150,\"pv\":\"PIDBenchmark.FOPDT_PV\",\"pvRange\":{\"max\":100,\"min\":0},\"sampleTime\":5,\"sp\":\"PIDBenchmark.FOPDT_SP\",\"spOperate\":{\"max\":100,\"min\":0},\"splitRangeControl\":false,\"suitCtrlTime\":100,\"td\":{},\"ti\":{}},\"identParam\":{\"delay\":10,\"modelType\":3},\"resetParam\":{\"ctrlMode\":0,\"ctrlStuc\":1}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_dcs_type":      "standard",
						"pid_loop_desc":          "Test For Terraform",
						"pid_loop_is_crucial":    "false",
						"pid_loop_name":          name,
						"pid_loop_type":          "0",
						"pid_project_id":         CHECKSET,
						"pid_loop_configuration": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pid_loop_configuration", "pid_loop_desc", "pid_project_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_desc": "Test For Terraform Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_project_id": CHECKSET,
						"pid_loop_desc":  "Test For Terraform Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_configuration": `{\"baseParam\":{\"forwardController\":false,\"integral\":false,\"kd\":{\"tagValue\":\"20\"},\"kp\":{},\"op\":\"PIDBenchmark.FOPDT_OP\",\"opParam\":{\"increment\":{\"max\":10},\"operate\":{\"max\":115,\"min\":-15},\"range\":{\"max\":115,\"min\":-15},\"trend\":0},\"openLoopTime\":150,\"pv\":\"PIDBenchmark.FOPDT_PV\",\"pvRange\":{\"max\":100,\"min\":0},\"sampleTime\":6,\"sp\":\"PIDBenchmark.FOPDT_SP\",\"spOperate\":{\"max\":100,\"min\":0},\"splitRangeControl\":false,\"suitCtrlTime\":100,\"td\":{},\"ti\":{}},\"identParam\":{\"delay\":10,\"modelType\":3},\"resetParam\":{\"ctrlMode\":0,\"ctrlStuc\":1}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_configuration": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_is_crucial": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_project_id":      CHECKSET,
						"pid_loop_is_crucial": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_desc":          "Test For Terraform",
					"pid_loop_name":          name,
					"pid_loop_type":          "0",
					"pid_loop_is_crucial":    "false",
					"pid_loop_configuration": `{\"baseParam\":{\"forwardController\":false,\"integral\":false,\"kd\":{\"tagValue\":\"20\"},\"kp\":{},\"op\":\"PIDBenchmark.FOPDT_OP\",\"opParam\":{\"increment\":{\"max\":10},\"operate\":{\"max\":115,\"min\":-15},\"range\":{\"max\":115,\"min\":-15},\"trend\":0},\"openLoopTime\":150,\"pv\":\"PIDBenchmark.FOPDT_PV\",\"pvRange\":{\"max\":100,\"min\":0},\"sampleTime\":5,\"sp\":\"PIDBenchmark.FOPDT_SP\",\"spOperate\":{\"max\":100,\"min\":0},\"splitRangeControl\":false,\"suitCtrlTime\":100,\"td\":{},\"ti\":{}},\"identParam\":{\"delay\":10,\"modelType\":3},\"resetParam\":{\"ctrlMode\":0,\"ctrlStuc\":1}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_desc":          "Test For Terraform",
						"pid_loop_name":          name,
						"pid_loop_type":          "0",
						"pid_loop_is_crucial":    "false",
						"pid_loop_configuration": CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudBrainIndustrialPidLoopMap = map[string]string{}

func AlicloudBrainIndustrialPidLoopBasicDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_brain_industrial_pid_organization" "default" {
		pid_organization_name = "%[1]s"
	}
	resource "alicloud_brain_industrial_pid_project" "default" {
		pid_organisation_id = alicloud_brain_industrial_pid_organization.default.id
		pid_project_name = "%[1]s"
	}`, name)
}
