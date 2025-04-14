package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Eflo Resource. >>> Resource test cases, automatically generated.
// Case Resource资源用例_接入CCAPI_线上 10576
func TestAccAliCloudEfloResource_basic10576(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_resource.default"
	ra := resourceAttrInit(resourceId, AliCloudEfloResourceMap10576)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEfloResourceBasicDependence10576)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EfloExperimentPlanTemplateSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_param": []map[string]interface{}{
						{
							"access_id":    os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ACCESS_ID"),
							"access_key":   os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ACCESS_KEY"),
							"workspace_id": os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_WORKSPACE_ID"),
							"endpoint":     os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ENDPOINT"),
						},
					},
					"cluster_id": name,
					"machine_types": []map[string]interface{}{
						{
							"cpu_info": "2x Intel Saphhire Rapid 8469C 48C CPU",
							"gpu_info": "8x OAM 810 GPU",
						},
					},
					"cluster_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":          name,
						"cluster_name":        name,
						"user_access_param.#": "1",
						"machine_types.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_param": []map[string]interface{}{
						{
							"access_id":    os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ACCESS_ID_UPDATE"),
							"access_key":   os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ACCESS_KEY_UPDATE"),
							"workspace_id": os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_WORKSPACE_ID_UPDATE"),
							"endpoint":     os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ENDPOINT_UPDATE"),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_param.#": "1",
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

func TestAccAliCloudEfloResource_basic10576_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_resource.default"
	ra := resourceAttrInit(resourceId, AliCloudEfloResourceMap10576)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEfloResourceBasicDependence10576)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EfloExperimentPlanTemplateSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_param": []map[string]interface{}{
						{
							"access_id":    os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ACCESS_ID"),
							"access_key":   os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ACCESS_KEY"),
							"workspace_id": os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_WORKSPACE_ID"),
							"endpoint":     os.Getenv("EFLO_CNP_USER_ACCESS_PARAM_ENDPOINT"),
						},
					},
					"cluster_id": name,
					"machine_types": []map[string]interface{}{
						{
							"memory_info":  "32x 64GB DDR4 4800 Memory",
							"type":         "Private",
							"bond_num":     "5",
							"node_count":   "1",
							"cpu_info":     "2x Intel Saphhire Rapid 8469C 48C CPU",
							"network_info": "1x 200Gbps Dual Port BF3 DPU for VPC 4x 200Gbps Dual Port EIC",
							"gpu_info":     "8x OAM 810 GPU",
							"disk_info":    "2x 480GB SATA SSD 4x 3.84TB NVMe SSD",
							"network_mode": "net",
							"name":         "lingjun",
						},
					},
					"cluster_name": name,
					"cluster_desc": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":          name,
						"cluster_name":        name,
						"user_access_param.#": "1",
						"machine_types.#":     "1",
						"cluster_desc":        name,
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

var AliCloudEfloResourceMap10576 = map[string]string{
	"resource_id": CHECKSET,
}

func AliCloudEfloResourceBasicDependence10576(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Eflo Resource. <<< Resource test cases, automatically generated.
