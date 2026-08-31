// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ecs DiskEncryptionByDefault. >>> Resource test cases, automatically generated.
// Case 测试账号下云盘默认加密_正式 12733
func TestAccAliCloudEcsDiskEncryptionByDefault_basic12733(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_encryption_by_default.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskEncryptionByDefaultMap12733)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDiskEncryptionByDefault")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskEncryptionByDefaultBasicDependence12733)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypted": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypted": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypted": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypted": "false",
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

func TestAccAliCloudEcsDiskEncryptionByDefault_basic12733_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_encryption_by_default.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskEncryptionByDefaultMap12733)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDiskEncryptionByDefault")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskEncryptionByDefaultBasicDependence12733)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypted": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypted": "true",
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

func TestAccAliCloudEcsDiskEncryptionByDefault_basic12735_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_encryption_by_default.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskEncryptionByDefaultMap12733)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDiskEncryptionByDefault")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskEncryptionByDefaultBasicDependence12733)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypted": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypted": "false",
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

var AliCloudEcsDiskEncryptionByDefaultMap12733 = map[string]string{
	"encrypted": CHECKSET,
}

func AliCloudEcsDiskEncryptionByDefaultBasicDependence12733(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Ecs DiskEncryptionByDefault. <<< Resource test cases, automatically generated.
