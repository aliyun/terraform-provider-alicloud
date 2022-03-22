package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCddcDedicatedHostGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHostGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scddcdedicatedhostgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostGroupBasicDependence0)
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
					"engine": "Redis",
					"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine": "Redis",
						"vpc_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_allocation_ratio": "110",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_allocation_ratio": "110",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_policy": "Intensively",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_policy": "Intensively",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_policy": "Evenly",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_policy": "Evenly",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_replace_policy": "Manual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_replace_policy": "Manual",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mem_allocation_ratio": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mem_allocation_ratio": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_allocation_ratio": "101",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_allocation_ratio": "101",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_allocation_ratio":     "111",
					"allocation_policy":         "Evenly",
					"host_replace_policy":       "Auto",
					"mem_allocation_ratio":      "61",
					"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
					"cpu_allocation_ratio":      "102",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_allocation_ratio":     "111",
						"allocation_policy":         "Evenly",
						"host_replace_policy":       "Auto",
						"mem_allocation_ratio":      "61",
						"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
						"cpu_allocation_ratio":      "102",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudCddcDedicatedHostGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHostGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scddcdedicatedhostgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostGroupBasicDependence0)
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
					"engine":          "SQLServer",
					"vpc_id":          "${data.alicloud_vpcs.default.ids.0}",
					"open_permission": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine": "SQLServer",
						"vpc_id": CHECKSET,
					}),
				),
			},
			// SQLServer does not support to update disk_allocation_ratio
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"disk_allocation_ratio": "110",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"disk_allocation_ratio": "110",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_policy": "Intensively",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_policy": "Intensively",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_replace_policy": "Manual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_replace_policy": "Manual",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mem_allocation_ratio": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mem_allocation_ratio": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_allocation_ratio": "101",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_allocation_ratio": "101",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_policy":         "Evenly",
					"host_replace_policy":       "Auto",
					"mem_allocation_ratio":      "61",
					"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
					"cpu_allocation_ratio":      "102",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_policy":         "Evenly",
						"host_replace_policy":       "Auto",
						"mem_allocation_ratio":      "61",
						"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
						"cpu_allocation_ratio":      "102",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudCddcDedicatedHostGroup_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHostGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scddcdedicatedhostgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostGroupBasicDependence0)
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
					"engine": "MySQL",
					"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine": "MySQL",
						"vpc_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_allocation_ratio": "110",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_allocation_ratio": "110",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_policy": "Intensively",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_policy": "Intensively",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_replace_policy": "Manual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_replace_policy": "Manual",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mem_allocation_ratio": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mem_allocation_ratio": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_allocation_ratio": "101",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_allocation_ratio": "101",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_allocation_ratio":     "111",
					"allocation_policy":         "Evenly",
					"host_replace_policy":       "Auto",
					"mem_allocation_ratio":      "61",
					"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
					"cpu_allocation_ratio":      "102",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_allocation_ratio":     "111",
						"allocation_policy":         "Evenly",
						"host_replace_policy":       "Auto",
						"mem_allocation_ratio":      "61",
						"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
						"cpu_allocation_ratio":      "102",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudCddcDedicatedHostGroup_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHostGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scddcdedicatedhostgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostGroupBasicDependence0)
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
					"engine": "PostgreSQL",
					"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine": "PostgreSQL",
						"vpc_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_allocation_ratio": "110",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_allocation_ratio": "110",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_policy": "Intensively",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_policy": "Intensively",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_replace_policy": "Manual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_replace_policy": "Manual",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mem_allocation_ratio": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mem_allocation_ratio": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_allocation_ratio": "101",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_allocation_ratio": "101",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_allocation_ratio":     "111",
					"allocation_policy":         "Evenly",
					"host_replace_policy":       "Auto",
					"mem_allocation_ratio":      "61",
					"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
					"cpu_allocation_ratio":      "102",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_allocation_ratio":     "111",
						"allocation_policy":         "Evenly",
						"host_replace_policy":       "Auto",
						"mem_allocation_ratio":      "61",
						"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
						"cpu_allocation_ratio":      "102",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudCddcDedicatedHostGroup_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHostGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scddcdedicatedhostgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostGroupBasicDependence0)
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
					"engine": "MongoDB",
					"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine": "MongoDB",
						"vpc_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_allocation_ratio": "110",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_allocation_ratio": "110",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_policy": "Intensively",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_policy": "Intensively",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_replace_policy": "Manual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_replace_policy": "Manual",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mem_allocation_ratio": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mem_allocation_ratio": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_desc": "DedicatedHostGroupDescAlone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_allocation_ratio": "101",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_allocation_ratio": "101",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_allocation_ratio":     "111",
					"allocation_policy":         "Evenly",
					"host_replace_policy":       "Auto",
					"mem_allocation_ratio":      "61",
					"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
					"cpu_allocation_ratio":      "102",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_allocation_ratio":     "111",
						"allocation_policy":         "Evenly",
						"host_replace_policy":       "Auto",
						"mem_allocation_ratio":      "61",
						"dedicated_host_group_desc": "DedicatedHostGroupDescAll",
						"cpu_allocation_ratio":      "102",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudCddcDedicatedHostGroup_basic5(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHostGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scddcdedicatedhostgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostGroupBasicDependence0)
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
					"engine":                    "SQLServer",
					"vpc_id":                    "${data.alicloud_vpcs.default.ids.0}",
					"allocation_policy":         "Evenly",
					"host_replace_policy":       "Auto",
					"dedicated_host_group_desc": name,
					"open_permission":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_policy":         "Evenly",
						"host_replace_policy":       "Auto",
						"dedicated_host_group_desc": name,
						"engine":                    "SQLServer",
						"vpc_id":                    CHECKSET,
						"open_permission":           "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCDDCDedicatedHostGroupMap0 = map[string]string{
	"engine": CHECKSET,
	"vpc_id": CHECKSET,
}

func AlicloudCDDCDedicatedHostGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}


`, name)
}
