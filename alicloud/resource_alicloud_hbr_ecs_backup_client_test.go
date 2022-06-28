package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudHBREcsBackupClient_basic0(t *testing.T) {
	checkoutAccount(t, true)
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	defer checkoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_ecs_backup_client.default"
	ra := resourceAttrInit(resourceId, AlicloudHBREcsBackupClientMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrEcsBackupClient")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrecsbackupclient%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBREcsBackupClientBasicDependence0)
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
					"instance_id": "${data.alicloud_instances.default.instances.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"use_https": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"use_https": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_network_type": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_network_type": "PUBLIC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_network_type": "VPC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_network_type": "VPC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_network_type": "CLASSIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_network_type": "CLASSIC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "STOPPED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "STOPPED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ACTIVATED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ACTIVATED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_cpu_core": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_cpu_core": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_worker": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_worker": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_proxy_setting": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_proxy_setting": "DISABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_proxy_setting": "USE_CONTROL_PROXY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_proxy_setting": "USE_CONTROL_PROXY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_proxy_setting": "CUSTOM",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_proxy_setting": "CUSTOM",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_host": "192.168.11.100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_host": "192.168.11.100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_port": "22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_port": "22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_user": "admin",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_user": "admin",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_password": "admin",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_password": "admin",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"use_https":          "false",
					"data_network_type":  "PUBLIC",
					"max_cpu_core":       "2",
					"max_worker":         "4",
					"data_proxy_setting": "USE_CONTROL_PROXY",
					"proxy_host":         "192.168.11.101",
					"proxy_port":         "80",
					"proxy_user":         "user",
					"proxy_password":     "password",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"use_https":          "false",
						"data_network_type":  "PUBLIC",
						"max_cpu_core":       "2",
						"max_worker":         "4",
						"data_proxy_setting": "USE_CONTROL_PROXY",
						"proxy_host":         "192.168.11.101",
						"proxy_port":         "80",
						"proxy_user":         "user",
						"proxy_password":     "password",
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

var AlicloudHBREcsBackupClientMap0 = map[string]string{
	"max_cpu_core":       CHECKSET,
	"proxy_port":         CHECKSET,
	"proxy_host":         "",
	"proxy_password":     "",
	"data_proxy_setting": "",
	"proxy_user":         "",
	"data_network_type":  CHECKSET,
	"status":             CHECKSET,
	"use_https":          CHECKSET,
	"max_worker":         CHECKSET,
	"instance_id":        CHECKSET,
}

func AlicloudHBREcsBackupClientBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-hbr-ecs-backup-plan"
  status     = "Running"
}

`, name)
}
