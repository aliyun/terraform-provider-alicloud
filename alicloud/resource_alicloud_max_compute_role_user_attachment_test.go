package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MaxCompute RoleUserAttachment. >>> Resource test cases, automatically generated.
// Case RoleUserAttachment_terraform测试(RamRole) 9962
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9962(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9962)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputeroleuserattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9962)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${var.role_name}",
					"user":         "${var.ram_role}",
					"project_name": "${var.project_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    CHECKSET,
						"user":         CHECKSET,
						"project_name": CHECKSET,
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

var AlicloudMaxComputeRoleUserAttachmentMap9962 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9962(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "aliyun_user" {
  default = "ALIYUN$openapiautomation@test.aliyunid.com"
}

variable "ram_user" {
  default = "RAM$openapiautomation@test.aliyunid.com:tf-example"
}

variable "ram_role" {
  default = "RAM$openapiautomation@test.aliyunid.com:role/terraform-no-ak-assumerole-no-deleting"
}

variable "role_name" {
  default = "role_project_admin"
}

variable "project_name" {
  default = "default_project_669886c"
}


`, name)
}

// Case RoleUserAttachment_terraform测试(RamUser) 9961
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9961(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9961)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputeroleuserattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9961)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${var.role_name}",
					"user":         "${var.ram_user}",
					"project_name": "${var.project_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    CHECKSET,
						"user":         CHECKSET,
						"project_name": CHECKSET,
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

var AlicloudMaxComputeRoleUserAttachmentMap9961 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9961(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "aliyun_user" {
  default = "ALIYUN$openapiautomation@test.aliyunid.com"
}

variable "ram_user" {
  default = "RAM$openapiautomation@test.aliyunid.com:tf-example"
}

variable "ram_role" {
  default = "RAM$openapiautomation@test.aliyunid.com:role/terraform-no-ak-assumerole-no-deleting"
}

variable "role_name" {
  default = "role_project_admin"
}

variable "project_name" {
  default = "default_project_669886c"
}


`, name)
}

// Case RoleUserAttachment_terraform测试(adminRole) 9966
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9966(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9966)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputeroleuserattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9966)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "admin",
					"user":         "${var.ram_user}",
					"project_name": "${var.project_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    "admin",
						"user":         CHECKSET,
						"project_name": CHECKSET,
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

var AlicloudMaxComputeRoleUserAttachmentMap9966 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9966(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "aliyun_user" {
  default = "ALIYUN$openapiautomation@test.aliyunid.com"
}

variable "ram_user" {
  default = "RAM$openapiautomation@test.aliyunid.com:tf-example"
}

variable "ram_role" {
  default = "RAM$openapiautomation@test.aliyunid.com:role/terraform-no-ak-assumerole-no-deleting"
}

variable "role_name" {
  default = "role_project_admin"
}

variable "project_name" {
  default = "default_project_669886c"
}


`, name)
}

// Case RoleUserAttachment_terraform测试(super_administratorRole) 9967
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9967(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9967)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputeroleuserattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9967)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${var.role_name}",
					"user":         "${var.ram_user}",
					"project_name": "${var.project_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    CHECKSET,
						"user":         CHECKSET,
						"project_name": CHECKSET,
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

var AlicloudMaxComputeRoleUserAttachmentMap9967 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9967(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "aliyun_user" {
  default = "ALIYUN$openapiautomation@test.aliyunid.com"
}

variable "ram_user" {
  default = "RAM$openapiautomation@test.aliyunid.com:tf-example"
}

variable "ram_role" {
  default = "RAM$openapiautomation@test.aliyunid.com:role/terraform-no-ak-assumerole-no-deleting"
}

variable "role_name" {
  default = "role_project_admin"
}

variable "project_name" {
  default = "default_project_669886c"
}


`, name)
}

// Case RoleUserAttachment_terraform测试(AliyunUser) 9960
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9960(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9960)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputeroleuserattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9960)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${var.role_name}",
					"user":         "${var.aliyun_user}",
					"project_name": "${var.project_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    CHECKSET,
						"user":         CHECKSET,
						"project_name": CHECKSET,
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

var AlicloudMaxComputeRoleUserAttachmentMap9960 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9960(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "aliyun_user" {
  default = "ALIYUN$openapiautomation@test.aliyunid.com"
}

variable "ram_user" {
  default = "RAM$openapiautomation@test.aliyunid.com:tf-example"
}

variable "ram_role" {
  default = "RAM$openapiautomation@test.aliyunid.com:role/terraform-no-ak-assumerole-no-deleting"
}

variable "role_name" {
  default = "role_project_admin"
}

variable "project_name" {
  default = "default_project_669886c"
}


`, name)
}

// Test MaxCompute RoleUserAttachment. <<< Resource test cases, automatically generated.
