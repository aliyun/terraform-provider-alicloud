package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens Disk. >>> Resource test cases, automatically generated.
// Case 5178
func TestAccAliCloudEnsDisk_basic5178(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_disk.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsDiskMap5178)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensdisk%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsDiskBasicDependence5178)
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
					"category":      "cloud_ssd",
					"payment_type":  "PayAsYouGo",
					"ens_region_id": "cn-chongqing-11",
					"size":          "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":      "cloud_ssd",
						"payment_type":  "PayAsYouGo",
						"ens_region_id": "cn-chongqing-11",
						"size":          "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category":      "cloud_ssd",
					"payment_type":  "PayAsYouGo",
					"ens_region_id": "cn-chongqing-11",
					"size":          "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":      "cloud_ssd",
						"payment_type":  "PayAsYouGo",
						"ens_region_id": "cn-chongqing-11",
						"size":          "60",
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

var AlicloudEnsDiskMap5178 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEnsDiskBasicDependence5178(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5179
func TestAccAliCloudEnsDisk_basic5179(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_disk.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsDiskMap5179)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-ensdisk%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsDiskBasicDependence5179)
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
					"size":          "20",
					"category":      "cloud_efficiency",
					"payment_type":  "PayAsYouGo",
					"ens_region_id": "ch-zurich-1",
					"disk_name":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size":          "20",
						"category":      "cloud_efficiency",
						"payment_type":  "PayAsYouGo",
						"ens_region_id": "ch-zurich-1",
						"disk_name":     name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_id":   "${alicloud_ens_snapshot.createsnapshot.id}",
					"category":      "cloud_efficiency",
					"kms_key_id":    "${alicloud_ens_disk.createdisk.kms_key_id}",
					"size":          "20",
					"encrypted":     "true",
					"payment_type":  "PayAsYouGo",
					"ens_region_id": "ch-zurich-1",
					"disk_name":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_id":   CHECKSET,
						"category":      "cloud_efficiency",
						"kms_key_id":    CHECKSET,
						"size":          "20",
						"encrypted":     "true",
						"payment_type":  "PayAsYouGo",
						"ens_region_id": "ch-zurich-1",
						"disk_name":     name,
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

var AlicloudEnsDiskMap5179 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEnsDiskBasicDependence5179(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_disk" "createdisk" {
  category      = "cloud_efficiency"
  size          = "20"
  encrypted     = true
  payment_type  = "PayAsYouGo"
  ens_region_id = "ch-zurich-1"
}

resource "alicloud_ens_snapshot" "createsnapshot" {
  description   = "DiskDescription_autotest"
  ens_region_id = "ch-zurich-1"
  snapshot_name = var.name

  disk_id = alicloud_ens_disk.createdisk.id
}


`, name)
}

// Case 5178  twin
func TestAccAliCloudEnsDisk_basic5178_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_disk.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsDiskMap5178)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensdisk%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsDiskBasicDependence5178)
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
					"category":      "cloud_ssd",
					"size":          "20",
					"payment_type":  "PayAsYouGo",
					"ens_region_id": "cn-chongqing-18",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":      "cloud_ssd",
						"size":          "20",
						"payment_type":  "PayAsYouGo",
						"ens_region_id": "cn-chongqing-18",
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

// Case 5179  twin
func TestAccAliCloudEnsDisk_basic5179_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_disk.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsDiskMap5179)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensdisk%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsDiskBasicDependence5179)
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
					"snapshot_id":   "${alicloud_ens_snapshot.createsnapshot.id}",
					"category":      "cloud_efficiency",
					"kms_key_id":    "${alicloud_ens_disk.createdisk.kms_key_id}",
					"size":          "30",
					"encrypted":     "true",
					"payment_type":  "PayAsYouGo",
					"ens_region_id": "ch-zurich-1",
					"disk_name":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_id":   CHECKSET,
						"category":      "cloud_efficiency",
						"kms_key_id":    CHECKSET,
						"size":          "30",
						"encrypted":     "true",
						"payment_type":  "PayAsYouGo",
						"ens_region_id": "ch-zurich-1",
						"disk_name":     name,
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

// Test Ens Disk. <<< Resource test cases, automatically generated.
