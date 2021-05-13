package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDdoscooPort_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_port.default"
	ra := resourceAttrInit(resourceId, AlicloudDdoscooPortMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdoscooPort")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sddoscooport%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdoscooPortBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"frontend_port":     `7001`,
					"backend_port":      `7002`,
					"instance_id":       "${data.alicloud_ddoscoo_instances.default.ids.0}",
					"frontend_protocol": "tcp",
					"real_servers":      []string{"192.168.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"frontend_port":     `7001`,
						"backend_port":      `7002`,
						"frontend_protocol": "tcp",
						"real_servers.#":    "1",
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
					"real_servers": []string{"192.168.0.1", "192.168.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{"192.168.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "1",
					}),
				),
			},
		},
	})
}

var AlicloudDdoscooPortMap0 = map[string]string{
	"backend_port":      "7002",
	"frontend_port":     "7001",
	"instance_id":       CHECKSET,
	"frontend_protocol": "tcp",
	"real_servers.#":    "1",
}

func AlicloudDdoscooPortBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

data "alicloud_ddoscoo_instances" "default" {}
`, name)
}

func TestAccAlicloudDdoscooPort_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_port.default"
	ra := resourceAttrInit(resourceId, AlicloudDdoscooPortMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdoscooPort")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sddoscooport%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdoscooPortBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"frontend_port":     `7001`,
					"backend_port":      `7002`,
					"instance_id":       "${data.alicloud_ddoscoo_instances.default.ids.0}",
					"frontend_protocol": "udp",
					"real_servers":      []string{"192.168.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"frontend_port":     `7001`,
						"backend_port":      `7002`,
						"frontend_protocol": "udp",
						"real_servers.#":    "1",
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
					"real_servers": []string{"192.168.0.1", "192.168.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{"192.168.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "1",
					}),
				),
			},
		},
	})
}

var AlicloudDdoscooPortMap1 = map[string]string{
	"backend_port":      "7002",
	"frontend_port":     "7001",
	"instance_id":       CHECKSET,
	"frontend_protocol": "udp",
	"real_servers.#":    "1",
}

func AlicloudDdoscooPortBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

data "alicloud_ddoscoo_instances" "default" {}
`, name)
}
