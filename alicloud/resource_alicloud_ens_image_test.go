package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens Image. >>> Resource test cases, automatically generated.
// Case 创建镜像 5662
func TestAccAliCloudEnsImage_basic5662(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_image.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsImageMap5662)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensimage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsImageBasicDependence5662)
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
					"image_name":  name,
					"instance_id": "${alicloud_ens_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":  name,
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_after_image_upload"},
			},
		},
	})
}

var AliCloudEnsImageMap5662 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"target_oss_region_id": CHECKSET,
}

func AliCloudEnsImageBasicDependence5662(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_instance" "default" {
  system_disk {
    size = "20"
  }
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "PayAsYouGo"
  password                   = "12345678ABCabc"
  amount                     = "1"
  internet_max_bandwidth_out = "10"
  public_ip_identification   = true
  ens_region_id              = "cn-chenzhou-telecom_unicom_cmcc"
  period_unit                = "Month"
  instance_type              = "ens.sn1.stiny"
  status                     = "Stopped"
}


`, name)
}

// Case 创建镜像_释放标识 5694
func TestAccAliCloudEnsImage_basic5694(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_image.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsImageMap5694)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensimage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsImageBasicDependence5694)
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
					"image_name":  name,
					"instance_id": "${alicloud_ens_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":  name,
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_after_image_upload"},
			},
		},
	})
}

var AliCloudEnsImageMap5694 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"target_oss_region_id": CHECKSET,
}

func AliCloudEnsImageBasicDependence5694(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_instance" "default" {
  system_disk {
    size = "20"
  }
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "PayAsYouGo"
  password                   = "12345678ABCabc"
  amount                     = "1"
  internet_max_bandwidth_out = "10"
  public_ip_identification   = true
  ens_region_id              = "cn-chenzhou-telecom_unicom_cmcc"
  period_unit                = "Month"
  instance_type              = "ens.sn1.stiny"
  status                     = "Stopped"
}


`, name)
}

// Case 创建镜像 5662  twin
func TestAccAliCloudEnsImage_basic5662_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_image.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsImageMap5662)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensimage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsImageBasicDependence5662)
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
					"image_name":                name,
					"instance_id":               "${alicloud_ens_instance.default.id}",
					"delete_after_image_upload": "true",
					"target_oss_region_id":      "ap-southeast-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":                name,
						"instance_id":               CHECKSET,
						"delete_after_image_upload": "true",
						"target_oss_region_id":      "ap-southeast-1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_after_image_upload"},
			},
		},
	})
}

// Case 创建镜像_释放标识 5694  twin
func TestAccAliCloudEnsImage_basic5694_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_image.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsImageMap5694)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensimage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsImageBasicDependence5694)
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
					"image_name":                name,
					"instance_id":               "${alicloud_ens_instance.default.id}",
					"delete_after_image_upload": "false",
					"target_oss_region_id":      "ap-southeast-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":                name,
						"instance_id":               CHECKSET,
						"delete_after_image_upload": "false",
						"target_oss_region_id":      "ap-southeast-1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_after_image_upload"},
			},
		},
	})
}

// Test Ens Image. <<< Resource test cases, automatically generated.
