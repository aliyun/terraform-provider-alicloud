package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	ecs "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The white-box TestUnitAlicloudECSCommand (which calls unexported
// resourceEcsCommand Create/Read/Delete and the helper.String alias) has been
// split into resource_alicloud_ecs_command_unit_test.go (package ecs); only the
// external acceptance tests remain here.

func TestAccAliCloudECSCommand_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_command.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsCommandMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsCommand")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEcsCommand%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsCommandBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"command_content": "bHMK",
					"description":     "For Terraform Test",
					"name":            name,
					"type":            "RunShellScript",
					"working_dir":     "/root",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_content": "bHMK",
						"description":     "For Terraform Test",
						"name":            name,
						"type":            "RunShellScript",
						"working_dir":     "/root",
					}),
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

func TestAccAliCloudECSCommand_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_command.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsCommandMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsCommand")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEcsCommand%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsCommandBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"command_content": "bHMK",
					"description":     "For Terraform Test",
					"name":            name,
					"type":            "RunShellScript",
					"working_dir":     "/root",
					"timeout":         "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_content": "bHMK",
						"description":     "For Terraform Test",
						"name":            name,
						"type":            "RunShellScript",
						"working_dir":     "/root",
						"timeout":         "60",
					}),
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

var AlicloudEcsCommandMap = map[string]string{
	"enable_parameter": "false",
}

func AlicloudEcsCommandBasicDependence(name string) string {
	return ""
}
