// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test ThreatDetection ServiceLinkedRole. >>> Resource test cases, automatically generated.
// Case resource_ServiceLinkedRole_test_1 12848
func TestAccAliCloudThreatDetectionServiceLinkedRole_basic12848(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_service_linked_role.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionServiceLinkedRoleMap12848)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionServiceLinkedRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionServiceLinkedRoleBasicDependence12848)
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
					"service_linked_role": "AliyunServiceRoleForAntiRansomwareMssp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_linked_role": "AliyunServiceRoleForAntiRansomwareMssp",
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

var AlicloudThreatDetectionServiceLinkedRoleMap12848 = map[string]string{
	"role_status": CHECKSET,
}

func AlicloudThreatDetectionServiceLinkedRoleBasicDependence12848(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_ServiceLinkedRole_test_2 12849
func TestAccAliCloudThreatDetectionServiceLinkedRole_basic12849(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_service_linked_role.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionServiceLinkedRoleMap12849)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionServiceLinkedRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionServiceLinkedRoleBasicDependence12849)
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
					"service_linked_role": "AliyunServiceRoleForSas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_linked_role": "AliyunServiceRoleForSas",
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

var AlicloudThreatDetectionServiceLinkedRoleMap12849 = map[string]string{
	"role_status": CHECKSET,
}

func AlicloudThreatDetectionServiceLinkedRoleBasicDependence12849(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_ServiceLinkedRole_test 12850
func TestAccAliCloudThreatDetectionServiceLinkedRole_basic12850(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_service_linked_role.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionServiceLinkedRoleMap12850)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionServiceLinkedRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionServiceLinkedRoleBasicDependence12850)
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
					"service_linked_role": "AliyunServiceRoleForSasSecllm",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_linked_role": "AliyunServiceRoleForSasSecllm",
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

var AlicloudThreatDetectionServiceLinkedRoleMap12850 = map[string]string{
	"role_status": CHECKSET,
}

func AlicloudThreatDetectionServiceLinkedRoleBasicDependence12850(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ThreatDetection ServiceLinkedRole. <<< Resource test cases, automatically generated.
