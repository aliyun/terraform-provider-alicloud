package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBAccountUpdate(t *testing.T) {
	var v *rds.DBInstanceAccount
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbaccount-%d", rand)
	var basicMap = map[string]string{
		"instance_id": CHECKSET,
		"name":        "tftestnormal",
		"password":    "YourPassword_123",
		"type":        string(DBAccountNormal),
	}
	resourceId := "alicloud_db_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_db_instance.instance.id}",
					"name":        "tftestnormal",
					"password":    "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "from terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YourPassword_1234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "YourPassword_1234",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf test",
					"password":    "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf test",
						"password":    "YourPassword_123",
					}),
				),
			},
		},
	})
}

func resourceDBAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "name" {
		default = "%v"
	}

	data "alicloud_db_instance_engines" "default" {
  		instance_charge_type = "PostPaid"
  		engine               = "MySQL"
  		engine_version       = "5.6"
	}

	data "alicloud_db_instance_classes" "default" {
 	 	engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
	}

	resource "alicloud_db_instance" "instance" {
		engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
		instance_type = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	    instance_name = "${var.name}"
	}
	`, RdsCommonTestCase, name)
}
