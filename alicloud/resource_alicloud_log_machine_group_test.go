package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudLogMachineGroup_basic(t *testing.T) {
	var v *sls.MachineGroup
	resourceId := "alicloud_log_machine_group.default"
	ra := resourceAttrInit(resourceId, logMachineGroupMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogmachinegroupip-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogMachineGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          name,
					"project":       "${alicloud_log_project.default.name}",
					"identify_list": []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":    name,
						"project": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"identify_type": "userdefined",
					"identify_list": []string{"terraform", "abc1234"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"identify_type":   "userdefined",
						"identify_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"topic": "terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic": "terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"identify_type": REMOVEKEY,
					"identify_list": []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
					"topic":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"identify_type":   "ip",
						"identify_list.#": "3",
						"topic":           REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudLogMachineGroup_multi(t *testing.T) {
	var v *sls.MachineGroup
	resourceId := "alicloud_log_machine_group.default.4"
	ra := resourceAttrInit(resourceId, logMachineGroupMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogmachinegroupip-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogMachineGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          name + "${count.index}",
					"project":       "${alicloud_log_project.default.name}",
					"identify_list": []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
					"count":         "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceLogMachineGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "default" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	`, name)
}

var logMachineGroupMap = map[string]string{
	"name":            CHECKSET,
	"project":         CHECKSET,
	"identify_list.#": "3",
}
