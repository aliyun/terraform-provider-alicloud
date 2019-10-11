package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"testing"
)

func TestAccAlicloudPolarDBCluster(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	var ips []map[string]interface{}
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   CHECKSET,
		"db_node_class": CHECKSET,
		"vswitch_id":    CHECKSET,
		"db_type":       CHECKSET,
		"db_version":    CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterConfigDependence)

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
					"db_type":       "MySQL",
					"db_version":    "8.0",
					"pay_type":      "PostPaid",
					"db_node_class": "polar.mysql.x4.large",
					"vswitch_id":    "${alicloud_vswitch.default.id}",
					"description":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testaccdnsrecordbasic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccdnsrecordbasic",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "16:00Z-17:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "16:00Z-17:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMapsForPolarDB(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   "tf-testaccdnsrecordbasic1",
					"maintain_time": "02:00Z-03:00Z",
					"security_ips":  []string{"10.168.1.13", "100.69.7.113"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "tf-testaccdnsrecordbasic1",
						"maintain_time": "02:00Z-03:00Z",
					}),
					testAccCheckKeyValueInMapsForPolarDB(ips, "security ip", "security_ips", "10.168.1.13,100.69.7.113"),
				),
			},
		},
	})

}

func TestAccAlicloudPolarDBClusterMulti(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "alicloud_polardb_cluster.default.2"
	var basicMap = map[string]string{
		"description":   CHECKSET,
		"db_node_class": CHECKSET,
		"vswitch_id":    CHECKSET,
		"db_type":       CHECKSET,
		"db_version":    CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterConfigDependence)

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
					"count":         "3",
					"db_type":       "MySQL",
					"db_version":    "8.0",
					"pay_type":      "PostPaid",
					"db_node_class": "polar.mysql.x4.large",
					"vswitch_id":    "${alicloud_vswitch.default.id}",
					"description":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckKeyValueInMapsForPolarDB(ps []map[string]interface{}, propName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range ps {
			if policy[key].(string) != value {
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

func resourcePolarDBClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

`, PolarDBCommonTestCase, name)
}
