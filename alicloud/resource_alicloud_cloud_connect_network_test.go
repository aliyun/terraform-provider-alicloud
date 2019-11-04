package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCloudConnectNetwork_basic(t *testing.T) {
	var ccn smartag.CloudConnectNetwork
	resourceId := "alicloud_cloud_connect_network.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &ccn, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetwork-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCcnDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"is_default": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_default": "true",
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
					"name": fmt.Sprintf("%s-Name", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("%s-Name", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("%s-Description", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("%s-Description", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_block": "192.168.1.0/24,192.168.2.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block": "192.168.1.0/24,192.168.2.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": name,
					"cidr_block":  "192.168.0.0/24",
					"is_default":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": name,
						"cidr_block":  "192.168.0.0/24",
						"is_default":  "true",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCloudConnectNetwork_multi(t *testing.T) {
	var ccn smartag.CloudConnectNetwork
	resourceId := "alicloud_cloud_connect_network.default.9"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &ccn, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetwork-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCcnDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}-${count.index}",
					"count":       "10",
					"description": "${var.name}-${count.index}",
					"cidr_block":  "192.168.0.0/24",
					"is_default":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("%s-9", name),
						"description": fmt.Sprintf("%s-9", name),
						"cidr_block":  "192.168.0.0/24",
						"is_default":  "true",
					}),
				),
			},
		},
	})
}

func resourceCcnDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
