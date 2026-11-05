package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	ecs "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudECSPrefixList_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_prefix_list.default"
	ra := resourceAttrInit(resourceId, AlicloudECSPrefixListMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsPrefixList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testaccrecsprefixlist")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSPrefixListBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"max_entries": "3",
					"entry": []map[string]interface{}{
						{
							"description": name,
							"cidr":        "192.168.0.0/24",
						},
					},
					"description":      name,
					"prefix_list_name": name,
					"address_family":   "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_entries":      "3",
						"entry.#":          "1",
						"description":      name,
						"prefix_list_name": name,
						"address_family":   "IPv4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entry": []map[string]interface{}{
						{
							"description": name,
							"cidr":        "192.168.0.0/24",
						},
						{
							"description": name + "update",
							"cidr":        "10.0.0.0/8",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entry.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entry": []map[string]interface{}{
						{
							"description": name,
							"cidr":        "192.168.0.0/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entry.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prefix_list_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix_list_name": name + "update",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      name + "updateall",
					"prefix_list_name": name + "updateall",
					"entry": []map[string]interface{}{
						{
							"description": name,
							"cidr":        "192.168.0.0/24",
						},
						{
							"description": name + "update",
							"cidr":        "10.0.0.0/8",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      name + "updateall",
						"prefix_list_name": name + "updateall",
						"entry.#":          "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudECSPrefixListMap0 = map[string]string{}

func AlicloudECSPrefixListBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
