package alb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudALBAcl_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_acl.default"
	ra := acctest.ResourceAttrInit(resourceId, AclMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAcl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccalbacl%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AclBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_entries": []map[string]interface{}{
						{
							"description": name,
							"entry":       "10.0.0.0/24",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/23",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/25",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/26",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/27",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/28",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/29",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/10",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/11",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/12",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/13",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/14",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/15",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/16",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/17",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/18",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/19",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/10",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/11",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/12",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/13",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/14",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/15",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/16",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/17",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/18",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/19",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/10",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/11",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/12",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/13",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/14",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/15",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/16",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/17",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/18",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/19",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_entries.#": "37",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_entries": []map[string]interface{}{
						{
							"description": name,
							"entry":       "10.0.0.0/23",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/22",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_entries.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc7",
						"For":     "Tftestacc7",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc7",
						"tags.For":     "Tftestacc7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": "${var.name}",
					"tags": map[string]string{
						"Created": "tfTestAcc99",
						"For":     "Tftestacc99",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name":     name,
						"tags.%":       "2",
						"tags.Created": "tfTestAcc99",
						"tags.For":     "Tftestacc99",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAliCloudALBAcl_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_acl.default"
	ra := acctest.ResourceAttrInit(resourceId, AclMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAcl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccalbacl%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AclBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name":          "${var.name}",
					"dry_run":           "false",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name":          name,
						"dry_run":           "false",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AclMap0 = map[string]string{
	"tags.%":            acctest.NOSET,
	"dry_run":           acctest.NOSET,
	"resource_group_id": acctest.CHECKSET,
	"acl_entries.#":     acctest.CHECKSET,
	"status":            acctest.CHECKSET,
}

func AclBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}
