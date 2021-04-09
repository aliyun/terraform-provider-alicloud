package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAdbDbCluster_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbDbClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudAdbDbClusterBasicDependence0)
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
					"db_cluster_category": "Cluster",
					"db_cluster_class":    "C8",
					"db_node_count":       "2",
					"db_node_storage":     "200",
					"mode":                "reserver",
					"db_cluster_version":  `3.0`,
					"payment_type":        "PayAsYouGo",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_category": "Cluster",
						"db_cluster_class":    "C8",
						"db_node_count":       "2",
						"db_node_storage":     "200",
						"mode":                "reserver",
						"db_cluster_version":  "3.0",
						"payment_type":        "PayAsYouGo",
						"vswitch_id":          CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_cluster_class", "mode", "modify_type", "period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class":   "C32",
					"db_node_count":   "4",
					"db_node_storage": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":   "C32",
						"db_node_count":   "4",
						"db_node_storage": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Test description update.",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Test description update.",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "23:00Z-00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "23:00Z-00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
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
					"description": "Test description.",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "Test description.",
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
		},
	})
}

var AlicloudAdbDbClusterMap0 = map[string]string{
	"db_cluster_version":  "3.0",
	"elastic_io_resource": "0",
}

func AlicloudAdbDbClusterBasicDependence0(name string) string {
	return fmt.Sprintf(`
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "terraformci"
}
%s
`, AdbCommonTestCase)
}
