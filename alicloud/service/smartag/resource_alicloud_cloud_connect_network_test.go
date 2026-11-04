package smartag_test

import (
	"fmt"
	"testing"

	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/smartag"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudConnectNetwork_basic(t *testing.T) {
	var ccn smartagsdk.CloudConnectNetwork
	resourceId := "alicloud_cloud_connect_network.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &ccn, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetwork-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCcnDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			acctest.PreCheck(t)
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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

func TestAccAliCloudCloudConnectNetwork_multi(t *testing.T) {
	var ccn smartagsdk.CloudConnectNetwork
	resourceId := "alicloud_cloud_connect_network.default.9"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &ccn, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetwork-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCcnDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			acctest.PreCheck(t)
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
