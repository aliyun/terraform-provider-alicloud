package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCcnInstance_basic(t *testing.T) {
	var ccn smartag.CloudConnectNetwork
	resourceId := "alicloud_ccn_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &ccn, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccCcnConfigName")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCcnBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": "tf-testAccCcnConfigDescription",
					"cidr_block":  "192.168.0.0/24",
					"is_default":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "tf-testAccCcnConfigDescription",
						"cidr_block":  "192.168.0.0/24",
						"is_default":  "true",
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
					"name": "tf-testAccCcnConfigName-Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccCcnConfigName-Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAccCcnConfigDescription-Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccCcnConfigDescription-Update",
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
					"cen_id":      "${alicloud_cen_instance.default.id}",
					"total_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":      CHECKSET,
						"total_count": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":      "${alicloud_cen_instance.default.id}",
					"total_count": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":      CHECKSET,
						"total_count": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": "tf-testAccCcnConfigDescription",
					"cidr_block":  "192.168.0.0/24",
					"is_default":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "tf-testAccCcnConfigDescription",
						"cidr_block":  "192.168.0.0/24",
						"is_default":  "true",
					}),
				),
			},
		},
	})
}

func resourceCcnBasicDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cen_instance" "default" {
 name = "tf-testAccCenConfigName"
 description = "tf-testAccCenConfigDescription"
}
	`)
}
