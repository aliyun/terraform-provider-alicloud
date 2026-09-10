package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Eflo Vpd. >>> Resource test cases, automatically generated.
// Case VPD测试 11252
func TestAccAliCloudEfloVpd_basic11252(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_vpd.default"
	ra := resourceAttrInit(resourceId, AliCloudEfloVpdMap11252)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloVpd")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEfloVpdBasicDependence11252)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr":     "10.0.0.0/8",
					"vpd_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr":     CHECKSET,
						"vpd_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpd_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpd_name": name + "_update",
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
					"secondary_cidr_blocks": []string{
						"192.168.0.0/18", "192.168.64.0/18", "192.168.128.0/17"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{
						"172.16.0.0/12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{
						"192.168.0.0/18", "192.168.64.0/18", "192.168.128.0/17"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudEfloVpd_basic11252_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_vpd.default"
	ra := resourceAttrInit(resourceId, AliCloudEfloVpdMap11252)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloVpd")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEfloVpdBasicDependence11252)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr":              "10.0.0.0/8",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"vpd_name":          name,
					"secondary_cidr_blocks": []string{
						"192.168.0.0/18", "192.168.64.0/18", "192.168.128.0/17"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr":                    CHECKSET,
						"resource_group_id":       CHECKSET,
						"vpd_name":                name,
						"secondary_cidr_blocks.#": "3",
						"tags.%":                  "2",
						"tags.Created":            "TF",
						"tags.For":                "Test",
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

var AliCloudEfloVpdMap11252 = map[string]string{
	"status":            CHECKSET,
	"create_time":       CHECKSET,
	"gmt_modified":      CHECKSET,
	"region_id":         CHECKSET,
	"resource_group_id": CHECKSET,
}

func AliCloudEfloVpdBasicDependence11252(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}


`, name)
}

// Test Eflo Vpd. <<< Resource test cases, automatically generated.
