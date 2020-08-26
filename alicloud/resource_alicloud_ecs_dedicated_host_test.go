package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEcsDedicatedHost_basic(t *testing.T) {
	var v ecs.DedicatedHost
	resourceId := "alicloud_ecs_dedicated_host.default"
	ra := resourceAttrInit(resourceId, EcsDedicatedHostMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDedicatedHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccEcsDedicatedHost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, EcsDedicatedHostBasicdependence)
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
					"dedicated_host_type": "ddh.g5",
					"description":         "From_Terraform",
					"dedicated_host_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_type": "ddh.g5",
						"description":         "From_Terraform",
						"dedicated_host_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "detail_fee", "dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "DDH_Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "DDH_Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_name": name + "ddh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_name": name + "ddh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_attributes": []map[string]interface{}{
						{
							"udp_timeout":     "70",
							"slb_udp_timeout": "70",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_attributes.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "Terraform",
						"For":     "DDH",
					},
					"dedicated_host_name": name,
					"description":         "From_Terraform",
					"network_attributes": []map[string]interface{}{
						{
							"udp_timeout":     "60",
							"slb_udp_timeout": "60",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":               "2",
						"tags.Created":         "Terraform",
						"tags.For":             "DDH",
						"dedicated_host_name":  name,
						"description":          "From_Terraform",
						"network_attributes.#": "1",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
		},
	})
}

var EcsDedicatedHostMap = map[string]string{
	"detail_fee": "false",
	"dry_run":    "false",
	"status":     CHECKSET,
}

func EcsDedicatedHostBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default"{
	}
`)
}
