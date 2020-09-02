package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCSKubernetesNodePool_basic(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                 name,
					"cluster_id":           "c9c9866ff0a914d8ba77a6569bfb36a7f",
					"vswitch_ids":          []string{"vsw-j6c8b2s5p4jq22mlij816"},
					"instance_types":       []string{"ecs.g6.large"},
					"node_count":           "0",
					"password":             "Test12345",
					"system_disk_category": "cloud_ssd",
					"system_disk_size":     "20",
					"tags":                 []map[string]string{{"key": "test-tag-key", "value": "test-tag-value"}},
					"data_disks":           []map[string]string{{"size": "100", "category": "cloud_ssd"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"node_count":            "0",
						"password":              "Test12345",
						"system_disk_category":  "cloud_ssd",
						"system_disk_size":      "20",
						"tags.#":                "1",
						"tags.0.key":            "test-tag-key",
						"tags.0.value":          "test-tag-value",
						"data_disks.#":          "1",
						"data_disks.0.size":     "100",
						"data_disks.0.category": "cloud_ssd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "80",
					"tags":             []map[string]string{{"key": "test-tag-key", "value": "test-tag-value"}, {"key": "new-tag-key", "value": "new-tag-value"}},
					"data_disks":       []map[string]string{{"size": "100", "category": "cloud_ssd"}, {"size": "200", "category": "cloud_efficiency"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size":      "80",
						"tags.#":                "2",
						"tags.0.key":            "test-tag-key",
						"tags.0.value":          "test-tag-value",
						"tags.1.key":            "new-tag-key",
						"tags.1.value":          "new-tag-value",
						"data_disks.#":          "2",
						"data_disks.0.size":     "100",
						"data_disks.0.category": "cloud_ssd",
						"data_disks.1.size":     "200",
						"data_disks.1.category": "cloud_efficiency",
					}),
				),
			},
		},
	})
}

var csdKubernetesNodePoolBasicMap = map[string]string{
	"node_count":           "0",
	"instance_types.0":     CHECKSET,
	"system_disk_size":     "40",
	"system_disk_category": "cloud_efficiency",
}

func resourceCSNodePoolConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}
